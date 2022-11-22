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

package strategy

import (
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	RefreshDuration = 30
)

const (
	STANDBY_SERVICE_NAME  = "{SERVICENAME}"
	STANDBY_PRODUCT_NAME  = "{PRODUCTNAME}"
	STANDBY_INSTANCE_ID   = "{INSTANCEID}"
	STANDBY_AGENT_ID      = "{AGENTID}"
	STANDBY_SIDECAR_ID    = "{SIDDECARD}"
	STANDBY_HOST_IP       = "{HOSTIP}"
	STANDBY_TASK_INTERVAL = "{INTERVAL}"
)

func createHostTask(s model.StrategyInfo) {
	log.Infof("create host task: %v", s.Name)
	as := StrategyCacheMap.GetAssign()
	rs := StrategyCacheMap.GetResource()
	infos, _ := model.DeployHostList.GetHostList(nil)
	ipToHostMap := make(map[string][]model.HostInfo)
	var healthHost []model.HostInfo
	for _, h := range infos {
		if time.Now().Sub(time.Time(h.UpdateDate)) >= 3*time.Minute {
			continue
		}
		if h.Status < 0 {
			continue
		}
		ipToHostMap[h.Ip] = append(ipToHostMap[h.Ip], h)
		healthHost = append(healthHost, h)
	}
	list, count := model.DeployProductList.GetProductList("", "", nil, nil)
	if count > 0 && len(list) > 0 {
		for _, pro := range list {
			if pro.Status == model.PRODUCT_STATUS_DEPLOYING || pro.Status == model.PRODUCT_STATUS_UNDEPLOYING {
				log.Infof("stop host task %v since product %v is deploying", s.Name, pro.ProductName)
				return
			}
		}
	}
	var targets []model.HostInfo
	if _, ok := as[s.ID]; !ok {
		for _, host := range healthHost {
			targets = append(targets, host)
		}
	} else {
		for _, a := range as[s.ID] {
			for _, ip := range strings.Split(a.Host, ",") {
				if _, ok := ipToHostMap[ip]; ok {
					targets = append(targets, ipToHostMap[ip]...)
				}
			}
		}
	}
	script := rs[s.ID].Content
	script = strings.Replace(script, STANDBY_TASK_INTERVAL, strconv.Itoa(s.CronInterval), -1)
	for _, h := range targets {
		if _, ok := rs[s.ID]; !ok {
			continue
		}
		task := TaskInfo{
			TaskType:     STR_TYPE_HOST,
			StrategyId:   s.ID,
			StrategyName: s.Name,
			SidecarId:    h.SidecarId,
			AgentId:      "",
			ExecScript:   script,
			Timeout:      strconv.Itoa(s.TimeOut) + "s",
			Parameter:    s.Params,
			Host:         h.Ip,
		}
		StrategyTaskMap.Add(task)
	}
}

func createServcieTask(s model.StrategyInfo) {
	log.Infof("create server task: %v", s.Name)
	as := StrategyCacheMap.GetAssign()
	rs := StrategyCacheMap.GetResource()
	infos, _ := model.DeployHostList.GetHostList(nil)
	ipToHostMap := make(map[string][]model.HostInfo)
	for _, h := range infos {
		if time.Now().Sub(time.Time(h.UpdateDate)) >= 3*time.Minute {
			continue
		}
		if h.Status < 0 {
			continue
		}
		ipToHostMap[h.Ip] = append(ipToHostMap[h.Ip], h)
	}
	instances, err := model.DeployInstanceList.GetInstanceList()
	if err != nil {
		log.Errorf("err :%v", err.Error())
		return
	}
	var targets []model.DeployInstanceInfo
	if _, ok := as[s.ID]; !ok {
		targets = instances
	} else {
		//TODO Later
		//策略关联到特定产品、服务
	}
	pidToPrroductInfo := make(map[int]model.DeployProductListInfo)
	list, _ := model.DeployProductList.GetProductList("", "", nil, nil)
	if len(list) > 0 {
		for _, pro := range list {
			//只要有产品在部署，停止策略脚本下发
			if pro.Status == model.PRODUCT_STATUS_DEPLOYING || pro.Status == model.PRODUCT_STATUS_UNDEPLOYING {
				log.Infof("stop service task %v since product %v is deploying/undeploying", s.Name, pro.ProductName)
				return
			}
			log.Infof("pid: %v", pro.ID)
			pidToPrroductInfo[pro.ID] = pro
		}
	}
	for _, inst := range targets {
		//屏蔽缓存失效
		if _, ok := rs[s.ID]; !ok {
			continue
		}
		//屏蔽非合法状态产品，部署中的和卸载中的
		if _, ok := pidToPrroductInfo[inst.Pid]; !ok {
			log.Infof("stop service %v-%v task: %v since productInfo is null", inst.ServiceName, inst.Ip, s.Name)
			continue
		}
		//屏蔽不在线主机
		if _, ok := ipToHostMap[inst.Ip]; !ok {
			log.Infof("stop service %v-%v task: %v, host is not available", inst.ServiceName, inst.Ip, s.Name)
			continue
		}
		product := pidToPrroductInfo[inst.Pid]
		//屏蔽部署失败的产品
		if product.Status != model.PRODUCT_STATUS_DEPLOYED {
			log.Infof("stop service %v-%v task: %v, product is not deployed", inst.ServiceName, inst.Ip, s.Name)
			continue
		}
		script := rs[s.ID].Content
		script = strings.Replace(script, STANDBY_PRODUCT_NAME, product.ProductName, -1)
		script = strings.Replace(script, STANDBY_SERVICE_NAME, inst.ServiceName, -1)
		script = strings.Replace(script, STANDBY_AGENT_ID, inst.AgentId, -1)
		script = strings.Replace(script, STANDBY_SIDECAR_ID, inst.Sid, -1)
		script = strings.Replace(script, STANDBY_INSTANCE_ID, strconv.Itoa(inst.ID), -1)
		script = strings.Replace(script, STANDBY_HOST_IP, inst.Ip, -1)
		script = strings.Replace(script, STANDBY_HOST_IP, inst.Ip, -1)
		script = strings.Replace(script, STANDBY_TASK_INTERVAL, strconv.Itoa(s.CronInterval), -1)

		task := TaskInfo{
			TaskType:          STR_TYPE_SERVICE,
			StrategyId:        s.ID,
			StrategyName:      s.Name,
			SidecarId:         inst.Sid,
			AgentId:           inst.AgentId,
			ExecScript:        script,
			Timeout:           strconv.Itoa(s.TimeOut) + "s",
			Parameter:         s.Params,
			Host:              inst.Ip,
			ServiceName:       inst.ServiceName,
			ProductName:       product.ProductName,
			ParentProductName: product.ParentProductName,
		}
		StrategyTaskMap.Add(task)
		//避免任务大量并发
		time.Sleep(1 * time.Second)
	}

}

func scheduleFunc(id int) {
	s := ScheduleCache.GetStategyCache(id)
	if len(s.Name) == 0 {
		return
	}
	log.Infof("do schedule, strategy : %v", s)
	switch s.Property {
	case STR_TYPE_HOST:
		createHostTask(s)
	case STR_TYPE_SERVICE:
		createServcieTask(s)
	default:
		log.Errorf("wrong strategy property: %v", s.Property)
	}
}

type safeScheduleCache struct {
	sync.RWMutex
	s map[int]model.StrategyInfo
}

var (
	ScheduleCache = &safeScheduleCache{
		s: make(map[int]model.StrategyInfo),
	}
)

func (this *safeScheduleCache) Refresh() {
	duration := time.Duration(RefreshDuration) * time.Second
	DefaultScheduler.Start()
	for {
		log.Infof("schedule refresh...")
		for key, _ := range this.s {
			if _, ok := StrategyCacheMap.GetStrategy()[key]; !ok {
				DefaultScheduler.Remove(strconv.Itoa(key), scheduleFunc)
			}
		}
		for key, value := range StrategyCacheMap.GetStrategy() {
			if _, ok := this.s[key]; !ok {
				this.addCron(*value)
			} else {
				if value.GmtModify != this.s[key].GmtModify {
					this.refreshCron(*value)
				}
			}
		}
		time.Sleep(duration)
	}
}

func (this *safeScheduleCache) addCron(value model.StrategyInfo) {
	this.Lock()
	defer this.Unlock()
	log.Infof("addcron, strategy id: %v", value.ID)
	this.s[value.ID] = value
	idstr := strconv.Itoa(value.ID)
	switch value.CronPeriod {
	case CRON_TYPE_MINUTES:
		DefaultScheduler.Every(uint64(value.CronInterval*60)).Seconds().Do(idstr, scheduleFunc, value.ID)
	case CRON_TYPE_HOURS:
		DefaultScheduler.Every(uint64(value.CronInterval)).Hours().Do(idstr, scheduleFunc, value.ID)
	case CRON_TYPE_DAYS:
		DefaultScheduler.Every(uint64(value.CronInterval)).Days().Do(idstr, scheduleFunc, value.ID)
	default:
		log.Errorf("wrong schedule period: %v", value.CronPeriod)
	}
}

func (this *safeScheduleCache) removeCron(value model.StrategyInfo) {
	log.Infof("remove cron, strategy id: %v", value.ID)
	DefaultScheduler.Remove(strconv.Itoa(value.ID), scheduleFunc)
}

func (this *safeScheduleCache) refreshCron(value model.StrategyInfo) {
	this.removeCron(value)
	this.addCron(value)
}

func (this *safeScheduleCache) GetStategyCache(id int) model.StrategyInfo {
	this.RLock()
	defer this.RUnlock()
	if _, ok := this.s[id]; ok {
		return this.s[id]
	}
	return model.StrategyInfo{}
}
