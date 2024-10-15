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
	"github.com/cloudwego-contrib/cwgo-pkg/telemetry/instrumentation/otellogrus"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"io"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/sirupsen/logrus"
)

var _ klog.FullLogger = (*Logger)(nil)

type Logger struct {
	l otellogrus.Logger
}

func NewLogger(opts ...Option) *Logger {

	return &Logger{
		l: *otellogrus.NewLogger(opts...),
	}
}

func (l *Logger) Logger() *logrus.Logger {
	return l.l.Logger()
}

func (l *Logger) Trace(v ...interface{}) {
	l.l.Trace(v...)
}

func (l *Logger) Debug(v ...interface{}) {
	l.l.Debug(v...)
}

func (l *Logger) Info(v ...interface{}) {
	l.l.Info(v...)
}

func (l *Logger) Notice(v ...interface{}) {
	l.l.Warn(v...)
}

func (l *Logger) Warn(v ...interface{}) {
	l.l.Warn(v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.l.Error(v...)
}

func (l *Logger) Fatal(v ...interface{}) {
	l.l.Fatal(v...)
}

func (l *Logger) Tracef(format string, v ...interface{}) {
	l.l.Tracef(format, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.l.Debugf(format, v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.l.Infof(format, v...)
}

func (l *Logger) Noticef(format string, v ...interface{}) {
	l.l.Warnf(format, v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.l.Warnf(format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.l.Errorf(format, v...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.l.Fatalf(format, v...)
}

func (l *Logger) CtxTracef(ctx context.Context, format string, v ...interface{}) {
	l.l.CtxTracef(ctx, format, v...)
}

func (l *Logger) CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	l.l.CtxDebugf(ctx, format, v...)
}

func (l *Logger) CtxInfof(ctx context.Context, format string, v ...interface{}) {
	l.l.CtxInfof(ctx, format, v...)
}

func (l *Logger) CtxNoticef(ctx context.Context, format string, v ...interface{}) {
	l.l.CtxNoticef(ctx, format, v...)
}

func (l *Logger) CtxWarnf(ctx context.Context, format string, v ...interface{}) {
	l.l.CtxNoticef(ctx, format, v...)
}

func (l *Logger) CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	l.l.CtxNoticef(ctx, format, v...)
}

func (l *Logger) CtxFatalf(ctx context.Context, format string, v ...interface{}) {
	l.l.CtxFatalf(ctx, format, v...)
}

func (l *Logger) SetLevel(level klog.Level) {
	var lv hlog.Level
	switch level {
	case klog.LevelTrace:
		lv = hlog.LevelTrace
	case klog.LevelDebug:
		lv = hlog.LevelDebug
	case klog.LevelInfo:
		lv = hlog.LevelInfo
	case klog.LevelWarn:
		lv = hlog.LevelWarn
	case klog.LevelNotice:
		lv = hlog.LevelWarn

	case klog.LevelError:
		lv = hlog.LevelError
	case klog.LevelFatal:
		lv = hlog.LevelFatal
	default:
		lv = hlog.LevelWarn
	}
	l.l.SetLevel(lv)
}

func (l *Logger) SetOutput(writer io.Writer) {
	l.l.SetOutput(writer)
}
