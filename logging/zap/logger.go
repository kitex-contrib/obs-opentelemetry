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

package zap

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/cloudwego/kitex/pkg/klog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ klog.FullLogger = (*Logger)(nil)

// Ref to https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/logs/README.md#json-formats
const (
	traceIDKey    = "trace_id"
	spanIDKey     = "span_id"
	traceFlagsKey = "trace_flags"
	logEventKey   = "log"
)

var (
	logSeverityTextKey = attribute.Key("otel.log.severity.text")
	logMessageKey      = attribute.Key("otel.log.message")
)

type Logger struct {
	*zap.SugaredLogger
	config *config
}

func NewLogger(opts ...Option) *Logger {
	config := defaultConfig()

	// apply options
	for _, opt := range opts {
		opt.apply(config)
	}

	logger := zap.New(
		zapcore.NewCore(config.coreConfig.enc, config.coreConfig.ws, config.coreConfig.lvl),
		config.zapOpts...)

	return &Logger{
		SugaredLogger: logger.Sugar(),
		config:        config,
	}
}

func (l *Logger) Log(level klog.Level, kvs ...interface{}) {
	logger := l.With()
	switch level {
	case klog.LevelTrace, klog.LevelDebug:
		logger.Debug(kvs...)
	case klog.LevelInfo:
		logger.Info(kvs...)
	case klog.LevelNotice, klog.LevelWarn:
		logger.Warn(kvs...)
	case klog.LevelError:
		logger.Error(kvs...)
	case klog.LevelFatal:
		logger.Fatal(kvs...)
	default:
		logger.Warn(kvs...)
	}
}

func (l *Logger) Logf(level klog.Level, format string, kvs ...interface{}) {
	logger := l.With()
	switch level {
	case klog.LevelTrace, klog.LevelDebug:
		logger.Debugf(format, kvs...)
	case klog.LevelInfo:
		logger.Infof(format, kvs...)
	case klog.LevelNotice, klog.LevelWarn:
		logger.Warnf(format, kvs...)
	case klog.LevelError:
		logger.Errorf(format, kvs...)
	case klog.LevelFatal:
		logger.Fatalf(format, kvs...)
	default:
		logger.Warnf(format, kvs...)
	}
}

func (l *Logger) CtxLogf(level klog.Level, ctx context.Context, format string, kvs ...interface{}) {
	var zlevel zapcore.Level
	var sl *zap.SugaredLogger

	span := trace.SpanFromContext(ctx)
	var traceKVs []interface{}
	if span.SpanContext().TraceID().IsValid() {
		traceKVs = append(traceKVs, traceIDKey, span.SpanContext().TraceID())
	}
	if span.SpanContext().SpanID().IsValid() {
		traceKVs = append(traceKVs, spanIDKey, span.SpanContext().SpanID())
	}
	if span.SpanContext().TraceFlags().IsSampled() {
		traceKVs = append(traceKVs, traceFlagsKey, span.SpanContext().TraceFlags())
	}
	if len(traceKVs) > 0 {
		sl = l.With(traceKVs...)
	} else {
		sl = l.With()
	}

	switch level {
	case klog.LevelDebug, klog.LevelTrace:
		zlevel = zap.DebugLevel
		sl.Debugf(format, kvs...)
	case klog.LevelInfo:
		zlevel = zap.InfoLevel
		sl.Infof(format, kvs...)
	case klog.LevelNotice, klog.LevelWarn:
		zlevel = zap.WarnLevel
		sl.Warnf(format, kvs...)
	case klog.LevelError:
		zlevel = zap.ErrorLevel
		sl.Errorf(format, kvs...)
	case klog.LevelFatal:
		zlevel = zap.FatalLevel
		sl.Fatalf(format, kvs...)
	default:
		zlevel = zap.WarnLevel
		sl.Warnf(format, kvs...)
	}

	if !span.IsRecording() {
		return
	}

	msg := getMessage(format, kvs)

	attrs := []attribute.KeyValue{
		logMessageKey.String(msg),
		logSeverityTextKey.String(OtelSeverityText(zlevel)),
	}
	span.AddEvent(logEventKey, trace.WithAttributes(attrs...))

	// set span status
	if zlevel <= l.config.traceConfig.errorSpanLevel {
		span.SetStatus(codes.Error, msg)
		span.RecordError(errors.New(msg), trace.WithStackTrace(l.config.traceConfig.recordStackTraceInSpan))
	}
}

func (l *Logger) Trace(v ...interface{}) {
	l.Log(klog.LevelTrace, v...)
}

func (l *Logger) Debug(v ...interface{}) {
	l.Log(klog.LevelDebug, v...)
}

func (l *Logger) Info(v ...interface{}) {
	l.Log(klog.LevelInfo, v...)
}

func (l *Logger) Notice(v ...interface{}) {
	l.Log(klog.LevelNotice, v...)
}

func (l *Logger) Warn(v ...interface{}) {
	l.Log(klog.LevelWarn, v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.Log(klog.LevelError, v...)
}

func (l *Logger) Fatal(v ...interface{}) {
	l.Log(klog.LevelFatal, v...)
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
	var lvl zapcore.Level
	switch level {
	case klog.LevelTrace, klog.LevelDebug:
		lvl = zap.DebugLevel
	case klog.LevelInfo:
		lvl = zap.InfoLevel
	case klog.LevelWarn, klog.LevelNotice:
		lvl = zap.WarnLevel
	case klog.LevelError:
		lvl = zap.ErrorLevel
	case klog.LevelFatal:
		lvl = zap.FatalLevel
	default:
		lvl = zap.WarnLevel
	}
	l.config.coreConfig.lvl.SetLevel(lvl)
}

func (l *Logger) SetOutput(writer io.Writer) {
	ws := zapcore.AddSync(writer)
	log := zap.New(
		zapcore.NewCore(l.config.coreConfig.enc, ws, l.config.coreConfig.lvl),
		l.config.zapOpts...,
	)
	l.config.coreConfig.ws = ws
	l.SugaredLogger = log.Sugar()
}

func (l *Logger) CtxKVLog(ctx context.Context, level klog.Level, format string, kvs ...interface{}) {
	if len(kvs) == 0 || len(kvs)%2 != 0 {
		l.Warn(fmt.Sprint("Keyvalues must appear in pairs:", kvs))
		return
	}

	span := trace.SpanFromContext(ctx)
	if span.SpanContext().TraceID().IsValid() {
		kvs = append(kvs, traceIDKey, span.SpanContext().TraceID())
	}
	if span.SpanContext().SpanID().IsValid() {
		kvs = append(kvs, spanIDKey, span.SpanContext().SpanID())
	}
	if span.SpanContext().TraceFlags().IsSampled() {
		kvs = append(kvs, traceFlagsKey, span.SpanContext().TraceFlags())
	}

	var zlevel zapcore.Level
	zl := l.With()
	switch level {
	case klog.LevelDebug, klog.LevelTrace:
		zlevel = zap.DebugLevel
		zl.Debugw(format, kvs...)
	case klog.LevelInfo:
		zlevel = zap.InfoLevel
		zl.Infow(format, kvs...)
	case klog.LevelNotice, klog.LevelWarn:
		zlevel = zap.WarnLevel
		zl.Warnw(format, kvs...)
	case klog.LevelError:
		zlevel = zap.ErrorLevel
		zl.Errorw(format, kvs...)
	case klog.LevelFatal:
		zlevel = zap.FatalLevel
		zl.Fatalw(format, kvs...)
	default:
		zlevel = zap.WarnLevel
		zl.Warnw(format, kvs...)
	}

	if !span.IsRecording() {
		return
	}

	msg := getMessage(format, kvs)
	attrs := []attribute.KeyValue{
		logMessageKey.String(msg),
		logSeverityTextKey.String(OtelSeverityText(zlevel)),
	}

	// notice: AddEvent,SetStatus,RecordError all have check span.IsRecording
	span.AddEvent(logEventKey, trace.WithAttributes(attrs...))

	// set span status
	if zlevel <= l.config.traceConfig.errorSpanLevel {
		span.SetStatus(codes.Error, msg)
		span.RecordError(errors.New(msg), trace.WithStackTrace(l.config.traceConfig.recordStackTraceInSpan))
	}
}
