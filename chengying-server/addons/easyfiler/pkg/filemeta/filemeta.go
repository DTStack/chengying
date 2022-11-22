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

package filemeta

import (
	"dtstack.com/dtstack/easymatrix/addons/easyfiler/pkg/model"
)

type FileMeta struct {
	FileSha1 string
	FileMD5  string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

func UpdateFileMeta(f FileMeta) bool {
	return model.UploadAndFinished(f.FileSha1, f.FileName, f.Location, f.FileSize)
}

func GetFileMeta(filesha1 string) (filemeta FileMeta, err error) {
	tablefile, err := model.GetFileMeta(filesha1)
	if err != nil {
		return FileMeta{}, err
	}
	filemeta = FileMeta{
		FileSha1: tablefile.FileHash,
		FileName: tablefile.FileName.String,
		FileSize: tablefile.FileSize.Int64,
		Location: tablefile.FileAddr.String,
	}
	return
}

func DelteFileMeta(fileSha1 string) {
	model.DeleteFileMeta(fileSha1)
}
