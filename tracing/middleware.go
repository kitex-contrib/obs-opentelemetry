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
	"errors"
	"time"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/tracing/internal"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var (
	errTraceCarrierNotFound = errors.New("tracer not found in context")
)

// ClientMiddleware inject span context into req meta
func ClientMiddleware(cfg *config) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req, resp interface{}) (err error) {
			span := oteltrace.SpanFromContext(ctx)
			if !span.IsRecording() {
				return next(ctx, req, resp)
			}

			// inject client service resource attributes (canonical service) to baggage
			readOnlySpan := span.(trace.ReadOnlySpan)
			bags, err := injectCanonicalServiceToBaggage(readOnlySpan.Resource().Attributes())
			if err != nil {
				return err
			}
			ctx = baggage.ContextWithBaggage(ctx, bags)

			// inject to meta
			md := metainfo.GetAllValues(ctx)
			if md == nil {
				md = make(map[string]string)
			}
			Inject(ctx, cfg, md)
			for k, v := range md {
				ctx = metainfo.WithValue(ctx, k, v)
			}

			if err = next(ctx, req, resp); err != nil {
				RecordErrorSpan(span, err, cfg.withStackTrace)
			}
			return err
		}
	}
}

// ServerMiddleware extract req meta into span context
func ServerMiddleware(cfg *config) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req, resp interface{}) (err error) {
			tc := internal.TraceCarrierFromContext(ctx)
			if tc == nil {
				return errTraceCarrierNotFound
			}

			// get tracer from carrier
			sTracer := tc.Tracer()

			var spanName string
			if tc.SpanNameFormatter() != nil {
				spanName = tc.SpanNameFormatter()(ctx)
			}

			ri := rpcinfo.GetRPCInfo(ctx)
			opts := []oteltrace.SpanStartOption{
				oteltrace.WithTimestamp(getStartTimeOrDefault(ri, time.Now())),
				oteltrace.WithSpanKind(oteltrace.SpanKindServer),
			}

			md := metainfo.GetAllValues(ctx)
			bags, spanCtx := Extract(ctx, cfg, md)
			ctx = baggage.ContextWithBaggage(ctx, bags)

			ctx, span := sTracer.Start(oteltrace.ContextWithRemoteSpanContext(ctx, spanCtx), spanName, opts...)

			// peer service resource attrs
			attrs := peerServiceAttributesFromBaggage(bags)
			span.SetAttributes(attrs...)

			// set span and attrs into tracer carrier for serverTracer finish
			tc.SetSpan(span)

			// reset service baggage
			bags = resetPeerServiceBaggageMember(bags)
			ctx = baggage.ContextWithBaggage(ctx, bags)

			if err = next(ctx, req, resp); err != nil {
				RecordErrorSpan(span, err, cfg.withStackTrace)
			}

			return err
		}
	}
}
