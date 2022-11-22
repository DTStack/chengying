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

package utils

import (
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

//http response
var (
	ResponseResult = &ResponseRes{}
	DefaultClient  = NewHTTPClient(GetTlsConfig())
)

type ResponseRes struct {
	ResultCode    int                    `json:"result_code"`
	ResultMessage string                 `json:"result_message"`
	Result        string                 `json:"result"`
	Data          map[string]interface{} `json:"data"`
}

func (responseRes *ResponseRes) SetResponseRes(resCode int, resMsg string, res string, data map[string]interface{}) {
	responseRes.Result = res
	responseRes.Data = data
	responseRes.ResultCode = resCode
	responseRes.ResultMessage = resMsg
}

//response body

type ResposeBody struct {
	Msg  string                 `json:"msg"`
	Code int                    `json:"code"`
	Data map[string]interface{} `json:"data"`
}

// use to create a http request
func GetTlsConfig() *tls.Config {
	var tlsConfig *tls.Config
	tlsConfig = &tls.Config{InsecureSkipVerify: true}
	return tlsConfig
}

// copyFile copies from src to dst until either EOF is reached
// on src or an error occurs. It verifies src exists and removes
// the dst if it exists.
func CopyFile(src, dst string) (int64, error) {
	cleanSrc := filepath.Clean(src)
	cleanDst := filepath.Clean(dst)
	if cleanSrc == cleanDst {
		return 0, nil
	}
	sf, err := os.Open(cleanSrc)
	if err != nil {
		return 0, err
	}
	defer sf.Close()
	if err := os.Remove(cleanDst); err != nil && !os.IsNotExist(err) {
		return 0, err
	}
	df, err := os.Create(cleanDst)
	if err != nil {
		return 0, err
	}
	defer df.Close()
	return io.Copy(df, sf)
}

func FoundIpIdx(ipList []string, ip string) uint {
	for idx, target := range ipList {
		if target == ip {
			return uint(idx)
		}
	}
	return 0
}

func StringContain(a []string, s string) bool {
	for _, str := range a {
		if str == s {
			return true
		}
	}
	return false
}

func ListFiles(dir string) ([]string, error) {
	var files []string
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return files, err
	}
	for _, file := range fs {
		if file.IsDir() {
			d := filepath.Join(dir, file.Name())
			f, _ := ListFiles(d)
			files = append(files, f...)
		} else {
			files = append(files, filepath.Join(dir, file.Name()))
		}
	}
	return files, nil
}

func Create(path string) (*os.File, error) {
	fi, err := os.Open(path)
	defer fi.Close()
	if err != nil {
		return nil, err
	}
	err = os.Rename(path, path+"."+strconv.FormatInt(time.Now().Unix(), 10)+"~")
	if err != nil {
		return nil, err
	}
	return os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_SYNC, 0600)
}

func Md5(raw string) string {
	h := md5.New()
	io.WriteString(h, raw)

	return fmt.Sprintf("%x", h.Sum(nil))
}

func IsPathExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}
