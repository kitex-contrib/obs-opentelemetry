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

package tracing

import (
	"context"
	"time"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/stats"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/tracing/internal"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var _ stats.Tracer = (*serverTracer)(nil)

type serverTracer struct {
	config            *config
	histogramRecorder map[string]syncfloat64.Histogram
}

func newServerOption(opts ...Option) (server.Option, *config) {
	cfg := newConfig(opts)
	st := &serverTracer{
		config: cfg,
	}

	st.createMeasures()

	return server.WithTracer(st), cfg
}

func (s *serverTracer) createMeasures() {
	serverDurationMeasure, err := s.config.meter.SyncFloat64().Histogram(ServerDuration)
	handleErr(err)

	s.histogramRecorder = map[string]syncfloat64.Histogram{
		ServerDuration: serverDurationMeasure,
	}
}

func (s *serverTracer) Start(ctx context.Context) context.Context {
	tc := &internal.TraceCarrier{}
	tc.SetTracer(s.config.tracer)

	return internal.WithTraceCarrier(ctx, tc)
}

func (s *serverTracer) Finish(ctx context.Context) {
	// trace carrier from context
	tc := internal.TraceCarrierFromContext(ctx)
	if tc == nil {
		return
	}

	// rpc info
	ri := rpcinfo.GetRPCInfo(ctx)
	if ri.Stats().Level() == stats.LevelDisabled {
		return
	}

	st := ri.Stats()
	rpcStart := st.GetEvent(stats.RPCStart)
	rpcFinish := st.GetEvent(stats.RPCFinish)
	duration := rpcFinish.Time().Sub(rpcStart.Time())
	elapsedTime := float64(duration) / float64(time.Millisecond)

	// span
	span := tc.Span()
	if span == nil || !span.IsRecording() {
		return
	}

	// span attributes
	attrs := []attribute.KeyValue{
		RPCSystemKitex,
		semconv.RPCMethodKey.String(ri.To().Method()),
		semconv.RPCServiceKey.String(ri.To().ServiceName()),
		RPCSystemKitexRecvSize.Int64(int64(st.RecvSize())),
		RPCSystemKitexSendSize.Int64(int64(st.SendSize())),
		RequestProtocolKey.String(ri.Config().TransportProtocol().String()),
	}

	// The source operation dimension maybe cause high cardinality issues
	if s.config.recordSourceOperation {
		attrs = append(attrs, SourceOperationKey.String(ri.From().Method()))
	}

	span.SetAttributes(attrs...)

	injectStatsEventsToSpan(span, st)

	if panicMsg, panicStack, rpcErr := parseRPCError(ri); rpcErr != nil || len(panicMsg) > 0 {
		recordErrorSpanWithStack(span, rpcErr, panicMsg, panicStack)
	}

	span.End(oteltrace.WithTimestamp(getEndTimeOrNow(ri)))

	metricsAttributes := extractMetricsAttributesFromSpan(span)
	s.histogramRecorder[ServerDuration].Record(ctx, elapsedTime, metricsAttributes...)
}
