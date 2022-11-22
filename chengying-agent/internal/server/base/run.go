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

package base

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"easyagent/internal/server/model"
	"easyagent/internal/server/publisher"
	"github.com/kataras/iris"
)

func Run() error {
	healthChecker := time.NewTimer(time.Second * 10)
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "FATAL!!! %v\n", r)
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
	if app := iris.Default(); app != nil && app.Shutdown != nil {
		fmt.Println("Stopping API server...")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		if err := app.Shutdown(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to stop API server: %v\n", err)
		}
	}

	if _RPCSERVER != nil {
		fmt.Println("Stopping RPC server...")
		_RPCSERVER.Stop()
	}

	if model.MYSQLDB != nil {
		fmt.Println("Disconnecting mysql-db...")
		if err := model.MYSQLDB.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to disconnect mysql-database: %v\n", err)
		}
	}

	//if model.UICDB != nil {
	//	fmt.Println("Disconnecting uic-db...")
	//	if err := model.UICDB.Close(); err != nil {
	//		fmt.Fprintf(os.Stderr, "Unable to disconnect uic-database: %v\n", err)
	//	}
	//}

	if publisher.Publish != nil {
		fmt.Println("Closing publish...")
		publisher.Publish.Close()
	}
}

func checkSystemHealth() {

}
