/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package asset

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"easyagent/internal/server/log"
)

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Errorf("[SetAssetWithLocalFile] getCurrentDirectory err: %v", err)
		return ""
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func ResetInstallSidecarShWithLocalFile() error {
	file := ""
	for _, typ := range []string{"easyagent_install.sh", "easyagent_install_4win.ps1"} {
		if len(getCurrentDirectory()) > 0 {
			file = getCurrentDirectory() + "/" + typ
		}
		if _, err := os.Stat(file); os.IsNotExist(err) {
			log.Errorf("[InitInstallSidecarShWithLocalFile] %v, err: %v", typ, err)
			continue
		}
		fi, err := os.Open(file)
		if err != nil {
			log.Errorf("[SetAssetWithLocalFile] %v, err: %v", typ, err)
			continue
		}
		content, err := ioutil.ReadAll(fi)
		fi.Close()
		if err != nil {
			log.Errorf("[SetAssetWithLocalFile] %v, err: %v", typ, err)
			continue
		}
		switch typ {
		case "easyagent_install.sh":
			_templatesInstallSidecarSh = content
		case "easyagent_install_4win.ps1":
			_templatesInstallSidecarWinPs1 = content

		}
	}
	return nil
}
