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
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/natefinch/lumberjack"
)

const (
	DEFAULT_LOG = iota
	DEFAULT_ERR_LOG
)

var debug = true

func SetDebug(b bool) {
	debug = b
}

var _LOGGERS = map[int]*log.Logger{}

func Infof(format string, args ...interface{}) {
	_LOGGERS[DEFAULT_LOG].Output(2, fmt.Sprintf(format, args...))
}

func Debugf(format string, args ...interface{}) {
	if debug {
		_LOGGERS[DEFAULT_LOG].Output(2, fmt.Sprintf(format, args...))
	}
}

func Errorf(format string, args ...interface{}) {
	_LOGGERS[DEFAULT_ERR_LOG].Output(2, fmt.Sprintf(format, args...))
}

func ConfigureLogger(dir string, maxSize, maxBackups, maxAge int) error {
	makeLogger := func(prefix string, tag string, flag int) *log.Logger {
		return log.New(&lumberjack.Logger{
			Filename:   filepath.Join(dir, prefix+".log"),
			MaxSize:    maxSize,
			MaxBackups: maxBackups,
			MaxAge:     maxAge,
		}, tag, flag)
	}

	// mkdir -p
	os.MkdirAll(dir, os.FileMode(0755))

	if debug {
		_LOGGERS[DEFAULT_LOG] = makeLogger("agent", "AGENT-DEBUG:", log.LstdFlags|log.Lshortfile)
	} else {
		_LOGGERS[DEFAULT_LOG] = makeLogger("agent", "AGENT:", log.LstdFlags)
	}
	_LOGGERS[DEFAULT_ERR_LOG] = makeLogger("agent-error", "AGENT-ERROR:", log.LstdFlags|log.Lshortfile)

	return nil
}
