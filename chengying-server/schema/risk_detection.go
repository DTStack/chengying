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

package schema

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"dtstack.com/dtstack/easymatrix/matrix/util"
)

const (
	fileRe = `[^\s/]+\.{{fileType}}`
)

var dangerLinuxCmdRe = []string{
	`\brm\b`,
	`\bmkfs(\.\w+)*\b`,
	`\bdd\b`,
}

var dangerSqlRe = []string{
	`\bdrop\b`,
	`\bDROP\b`,
	`\bdelete\b`,
	`\bDELETE\b`,
}

func GetRisks(dir, cmd string) []string {
	var risks []string
	risks = append(risks, GetRiskCmds(cmd)...)
	shFiles := getFiles("sh", cmd)
	sqlFiles := getFiles("sql", cmd)
	for _, file := range shFiles {
		if data, err := ioutil.ReadFile(filepath.Join(dir, file)); err == nil {
			r := GetRiskCmds(string(data))
			risks = append(risks, r...)
		}
	}

	for _, file := range sqlFiles {
		if data, err := ioutil.ReadFile(filepath.Join(dir, file)); err == nil {
			r := GetRiskSqls(string(data))
			risks = append(risks, r...)
		}
	}

	return util.ArrayUniqueStr(risks)
}

func GetRiskCmds(cmd string) []string {
	reg := regexp.MustCompile(joinRes(dangerLinuxCmdRe))
	matches := reg.FindAllString(cmd, -1)

	for i := range matches {
		matches[i] = strings.ToLower(matches[i])
	}
	return matches
}

func GetRiskSqls(sqls string) []string {
	reg := regexp.MustCompile(joinRes(dangerSqlRe))
	matches := reg.FindAllString(sqls, -1)

	for i := range matches {
		matches[i] = strings.ToLower(matches[i])
	}
	return matches
}

func getFiles(fileType, str string) []string {
	re := strings.Replace(fileRe, "{{fileType}}", fileType, -1)
	exp := regexp.MustCompile(re)
	return exp.FindAllString(str, -1)
}

func joinRes(exps []string) string {
	var ret = ""
	for _, exp := range exps {
		if ret == "" {
			ret += "(" + exp + ")"
		} else {
			ret += "|(" + exp + ")"
		}
	}

	return ret
}
