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

package logrus

import (
	"context"
	"io"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/sirupsen/logrus"
)

var _ klog.FullLogger = (*Logger)(nil)

type Logger struct {
	L *logrus.Logger
}

func NewLogger(opts ...Option) *Logger {
	cfg := defaultConfig()

	// apply options
	for _, opt := range opts {
		opt.apply(cfg)
	}

	// default trace hooks
	cfg.hooks = append(cfg.hooks, NewTraceHook(cfg.traceHookConfig))

	// attach hook
	for _, hook := range cfg.hooks {
		cfg.logger.AddHook(hook)
	}

	return &Logger{
		L: cfg.logger,
	}
}

func (l *Logger) Logger() *logrus.Logger {
	return l.L
}

func (l *Logger) Trace(v ...interface{}) {
	l.L.Trace(v...)
}

func (l *Logger) Debug(v ...interface{}) {
	l.L.Debug(v...)
}

func (l *Logger) Info(v ...interface{}) {
	l.L.Info(v...)
}

func (l *Logger) Notice(v ...interface{}) {
	l.L.Warn(v...)
}

func (l *Logger) Warn(v ...interface{}) {
	l.L.Warn(v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.L.Error(v...)
}

func (l *Logger) Fatal(v ...interface{}) {
	l.L.Fatal(v...)
}

func (l *Logger) Tracef(format string, v ...interface{}) {
	l.L.Tracef(format, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.L.Debugf(format, v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.L.Infof(format, v...)
}

func (l *Logger) Noticef(format string, v ...interface{}) {
	l.L.Warnf(format, v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.L.Warnf(format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.L.Errorf(format, v...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.L.Fatalf(format, v...)
}

func (l *Logger) CtxTracef(ctx context.Context, format string, v ...interface{}) {
	l.L.WithContext(ctx).Tracef(format, v...)
}

func (l *Logger) CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	l.L.WithContext(ctx).Debugf(format, v...)
}

func (l *Logger) CtxInfof(ctx context.Context, format string, v ...interface{}) {
	l.L.WithContext(ctx).Infof(format, v...)
}

func (l *Logger) CtxNoticef(ctx context.Context, format string, v ...interface{}) {
	l.L.WithContext(ctx).Warnf(format, v...)
}

func (l *Logger) CtxWarnf(ctx context.Context, format string, v ...interface{}) {
	l.L.WithContext(ctx).Warnf(format, v...)
}

func (l *Logger) CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	l.L.WithContext(ctx).Errorf(format, v...)
}

func (l *Logger) CtxFatalf(ctx context.Context, format string, v ...interface{}) {
	l.L.WithContext(ctx).Fatalf(format, v...)
}

func (l *Logger) SetLevel(level klog.Level) {
	var lv logrus.Level
	switch level {
	case klog.LevelTrace:
		lv = logrus.TraceLevel
	case klog.LevelDebug:
		lv = logrus.DebugLevel
	case klog.LevelInfo:
		lv = logrus.InfoLevel
	case klog.LevelWarn, klog.LevelNotice:
		lv = logrus.WarnLevel
	case klog.LevelError:
		lv = logrus.ErrorLevel
	case klog.LevelFatal:
		lv = logrus.FatalLevel
	default:
		lv = logrus.WarnLevel
	}
	l.L.SetLevel(lv)
}

func (l *Logger) SetOutput(writer io.Writer) {
	l.L.SetOutput(writer)
}
