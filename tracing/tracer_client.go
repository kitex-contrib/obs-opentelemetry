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

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/stats"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var _ stats.Tracer = (*clientTracer)(nil)

type clientTracer struct {
	config            *config
	histogramRecorder map[string]syncfloat64.Histogram
}

func newClientOption(opts ...Option) (client.Option, *config) {
	cfg := newConfig(opts)
	ct := &clientTracer{config: cfg}

	ct.createMeasures()

	return client.WithTracer(ct), cfg
}

func (c *clientTracer) createMeasures() {
	clientDurationMeasure, err := c.config.meter.SyncFloat64().Histogram(ClientDuration)
	handleErr(err)

	c.histogramRecorder = map[string]syncfloat64.Histogram{
		ClientDuration: clientDurationMeasure,
	}
}

func (c *clientTracer) Start(ctx context.Context) context.Context {
	ri := rpcinfo.GetRPCInfo(ctx)
	ctx, _ = c.config.tracer.Start(
		ctx,
		spanNaming(ri),
		oteltrace.WithTimestamp(getStartTimeOrNow(ri)),
		oteltrace.WithSpanKind(oteltrace.SpanKindClient),
	)

	return ctx
}

func (c *clientTracer) Finish(ctx context.Context) {
	span := oteltrace.SpanFromContext(ctx)
	if span == nil || !span.IsRecording() {
		return
	}

	ri := rpcinfo.GetRPCInfo(ctx)
	if ri.Stats().Level() == stats.LevelDisabled {
		return
	}

	st := ri.Stats()
	rpcStart := st.GetEvent(stats.RPCStart)
	rpcFinish := st.GetEvent(stats.RPCFinish)
	duration := rpcFinish.Time().Sub(rpcStart.Time())
	elapsedTime := float64(duration) / float64(time.Millisecond)

	attrs := []attribute.KeyValue{
		RPCSystemKitex,
		semconv.RPCMethodKey.String(ri.To().Method()),
		semconv.RPCServiceKey.String(ri.To().ServiceName()),
		RPCSystemKitexRecvSize.Int64(int64(st.RecvSize())),
		RPCSystemKitexSendSize.Int64(int64(st.SendSize())),
		RequestProtocolKey.String(ri.Config().TransportProtocol().String()),
	}

	// The source operation dimension maybe cause high cardinality issues
	if c.config.recordSourceOperation {
		attrs = append(attrs, SourceOperationKey.String(ri.From().Method()))
	}

	span.SetAttributes(attrs...)

	injectStatsEventsToSpan(span, st)

	if panicMsg, panicStack, rpcErr := parseRPCError(ri); rpcErr != nil || len(panicMsg) > 0 {
		recordErrorSpanWithStack(span, rpcErr, panicMsg, panicStack)
	}

	span.End(oteltrace.WithTimestamp(getEndTimeOrNow(ri)))

	metricsAttributes := extractMetricsAttributesFromSpan(span)
	c.histogramRecorder[ClientDuration].Record(ctx, elapsedTime, metricsAttributes...)
}
