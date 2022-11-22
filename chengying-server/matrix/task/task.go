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

package task

import (
	"dtstack.com/dtstack/easymatrix/matrix/agent"
	"dtstack.com/dtstack/easymatrix/matrix/base"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"fmt"
	"github.com/gorhill/cronexpr"
	"github.com/jakecoffman/cron"
	"github.com/ouqiang/goutil"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	ServiceTask Task
)

type Task struct{}

var (
	//定时任务调度管理器
	serviceCron *cron.Cron

	// 任务计数-正在运行的任务
	taskCount TaskCount

	// 并发队列, 限制同时运行的任务数量
	concurrencyQueue ConcurrencyQueue
)

type TaskResult struct {
	taskLogId int64
	HostIp    string
	Result    string
	Error     error
}

//TaskCount 任务计数
type TaskCount struct {
	wg   sync.WaitGroup
	exit chan struct{}
}

func (tc *TaskCount) Add() {
	tc.wg.Add(1)
}

func (tc *TaskCount) Done() {
	tc.wg.Done()
}

func (tc *TaskCount) Wait() {
	tc.Add()
	tc.wg.Wait()
	close(tc.exit)
}

//ConcurrencyQueue 并发队列
type ConcurrencyQueue struct {
	queue chan struct{}
}

func (cq *ConcurrencyQueue) Add() {
	cq.queue <- struct{}{}
}

func (cq *ConcurrencyQueue) Done() {
	<-cq.queue
}

//Initialize 初始化任务
func (task Task) Initialize() {
	serviceCron = cron.New()
	serviceCron.Start()
	concurrencyQueue = ConcurrencyQueue{queue: make(chan struct{}, 500)}
	taskCount = TaskCount{sync.WaitGroup{}, make(chan struct{})}
	go taskCount.Wait()

	log.Debugf("[Task] 开始初始化定时任务...")
	taskNum := 0
	taskList, err := model.TaskList.GetTaskInfoList()
	if err != nil {
		log.Errorf("[Task] 定时任务初始化，获取任务列表错误: %s", err)
	}
	task.clearTaskLogs()
	for _, taskInfo := range taskList {
		if _, err = checkCrontabExpr(taskInfo.Spec); err != nil {
			log.Errorf("[Task] 定时任务初始化，%s，cron表达式填写错误：%v", taskInfo.Name, err)
			continue
		}
		err, hostList := model.TaskHostList.GetTaskHostInfoByTaskId(taskInfo.ID)
		if err != nil {
			log.Errorf("[Task] 定时任务初始化，%s，获取任务主机列表错误: %s", taskInfo.Name, err)
		}
		if taskInfo.Status == model.TASK_STATUS_DISABLE || len(hostList) == 0 {
			continue
		}
		taskInfo.Hosts = hostList
		taskInfo.ExecType = model.TASK_AUTO_RUN
		task.Add(taskInfo)
		log.Debugf("[Task] 定时任务初始化，%s，已添加到调度器", taskInfo.Name)
		taskNum++
	}
	log.Debugf("[Task] 定时任务初始化完成, 共%d个定时任务添加到调度器", taskNum)
}

//Run 手动运行任务
func (task Task) Run(taskModel model.TaskInfo) {
	go createJob(taskModel)()
}

//Remove 删除任务
func (task Task) Remove(taskId int) {
	serviceCron.RemoveJob(strconv.Itoa(taskId))
}

//RemoveAndAdd 删除任务后添加
func (task Task) RemoveAndAdd(taskModel model.TaskInfo) {
	task.Remove(taskModel.ID)
	task.Add(taskModel)
}

//Add 添加任务
func (task Task) Add(taskModel model.TaskInfo) {
	taskFunc := createJob(taskModel)
	if taskFunc == nil {
		log.Errorf("[Task] 创建任务处理Job失败")
		return
	}
	cronName := strconv.Itoa(taskModel.ID)
	err := goutil.PanicToError(func() {
		serviceCron.AddFunc(taskModel.Spec, taskFunc, cronName)
	})
	if err != nil {
		log.Errorf("[Task] 添加任务到调度器失败,err: %v", err)
	}
}

//createJob 创建任务
func createJob(taskModel model.TaskInfo) cron.FuncJob {
	taskFunc := func() {
		taskCount.Add()
		defer taskCount.Done()
		concurrencyQueue.Add()
		defer concurrencyQueue.Done()
		log.Debugf("[Task] 开始执行任务%s", taskModel.Name)
		taskResultList := execJob(taskModel)
		afterExecJob(taskResultList)
	}
	return taskFunc
}

//beforeExecJob 任务前置操作
func beforeExecJob(taskLogModel model.TaskLog) (taskLogId int64) {
	taskLogId, err := createTaskLog(taskLogModel)
	if err != nil {
		log.Errorf("[Task] create task log fail,err: %v", err)
		return
	}
	log.Debugf("[Task] 任务命令-%s", taskLogModel.Command)

	return taskLogId
}

//afterExecJob 任务执行后置操作
func afterExecJob(taskResultList []TaskResult) {
	for _, taskResult := range taskResultList {
		err := updateTaskLog(taskResult)
		if err != nil {
			log.Errorf("[Task] 任务结束#更新任务日志失败,err: %v", err)
		}
	}
}

//createTaskLog 创建任务日志
func createTaskLog(taskLogModel model.TaskLog) (int64, error) {
	insertId, err := model.TaskLogList.InsertTaskLog(taskLogModel)
	return insertId, err
}

//updateTaskLog 更新任务日志
func updateTaskLog(taskResult TaskResult) error {
	var status model.Status
	var result = ""
	taskLogId := taskResult.taskLogId
	if taskResult.Error != nil {
		status = model.FAILURE
		//只显示异常详情
		result = taskResult.Error.Error()
	} else {
		status = model.FINISH
	}
	return model.TaskLogList.UpdateTaskLogById(taskLogId, status, result)
}

//execJob 执行具体任务
func execJob(taskModel model.TaskInfo) []TaskResult {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("[Task] panic: task/task.go:execJob: ", err)
		}
	}()
	var err error
	targetFile := filepath.Join(base.WebRoot, model.TASK_FILES_DIR, taskModel.Name)
	fi, err := os.Open(targetFile)
	defer fi.Close()
	var cmd string
	if err == nil {
		var content []byte
		content, err = ioutil.ReadAll(fi)
		cmd = fmt.Sprintf("#!/bin/sh\n %s", string(content[:]))
	}
	//调用agent执行任务
	taskLog := model.TaskLog{
		TaskId:      taskModel.ID,
		Name:        taskModel.Name,
		Spec:        taskModel.Spec,
		OperationId: uuid.NewV4().String(),
		Command:     cmd,
		ExecType:    taskModel.ExecType,
		ExecStatus:  model.RUNNING,
	}
	resultChan := make(chan TaskResult, len(taskModel.Hosts))
	for _, taskHost := range taskModel.Hosts {
		taskLog.Ip = taskHost.Ip
		taskLog.StartTime.Time = time.Now()
		taskLogId := beforeExecJob(taskLog)
		if taskLogId <= 0 {
			continue
		}
		go func(host model.ResHostInfo) {
			var output, errorMsg string
			if err != nil {
				errorMsg = err.Error()
			} else {
				execTimeout := strconv.Itoa(taskModel.ExecTimeout) + "s"
				output, err = agent.AgentClient.ToExecCmdWithTimeout(host.Sid, "", cmd, execTimeout, "", "")
				if err != nil {
					errorMsg = err.Error()
				}
			}
			outputMessage := fmt.Sprintf("主机: [%s-%s]\n%s\n%s\n\n", host.Ip, host.HostName, errorMsg, output)
			resultChan <- TaskResult{taskLogId: taskLogId, HostIp: host.Ip, Result: outputMessage, Error: err}
		}(taskHost)
	}
	var taskResultList []TaskResult
	for i := 0; i < (len(taskModel.Hosts)); i++ {
		taskResult := <-resultChan
		taskResultList = append(taskResultList, taskResult)
	}
	return taskResultList
}

//checkCrontabExpr 校验cron表达式
func checkCrontabExpr(cronStr string) (*cronexpr.Expression, error) {
	var build strings.Builder
	cronStr = strings.Trim(cronStr, " ")
	if len(strings.Split(cronStr, " ")) != 6 {
		return nil, fmt.Errorf("cron表达式需为6个空格分隔的字段")
	}
	build.WriteString(cronStr)
	build.WriteString(" *")
	cronStr = build.String()

	expr, err := cronexpr.Parse(cronStr)
	if err != nil {
		return nil, err
	}
	return expr, nil
}

//clearTaskLogs 日志清理（每天0点执行）
func (task Task) clearTaskLogs() {
	spec := "0 0 0 * * *"
	cronName := strconv.Itoa(0)
	taskFunc := func() {
		log.Debugf("[Task] 开始执行日志清理任务...")
		taskList, err := model.TaskList.GetTaskInfoList()
		if err != nil {
			log.Errorf("[Task] 执行日志清理任务，获取任务列表错误: %s", err)
		}
		for _, task := range taskList {
			query := fmt.Sprintf("DELETE FROM %s WHERE task_id = ? AND start_time <= DATE_SUB(NOW(), INTERVAL ? DAY)", model.TaskLogList.TableName)
			if _, err := model.USE_MYSQL_DB().Exec(query, task.ID, task.LogRetention); err != nil {
				log.Errorf("[Task] 执行日志清理任务，日志清理错误: %v", err)
			}
		}
		log.Debugf("[Task] 执行日志清理任务完成")
	}
	err := goutil.PanicToError(func() {
		serviceCron.AddFunc(spec, taskFunc, cronName)
	})
	if err != nil {
		log.Errorf("[Task] 添加日志清理任务到调度器失败,err: %v", err)
	}
}
