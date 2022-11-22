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
	"bytes"
	basicAES "crypto/aes"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"encoding/base64"
	"fmt"

	"crypto/cipher"
)

//前端无法处理空串加解密；所以用 uuid 代替空串作为密码加解密
var replaceBlank = "45d3161f8aa54588ab3901fa954cc01a"

//加密
func AesEncryptByPassword(original, password string) (string, error) {

	//加密前如果是空串，那么用特定 uuid 代替空串加密
	if original == "" {
		original = replaceBlank
	}

	const iv = "1234567890123456"
	aes := aesTool(password, iv)
	cipherText, err := aes.encrypt(original)
	if err != nil {
		log.Errorf("加密失败")
		return "", fmt.Errorf("加密失败")
	}
	return cipherText, nil
	//outPlainText, _ := aes.decrypt(cipherText)
	//fmt.Println("解密后明文：" + outPlainText)
}

//解密
func AesDecryptByPassword(cipherText, password string) (string, error) {
	const iv = "1234567890123456"
	aes := aesTool(password, iv)
	pwd, err := aes.decrypt(cipherText)
	if err != nil {
		log.Errorf("解密失败")
		return "", fmt.Errorf("解密失败")
	}

	//解密后如果是特定 uuid,那么密码其实为空串
	if pwd == replaceBlank {
		pwd = ""
	}
	return pwd, nil
}

type aes struct {
	securityKey []byte
	iv          []byte
}

/**
 * constructor
 */
func aesTool(securityKey string, iv string) *aes {
	return &aes{[]byte(securityKey), []byte(iv)}
}

/**
 * 加密
 * @param string $plainText 明文
 * @return bool|string
 */
func (a aes) encrypt(plainText string) (string, error) {
	block, err := basicAES.NewCipher(a.securityKey)
	if err != nil {
		return "", err
	}
	plainTextByte := []byte(plainText)
	blockSize := block.BlockSize()
	plainTextByte = addPKCS7Padding(plainTextByte, blockSize)
	cipherText := make([]byte, len(plainTextByte))
	mode := cipher.NewCBCEncrypter(block, a.iv)
	mode.CryptBlocks(cipherText, plainTextByte)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

/**
 * 解密
 * @param string $cipherText 密文
 * @return bool|string
 */
func (a aes) decrypt(cipherText string) (string, error) {
	block, err := basicAES.NewCipher(a.securityKey)
	if err != nil {
		return "", err
	}
	cipherDecodeText, decodeErr := base64.StdEncoding.DecodeString(cipherText)
	if decodeErr != nil {
		return "", decodeErr
	}
	mode := cipher.NewCBCDecrypter(block, a.iv)
	originCipherText := make([]byte, len(cipherDecodeText))
	mode.CryptBlocks(originCipherText, cipherDecodeText)
	originCipherText = stripPKSC7Padding(originCipherText)
	return string(originCipherText), nil
}

/**
 * 填充算法
 * @param string $source
 * @return string
 */
func addPKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	paddingText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, paddingText...)
}

/**
 * 移去填充算法
 * @param string $source
 * @return string
 */
func stripPKSC7Padding(cipherText []byte) []byte {
	length := len(cipherText)
	unpadding := int(cipherText[length-1])
	return cipherText[:(length - unpadding)]
}
