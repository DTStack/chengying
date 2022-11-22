// Licensed to Apache Software Foundation(ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation(ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package impl

import (
	"bytes"
	"database/sql"
	dbhelper "dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"dtstack.com/dtstack/easymatrix/matrix/cache"
	"dtstack.com/dtstack/easymatrix/matrix/encrypt"
	"dtstack.com/dtstack/easymatrix/matrix/util"
	"errors"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"dtstack.com/dtstack/easymatrix/go-common/api-base"
	"dtstack.com/dtstack/easymatrix/matrix/base"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"github.com/dchest/captcha"
	"github.com/kataras/iris/context"
	pwd "github.com/sethvargo/go-password/password"
)

const (
	ENABLE                       = 0
	DISABLE                      = 1
	LOGIN_LOCK_DISABLE           = 0
	LOGIN_LOCK_ENABLE            = 1
	FORCE_RESET_PASSWORD_DISABLE = 0
	FORCE_RESET_PASSWORD_ENABLE  = 1
	PASSWORD_NOT_CHANGED         = 0
	PASSWORD_CHANGED             = 1
)

var (
	VerifyCode     bool
	VerifyIdentity bool
)

func Login(ctx context.Context) apibase.Result {
	log.Debugf("[Login] Login from EasyMatrix API ")

	paramErrs := apibase.NewApiParameterErrors()
	username := ctx.FormValue("username")
	if username == "" {
		paramErrs.AppendError("$", fmt.Errorf("username is empty"))
	}

	password := ctx.FormValue("password")

	if password == "" {
		paramErrs.AppendError("$", fmt.Errorf("password is empty"))
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	b, err := encrypt.PlatformEncrypt.CommonDecrypt([]byte(password))
	if err != nil {
		return err
	}
	err, info := model.UserList.GetInfoByUserName(username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Errorf("%v", err)
		return err
	}

	//禁用状态
	if info.Status == DISABLE {
		return fmt.Errorf("用户已经被禁用，请联系管理员")
	}

	//锁定状态
	lockSwitch := cache.SysConfig.PlatFormSecurity.AccountLoginLockSwitch
	maxErrLogin := cache.SysConfig.PlatFormSecurity.AccountLoginLimitError
	cache.LoginErrorCount.Lock.RLock()
	defer cache.LoginErrorCount.Lock.RUnlock()
	if errCount, expiration, found := cache.LoginLockStatus.GetWithExpiration(username); found {
		if lockSwitch == LOGIN_LOCK_ENABLE && errCount.(int) >= maxErrLogin {
			now := time.Now()
			if minuteSpan := int(expiration.Sub(now).Minutes()); minuteSpan >= 1 {
				return fmt.Errorf("账户暂被锁定，请%v分钟后重试", minuteSpan)
			} else {
				if secondSpan := int(expiration.Sub(now).Seconds()); secondSpan == 0 {
					cache.LoginErrorCount.Data[username] = 0
				} else {
					return fmt.Errorf("账户暂被锁定，请%v秒钟后重试", secondSpan)
				}
			}
		} else if lockSwitch == LOGIN_LOCK_DISABLE {
			cache.LoginErrorCount.Data[username] = 0
			cache.LoginLockStatus.Flush()
		}
	} else {
		cache.LoginErrorCount.Data[username] = 0
	}

	//登录错误
	if errors.Is(err, sql.ErrNoRows) || info.PassWord != util.Md5(string(b)) {
		if lockSwitch == LOGIN_LOCK_ENABLE {
			cache.LoginErrorCount.Data[username] = cache.LoginErrorCount.Data[username] + 1
			cache.SetLoginCache(username, cache.LoginErrorCount.Data[username])
			if cache.LoginErrorCount.Data[username] == maxErrLogin {
				return fmt.Errorf("登录失败次数超过限制，账户暂被锁定")
			}
			return fmt.Errorf("用户名或密码错误，剩余%v次尝试", maxErrLogin-cache.LoginErrorCount.Data[username])
		}
		return fmt.Errorf("用户名或密码错误")
	}
	cache.LoginErrorCount.Data[username] = 0

	err, userRole := model.UserRoleList.ListUserRoleByUserId(info.ID)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	userPermission := 0
	for _, uRole := range userRole {
		err, info := model.RoleList.GetInfoByRoleId(uRole.RoleId)
		if err != nil {
			log.Errorf("Database err: %v", err)
			return err
		}
		roleValue, _ := strconv.Atoi(info.RoleValue)
		userPermission |= roleValue
	}
	dtToken := apibase.CreateToken(username, info.ID, userPermission)

	isAdmin := "false"
	for _, v := range userRole {
		if v.RoleId == model.ROLE_ADMIN_ID {
			isAdmin = "true"
			break
		}
	}

	apibase.SetCookie(ctx, "em_token", dtToken)
	apibase.SetCookie(ctx, "em_admin", isAdmin)
	if err := model.SafetyAuditList.InsertSafetyAuditRecord(username, "产品访问", "进入EM", ctx.RemoteAddr(), ""); err != nil {
		log.Errorf("failed to add safety audit record\n")
	}
	return true
}

func GetPublicKey(ctx context.Context) apibase.Result {
	log.Debugf("[GetPublicKey] Get Public Key from EasyMatrix API ")
	kind := cache.SysConfig.PlatFormSecurity.LoginEncrypt
	response := map[string]string{
		"encrypt_type":       kind,
		"encrypt_public_key": encrypt.PlatformEncrypt.CommonGetPublicKey(),
	}
	return response
}

func LogOut(ctx context.Context) apibase.Result {
	log.Debugf("[LogOut] LogOut from EasyMatrix API ")

	userId, err := ctx.Values().GetInt("userId")
	if err != nil {
		log.Errorf("%v", err)
		return err
	}

	if err := addSafetyAuditRecord(ctx, "产品访问", "退出EM", ""); err != nil {
		log.Errorf("failed to add safety audit record\n")
	}

	apibase.DeleteUserCache(userId)
	ctx.RemoveCookie("em_token")

	return true
}

func UserInfo(ctx context.Context) apibase.Result {

	username := sqlTransfer(ctx.URLParam("username"))
	pagination := apibase.GetPaginationFromQueryParameters(nil, ctx, model.UserInfo{})
	status := ctx.URLParam("status")
	roleIds := ctx.URLParam("role_id")
	info, count := model.UserList.UserInfoList(username, status, roleIds, pagination)

	for i, v := range info {
		info[i].UpdateTimeFormat = v.CreateTime.Time.Format(base.TsLayout)

		userId := v.ID
		resClusterList := make([]model.ResClusterInfo, 0)
		clusterList, err := model.DeployClusterList.GetClusterListByUserId(userId)
		if err != nil {
			return fmt.Errorf("Database err: %v", err)
		}
		for _, cluster := range clusterList {
			resClusterList = append(resClusterList, model.ResClusterInfo{
				ClusterId:   cluster.Id,
				ClusterName: cluster.Name,
				ClusterType: cluster.Type},
			)
		}
		info[i].ClusterList = resClusterList
	}
	return map[string]interface{}{
		"list":  info,
		"count": count,
	}
}

func Register(ctx context.Context) (rlt apibase.Result) {
	log.Debugf("[AddUser] AddUser from EasyMatrix API ")

	paramErrs := apibase.NewApiParameterErrors()
	email := ctx.FormValue("email")
	if email == "" {
		paramErrs.AppendError("$", fmt.Errorf("email is empty"))
	}
	userName := ctx.FormValue("username")
	if userName == "" {
		paramErrs.AppendError("$", fmt.Errorf("userName is empty"))
	}
	password := ctx.FormValue("password")

	if password == "" {
		var err error
		gen, _ := pwd.NewGenerator(&pwd.GeneratorInput{
			Symbols: "~!@#$%^&*",
		})
		password, err = gen.Generate(8, 2, 2, false, true)
		if err != nil {
			log.Errorf(err.Error())
			return fmt.Errorf("%v", err.Error())
		}
	}

	paramErrs.CheckAndThrowApiParameterErrors()

	company := ctx.FormValue("company")
	fullName := ctx.FormValue("fullName")
	phone := ctx.FormValue("phone")
	roleIds := ctx.FormValue("roleId")
	roleId := model.ROLE_READER_ID
	for _, v := range strings.Split(roleIds, ",") {
		value, _ := strconv.Atoi(strings.TrimSpace(v))
		if value < roleId {
			roleId = value
		}
	}
	clusterList := ctx.FormValue("clusterList")

	err, userId := model.UserList.InsertUserIfNotExist(userName, util.Md5(password), company, fullName, email, phone)
	if err != nil {
		log.Errorf(err.Error())
		return fmt.Errorf("%v", err.Error())
	}
	defer func() {
		if err := addSafetyAuditRecord(ctx, "用户管理", "创建账号", "账号："+userName); err != nil {
			log.Errorf("failed to add safety audit record\n")
		}
	}()

	err, _ = model.UserRoleList.InsertUserRole(userId, roleId)
	if err != nil {
		log.Errorf(err.Error())
		return fmt.Errorf("can not insert UserRole , err : %v", err.Error())
	}
	for _, clusterId := range strings.Split(clusterList, ",") {
		clusterId, _ := strconv.Atoi(clusterId)
		err, _ = model.ClusterRightList.InsertUserClusterRight(userId, clusterId)
		if err != nil {
			log.Errorf(err.Error())
			return fmt.Errorf("can not insert ClusterRight, err : %v", err.Error())
		}
	}

	return map[string]interface{}{
		"username": userName,
		"password": password,
	}
}

func ResetPwdByAdmin(ctx context.Context) apibase.Result {
	log.Debugf("[AddUser] AddUser from EasyMatrix API ")

	if err := isAdmin(ctx); err != nil {
		return err
	}

	paramErrs := apibase.NewApiParameterErrors()
	targetUserIdStr := ctx.FormValue("targetUserId")
	if targetUserIdStr == "" {
		paramErrs.AppendError("$", fmt.Errorf("targetUserId is empty"))
	}
	targetUserId, err := strconv.Atoi(targetUserIdStr)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	password := ctx.FormValue("password")
	if password == "" {
		paramErrs.AppendError("$", fmt.Errorf("password is empty"))
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	b, err := encrypt.PlatformEncrypt.CommonDecrypt([]byte(password))
	if err != nil {
		return err
	}

	//校验新密码是否与旧密码相同
	err, info := model.UserList.GetInfoByUserId(targetUserId)
	if err != nil {
		log.Errorf("GetInfoByUserId %v", err)
		return err
	}
	if info.PassWord == util.Md5(string(b)) {
		return fmt.Errorf("新密码不能与旧密码相同")
	}

	defer func() {
		err, info := model.UserList.GetInfoByUserId(targetUserId)
		if err != nil {
			log.Errorf("%v", err)
			return
		}
		if err := addSafetyAuditRecord(ctx, "用户管理", "重置密码", "账号："+info.UserName); err != nil {
			log.Errorf("failed to add safety audit record\n")
		}
	}()

	res := model.UserList.UpdatePwdByUserId(util.Md5(string(b)), targetUserId)
	return res
}

func ModifyInfoById(ctx context.Context) apibase.Result {
	log.Debugf("[ModifyUserById] ModifyUserById from EasyMatrix API ")

	userId, err := ctx.Values().GetInt("userId")
	if err != nil {
		log.Errorf("%v", err)
		return err
	}

	paramErrs := apibase.NewApiParameterErrors()
	email := ctx.FormValue("email")
	if email == "" {
		paramErrs.AppendError("$", fmt.Errorf("email is empty"))
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	company := ctx.FormValue("company")
	fullName := ctx.FormValue("fullName")
	phone := ctx.FormValue("phone")

	res := model.UserList.UpdateInfoByUserId(company, fullName, email, phone, userId)

	return res
}

func ModifyInfoByAdmin(ctx context.Context) apibase.Result {
	log.Debugf("[ModifyInfoByAdmin] ModifyInfoByAdmin from EasyMatrix API ")

	userId, err := strconv.Atoi(ctx.FormValue("userId"))
	if err != nil {
		log.Errorf("%v", err)
		return fmt.Errorf("param userId must be a number")
	}

	paramErrs := apibase.NewApiParameterErrors()
	email := ctx.FormValue("email")
	if email == "" {
		paramErrs.AppendError("$", fmt.Errorf("email is empty"))
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	company := ctx.FormValue("company")
	fullName := ctx.FormValue("fullName")
	phone := ctx.FormValue("phone")
	roleId := ctx.FormValue("roleId")
	clusterList := ctx.FormValue("clusterList")
	if err := model.UserRoleList.UpdateWhere(dbhelper.MakeWhereCause().Equal("user_id", userId), dbhelper.UpdateFields{"role_id": roleId}, false); err != nil {
		log.Errorf("failed to update role info error: %v", err)
		return err
	}

	//修改集群权限
	qClusterList, err := model.ClusterRightList.GetUserClusterRightByUserId(userId)
	if err != nil {
		log.Errorf(err.Error())
		return fmt.Errorf("can not insert ClusterRight, err : %v", err.Error())
	}
	qClusterMap := map[int]int{}
	for _, qCluster := range qClusterList {
		qClusterMap[qCluster.ClusterId] = qCluster.ID
	}
	clusterMap := map[int]interface{}{}
	for _, clusterId := range strings.Split(clusterList, ",") {
		clusterId, _ := strconv.Atoi(clusterId)
		clusterMap[clusterId] = ""
		//增加权限
		if _, ok := qClusterMap[clusterId]; !ok {
			err, _ := model.ClusterRightList.InsertUserClusterRight(userId, clusterId)
			if err != nil {
				log.Errorf(err.Error())
				return fmt.Errorf("can not insert ClusterRight, err : %v", err.Error())
			}
		}
	}
	//删除权限
	for qClusterId, id := range qClusterMap {
		if _, ok := clusterMap[qClusterId]; !ok {
			err := model.ClusterRightList.DeleteById(id)
			if !errors.Is(err, sql.ErrNoRows) && err != nil {
				return fmt.Errorf("Database err: %v", err)
			}
		}
	}
	res := model.UserList.UpdateInfoByUserId(company, fullName, email, phone, userId)

	return res
}

func ModifyPwdById(ctx context.Context) apibase.Result {
	log.Debugf("[ModifyUserById] ModifyUserById from EasyMatrix API ")

	userId, err := ctx.Values().GetInt("userId")
	if err != nil {
		log.Errorf("%v", err)
		return err
	}

	paramErrs := apibase.NewApiParameterErrors()
	oldPassword := ctx.FormValue("old_password")
	if oldPassword == "" {
		paramErrs.AppendError("$", fmt.Errorf("old_password is empty"))
	}
	//认证
	password := ctx.FormValue("new_password")
	if password == "" {
		paramErrs.AppendError("$", fmt.Errorf("new_password is empty"))
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	b, err := encrypt.PlatformEncrypt.CommonDecrypt([]byte(password))
	oldB, oldErr := encrypt.PlatformEncrypt.CommonDecrypt([]byte(oldPassword))
	if err != nil || oldErr != nil {
		return fmt.Sprintf("%s %s", err, oldErr)
	}

	err, info := model.UserList.GetInfoByUserId(userId)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	if info.PassWord != util.Md5(string(oldB)) {
		log.Errorf("Password Verify not pass, userId:%v password:%v old_password:%v", userId, info.PassWord, oldPassword)
		return fmt.Errorf("旧密码不正确")
	}
	//校验新密码是否与旧密码相同
	if info.PassWord == util.Md5(string(b)) {
		return fmt.Errorf("新密码不能与旧密码相同")
	}

	res := model.UserList.UpdatePwdByUserId(util.Md5(string(b)), userId)
	return res
}

func RemoveUserById(ctx context.Context) apibase.Result {
	log.Debugf("[RemoveUserById] RemoveUserById from EasyMatrix API ")

	if err := isAdmin(ctx); err != nil {
		return err
	}

	paramErrs := apibase.NewApiParameterErrors()
	targetUserIdStr := ctx.FormValue("targetUserId")
	if targetUserIdStr == "" {
		paramErrs.AppendError("$", fmt.Errorf("targetUserId is empty"))
	}
	targetUserId, err := strconv.Atoi(targetUserIdStr)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	err, user := model.UserList.GetInfoByUserId(targetUserId)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	if user.Status != DISABLE {
		return fmt.Errorf("只能删除停用的用户")
	}

	res := model.UserList.DeleteByUserId(targetUserId)
	if res {
		res = model.UserRoleList.DeleteByUserId(targetUserId)
		res = model.ClusterRightList.DeleteByUserId(targetUserId)
	}

	apibase.DeleteUserCache(targetUserId)
	defer func() {
		if err := addSafetyAuditRecord(ctx, "用户管理", "移除账号", "账号："+user.UserName); err != nil {
			log.Errorf("failed to add safety audit record\n")
		}
	}()

	return res
}

func GetCaptcha(ctx context.Context) apibase.Result {
	if VerifyCode {
		return captcha.NewLen(4)
	}
	return ""
}

func ProcessCaptcha(ctx context.Context) apibase.Result {
	id := ctx.FormValue("captchaId")
	digits := ctx.FormValue("captchaSolution")
	if !captcha.VerifyString(id, digits) {
		return errors.New("验证码错误\n")
	} else {
		return "验证码正确.\n"
	}
}

func ShowCaptcha(ctx context.Context) apibase.Result {
	file := ctx.Params().Get("captcha")
	if file == "" {
		http.NotFound(ctx.ResponseWriter(), ctx.Request())
		return nil
	}
	//dir, file := path.Split(ctx.Path())
	ext := path.Ext(file)
	id := file[:len(file)-len(ext)]
	if ext == "" || id == "" {
		http.NotFound(ctx.ResponseWriter(), ctx.Request())
		return nil
	}
	if ctx.FormValue("reload") != "" {
		captcha.Reload(id)
	}
	lang := strings.ToLower(ctx.FormValue("lang"))
	if serve(ctx, id, ext, lang) == captcha.ErrNotFound {
		http.NotFound(ctx.ResponseWriter(), ctx.Request())
	}
	return nil
}

func serve(ctx context.Context, id, ext, lang string) error {
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Header("Pragma", "no-cache")
	ctx.Header("Expires", "0")

	var content bytes.Buffer
	switch ext {
	case ".png":
		ctx.Header("Content-Type", "image/png")
		captcha.WriteImage(&content, id, captcha.StdWidth, captcha.StdHeight)
	case ".wav":
		ctx.Header("Content-Type", "audio/x-wav")
		captcha.WriteAudio(&content, id, lang)
	default:
		return captcha.ErrNotFound
	}

	ctx.ServeContent(bytes.NewReader(content.Bytes()), id+ext, time.Time{}, false)
	ctx.StatusCode(http.StatusOK)
	ctx.Done()
	return nil
}

func Disable(ctx context.Context) apibase.Result {
	log.Debugf("[Disable] Disable from EasyMatrix API ")

	if err := isAdmin(ctx); err != nil {
		return err
	}

	paramErrs := apibase.NewApiParameterErrors()
	targetUserIdStr := ctx.FormValue("targetUserId")
	if targetUserIdStr == "" {
		paramErrs.AppendError("$", fmt.Errorf("targetUserId is empty"))
	}
	targetUserId, err := strconv.Atoi(targetUserIdStr)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	res := model.UserList.UpdateStatusByUserId(DISABLE, targetUserId)

	apibase.DeleteUserCache(targetUserId)

	defer func() {
		err, info := model.UserList.GetInfoByUserId(targetUserId)
		if err != nil {
			log.Errorf("%v", err)
			return
		}
		if err := addSafetyAuditRecord(ctx, "用户管理", "禁用账号", "账号："+info.UserName); err != nil {
			log.Errorf("failed to add safety audit record\n")
		}
	}()

	return res
}

func Personal(ctx context.Context) apibase.Result {
	log.Debugf("[Personal] Personal from EasyMatrix API ")
	userId, err := ctx.Values().GetInt("userId")
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	err, info := model.UserList.GetInfoByUserId(userId)
	if err != nil {
		return err
	}
	info.PassWord = ""

	var resetPwd bool
	if cache.SysConfig.ForceResetPassword == FORCE_RESET_PASSWORD_ENABLE &&
		info.ResetPasswordStatus == PASSWORD_NOT_CHANGED {
		resetPwd = true
	}
	return map[string]interface{}{
		"info":           info,
		"reset_password": resetPwd,
	}
}

func Enable(ctx context.Context) apibase.Result {
	log.Debugf("[Unable] Unable from EasyMatrix API ")

	if err := isAdmin(ctx); err != nil {
		return err
	}

	paramErrs := apibase.NewApiParameterErrors()
	targetUserIdStr := ctx.FormValue("targetUserId")
	if targetUserIdStr == "" {
		paramErrs.AppendError("$", fmt.Errorf("targetUserId is empty"))
	}
	targetUserId, err := strconv.Atoi(targetUserIdStr)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	res := model.UserList.UpdateStatusByUserId(ENABLE, targetUserId)

	defer func() {
		err, info := model.UserList.GetInfoByUserId(targetUserId)
		if err != nil {
			log.Errorf("%v", err)
			return
		}
		if err := addSafetyAuditRecord(ctx, "用户管理", "启用账号", "账号："+info.UserName); err != nil {
			log.Errorf("failed to add safety audit record\n")
		}
	}()

	return res
}

func isAdmin(ctx context.Context) error {
	userId, err := ctx.Values().GetInt("userId")
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	err, info := model.UserRoleList.ListUserRoleByUserId(userId)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}

	isAdmin := false
	for _, v := range info {
		if v.RoleId == model.ROLE_ADMIN_ID {
			isAdmin = true
			break
		}
	}
	if !isAdmin {
		return fmt.Errorf("permission limit")
	}
	return nil
}

func Identity(ctx context.Context) apibase.Result {
	return VerifyIdentity
}

func GetSysconfigPlatformSecurity(ctx context.Context) apibase.Result {
	return cache.SysConfig.PlatFormSecurity
}

func ModifySysconfigPlatformSecurity(ctx context.Context) apibase.Result {
	var (
		request = cache.PlatFormSecurity{}
	)
	if err := ctx.ReadJSON(&request); err != nil {
		return errors.New("json body illegal")
	}
	if request.LoginEncrypt != "sm2" && request.LoginEncrypt != "rsa" {
		return fmt.Errorf("login_encrypt type need sm2 or rsa")
	}
	if request.AccountLogoutSleepTime <= 0 {
		return fmt.Errorf("account_logout_sleep_time need greater 0")
	}
	if request.ForceResetPassword != FORCE_RESET_PASSWORD_DISABLE &&
		request.ForceResetPassword != FORCE_RESET_PASSWORD_ENABLE {
		return fmt.Errorf("force_reset_password need 0 or 1")
	}
	if request.AccountLoginLockSwitch != LOGIN_LOCK_DISABLE &&
		request.AccountLoginLockSwitch != LOGIN_LOCK_ENABLE {
		return fmt.Errorf("account_login_lock_switch need 0 or 1")
	}
	if request.AccountLoginLockTime <= 0 {
		return fmt.Errorf("account_login_lock_time need greater 0")
	}
	if request.AccountLoginLimitError <= 0 {
		return fmt.Errorf("account_login_limit_error need greater 0")
	}
	err := cache.SysConfig.UpdatePlatFormSecurity(request)
	return err
}

func GetGlobalConfig(ctx context.Context) apibase.Result {
	return cache.SysConfig.GlobalConfig
}

func ModifyGlobalConfig(ctx context.Context) apibase.Result {
	var (
		request = cache.GlobalConfig{}
	)
	if err := ctx.ReadJSON(&request); err != nil {
		return errors.New("json body illegal")
	}
	if request.ServiceInstallTimeoutLimit == 0 {
		return errors.New("time can't set zero")
	}
	err := cache.SysConfig.UpdateGloablConfig(request)
	return err
}

// ModifyInspectConfig 	godoc
// @Summary      	修改巡检报告统计配置
// @Description  	修改巡检报告统计配置
// @Tags         	Inspect
// @Produce      	json
// @Success      	200		{object} string	"{"msg":"ok","code":0,"data":{}}"
// @Router       	/api/v2/platform/inspect/graph/config/update [post]
func ModifyInspectConfig(ctx context.Context) apibase.Result {
	var (
		request = cache.InspectConfig{}
	)
	if err := ctx.ReadJSON(&request); err != nil {
		return errors.New("json body illegal")
	}
	if request.FullGCTime <= 0 {
		return errors.New("full gc time must be a positive integer")
	}
	err := cache.SysConfig.UpdateInspectConfig(request)
	return err
}

// GetInspectConfig 	godoc
// @Summary      	获取巡检报告统计配置
// @Description  	获取巡检报告统计配置
// @Tags         	Inspect
// @Produce      	json
// @Success      	200		{object} 	cache.InspectConfig
// @Router       	/api/v2/platform/inspect/statisticsConfig [get]
func GetInspectConfig(ctx context.Context) apibase.Result {
	return cache.SysConfig.InspectConfig
}
