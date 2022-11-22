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
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const (
	Major = "6000"
)

var devs []string

func GetDevs() []string {
	return devs
}

var reDev = regexp.MustCompile(fmt.Sprintf(`qdisc htb %s: dev (\w+)`, Major))

func ClassidToHandle(classid uint32) string {
	major, _ := strconv.ParseUint(Major, 16, 16)
	minor := uint64(classid) - major<<16
	return Major + ":" + strconv.FormatUint(minor, 16)
}

func DeleteTcDev() error {
	Debugf("delete old tc configuration")

	out, err := exec.Command("tc", "qdisc", "show").CombinedOutput()
	if err != nil {
		if len(out) > 0 {
			err = fmt.Errorf("tc qdisc show fail: %q", out)
		} else {
			err = fmt.Errorf("tc qdisc show fail: %v", err)
		}
		return err
	}
	handle := Major + ":"
	for _, match := range reDev.FindAllSubmatch(out, -1) {
		dev := string(match[1])
		out, err = exec.Command("tc", "qdisc", "del", "dev", dev, "root", "handle", handle, "htb").CombinedOutput()
		if err != nil {
			if len(out) > 0 {
				Errorf("tc del %v fail: %q", dev, out)
			} else {
				Errorf("tc del %v fail: %v", dev, err)
			}
		}
	}
	return nil
}

func InitTC(networks []string) error {
	if _, err := exec.LookPath("tc"); err != nil {
		return err
	}
	Infof("initialize tc configuration")

	if err := DeleteTcDev(); err != nil {
		return err
	}

	handle := Major + ":"
	for _, dev := range networks {
		iface, err := net.InterfaceByName(dev)
		if err != nil {
			Errorf("get network interface %v err: %v", dev, err)
			continue
		}
		if (iface.Flags&net.FlagLoopback) == 0 &&
			(iface.Flags&net.FlagPointToPoint) == 0 &&
			!strings.Contains(iface.Name, "docker") &&
			!strings.Contains(iface.Name, "veth") &&
			!strings.Contains(iface.Name, "flannel") &&
			!strings.Contains(iface.Name, "cni") &&
			!strings.Contains(iface.Name, "br-") &&
			!strings.Contains(iface.Name, ":") {
			cmd := exec.Command("tc", "qdisc", "add", "dev", iface.Name, "root", "handle", handle, "htb")
			if out, err := cmd.CombinedOutput(); err != nil {
				if len(out) > 0 {
					Errorf("add qdisc to %v fail: %q", iface.Name, out)
				} else {
					Errorf("add qdisc to %v fail: %v", iface.Name, err)
				}
				continue
			}
			Infof("add qdisc to %v successful", iface.Name)

			cmd = exec.Command("tc", "filter", "add", "dev", iface.Name, "protocol", "ip", "parent", handle, "prio", "1", "handle", handle, "cgroup")
			if out, err := cmd.CombinedOutput(); err != nil {
				if len(out) > 0 {
					Errorf("add filter to %v fail: %q", iface.Name, out)
				} else {
					Errorf("add filter to %v fail: %v", iface.Name, err)
				}
				continue
			}
			devs = append(devs, iface.Name)
			Infof("add filter to %v successful", iface.Name)
		}
	}

	return nil
}
