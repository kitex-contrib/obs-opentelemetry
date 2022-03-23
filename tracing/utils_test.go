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
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func Test_getServiceFromResourceAttributes(t *testing.T) {
	type args struct {
		attrs []attribute.KeyValue
	}
	tests := []struct {
		name                 string
		args                 args
		wantServiceName      string
		wantServiceNamespace string
		wantDeploymentEnv    string
	}{
		{
			name: "valid",
			args: args{
				attrs: []attribute.KeyValue{
					semconv.ServiceNameKey.String("foo"),
				},
			},
			wantServiceName:      "foo",
			wantServiceNamespace: "",
			wantDeploymentEnv:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotServiceName, gotServiceNamespace, gotDeploymentEnv := getServiceFromResourceAttributes(tt.args.attrs)
			if gotServiceName != tt.wantServiceName {
				t.Errorf("getServiceFromResourceAttributes() gotServiceName = %v, want %v", gotServiceName, tt.wantServiceName)
			}
			if gotServiceNamespace != tt.wantServiceNamespace {
				t.Errorf("getServiceFromResourceAttributes() gotServiceNamespace = %v, want %v", gotServiceNamespace, tt.wantServiceNamespace)
			}
			if gotDeploymentEnv != tt.wantDeploymentEnv {
				t.Errorf("getServiceFromResourceAttributes() gotDeploymentEnv = %v, want %v", gotDeploymentEnv, tt.wantDeploymentEnv)
			}
		})
	}
}

func Test_recordErrorSpan(t *testing.T) {
	sr := tracetest.NewSpanRecorder()
	tp := trace.NewTracerProvider(trace.WithSpanProcessor(sr))
	defer tp.Shutdown(context.Background())

	type args struct {
		err            error
		withStackTrace bool
		attributes     []attribute.KeyValue
	}
	tests := []struct {
		name                   string
		args                   args
		wantEventsLen          int
		wantEventAttributesLen int
	}{
		{
			name: "empty attributes",
			args: args{
				err:            errors.New("mock error"),
				withStackTrace: true,
				attributes:     nil,
			},
			wantEventsLen:          1,
			wantEventAttributesLen: 3,
		},
		{
			name: "with attributes",
			args: args{
				err:            errors.New("mock error"),
				withStackTrace: true,
				attributes: []attribute.KeyValue{
					RPCSystemKitex,
				},
			},
			wantEventsLen:          1,
			wantEventAttributesLen: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, testSpan := tp.Tracer("test").Start(context.Background(), "test-span")
			defer testSpan.End()

			recordErrorSpanWithStack(testSpan, tt.args.err, "mock panic", "mock stack")

			readOnlySpan := testSpan.(trace.ReadOnlySpan)

			assert.Equal(t, trace.Status{Code: codes.Error, Description: "mock error"}, readOnlySpan.Status())
			assert.Equal(t, tt.wantEventsLen, len(readOnlySpan.Events()))
			for _, event := range readOnlySpan.Events() {
				assert.Equal(t, "exception", event.Name)
				assert.Equal(t, tt.wantEventAttributesLen, len(event.Attributes))
			}
		})
	}
}
