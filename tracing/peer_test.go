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
	"reflect"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func Test_injectCanonicalServiceToBaggage(t *testing.T) {
	canonicalServiceBaggage, _ := baggage.Parse("service.name=foo,service.namespace=test-ns,deployment.environment=dev")
	type args struct {
		attrs []attribute.KeyValue
	}
	tests := []struct {
		name    string
		args    args
		want    baggage.Baggage
		wantErr bool
	}{
		{
			name: "inject valid canonical service",
			args: args{
				attrs: []attribute.KeyValue{
					semconv.ServiceNameKey.String("foo"),
					semconv.ServiceNamespaceKey.String("test-ns"),
					semconv.DeploymentEnvironmentKey.String("dev"),
				},
			},
			want:    canonicalServiceBaggage,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := injectCanonicalServiceToBaggage(tt.args.attrs)
			if (err != nil) != tt.wantErr {
				t.Errorf("injectCanonicalServiceToBaggage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("injectCanonicalServiceToBaggage() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_peerServiceAttributesFromBaggage(t *testing.T) {
	canonicalServiceBaggage, _ := baggage.Parse("service.name=foo,service.namespace=test-ns,deployment.environment=dev")
	type args struct {
		bags baggage.Baggage
	}
	tests := []struct {
		name string
		args args
		want []attribute.KeyValue
	}{
		{
			name: "peer service attrs from baggage",
			args: args{
				bags: canonicalServiceBaggage,
			},
			want: []attribute.KeyValue{
				semconv.PeerServiceKey.String("foo"),
				PeerServiceNamespaceKey.String("test-ns"),
				PeerDeploymentEnvironmentKey.String("dev"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := peerServiceAttributesFromBaggage(tt.args.bags); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("peerServiceAttributesFromBaggage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_resetPeerServiceBaggageMember(t *testing.T) {
	canonicalServiceBaggage, _ := baggage.Parse("service.name=foo,service.namespace=test-ns,deployment.environment=dev,foo=bar")
	resetBaggage, _ := baggage.Parse("foo=bar")
	type args struct {
		bags baggage.Baggage
	}
	tests := []struct {
		name string
		args args
		want baggage.Baggage
	}{
		{
			name: "reset successful",
			args: args{
				bags: canonicalServiceBaggage,
			},
			want: resetBaggage,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resetPeerServiceBaggageMember(tt.args.bags); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resetPeerServiceBaggageMember() = %v, want %v", got, tt.want)
			}
		})
	}
}
