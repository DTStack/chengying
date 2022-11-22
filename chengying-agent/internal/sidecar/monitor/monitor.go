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

package monitor

import (
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"easyagent/internal/proto"
	"easyagent/internal/sidecar/base"
	"easyagent/internal/sidecar/event"
	"github.com/elastic/gosigar"
	"github.com/satori/go.uuid"
	psnet "github.com/shirou/gopsutil/net"
)

var (
	// set by SetMonitorInterval
	monitorInterval = 3 * time.Second
)

var (
	pMu        sync.RWMutex
	processMap = map[uuid.UUID]*process{}
)

func SetMonitorInterval(d time.Duration) {
	monitorInterval = d
}

type process struct {
	agentId uuid.UUID
	pid     int
	classid uint32

	lastProcTime *gosigar.ProcTime
	lastTime     time.Time

	quit chan struct{}
}

func (p *process) collectMetrics() {
	base.Debugf("start collectMetrics PID: %v...", p.pid)

	p.lastProcTime = &gosigar.ProcTime{}
	if err := p.lastProcTime.Get(p.pid); err != nil {
		base.Errorf("get %v procTime error: %v", p.pid, err)
		p.lastProcTime = nil
	}
	p.lastTime = time.Now()

	if err := setTrafficEnable(uint32(p.pid)); err != nil {
		base.Errorf("%v setTrafficEnable error: %v", p.pid, err)
	}

	ticker := time.NewTicker(monitorInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
		case <-p.quit:
			base.Debugf("quit collectMetrics PID: %v", p.pid)
			return
		}

		ev := &proto.Event_ProcessResourceUsages{AgentId: p.agentId.Bytes()}

		currentProcTime := &gosigar.ProcTime{}
		if err := currentProcTime.Get(p.pid); err != nil {
			base.Errorf("get %v procTime error: %v", p.pid, err)
		} else {
			currentTime := time.Now()
			if p.lastProcTime != nil {
				// convert Millisecond to Nanosecond
				deltaProcTime := (currentProcTime.Total - p.lastProcTime.Total) / uint64(runtime.NumCPU()) * 1e6
				deltaTime := currentTime.Sub(p.lastTime).Nanoseconds()
				ev.CpuUsage = float32(float64(deltaProcTime) / float64(deltaTime))
				//base.Debugf("get %v cpu: %.2f", p.pid, ev.CpuUsage)
			}
			p.lastProcTime = currentProcTime
			p.lastTime = currentTime
		}

		mem := gosigar.ProcMem{}
		if err := mem.Get(p.pid); err != nil {
			base.Errorf("get %v memory error: %v", p.pid, err)
		} else {
			ev.Memory = mem.Resident
			//base.Debugf("get %v memory: %d", p.pid, ev.Memory)
		}

		cmd := gosigar.ProcArgs{}
		if err := cmd.Get(p.pid); err != nil {
			base.Errorf("get %v cmd error: %v", p.pid, err)
		} else {
			ev.Cmd = strings.Join(cmd.List, " ")
			//base.Debugf("get %v cmd: %q", p.pid, ev.Cmd)
		}

		var cp uint32
		if runtime.GOOS == "windows" {
			cp = uint32(p.pid)
		} else {
			cp = p.classid
		}
		if bytesSent, bytesRecv, err := getTraffic(cp); err != nil {
			base.Errorf("get %v traffic error: %v", p.pid, err)
		} else {
			ev.BytesSent, ev.BytesRecv = bytesSent, bytesRecv
			//base.Debugf("get %v bytesSent: %d, bytesRecv: %d", p.pid, bytesSent, bytesRecv)
		}

		event.ReportEvent(ev)
	}
}

func StartMonitAgent(agentId uuid.UUID, pid int, classid uint32) {
	p := &process{
		agentId: agentId,
		pid:     pid,
		classid: classid,
		quit:    make(chan struct{}),
	}

	pMu.Lock()
	if pOld, ok := processMap[agentId]; ok {
		delete(processMap, agentId)
		close(pOld.quit)
	}
	processMap[agentId] = p
	pMu.Unlock()

	go p.collectMetrics()
}

func StopMonitAgent(agentId uuid.UUID) {
	pMu.Lock()
	p, ok := processMap[agentId]
	if !ok {
		pMu.Unlock()
		return
	}
	delete(processMap, agentId)
	pMu.Unlock()

	close(p.quit)
}

func getAgents() []*process {
	ps := make([]*process, 0, 1)

	pMu.RLock()
	defer pMu.RUnlock()

	for _, p := range processMap {
		ps = append(ps, p)
	}
	return ps
}

func getWindowsDrives() (drives []gosigar.FileSystem) {
	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		dirName := string(drive) + ":\\"
		_, err := os.Open(dirName)
		if err == nil {
			fs := gosigar.FileSystem{DirName: dirName}
			drives = append(drives, fs)
		}
	}
	return
}

func getIfIp(name string) []string {
	inter, err := net.InterfaceByName(name)
	if err != nil {
		base.Errorf("get %v interface error: %v", name, err)
		return nil
	}
	addrs, err := inter.Addrs()
	if err != nil {
		base.Errorf("get %v Addrs error: %v", name, err)
		return nil
	}
	ips := make([]string, 0, len(addrs))
	for _, addr := range addrs {
		ips = append(ips, addr.String())
	}
	return ips
}

func collectSystemMetrics() {
	base.Debugf("start collectSystemMetrics...")

	last := &gosigar.Cpu{}
	if err := last.Get(); err != nil {
		base.Errorf("system cpu get error: %v", err)
		last = nil
	}
	for {
		time.Sleep(monitorInterval)

		ev := &proto.Event_OsResourceUsages{}

		current := &gosigar.Cpu{}
		if err := current.Get(); err != nil {
			base.Errorf("system cpu get error: %v", err)
		} else {
			if last != nil {
				deltaAll := float32(current.Total() - last.Total())
				deltaIdle := float32(current.Idle - last.Idle)
				ev.CpuUsage = 1 - deltaIdle/deltaAll
				//base.Debugf("system cpu: %.2f", ev.CpuUsage)
			}
			last = current
		}

		mem := gosigar.Mem{}
		if err := mem.Get(); err != nil {
			base.Errorf("system memory get error: %v", err)
		} else {
			ev.MemUsage = mem.ActualUsed
			//base.Debugf("system memory used: %d", ev.MemUsage)
		}

		swap := gosigar.Swap{}
		if err := swap.Get(); err != nil {
			base.Errorf("system swap get error: %v", err)
		} else {
			ev.SwapUsage = swap.Used
			//base.Debugf("system swap used: %d", ev.SwapUsage)
		}

		load := gosigar.LoadAverage{}
		if err := load.Get(); err != nil {
			base.Errorf("system load get error: %v", err)
		} else {
			ev.Load1 = float32(load.One)
			//base.Debugf("system load1: %.2f", ev.Load1)
		}

		uptime := gosigar.Uptime{}
		if err := uptime.Get(); err != nil {
			base.Errorf("system uptime get error: %v", err)
		} else {
			ev.Uptime = uptime.Length
			//base.Debugf("system uptime: %.2f(%v)", ev.Uptime, uptime.Format())
		}

		volumes := []gosigar.FileSystem{}
		if runtime.GOOS == "windows" {
			volumes = getWindowsDrives()
		} else {
			fslist := gosigar.FileSystemList{}
			if err := fslist.Get(); err != nil {
				base.Errorf("system filesystem get error: %v", err)
			} else {
				volumes = fslist.List
			}
		}
		for _, volume := range volumes {
			if strings.HasPrefix(volume.DevName, "/") {
				usage := gosigar.FileSystemUsage{}
				if err := usage.Get(volume.DirName); err != nil {
					base.Errorf("filesystem %v usage get error: %v", volume.DirName, err)
				} else {
					ev.DiskUsage = append(ev.DiskUsage, proto.Event_DiskUsage{
						MountPoint: volume.DirName,
						UsedSpace:  usage.Used,
						TotalSpace: usage.Total,
					})
				}
			}
		}
		//base.Debugf("system volumes: %v", ev.DiskUsage)

		netStats, err := psnet.IOCounters(true)
		if err != nil {
			base.Errorf("system netstats get error: %v", err)
		}
		for _, netStat := range netStats {
			if netStat.Name == "lo" {
				continue
			}
			ev.NetStats = append(ev.NetStats, proto.Event_NetStat{
				IfName:    netStat.Name,
				IfIp:      getIfIp(netStat.Name),
				BytesSent: netStat.BytesSent,
				BytesRecv: netStat.BytesRecv,
			})
		}
		//base.Debugf("system netstats: %v", ev.NetStats)

		event.ReportEvent(ev)
	}
}

func StartMonitSystem() {
	self := &process{
		pid:  os.Getpid(),
		quit: make(chan struct{}),
	}
	go self.collectMetrics()
	go collectSystemMetrics()
	go tcStatistic()
}
