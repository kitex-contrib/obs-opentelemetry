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
	"os"
	"testing"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/rs/zerolog"
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
	shutdown := stdoutProvider(ctx)
	defer shutdown()

	logger := NewLogger(
		WithTraceErrorSpanLevel(zerolog.WarnLevel),
		WithRecordStackTraceInSpan(true),
	)

	logger.Logger().Info().Msg("log from origin zerolog")

	klog.SetLogger(logger)

	klog.SetLevel(klog.LevelDebug)

	klog.SetOutput(os.Stderr)

	tracer := otel.Tracer("test otel std logger")

	ctx, span := tracer.Start(ctx, "root")

	klog.CtxInfof(ctx, "hello %s", "world")

	defer span.End()

	ctx, child := tracer.Start(ctx, "child")

	klog.CtxDebugf(ctx, "foo %s", "bar")
	klog.CtxTracef(ctx, "foo %s", "bar")
	klog.CtxInfof(ctx, "foo %s", "bar")
	klog.CtxNoticef(ctx, "foo %s", "bar")
	klog.CtxWarnf(ctx, "foo %s", "bar")
	klog.CtxErrorf(ctx, "foo %s", "bar")
	klog.Debugf("foo %s", "bar")
	klog.Tracef("foo %s", "bar")
	klog.Infof("foo %s", "bar")
	klog.Noticef("foo %s", "bar")
	klog.Warnf("foo %s", "bar")
	klog.Errorf("foo %s", "bar")
	klog.Debug("foo bar")
	klog.Trace("foo bar")
	klog.Info("foo bar")
	klog.Notice("foo bar")
	klog.Warn("foo bar")
	klog.Error("foo bar")

	child.End()

	ctx, errSpan := tracer.Start(ctx, "error")

	klog.CtxErrorf(ctx, "error %s", "this is a error")

	klog.Info("no trace context")

	errSpan.End()
}
