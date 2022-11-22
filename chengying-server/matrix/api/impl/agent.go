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
	"database/sql"
	"dtstack.com/dtstack/easymatrix/matrix/encrypt"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"dtstack.com/dtstack/easymatrix/addons/easykube/pkg/view/response"
	apibase "dtstack.com/dtstack/easymatrix/go-common/api-base"
	"dtstack.com/dtstack/easymatrix/matrix/agent"
	"dtstack.com/dtstack/easymatrix/matrix/base"
	"dtstack.com/dtstack/easymatrix/matrix/discover"
	"dtstack.com/dtstack/easymatrix/matrix/enums"
	"dtstack.com/dtstack/easymatrix/matrix/grafana"
	"dtstack.com/dtstack/easymatrix/matrix/host"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/node"
	kutil "dtstack.com/dtstack/easymatrix/matrix/k8s/util"
	xke_service "dtstack.com/dtstack/easymatrix/matrix/k8s/xke-service"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"dtstack.com/dtstack/easymatrix/matrix/util"
	"github.com/kataras/iris/context"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
)

const (
	LINE_LOG = "--------------------------------------------------------------"
)

func InstallInit(ctx context.Context) apibase.Result {
	return nil
}

func isLegalIp(ip string) bool {
	r, _ := regexp.Compile(`^((0|[1-9]\d?|1\d\d|2[0-4]\d|25[0-5])\.){3}(0|[1-9]\d?|1\d\d|2[0-4]\d|25[0-5])$`)
	return r.MatchString(ip)
}
func parseIpHostsWithRange(hostRange string) ([]string, error) {
	ips := []string{}
	hosts := strings.Split(hostRange, ",")
	for i, v := range hosts {
		if strings.Contains(v, "-") {
			lastDotIdx := strings.LastIndex(v, ".")
			if lastDotIdx > strings.Index(v, "-") {
				fmt.Println(i, v)
				return ips, fmt.Errorf("主机IP：%v 不符合规范，仅支持最后一个ip段为范围（例如：192.168.0.1-111）", v)
			}
			ipRanges := strings.Split(v[lastDotIdx+1:], "-")
			endIp, _ := strconv.Atoi(ipRanges[1])
			for i, _ := strconv.Atoi(ipRanges[0]); i <= endIp; i++ {
				ips = append(ips, v[:lastDotIdx+1]+strconv.Itoa(i))
			}
		} else {
			if !isLegalIp(v) {
				return ips, fmt.Errorf("主机IP：%v 不符合规范", v)
			}
			ips = append(ips, v)
		}
	}
	return ips, nil
}

type hostInfo struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password,omitempty"`
	Pk       string `json:"pk,omitempty"`
}

func execHostInfoCheck(params interface{}, cid int, suffixUrl string, isPwd bool) (interface{}, error) {
	var (
		result                = make(map[string][]string)
		hoststr               = ""
		connectErrorIps       = []string{}
		currentClusterExistIp = []string{}
		otherClusterExistIp   = []string{}
		wg                    = sync.WaitGroup{}
		baseUrl, _            = url.Parse("http://" + host.AgentInstall.AgentHost)
	)

	if isPwd {
		hoststr = params.(*util.PwdConnectParams).Host
	} else {
		hoststr = params.(*util.PkConnectParams).Host
	}

	ipHosts, err := parseIpHostsWithRange(hoststr)
	if err != nil {
		return result, err
	}
	//roles := strings.Split(params.(*util.PwdConnectParams).Role, ",")

	for _, ipHost := range ipHosts {
		_, nodeInfo := model.DeployNodeList.GetNodeInfoByNodeIp(ipHost) // 判断导入k8s集群的情况
		if nodeInfo.ID > 0 {
			otherClusterExistIp = append(otherClusterExistIp, ipHost)
		}

		hostRel, err := model.DeployClusterHostRel.GetClusterHostRelByIp(ipHost)
		if err == sql.ErrNoRows {
			log.Debugf("[skip exist cluster]", err)
			continue
		}
		//hostRoles := strings.Split(hostRel.Roles, ",")
		//if util.DiffNoOrderStringSlice(roles, hostRoles) {
		//
		//}
		if hostRel.ClusterId == cid {
			currentClusterExistIp = append(currentClusterExistIp, ipHost)
			continue
		}
		otherClusterExistIp = append(otherClusterExistIp, ipHost)
	}

	wg.Add(len(ipHosts))
	for _, v := range ipHosts {
		ip := v
		if isPwd {
			port, _ := strconv.Atoi(params.(*util.PwdConnectParams).Port)
			hostInfo := &hostInfo{
				Host:     ip,
				Port:     port,
				User:     params.(*util.PwdConnectParams).User,
				Password: params.(*util.PwdConnectParams).Password,
			}
			go func() {
				defer wg.Done()
				if !hostInfo.checkConnect(baseUrl, suffixUrl) {
					connectErrorIps = append(connectErrorIps, ip)
				}
			}()
		} else {
			port, _ := strconv.Atoi(params.(*util.PkConnectParams).Port)
			hostInfo := &hostInfo{
				Host: ip,
				Port: port,
				User: params.(*util.PkConnectParams).User,
				Pk:   params.(*util.PkConnectParams).Pk,
			}
			go func() {
				defer wg.Done()
				if !hostInfo.checkConnect(baseUrl, suffixUrl) {
					connectErrorIps = append(connectErrorIps, ip)
				}
			}()
		}
	}
	wg.Wait()
	result["connectErrorIps"] = connectErrorIps
	result["currentClusterExistIp"] = currentClusterExistIp
	result["otherClusterExistIp"] = otherClusterExistIp
	return result, nil
}

func (h *hostInfo) checkConnect(baseUrl *url.URL, suffixUrl string) bool {

	c := util.NewClient(util.DefaultClient)
	c.BaseURL = baseUrl
	// post to the api em[Agent-Server] http://agent-server:8889/api/v1/ssh/checkByUserPwd
	//c struct 内部已经实现 BaseUrl [类似http://xxxxx] ,所以提供的接口地址仅仅需要后缀 路由地址
	r, err := c.NewRequest("POST", suffixUrl, nil, h, "")
	if err != nil {
		log.Errorf("[batchConnect] %v Can not initialize http request, err : %v", h.Host, err.Error())
		return false
	}

	resBody := new(util.ResposeBody)
	resp, err := c.Do(r, resBody)

	if resp != nil && resp.StatusCode == 200 {
		if resBody.Code != 0 {
			log.Errorf("[batchConnect] %v error result response from EasyMatrix Agent-server err: %v:", h.Host, resBody.Msg)
			return false
		}
		return true
	} else if resp != nil {
		log.Errorf("[batchConnect] %v Bad response from EasyMatrix Agent-server http code : %v , err: %v, ", h.Host, resp.StatusCode, resBody.Msg)
	} else if err != nil {
		log.Errorf("[batchConnect] %v error from EasyMatrix Agent-server  error: %v, ", h.Host, err.Error())
	} else {
		log.Errorf("[batchConnect] %v unknown error from EasyMatrix Agent-server", h.Host)
	}
	return false
}

// CheckPwdConnect
// @Description  	通过password密码检查ssh连通性检查
// @Summary      	密码连通性测试
// @Tags         	agent
// @Accept          application/json
// @Produce 		application/json
// @Param           message body util.PwdConnectParams true "主机密码信息"
// @Success         200 {string}  string "{"msg":"ok","code":0,"data":null}"
// @Router          /api/v2/agent/install/pwdconnect [post]
func CheckPwdConnect(ctx context.Context) apibase.Result {
	log.Debugf("[CheckPwdConnect] check ssh connect by user password from EasyMatrix API ")
	paramErrs := apibase.NewApiParameterErrors()
	params := &util.PwdConnectParams{}
	cid, _ := GetCurrentClusterId(ctx)
	if err := ctx.ReadJSON(&params); err != nil {
		paramErrs.AppendError("$", err)
	}
	if password, err := encrypt.PlatformEncrypt.CommonDecrypt([]byte(params.Password)); err != nil {
		return err
	} else {
		params.Password = string(password)
	}

	paramErrs.CheckAndThrowApiParameterErrors()

	if res, err := execHostInfoCheck(params, cid, "/api/v1/ssh/checkByUserPwd", true); err != nil {
		return err
	} else {
		return res
	}
}

// CheckPkConnect
// @Description  	通过秘钥检查ssh连通性检查
// @Summary      	秘钥连通性测试
// @Tags         	agent
// @Accept          application/json
// @Produce 		application/json
// @Param           message body util.PkConnectParams  true "主机秘钥信息"
// @Success         200 {object}  string "{"msg":"ok","code":0,"data":null}"
// @Router          /api/v2/agent/install/pkconnect [post]
func CheckPkConnect(ctx context.Context) apibase.Result {
	//check post params
	log.Debugf("[CheckPkConnect] check ssh connect by user pk  from EasyMatrix API ")
	paramErrs := apibase.NewApiParameterErrors()
	params := &util.PkConnectParams{}
	cid, _ := GetCurrentClusterId(ctx)
	if err := ctx.ReadJSON(&params); err != nil {
		paramErrs.AppendError("$", err)
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	if res, err := execHostInfoCheck(params, cid, "/api/v1/ssh/checkByUserPk", false); err != nil {
		return err
	} else {
		return res
	}
}

// AgentInstallByPwd
// @Description  	install agent by pwd
// @Summary      	安装agent
// @Tags         	agent
// @Accept          application/json
// @Produce 		application/json
// @Param           message body util.PwdInstallParams true "主机密码信息"
// @Success         200 {object}  string "{"msg":"ok","code":0,"data":null}"
// @Router          /apv/v2/agent/install/pwdinstall [post]
func AgentInstallByPwd(ctx context.Context) apibase.Result {
	// check post params
	log.Debugf("[AgentInstallByPwd] install agent by user password from EasyMatrix API ")

	paramErrs := apibase.NewApiParameterErrors()
	params := &util.PwdInstallParams{}

	if err := ctx.ReadJSON(&params); err != nil {
		paramErrs.AppendError("$", err)
	}
	if password, err := encrypt.PlatformEncrypt.CommonDecrypt([]byte(params.Password)); err != nil {
		return err
	} else {
		params.Password = string(password)
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	cluster, _ := model.DeployClusterList.GetClusterInfoById(params.ClusterId)
	defer func() {
		if err := addSafetyAuditRecord(ctx, "集群管理", "添加主机", "集群名称："+cluster.Name+", 主机IP："+params.Host); err != nil {
			log.Errorf("failed to add safety audit record\n")
		}
	}()

	res := map[string]string{}
	wg := sync.WaitGroup{}
	for _, installHost := range strings.Split(params.Host, ",") {
		wg.Add(1)
		go execAgentInstallByPwd(installHost, params, ctx.Request().Host, res, &wg)
	}
	wg.Wait()
	return res
}

func execAgentInstallByPwd(installHost string, params *util.PwdInstallParams, requestHost string, res map[string]string, wg *sync.WaitGroup) (rt string) {
	defer func() {
		res[installHost] = rt
		wg.Done()
	}()

	// if node is installed, change roles
	err, info := model.DeployHostList.GetHostInfoByIp(installHost)
	if err == nil && info.Status == 7 {

		cluster, _ := model.DeployClusterList.GetClusterInfoById(params.ClusterId)
		rkeConfig, _ := xke_service.BuildRKEConfigFromRaw(cluster.Yaml.String)

		xke, err := xke_service.NewXkeService()
		if err != nil {
			return fmt.Sprintf("%v xke init error:", err)
		}
		// update role
		for i := range rkeConfig.Nodes {
			if rkeConfig.Nodes[i].Address == installHost {
				var rkeRoles []string
				for _, r := range strings.Split(params.Role, ",") {
					if _, ok := node.RoleToRkeRole[r]; ok {
						rkeRoles = append(rkeRoles, node.RoleToRkeRole[r])
					}
				}
				rkeConfig.Nodes[i].Role = rkeRoles
				break
			}
		}
		log.Infof("update k8s with rke config: %v", rkeConfig)

		rkeConfigYaml, err := yaml.Marshal(rkeConfig)
		if err != nil {
			return fmt.Sprintf("yaml marshal error:%v ", err)
		}
		err = xke.Create(rkeConfig.ClusterName, string(rkeConfigYaml), cluster.Id)
		if err != nil {
			return fmt.Sprintf("create rke error:%v ", err)
		} else {
			_ = model.DeployClusterHostRel.UpdateRolesWithSid(info.SidecarId, params.Role)
			return fmt.Sprintf("update roles success:%v ", installHost)
		}
	}

	err, _ = model.DeployHostList.GetHostInfoByIpAndStatus(installHost, host.InitStatus, host.InitInitializeShOk)
	if err == nil {
		rt = fmt.Sprintf("%v 主机正在接入", installHost)
		return
	}

	easyMatrixAgentServerHost := host.AgentInstall.AgentHost
	httpProtocol := "http://"
	flagUrl := httpProtocol + easyMatrixAgentServerHost
	baseUrl, _ := url.Parse(flagUrl)
	c := util.NewClient(util.DefaultClient)
	c.BaseURL = baseUrl

	group := params.Group
	if group == "" {
		group = "default"
	}
	err, aid := model.DeployHostList.AutoCreateAid(installHost, group)
	if err != nil {
		msg := "[AgentInstallByPwd] Can not initialize aid"
		log.Errorf("%s: %v", msg, err)
		rt = err.Error()
		return
	}
	err, cmd := host.AgentInstall.GetAgentInstallCmd(aid, requestHost, params.ClusterType, strconv.Itoa(params.ClusterId), params.Role)
	if err != nil {
		msg := "[AgentInstallByPwd] Can not initialize cmd"
		log.Errorf("%s: %v", msg, err)
		rt = err.Error()
		return
	}

	//installInfo := &util.PwdInstallParams{}
	installInfo := struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Cmd      string `json:"cmd"`
	}{}

	installInfo.Host = installHost
	installInfo.Port, _ = strconv.Atoi(params.Port)
	installInfo.User = params.User
	installInfo.Password = params.Password
	installInfo.Cmd = cmd

	// post to the api em[Agent-Server] http://agent-server:8889/api/v1/ssh/runWithUserPwd
	r, err := c.NewRequest("POST", "/api/v1/ssh/runWithUserPwd", nil, installInfo, "")
	if err != nil {
		msg := "[AgentInstallByPwd] Can not initialize http request"
		util.ResponseResult.SetResponseRes(base.StatusError, msg, "false", nil)
		log.Errorf("%s: %v", msg, err)
		rt = err.Error()
		return
	}

	resBody := new(util.ResposeBody)

	returnRespons := new(struct {
		Aid    int         `json:"aid"`
		Result interface{} `json:"result"`
	})

	resp, err := c.Do(r, resBody)
	//proxy pass the return data to EM2.0 client

	if resp != nil && resp.StatusCode == 200 {
		if resBody.Code != 0 {
			log.Errorf("[AgentInstallByPwd] error result response from EasyMatrix Agent-server err: %v:", resBody.Msg)
			model.DeployHostList.UpdateStatus(aid, -1, host.ERROR_SIDECAR_INSTALL+resBody.Msg)
			rt = fmt.Sprintf("error:%v", resBody.Msg)
			return
		}
		returnRespons.Aid = aid
		returnRespons.Result = resBody.Data
		model.DeployHostList.UpdateStatus(aid, 1, host.SUCCESS_SIDECAR_INSTALL)
		model.DeployHostList.UpdateSteps(aid, 1)
		rt = "ok"
		return
	} else if resp != nil {
		model.DeployHostList.UpdateStatus(aid, -1, host.ERROR_SIDECAR_INSTALL+"wrong status code:"+strconv.Itoa(resp.StatusCode))
		log.Errorf("[AgentInstallByPwd] Bad response from EasyMatrix Agent-server http code : %v , err: %v, ", resp.StatusCode, resBody.Msg)
		rt = fmt.Sprintf("http error : %v , msg: %v, ", resp.StatusCode, resBody.Msg)
		return

	} else if err != nil {
		model.DeployHostList.UpdateStatus(aid, -1, host.ERROR_SIDECAR_INSTALL+err.Error())
		log.Errorf("[AgentInstallByPwd] error from EasyMatrix Agent-server  err: %v, ", err.Error())
		rt = fmt.Sprintf("error: %v", err.Error())
		return
	} else {
		model.DeployHostList.UpdateStatus(aid, -1, host.ERROR_SIDECAR_INSTALL+" unknow error from EasyMatrix Agent-server")
		log.Errorf("[AgentInstallByPwd] unknow error from EasyMatrix Agent-server")
		rt = fmt.Sprintf("unknown error")
		return
	}
}

// AgentInstallCmd
// @Description  	install agent by cmd
// @Summary      	安装agent
// @Tags         	agent
// @Accept          application/json
// @Produce 		application/json
// @Param           type query string true "集群类型"    default(host)
// @Param           clusterId query string true "集群id" default(1)
// @Param           role query string true "集群角色"    default(Etcd,Control,Worker)
// @Success         200 {object}  string "curl -s 'http://172.16.82.176:8889/api/v1/deploy/sidecar/install/shell?TargetPath=/opt/dtstack/easymanager/easyagent&CallBack=aHR0cDovLzE3Mi4xNi44Mi4xNzY6ODg2NC9hcGkvdjIvYWdlbnQvaW5zdGFsbC9jYWxsYmFjaz9haWQ9LTE=&Type=hosts&ClusterId=1&Roles=Etcd,Control,Worker' | sh"
// @Router          /api/v2/agent/install/installCmd [post]
func AgentInstallCmd(ctx context.Context) apibase.Result {
	// check post params
	log.Debugf("[AgentInstallCmd] get install cmd from EasyMatrix API ")

	ctype := ctx.URLParam("type")
	clusterId := ctx.URLParam("clusterId")
	roles := ctx.URLParam("role")

	cid, _ := strconv.Atoi(clusterId)
	cluster, _ := model.DeployClusterList.GetClusterInfoById(cid)
	defer func() {
		if err := addSafetyAuditRecord(ctx, "集群管理", "添加主机", "集群名称："+cluster.Name+", 命令行接入"); err != nil {
			log.Errorf("failed to add safety audit record\n")
		}
	}()
	index := strings.Index(host.AgentInstall.StaticHost, "://")
	if index != -1 {
		index = index + 3
	}
	matrixHost := host.AgentInstall.StaticHost[index:]
	err, cmd := host.AgentInstall.GetAgentInstallCmd(-1, matrixHost, ctype, clusterId, roles)
	if err != nil {
		msg := "[AgentInstallCmd] Can not initialize cmd"
		log.Errorf("%s: %v", msg, err)
		return err
	}
	return cmd
}

// AgentInstallByPk
// @Description  	install agent by pk
// @Summary      	安装agent
// @Tags         	agent
// @Accept          application/json
// @Produce 		application/json
// @Param           message body util.PkInstallParams true "秘钥信息"
// @Success         200 {object}  string "{"msg":"ok","code":0,"data":null}"
// @Router          /api/v2/agent/install/pkinstall [post]
func AgentInstallByPk(ctx context.Context) apibase.Result {
	// check post params
	log.Debugf("[AgentInstallByPk] install agent by user pk from EasyMatrix API ")

	paramErrs := apibase.NewApiParameterErrors()
	params := &util.PkInstallParams{}

	if err := ctx.ReadJSON(&params); err != nil {
		paramErrs.AppendError("$", err)
	}

	paramErrs.CheckAndThrowApiParameterErrors()

	res := map[string]string{}
	wg := sync.WaitGroup{}
	for _, installHost := range strings.Split(params.Host, ",") {
		wg.Add(1)
		go execAgentInstallByPk(installHost, params, ctx.Request().Host, res, &wg)
	}
	wg.Wait()
	return res
}

func execAgentInstallByPk(installHost string, params *util.PkInstallParams, requestHost string, res map[string]string, wg *sync.WaitGroup) (rt string) {
	defer func() {
		res[installHost] = rt
		wg.Done()
	}()

	err, _ := model.DeployHostList.GetHostInfoByIpAndStatus(installHost, host.InitStatus, host.InitInitializeShOk)
	if err == nil {
		rt = fmt.Sprintf("%v 主机正在接入", installHost)
	}

	easyMatrixAgentServerHost := host.AgentInstall.AgentHost
	httpProtocol := "http://"
	flagUrl := httpProtocol + easyMatrixAgentServerHost
	baseUrl, _ := url.Parse(flagUrl)
	c := util.NewClient(util.DefaultClient)
	c.BaseURL = baseUrl

	group := params.Group
	if group == "" {
		group = "default"
	}
	err, aid := model.DeployHostList.AutoCreateAid(installHost, group)
	if err != nil {
		msg := "[AgentInstallByPk] Can not initialize aid"
		log.Errorf(msg)
		rt = fmt.Sprintf("can not initialize aid , err : %v", err.Error())
		return
	}

	err, cmd := host.AgentInstall.GetAgentInstallCmd(aid, requestHost, params.ClusterType, strconv.Itoa(params.ClusterId), params.Role)
	if err != nil {
		msg := "[AgentInstallByPk] Can not initialize cmd"
		log.Errorf(msg)
		rt = fmt.Sprintf("can not initialize cmd , err : %v", err.Error())
		return
	}

	//installInfo := &util.PkInstallParams{}
	installInfo := struct {
		Host string `json:"host"`
		Port int    `json:"port"`
		User string `json:"user"`
		Pk   string `json:"pk"`
		Cmd  string `json:"cmd"`
	}{}

	installInfo.Host = installHost
	installInfo.Port, _ = strconv.Atoi(params.Port)
	installInfo.User = params.User
	installInfo.Pk = params.Pk
	installInfo.Cmd = cmd

	// post to the api em[Agent-Server] http://agent-server:8889/api/v1/ssh/runWithUserPk
	r, err := c.NewRequest("POST", "/api/v1/ssh/runWithUserPk", nil, installInfo, "")
	if err != nil {
		msg := "[AgentInstallByPk] Can not initialize http request"
		util.ResponseResult.SetResponseRes(base.StatusError, msg, "false", nil)
		log.Errorf(msg)
		rt = fmt.Sprintf("can not initialize http request, err : %v", err.Error())
		return
	}

	resBody := new(util.ResposeBody)

	returnResponse := new(struct {
		Aid    int         `json:"aid"`
		Result interface{} `json:"result"`
	})

	resp, err := c.Do(r, resBody)
	//proxy pass the return data to EM2.0 client
	if resp != nil && resp.StatusCode == 200 {
		if resBody.Code != 0 {
			log.Errorf("[AgentInstallByPk] error result response from EasyMatrix Agent-server err: %v:", resBody.Msg)
			model.DeployHostList.UpdateStatus(aid, -1, host.ERROR_SIDECAR_INSTALL+resBody.Msg)
			log.Errorf("[AgentInstallByPk] error result response from EasyMatrix Agent-server err: %v", resBody.Msg)
			rt = fmt.Sprintf("error:%v", resBody.Msg)
			return
		}
		returnResponse.Aid = aid
		returnResponse.Result = resBody.Data
		model.DeployHostList.UpdateStatus(aid, 1, host.SUCCESS_SIDECAR_INSTALL)
		model.DeployHostList.UpdateSteps(aid, 1)
		rt = "ok"
		return
	} else if resp != nil {
		model.DeployHostList.UpdateStatus(aid, -1, host.ERROR_SIDECAR_INSTALL+"wrong status code:"+strconv.Itoa(resp.StatusCode))
		log.Errorf("[AgentInstallByPk] Bad response from EasyMatrix Agent-server http code : %v , err: %v, ", resp.StatusCode, resBody.Msg)
		rt = fmt.Sprintf("http error : %v , msg: %v, ", resp.StatusCode, resBody.Msg)
		return
	} else if err != nil {
		model.DeployHostList.UpdateStatus(aid, -1, host.ERROR_SIDECAR_INSTALL+err.Error())
		log.Errorf("[AgentInstallByPk] error from EasyMatrix Agent-server  err: %v, ", err.Error())
		rt = fmt.Sprintf("error: %v", err.Error())
		return
	} else {
		model.DeployHostList.UpdateStatus(aid, -1, host.ERROR_SIDECAR_INSTALL+" unknow error from EasyMatrix Agent-server")
		log.Errorf("[AgentInstallByPk] unknown error from EasyMatrix Agent-server")
		rt = fmt.Sprintf("unknown error")
		return
	}
}

// AgentInstallCheck
// @Description  	check status by aid
// @Summary      	检查agent
// @Tags         	agent
// @Accept          application/json
// @Produce 		application/json
// @Param           message body util.ApiShipperCheck true "agent id"
// @Success         200 {string} string  "{"msg":"ok","code":0,"data":{"aid":"","ip":"","status_msg":"","status":""}}"
// @Router          /api/v2/agent/install/checkinstall [post]
func AgentInstallCheck(ctx context.Context) apibase.Result {
	// check post params
	log.Debugf("[AgentInstallCheck] check agent install from EasyMatrix API ")

	paramErrs := apibase.NewApiParameterErrors()
	params := &util.ApiShipperCheck{}

	if err := ctx.ReadJSON(&params); err != nil {
		paramErrs.AppendError("$", err)
	}

	paramErrs.CheckAndThrowApiParameterErrors()

	err, checkRes := model.DeployHostList.GetHostInfoById(strconv.Itoa(params.Aid))

	if err == nil {
		installStatus := checkRes.Status
		msg := checkRes.ErrorMsg
		return map[string]interface{}{
			"aid":        checkRes.ID,
			"ip":         checkRes.Ip,
			"status_msg": msg,
			"status":     installStatus,
		}
	}

	return fmt.Errorf("Agent install failed! ")
}

// AgentInstallCheckAll
// @Description  	check status all
// @Summary      	检查agent
// @Tags         	agent
// @Accept          application/json
// @Produce 		application/json
// @Success         200 {object} string  "{"msg":"ok","code":0,"data":{"list":"","total":""}}"
// @Router          /api/v2/agent/install/checkinstallall [get]
func AgentInstallCheckAll(ctx context.Context) apibase.Result {
	// check post params
	log.Debugf("[AgentInstallCheckAll] check agent install from EasyMatrix API ")

	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	type hostInfo struct {
		Aid       string    `json:"aid"`
		Sid       string    `json:"sid"`
		HostName  string    `json:"host_name"`
		Ip        string    `json:"ip"`
		Status    int       `json:"status"`
		StatusMsg string    `json:"status_msg"`
		Updated   base.Time `json:"updated"`
		Created   base.Time `json:"created"`
		Group     string    `json:"group"`
	}
	var hosts []hostInfo
	values := []interface{}{clusterId}
	query := "SELECT h.id as aid, h.sid, h.hostname as host_name, h.ip, h.status, h.errorMsg as status_msg, h.updated, h.created, h.group " +
		"FROM deploy_cluster_list " +
		"LEFT JOIN deploy_cluster_host_rel ON deploy_cluster_list.id = deploy_cluster_host_rel.clusterId " +
		"LEFT JOIN deploy_host as h ON deploy_cluster_host_rel.sid = h.sid " +
		"WHERE deploy_cluster_list.id = ? AND deploy_cluster_list.is_deleted=0 "
	if err := model.USE_MYSQL_DB().Select(&hosts, query, values...); err != nil {
		apibase.ThrowDBModelError(err)
	}
	return map[string]interface{}{
		"list":  hosts,
		"total": len(hosts),
	}
}

// AgentInstallCheckByIp
// @Description  	check status by ip
// @Summary      	检查agent
// @Tags         	agent
// @Accept          application/json
// @Produce 		application/json
// @Success         200 {object} string  "{"msg":"ok","code":0,"data":{"ip":"","status":"","status_msg":""}}"
// @Router          /api/v2/agent/install/checkinstallbyip [post]
func AgentInstallCheckByIp(ctx context.Context) apibase.Result {
	// check post params
	log.Debugf("[AgentInstallCheckByIp] check agent install from EasyMatrix API ")

	paramErrs := apibase.NewApiParameterErrors()
	params := &util.ApiShipperCheckByIp{}

	if err := ctx.ReadJSON(&params); err != nil {
		paramErrs.AppendError("$", err)
	}

	paramErrs.CheckAndThrowApiParameterErrors()

	err, checkRes := model.DeployHostList.GetHostInfoByIp(params.Ip)

	if err == nil {
		installStatus := checkRes.Status
		if installStatus > 0 {
			msg := checkRes.ErrorMsg
			return map[string]interface{}{
				"ip":         checkRes.Ip,
				"status":     checkRes.Status,
				"status_msg": msg,
			}
		} else {
			log.Errorf("[AgentInstallCheckByIp]Failed, msg: %v, ip %v ", base.AgentInstallStateName[installStatus], checkRes.Ip)
			return fmt.Errorf("failed, msg: %v, ip %v ", base.AgentInstallStateName[installStatus], checkRes.Ip)
		}
	}

	return fmt.Errorf("Agent install failed! ")
}

// AgentInstallCheckBySid
// @Description  	check status by ip
// @Summary      	检查agent
// @Tags         	agent
// @Accept          application/json
// @Produce 		application/json
// @Param           sid  query   string   true  "7bfad86f-576b-474a-8697-f66b9fcbac9d"
// @Success         200 {object} string  "{"msg":"ok","code":0,"data":{"aid":"","ip":"","status":"","status_msg":""}}"
// @Router          /api/v2/agent/install/checkinstallbysid [get]
func AgentInstallCheckBySid(ctx context.Context) apibase.Result {
	// check get params
	log.Debugf("[AgentInstallCheckBySid] check agent install from EasyMatrix API ")

	paramErrs := apibase.NewApiParameterErrors()

	sidStr := ctx.FormValue("sid")

	log.Debugf("sid: %v", sidStr)

	_, err := uuid.FromString(sidStr)

	if err != nil {
		paramErrs.AppendError("sid", "sid is not uuid format")
	}

	paramErrs.CheckAndThrowApiParameterErrors()

	err, checkRes := model.DeployHostList.GetHostInfoBySid(sidStr)

	if err == nil {
		installStatus := checkRes.Status
		msg := checkRes.ErrorMsg
		return map[string]interface{}{
			"aid":        checkRes.ID,
			"ip":         checkRes.Ip,
			"status_msg": msg,
			"status":     installStatus,
		}
	}

	return fmt.Errorf(" Agent install failed! ")
}

func HostInstallHostRegister(sid, hostName, ip string, aid int, paramErrs *apibase.ApiParameterErrors) int {
	log.Infof("-->Step HostInstallHostRegister host register, sid: %v, aid: %v", sid, aid)
	log.OutputInfof(sid, "-->Step HostInstallHostRegister host register, sid: %v, aid: %v", sid, aid)
	var err error
	if aid > 0 {
		err = model.DeployHostList.UpdateWithAid(aid, sid, hostName, ip)
		if err != nil {
			log.Errorf("\t\tupdate host record error: %v", err)
			paramErrs.AppendError("AgentInstallCallBack update host", "rollingUpdate_ error: %v", err)
		}
	} else {
		err, aid = model.DeployHostList.InsertHostRecord(sid, hostName, ip)
		if err != nil {
			log.Errorf("\t\tinsert host record error: %v", err)
			paramErrs.AppendError("AgentInstallCallBack insert host", "insert error: %v", err)
		}
	}
	if aid > 0 {
		model.DeployHostList.UpdateStatus(aid, host.InstallSidecarOk, host.SUCCESS_SIDECAR_INSTALL)
	} else {
		log.Errorf("[\t\tHostInstallHostRegister error with aid: %v", aid)
	}
	log.Infof("<--Step HostInstallHostRegiste success, sid: %v, aid: %v", sid, aid)
	log.OutputInfof(sid, "<--Step HostInstallHostRegiste success, sid: %v, aid: %v", sid, aid)
	return aid
}

func NodeInstallNodeRegister(sid, hostName, ip string, aid int, ctx context.Context) error {
	var err error
	log.Infof("-->Step NodeInstallNodeRegister NodeInstallNodeRegister sid: %v", sid)
	log.OutputInfof(sid, "-->Step NodeInstallNodeRegister NodeInstallNodeRegister sid: %v", sid)
	mode := ctx.Request().Header.Get("Mode")
	clusterId := ctx.Request().Header.Get("ClusterId")
	if mode == "" || clusterId == "" {
		return fmt.Errorf("\t\t ClusterId or mode is null")
	}
	cid, err := strconv.Atoi(clusterId)
	if err != nil {
		return fmt.Errorf("\t\t clusterId is not number")
	}

	err, aid = model.DeployNodeList.InsertNodeRecord(sid, hostName, ip)
	if err != nil {
		log.Errorf("insert node record error %v", err.Error())
		return err
	}
	if aid > 0 {
		model.DeployNodeList.UpdateStatus(aid, host.K8SNodeDeploymentOk, host.K8S_SUCCESS_NODE_DEPLOYMENT_INIT)
		if mode != host.KUBERNETES_MODE {
			err, hostInfo := model.DeployHostList.GetHostInfoByIp(ip)
			if err != nil {
				log.Errorf("get host info by ip error %v", err)
			}
			model.DeployHostList.UpdateStatus(hostInfo.ID, host.K8SNodeDeploymentOk, host.K8S_SUCCESS_NODE_DEPLOYMENT_INIT)
		}
	}
	log.Infof("<--Step NodeInstallNodeRegister success")
	log.OutputInfof(sid, "<--Step NodeInstallNodeRegister success")

	// add role、k8s_version for import k8s cluster
	if mode == host.KUBERNETES_MODE {

		var content apibase.ApiResult
		var extraInfo response.ExtraInfoResponse

		params := agent.ExecRestParams{
			Method:  "GET",
			Path:    "clientgo/extraInfo?hostIp=" + ip,
			Timeout: "5s",
		}
		// 通过client-go取得 roles 和 version 信息
		resp, err := agent.AgentClient.ToExecRest(sid, &params, "")
		log.Infof("ExecRest ExtraInfo Response:%v", resp)
		if err != nil {
			return fmt.Errorf("ToExecRest extraInfo err:%v", err)
		}
		decodeResp, err := base64.StdEncoding.DecodeString(resp)
		if err != nil {
			log.Errorf("client-go response decode err:%v", err)
		}
		_ = json.Unmarshal(decodeResp, &content)
		data, _ := json.Marshal(content.Data)
		_ = json.Unmarshal(data, &extraInfo)

		err = agentInstallUpdateClusterHostRel(cid, sid, extraInfo.Roles)
		if err != nil {
			log.Errorf("DeployClusterHostRel InsertClusterHostRel err: %v", err)
		}
		hostRel, err := model.DeployClusterHostRel.GetClusterHostRelBySid(sid)
		if err != nil {
			log.Errorf("DeployClusterHostRel GetClusterHostRelBySid err: %v", err)
		}
		err = model.DeployClusterList.UpdateVersionById(hostRel.ClusterId, extraInfo.Version)
		if err != nil {
			log.Errorf("DeployClusterList UpdateVersionById err: %v", err)
		}
		discover.FlushNodeDiscover()
		grafana.Register(sid)
	}

	return nil
}

func AgentInstallStepTwo(sid string, aid int, targetPath, agentServer, debug string) bool {
	log.Debugf("[AgentInstallCallBack] AgentInstallStepTwo sid: %v", sid)

	err, records, _ := model.DeployHostList.GetHostListBySid(sid)
	if err == nil {
		for _, inst := range records {
			if inst.Steps >= host.InstallScriptWrapperOk {
				log.Debugf("[AgentInstallCallBack] AgentInstallStepTwo sid: %v already install script wrapper", sid)
				model.DeployHostList.UpdateStatus(aid, host.InstallScriptWrapperOk, host.SUCCESS_SCRIPT_WRAPPER_INSTALL)
				model.DeployHostList.UpdateSteps(aid, host.InstallScriptWrapperOk)
				return true
			}
		}
	} else {
		log.Errorf("[AgentInstallCallBack] AgentInstallStepTwo GetHostListBySid error: %v", err)
	}
	err, agentId := host.AgentInstall.LoopInstallScriptWrapper(sid, agentServer, targetPath, debug)
	if err != nil {
		model.DeployHostList.UpdateStatus(aid, host.InstallScriptWrapperFail, host.ERROR_SCRIPT_WRAPPER_INSTALL+err.Error())
		log.Errorf("[AgentInstallCallBack] iLoopInstallScriptWrapper error: %v", err)
		return false
	}
	if agentId != "" && host.AgentInstall.StartScriptWrapper(agentId) == nil {
		model.DeployHostList.UpdateStatus(aid, host.InstallScriptWrapperOk, host.SUCCESS_SCRIPT_WRAPPER_INSTALL)
		model.DeployHostList.UpdateSteps(aid, host.InstallScriptWrapperOk)
		return true
	}
	model.DeployHostList.UpdateStatus(aid, host.InstallScriptWrapperFail, host.ERROR_SCRIPT_WRAPPER_START)
	return false
}

func HostInstallStepEnvInit(sid string, aid int) error {
	log.Infof("-->Step HostInstallStepEnvInit sid: %v", sid)
	log.OutputInfof(sid, "-->Step HostInstallStepEnvInit sid: %v", sid)
	clusterHostRel, err := model.DeployClusterHostRel.GetClusterHostRelBySid(sid)
	if err != nil {
		return err
	}
	err, info := model.DeployHostList.GetHostInfoBySid(sid)
	if err != nil {
		log.Errorf("DeployHostList.GetHostInfoBySid error:%v", err)
	}
	//生成 operationid 并且落库
	operationId := uuid.NewV4().String()
	err = model.OperationList.Insert(model.OperationInfo{
		ClusterId:       clusterHostRel.ClusterId,
		OperationId:     operationId,
		OperationType:   enums.OperationType.HostInit.Code,
		OperationStatus: enums.ExecStatusType.Running.Code,
		ObjectType:      enums.OperationObjType.Host.Code,
		ObjectValue:     info.Ip,
	})
	if err != nil {
		log.Errorf("OperationList Insert err:%v", err)
	}
	execId := uuid.NewV4().String()
	err = model.ExecShellList.InsertExecShellInfo(clusterHostRel.ClusterId, operationId, execId, "", "", sid, enums.ShellType.Exec.Code)
	if err != nil {
		log.Errorf("InsertExecShellInfo err:%v", err)
	}
	if err := host.AgentInstall.EnvironmentInit(sid, execId); err != nil {
		log.Errorf("\t\tHostInstallStepEnvInit error: %v", err)
		model.DeployHostList.UpdateStatus(aid, host.InitInitializeShFail, host.ERROR_HOST_INIT+","+err.Error())
		return err
	}
	model.DeployHostList.UpdateStatus(aid, host.InitInitializeShOk, host.SUCCESS_HOST_INIT)
	model.DeployHostList.UpdateSteps(aid, host.InitInitializeShOk)
	log.Infof("<--Step HostInstallStepEnvInit success, sid %v", sid)
	log.OutputInfof(sid, "<--Step HostInstallStepEnvInit success, sid %v", sid)
	return nil
}

func HostInstallStepCluster(sid string, aid int, ctx context.Context) {
	log.Infof("-->Step HostInstallStepCluster sid: %v", sid)
	log.OutputInfof(sid, "-->Step HostInstallStepCluster sid: %v", sid)
	ctype := ctx.Request().Header.Get("Type")
	clusterId := ctx.Request().Header.Get("ClusterId")
	roles := ctx.Request().Header.Get("Roles")
	if ctype == "" || clusterId == "" {
		log.Errorf("\t\tClusterId or Types is null")
		return
	}
	cid, err := strconv.Atoi(clusterId)
	if err != nil {
		log.Errorf("\t\tclusterId is not number")
		return
	}
	if err := agentInstallUpdateClusterHostRel(cid, sid, roles); err != nil {
		log.Errorf("\t\t%v", err.Error())
		return
	}
	if ctype == model.DEPLOY_CLUSTER_TYPE_KUBERNETES {
		agentInstallProcessK8SCluster(cid, aid, sid, roles)
	}
	log.Infof("<--Step HostInstallStepCluster Success sid: %v", sid)
	log.OutputInfof(sid, "<--Step HostInstallStepCluster Success sid: %v", sid)
}

func agentInstallUpdateClusterHostRel(clusterId int, sid, roles string) error {
	//rel, err := model.DeployClusterHostRel.GetClusterHostRelBySid(sid)
	//if err == nil &&  rel.Roles == roles {
	//	return fmt.Errorf("\t\tcluster host rel existed, sid: %v!", sid)
	//}
	roles = strings.Replace(roles, "controlplane", "Control", -1)
	roles = strings.Replace(roles, "etcd", "Etcd", -1)
	roles = strings.Replace(roles, "worker", "Worker", -1)
	_, err := model.DeployClusterHostRel.InsertClusterHostRel(clusterId, sid, roles)
	if err != nil {
		return err
	}
	return nil
}

func agentInstallProcessK8SCluster(clusterId, aid int, sid, roles string) error {
	info, err := model.DeployClusterList.GetClusterInfoById(clusterId)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	if info.Mode == model.DEPLOY_CLUSTER_MODE_IMPORT {
		model.DeployHostList.UpdateStatus(aid, host.K8SNodeInitializeOk, host.K8S_SUCCESS_NODE_INIT)
		model.DeployHostList.UpdateSteps(aid, host.K8SNodeInitializeOk)
		return nil
	}
	err = host.AgentInstall.DockerEnvironmentInit(sid)
	if err != nil {
		log.Errorf("%v", err.Error())
		model.DeployHostList.UpdateStatus(aid, host.K8SDockerInitializeFail, err.Error())
		return err
	}
	// update host status
	model.DeployHostList.UpdateStatus(aid, host.K8SDockerInitializeOk, host.K8S_SUCCESS_DOCKCER_INIT)
	model.DeployHostList.UpdateSteps(aid, host.K8SDockerInitializeOk)

	//do rke create in queue
	node.NodeManager.AddNode(clusterId, aid, sid, info.Name, roles)

	return nil
}

//1:管控安装成功, -1:管控安装失败, 2:script安装成功, -2:script安装失败, 3:主机初始化成功, -3:主机初始化失败
func AgentInstallCallBack(ctx context.Context) apibase.Result {

	log.Infof("%v", LINE_LOG)
	log.Infof("->AgentInstallCallBack: %v", ctx.Request().RequestURI)
	log.Infof("\t\tHeaders: %v", ctx.Request().Header)

	paramErrs := apibase.NewApiParameterErrors()
	aid := ctx.URLParam("aid")
	if aid == "" {
		log.Errorf("\t\taid is empty")
		paramErrs.AppendError("AgentInstallCallBack aid check", " aid is empty")
	}
	sid := ctx.Request().Header.Get("SID")
	if sid == "" {
		log.Errorf("\t\tsid is empty")
		paramErrs.AppendError("AgentInstallCallBack sid check", " sid is empty")
	}
	hostName := ctx.Request().Header.Get("HostName")
	if hostName == "" {
		log.Errorf("\t\thostName is empty")
		paramErrs.AppendError("AgentInstallCallBack hostName check", " hostName is empty")
	}
	ip := ctx.Request().Header.Get("IP")
	if ip == "" {
		log.Errorf("\t\tip is empty")
		paramErrs.AppendError("AgentInstallCallBack ip check", " ip is empty")
	}
	iid, err := strconv.Atoi(aid)
	if err != nil {
		log.Errorf("\t\taid is not number")
		paramErrs.AppendError("AgentInstallCallBack aid format", " aid is not number")
	}
	clusterId := ctx.Request().Header.Get("ClusterId")
	cid, err := strconv.Atoi(clusterId)
	if err != nil {
		log.Errorf("\t\tclusterId is not number")
		paramErrs.AppendError("AgentInstallCallBack clusterId format", " clusterId is not number")
	}
	model.DeployClusterHostRel.InsertClusterHostRel(cid, sid, "")
	deployment := ctx.Request().Header.Get("Deploy")

	// redirect log output
	cluster, _ := model.DeployClusterList.GetClusterInfoById(cid)
	fileName := kutil.BuildClusterLogName(cluster.Name, cid)
	logf, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	defer logf.Close()
	if err == nil {
		log.NewOutputPath(sid, logf)
		// defer log.CloseOutputPath(sid)
	} else {
		log.Errorf(err.Error())
	}

	log.OutputInfof(sid, "%v", LINE_LOG)
	log.OutputInfof(sid, "AgentInstallCallBack: %v", ctx.Request().RequestURI)
	log.OutputInfof(sid, "Headers: %v", ctx.Request().Header)

	switch deployment {
	case "deployment":
	case "daemonset":
		// Handler for k8s node deployment/daemonset
		err := NodeInstallNodeRegister(sid, hostName, ip, iid, ctx)
		if err != nil {
			log.Errorf("\t\t%v", err.Error())
			log.OutputInfof(sid, "error: %v", err.Error())
		}
	default:
		// Handler for host or k8s cluster create callback
		// 根据sid获取主机是否已经初始化过，如果初始化过，则不执行主机初始化脚本
		err, hostInfo := model.DeployHostList.GetHostInfoBySid(sid)
		if err != nil && err != sql.ErrNoRows {
			log.OutputInfof(sid, "error: %v", err.Error())
			return err
		}
		isInit := hostInfo.Status == host.InitInitializeShOk && hostInfo.ErrorMsg == host.SUCCESS_HOST_INIT
		iid = HostInstallHostRegister(sid, hostName, ip, iid, paramErrs)
		paramErrs.CheckAndThrowApiParameterErrors()
		if !isInit {
			if err := HostInstallStepEnvInit(sid, iid); err != nil {
				log.OutputInfof(sid, "error: %v", err.Error())
				return err
			}
		} else {
			// 已经初始化过，直接将状态更新为主机初始化成功
			log.Infof("sid %s already initialized success, skip init environment", sid)
			model.DeployHostList.UpdateStatus(iid, host.InitInitializeShOk, host.SUCCESS_HOST_INIT)
			model.DeployHostList.UpdateSteps(iid, host.InitInitializeShOk)
		}
		HostInstallStepCluster(sid, iid, ctx)
		discover.FlushNodeDiscover()
		grafana.Register(sid)
	}
	//if err := addSafetyAuditRecord(ctx, "集群管理", "添加主机", "集群名称："+cluster.Name+", 主机IP："+ip); err != nil {
	//	log.Errorf("failed to add safety audit record\n")
	//}
	return nil
}

// AgentHosts
// @Description  	list agent
// @Summary      	查询所有agent
// @Tags         	agent
// @Accept          application/json
// @Produce 		application/json
// @Param           message body util.PwdInstallParams true "主机密码信息"
// @Success         200 {object}  string "{"msg":"ok","code":0,"data":{"hosts":"","count":""}}"
// @Router          /api/v2/agent/install/pwdinstall [get]
func AgentHosts(ctx context.Context) apibase.Result {
	log.Debugf("AgentHosts: %v", ctx.Request().RequestURI)

	var baseQuery, whereCause string
	var values []interface{}
	var err error

	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return map[string]interface{}{
			"hosts": []hostInfo{},
			"count": 0,
		}
	}

	hostOrIp := ctx.URLParam("host_or_ip")
	productName := ctx.URLParam("product_name")
	parentProductName := ctx.URLParam("parent_product_name")
	group := ctx.URLParam("group")
	values = []interface{}{"%" + hostOrIp + "%", "%" + hostOrIp + "%"}
	if productName != "" {
		whereCause = ` AND deploy_product_list.product_name IN (`
		for i, v := range strings.Split(productName, ",") {
			if i > 0 {
				whereCause += `,`
			}
			whereCause += `?`
			values = append(values, v)
		}
		whereCause += `)`
	}
	if parentProductName != "" {
		whereCause += ` AND deploy_product_list.parent_product_name=?`
		values = append(values, parentProductName)
	}
	if group != "" {
		whereCause += ` AND deploy_host.group=?`
		values = append(values, group)
	}
	//add host deleted fix
	whereCause += ` AND deploy_host.isDeleted=?`
	values = append(values, 0)

	if clusterId > 0 {
		whereCause += ` AND deploy_instance_list.cluster_id=?`
		values = append(values, clusterId)
	}

	baseQuery = fmt.Sprintf(`FROM deploy_host
LEFT JOIN deploy_instance_list ON deploy_host.sid = deploy_instance_list.sid
LEFT JOIN deploy_product_list ON deploy_instance_list.pid = deploy_product_list.id
LEFT JOIN sidecar_list ON sidecar_list.id = deploy_host.sid
WHERE deploy_host.sid != '' AND (deploy_host.hostname LIKE ? OR deploy_host.ip LIKE ?)%s`, whereCause)

	//TODO

	type hostInfo struct {
		model.HostInfo
		RunUser                string                  `json:"run_user"`
		ProductNameList        string                  `json:"product_name_list" db:"product_name_list"`
		ProductNameDisplayList string                  `json:"product_name_display_list" db:"product_name_display_list"`
		ProductIdList          string                  `json:"pid_list" db:"pid_list"`
		MemSize                int64                   `json:"mem_size" db:"mem_size"`
		MemUsage               int64                   `json:"mem_usage" db:"mem_usage"`
		DiskUsage              sql.NullString          `json:"disk_usage" db:"disk_usage"`
		NetUsage               sql.NullString          `json:"net_usage" db:"net_usage"`
		MemSizeDisplay         string                  `json:"mem_size_display"`
		MemUsedDisplay         string                  `json:"mem_used_display"`
		DiskSizeDisplay        string                  `json:"disk_size_display"`
		DiskUsedDisplay        string                  `json:"disk_used_display"`
		FileSizeDisplay        string                  `json:"file_size_display"`
		FileUsedDisplay        string                  `json:"file_used_display"`
		NetUsageDisplay        []model.NetUsageDisplay `json:"net_usage_display,omitempty"`
		IsRunning              bool                    `json:"is_running"`
	}

	var count int
	var hostsList []hostInfo
	query := "SELECT COUNT(DISTINCT deploy_host.sid) " + baseQuery
	whoamiCmd := "#!/bin/sh\n whoami"
	if err = model.USE_MYSQL_DB().Get(&count, query, values...); err != nil {
		log.Errorf("AgentHosts count query: %v, values %v, err: %v", query, values, err)
		apibase.ThrowDBModelError(err)
	}
	if count > 0 {
		query = "SELECT deploy_host.*, " +
			"IFNULL(GROUP_CONCAT(DISTINCT(deploy_product_list.product_name)),'') AS product_name_list, " +
			"IFNULL(GROUP_CONCAT(DISTINCT(deploy_product_list.product_name_display)),'') AS product_name_display_list, " +
			"IFNULL(GROUP_CONCAT(DISTINCT(deploy_product_list.id)),'') AS pid_list," +
			"sidecar_list.mem_size, sidecar_list.mem_usage, sidecar_list.disk_usage, sidecar_list.net_usage " +
			baseQuery + " GROUP BY deploy_host.sid " + apibase.GetPaginationFromQueryParameters(nil, ctx, model.HostInfo{}).AsQuery()
		if err = model.USE_MYSQL_DB().Select(&hostsList, query, values...); err != nil {
			log.Errorf("AgentHosts query: %v, values %v, err: %v", query, values, err)
			apibase.ThrowDBModelError(err)
		}
		for i, list := range hostsList {
			if time.Now().Sub(time.Time(list.UpdateDate)) < 3*time.Minute {
				hostsList[i].IsRunning = true
			}
			hostsList[i].MemSizeDisplay = sizeConvert(list.MemSize)
			hostsList[i].MemUsedDisplay = sizeConvert(list.MemUsage)
			if list.DiskUsage.Valid {
				hostsList[i].DiskSizeDisplay, hostsList[i].DiskUsedDisplay, hostsList[i].FileSizeDisplay, hostsList[i].FileUsedDisplay = diskUsageConvert(list.DiskUsage.String)
			}
			if list.NetUsage.Valid {
				hostsList[i].NetUsageDisplay = netUsageConvert(list.NetUsage.String)
			}
			if list.IsDeleted == 0 && list.Status > 0 && hostsList[i].IsRunning {
				content, err := agent.AgentClient.ToExecCmdWithTimeout(list.SidecarId, "", whoamiCmd, "5s", "", "")
				if err != nil {
					//exec failed
					content = err.Error()
				}
				user := strings.Replace(content, LINUX_SYSTEM_LINES, "", -1)
				hostsList[i].RunUser = user
			}
		}
	}

	return map[string]interface{}{
		"hosts": hostsList,
		"count": count,
	}
}

// AgentHostGroups
// @Description  	通过密码安装agent
// @Summary      	安装agent
// @Tags         	agent
// @Accept          application/json
// @Produce 		application/json
// @Param           host_or_ip  query  string  true  "host ip"
// @Param           product_name  query  string  true  "product name"
// @Param           parent_product_name  query  string  true  "parent product name"
// @Param           group  query  string  false  "agent group"
// @Success         200 {object}  string "{"msg":"ok","code":0,"data":["default"]}"
// @Router          /api/v2/agent/install/hostgroups [post]
func AgentHostGroups(ctx context.Context) apibase.Result {
	log.Debugf("AgentHostGroups: %v", ctx.Request().RequestURI)

	var whereCause string
	var values []interface{}
	var err error

	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return []string{}
	}

	hostOrIp := ctx.URLParam("host_or_ip")
	productName := ctx.URLParam("product_name")
	parentProductName := ctx.URLParam("parent_product_name")
	group := ctx.URLParam("group")
	values = []interface{}{"%" + hostOrIp + "%", "%" + hostOrIp + "%"}
	if productName != "" {
		whereCause = ` AND deploy_product_list.product_name IN (`
		for i, v := range strings.Split(productName, ",") {
			if i > 0 {
				whereCause += `,`
			}
			whereCause += `?`
			values = append(values, v)
		}
		whereCause += `)`
	}
	if parentProductName != "" {
		whereCause += ` AND deploy_product_list.parent_product_name=?`
		values = append(values, parentProductName)
	}
	if group != "" {
		whereCause += ` AND deploy_host.group=?`
		values = append(values, group)
	}

	// add clusterId
	whereCause += ` AND deploy_cluster_host_rel.clusterId=?`
	values = append(values, clusterId)

	query := fmt.Sprintf(`SELECT DISTINCT deploy_host.group FROM deploy_host
LEFT JOIN deploy_cluster_host_rel ON deploy_host.sid = deploy_cluster_host_rel.sid
LEFT JOIN deploy_instance_list ON deploy_cluster_host_rel.sid = deploy_instance_list.sid
LEFT JOIN deploy_product_list ON deploy_instance_list.pid = deploy_product_list.id
LEFT JOIN sidecar_list ON sidecar_list.id = deploy_host.sid
WHERE deploy_host.isDeleted=0 AND deploy_cluster_host_rel.is_deleted=0 AND deploy_host.sid != '' AND (deploy_host.hostname LIKE ? OR deploy_host.ip LIKE ?)%s GROUP BY deploy_host.sid`, whereCause)

	var groups []string
	if err = model.USE_MYSQL_DB().Select(&groups, query, values...); err != nil {
		log.Errorf("AgentHostGroups query: %v, values %v, err: %v", query, values, err)
		apibase.ThrowDBModelError(err)
	}

	return groups
}

type hostGroupRenameParam struct {
	Old string `json:"old"`
	New string `json:"new"`
}

// AgentHostGroupRename
// @Description  	rename group name
// @Summary      	重命名组
// @Tags         	agent
// @Accept          application/json
// @Produce 		application/json
// @Param           message body hostGroupRenameParam true "命名参数"
// @Success         200 {object}  string "{"msg":"ok","code":0,"data":null}"
// @Router          /api/v2/agent/install/hostgroup_rename [post]
func AgentHostGroupRename(ctx context.Context) apibase.Result {
	log.Debugf("AgentHostGroupRename: %v", ctx.Request().RequestURI)

	param := hostGroupRenameParam{}
	if err := ctx.ReadJSON(&param); err != nil {
		return fmt.Errorf("ReadJSON err: %v", err)
	}
	if param.Old == "" || param.New == "" {
		return fmt.Errorf("old or new param is empty")
	}

	return model.DeployHostList.UpdateGroup(param.Old, param.New)
}

type hostMoveParam struct {
	Aid     []int  `json:"aid"`
	ToGroup string `json:"to_group"`
}

// AgentHostMove
// @Description  	move host group
// @Summary      	改变主机组名
// @Tags         	agent
// @Accept          application/json
// @Produce 		application/json
// @Param           message body hostMoveParam true "主机密码信息"
// @Success         200 {object}  string "{"msg":"ok","code":0,"data":null}"
// @Router          /api/v2/agent/install/hostmove [post]
func AgentHostMove(ctx context.Context) apibase.Result {
	log.Debugf("AgentHostMove: %v", ctx.Request().RequestURI)

	param := hostMoveParam{}
	if err := ctx.ReadJSON(&param); err != nil {
		return fmt.Errorf("ReadJSON err: %v", err)
	}
	if len(param.Aid) == 0 || param.ToGroup == "" {
		return fmt.Errorf("aid or to_group param is empty")
	}

	for _, aid := range param.Aid {
		if err := model.DeployHostList.UpdateGroupWithAid(aid, param.ToGroup); err != nil {
			return err
		}
	}

	return nil
}

func diskUsageConvert(diskInfoStr string) (diskSizeDisplay string, diskUsedDisplay string, fileSizeDisplay string, fileUsedDisplay string) {
	var diskUsages []model.DiskUsage
	if err := json.Unmarshal([]byte(diskInfoStr), &diskUsages); err != nil {
		return
	}
	var diskSize, diskUsed, fileSize, fileUsed int64
	for _, diskUsage := range diskUsages {
		if diskUsage.MountPoint == "/" {
			fileSize += int64(diskUsage.TotalSpace)
			fileUsed += int64(diskUsage.UsedSpace)
		} else {
			diskSize += int64(diskUsage.TotalSpace)
			diskUsed += int64(diskUsage.UsedSpace)
		}
	}
	diskSizeDisplay = sizeConvert(diskSize)
	diskUsedDisplay = sizeConvert(diskUsed)
	fileSizeDisplay = sizeConvert(fileSize)
	fileUsedDisplay = sizeConvert(fileUsed)
	return
}

func netUsageConvert(netInfoStr string) []model.NetUsageDisplay {
	var netUsages []model.NetUsage
	if err := json.Unmarshal([]byte(netInfoStr), &netUsages); err != nil {
		return nil
	}

	netUsagesDisplay := make([]model.NetUsageDisplay, len(netUsages))
	for i := 0; i < len(netUsages); i++ {
		netUsagesDisplay[i].IfName = netUsages[i].IfName
		netUsagesDisplay[i].IfIp = netUsages[i].IfIp
		netUsagesDisplay[i].BytesSent = sizeConvert(int64(netUsages[i].BytesSent))
		netUsagesDisplay[i].BytesRecv = sizeConvert(int64(netUsages[i].BytesRecv))
	}

	return netUsagesDisplay
}

func sizeConvert(size int64) string {
	if size <= 0 {
		return fmt.Sprintf("%d", size)
	}

	sizeUnits := [...]string{"B", "KB", "MB", "GB", "TB"}
	f := float32(size)
	for _, v := range sizeUnits {
		if f < 1024 {
			return fmt.Sprintf("%.2f"+v, f)
		} else {
			f = f / 1024
		}
	}
	return fmt.Sprintf("%.2f"+sizeUnits[len(sizeUnits)-1], f)
}

// AgentHostService
// @Description  	list service by pid
// @Summary      	查看主机服务
// @Tags         	agent
// @Accept          application/json
// @Produce 		application/json
// @Param           pid_list  query     string     false  "["1","2"]"
// @Param           ip  query     string   false  "主机ip"
// @Success         200 {object}  string "{"msg":"ok","code":0,"data":[{"product_name":"","product_name_display":"","group":"","service_name_list":"","service_name_display_list":"service_name_display_list"}]}"
// @Router          /api/v2/agent/install/hostService [get]
func AgentHostService(ctx context.Context) apibase.Result {
	log.Debugf("AgentHostService: %v", ctx.Request().RequestURI)

	paramErrs := apibase.NewApiParameterErrors()
	pidList := ctx.URLParam("pid_list")
	if pidList == "" {
		log.Errorf("[AgentHostService] pid_list is not null")
		paramErrs.AppendError("AgentHostService params error", " pid_list is not null")
	}
	ip := ctx.URLParam("ip")
	if ip == "" {
		paramErrs.AppendError("AgentHostService params error", " ip is not null")
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	pids := strings.Split(pidList, ",")
	var values []interface{}
	var whereCause string
	//values = append(values, ip)
	for i, pid := range pids {
		if i > 0 {
			whereCause += ","
		}
		whereCause += "?"
		values = append(values, pid)
	}

	type serviceInfo struct {
		ProductName            string `db:"product_name" json:"product_name"`
		ProductNameDisplay     string `db:"product_name_display" json:"product_name_display"`
		Group                  string `db:"group" json:"group"`
		ServiceNameList        string `db:"service_name_list" json:"service_name_list"`
		ServiceNameDisplayList string `db:"service_name_display_list" json:"service_name_display_list"`
	}
	var services []serviceInfo
	query := "SELECT deploy_product_list.product_name,deploy_product_list.product_name_display,deploy_instance_list.group, " +
		"IFNULL(GROUP_CONCAT(DISTINCT(deploy_instance_list.service_name)),'') AS service_name_list, " +
		"IFNULL(GROUP_CONCAT(DISTINCT(deploy_instance_list.service_name_display)),'') AS service_name_display_list " +
		"FROM deploy_instance_list " +
		"LEFT JOIN deploy_product_list ON deploy_product_list.id = deploy_instance_list.pid " +
		"WHERE deploy_instance_list.ip LIKE '%s%%' AND deploy_instance_list.pid IN (%s)" +
		"GROUP BY deploy_instance_list.pid,deploy_instance_list.group " +
		"ORDER BY deploy_instance_list.pid,deploy_instance_list.group"
	if err := model.USE_MYSQL_DB().Select(&services, fmt.Sprintf(query, ip, whereCause), values...); err != nil {
		apibase.ThrowDBModelError(err)
	}

	return services
}

type hostDeleteParam struct {
	Aid  []int  `json:"aid"`
	Type string `json:"type"`
}

// AgentHostDelete
// @Description  	delete agent
// @Summary      	删除agent
// @Tags         	agent
// @Accept          application/json
// @Produce 		application/json
// @Param           message body hostDeleteParam true "主机信息"
// @Success         200 {object}  string "{"msg":"ok","code":0,"data":null}"
// @Router          /agent/install/hostdelete [post]
func AgentHostDelete(ctx context.Context) apibase.Result {
	log.Infof("AgentHostDelete: %v", ctx.Request().RequestURI)
	clusterId, err := GetCurrentClusterId(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	param := hostDeleteParam{}
	if err := ctx.ReadJSON(&param); err != nil {
		return fmt.Errorf("ReadJSON err: %v", err)
	}
	if len(param.Aid) == 0 {
		return fmt.Errorf("aid param is empty")
	}
	//tx := model.USE_MYSQL_DB().MustBegin()
	//defer func() {
	//	if r := recover(); r != nil {
	//		tx.Rollback()
	//	}
	//}()
	for _, aid := range param.Aid {
		_, info := model.DeployHostList.GetHostInfoById(strconv.Itoa(aid))
		model.DeployHostList.DeleteWithAid(aid)

		model.DeployClusterHostRel.DeleteWithSid(info.SidecarId)

		//先获取信息后再删除
		err, nodeInfo := model.DeployNodeList.GetNodeInfoByNodeIp(info.Ip) // 判断导入k8s集群的情况
		//server_ip_node
		if err != nil {
			if err != sql.ErrNoRows {
				return err
			}
			model.DeployInstanceList.DeleteByIp(info.Ip) //主机
		} else {
			model.DeployInstanceList.DeleteBySid(nodeInfo.SidecarId) //k8s自建
		}

		model.DeployNodeList.DeleteWithIp(info.Ip)

		//server_ip_list
		model.DeployServiceIpList.HostOffByIp(info.Ip)

		//server_ip_node
		model.ServiceIpNode{}.DeleteByIp(info.Ip)

		//service_health_check
		err = model.HealthCheck.DeleteByIp(info.Ip)
		if err != nil {
			return err
		}

		//schema_multi_fields
		productServiceTupleList, err := model.SchemaMultiField.GetProductServiceByIp(clusterId, info.Ip)
		if err != nil {
			log.Errorf("%v", err)
		}
		model.SchemaMultiField.DeleteByIp(clusterId, info.Ip)
		for _, tuple := range productServiceTupleList {
			ValidateSchemaFields(clusterId, tuple.ProductName, tuple.ServiceName)
		}
	}
	return nil
}
