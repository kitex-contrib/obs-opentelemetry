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
	"errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

// recordErrorSpanWithStack record error with stack
func recordErrorSpanWithStack(span trace.Span, err error, stackMessage, stackTrace string, attributes ...attribute.KeyValue) {
	if span == nil {
		return
	}

	// compatible with the case where error is empty
	if err == nil {
		err = errors.New(stackMessage)
	}

	// stack trace
	attributes = append(attributes,
		semconv.ExceptionStacktraceKey.String(stackTrace),
	)

	span.SetStatus(codes.Error, err.Error())
	span.RecordError(
		err,
		trace.WithAttributes(attributes...),
	)
}

func getServiceFromResourceAttributes(attrs []attribute.KeyValue) (serviceName, serviceNamespace, deploymentEnv string) {
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
