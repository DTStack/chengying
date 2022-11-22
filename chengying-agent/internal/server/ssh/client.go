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

package sshs

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"net"
	"strconv"
	"time"

	"easyagent/internal/server/log"
	. "easyagent/internal/server/tracy"
	"golang.org/x/crypto/ssh"
)

var (
	ErrIsNull = errors.New("ssh client is nil")
)

//mode
//1 user pwd
//2 user pem
//3 user rsa
//4 user dsa
//5 user ecdsa

type SshParam struct {
	Host string `json:"host"`
	Port int    `json:"int"`
	User string `json:"user"`
	Pass string `json:"pass"`
	Pk   string `json:"pk"`
	Mode int    `json:"mode"`
	Cmd  string `json:"cmd"`
}

type sshClient struct {
	sClient  *ssh.Client
	sSession *ssh.Session
}

func CreateWithParam(param *SshParam) (*sshClient, error) {
	if param == nil {
		return nil, errors.New("ssh connect param is nil")
	}
	switch param.Mode {
	case 1:
		return CreateWithUserPwd(param.Host, param.User, param.Pass, param.Port)
	case 2:
		block, _ := pem.Decode([]byte(param.Pk))
		if block != nil {
			log.Debugf("CreateWithParam pem block type: %v", block.Type)
			InstallProgressLog("[CreateWithParam] pem block type: %v", block.Type)
		}
		return CreateWithUserPem(param.Host, param.User, param.Pk, param.Port)
	case 3:
		return CreateWithUserRsa(param.Host, param.User, param.Pk, param.Port)
	case 4:
		return CreateWithUserDsa(param.Host, param.User, param.Pk, param.Port)
	case 5:
		return CreateWithUserEcdsa(param.Host, param.User, param.Pk, param.Port)
	default:
		return nil, errors.New("ssh connect param with unsupport mode")
	}
}

func CreateWithUserPwd(host, user, pass string, port int) (*sshClient, error) {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(pass)},
		//验证服务端，不做验证返回nil就可以
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: time.Second * 3,
	}

	cli, err := ssh.Dial("tcp", net.JoinHostPort(host, strconv.Itoa(port)), sshConfig)
	if err != nil {
		log.Errorf("connect ssh with pwd err:%v", err)
		return nil, err
	}

	session, err := cli.NewSession()
	if err != nil {
		log.Errorf("new sshsession with pwd err:%v", err)
		return nil, err
	}
	/*
		modes := ssh.TerminalModes{
			ssh.ECHO:          0,
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		}
		// request tty -- fixes error with hosts that use
		// "Defaults requiretty" in /etc/sudoers - I'm looking at you RedHat
		if err = session.RequestPty("xterm", 24, 80, modes); err != nil {
			log.Errorf("RequestPty with pwd err:%v", err)
			return nil, err
		}
	*/

	return &sshClient{sClient: cli, sSession: session}, nil
}

func CreateWithUserPem(host, user, pem string, port int) (*sshClient, error) {
	signer, err := ssh.ParsePrivateKey([]byte(pem))
	if err != nil {
		log.Errorf("parse key failed:%v", err)
		return nil, err
	}
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
		//验证服务端，不做验证返回nil就可以
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: time.Second * 3,
	}

	cli, err := ssh.Dial("tcp", net.JoinHostPort(host, strconv.Itoa(port)), sshConfig)
	if err != nil {
		log.Errorf("connect ssh with pem err:%v", err)
		return nil, err
	}
	session, err := cli.NewSession()
	if err != nil {
		log.Errorf("new sshsession with pem err:%v", err)
		return nil, err
	}
	/*
		modes := ssh.TerminalModes{
			ssh.ECHO:          0,
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		}
		// request tty -- fixes error with hosts that use
		// "Defaults requiretty" in /etc/sudoers - I'm looking at you RedHat
		if err = session.RequestPty("xterm", 24, 80, modes); err != nil {
			log.Errorf("RequestPty with pem err:%v", err)
			return nil, err
		}
	*/
	return &sshClient{sClient: cli, sSession: session}, nil
}

func CreateWithUserRsa(host, user, sa string, port int) (*sshClient, error) {
	key, err := x509.ParsePKCS1PrivateKey([]byte(sa))
	if err != nil {
		log.Errorf("ParsePKCS1PrivateKey failed:%v", err)
		return nil, err
	}
	signer, err := ssh.NewSignerFromKey(key)
	if err != nil {
		log.Errorf("ParsePKCS1PrivateKey NewSignerFromKey failed:%v", err)
		return nil, err
	}
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
		//验证服务端，不做验证返回nil就可以
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: time.Second * 3,
	}
	cli, err := ssh.Dial("tcp", net.JoinHostPort(host, strconv.Itoa(port)), sshConfig)
	if err != nil {
		log.Errorf("connect ssh with rsa err:%v", err)
		return nil, err
	}

	session, err := cli.NewSession()
	if err != nil {
		log.Errorf("new sshsession with rsa err:%v", err)
		return nil, err
	}
	/*
		modes := ssh.TerminalModes{
			ssh.ECHO:          0,
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		}
		// request tty -- fixes error with hosts that use
		// "Defaults requiretty" in /etc/sudoers - I'm looking at you RedHat
		if err = session.RequestPty("xterm", 24, 80, modes); err != nil {
			log.Errorf("terminal.GetSize with rsa err:%v", err)
			return nil, err
		}
	*/
	return &sshClient{sClient: cli, sSession: session}, nil
}

func CreateWithUserDsa(host, user, dsa string, port int) (*sshClient, error) {
	key, err := ssh.ParseDSAPrivateKey([]byte(dsa))
	if err != nil {
		log.Errorf("ParseDSAPrivateKey failed:%v", err)
		return nil, err
	}
	signer, err := ssh.NewSignerFromKey(key)
	if err != nil {
		log.Errorf("ParseDSAPrivateKey NewSignerFromKey failed:%v", err)
		return nil, err
	}
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
		//验证服务端，不做验证返回nil就可以
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: time.Second * 3,
	}
	cli, err := ssh.Dial("tcp", net.JoinHostPort(host, strconv.Itoa(port)), sshConfig)
	if err != nil {
		log.Errorf("connect ssh with dsa err:%v", err)
		return nil, err
	}

	session, err := cli.NewSession()
	if err != nil {
		log.Errorf("new sshsession with dsa err:%v", err)
		return nil, err
	}
	/*
		modes := ssh.TerminalModes{
			ssh.ECHO:          0,
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		}
		// request tty -- fixes error with hosts that use
		// "Defaults requiretty" in /etc/sudoers - I'm looking at you RedHat
		if err = session.RequestPty("xterm", 24, 80, modes); err != nil {
			log.Errorf("RequestPty with dsa err:%v", err)
			return nil, err
		}
	*/
	return &sshClient{sClient: cli, sSession: session}, nil
}

func CreateWithUserEcdsa(host, user, dsa string, port int) (*sshClient, error) {
	key, err := x509.ParseECPrivateKey([]byte(dsa))
	if err != nil {
		log.Errorf("ParseECPrivateKey failed:%v", err)
		return nil, err
	}
	signer, err := ssh.NewSignerFromKey(key)
	if err != nil {
		log.Errorf("ParseECPrivateKey NewSignerFromKey failed:%v", err)
		return nil, err
	}
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
		//验证服务端，不做验证返回nil就可以
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: time.Second * 3,
	}
	cli, err := ssh.Dial("tcp", net.JoinHostPort(host, strconv.Itoa(port)), sshConfig)
	if err != nil {
		log.Errorf("connect ssh with ecdsa err:%v", err)
		return nil, err
	}

	session, err := cli.NewSession()
	if err != nil {
		log.Errorf("new sshsession with ecdsa err:%v", err)
		return nil, err
	}
	/*
		modes := ssh.TerminalModes{
			ssh.ECHO:          0,
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		}
		// request tty -- fixes error with hosts that use
		// "Defaults requiretty" in /etc/sudoers - I'm looking at you RedHat
		if err = session.RequestPty("xterm", 24, 80, modes); err != nil {
			log.Errorf("RequestPty with ecdsa err:%v", err)
			return nil, err
		}
	*/
	return &sshClient{sClient: cli, sSession: session}, nil
}

func (sh *sshClient) isValid() bool {
	if sh.sSession == nil || sh.sClient == nil {
		return false
	}
	return true
}

func (sh *sshClient) RunSync(cmd string) (string, error) {
	if !sh.isValid() {
		return "", ErrIsNull
	}
	log.Debugf("Run cmd sync:%v", cmd)
	out, err := sh.sSession.CombinedOutput(cmd)
	if err != nil {
		log.Errorf("RunSync Cmd error: %v, %v", err.Error(), out)
		return string(out), err
	}
	log.Debugf("RunSync cmd output::%v", string(out))
	return string(out), nil
}

func (sh *sshClient) Run(cmd string) error {
	if !sh.isValid() {
		return ErrIsNull
	}
	log.Debugf("Run cmd:%v", cmd)
	err := sh.sSession.Start(cmd)
	if err != nil {
		log.Errorf("Run Cmd error: %v", err.Error())
		return err
	}
	return nil
}

func (sh *sshClient) Close() {
	if !sh.isValid() {
		return
	}
	if sh.sClient != nil {
		err := sh.sClient.Close()
		if err != nil {
			log.Errorf("Unable to Close ssh client: %v", err.Error())
		}
	}
}
