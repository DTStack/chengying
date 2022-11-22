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
	"dtstack.com/dtstack/easymatrix/matrix/api/k8s/view"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	modelkube "dtstack.com/dtstack/easymatrix/matrix/model/kube"
	"strconv"
)
func NamsapceEvent(clusterid,namespace string, limit,start int) (*view.NamespaceEventRsp,error){
	cid,err := strconv.Atoi(clusterid)
	if err != nil{
		log.Errorf("[namespace_event]: convert %s to int error: %v",clusterid,err)
		return nil,err
	}
	tbsc,err := modelkube.DeployNamespaceList.Get(namespace,cid)
	if err != nil{
		return  nil,err
	}
	if tbsc == nil{
		return nil,nil
	}
	worksapceid := tbsc.Id
	events,err := namsapceEventItems(worksapceid,limit,start)
	if err!= nil{
		return nil,err
	}
	count,err := modelkube.DeployNamespaceEvent.SelectCount(worksapceid)
	if err != nil{
		return nil,err
	}
	return &view.NamespaceEventRsp{
		Size:   count,
		Events: events,
	},nil
}
func namsapceEventItems(workspaceid,limit,start int) ([]view.NamespaceEvent,error){

	tbscs, err := modelkube.DeployNamespaceEvent.PageSelect(workspaceid,start,limit)
	if err != nil{
		return nil,err
	}
	events := make([]view.NamespaceEvent,0,len(tbscs))
	for _,tbsc := range tbscs{
		e := view.NamespaceEvent{
			Id: 	  tbsc.Id,
			Time:     tbsc.Time.Format("2006-01-02 15:04:05"),
			Type:     tbsc.Type,
			Reason:   tbsc.Reason,
			Resource: tbsc.Resource,
			Message:  tbsc.Message,
		}
		events = append(events,e)
	}
	return events,nil

}
