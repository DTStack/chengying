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

package log

import (
	"dtstack.com/dtstack/easymatrix/addons/operator/cmd/options"
	"github.com/go-logr/logr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	sigszap "sigs.k8s.io/controller-runtime/pkg/log/zap"
	"time"
)

const (
	DebugLevel int8 = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
)

type noopInfoLogger struct{}

func (l *noopInfoLogger) Enabled() bool                   { return false }
func (l *noopInfoLogger) Info(_ string, _ ...interface{}) {}

var disabledInfoLogger = &noopInfoLogger{}

type infoLogger struct {
	lvl zapcore.Level
	l   *zap.Logger
}

//it's same with github.com/go-logr/zapr
//check the level, and use the zapcore's writer to write
func (i *infoLogger) Info(msg string, keysAndValues ...interface{}) {
	if checkedEntry := i.l.Check(i.lvl, msg); checkedEntry != nil {
		checkedEntry.Write(handleFields(i.l, keysAndValues)...)
	}
}

func (i *infoLogger) Enabled() bool {
	return true
}

type logger struct {
	infoLogger
	l *zap.Logger
}

func (l *logger) Error(err error, msg string, keysAndValues ...interface{}) {

	if checkedEntry := l.l.Check(zap.ErrorLevel, msg); checkedEntry != nil {
		checkedEntry.Write(handleFields(l.l, keysAndValues, zap.Error(err))...)
	}
}

//it's a little bit diffrenet with github.com/go-logr/zapr
//use the level directory to compare with the logger's level.
// if the level is higher than the level of logger, print
func (l *logger) V(level int) logr.InfoLogger {
	lvl := zapcore.Level(level)
	if l.l.Core().Enabled(lvl) {
		return &infoLogger{
			lvl: lvl,
			l:   l.infoLogger.l,
		}
	}
	return disabledInfoLogger
}

func (l *logger) WithValues(keysAndValues ...interface{}) logr.Logger {
	fields := handleFields(l.l, keysAndValues)
	el := l.l.With(fields...)
	il := l.infoLogger.l.With(fields...)
	return l.newlogger(il, el)
}

func (l *logger) WithName(name string) logr.Logger {
	el := l.l.Named(name)
	il := l.infoLogger.l.Named(name)
	return l.newlogger(il, el)
}

func (l *logger) newlogger(il, el *zap.Logger) logr.Logger {
	return &logger{
		infoLogger: infoLogger{
			l.infoLogger.lvl,
			il,
		},
		l: el,
	}
}

func newOptions(options *options.Options) []zap.Option {
	zops := []zap.Option{}
	zops = append(zops, zap.AddStacktrace(zap.ErrorLevel))
	if options.LogDebug {
		zops = append(zops, zap.Development())
	}
	return zops
}

//the logger is set two zap-logger, one is finally redirect to stdout, one is finally redirect to stderr
func NewLogger(options *options.Options) logr.Logger {

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = TimeEncode(options.LogTime)
	stdout := zapcore.AddSync(os.Stdout)
	stderr := zapcore.AddSync(os.Stderr)
	lvl := zap.NewAtomicLevelAt(zap.InfoLevel)
	if options.LogDebug {
		lvl = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	il := zap.New(zapcore.NewCore(&sigszap.KubeAwareEncoder{
		Encoder: zapcore.NewConsoleEncoder(config),
		Verbose: options.LogDebug,
	}, stdout, lvl))

	el := zap.New(zapcore.NewCore(&sigszap.KubeAwareEncoder{
		Encoder: zapcore.NewConsoleEncoder(config),
		Verbose: options.LogDebug,
	}, stderr, lvl))

	zops := newOptions(options)
	il = il.WithOptions(zops...)
	el = el.WithOptions(zops...)
	return &logger{
		infoLogger: infoLogger{
			lvl: lvl.Level(),
			l:   il,
		},
		l: el,
	}
}

func TimeEncode(layout string) func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(layout))
	}
}

//it's same with github.com/go-logr/zapr
//it's means the args is one key, one value. and switch it to zap.field
func handleFields(l *zap.Logger, args []interface{}, additional ...zap.Field) []zap.Field {
	// a slightly modified version of zap.SugaredLogger.sweetenFields
	if len(args) == 0 {
		// fast-return if we have no suggared fields.
		return additional
	}

	// unlike Zap, we can be pretty sure users aren't passing structured
	// fields (since logr has no concept of that), so guess that we need a
	// little less space.
	fields := make([]zap.Field, 0, len(args)/2+len(additional))
	for i := 0; i < len(args); {
		// check just in case for strongly-typed Zap fields, which is illegal (since
		// it breaks implementation agnosticism), so we can give a better error message.
		if _, ok := args[i].(zap.Field); ok {
			l.DPanic("strongly-typed Zap Field passed to logr", zap.Any("zap field", args[i]))
			break
		}

		// make sure this isn't a mismatched key
		if i == len(args)-1 {
			l.DPanic("odd number of arguments passed as key-value pairs for logging", zap.Any("ignored key", args[i]))
			break
		}

		// process a key-value pair,
		// ensuring that the key is a string
		key, val := args[i], args[i+1]
		keyStr, isString := key.(string)
		if !isString {
			// if the key isn't a string, DPanic and stop logging
			l.DPanic("non-string key argument passed to logging, ignoring all later arguments", zap.Any("invalid key", key))
			break
		}

		fields = append(fields, zap.Any(keyStr, val))
		i += 2
	}

	return append(fields, additional...)
}
