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

package sm2

import (
	"crypto/rand"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"encoding/hex"
	"fmt"
	"github.com/tjfoc/gmsm/sm2"
)

// https://juejin.cn/post/6966639973340545032

var (
	mode = 0 // C1C3C2=0 C1C2C3=1
)

type SM2encrypt interface {
	Encrypt(msg []byte) ([]byte, error)          // publickey 加密
	DecryptHexString(msg []byte) ([]byte, error) // privatekey 16进制字符串解密
	DecryptByte(msg []byte) ([]byte, error)      // privatekey 字节解密
	GetPubliKey() string
	GetPrivateKey() string
	Sign(msg string) (sign []byte, ok bool) // privatekey签名
	Verify(msg, sign string) bool           // publickey 验签
}

type sm struct {
	privatekey *sm2.PrivateKey
	publickey  *sm2.PublicKey
}

func (s *sm) DecryptHexString(msg []byte) ([]byte, error) {
	msgHexDecode, _ := hex.DecodeString(string(msg))
	planiText, err := sm2.Decrypt(s.privatekey, msgHexDecode, mode)
	if err != nil {
		log.Errorf("%v", err)
	}
	return planiText, err
}

func (s *sm) DecryptByte(msg []byte) ([]byte, error) {
	planiText, err := sm2.Decrypt(s.privatekey, msg, mode)
	return planiText, err
}

func (s *sm) GetPrivateKey() string {
	d := fmt.Sprintf("%x", s.privatekey.D)
	return d
}

func (s *sm) GetPubliKey() string {
	xy := fmt.Sprintf("04%x%x", s.publickey.X, s.publickey.Y)
	return xy
}

func (s *sm) Sign(msg string) (sign []byte, ok bool) {
	sign, err := s.privatekey.Sign(rand.Reader, []byte(msg), nil)
	if err != nil {
		log.Errorf("[sm.Sign] %s", err)
		return sign, false
	}
	return sign, true
}

func (s *sm) Verify(msg, sign string) bool {
	return s.publickey.Verify([]byte(msg), []byte(sign)) //sm2验签
}

func (s *sm) Encrypt(msg []byte) ([]byte, error) {
	encrptData, err := sm2.Encrypt(&s.privatekey.PublicKey, msg, rand.Reader, mode)
	return encrptData, err
}

// 这个库有点坑，有时候生成的key会有问题，自己加密后的东西自己都无法解密
// 自己先测试一下，如果有问题，重新再生成一个
func available(privatekey *sm2.PrivateKey) bool {
	var (
		value = "test"
	)
	encrypts1Str, err1 := privatekey.EncryptAsn1([]byte(value), rand.Reader)
	decrypts1Str, err2 := privatekey.DecryptAsn1(encrypts1Str)
	if err1 != nil || err2 != nil || string(decrypts1Str) != value {
		log.Debugf("%v %v", err1, err2)
		return false
	}
	log.Debugf("sm2 加密测试：test == %s\n", string(decrypts1Str))
	return true
}

func NewSm2Encrypt() *sm {
	var (
		priv  *sm2.PrivateKey
		err   error
		retry = 100
	)

	for i := 0; i < retry; i++ {
		priv, err = sm2.GenerateKey(rand.Reader) // 生成密钥对
		if err != nil {
			log.Debugf("[sm2.GenerateKey] err: %v", err)
			continue
		}
		if available(priv) {
			break
		}
	}
	log.Debugf("[sm2 PubliKey]: 04%x%x\n", priv.PublicKey.X, priv.PublicKey.Y)
	log.Debugf("[sm2 PrivateKey]: %x\n", priv.D)
	return &sm{
		privatekey: priv,
		publickey:  &priv.PublicKey,
	}
}
