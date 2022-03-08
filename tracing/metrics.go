package tracing

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
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
		SourceCanonicalServiceKey,
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
	readOnlySpan := span.(trace.ReadOnlySpan)

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
