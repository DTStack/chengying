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

package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type DBconf struct {
	User     string
	Password string
	Host     string
	Port     int
	DB       string
}

func InitDB(d DBconf) error {
	var err error
	dsn := `%s:%s@tcp(%s:3306)/%s`
	db, err = sql.Open("mysql", fmt.Sprintf(dsn, d.User, d.Password, d.Host, d.DB))
	if err != nil {
		fmt.Printf("Failed to open connection to mysql, error:%v\n", err)
		return err
	}
	err = db.Ping()
	if err != nil {
		fmt.Printf("Failed to connect to mysql, error:%v\n", err)
		return err
	}
	fmt.Println("connected to mysql")
	db.SetMaxOpenConns(1000)
	return nil
}

func DBConn() *sql.DB {
	return db
}
