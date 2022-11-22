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
	"dtstack.com/dtstack/easymatrix/go-common/api-base"
	"dtstack.com/dtstack/easymatrix/matrix/base"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"dtstack.com/dtstack/easymatrix/matrix/task"
	"dtstack.com/dtstack/easymatrix/matrix/util"
	"errors"
	"fmt"
	"github.com/gorhill/cronexpr"
	"github.com/kataras/iris/context"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	pythonExt = ".py"
	shellExt  = ".sh"
)

var taskLock sync.Mutex

type ResTaskInfo struct {
	Id         int                 `json:"id"`
	Name       string              `json:"name"`
	Describe   string              `json:"describe"`
	CreateTime string              `json:"create_time"`
	Status     int                 `json:"status"`
	Spec       string              `json:"spec"`
	Hosts      []model.ResHostInfo `json:"hosts"`
	EndTime    string              `json:"end_time"`
	ExecStatus model.Status        `json:"exec_status"`
	RunStatus  model.Status        `json:"run_status"`
}

type ResTaskLog struct {
	EndTime    string        `json:"end_time"`
	ExecType   int           `json:"exec_type"`
	ExecStatus model.Status  `json:"exec_status"`
	Children   []LogChildren `json:"children"`
}

type LogChildren struct {
	EndTime    string       `json:"end_time"`
	Ip         string       `json:"ip"`
	ExecType   int          `json:"exec_type"`
	ExecStatus model.Status `json:"exec_status"`
	ExecResult string       `json:"exec_result"`
}

func TaskList(ctx context.Context) apibase.Result {
	log.Debugf("[Task->TaskList] TaskList from EasyMatrix API ")

	type res struct {
		List  []ResTaskInfo `json:"list"`
		Count int           `json:"count"`
	}
	name := sqlTransfer(ctx.URLParam("name"))
	execStatus := sqlTransfer(ctx.URLParam("exec-status"))

	resList := make([]ResTaskInfo, 0)
	taskList, _ := model.TaskList.GetTaskInfoListByName(name)
	for _, info := range taskList {
		var m ResTaskInfo
		m.Id = info.ID
		m.Name = info.Name
		m.Describe = info.Describe
		if info.CreateTime.Valid == true {
			m.CreateTime = info.CreateTime.Time.Format(base.TsLayout)
		} else {
			m.CreateTime = ""
		}
		m.Status = info.Status
		m.Spec = info.Spec
		//下发主机列表
		err, hostList := model.TaskHostList.GetTaskHostInfoByTaskId(info.ID)
		if err != nil {
			log.Errorf("[Task->TaskList] get task host info err: %v", err)
			return err
		}
		m.Hosts = hostList
		//最近一次执行状态和时间
		query := fmt.Sprintf("SELECT exec_status,end_time FROM %s WHERE operation_id = (SELECT DISTINCT operation_id  FROM task_log "+
			"WHERE task_id = ? ORDER BY id DESC LIMIT 1) ORDER BY end_time DESC", model.TBL_TASK_LOG)
		var taskLogList []model.TaskLog
		if err := model.TaskLogList.GetDB().Select(&taskLogList, query, info.ID); err != nil {
			return err
		}
		if len(taskLogList) > 0 {
			m.ExecStatus, m.EndTime = getExecStatusAndEndTime(taskLogList)
		} else {
			m.ExecStatus = model.NOT_RUNNING
			m.EndTime = ""
		}
		//最近一次手动执行状态
		queryRun := fmt.Sprintf("SELECT exec_status,end_time FROM %s WHERE operation_id = (SELECT DISTINCT operation_id  FROM task_log "+
			"WHERE task_id = ? AND exec_type = ? ORDER BY id DESC LIMIT 1) ORDER BY end_time DESC", model.TBL_TASK_LOG)
		var taskRunLogList []model.TaskLog
		if err := model.TaskLogList.GetDB().Select(&taskRunLogList, queryRun, info.ID, model.TASK_MANUAL_RUN); err != nil {
			return err
		}
		if len(taskRunLogList) > 0 {
			m.RunStatus, _ = getExecStatusAndEndTime(taskRunLogList)
		} else {
			m.RunStatus = model.NOT_RUNNING
		}
		if execStatus == "" {
			resList = append(resList, m)
		} else if execStatus == strconv.Itoa(int(m.ExecStatus)) {
			resList = append(resList, m)
		}
	}
	// 重写排序
	pagination := apibase.GetPaginationFromQueryParameters(nil, ctx, nil)
	switch pagination.SortBy {
	case "create_time":
		sort.SliceStable(resList, func(i, j int) bool {
			if pagination.SortDesc {
				return resList[i].CreateTime > resList[j].CreateTime
			} else {
				return resList[i].CreateTime < resList[j].CreateTime
			}
		})
	case "end_time":
		sort.SliceStable(resList, func(i, j int) bool {
			if pagination.SortDesc {
				return resList[i].EndTime > resList[j].EndTime
			} else {
				return resList[i].EndTime < resList[j].EndTime
			}
		})
	}
	// 重写分页
	total := len(resList)
	if pagination.Start > 0 {
		if pagination.Start+pagination.Limit < total {
			resList = resList[pagination.Start : pagination.Start+pagination.Limit]
		} else if pagination.Start > total {
			resList = nil
		} else {
			resList = resList[pagination.Start:total]
		}
	} else {
		if pagination.Limit == 0 {
			resList = resList[:total]
		} else if pagination.Limit < total {
			resList = resList[:pagination.Limit]
		}
	}

	return res{
		List:  resList,
		Count: total,
	}
}

func getExecStatusAndEndTime(taskLogList []model.TaskLog) (model.Status, string) {
	var runningCount, failureCount, finishCount int
	var failureTime, finishTime []time.Time
	for _, taskLog := range taskLogList {
		switch taskLog.ExecStatus {
		case model.RUNNING:
			runningCount++
		case model.FAILURE:
			failureCount++
			if taskLog.EndTime.Valid {
				failureTime = append(failureTime, taskLog.EndTime.Time)
			}
		case model.FINISH:
			finishCount++
			if taskLog.EndTime.Valid {
				finishTime = append(finishTime, taskLog.EndTime.Time)
			}
		}
	}
	//状态优先级：运行中 > 异常 > 正常
	var execStatus model.Status
	var endTime = ""
	if runningCount != 0 {
		execStatus = model.RUNNING
	} else if failureCount != 0 {
		execStatus = model.FAILURE
		if len(failureTime) > 0 {
			sort.SliceStable(failureTime, func(i, j int) bool {
				return failureTime[i].Unix() > failureTime[j].Unix()
			})
			endTime = failureTime[0].Format(base.TsLayout)
		}
	} else if finishCount != 0 {
		execStatus = model.FINISH
		if len(finishTime) > 0 {
			sort.SliceStable(finishTime, func(i, j int) bool {
				return finishTime[i].Unix() > finishTime[j].Unix()
			})
			endTime = finishTime[0].Format(base.TsLayout)
		}
	}
	return execStatus, endTime
}

func TaskUpload(ctx context.Context) apibase.Result {
	log.Debugf("[Task->TaskUpload] TaskUpload from EasyMatrix API ")

	taskLock.Lock()
	defer taskLock.Unlock()
	file, head, err := ctx.FormFile("file")
	if errors.Is(err, http.ErrMissingFile) {
		return fmt.Errorf("请上传脚本")
	} else if err != nil {
		return err
	}
	defer file.Close()
	describe := ctx.FormValue("describe")
	execTimeoutStr := ctx.FormValue("exec_timeout")
	if execTimeoutStr == "" || execTimeoutStr == "0" {
		return fmt.Errorf("未配置'超时设置'")
	}
	execTimeout, err := strconv.Atoi(execTimeoutStr)
	if err != nil {
		return fmt.Errorf("exec_timeout is not number")
	}
	logRetentionStr := ctx.FormValue("log_retention")
	if logRetentionStr == "" || logRetentionStr == "0" {
		return fmt.Errorf("未配置'执行历史保存周期'")
	}
	logRetention, err := strconv.Atoi(logRetentionStr)
	if err != nil {
		return fmt.Errorf("log_retention is not number")
	}

	//限制文件后缀
	fileExt := path.Ext(head.Filename)
	var FileAllow = map[string]bool{
		pythonExt: true,
		shellExt:  true,
	}
	if _, ok := FileAllow[fileExt]; !ok {
		return fmt.Errorf("仅支持 %s，%s 格式文件", pythonExt, shellExt)
	}
	//上传文件
	var (
		taskPath    = filepath.Join(base.WebRoot, model.TASK_FILES_DIR)
		absTaskFile = filepath.Join(taskPath, head.Filename)
	)
	err, taskId := model.TaskList.InsertTaskIfNotExist(head.Filename, describe, "", execTimeout, logRetention)
	if err != nil {
		log.Errorf("[Task->TaskUpload] insert task err: %v", err)
		return fmt.Errorf("上传失败，%v", err)
	}
	defer func() {
		if err := addSafetyAuditRecord(ctx, "平台管理", "脚本管理", "上传脚本: "+head.Filename+", TaskId: "+strconv.Itoa(taskId)); err != nil {
			log.Errorf("[Task->TaskUpload] failed to add safety audit record\n")
		}
	}()
	if err := os.MkdirAll(taskPath, 0755); err != nil {
		return err
	}
	fw, err := os.Create(absTaskFile)
	if err != nil {
		return fmt.Errorf("failed to create file, err: %v", err)
	}
	if _, err = io.Copy(fw, file); err != nil {
		fw.Close()
		return err
	}
	defer fw.Close()

	return nil
}

func TaskFileContent(ctx context.Context) apibase.Result {
	log.Debugf("[Task->TaskFileContent] TaskFileContent from EasyMatrix API ")

	paramErrs := apibase.NewApiParameterErrors()
	type res struct {
		ScriptContent string `json:"script_content"`
		Describe      string `json:"describe"`
		ExecTimeout   int    `json:"exec_timeout"`
		LogRetention  int    `json:"log_retention"`
	}
	id := ctx.Params().Get("id")
	if id == "" {
		paramErrs.AppendError("$", fmt.Errorf("id is empty"))
	}
	taskId, err := strconv.Atoi(id)
	if err != nil {
		paramErrs.AppendError("$", fmt.Errorf("id is invalid"))
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	info, err := model.TaskList.GetTaskInfoByTaskId(taskId)
	if err != nil {
		log.Errorf("[Task->TaskFileContent] get task info err: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("task id is not exist: %v", taskId)
		}
		return err
	}
	fileName := info.Name
	targetFile := filepath.Join(base.WebRoot, model.TASK_FILES_DIR, fileName)
	if !util.FileIsExist(targetFile) {
		return fmt.Errorf("文件不存在，文件路径: %v", targetFile)
	}
	fi, err := os.Open(targetFile)
	if err != nil {
		log.Errorf("[Task->TaskFileContent] get task file err: %v", err)
		return err
	}
	defer fi.Close()
	content, err := ioutil.ReadAll(fi)
	if err != nil {
		log.Errorf("[Task->TaskFileContent] read task file err: %v", err)
		return err
	}

	return res{
		ScriptContent: string(content[:]),
		Describe:      info.Describe,
		ExecTimeout:   info.ExecTimeout,
		LogRetention:  info.LogRetention,
	}
}

func TaskEdit(ctx context.Context) apibase.Result {
	log.Debugf("[Task->TaskEdit] TaskEdit from EasyMatrix API ")

	param := struct {
		Describe     string `json:"describe"`
		ExecTimeout  int    `json:"exec_timeout"`
		LogRetention int    `json:"log_retention"`
	}{}
	id := ctx.Params().Get("id")
	if id == "" {
		return fmt.Errorf("id is empty")
	}
	taskId, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("id is invalid")
	}
	if err := ctx.ReadJSON(&param); err != nil {
		log.Errorf("parse param error: %v", err)
		return err
	}
	describe := param.Describe
	execTimeout := param.ExecTimeout
	if execTimeout == 0 {
		return fmt.Errorf("未配置'超时设置'")
	}
	logRetention := param.LogRetention
	if logRetention == 0 {
		return fmt.Errorf("未配置'执行历史保存周期'")
	}

	info, err := model.TaskList.GetTaskInfoByTaskId(taskId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Errorf("[Task->TaskEdit] get task info err: %v", err)
		return err
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("task id is not exist: %v", taskId)
	}
	query := "UPDATE " + model.TaskList.TableName + " SET `describe`=?, exec_timeout=?, log_retention=?, update_time=NOW() WHERE id=? AND is_deleted=0"
	if _, err := model.TaskList.GetDB().Exec(query, describe, execTimeout, logRetention, taskId); err != nil {
		log.Errorf("[Task->TaskEdit] update task info err: %v", err)
		return err
	}
	//重新添加任务
	if info.Status == model.TASK_STATUS_ENABLE {
		if err = addTaskToTimer(info.ID); err != nil {
			log.Errorf("[Task->TaskEdit] add task to timer err: %v", err)
			return err
		}
	}

	return nil
}

func ModifyStatus(ctx context.Context) apibase.Result {
	log.Debugf("[Task->ModifyStatus] ModifyStatus from EasyMatrix API ")

	paramErrs := apibase.NewApiParameterErrors()
	var param struct {
		TaskId string `json:"task_id"`
		Status int    `json:"status"`
	}
	if err := ctx.ReadJSON(&param); err != nil {
		return fmt.Errorf("ReadJSON err: %v", err)
	}
	taskIdStr := param.TaskId
	if taskIdStr == "" {
		paramErrs.AppendError("$", fmt.Errorf("task_id is empty"))
	}
	status := param.Status
	if status != model.TASK_STATUS_DISABLE && status != model.TASK_STATUS_ENABLE {
		paramErrs.AppendError("$", fmt.Errorf("status need 0 or 1"))
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	//校验下发主机和cron表达式
	taskIds := strings.Split(taskIdStr, ",")
	taskList := make([]model.TaskInfo, 0)
	for _, id := range taskIds {
		taskId, err := strconv.Atoi(id)
		if err != nil {
			return fmt.Errorf("task_id is invalid")
		}
		taskInfo, err := model.TaskList.GetTaskInfoByTaskId(taskId)
		if err != nil {
			log.Errorf("[Task->ModifyStatus] get task info err: %v", err)
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("task id is not exist: %v", taskId)
			}
			return err
		}
		//打开定时（需校验），关闭定时无需校验
		if status == model.TASK_STATUS_ENABLE {
			err, hostList := model.TaskHostList.GetTaskHostInfoByTaskId(taskInfo.ID)
			if err != nil {
				log.Errorf("[Task->ModifyStatus] get task host info err: %v", err)
				return err
			}
			if taskInfo.Spec == "" || len(hostList) == 0 {
				return fmt.Errorf("开启失败，请先完善定时设置")
			}
			if _, err = checkCrontabExpr(taskInfo.Spec); err != nil {
				return err
			}
			taskInfo.Hosts = hostList
			taskInfo.ExecType = model.TASK_AUTO_RUN
		}
		taskList = append(taskList, taskInfo)
	}
	//批量开启、关闭定时任务
	for _, taskInfo := range taskList {
		if status == model.TASK_STATUS_ENABLE {
			task.ServiceTask.RemoveAndAdd(taskInfo)
		} else {
			task.ServiceTask.Remove(taskInfo.ID)
		}
	}
	if err := model.TaskList.UpdateStatusByTaskIds(taskIds, status); err != nil {
		log.Errorf("[Task->ModifyStatus] update task status err: %v", err)
		return err
	}
	defer func() {
		if err := addSafetyAuditRecord(ctx, "平台管理", "脚本管理", "修改定时状态，TaskId: "+taskIdStr); err != nil {
			log.Errorf("[Task->ModifyStatus] failed to add safety audit record\n")
		}
	}()
	return nil
}

func addTaskToTimer(taskId int) error {
	taskInfo, err := model.TaskList.GetTaskInfoByTaskId(taskId)
	if err != nil {
		return err
	}
	err, hostList := model.TaskHostList.GetTaskHostInfoByTaskId(taskInfo.ID)
	if err != nil {
		return err
	}
	taskInfo.Hosts = hostList
	taskInfo.ExecType = model.TASK_AUTO_RUN
	task.ServiceTask.RemoveAndAdd(taskInfo)
	return nil
}

func TaskUpdate(ctx context.Context) apibase.Result {
	log.Debugf("[Task->TaskUpdate] TaskUpdate from EasyMatrix API ")

	paramErrs := apibase.NewApiParameterErrors()
	var param struct {
		HostId string `json:"host_id"`
		Spec   string `json:"spec"`
	}
	if err := ctx.ReadJSON(&param); err != nil {
		return fmt.Errorf("ReadJSON err: %v", err)
	}
	id := ctx.Params().Get("id")
	if id == "" {
		paramErrs.AppendError("$", fmt.Errorf("id is empty"))
	}
	taskId, err := strconv.Atoi(id)
	if err != nil {
		paramErrs.AppendError("$", fmt.Errorf("id is invalid"))
	}
	cronStr := param.Spec
	if cronStr == "" {
		paramErrs.AppendError("$", fmt.Errorf("spec is empty"))
	}
	if _, err = checkCrontabExpr(cronStr); err != nil {
		return err
	}
	hostIdStr := param.HostId
	if hostIdStr == "" {
		paramErrs.AppendError("$", fmt.Errorf("host_id is empty"))
	}
	hostIdStrList := strings.Split(hostIdStr, ",")
	taskHosts := make([]model.TaskHost, len(hostIdStrList))
	for i, hostIdStr := range hostIdStrList {
		taskHosts[i].TaskId = taskId
		taskHosts[i].HostId, err = strconv.Atoi(hostIdStr)
		if err != nil {
			paramErrs.AppendError("$", fmt.Errorf("host_id is invalid"))
		}
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	info, err := model.TaskList.GetTaskInfoByTaskId(taskId)
	if err != nil {
		log.Errorf("[Task->TaskUpdate] get task info err: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("task id is not exist: %v", taskId)
		}
		return err
	}
	//插入主机记录
	if err := model.TaskHostList.InsertTaskHost(taskId, taskHosts); err != nil {
		log.Errorf("[Task->TaskUpdate] insert task host err: %v", err)
		return err
	}
	//更新cron表达式
	if cronStr != info.Spec {
		if err := model.TaskList.UpdateSpecByTaskId(taskId, cronStr); err != nil {
			log.Errorf("[Task->TaskUpdate] update task spec err: %v", err)
			return err
		}
	}
	//更新定时状态-打开
	if info.Status == model.TASK_STATUS_DISABLE {
		taskIdList := []string{strconv.Itoa(taskId)}
		if err := model.TaskList.UpdateStatusByTaskIds(taskIdList, model.TASK_STATUS_ENABLE); err != nil {
			log.Errorf("[Task->TaskUpdate] update task status err: %v", err)
			return err
		}
	}
	//重新添加任务
	if err = addTaskToTimer(taskId); err != nil {
		log.Errorf("[Task->TaskUpdate] add task to timer err: %v", err)
		return err
	}

	return nil
}

func ParseSpec(ctx context.Context) apibase.Result {
	log.Debugf("[Task->ParseSpec] ParseSpec from EasyMatrix API ")

	type res struct {
		NextTime []string `json:"next_time"`
	}
	var param struct {
		Spec string `json:"spec"`
		Next int    `json:"next"`
	}
	if err := ctx.ReadJSON(&param); err != nil {
		return fmt.Errorf("ReadJSON err: %v", err)
	}
	cronStr := param.Spec
	outTimeCount := param.Next
	if outTimeCount < 1 {
		outTimeCount = 1
	}

	var nextTimeArr []string
	expr, err := checkCrontabExpr(cronStr)
	if err != nil {
		return err
	}
	outTimes := expr.NextN(time.Now(), uint(outTimeCount))
	for _, outTime := range outTimes {
		nextTimeArr = append(nextTimeArr, outTime.Format(base.TsLayout))
	}

	return res{
		NextTime: nextTimeArr,
	}
}

func checkCrontabExpr(cronStr string) (*cronexpr.Expression, error) {
	var build strings.Builder
	cronStr = strings.Trim(cronStr, " ")
	if len(strings.Split(cronStr, " ")) != 6 {
		return nil, fmt.Errorf("cron表达式需为6个空格分隔的字段，官方文档：https://pkg.go.dev/github.com/robfig/cron")
	}
	build.WriteString(cronStr)
	build.WriteString(" *")
	cronStr = build.String()

	expr, err := cronexpr.Parse(cronStr)
	if err != nil {
		return nil, fmt.Errorf("cron表达式填写错误：%v", err)
	}
	return expr, nil
}

func TaskDelete(ctx context.Context) apibase.Result {
	log.Debugf("[Task->TaskDelete] TaskDelete from EasyMatrix API ")

	paramErrs := apibase.NewApiParameterErrors()
	id := ctx.Params().Get("id")
	if id == "" {
		paramErrs.AppendError("$", fmt.Errorf("id is empty"))
	}
	taskId, err := strconv.Atoi(id)
	if err != nil {
		paramErrs.AppendError("$", fmt.Errorf("id is invalid"))
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	if err := model.TaskList.DeleteTaskByTaskId(taskId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("task id is not exist: %v", taskId)
		}
		log.Errorf("[Task->TaskDelete] delete task err: %v", err)
		return err
	}
	if err := model.TaskHostList.DeleteTaskHostByTaskId(taskId); err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Errorf("[Task->TaskDelete] delete task host err: %v", err)
		return err
	}
	task.ServiceTask.Remove(taskId)
	defer func() {
		if err := addSafetyAuditRecord(ctx, "平台管理", "脚本管理", "删除脚本, TaskId: "+id); err != nil {
			log.Errorf("[Task->TaskDelete] failed to add safety audit record\n")
		}
	}()

	return nil
}

func TaskRun(ctx context.Context) apibase.Result {
	log.Debugf("[Task->TaskRun] TaskRun from EasyMatrix API ")

	paramErrs := apibase.NewApiParameterErrors()
	var param struct {
		HostId string `json:"host_id"`
	}
	if err := ctx.ReadJSON(&param); err != nil {
		return fmt.Errorf("ReadJSON err: %v", err)
	}
	id := ctx.Params().Get("id")
	if id == "" {
		paramErrs.AppendError("$", fmt.Errorf("id is empty"))
	}
	taskId, err := strconv.Atoi(id)
	if err != nil {
		paramErrs.AppendError("$", fmt.Errorf("id is invalid"))
	}
	hostIdStr := param.HostId
	if hostIdStr == "" {
		paramErrs.AppendError("$", fmt.Errorf("host_id is empty"))
	}
	hostIdStrList := strings.Split(hostIdStr, ",")
	hostIds := make([]int, len(hostIdStrList))
	for i, hostIdStr := range hostIdStrList {
		hostIds[i], err = strconv.Atoi(hostIdStr)
		if err != nil {
			paramErrs.AppendError("$", fmt.Errorf("host_id is invalid"))
		}
	}
	paramErrs.CheckAndThrowApiParameterErrors()

	info, err := model.TaskList.GetTaskInfoByTaskId(taskId)
	if err != nil {
		log.Errorf("[Task->TaskRun] get task info err: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("task id is not exist: %v", taskId)
		}
		return err
	}
	err, hostList := model.TaskHostList.GetTaskHostByHostIds(hostIdStrList)
	if err != nil {
		log.Errorf("[Task->TaskRun] get task host info err: %v", err)
		return err
	}
	info.Hosts = hostList
	info.ExecType = model.TASK_MANUAL_RUN
	task.ServiceTask.Run(info)

	return nil
}

func TaskLogs(ctx context.Context) apibase.Result {
	log.Debugf("[Task->TaskLogs] TaskRun from EasyMatrix API ")

	paramErrs := apibase.NewApiParameterErrors()
	type res struct {
		List  []ResTaskLog `json:"list"`
		Count int          `json:"count"`
	}
	pagination := apibase.GetPaginationFromQueryParameters(nil, ctx, model.TaskLog{})
	id := ctx.Params().Get("id")
	if id == "" {
		paramErrs.AppendError("$", fmt.Errorf("id is empty"))
	}
	taskId, err := strconv.Atoi(id)
	if err != nil {
		paramErrs.AppendError("$", fmt.Errorf("id is invalid"))
	}
	execStatus := ctx.URLParam("exec-status")
	paramErrs.CheckAndThrowApiParameterErrors()

	resList := make([]ResTaskLog, 0)
	opList, count := model.TaskLogList.GetOperationIdByTaskId(taskId, execStatus, pagination)
	for _, v := range opList {
		var taskLogList []model.TaskLog
		query := fmt.Sprintf("SELECT ip,exec_type,exec_status,end_time,exec_result FROM %s "+
			"WHERE operation_id = ? ORDER BY ip DESC", model.TBL_TASK_LOG)
		if err := model.TaskLogList.GetDB().Select(&taskLogList, query, v.OperationId); err != nil {
			return err
		}
		r := ResTaskLog{}
		r.ExecStatus, r.EndTime = getExecStatusAndEndTime(taskLogList)
		for _, item := range taskLogList {
			c := LogChildren{}
			if item.EndTime.Valid {
				c.EndTime = item.EndTime.Time.Format(base.TsLayout)
			} else {
				c.EndTime = ""
			}
			c.Ip = item.Ip
			c.ExecType = item.ExecType
			c.ExecStatus = item.ExecStatus
			c.ExecResult = item.ExecResult
			r.ExecType = c.ExecType
			if execStatus == "" {
				r.Children = append(r.Children, c)
			} else if execStatus == strconv.Itoa(int(c.ExecStatus)) {
				r.Children = append(r.Children, c)
			}

		}
		resList = append(resList, r)
	}

	return res{
		List:  resList,
		Count: count,
	}
}
