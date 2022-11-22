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

package base

import (
	serve "dtstack.com/dtstack/easymatrix/addons/easyfiler/pkg/rpc-server"
	"dtstack.com/dtstack/easymatrix/go-common/log"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(srv *serve.Server) error {
	healthChecker := time.NewTimer(time.Second * 10)
	signals := make(chan os.Signal, 1)
	log.Infof("easyfiler start")
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go srv.Start()
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "FATAL %v\n", r)
		}
		healthChecker.Stop()
		exitGracefully()
	}()

LOOP:
	for {
		select {
		case sig := <-signals:
			fmt.Printf("Quit according to signal '%s'\n", sig.String())
			break LOOP
		case systemFailure := <-_SYSTEM_FAIL:
			if systemFailure.ExitCode > 0 {
				return fmt.Errorf("SYSTEM FAILURE: %d\nREASON: %s", systemFailure.ExitCode, systemFailure.Reason)
			}
		case <-healthChecker.C:
			checkSystemHealth()
		}
	}
	return nil
}

func exitGracefully() {

}

func checkSystemHealth() {

}
