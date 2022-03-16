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
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

const (
	// RequestProtocolKey protocol of the request.
	//
	// Type: string
	// Required: Always
	// Examples:
	// http: 'http'
	// rpc: 'grpc', 'java_rmi', 'wcf', 'kitex'
	// db: mysql, postgresql
	// mq: 'rabbitmq', 'activemq', 'AmazonSQS'
	RequestProtocolKey = attribute.Key("request.protocol")
)

const (
	// RPCSystemKitexRecvSize recv_size
	RPCSystemKitexRecvSize = attribute.Key("kitex.recv_size")
	// RPCSystemKitexSendSize send_size
	RPCSystemKitexSendSize = attribute.Key("kitex.send_size")
)

const (
	// PeerServiceNamespaceKey peer.service.namespace
	PeerServiceNamespaceKey = attribute.Key("peer.service.namespace")
	// PeerDeploymentEnvironmentKey peer.deployment.environment
	PeerDeploymentEnvironmentKey = attribute.Key("peer.deployment.environment")
)

const (
	SourceCanonicalServiceKey = attribute.Key("source_canonical_service")
	SourceOperationKey        = attribute.Key("source_operation")
)

const (
	StatusKey = attribute.Key("status.code")
)

var (
	// RPCSystemKitex Semantic convention for kitex as the remoting system.
	RPCSystemKitex = semconv.RPCSystemKey.String("kitex")
)
