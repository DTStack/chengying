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
	"encoding/hex"
	"fmt"
	"github.com/tjfoc/gmsm/sm2"
	"math/big"
	"testing"
)

func newSm2() (x, y string) {
	var (
		priv *sm2.PrivateKey
		err  error
	)
	priv, err = sm2.GenerateKey(rand.Reader) // 生成密钥对
	if err != nil {
		fmt.Printf("[sm2.GenerateKey] err: %v", err)
	}
	return fmt.Sprintf("04%x%x", priv.PublicKey.X, priv.PublicKey.Y), fmt.Sprintf("%x", priv.D)
}

func hexCode2Sm2(pb, pk string) *sm2.PrivateKey {
	x, _ := new(big.Int).SetString(pb[2:66], 16)
	y, _ := new(big.Int).SetString(pb[66:], 16)
	d, _ := new(big.Int).SetString(pk[:], 16)
	privatekey := &sm2.PrivateKey{
		PublicKey: sm2.PublicKey{
			Curve: sm2.P256Sm2(),
			X:     x,
			Y:     y,
		},
		D: d,
	}
	return privatekey
}
func encode(privatekey *sm2.PrivateKey) bool {
	var (
		value = "test"
	)
	encrypts1Str, err1 := privatekey.EncryptAsn1([]byte(value), rand.Reader)
	decrypts1Str, err2 := privatekey.DecryptAsn1(encrypts1Str)

	msgHexDecode, _ := hex.DecodeString("04f0804c7844e7c5560052f7b8a65c12ab3a314a11e6cb78e019b9ff398b479d71bb75692bc3559c66e8de04fb5a37e0dd550856eef59de372924ead8f01fe2df7429f827066638b4796eac026c2e8f81115e7ffc44228c51b646ba6180816c589c5caaa40aceacc26af1aaa05c60899")
	planiText, e := sm2.Decrypt(privatekey, msgHexDecode, 0)
	fmt.Println(string(planiText), e)
	enValue := string(decrypts1Str)
	if err1 != nil || err2 != nil {
		fmt.Printf("%s %v %v\n", enValue, err1, err2)
		return false
	}
	return true
}

func TestBenchMark(t *testing.T) {
	for k := 0; k < 100; k++ {
		xS, yS := newSm2()
		privatekey := hexCode2Sm2(xS, yS)
		if !encode(privatekey) {
			xS, yS := newSm2()
			nPrivatekey := hexCode2Sm2(xS, yS)
			if encode(nPrivatekey) {
				fmt.Println("retry success")
			} else {
				fmt.Println("retry fail")
			}
		}
	}
}

func TestName(t *testing.T) {

}
