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
	"testing"
	"time"
	"unsafe"
)

func TestGetTcpTable2(t *testing.T) {
	tcpTable := make([]byte, 1)
	var sizePointer uint32
	psizePointer := &sizePointer
	rlt, err := GetTcpTable2_(tcpTable, psizePointer, false)
	if rlt != ERROR_INSUFFICIENT_BUFFER {
		t.Fatalf("expect ERROR_INSUFFICIENT_BUFFER(%v), but its %v", ERROR_INSUFFICIENT_BUFFER, rlt)
	}
	t.Logf("---%v---%v---%v", rlt, sizePointer, err)
	tcpTable = make([]byte, sizePointer)
	rlt, err = GetTcpTable2_(tcpTable, psizePointer, false)
	if rlt != NO_ERROR {
		t.Fatalf("expect NO_ERROR(%v), but its %v", NO_ERROR, rlt)
	}
	sizeTcpRow := unsafe.Sizeof(MIB_TCPROW2{})
	offsetTable := unsafe.Offsetof(MIB_TCPTABLE2{}.table)
	for i := offsetTable; i < uintptr(*psizePointer); i += sizeTcpRow {
		tcpRow := (*MIB_TCPROW2)(unsafe.Pointer(&tcpTable[i]))
		t.Logf("%+v---%v", tcpRow, int2ip(tcpRow.dwRemoteAddr))
	}
	t.Logf("---%v---%v---%v", rlt, sizePointer, err)
}

func TestGetTcpRows(t *testing.T) {
	pidTcpRows, err := GetTcpRows()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v", pidTcpRows)
}

func TestGetPerTcpConnectionEStats(t *testing.T) {
	pidTcpRows, err := GetTcpRows()
	if err != nil {
		t.Error(err)
		return
	}
	dataRod := TCP_ESTATS_DATA_ROD_v0{}
	sizeDataRod := uint32(unsafe.Sizeof(dataRod))
	for pid, tcpRows := range pidTcpRows {
		for i := 0; i < len(tcpRows); i++ {
			rlt, err := GetPerTcpConnectionEStats(
				&tcpRows[i], TcpConnectionEstatsData,
				0, 0, 0,
				0, 0, 0,
				uintptr(unsafe.Pointer(&dataRod)), 0, sizeDataRod,
			)
			t.Logf("---%v---%v---%+v---%v", pid, rlt, dataRod, err)
		}
	}
}

func TestSetPerTcpConnectionEStats(t *testing.T) {
	pidTcpRows, err := GetTcpRows()
	if err != nil {
		t.Error(err)
		return
	}
	dataRw := TCP_ESTATS_DATA_RW_v0{true}
	sizeDataRw := uint32(unsafe.Sizeof(dataRw))
	for pid, tcpRows := range pidTcpRows {
		for i := 0; i < len(tcpRows); i++ {
			rlt, err := SetPerTcpConnectionEStats(
				&tcpRows[i], TcpConnectionEstatsData,
				uintptr(unsafe.Pointer(&dataRw)), 0, sizeDataRw,
				0,
			)
			t.Logf("---%v---%v---%v---", pid, rlt, err)
		}
	}
}

func TestGetPidTcpEstats(t *testing.T) {
	pid := 5148
	pidTcpRows, err := GetTcpRows()
	if err != nil {
		t.Error(err)
		return
	}
	tcpRows := make([]MIB_TCPROW, 0)
	dataRw := TCP_ESTATS_DATA_RW_v0{true}
	sizeDataRw := uint32(unsafe.Sizeof(dataRw))
	for _, row := range pidTcpRows[pid] {
		if row.dwState == MIB_TCP_STATE_ESTAB {
			rlt, err := SetPerTcpConnectionEStats(
				&row, TcpConnectionEstatsData,
				uintptr(unsafe.Pointer(&dataRw)), 0, sizeDataRw,
				0,
			)
			if rlt != NO_ERROR {
				t.Logf("SetPerTcpConnectionEStats error: %v(%v)", err, rlt)
				continue
			}
			tcpRows = append(tcpRows, row)
		}
	}

	time.Sleep(10 * time.Second)

	dataRod := TCP_ESTATS_DATA_ROD_v0{}
	sizeDataRod := uint32(unsafe.Sizeof(dataRod))
	for _, row := range tcpRows {
		rlt, err := GetPerTcpConnectionEStats(
			&row, TcpConnectionEstatsData,
			0, 0, 0,
			0, 0, 0,
			uintptr(unsafe.Pointer(&dataRod)), 0, sizeDataRod,
		)
		if rlt != NO_ERROR {
			t.Logf("GetPerTcpConnectionEStats error: %v(%v)", err, rlt)
			continue
		}
		t.Logf("-----%+v", dataRod)
	}
}
