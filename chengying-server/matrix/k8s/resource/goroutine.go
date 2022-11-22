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

package resource

import (
	"context"
	easymonitor "dtstack.com/dtstack/easymatrix/addons/easymonitor/pkg"
	monitorevents "dtstack.com/dtstack/easymatrix/addons/easymonitor/pkg/monitor/events"
	"dtstack.com/dtstack/easymatrix/matrix/api/k8s/view"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/constant"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/kube"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/monitor"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	modelkube "dtstack.com/dtstack/easymatrix/matrix/model/kube"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// load the conncet info for namescpcae client talk with api-server
func InitResource() error{

	tbscs,err := modelkube.DeployNamespaceList.SelectAll()
	if err != nil{
		return err
	}
	for _,tbsc := range tbscs{
		cache,err:= kube.ClusterNsClientCache.GetClusterNsClient(strconv.Itoa(tbsc.ClusterId)).GetClientCache(kube.ImportType(tbsc.Type))
		if err != nil{
			return err
		}
		if tbsc.Status == constant.NAMESPACE_NOT_CONNECT || tbsc.Status == constant.NAMESPACE_INVALID{
			continue
		}
		nsSave := tbsc.NamespaceSaveReq
		if tbsc.Type == kube.IMPORT_KUBECONFIG.String(){
			connectInfo,err := modelkube.DeployNamespaceClient.Get(tbsc.Id)
			if err!=nil{
				return err
			}
			if connectInfo == nil{
				continue
			}
			err = cache.Connect(connectInfo.Yaml,tbsc.Namespace)
			if err != nil{
				return err
			}
			nsSave.Yaml = connectInfo.Yaml
		}else{
			ip := tbsc.Ip
			port := tbsc.Port
			if len(ip) ==0 || len(port) ==0{
				continue
			}
			if !strings.HasPrefix(ip,"http"){
				ip = "http://"+ip
			}
			host := ip + ":" +port
			err = cache.Connect(host,tbsc.Namespace)
			if err != nil{
				return err
			}
		}
		err = StartGoroutines(context.Background(),strconv.Itoa(tbsc.ClusterId),&nsSave)
		if err != nil{
			return err
		}
	}

	go regularCleanEvents()
	return nil
}

func StartGoroutines(ctx context.Context, clusterid string, vo *view.NamespaceSaveReq) error{
	cache,err := kube.ClusterNsClientCache.GetClusterNsClient(clusterid).GetClientCache(kube.ImportType(vo.Type))
	if err != nil{
		return err
	}
	cid,_ := strconv.Atoi(clusterid)
	if client := cache.GetClient(vo.Namespace); client != nil{
		go WatchNamspaceStatus(client,vo.Namespace,cid)
	}

	if vo.Type == kube.IMPORT_KUBECONFIG.String(){
		//run monitor
		trasmitor := &KubeconfigModeTransmitor{vo.Id}
		stopCh := make(chan struct{})
		err := easymonitor.StartMonitorControllerWithTransmitor("",vo.Yaml,vo.Namespace,stopCh,trasmitor)
		if err != nil{
			log.Errorf("start monitor with kubeconfig fail, error %v",err)
		}
		go monitorCheck(vo.Namespace,cid,stopCh)
	}else if vo.Type == kube.IMPORT_AGENT.String(){
		if c := cache.GetClient(vo.Namespace); c != nil{
			client := c.(*kube.RestClient)
			go monitorAgent(vo.Namespace,cid,client)
		}
	}

	return nil
}

type KubeconfigModeTransmitor struct {
	NamespaceId int
}

func (k *KubeconfigModeTransmitor) Push(event monitorevents.Eventer){
	e := event.(*monitorevents.Event)
	e.Workspaceid = k.NamespaceId
	monitor.HandleResourceM(e)
}

func (k *KubeconfigModeTransmitor)Process(){

}

func monitorAgent(namespace string, clusterid int,client *kube.RestClient){
	tick := time.NewTicker(time.Second*5)
	defer tick.Stop()
	tbsc,err := modelkube.DeployNamespaceList.Get(namespace,clusterid)
	if err != nil && tbsc == nil{
		return
	}
	updateTime := tbsc.UpdateTime
	for{
		<- tick.C
		tbsc,err = modelkube.DeployNamespaceList.Get(namespace,clusterid)
		if err != nil {
			log.Errorf("[goroutine]: monitorAgent task database err %v",err)
			continue
		}
		if tbsc.UpdateTime != updateTime{
			break
		}
		events := []monitorevents.Event{}
		err := client.Events(&events)
		for _,event := range events{
			event.Workspaceid = tbsc.Id
			monitor.HandleResourceM(&event)
		}
		if err != nil{
			continue
		}

	}
}


func monitorCheck(namespace string, clusterid int, stopCh chan struct{}){
	tick := time.NewTicker(time.Minute)
	defer tick.Stop()
	tbsc,err := modelkube.DeployNamespaceList.Get(namespace,clusterid)
	if err != nil && tbsc == nil{
		return
	}
	updateTime := tbsc.UpdateTime
	for {
		<- tick.C
		tbsc,err = modelkube.DeployNamespaceList.Get(namespace,clusterid)
		if err != nil {
			log.Errorf("[goroutine]: monitorCheck task database err %v",err)
			continue
		}
		if tbsc.UpdateTime != updateTime{
			close(stopCh)
			break
		}
	}
}

func WatchNamspaceStatus(client kube.Client, namespace string, clusterid int){
	tick := time.NewTicker(time.Minute)
	defer tick.Stop()
	defer func() {
		log.Infof("namespace info is updated")
	}()
	var err error
	tbsc, err := modelkube.DeployNamespaceList.Get(namespace,clusterid)
	if err != nil && tbsc == nil{
		return
	}
	updateTime := tbsc.UpdateTime
	for {
		<- tick.C
		tbsc, err = modelkube.DeployNamespaceList.Get(namespace,clusterid)
		if err != nil{
			log.Errorf("[goroutine]: watchNamspaceStatus task database err %v",err)
			continue
		}
		//if updatetime changed, maybe the client info has changed.exit the loop
		if tbsc.UpdateTime != updateTime{
			break
		}
		err := ping(client,namespace)
		if err != nil{
			tbsc.Status = constant.NAMESPACE_INVALID
			modelkube.DeployNamespaceList.UpdateStatus(tbsc)
			continue
		}

		if tbsc.Status == constant.NAMESPACE_VALID{
			continue
		}
		tbsc.Status = constant.NAMESPACE_VALID
		modelkube.DeployNamespaceList.UpdateStatus(tbsc)
	}

}

//clean event history before 14 days ago, every 4:00:00
func regularCleanEvents(){
	regularHour := 3
	now := time.Now()
	waitHour := time.Duration(regularHour - 1 - now.Hour())
	waitMinute := time.Duration(59-now.Minute())
	waitSecond := time.Duration(60-now.Second())
	if waitHour < 0{
		waitHour = waitHour + 24
	}
	duration := time.Hour * waitHour + time.Minute * waitMinute + time.Second * waitSecond
	//regularHour - nowHour
	timer := time.NewTimer(duration)
	for {
		<- timer.C
		fmt.Println("time to clean events",time.Now().Format("2006-01-02 15:04:05"))
		timer.Reset(time.Hour * 24)
		modelkube.DeployNamespaceEvent.CleanHistory(14)
	}
}
