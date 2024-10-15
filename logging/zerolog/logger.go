// Copyright 2024 CloudWeGo Authors.
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

package zerolog

import (
	"context"
	cwotelzero "github.com/cloudwego-contrib/cwgo-pkg/telemetry/instrumentation/otelzerolog"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/rs/zerolog"
	"io"
)

var _ klog.FullLogger = (*Logger)(nil)

// Ref to https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/logs/README.md#json-formats
const (
	traceIDKey    = "trace_id"
	spanIDKey     = "span_id"
	traceFlagsKey = "trace_flags"
)

type Logger struct {
	l cwotelzero.Logger
}

func NewLogger(opts ...Option) *Logger {
	return &Logger{
		l: *cwotelzero.NewLogger(opts...),
	}
}

func (l *Logger) Logger() *zerolog.Logger {
	logger := l.l.Logger.GetLogger()
	return &logger
}

// Log log using zerolog logger with specified level
func (l *Logger) Log(level klog.Level, kvs ...any) {
	l.l.Log(levelSwitch(level), kvs...)
}

// Logf log using zerolog logger with specified level and formatting
func (l *Logger) Logf(level klog.Level, format string, kvs ...any) {
	l.l.Logf(levelSwitch(level), format, kvs...)
}

// CtxLogf log with logger associated with context.
// If no logger is associated, DefaultContextLogger is used, unless DefaultContextLogger is nil, in which case a disabled logger is used.
func (l *Logger) CtxLogf(level klog.Level, ctx context.Context, format string, kvs ...any) {
	l.l.CtxLogf(levelSwitch(level), ctx, format, kvs...)
}

// Trace logs a message at trace level.
func (l *Logger) Trace(v ...any) {
	l.Log(klog.LevelTrace, v...)
}

// Debug logs a message at debug level.
func (l *Logger) Debug(v ...any) {
	l.Log(klog.LevelDebug, v...)
}

// Info logs a message at info level.
func (l *Logger) Info(v ...any) {
	l.Log(klog.LevelInfo, v...)
}

// Notice logs a message at notice level.
func (l *Logger) Notice(v ...any) {
	l.Log(klog.LevelNotice, v...)
}

// Warn logs a message at warn level.
func (l *Logger) Warn(v ...any) {
	l.Log(klog.LevelWarn, v...)
}

// Error logs a message at error level.
func (l *Logger) Error(v ...any) {
	l.Log(klog.LevelError, v...)
}

// Fatal logs a message at fatal level.
func (l *Logger) Fatal(v ...any) {
	l.Log(klog.LevelFatal, v...)
}

// Tracef logs a formatted message at trace level.
func (l *Logger) Tracef(format string, v ...any) {
	l.Logf(klog.LevelTrace, format, v...)
}

// Debugf logs a formatted message at debug level.
func (l *Logger) Debugf(format string, v ...any) {
	l.Logf(klog.LevelDebug, format, v...)
}

// Infof logs a formatted message at info level.
func (l *Logger) Infof(format string, v ...any) {
	l.Logf(klog.LevelInfo, format, v...)
}

// Noticef logs a formatted message at notice level.
func (l *Logger) Noticef(format string, v ...any) {
	l.Logf(klog.LevelWarn, format, v...)
}

// Warnf logs a formatted message at warn level.
func (l *Logger) Warnf(format string, v ...any) {
	l.Logf(klog.LevelWarn, format, v...)
}

// Errorf logs a formatted message at error level.
func (l *Logger) Errorf(format string, v ...any) {
	l.Logf(klog.LevelError, format, v...)
}

// Fatalf logs a formatted message at fatal level.
func (l *Logger) Fatalf(format string, v ...any) {
	l.Logf(klog.LevelError, format, v...)
}

// CtxTracef logs a message at trace level with logger associated with context.
// If no logger is associated, DefaultContextLogger is used, unless DefaultContextLogger is nil, in which case a disabled logger is used.
func (l *Logger) CtxTracef(ctx context.Context, format string, v ...any) {
	l.CtxLogf(klog.LevelTrace, ctx, format, v...)
}

// CtxDebugf logs a message at debug level with logger associated with context.
// If no logger is associated, DefaultContextLogger is used, unless DefaultContextLogger is nil, in which case a disabled logger is used.
func (l *Logger) CtxDebugf(ctx context.Context, format string, v ...any) {
	l.CtxLogf(klog.LevelDebug, ctx, format, v...)
}

// CtxInfof logs a message at info level with logger associated with context.
// If no logger is associated, DefaultContextLogger is used, unless DefaultContextLogger is nil, in which case a disabled logger is used.
func (l *Logger) CtxInfof(ctx context.Context, format string, v ...any) {
	l.CtxLogf(klog.LevelInfo, ctx, format, v...)
}

// CtxNoticef logs a message at notice level with logger associated with context.
// If no logger is associated, DefaultContextLogger is used, unless DefaultContextLogger is nil, in which case a disabled logger is used.
func (l *Logger) CtxNoticef(ctx context.Context, format string, v ...any) {
	l.CtxLogf(klog.LevelNotice, ctx, format, v...)
}

// CtxWarnf logs a message at warn level with logger associated with context.
// If no logger is associated, DefaultContextLogger is used, unless DefaultContextLogger is nil, in which case a disabled logger is used.
func (l *Logger) CtxWarnf(ctx context.Context, format string, v ...any) {
	l.CtxLogf(klog.LevelWarn, ctx, format, v...)
}

// CtxErrorf logs a message at error level with logger associated with context.
// If no logger is associated, DefaultContextLogger is used, unless DefaultContextLogger is nil, in which case a disabled logger is used.
func (l *Logger) CtxErrorf(ctx context.Context, format string, v ...any) {
	l.CtxLogf(klog.LevelError, ctx, format, v...)
}

// CtxFatalf logs a message at fatal level with logger associated with context.
// If no logger is associated, DefaultContextLogger is used, unless DefaultContextLogger is nil, in which case a disabled logger is used.
func (l *Logger) CtxFatalf(ctx context.Context, format string, v ...any) {
	l.CtxLogf(klog.LevelFatal, ctx, format, v...)
}

func (l *Logger) SetLevel(level klog.Level) {
	l.l.SetLevel(levelSwitch(level))
}

func (l *Logger) SetOutput(writer io.Writer) {
	l.l.SetOutput(writer)
}

func levelSwitch(klvel klog.Level) hlog.Level {
	var lv hlog.Level
	switch klvel {
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
	return lv
}
