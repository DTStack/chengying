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
	"encoding/binary"
	"net"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"easyagent/internal/sidecar/base"
)

var (
	getTcpTable2, getPerTcpConnectionEStats, setPerTcpConnectionEStats *syscall.Proc

	pmMu   sync.Mutex
	pidMap = map[int][]MIB_TCPROW{}
)

func init() {
	iphlpapi, err := syscall.LoadDLL("Iphlpapi.dll")
	if err != nil {
		base.Errorf("LoadDll error: %v", err)
		return
	}
	getTcpTable2, err = iphlpapi.FindProc("GetTcpTable2")
	if err != nil {
		base.Errorf("FindProc error: %v", err)
		return
	}
	getPerTcpConnectionEStats, err = iphlpapi.FindProc("GetPerTcpConnectionEStats")
	if err != nil {
		base.Errorf("FindProc error: %v", err)
		return
	}
	setPerTcpConnectionEStats, err = iphlpapi.FindProc("SetPerTcpConnectionEStats")
	if err != nil {
		base.Errorf("FindProc error: %v", err)
		return
	}
}

func int2ip(i uint32) net.IP {
	ip := make(net.IP, 4)
	binary.LittleEndian.PutUint32(ip, i)
	return ip
}

func getTraffic(pid uint32) (uint64, uint64, error) {
	if getTcpTable2 == nil || getPerTcpConnectionEStats == nil {
		return 0, 0, nil
	}
	pmMu.Lock()
	tcpRows := pidMap[int(pid)]
	pmMu.Unlock()

	var bytesSent, bytesRecv uint64
	dataRod := TCP_ESTATS_DATA_ROD_v0{}
	sizeDataRod := uint32(unsafe.Sizeof(dataRod))
	for i := 0; i < len(tcpRows); i++ {
		rlt, err := GetPerTcpConnectionEStats(
			&tcpRows[i], TcpConnectionEstatsData,
			0, 0, 0,
			0, 0, 0,
			uintptr(unsafe.Pointer(&dataRod)), 0, sizeDataRod,
		)
		if rlt != NO_ERROR {
			base.Errorf("GetPerTcpConnectionEStats error: %v(%v)", err, rlt)
			continue
		}
		bytesSent += dataRod.DataBytesOut
		bytesRecv += dataRod.DataBytesIn
	}
	return bytesSent, bytesRecv, nil
}

func setTrafficEnable(pid uint32) error {
	if setPerTcpConnectionEStats == nil {
		return nil
	}
	pmMu.Lock()
	tcpRows := pidMap[int(pid)]
	pmMu.Unlock()

	dataRw := TCP_ESTATS_DATA_RW_v0{true}
	sizeDataRw := uint32(unsafe.Sizeof(dataRw))
	for i := 0; i < len(tcpRows); i++ {
		rlt, err := SetPerTcpConnectionEStats(
			&tcpRows[i], TcpConnectionEstatsData,
			uintptr(unsafe.Pointer(&dataRw)), 0, sizeDataRw,
			0,
		)
		if rlt != NO_ERROR {
			base.Errorf("SetPerTcpConnectionEStats error: %v(%v)", err, rlt)
		}
	}
	return nil
}

func tcStatistic() {
	for {
		time.Sleep(monitorInterval)

		pidTcpRows, err := GetTcpRows()
		if err != nil {
			base.Errorf("get TCP map error: %v", err)
			continue
		}

		pmMu.Lock()
		pidMap = pidTcpRows
		pmMu.Unlock()
	}
}
