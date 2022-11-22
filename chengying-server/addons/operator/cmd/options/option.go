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

package options

import (
	"github.com/spf13/pflag"
	"sync"
)

var opt *Options
var once sync.Once

const (
	DefaultLogDebug         = false
	DefaultLogTiemLayout    = "2006-01-02 15:04:05"
	DefaultElectionLockName = "dtstack-operator"
)

//option will be used
type Options struct {
	LogDebug         bool
	LogTime          string
	WatchNamespace   string
	ElectionLockName string
}

func GetOptions() *Options {
	once.Do(func() {
		opt = &Options{
			LogDebug:         DefaultLogDebug,
			LogTime:          DefaultLogTiemLayout,
			ElectionLockName: DefaultElectionLockName,
		}
	})
	return opt
}

func (opt *Options) AddToFlagSet(fs *pflag.FlagSet) {
	fs.BoolVar(&opt.LogDebug, "log-debug", DefaultLogDebug, "set logger in debug mode")
	fs.StringVar(&opt.LogTime, "log-time", DefaultLogTiemLayout, "set log time layout format")
	fs.StringVar(&opt.WatchNamespace, "watch-namespace", "", "namespace the operator is watched")
	fs.StringVar(&opt.ElectionLockName, "election-lock-name", DefaultElectionLockName, "name of configmap lock name")
}
