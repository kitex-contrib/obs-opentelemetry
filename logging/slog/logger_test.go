// Copyright 2023 CloudWeGo Authors.
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
	"bytes"
	"context"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func stdoutProvider(ctx context.Context) func() {
	provider := sdktrace.NewTracerProvider()
	otel.SetTracerProvider(provider)

	exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		panic(err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(exp)
	provider.RegisterSpanProcessor(bsp)

	return func() {
		if err := provider.Shutdown(ctx); err != nil {
			panic(err)
		}
	}
}

func TestLogger(t *testing.T) {
	ctx := context.Background()

	buf := new(bytes.Buffer)

	shutdown := stdoutProvider(ctx)
	defer shutdown()

	logger := NewLogger(
		WithTraceErrorSpanLevel(slog.LevelWarn),
		WithRecordStackTraceInSpan(true),
	)

	klog.SetLogger(logger)
	klog.SetOutput(buf)
	klog.SetLevel(klog.LevelDebug)

	logger.Info("log from origin zap")
	assert.True(t, strings.Contains(buf.String(), "log from origin zap"))
	buf.Reset()

	tracer := otel.Tracer("test otel std logger")

	ctx, span := tracer.Start(ctx, "root")

	klog.CtxInfof(ctx, "hello %s", "you")
	assert.True(t, strings.Contains(buf.String(), "trace_id"))
	assert.True(t, strings.Contains(buf.String(), "span_id"))
	assert.True(t, strings.Contains(buf.String(), "trace_flags"))
	buf.Reset()

	span.End()

	ctx, child := tracer.Start(ctx, "child")

	klog.CtxWarnf(ctx, "foo %s", "bar")

	klog.CtxTracef(ctx, "trace %s", "this is a trace log")
	klog.CtxDebugf(ctx, "debug %s", "this is a debug log")
	klog.CtxInfof(ctx, "info %s", "this is a info log")
	klog.CtxNoticef(ctx, "notice %s", "this is a notice log")
	klog.CtxWarnf(ctx, "warn %s", "this is a warn log")
	klog.CtxErrorf(ctx, "error %s", "this is a error log")

	child.End()

	_, errSpan := tracer.Start(ctx, "error")

	klog.Info("no trace context")

	errSpan.End()
}

func TestLogLevel(t *testing.T) {
	buf := new(bytes.Buffer)

	logger := NewLogger(
		WithTraceErrorSpanLevel(slog.LevelWarn),
		WithRecordStackTraceInSpan(true),
	)

	// output to buffer
	logger.SetOutput(buf)

	logger.Debug("this is a debug log")
	assert.False(t, strings.Contains(buf.String(), "this is a debug log"))

	logger.SetLevel(klog.LevelDebug)

	logger.Debugf("this is a debug log %s", "msg")
	assert.True(t, strings.Contains(buf.String(), "this is a debug log"))
}

func TestLogOption(t *testing.T) {
	buf := new(bytes.Buffer)

	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelDebug)
	logger := NewLogger(
		WithLevel(lvl),
		WithWriter(buf),
		WithHandlerOptions(&slog.HandlerOptions{
			AddSource: true,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.MessageKey {
					msg := a.Value.Any().(string)
					msg = strings.ReplaceAll(msg, "log", "new log")
					a.Value = slog.StringValue(msg)
				}
				return a
			},
		}),
	)

	logger.Debug("this is a debug log")
	assert.True(t, strings.Contains(buf.String(), "this is a debug new log"))

	dir, _ := os.Getwd()
	assert.True(t, strings.Contains(buf.String(), dir))
}
