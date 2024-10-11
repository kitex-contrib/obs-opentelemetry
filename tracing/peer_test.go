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

/*func Test_extractPeerServiceAttributesFromMetaInfo(t *testing.T) {
	type args struct {
		md map[string]string
	}
	tests := []struct {
		name                 string
		args                 args
		want                 []attribute.KeyValue
		wantCanonicalService string
	}{
		{
			name: "peer service",
			args: args{
				md: map[string]string{
					string(semconv.ServiceNameKey): "foo",
				},
			},
			want: []attribute.KeyValue{
				semconv.PeerServiceKey.String("foo"),
			},
		},
		{
			name: "full peer",
			args: args{
				md: map[string]string{
					string(semconv.ServiceNameKey):           "foo",
					string(semconv.ServiceNamespaceKey):      "test-ns",
					string(semconv.DeploymentEnvironmentKey): "test-env",
				},
			},
			want: []attribute.KeyValue{
				semconv.PeerServiceKey.String("foo"),
				PeerServiceNamespaceKey.String("test-ns"),
				PeerDeploymentEnvironmentKey.String("test-env"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractPeerServiceAttributesFromMetaInfo(tt.args.md)
			assert.ElementsMatch(t, got, tt.want)
		})
	}
}

func Test_injectPeerServiceToMetaInfo(t *testing.T) {
	type args struct {
		ctx   context.Context
		attrs []attribute.KeyValue
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "peer service",
			args: args{
				ctx: context.Background(),
				attrs: []attribute.KeyValue{
					semconv.ServiceNameKey.String("foo"),
				},
			},
			want: map[string]string{
				"service.name": "foo",
			},
		},
		{
			name: "full peer",
			args: args{
				ctx: context.Background(),
				attrs: []attribute.KeyValue{
					semconv.ServiceNameKey.String("foo"),
					semconv.ServiceNamespaceKey.String("test-ns"),
					semconv.DeploymentEnvironmentKey.String("test-env"),
				},
			},
			want: map[string]string{
				string(semconv.ServiceNameKey):           "foo",
				string(semconv.ServiceNamespaceKey):      "test-ns",
				string(semconv.DeploymentEnvironmentKey): "test-env",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := injectPeerServiceToMetaInfo(tt.args.ctx, tt.args.attrs)
			assert.Equal(t, tt.want, got)
		})
	}
}
*/
