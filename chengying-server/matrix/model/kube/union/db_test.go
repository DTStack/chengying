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
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"fmt"
	"github.com/jmoiron/sqlx"
)

func getConn(){
	user := "root"
	password := "dtstack"
	host:= "172.16.10.37"
	port := 3306
	dbname := "dtagent_test"
	log.ConfigureLogger("/tmp/matrix",100,3,1)
	db,_ := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&loc=Local&parseTime=true", user, password, host, port, dbname))
	model.MYSQLDB = db
	err := Build()
	if err != nil{
		fmt.Printf("err: %v \n",err)
	}
}
//func Test1(t *testing.T) {
//	getConn()
//	tbscs,err := UnionT4T7.Select(58,"dtstack-system")
//	if err != nil{
//		fmt.Println("err",err.Error())
//		return
//	}
//	for _, sc := range tbscs{
//		fmt.Printf("sc %+v \n",sc)
//	}
//}
