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

package monitor

import (
	"dtstack.com/dtstack/easymatrix/addons/easymonitor/pkg/monitor/events"
	workloadv1beta1 "dtstack.com/dtstack/easymatrix/addons/operator/pkg/apis/workload/v1beta1"
	"dtstack.com/dtstack/easymatrix/go-common/db-helper"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	modelkube "dtstack.com/dtstack/easymatrix/matrix/model/kube"
	"dtstack.com/dtstack/easymatrix/matrix/util"
	"dtstack.com/dtstack/easymatrix/schema"
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"strconv"
	"strings"
)

var (
	sigc                       chan interface{}
	OUT_OF_MONITOR_STATUS_LIST = []string{model.PRODUCT_STATUS_UNDEPLOYED, model.PRODUCT_STATUS_DEPLOYED}
)

func init() {
	sigc = make(chan interface{}, 1)
}
func isClusterProductStatusOver(product *model.ClusterProductRel, statusList []string) bool {
	for _, status := range statusList {
		if status == product.Status {
			return true
		}
	}
	return false
}

func getClusterKubeLabels(meta *metav1.ObjectMeta) (error, *model.ClusterKubeLabels) {
	var pid, clusterId int
	if _, ok := meta.GetLabels()[events.RESOURCE_LABLE_PRODUCT_ID]; !ok {
		return fmt.Errorf("lable %v not existed!", events.RESOURCE_LABLE_PRODUCT_ID), nil
	}
	pid, err := strconv.Atoi(meta.GetLabels()[events.RESOURCE_LABLE_PRODUCT_ID])
	if err != nil {
		return err, nil
	}
	if _, ok := meta.GetLabels()[events.RESOURCE_LABLE_CLUSTER_ID]; !ok {
		return fmt.Errorf("lable %v not existed!", events.RESOURCE_LABLE_CLUSTER_ID), nil
	}
	clusterId, err = strconv.Atoi(meta.GetLabels()[events.RESOURCE_LABLE_CLUSTER_ID])
	if err != nil {
		return err, nil
	}
	if _, ok := meta.GetLabels()[events.RESOURCE_LABLE_PRODUCT_NAME]; !ok {
		return fmt.Errorf("lable %v not existed!", events.RESOURCE_LABLE_PRODUCT_NAME), nil
	}
	if _, ok := meta.GetLabels()[events.RESOURCE_LABLE_PRODUCT_VERSION]; !ok {
		return fmt.Errorf("lable %v not existed!", events.RESOURCE_LABLE_PRODUCT_VERSION), nil
	}
	if _, ok := meta.GetLabels()[events.RESOURCE_LABLE_SERVICE_NAME]; !ok {
		return fmt.Errorf("lable %v not existed!", events.RESOURCE_LABLE_SERVICE_NAME), nil
	}
	if _, ok := meta.GetLabels()[events.RESOURCE_LABLE_SERVICE_VERSION]; !ok {
		return fmt.Errorf("lable %v not existed!", events.RESOURCE_LABLE_SERVICE_VERSION), nil
	}
	if _, ok := meta.GetLabels()[events.RESOURCE_LABLE_SERVICE_GROUP]; !ok {
		return fmt.Errorf("lable %v not existed!", events.RESOURCE_LABLE_SERVICE_GROUP), nil
	}

	return nil, &model.ClusterKubeLabels{
		DeployUuid:     meta.GetLabels()[events.RESOURCE_LABLE_DEPLOY_UUID],
		Pid:            pid,
		ClusterId:      clusterId,
		ProductName:    meta.GetLabels()[events.RESOURCE_LABLE_PRODUCT_NAME],
		ProductVersion: meta.GetLabels()[events.RESOURCE_LABLE_PRODUCT_VERSION],
		ServiceName:    meta.GetLabels()[events.RESOURCE_LABLE_SERVICE_NAME],
		ServiceVersion: meta.GetLabels()[events.RESOURCE_LABLE_SERVICE_VERSION],
		ServiceGroup:   meta.GetLabels()[events.RESOURCE_LABLE_SERVICE_GROUP],
	}
}

func checkPodRunningStatus(pod *v1.Pod) (error, bool, string) {
	var ready bool
	var messages []string
	if len(pod.Status.ContainerStatuses) == 0 {
		return nil, false, "no container status"
	}
	ready = true
	for _, condition := range pod.Status.Conditions {
		if string(condition.Status) == "False" {
			ready = false
			messages = append(messages, string(condition.Type))
			messages = append(messages, "\n")
			messages = append(messages, condition.Message)
		}
	}
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if !containerStatus.Ready {
			ready = false
			if containerStatus.LastTerminationState.Terminated != nil {
				last := containerStatus.LastTerminationState.Terminated.Message
				messages = append(messages, last)
			}
		}
	}
	return nil, ready, strings.Join(messages, "\n")
}
func triggerInstanceListUpdate(pod *v1.Pod, labels *model.ClusterKubeLabels) error {
	var healthState int
	var progress uint
	var status string
	var message string
	switch pod.Status.Phase {
	case v1.PodPending:
		status = model.INSTANCE_STATUS_INSTALLING
		healthState = model.INSTANCE_HEALTH_WAITING
		progress = 30
	case v1.PodRunning:
		_, ready, msg := checkPodRunningStatus(pod)
		if ready {
			status = model.INSTANCE_STATUS_RUNNING
			healthState = model.INSTANCE_HEALTH_OK
			progress = 100
			message = msg
		} else {
			status = model.INSTANCE_STATUS_INSTALLING
			healthState = model.INSTANCE_HEALTH_WAITING
			progress = 70
			message = msg
		}
	case v1.PodSucceeded:
		status = model.INSTANCE_STATUS_RUNNING
		healthState = model.INSTANCE_HEALTH_OK
		progress = 100
	case v1.PodFailed:
		status = model.INSTANCE_STATUS_RUN_FAIL
		healthState = model.INSTANCE_HEALTH_BAD
	default:
		status = model.INSTANCE_STATUS_INSTALL_FAIL
		healthState = model.INSTANCE_HEALTH_BAD
	}
	hostIp := pod.Status.HostIP
	err, hostInfo := model.DeployNodeList.GetNodeInfoByNodeIp(hostIp)
	if err != nil {
		log.Errorf("%v, nodeIp: %v", err.Error(), hostIp)
	}
	productRel, err := model.DeployClusterProductRel.GetByPidAndClusterIdNamespacce(labels.Pid, labels.ClusterId, pod.GetNamespace())
	if err != nil {
		log.Errorf("%v, productId: %v, clusterId: %v, ns: %v", err.Error(), labels.Pid, labels.ClusterId, pod.GetNamespace())
		return err
	}

	sid := hostInfo.SidecarId
	if len(sid) == 0 {
		sid = string(pod.GetUID())
	}
	group := labels.ServiceGroup
	sc, err := schema.Unmarshal(productRel.ProductParsed)
	if err != nil {
		log.Errorf("schema unmarshal err:%v", err)
		return err
	}
	schemaByte, err := json.Marshal(sc.Service[labels.ServiceName])
	ip := hostIp + "/" + pod.Status.PodIP

	err, id, _ := model.DeployInstanceList.NewPodInstanceRecord(
		labels.ClusterId,
		labels.Pid,
		0,
		healthState,
		pod.GetNamespace(),
		ip,
		sid,
		group,
		labels.ServiceName,
		string(pod.GetUID()),
		labels.ServiceVersion,
		status,
		message,
		schemaByte,
	)
	if err != nil {
		log.Errorf("%v", err.Error())
		return err
	}
	// 当前不在部署过程中，不更新状态
	if isClusterProductStatusOver(&productRel, OUT_OF_MONITOR_STATUS_LIST) {
		log.Infof("old pod status changed, pod key: %v", pod.SelfLink)
		return nil
	}
	deployUuid, err := uuid.FromString(labels.DeployUuid)
	if err != nil {
		log.Errorf("%v", err.Error())
		return err
	}
	//update deploy  instance record
	instanceRecordInfo := &model.DeployInstanceRecordInfo{
		DeployUUID:         deployUuid,
		InstanceId:         int(id),
		Sid:                sid,
		Ip:                 ip,
		ProductName:        labels.ProductName,
		ProductVersion:     labels.ProductVersion,
		ProductNameDisplay: labels.ProductName,
		Group:              labels.ServiceGroup,
		ServiceName:        labels.ServiceName,
		ServiceVersion:     labels.ServiceVersion,
		ServiceNameDisplay: labels.ServiceName,
		Status:             status,
		StatusMessage:      message,
		Progress:           progress,
	}
	err, _, _ = model.DeployInstanceRecord.CreateOrUpdate(instanceRecordInfo)
	if err != nil {
		log.Errorf("%v", err.Error())
		return err
	}
	return nil
}

func triggerInstanceListDelete(key string) error {
	pod, err := model.DeployKubePodList.GetByPodKey(key)
	if err != nil {
		log.Errorf("%v, key: %v", err.Error(), key)
		return err
	}
	agentId := string(pod.PodId)
	err, instance := model.DeployInstanceList.GetInstanceInfoByAgentId(agentId)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	// 确保唯一
	record, err := model.DeployInstanceRecord.GetDeployInstanceRecordByInstanceIdAndStatus(instance.ID, model.INSTANCE_STATUS_UNINSTALLING)
	if err == nil {
		updateQuery := "UPDATE " + model.DeployInstanceRecord.TableName + " SET `status`=?, status_message=?, progress=?, update_time=NOW() WHERE id=?"
		if _, err := model.DeployInstanceRecord.GetDB().Exec(updateQuery, model.INSTANCE_STATUS_UNINSTALLED, "", 100, record.ID); err != nil {
			log.Errorf("%v", err)
			return err
		}
	}

	// delete instance list
	err = model.DeployInstanceList.DeleteByagentId(agentId)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	// delete service ip list
	err = model.DeployServiceIpList.Delete(pod.NameSpace, record.ProductName, record.ServiceName, instance.ClusterId)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	productRel, err := model.DeployClusterProductRel.GetByPidAndClusterIdNamespacce(instance.Pid, instance.ClusterId, instance.Namespace)
	if err != nil {
		log.Errorf("clusterId: %v, productId: %v, err %v", instance.ClusterId, instance.Pid, err.Error())
		return err
	}
	// check undeploying progress
	if productRel.Status != model.PRODUCT_STATUS_UNDEPLOYING {
		return nil
	}
	whereCause := dbhelper.WhereCause{}
	whereCause = whereCause.Equal("pid", instance.Pid)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("cluster_id", pod.ClusterId)
	whereCause = whereCause.And()
	whereCause = whereCause.Equal("namespace", pod.NameSpace)
	list, err := model.DeployInstanceList.GetInstanceListByWhere(whereCause)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	log.Infof("delete deploy uuid %v, len list %v", record.DeployUUID, len(list))
	if len(list) == 0 {
		//关联服务delete完毕，卸载成功
		query := "UPDATE " + model.DeployProductHistory.TableName + " SET `status`=?, deploy_end_time=NOW() WHERE deploy_uuid=? and status=?"
		if _, err := model.DeployProductHistory.GetDB().Exec(query, model.PRODUCT_STATUS_UNDEPLOYED, record.DeployUUID, model.PRODUCT_STATUS_UNDEPLOYING); err != nil {
			log.Errorf("%v", err)
		}
		query = "UPDATE " + model.DeployClusterProductRel.TableName + " SET status=?, is_deleted=?, update_time=NOW() WHERE pid=? AND clusterId=? AND namespace=?"
		if _, err := model.DeployClusterProductRel.GetDB().Exec(query, model.PRODUCT_STATUS_UNDEPLOYED, 1, instance.Pid, instance.ClusterId, instance.Namespace); err != nil {
			log.Errorf("%v", err)
		}
	}
	return nil
}

func triggerServiceIpListUpdate(pod *v1.Pod, labels *model.ClusterKubeLabels) error {
	var ipList []string
	if err := model.USE_MYSQL_DB().Select(&ipList, fmt.Sprintf("select host_ip "+
		"from deploy_cluster_kube_pod_list where "+
		"clusterId=%v and namespace='%v' and product_name='%v' and service_name='%v' and is_deleted=0", labels.ClusterId, pod.GetNamespace(), labels.ProductName, labels.ServiceName)); err != nil {
		return err
	}
	uniqList := util.ArrayUniqueNotNullStr(ipList)
	if err := model.DeployServiceIpList.SetPodServiceIp(pod.GetNamespace(), labels.ProductName, labels.ServiceName, strings.Join(uniqList, ","), labels.ClusterId); err != nil {
		return err
	}
	return nil
}

func handlePodE(event *events.Event) error {
	obj, err := json.Marshal(event.Object)
	if err != nil {
		log.Errorf("%v", err.Error())
		return err
	}
	pod := &v1.Pod{}
	err = json.Unmarshal(obj, pod)
	if err != nil {
		log.Errorf("get event pod err %v", err.Error())
		return err
	}
	err, labels := getClusterKubeLabels(&pod.ObjectMeta)
	if err != nil {
		log.Errorf("%v", err.Error())
		return err
	}
	switch event.Operation {
	case events.OPERATION_CREATE_OR_UPDATE:
		log.Infof("updating pod/instance/service ip list, id: %v", pod.GetSelfLink())
		err, _ := model.DeployKubePodList.UpdateOrCreate(event.Key, pod, labels)
		if err != nil {
			log.Errorf("%v", err.Error())
			return err
		}
		if len(pod.Status.HostIP) == 0 {
			return nil
		}
		err = triggerInstanceListUpdate(pod, labels)
		if err != nil {
			log.Errorf("%v", err.Error())
			return err
		}
		if len(pod.Status.PodIP) == 0 {
			return nil
		}
		err = triggerServiceIpListUpdate(pod, labels)
		if err != nil {
			log.Errorf("%v", err.Error())
			return err
		}
	case events.OPERATION_DELETE:
		log.Infof("deleting pod/instance/service ip list, id: %v", pod.GetSelfLink())
		err = triggerInstanceListDelete(event.Key)
		if err != nil {
			log.Errorf("%v", err.Error())
			return err
		}
		err := model.DeployKubePodList.DeleteByKey(event.Key)
		if err != nil {
			log.Errorf("%v", err.Error())
		}
	default:
	}
	return nil
}

func handleServiceE(event *events.Event) error {
	//err, labels := getClusterKubeLabels(&event.Object.(*v1.Service).ObjectMeta)
	//if err != nil {
	//	return err
	//}
	//switch event.Operation {
	//case events.OPERATION_CREATE_OR_UPDATE:
	//	err, _ := model.DeployKubeServiceList.UpdateOrCreate(event.Object.(*v1.Service), labels)
	//	if err != nil {
	//		return err
	//	}
	//case events.OPERATION_DELETE:
	//	return model.DeployKubeServiceList.Delete(event.Object.(*v1.Service), labels)
	//default:
	//}
	return nil
}

func handleDeploymentE(event *events.Event) error {
	return nil
}

func handleIngressE(event *events.Event) error {
	return nil
}

func handleMoleE(event *events.Event) error {
	obj, err := json.Marshal(event.Object)
	if err != nil {
		log.Errorf("%v", err.Error())
		return err
	}
	mole := &unstructured.Unstructured{}
	err = json.Unmarshal(obj, mole)
	if err != nil {
		log.Errorf("%v", err.Error())
		return err
	}
	labels := mole.GetLabels()
	if len(labels) == 0 {
		log.Infof("found mole without lables, ignore it %v", mole.GetSelfLink())
		return nil
	}
	phase := mole.Object["status"].(map[string]interface{})["phase"]
	switch event.Operation {
	case events.OPERATION_CREATE_OR_UPDATE:
		log.Infof("updating mole id: %v, status: %v", mole.GetSelfLink(), phase)
		var pid, clusterId int
		pid, err = strconv.Atoi(labels[events.RESOURCE_LABLE_PRODUCT_ID])
		if err != nil {
			log.Errorf("%v", err.Error())
			return err
		}
		if _, ok := labels[events.RESOURCE_LABLE_CLUSTER_ID]; !ok {
			return fmt.Errorf("lable %v not existed!", events.RESOURCE_LABLE_CLUSTER_ID)
		}
		clusterId, err = strconv.Atoi(labels[events.RESOURCE_LABLE_CLUSTER_ID])
		if err != nil {
			log.Errorf("%v", err.Error())
			return err
		}
		namespace := mole.GetNamespace()
		rel, err := model.DeployClusterProductRel.GetByPidAndClusterIdNamespacce(pid, clusterId, namespace)
		if err != nil {
			log.Errorf("%v", err.Error())
			return err
		}
		if rel.DeployUUID != labels[events.RESOURCE_LABLE_DEPLOY_UUID] {
			log.Infof("deploy_uuid is out of date %v, now is %v", labels[events.RESOURCE_LABLE_DEPLOY_UUID], rel.DeployUUID)
			return nil
		}
		if rel.Status != model.PRODUCT_STATUS_DEPLOYING {
			//not in deploy process, maybe just concile
			return nil
		}
		var status string
		switch phase {
		case "Pending":
			status = model.PRODUCT_STATUS_DEPLOYING
		case "Running":
			status = model.PRODUCT_STATUS_DEPLOYED
		case "Failed":
			status = model.PRODUCT_STATUS_DEPLOY_FAIL
		}
		err = model.DeployClusterProductRel.UpdateStatusWithNamespace(clusterId, pid, namespace, status)
		if err != nil {
			log.Errorf("%v", err.Error())
		}
		log.Infof("deploy uuid %v", labels[events.RESOURCE_LABLE_DEPLOY_UUID])
		query := "UPDATE " + model.DeployProductHistory.TableName + " SET `status`=?, deploy_end_time=NOW() WHERE deploy_uuid=? AND cluster_id=? AND namespace=?"
		if _, err := model.DeployProductHistory.GetDB().Exec(query, status, labels[events.RESOURCE_LABLE_DEPLOY_UUID], clusterId, namespace); err != nil {
			log.Errorf("%v", err)
		}
		return err
	case events.OPERATION_DELETE:
		log.Infof("deleting mole id: %v, status: %v", mole.GetSelfLink(), phase)
		var pid, clusterId int
		pid, err = strconv.Atoi(labels[events.RESOURCE_LABLE_PRODUCT_ID])
		if err != nil {
			log.Errorf("%v", err.Error())
			return err
		}
		if _, ok := labels[events.RESOURCE_LABLE_CLUSTER_ID]; !ok {
			return fmt.Errorf("lable %v not existed!", events.RESOURCE_LABLE_CLUSTER_ID)
		}
		clusterId, err = strconv.Atoi(labels[events.RESOURCE_LABLE_CLUSTER_ID])
		if err != nil {
			log.Errorf("%v", err.Error())
			return err
		}
		namespace := mole.GetNamespace()
		rel, err := model.DeployClusterProductRel.GetByPidAndClusterIdNamespacce(pid, clusterId, namespace)
		if err != nil {
			log.Errorf("%v", err.Error())
			return err
		}
		if rel.Status != model.PRODUCT_STATUS_UNDEPLOYING {
			return nil
		}
		whereCause := dbhelper.WhereCause{}
		whereCause = whereCause.Equal("pid", pid)
		whereCause = whereCause.And()
		whereCause = whereCause.Equal("cluster_id", clusterId)
		whereCause = whereCause.And()
		whereCause = whereCause.Equal("namespace", namespace)
		list, _ := model.DeployInstanceList.GetInstanceListByWhere(whereCause)
		if len(list) == 0 {
			query := "UPDATE " + model.DeployProductHistory.TableName + " SET `status`=?, deploy_end_time=NOW() WHERE deploy_uuid=? and status=?"
			if _, err := model.DeployProductHistory.GetDB().Exec(query, model.PRODUCT_STATUS_UNDEPLOYED, labels[events.RESOURCE_LABLE_DEPLOY_UUID], model.PRODUCT_STATUS_UNDEPLOYING); err != nil {
				log.Errorf("%v", err)
			}
			query = "UPDATE " + model.DeployClusterProductRel.TableName + " SET status=?, is_deleted=?, update_time=NOW() WHERE pid=? AND clusterId=? AND namespace=?"
			if _, err := model.DeployClusterProductRel.GetDB().Exec(query, model.PRODUCT_STATUS_UNDEPLOYED, 1, pid, clusterId, namespace); err != nil {
				log.Errorf("%v", err)
			}
		}
		return nil
	}
	return nil
}
func handleEvent(event *events.Event) error{
	bts,err := json.Marshal(event.Object)
	if err != nil{
		log.Errorf("[resource_monitor]: mashal obj err %v",err)
		return err
	}
	e := &v1.Event{}
	err = json.Unmarshal(bts,e)
	if err != nil{
		log.Errorf("[resource_monitor]: unmashal obj to corev1 event error %v",err)
		return err
	}
	resource := e.InvolvedObject.Kind+"/"+e.InvolvedObject.Name
	if err != nil{
		log.Errorf("[resouce_monitor]: convert string %s to int error %v",event.Workspaceid,err )
		return err
	}

	saveEvent := &modelkube.DeployNamespaceEventSchema{
		Type:     e.Type,
		Reason:   e.Reason,
		Resource: resource,
		Message: e.Message,
		NamespaceId: event.Workspaceid,
		Time:     e.LastTimestamp.Local(),
	}
	return modelkube.DeployNamespaceEvent.Insert(saveEvent)
}

func handleWorkload(event *events.Event) error{

	obj, err := json.Marshal(event.Object)
	if err != nil {
		log.Errorf("[resource_monitor]: mashal obj err %v",err)
		return err
	}
	process := &workloadv1beta1.WorkloadProcess{}
	err = json.Unmarshal(obj, process)
	if err != nil {
		log.Errorf("[resource_monitor]: unmashal obj to workload error %v",err)
		return err
	}
	phase := process.Status.Phase
	fmt.Println("=====-----+++",phase)
	switch event.Operation {
	case events.OPERATION_CREATE_OR_UPDATE:
		var pid, clusterId int
		pid = process.Spec.ProductId
		clusterId = process.Spec.ClusterId
		namespace := process.Namespace
		rel, err := model.DeployClusterProductRel.GetByPidAndClusterIdNamespacce(pid, clusterId, namespace)
		if err != nil {
			log.Errorf("get from table deploy_cluster_product_rel fail, error %v", err)
			return err
		}
		if rel.DeployUUID != process.Spec.DeployUUId {
			log.Infof("deploy_uuid is out of date %v, now is %v", process.Spec.DeployUUId, rel.DeployUUID)
			return nil
		}
		if rel.Status != model.PRODUCT_STATUS_DEPLOYING {
			//not in deploy process, maybe just concile
			return nil
		}
		var status string
		switch phase {
		case "":
			status = model.PRODUCT_STATUS_DEPLOYING
		case workloadv1beta1.ProcessPending:
			status = model.PRODUCT_STATUS_DEPLOYING
		case workloadv1beta1.ProcessFinish:
			status = model.PRODUCT_STATUS_DEPLOYED
		case workloadv1beta1.WorkloadFail:
			status = model.PRODUCT_STATUS_DEPLOY_FAIL
		}
		err = model.DeployClusterProductRel.UpdateStatusWithNamespace(clusterId, pid, namespace, status)
		if err != nil {
			log.Errorf("update deploy_cluster_product_rel status with namespace fail, error %v", err)
		}
		log.Infof("deploy uuid %v", process.Spec.DeployUUId)
		query := "UPDATE " + model.DeployProductHistory.TableName + " SET `status`=?, deploy_end_time=NOW() WHERE deploy_uuid=? AND cluster_id=? AND namespace=?"
		if _, err := model.DeployProductHistory.GetDB().Exec(query, status, process.Spec.DeployUUId, clusterId, namespace); err != nil {
			log.Errorf("%v", err)
		}
		return err
	case events.OPERATION_DELETE:
		var pid, clusterId int
		pid = process.Spec.ProductId
		clusterId = process.Spec.ClusterId
		namespace := process.Namespace
		rel, err := model.DeployClusterProductRel.GetByPidAndClusterIdNamespacce(pid, clusterId, namespace)
		if err != nil {
			log.Errorf("get from table deploy_cluster_product_rel fail, error %v", err)
			return err
		}
		if rel.Status != model.PRODUCT_STATUS_UNDEPLOYING {
			return nil
		}
		whereCause := dbhelper.WhereCause{}
		whereCause = whereCause.Equal("pid", pid)
		whereCause = whereCause.And()
		whereCause = whereCause.Equal("cluster_id", clusterId)
		whereCause = whereCause.And()
		whereCause = whereCause.Equal("namespace", namespace)
		list, _ := model.DeployInstanceList.GetInstanceListByWhere(whereCause)
		if len(list) == 0 {
			query := "UPDATE " + model.DeployProductHistory.TableName + " SET `status`=?, deploy_end_time=NOW() WHERE deploy_uuid=? and status=?"
			if _, err := model.DeployProductHistory.GetDB().Exec(query, model.PRODUCT_STATUS_UNDEPLOYED, process.Spec.DeployUUId, model.PRODUCT_STATUS_UNDEPLOYING); err != nil {
				log.Errorf("%v", err)
			}
			query = "UPDATE " + model.DeployClusterProductRel.TableName + " SET status=?, is_deleted=?, update_time=NOW() WHERE pid=? AND clusterId=? AND namespace=?"
			if _, err := model.DeployClusterProductRel.GetDB().Exec(query, model.PRODUCT_STATUS_UNDEPLOYED, 1, pid, clusterId, namespace); err != nil {
				log.Errorf("%v", err)
			}
		}
		return nil
	}
	return nil
}

func hanldeDefaultE(obj interface{}) error {
	return nil
}

func HandleResourceM(event *events.Event) error {
	sigc <- 0
	defer func() {
		<-sigc
	}()
	log.Infof("k8s resourcce monitor get event,  key: %v", event.Key)
	switch event.Resource {
	case "Pod":
		return handlePodE(event)
	case "Event":
		return handleEvent(event)
	case "Service":
		return handleServiceE(event)
	case "Deployment":
		return handleDeploymentE(event)
	case "Ingress":
		return handleIngressE(event)
	case "WorkloadProcess":
		return handleWorkload(event)
	default:
		return handleMoleE(event)
	}
	return nil
}
