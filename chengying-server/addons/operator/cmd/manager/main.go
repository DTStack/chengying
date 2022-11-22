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

package main

import (
	"dtstack.com/dtstack/easymatrix/addons/operator/cmd/options"
	statcklog "dtstack.com/dtstack/easymatrix/addons/operator/log"
	"dtstack.com/dtstack/easymatrix/addons/operator/pkg/apis"
	"dtstack.com/dtstack/easymatrix/addons/operator/pkg/controller"
	"flag"
	"fmt"
	"k8s.io/klog"
	"os"
	"runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	"github.com/spf13/pflag"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

var log = logf.Log.WithName("cmd")

func startPrint() {
	log.Info(fmt.Sprintf("Go Version: %s", runtime.Version()))
	log.Info(fmt.Sprintf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH))
	log.Info("operator start")
}

func main() {

	fs := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	opt := options.GetOptions()
	opt.AddToFlagSet(fs)
	fs.Parse(os.Args)

	//set klog log to stderr
	klogFlagSet := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(klogFlagSet)

	klogFlagSet.Set("logtostderr", "true")
	klogFlagSet.Set("alsologtostderr", "false")

	logf.SetLogger(statcklog.NewLogger(opt))

	startPrint()

	namespace, err := getNamespace(opt.WatchNamespace)
	if err != nil {
		log.Error(err, "Failed to get watch namespace")
		os.Exit(1)
	}
	// use incluster config to talk with api-server
	cfg, err := config.GetConfig()
	if err != nil {
		log.Error(err, "unable to get rest configuration")
	}

	// would elect a leader to start the informer,
	//if the operator pod is not leader any more, the pod will exist and than restart
	mgr, err := manager.New(cfg, manager.Options{
		LeaderElection:          true,
		LeaderElectionNamespace: namespace,
		LeaderElectionID:        opt.ElectionLockName,
		Namespace:               namespace,
	})
	if err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	// Setup Scheme for all resources
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	// Setup all Controllers
	if err := controller.AddToManager(mgr); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	// Start the Cmd
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Error(err, "Manager exited non-zero")
		os.Exit(1)
	}
}

// use the namespace from option if is set in commandline
func getNamespace(namespace string) (string, error) {
	if len(namespace) != 0 {
		return namespace, nil
	}
	ns, found := os.LookupEnv("WATCH_NAMESPACE")
	if !found {
		return "", fmt.Errorf("%s must be set", "WATCH_NAMESPACE")
	}
	return ns, nil
}
