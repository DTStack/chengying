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

package util

import (
	"context"
	"io/ioutil"
	"os"
	"os/exec"
)

type Cmd struct {
	*exec.Cmd
	ctx context.Context
}

type Cgroup struct{}

func (c *Cgroup) GetInitStub() string {
	return ""
}

func CommandContext(ctx context.Context, user string, isSeniorKill bool, cg *Cgroup, name string, arg ...string) *Cmd {
	return &Cmd{Cmd: exec.CommandContext(ctx, name, arg...)}
}

func CreateTempScript(content string, prefix string) (path string, err error) {
	f, err := ioutil.TempFile("", prefix)
	if err != nil {
		return "", err
	}

	if _, err = f.WriteString(content); err != nil {
		f.Close()
		os.Remove(f.Name())
		return "", err
	}

	f.Close()
	newpath := f.Name() + ".bat"
	if err = os.Rename(f.Name(), newpath); err != nil {
		os.Remove(f.Name())
		return "", err
	}

	return newpath, nil
}
