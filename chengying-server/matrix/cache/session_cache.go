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

package cache

import (
	"crypto/md5"
	"fmt"
	"io"
	"math"
	"regexp"
	"strconv"
	"time"
)

var (
	//	/api/v2/product/xxx/xxx/history
	//	/api/v2/cluster/restartServices
	//	/api/v2/product/anomalyService
	//	/api/v2/instance/product/xxx/service/xxx
	//	/api/v2/instance/product/xxx/service/xxx/healthCheck
	//	/api/v2/cluster/restartServices
	//	/api/v2/cluster/hosts/hosts
	//	/api/v2/cluster/hostgroups
	//	/api/v2/cluster/orderList
	//	/api/v2/cluster/orderDetail
	noFlushApiReg = []*regexp.Regexp{
		regexp.MustCompile("/api/v2/product/.*/.*/history"),
		regexp.MustCompile("/api/v2/cluster/restartServices"),
		regexp.MustCompile("/api/v2/product/anomalyService"),
		regexp.MustCompile("/api/v2/instance/product/.*/service/.*"),
		regexp.MustCompile("/api/v2/instance/product/.*/service/.*/healthCheck"),
		regexp.MustCompile("/api/v2/cluster/restartServices"),
		regexp.MustCompile("/api/v2/cluster/hosts/hosts"),
		regexp.MustCompile("/api/v2/cluster/hostgroups"),
		regexp.MustCompile("/api/v2/cluster/orderList"),
		regexp.MustCompile("/api/v2/cluster/orderDetail"),
	}
)

func SetSessionCache(sessionStr string) {
	cleanTime := SysConfig.PlatFormSecurity.AccountLogoutSleepTime
	sessionIdCache.Set(sessionStr, true, time.Minute*time.Duration(cleanTime))
}

func ValidationSessionId(path, hashId string) bool {
	_, ok := sessionIdCache.Get(hashId)
	if !ok {
		return false
	}
	for _, api := range noFlushApiReg {
		if api.MatchString(path) {
			return true
		}
	}
	SetSessionCache(hashId)
	return true
}

func createSessionId() (f func() string) {
	i := 0
	return func() string {
		w := md5.New()
		io.WriteString(w, strconv.Itoa(i))
		i++
		if i == math.MaxInt32 {
			i = 0
		}
		sessionStr := fmt.Sprintf("%x", w.Sum(nil))
		return sessionStr
	}
}
