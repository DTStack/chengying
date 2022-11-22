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

package base

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/containerd/cgroups"
)

const (
	cgroupDir = "cgroup"
)

// cgroupInfo returns cgroup root and all mount points
// where the cgroup mountpoints are mounted in a single hiearchy
func cgroupInfo() (root string, cp map[string]string, err error) {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return
	}
	defer f.Close()
	cp = make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if err = scanner.Err(); err != nil {
			return
		}
		var (
			text   = scanner.Text()
			fields = strings.Split(text, " ")
			// safe as mountinfo encodes mountpoints with spaces as \040.
			index               = strings.Index(text, " - ")
			postSeparatorFields = strings.Fields(text[index+3:])
			numPostFields       = len(postSeparatorFields)
		)
		// this is an error as we can't detect if the mount is for "cgroup"
		if numPostFields < 2 {
			err = fmt.Errorf("Found less fields post '-' in %q", text)
			return
		}
		if postSeparatorFields[0] == "tmpfs" && postSeparatorFields[1] == "easyagent" {
			root = fields[4]
		}
		if postSeparatorFields[0] == "cgroup" {
			// check that the mount is properly formated.
			if numPostFields < 3 {
				err = fmt.Errorf("Error found less than 3 fields post '-' in %q", text)
				return
			}
			root = filepath.Dir(fields[4])
			option := strings.SplitN(postSeparatorFields[numPostFields-1], ",", -1)
			for i := 1; i < len(option); i++ {
				cp[option[i]] = fields[4]
			}
		}
	}
	if root != "" {
		return
	}

	err = cgroups.ErrMountPointNotExist
	return
}

func CpuMemNetCLS() ([]cgroups.Subsystem, error) {
	subsys, err := cgroups.V1()
	if err != nil {
		return nil, err
	}

	newSubsys := make([]cgroups.Subsystem, 0, 3)
	for _, v := range subsys {
		switch v.Name() {
		case cgroups.Cpu:
			newSubsys = append(newSubsys, v)
		case cgroups.Memory:
			newSubsys = append(newSubsys, v)
		case cgroups.NetCLS:
			newSubsys = append(newSubsys, v)
			//case cgroups.Freezer:
			//	newSubsys = append(newSubsys, v)
		}
	}
	if len(newSubsys) == 0 {
		return nil, errors.New("No cgroup subsystem can use")
	}
	return newSubsys, nil
}

func mountCPU(path string) error {
	cpuPath := filepath.Join(path, string(cgroups.Cpu))
	if err := os.MkdirAll(cpuPath, 0755); err != nil {
		Errorf("%v", err)
		return err
	}
	Infof("%v", "we need mount cpu subsystem")
	err := syscall.Mount("easyagent", cpuPath, "cgroup", 0, "cpu")
	if err != nil {
		Errorf("mount error: %v", err)
	}
	return err
}

func mountMemory(path string) error {
	memPath := filepath.Join(path, string(cgroups.Memory))
	if err := os.MkdirAll(memPath, 0755); err != nil {
		Errorf("%v", err)
		return err
	}
	Infof("%v", "we need mount memory subsystem")
	err := syscall.Mount("easyagent", memPath, "cgroup", 0, "memory")
	if err != nil {
		Errorf("mount error: %v", err)
	}
	return err
}

func mountNetCLS(path string) error {
	netClsPath := filepath.Join(path, string(cgroups.NetCLS))
	if err := os.MkdirAll(netClsPath, 0755); err != nil {
		Errorf("%v", err)
		return err
	}
	Infof("%v", "we need mount net_cls subsystem")
	err := syscall.Mount("easyagent", netClsPath, "cgroup", 0, "net_cls")
	if err != nil {
		Errorf("mount error: %v", err)
	}
	return err
}

func mountFreezer(path string) error {
	freezerPath := filepath.Join(path, string(cgroups.Freezer))
	if err := os.MkdirAll(freezerPath, 0755); err != nil {
		Errorf("%v", err)
		return err
	}
	Infof("%v", "we need mount freezer subsystem")
	err := syscall.Mount("easyagent", freezerPath, "cgroup", 0, "freezer")
	if err != nil {
		Errorf("mount error: %v", err)
	}
	return err
}

func mountTmpfs(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		Errorf("%v", err)
		return err
	}
	Infof("%v", "we need mount tmpfs for cgroup")
	err := syscall.Mount("easyagent", path, "tmpfs", 0, "")
	if err != nil {
		Errorf("mount error: %v", err)
	}
	return err
}

func MountCgroup() error {
	cgroupRoot, cp, err := cgroupInfo()
	if err == cgroups.ErrMountPointNotExist {
		cgroupRoot = filepath.Join(os.TempDir(), cgroupDir)
		if err = mountTmpfs(cgroupRoot); err != nil {
			return err
		}
		err = mountCPU(cgroupRoot)
		err = mountMemory(cgroupRoot)
		err = mountNetCLS(cgroupRoot)
		//err = mountFreezer(cgroupRoot)
		return err
	} else if err != nil {
		return fmt.Errorf("get cgroup root error: %v", err)
	}

	if _, ok := cp[string(cgroups.Cpu)]; !ok {
		err = mountCPU(cgroupRoot)
	}
	if _, ok := cp[string(cgroups.Memory)]; !ok {
		err = mountMemory(cgroupRoot)
	}
	if _, ok := cp[string(cgroups.NetCLS)]; !ok {
		err = mountNetCLS(cgroupRoot)
	}
	//if _, ok := cp[string(cgroups.Freezer)]; !ok {
	//	err = mountFreezer(cgroupRoot)
	//}

	return err
}
