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

package register

import (
	"net"
	"os"
	"runtime"
	"time"

	"easyagent/internal/proto"
	"easyagent/internal/sidecar/base"
	"easyagent/internal/sidecar/client"
	"easyagent/internal/sidecar/controller"
	"github.com/elastic/gosigar"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"strings"
)

func GetRealIp(url string) string {
	ip := GetHostIP()
	conn, err := net.Dial("udp", url)
	if err == nil {
		ips := strings.Split(conn.LocalAddr().String(), ":")
		if len(ips) == 0 {
			return ip
		}
		ip = ips[0]
		conn.Close()
	}
	return ip
}

func GetHostIP() string {
	defaultIp := "0.0.0.0"
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return defaultIp
	}

	hostIP := ""
	// try to find IPv4 address first
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				//return ipnet.IP.String()
				hostIP += ipnet.IP.String() + ","
			}
			//if ipnet.IP.To16() != nil {
			//	return ipnet.IP.String()
			//}
		}
	}
	if hostIP != "" {
		return hostIP[:len(hostIP)-1]
	}
	// if nothing found try IPv6 address
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To16() != nil {
				//return ipnet.IP.String()
				hostIP += ipnet.IP.String() + ","
			}
		}
	}
	if hostIP != "" {
		return hostIP[:len(hostIP)-1]
	}
	// in doubt return default address
	return defaultIp
}

// It send grpc call RegisterSidecar
func RegisterSidecar(client client.EaClienter, controller *controller.Controller, callback []string) {
	mem, swap := gosigar.Mem{}, gosigar.Swap{}
	mem.Get()
	swap.Get()

	platform, _, version, err := host.PlatformInformation()
	if err != nil {
		base.Errorf("PlatformInformation error: %v", err)
	}

	var cpuSerial string
	if cpuInfo, err := cpu.Info(); err != nil {
		base.Errorf("get cpu info error: %v", err)
	} else {
		cpuSerial = cpuInfo[0].ModelName
	}

	host, _ := os.Hostname()

	go func() {
		req := &proto.RegisterRequest{
			OsType:     runtime.GOOS,
			OsPlatform: platform,
			OsVersion:  version,
			CpuSerial:  cpuSerial,
			CpuCores:   uint32(runtime.NumCPU()),
			MemSize:    mem.Total,
			SwapSize:   swap.Total,
			Host:       host,
			LocalIp:    GetRealIp(client.GetServerAddress()),
			CallBack:   callback,
		}
		for {
			err := client.RegisterSidecar(req)
			if err == nil {
				break
			}
			base.Errorf("RegisterSidecar error: %v", err)
			time.Sleep(3 * time.Second)
		}
	}()
}
