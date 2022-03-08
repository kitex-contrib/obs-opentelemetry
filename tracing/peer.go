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
	"go.opentelemetry.io/otel/baggage"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func resetPeerServiceBaggageMember(bags baggage.Baggage) baggage.Baggage {
	return bags.DeleteMember(string(semconv.ServiceNameKey)).
		DeleteMember(string(semconv.ServiceNamespaceKey)).
		DeleteMember(string(semconv.DeploymentEnvironmentKey))
}

func peerServiceAttributesFromBaggage(bags baggage.Baggage) []attribute.KeyValue {
	var attrs []attribute.KeyValue

	serviceName := bags.Member(string(semconv.ServiceNameKey))
	if serviceName.Value() != "" {
		attrs = append(attrs, semconv.PeerServiceKey.String(serviceName.Value()))
	}

	serviceNamespace := bags.Member(string(semconv.ServiceNamespaceKey))
	if serviceNamespace.Value() != "" {
		attrs = append(attrs, PeerServiceNamespaceKey.String(serviceNamespace.Value()))
	}

	deploymentEnvironment := bags.Member(string(semconv.DeploymentEnvironmentKey))
	if deploymentEnvironment.Value() != "" {
		attrs = append(attrs, PeerDeploymentEnvironmentKey.String(deploymentEnvironment.Value()))
	}

	return attrs
}

func injectCanonicalServiceToBaggage(attrs []attribute.KeyValue) (baggage.Baggage, error) {
	var bags []baggage.Member

	serviceName, serviceNamespace, deploymentEnv := getServiceFromResourceAttributes(attrs)

	if serviceName != "" {
		serviceNameM, _ := baggage.NewMember(string(semconv.ServiceNameKey), serviceName)
		bags = append(bags, serviceNameM)
	}

	if serviceNamespace != "" {
		serviceNamespaceM, _ := baggage.NewMember(string(semconv.ServiceNamespaceKey), serviceNamespace)
		bags = append(bags, serviceNamespaceM)
	}

	if deploymentEnv != "" {
		deploymentEnvM, _ := baggage.NewMember(string(semconv.DeploymentEnvironmentKey), deploymentEnv)
		bags = append(bags, deploymentEnvM)
	}

	return baggage.New(bags...)
}
