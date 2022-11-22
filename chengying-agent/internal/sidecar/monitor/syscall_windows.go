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
	"unsafe"
)

const (
	NO_ERROR                  = 0
	ERROR_SUCCESS             = 0
	ERROR_INSUFFICIENT_BUFFER = 122
	ERROR_INVALID_PARAMETER   = 87
	ERROR_NOT_SUPPORTED       = 50
	ERROR_INVALID_USER_BUFFER = 1784
	ERROR_NOT_FOUND           = 1168
	ERROR_ACCESS_DENIED       = 5

	MIB_TCP_STATE_CLOSED     = 1
	MIB_TCP_STATE_LISTEN     = 2
	MIB_TCP_STATE_SYN_SENT   = 3
	MIB_TCP_STATE_SYN_RCVD   = 4
	MIB_TCP_STATE_ESTAB      = 5
	MIB_TCP_STATE_FIN_WAIT1  = 6
	MIB_TCP_STATE_FIN_WAIT2  = 7
	MIB_TCP_STATE_CLOSE_WAIT = 8
	MIB_TCP_STATE_CLOSING    = 9
	MIB_TCP_STATE_LAST_ACK   = 10
	MIB_TCP_STATE_TIME_WAIT  = 11
	MIB_TCP_STATE_DELETE_TCB = 12

	TcpConnectionEstatsData TCP_ESTATS_TYPE = 1
)

type TCP_ESTATS_TYPE uint8

type TCP_ESTATS_DATA_RW_v0 struct {
	EnableCollection bool
}

type TCP_ESTATS_DATA_ROD_v0 struct {
	DataBytesOut      uint64
	DataSegsOut       uint64
	DataBytesIn       uint64
	DataSegsIn        uint64
	SegsOut           uint64
	SegsIn            uint64
	SoftErrors        uint32
	SoftErrorReason   uint32
	SndUna            uint32
	SndNxt            uint32
	SndMax            uint32
	ThruBytesAcked    uint64
	RcvNxt            uint32
	ThruBytesReceived uint64
}

type MIB_TCPTABLE2 struct {
	dwNumEntries uint32
	table        [1]MIB_TCPROW2
}

type MIB_TCPROW struct {
	dwState      uint32
	dwLocalAddr  uint32
	dwLocalPort  uint32
	dwRemoteAddr uint32
	dwRemotePort uint32
}

type MIB_TCPROW2 struct {
	dwState        uint32
	dwLocalAddr    uint32
	dwLocalPort    uint32
	dwRemoteAddr   uint32
	dwRemotePort   uint32
	dwOwningPid    uint32
	dwOffloadState uint32
}

func GetTcpTable2_(buf []byte, SizePointer *uint32, Order bool) (uint32, error) {
	var _p0 *byte
	if len(buf) > 0 {
		_p0 = &buf[0]
	}
	var _p1 uint32
	if Order {
		_p1 = 1
	}
	r1, _, e1 := getTcpTable2.Call(uintptr(unsafe.Pointer(_p0)), uintptr(unsafe.Pointer(SizePointer)), uintptr(_p1))
	if r1 != NO_ERROR {
		return uint32(r1), e1
	}
	return uint32(r1), nil
}

func GetPerTcpConnectionEStats(
	Row *MIB_TCPROW,
	EstatsType TCP_ESTATS_TYPE,
	Rw uintptr, RwVersion uint32, RwSize uint32,
	Ros uintptr, RosVersion uint32, RosSize uint32,
	Rod uintptr, RodVersion uint32, RodSize uint32,
) (uint32, error) {
	r1, _, e1 := getPerTcpConnectionEStats.Call(
		uintptr(unsafe.Pointer(Row)),
		uintptr(EstatsType),
		Rw, uintptr(RwVersion), uintptr(RwSize),
		Ros, uintptr(RosVersion), uintptr(RosSize),
		Rod, uintptr(RodVersion), uintptr(RodSize),
	)
	if r1 != NO_ERROR {
		return uint32(r1), e1
	}
	return uint32(r1), nil
}

func SetPerTcpConnectionEStats(
	Row *MIB_TCPROW,
	EstatsType TCP_ESTATS_TYPE,
	Rw uintptr, RwVersion uint32, RwSize uint32,
	Offset uint32,
) (uint32, error) {
	r1, _, e1 := setPerTcpConnectionEStats.Call(
		uintptr(unsafe.Pointer(Row)),
		uintptr(EstatsType),
		Rw, uintptr(RwVersion), uintptr(RwSize),
		uintptr(Offset),
	)
	if r1 != NO_ERROR {
		return uint32(r1), e1
	}
	return uint32(r1), nil
}

// GetTcpRows only return tcp connection established to remote
func GetTcpRows() (map[int][]MIB_TCPROW, error) {
	tcpTable := make([]byte, 1)
	var sizePointer uint32
	psizePointer := &sizePointer
	rlt, err := GetTcpTable2_(tcpTable, psizePointer, false)
	if rlt == ERROR_INSUFFICIENT_BUFFER {
		tcpTable = make([]byte, sizePointer)
		rlt, err = GetTcpTable2_(tcpTable, psizePointer, false)
		if rlt != NO_ERROR {
			return nil, err
		}
	} else if rlt != NO_ERROR {
		return nil, err
	}

	pidTcpRows := make(map[int][]MIB_TCPROW, 64)
	sizeTcpRow := unsafe.Sizeof(MIB_TCPROW2{})
	offsetTable := unsafe.Offsetof(MIB_TCPTABLE2{}.table)
	for i := offsetTable; i < uintptr(*psizePointer); i += sizeTcpRow {
		tcpRow2 := (*MIB_TCPROW2)(unsafe.Pointer(&tcpTable[i]))
		if tcpRow2.dwState != MIB_TCP_STATE_ESTAB {
			continue
		}
		if ip := int2ip(tcpRow2.dwRemoteAddr); ip.IsLoopback() || ip.IsUnspecified() {
			continue
		}
		pid := int(tcpRow2.dwOwningPid)
		pidTcpRows[pid] = append(pidTcpRows[pid], MIB_TCPROW{
			dwState:      tcpRow2.dwState,
			dwLocalAddr:  tcpRow2.dwLocalAddr,
			dwLocalPort:  tcpRow2.dwLocalPort,
			dwRemoteAddr: tcpRow2.dwRemoteAddr,
			dwRemotePort: tcpRow2.dwRemotePort,
		})
	}
	return pidTcpRows, nil
}
