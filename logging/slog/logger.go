// Copyright 2022 CloudWeGo Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package slog

import (
	"context"
	"io"
	"log/slog"

	"github.com/cloudwego/kitex/pkg/klog"
)

const (
	LevelTrace  = slog.Level(-8)
	LevelNotice = slog.Level(2)
	LevelFatal  = slog.Level(12)
)

var _ klog.FullLogger = (*Logger)(nil)

type Logger struct {
	l      *slog.Logger
	config *config
}

func NewLogger(opts ...Option) *Logger {

	config := defaultConfig()

	for _, opt := range opts {
		opt.apply(config)
	}

	logger := slog.New(NewTraceHandler(config.coreConfig.writer, &config.coreConfig.opt, config.traceConfig))
	return &Logger{
		l:      logger,
		config: config,
	}
}

func (l *Logger) Log(level klog.Level, msg string, kvs ...interface{}) {
	logger := l.l.With()
	logger.Log(context.TODO(), tranSLevel(level), msg, kvs...)
}

func (l *Logger) Logf(level klog.Level, format string, kvs ...interface{}) {
	logger := l.l.With()
	msg := getMessage(format, kvs)
	logger.Log(context.TODO(), tranSLevel(level), msg)
}

func (l *Logger) CtxLogf(level klog.Level, ctx context.Context, format string, kvs ...interface{}) {
	logger := l.l.With()
	msg := getMessage(format, kvs)
	logger.Log(ctx, tranSLevel(level), msg)
}

func (l *Logger) Trace(v ...interface{}) {
	msg, kvs := tranAtrr(v)
	l.Log(klog.LevelTrace, msg, kvs...)
}

func (l *Logger) Debug(v ...interface{}) {
	msg, kvs := tranAtrr(v)
	l.Log(klog.LevelDebug, msg, kvs...)
}

func (l *Logger) Info(v ...interface{}) {
	msg, kvs := tranAtrr(v)
	l.Log(klog.LevelInfo, msg, kvs...)
}

func (l *Logger) Notice(v ...interface{}) {
	msg, kvs := tranAtrr(v)
	l.Log(klog.LevelNotice, msg, kvs...)
}

func (l *Logger) Warn(v ...interface{}) {
	msg, kvs := tranAtrr(v)
	l.Log(klog.LevelWarn, msg, kvs...)
}

func (l *Logger) Error(v ...interface{}) {
	msg, kvs := tranAtrr(v)
	l.Log(klog.LevelError, msg, kvs...)
}

func (l *Logger) Fatal(v ...interface{}) {
	msg, kvs := tranAtrr(v)
	l.Log(klog.LevelFatal, msg, kvs...)
}

func (l *Logger) Tracef(format string, v ...interface{}) {
	l.Logf(klog.LevelTrace, format, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Logf(klog.LevelDebug, format, v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.Logf(klog.LevelInfo, format, v...)
}

func (l *Logger) Noticef(format string, v ...interface{}) {
	l.Logf(klog.LevelInfo, format, v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.Logf(klog.LevelWarn, format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Logf(klog.LevelError, format, v...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Logf(klog.LevelFatal, format, v...)
}

func (l *Logger) CtxTracef(ctx context.Context, format string, v ...interface{}) {
	l.CtxLogf(klog.LevelDebug, ctx, format, v...)
}

func (l *Logger) CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	l.CtxLogf(klog.LevelDebug, ctx, format, v...)
}

func (l *Logger) CtxInfof(ctx context.Context, format string, v ...interface{}) {
	l.CtxLogf(klog.LevelInfo, ctx, format, v...)
}

func (l *Logger) CtxNoticef(ctx context.Context, format string, v ...interface{}) {
	l.CtxLogf(klog.LevelWarn, ctx, format, v...)
}

func (l *Logger) CtxWarnf(ctx context.Context, format string, v ...interface{}) {
	l.CtxLogf(klog.LevelWarn, ctx, format, v...)
}

func (l *Logger) CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	l.CtxLogf(klog.LevelError, ctx, format, v...)
}

func (l *Logger) CtxFatalf(ctx context.Context, format string, v ...interface{}) {
	l.CtxLogf(klog.LevelFatal, ctx, format, v...)
}

func (l *Logger) SetLevel(level klog.Level) {
	lvl := tranSLevel(level)
	l.config.coreConfig.level.Set(lvl)
}

func (l *Logger) SetOutput(writer io.Writer) {

	log := slog.New(NewTraceHandler(writer, &l.config.coreConfig.opt, l.config.traceConfig))
	l.config.coreConfig.writer = writer
	l.l = log
}
