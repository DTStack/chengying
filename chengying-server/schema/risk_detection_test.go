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
	"testing"
)

func TestGetRiskCmd(t *testing.T) {
	cmd := "rm bar && mkfs.ext3 && dd if=/devhda of=/dev/hdb && ./abc>123.txt "
	risks := GetRiskCmds(cmd)
	if len(risks) != 3 {
		t.Errorf("not all risk detected: %v\n", risks)
	}
}

func TestGetRiskCmdFile(t *testing.T) {
	sh := `#!/bin/sh
    if [ -z "$STATIC_URL" ];then
        rm bar
        echo "STATIC_URL is empty!"
        exit 1
    fi
    mkfs.ext3
    dd if=/devhda of=/dev/hdb
    ./abc>123.txt
`
	risks := GetRiskCmds(sh)
	if len(risks) != 3 {
		t.Errorf("not all risk detected: %v\n", risks)
	}
}

func TestGetRiskSqls(t *testing.T) {
	sqls := "DROP users; SELECT * FROM books; delete from users where 1"
	risks := GetRiskSqls(sqls)
	if len(risks) != 2 {
		t.Errorf("not all risk detected: %v\n", risks)
	}
}

func TestGetFiles(t *testing.T) {
	cmd := "./post_deploy.sh post_deploy.ini post_deploy.sql"

	shFiles := getFiles("sh", cmd)
	if fn := len(shFiles); fn != 1 {
		t.Fatalf("bad sh file nums: %d", fn)
	}
	if shFiles[0] != "post_deploy.sh" {
		t.Fatalf("bad sh filename: %s", shFiles[0])
	}

	sqlFiles := getFiles("sql", cmd)
	if fn := len(sqlFiles); fn != 1 {
		t.Fatalf("bad sql file nums: %d", fn)
	}
	if sqlFiles[0] != "post_deploy.sql" {
		t.Fatalf("bad sql filename: %s", sqlFiles[0])
	}

}
