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

package logrus_test

import (
	"context"
	"testing"

	"github.com/cloudwego/kitex/pkg/klog"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	"github.com/sirupsen/logrus"
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

	logger := kitexlogrus.NewLogger(
		kitexlogrus.WithTraceHookErrorSpanLevel(logrus.WarnLevel),
		kitexlogrus.WithTraceHookLevels(logrus.AllLevels),
		kitexlogrus.WithRecordStackTraceInSpan(true),
	)

	logger.Info("log from origin logrus")

	klog.SetLogger(logger)

	klog.SetLevel(klog.LevelDebug)

	tracer := otel.Tracer("test otel std logger")

	ctx, span := tracer.Start(ctx, "root")

	klog.CtxInfof(ctx, "hello %s", "world")

	span.End()

	ctx, child := tracer.Start(ctx, "child")

	klog.CtxWarnf(ctx, "foo %s", "bar")

	child.End()

	ctx, errSpan := tracer.Start(ctx, "error")

	klog.CtxErrorf(ctx, "error %s", "this is a error")

	klog.Info("no trace context")

	errSpan.End()
}
