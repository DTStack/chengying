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

package union

import (
	"database/sql"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	modelkube "dtstack.com/dtstack/easymatrix/matrix/model/kube"
	"github.com/jmoiron/sqlx"
)

var (
	_getProductInfoSql = "select p.parent_product_name, p.product_name, p.product_version, p.product from deploy_product_list as p LEFT JOIN deploy_cluster_product_rel as r on p.id = r.pid where r.`status` = 'deployed'"
	//r.clusterId = :clusterId and r.namespace = :namespace and
	getparentProductsSql = _getProductInfoSql + " and r.clusterId = :clusterId and r.namespace = :namespace"
	getProductsSql = getparentProductsSql + " and p.parent_product_name = :parent_product_name"
	getServciesSql = getProductsSql + " and p.product_name = :product_name"
	getInstalledProductSql = "select p.product from deploy_product_list as p LEFT JOIN deploy_cluster_product_rel as r on p.id = r.pid where r.`status` = 'deployed' and r.clusterId = :clusterId and r.namespace = :namespace and p.product_name = :product_name"
	getparentProductsSts *sqlx.NamedStmt
	getProductsSts *sqlx.NamedStmt
	getServciesSts *sqlx.NamedStmt
	getInstalledProductSts *sqlx.NamedStmt
	UnionT4T7 = unionT4T7{
		PrepareFunc: prepareUnionT4T7,
	}
)

type unionT4T7 struct {
	modelkube.PrepareFunc
}

func prepareUnionT4T7() error{
	var err error
	getparentProductsSts,err = model.USE_MYSQL_DB().PrepareNamed(getparentProductsSql)
	if err != nil{
		log.Errorf("[union_T4_T7]: init sql: %s , error %v",getparentProductsSql,err)
		return err
	}
	getProductsSts,err = model.USE_MYSQL_DB().PrepareNamed(getProductsSql)
	if err != nil{
		log.Errorf("[union_T4_T7]: init sql %s, error %v",getProductsSql,err)
		return err
	}
	getServciesSts,err = model.USE_MYSQL_DB().PrepareNamed(getServciesSql)
	if err != nil{
		log.Errorf("[union_T4_T7]: init sql %s, error %v",getServciesSql,err)
		return err
	}
	getInstalledProductSts,err = model.USE_MYSQL_DB().PrepareNamed(getInstalledProductSql)
	if err != nil{
		log.Errorf("[union_T4_T7]: init sql %s, error %v",getInstalledProductSql,err)
		return err
	}
	return nil
}

type UnionT4T7Schema struct {
	modelkube.T4
	modelkube.T7
}

func (u *unionT4T7)SelectParentProduct(clusterid int, namespace string) ([]UnionT4T7Schema,error){
	unionSc := &UnionT4T7Schema{
		T4: modelkube.T4{
			ClusterId: clusterid,
			Namespace: namespace,
		},
	}
	unionscs := []UnionT4T7Schema{}
	rows,err := getparentProductsSts.Queryx(unionSc)
	if err != nil && err != sql.ErrNoRows{
		log.Errorf("[union_T4_T7]: sql %s, value %+v, error %v",getparentProductsSql,unionSc,err)
		return nil,err
	}
	for rows.Next(){
		err := rows.StructScan(unionSc)
		if err != nil{
			log.Errorf("[union_T4_T7]: struct scan to deploy_cluster_product_rel and deploy_product_list error :%v",err)
			return nil,err
		}
		unionscs = append(unionscs,*unionSc)
	}
	return unionscs,nil
}

func (u *unionT4T7)SelectProducts(clusterid int, namespace string, parentProduct string) ([]UnionT4T7Schema,error){
	unionSc := &UnionT4T7Schema{
		T7: modelkube.T7{
			ParentProductName: parentProduct,
		},
		T4: modelkube.T4{
			ClusterId: clusterid,
			Namespace: namespace,
		},
	}
	unionscs := []UnionT4T7Schema{}
	rows,err := getparentProductsSts.Queryx(unionSc)
	if err != nil && err != sql.ErrNoRows{
		log.Errorf("[union_T4_T7]: sql %s, value %+v, error %v",getparentProductsSql,unionSc,err)
		return nil,err
	}
	for rows.Next(){
		err := rows.StructScan(unionSc)
		if err != nil{
			log.Errorf("[union_T4_T7]: struct scan to deploy_cluster_product_rel and deploy_product_list error :%v",err)
			return nil,err
		}
		unionscs = append(unionscs,*unionSc)
	}
	return unionscs,nil
}

func (u *unionT4T7) GetProduct(clusterid int, namespace string,productname string) (*UnionT4T7Schema,error){
	unionSc := &UnionT4T7Schema{
		T7: modelkube.T7{
			ProductName: productname,
		},
		T4: modelkube.T4{
			ClusterId: clusterid,
			Namespace: namespace,
		},
	}
	err := getInstalledProductSts.Get(unionSc,unionSc)
	if err != nil{
		if err == sql.ErrNoRows{
			return nil,nil
		}
		log.Errorf("[union_T4_T7]: sql %s exec fail, clusterid %d, productname %s, namespace %s",clusterid,productname,namespace)
		return nil, err
	}
	return unionSc,nil
}
