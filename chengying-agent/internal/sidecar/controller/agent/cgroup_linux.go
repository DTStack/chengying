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

package agent

import (
	"io/ioutil"
	"runtime"
	"strconv"

	"easyagent/internal/sidecar/base"
	"easyagent/internal/sidecar/controller/util"
	"github.com/containerd/cgroups"
	"github.com/opencontainers/runtime-spec/specs-go"
)

type Cgroup struct {
	cgroups.Cgroup
}

func (ag *agent) installCgroup() error {
	classid := allocClassid()
	cg, err := cgroups.New(base.CpuMemNetCLS, cgroups.StaticPath(ag.agentId.String()), &specs.LinuxResources{
		Network: &specs.LinuxNetwork{ClassID: &classid},
	})
	if err != nil {
		freeClassid(classid)
		return err
	}
	ag.cg = util.NewCgroup(cg, ag.agentId.String())
	ag.classid = classid
	return nil
}

func (ag *agent) unInstallCgroup() error {
	var err error

	if ag.cg != nil {
		freeClassid(ag.classid)
		err = ag.cg.Delete()
	}
	return err
}

func (ag *agent) updateCgroup(cpuLimit float32, memLimit uint64) error {
	var err error
	if ag.cg != nil {
		var defaultCpu int64 = -1
		var defaultMem int64 = -1
		var defaultSwappiness uint64 = getDefaultSwappiness()
		limit := &specs.LinuxResources{
			CPU:    &specs.LinuxCPU{Quota: &defaultCpu},
			Memory: &specs.LinuxMemory{Limit: &defaultMem, Swappiness: &defaultSwappiness},
		}
		if cpuLimit > 0 {
			defaultCpu = int64(100000*cpuLimit) * int64(runtime.NumCPU())
		}
		if memLimit > 0 {
			defaultMem = int64(memLimit)
			defaultSwappiness = 0
		}
		err = ag.cg.Update(limit)
	}
	return err
}

func getDefaultSwappiness() uint64 {
	buf, err := ioutil.ReadFile("/proc/sys/vm/swappiness")
	if err != nil {
		return 60
	}
	swappiness, err := strconv.ParseUint(string(buf), 10, 0)
	if err != nil {
		return 60
	}
	return swappiness
}
