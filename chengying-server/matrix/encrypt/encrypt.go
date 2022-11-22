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

package encrypt

import (
	"dtstack.com/dtstack/easymatrix/matrix/cache"
	"dtstack.com/dtstack/easymatrix/matrix/encrypt/aes"
	"dtstack.com/dtstack/easymatrix/matrix/encrypt/rsa"
	"dtstack.com/dtstack/easymatrix/matrix/encrypt/sm2"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"errors"
)

type commonEncrypt interface {
	CommonEncrypt(msg []byte) ([]byte, error)
	CommonDecrypt(msg []byte) ([]byte, error)
	SchemeDecrypt(msg string, aesPassword string) (string, error)
	CommonGetPublicKey() string
}

var (
	PlatformEncrypt commonEncrypt
	sm2Retry        = 0
	sm2MaxRetry     = 3
)

const (
	smDecryptErr = "Decrypt: failed to decrypt"
)

func InitPlatformEncrypt() {
	PlatformEncrypt = newEncrypt()
}

type platFormEncrypt struct {
	sm  sm2.SM2encrypt
	rsa rsa.Cipher
}

func (a *platFormEncrypt) CommonGetPublicKey() string {
	switch cache.SysConfig.PlatFormSecurity.LoginEncrypt {
	case "sm2":
		return a.sm.GetPubliKey()
	case "rsa":
		return a.rsa.PublicKeyToString()
	}
	return ""
}

func (a *platFormEncrypt) CommonEncrypt(msg []byte) ([]byte, error) {

	panic("implement me")
}

func (a *platFormEncrypt) CommonDecrypt(msg []byte) ([]byte, error) {
	var plain []byte
	switch cache.SysConfig.PlatFormSecurity.LoginEncrypt {
	case "sm2":
		decryptMsg, err := a.sm.DecryptHexString(msg)
		// 因为这个库和前端存在不兼容的问题，第一次需要先看看前端发来是不是能解密
		if sm2Retry < sm2MaxRetry && err != nil && err.Error() == smDecryptErr {
			log.Debugf("web decrypt err %v ,retry create sm2", err)
			InitPlatformEncrypt()
			sm2Retry++
		}
		if err != nil {
			return nil, err
		}
		plain = decryptMsg
	case "rsa":
		decryptMsg, err := a.rsa.Decrypt(msg)
		if err != nil {
			return nil, err
		}
		plain = decryptMsg
	}
	return plain, nil
}

// aesPassword == ""
// scheme 加密 rsa类型时候，默认使用aes
func (a *platFormEncrypt) SchemeDecrypt(msg string, aesPassword string) (string, error) {
	switch cache.SysConfig.PlatFormSecurity.LoginEncrypt {
	case "sm2":
		decryptMsg, err := a.sm.DecryptHexString([]byte(msg))
		if err != nil {
			return "", err
		}
		log.Debugf("[platFormEncrypt.SchemeDecrypt] %s %s", msg, decryptMsg)
		return string(decryptMsg), nil
	case "rsa":
		password, err := aes.AesDecryptByPassword(msg, aesPassword)
		if err != nil {
			return "", err
		}
		return password, nil
	}
	return "", errors.New("no support encrypt type")
}

func newEncrypt() commonEncrypt {
	return &platFormEncrypt{
		sm:  sm2.NewSm2Encrypt(),
		rsa: rsa.NewRsaEncrypt(),
	}
}
