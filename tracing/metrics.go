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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// RPC Server metrics
// ref to https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/metrics/semantic_conventions/rpc.md#rpc-server
const (
	ServerDuration        = "rpc.server.duration"          // measures duration of inbound RPC
	ServerRequestSize     = "rpc.server.request.size"      // measures size of RPC request messages (uncompressed)
	ServerResponseSize    = "rpc.server.response.size"     // measures size of RPC response messages (uncompressed)
	ServerRequestsPerRPC  = "rpc.server.requests_per_rpc"  // measures the number of messages received per RPC. Should be 1 for all non-streaming RPCs
	ServerResponsesPerRPC = "rpc.server.responses_per_rpc" // measures the number of messages sent per RPC. Should be 1 for all non-streaming RPCs
)

// RPC Client metrics
// ref to https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/metrics/semantic_conventions/rpc.md#rpc-client
const (
	ClientDuration        = "rpc.client.duration"          // measures duration of outbound RPC
	ClientRequestSize     = "rpc.client.request.size"      // measures size of RPC request messages (uncompressed)
	ClientResponseSize    = "rpc.client.response.size"     // measures size of RPC response messages (uncompressed)
	ClientRequestsPerRPC  = "rpc.client.requests_per_rpc"  // measures the number of messages received per RPC. Should be 1 for all non-streaming RPCs
	ClientResponsesPerRPC = "rpc.client.responses_per_rpc" // measures the number of messages sent per RPC. Should be 1 for all non-streaming RPCs
)

var (
	// RPCMetricsAttributes rpc metrics attributes
	// ref to https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/metrics/semantic_conventions/rpc.md#attributes
	RPCMetricsAttributes = []attribute.Key{
		semconv.RPCServiceKey,
		semconv.RPCSystemKey,
		semconv.RPCMethodKey,
		semconv.NetPeerNameKey,
		semconv.NetTransportKey,
	}

	PeerMetricsAttributes = []attribute.Key{
		semconv.PeerServiceKey,
		PeerServiceNamespaceKey,
		PeerDeploymentEnvironmentKey,
		RequestProtocolKey,
		SourceOperationKey,
	}

	// MetricResourceAttributes resource attributes
	MetricResourceAttributes = []attribute.Key{
		semconv.ServiceNameKey,
		semconv.ServiceNamespaceKey,
		semconv.DeploymentEnvironmentKey,
		semconv.ServiceInstanceIDKey,
		semconv.ServiceVersionKey,
		semconv.TelemetrySDKLanguageKey,
		semconv.TelemetrySDKVersionKey,
		semconv.ProcessPIDKey,
		semconv.HostNameKey,
		semconv.HostIDKey,
	}
)

func extractMetricsAttributesFromSpan(span oteltrace.Span) []attribute.KeyValue {
	var attrs []attribute.KeyValue
	readOnlySpan, ok := span.(trace.ReadOnlySpan)
	if !ok {
		return attrs
	}

	// span attributes
	for _, attr := range readOnlySpan.Attributes() {
		if matchAttributeKey(attr.Key, RPCMetricsAttributes) {
			attrs = append(attrs, attr)
		}

		if matchAttributeKey(attr.Key, PeerMetricsAttributes) {
			attrs = append(attrs, attr)
		}
	}

	// span resource attributes
	for _, attr := range readOnlySpan.Resource().Attributes() {
		if matchAttributeKey(attr.Key, MetricResourceAttributes) {
			attrs = append(attrs, attr)
		}
	}

	// status code
	attrs = append(attrs, StatusKey.String(readOnlySpan.Status().Code.String()))

	return attrs
}

func matchAttributeKey(key attribute.Key, toMatchKeys []attribute.Key) bool {
	for _, attrKey := range toMatchKeys {
		if attrKey == key {
			return true
		}
	}
	return false
}
