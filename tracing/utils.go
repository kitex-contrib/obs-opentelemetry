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
	"time"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/stats"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

// recordErrorSpan log error to span
func recordErrorSpan(span trace.Span, err error, withStackTrace bool, attributes ...attribute.KeyValue) {
	if span == nil || err == nil {
		return
	}

	span.SetStatus(codes.Error, err.Error())
	span.RecordError(
		err,
		trace.WithAttributes(attributes...),
		trace.WithStackTrace(withStackTrace),
	)
}

func handleErr(err error) {
	if err != nil {
		otel.Handle(err)
	}
}

func getStartTimeOrNow(ri rpcinfo.RPCInfo) time.Time {
	if event := ri.Stats().GetEvent(stats.RPCStart); event != nil {
		return event.Time()
	}
	return time.Now()
}

func getEndTimeOrNow(ri rpcinfo.RPCInfo) time.Time {
	if event := ri.Stats().GetEvent(stats.RPCFinish); event != nil {
		return event.Time()
	}
	return time.Now()
}

func getServiceFromResourceAttributes(attrs []attribute.KeyValue) (serviceName string, serviceNamespace string, deploymentEnv string) {
	for _, attr := range attrs {
		switch attr.Key {
		case semconv.ServiceNameKey:
			serviceName = attr.Value.AsString()
		case semconv.ServiceNamespaceKey:
			serviceNamespace = attr.Value.AsString()
		case semconv.DeploymentEnvironmentKey:
			deploymentEnv = attr.Value.AsString()
		}
	}
	return
}
