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
	"dtstack.com/dtstack/easymatrix/matrix/api/k8s/view"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/constant"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/kube"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/resource/deployment"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	modelkube "dtstack.com/dtstack/easymatrix/matrix/model/kube"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
)

func NamespaceGet(namespace,clusterid string) (*view.NamespaceGetRsp,error){
	cid,err := strconv.Atoi(clusterid)
	if err != nil{
		log.Errorf("[namespace_client_info]: convert clusterid %s to int error %v",clusterid,err)
		return nil,err
	}
	tbsc,err := modelkube.DeployNamespaceList.Get(namespace,cid)
	if err != nil && tbsc == nil{
		return nil, err
	}
	clientTbsc,err := modelkube.DeployNamespaceClient.Get(tbsc.Id)
	if err != nil && tbsc == nil{
		return nil,err
	}
	imageTbsc,err := modelkube.DeployClusterImageStore.GetById(tbsc.RegistryId)
	if err != nil{
		return nil,err
	}
	reigistryId := tbsc.RegistryId
	if imageTbsc == nil{
		reigistryId = -1
	}
	return &view.NamespaceGetRsp{
		Id: 	   tbsc.Id,
		Type:      tbsc.Type,
		Namespace: namespace,
		Registry:  reigistryId,
		Yaml:      clientTbsc.Yaml,
		FileName:  clientTbsc.Filename,
	},nil
}

func NamespaceDeleteConfirm(namesapce,clusterid string) error{
	cid,err := strconv.Atoi(clusterid)
	if err != nil{
		log.Errorf("[namespace_client_info]: convert clusterid %s to int error %v",clusterid,err)
		return err
	}
	list,err := modelkube.DeployClusterProductRel.SelectNamespacedDeployed(cid,namesapce)
	if err != nil{
		return err
	}
	if len(list) !=0{
		return fmt.Errorf("deployed produtct exist in namespace %s",namesapce)
	}
	return nil
}

func NamespaceDelete(ctx context.Context,namespace,clusterid string) error{
	if err := NamespaceDeleteConfirm(namespace,clusterid);err != nil{
		return err
	}
	cid,err := strconv.Atoi(clusterid)
	if err != nil{
		log.Errorf("[namespace_client_info]: convert clusterid %s to int error %v",clusterid,err)
		return err
	}
	tbsc,err := modelkube.DeployNamespaceList.Get(namespace,cid)
	if err != nil{
		return err
	}
	if tbsc == nil{
		return nil
	}
	if tbsc.Status == constant.NAMESPACE_VALID{
		err = deleteInitResource(ctx,namespace,clusterid)
		if err != nil{
			return err
		}

		kube.ClusterNsClientCache.GetClusterNsClient(clusterid).DeleteNsClient(namespace)
	}

	return modelkube.DeployNamespaceList.Delete(namespace,cid)
}

//delete the deployment mole-operator
//other init moudles should be deleted, do that when have time
func deleteInitResource(ctx context.Context,namespace,clusterid string)error{
	cid,err := strconv.Atoi(clusterid)
	if err != nil{
		log.Errorf("[namespace_client_info]: convert clusterid %s to int error %v",clusterid,err)
		return err
	}
	tbsc,err := modelkube.DeployNamespaceList.Get(namespace,cid)
	if err != nil{
		return err
	}
	if tbsc == nil{
		return nil
	}

	cache,err := kube.ClusterNsClientCache.GetClusterNsClient(clusterid).GetClientCache(kube.ImportType(tbsc.Type))
	if err != nil{
		return err
	}
	d := deployment.New()
	d.Namespace = namespace
	//set the fixed value for mole operator.
	//it should be change when have time
	d.Name = "mole-operator"
	err = cache.GetClient(namespace).Delete(ctx,d)
	if err != nil{
		if err, b := err.(*errors.StatusError); b {
			if err.Status().Reason == metav1.StatusReasonNotFound {
				return nil
			}
		}
		return err
	}
	return nil
}
