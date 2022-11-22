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

package aes

import (
	"fmt"
	"testing"
)

func TestAes(t *testing.T) {
	plainText := ""
	securityKey16 := "ca6590a271539cc89e2cc20bd6b58518"
	iv := "1234567890123456"
	aes := aesTool(securityKey16, iv)
	cipherText, _ := aes.encrypt(plainText)
	fmt.Println("加密后的密文：" + cipherText)
	outPlainText, _ := aes.decrypt(cipherText)
	fmt.Println("解密后明文：" + outPlainText)

}

//707a161834ea0e5a3f1cc419a37bc030
//cca690a271539cc89e2cc20bd6b58518