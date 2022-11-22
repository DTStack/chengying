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

package kube

import (
	"database/sql"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"encoding/json"
	"github.com/jmoiron/sqlx"
)

var (
	getDefLatestSql    = "select * from workload_definition where name = :name and latest = 1"
	getDefSql          = "select * from workload_definition where name = :name and version = :version"
	getDefLatestSts    *sqlx.NamedStmt
	getDefSts          *sqlx.NamedStmt
	WorkloadDefinition = &workloadDefinition{
		PrepareFunc: prepareWorkloadDefinition,
	}
)

func prepareWorkloadDefinition() error {
	var err error
	getDefLatestSts, err = model.USE_MYSQL_DB().PrepareNamed(getDefLatestSql)
	if err != nil {
		log.Errorf("[kube workload_definition]: init sql: %s , error %v", getDefLatestSql, err)
		return err
	}
	getDefSts, err = model.USE_MYSQL_DB().PrepareNamed(getDefSql)
	if err != nil {
		log.Errorf("[kube workload_definition]: init sql: %s , error %v", getDefSql, err)
	}
	return nil
}

type workloadDefinition struct {
	PrepareFunc
}
type WorkloadDefinitionSchema struct {
	Id      int       `db:"id"`
	Name    string    `db:"name"`
	Version string    `db:"version"`
	Params  ParamsTyp `db:"params"`
	Latest  int       `db:"latest"`
}
type ParamsTyp string

func (p ParamsTyp) ToStruct() []ParamStruct {

	params := []ParamStruct{}
	json.Unmarshal([]byte(p), &params)
	return params
}

type ParamStruct struct {
	Key string `json:"key"`
	Ref string `json:"ref"`
}

// if the version is not specified, get the latest
func (w *workloadDefinition) Get(name, version string) (*WorkloadDefinitionSchema, error) {
	tbsc := &WorkloadDefinitionSchema{
		Name:    name,
		Version: version,
	}
	var err error
	if len(version) == 0 {
		err = getDefLatestSts.Get(tbsc, tbsc)
	} else {
		err = getDefSts.Get(tbsc, tbsc)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Errorf("[workload_definition]: get workload, sql %s, name %s, version %s, err %v", getDefLatestSql, name, version, err)
		return nil, err
	}
	return tbsc, nil
}

type WorkloadDefinitionYaml struct {
	ApiVersion string   `yaml:"apiVersion" jason:"apiVersion"`
	Kind       string   `yaml: "kind" json:"kind"`
	Metadata   Metadata `yaml: "metadata" json:"metadata"`
	Spec       Spec     `yaml: "spec" json:"spec"`
}

type Metadata struct {
	Name string `yaml:"name" json:"name"`
}

type Spec struct {
	Params        []Params        `yaml: "params" json:"params"`
	WorkloadPatrs []WorkloadPatrs `yaml: "workloadpatrs" json:"workloadpatrs"`
}

type Params struct {
	Key string `yaml:"key" json:"key"`
	Ref string `yaml:"ref" json:"ref"`
}

type WorkloadPatrs struct {
	Baseworkload Baseworkload `yaml: "baseworkload" json:"baseworkload"`
	Steps        []Steps      `yaml: "steps" json:"steps"`
}

type Baseworkload struct {
	Type       string `yaml: "type" json:"type"`
	Name       string `yaml: "name" json:"name"`
	Parameters struct {
		Spec struct {
			Replicas    int    `yaml:"replicas" json:"replicas,omitempty"`
			ServiceName string `yaml:"serviceName" json:"serviceName,omitempty"`
		} `yaml:"spec" json:"spec,omitempty"`
	} `yaml: "parameters" json:"parameters"`
}

type Steps struct {
	Name   string      `yaml: "name" json:"name"`
	Type   string      `yaml: "type" json:"type"`
	Action string      `yaml: "action" json:"action"`
	Object interface{} `yaml: "object" json:"object"`
}

func (w *workloadDefinition) InsertOrUpdate(db *sqlx.Tx, name, version, params string) (workloadId int64, err error) {
	workloaddef := WorkloadDefinitionSchema{}
	var workload_id int64

	query := "select id,name,version from workload_definition where name=? and version=?"
	wkdef_insert := "insert into workload_definition (name,version,params,latest) values(?,?,?,1)"
	wkdef_update := "update workload_definition set params=?  where name=? and version=?"
	err = db.Get(&workloaddef, query, name, version)
	if err == sql.ErrNoRows {
		relid, err := db.Exec(wkdef_insert, name, version, params)
		if err != nil {
			log.Errorf("Database insert workloaddefination error:%v, sql:%v\n", err, wkdef_insert)
			return 0, err
		}
		workload_id, _ = relid.LastInsertId()

	} else {
		_, err := db.Exec(wkdef_update, params, name, version)
		if err != nil {
			log.Errorf("Database update workloaddefination error:%v, sql:%v\n", err, wkdef_update)
			return 0, err
		}
		workload_id = int64(workloaddef.Id)
	}
	return workload_id, nil
}
