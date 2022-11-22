// Licensed to Apache Software Foundation(ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation(ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package model

//func init() {
//	err := ConfigureMysqlDatabase("172.16.8.165", 3306, "root", "dtstack", "dtagent")
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//}

//func TestCreateServiceIpNode(t *testing.T) {
//	node := ServiceIpNode{
//		ClusterId:   1,
//		ProductName: "product test",
//		ServiceName: "service 1",
//		Ip:          "127.0.0.1",
//		NodeId:      3,
//	}
//
//	if err := node.Create(); err != nil {
//		t.Error(err)
//	}
//}
//
//func TestGetServiceIpNodes(t *testing.T) {
//	nodes, err := GetServiceNodes(1, "product test", "service 1")
//	if err != nil {
//		t.Error(err)
//	}
//	bytes, _ := json.Marshal(nodes)
//	fmt.Printf("nodes: %s\n", bytes)
//}
//
//func TestGetServiceIpNode(t *testing.T) {
//	node, err := GetServiceIpNode(1, "product test", "service 1", "127.0.0.2")
//	if err != nil {
//		t.Error(err)
//	}
//	fmt.Printf("%v", node)
//}
