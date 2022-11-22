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

package model

import (
	"database/sql"
	mydb "dtstack.com/dtstack/easymatrix/addons/easyfiler/pkg/model/mysql"
	"fmt"
)

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

func UploadAndFinished(filehash, filename, fileaddr string, filesize int64) bool {
	const sqlStr = `insert into easyft(file_sha1,file_name,file_size,file_addr,status) values(?,?,?,?,1)`
	stmt, err := mydb.DBConn().Prepare(sqlStr)
	if err != nil {
		fmt.Printf("Failed to prepare statement, error:%v\n", err)
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Println("Failed upload file")
		return false
	}
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Printf("File with hash:%s has been uploaded before", filehash)
		}
		return true
	}
	return false
}

func GetFileMeta(filehash string) (*TableFile, error) {
	const sqlStr = `select file_sha1,file_addr,file_name,file_size from eastft where file_sha1=? and status=1 limit 1`
	stmt, err := mydb.DBConn().Prepare(sqlStr)
	if err != nil {
		fmt.Printf(", error:%v\n", err)
		return nil, err
	}
	defer stmt.Close()
	tablefile := TableFile{}
	err = stmt.QueryRow(filehash).Scan(&tablefile.FileHash, &tablefile.FileAddr,
		&tablefile.FileName.String, &tablefile.FileSize)
	if err != nil {
		fmt.Println("error is here")
		fmt.Printf(err.Error())
		return nil, err
	}
	return &tablefile, nil
}

func DeleteFileMeta(filehash string) (ok bool, err error) {
	const sqlStr = `delete from easyft where file_sha1=? and status=1`
	stmt, err := mydb.DBConn().Prepare(sqlStr)
	if err != nil {
		fmt.Printf(", error:%v\n", err)
		return false, err
	}
	defer stmt.Close()
	_, err = stmt.Exec(filehash)
	if err != nil {
		fmt.Printf("Failed to delete this filemeta, error:%v\n", err)
		return false, err
	}
	return true, nil
}
