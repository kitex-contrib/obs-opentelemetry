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

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/tracing/internal"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// ClientMiddleware inject span context into req meta
func ClientMiddleware(cfg *config) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req, resp interface{}) (err error) {
			span := oteltrace.SpanFromContext(ctx)
			if !span.IsRecording() {
				return next(ctx, req, resp)
			}

			readOnlySpan := span.(trace.ReadOnlySpan)

			// inject client service resource attributes (canonical service) to meta info
			md := injectPeerServiceToMetaInfo(ctx, readOnlySpan.Resource().Attributes())

			Inject(ctx, cfg, md)

			for k, v := range md {
				ctx = metainfo.WithValue(ctx, k, v)
			}

			return next(ctx, req, resp)
		}
	}
}

// ServerMiddleware extract req meta into span context
func ServerMiddleware(cfg *config) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req, resp interface{}) (err error) {
			tc := internal.TraceCarrierFromContext(ctx)
			if tc == nil {
				klog.CtxWarnf(ctx, "TraceCarrier not found in context")
				return next(ctx, req, resp)
			}

			// get tracer from carrier
			sTracer := tc.Tracer()

			ri := rpcinfo.GetRPCInfo(ctx)
			opts := []oteltrace.SpanStartOption{
				oteltrace.WithTimestamp(getStartTimeOrNow(ri)),
				oteltrace.WithSpanKind(oteltrace.SpanKindServer),
			}

			md := metainfo.GetAllValues(ctx)
			peerServiceAttributes := extractPeerServiceAttributesFromMetaInfo(md)

			bags, spanCtx := Extract(ctx, cfg, md)
			ctx = baggage.ContextWithBaggage(ctx, bags)

			ctx, span := sTracer.Start(oteltrace.ContextWithRemoteSpanContext(ctx, spanCtx), spanNaming(ri), opts...)

			// peer service attributes
			span.SetAttributes(peerServiceAttributes...)

			// set span and attrs into tracer carrier for serverTracer finish
			tc.SetSpan(span)

			return next(ctx, req, resp)
		}
	}
}
