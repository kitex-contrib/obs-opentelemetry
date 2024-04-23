// Copyright 2024 CloudWeGo Authors.
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

	"github.com/cloudwego/kitex/pkg/remote/trans/nphttp2/metadata"

	"go.opentelemetry.io/otel/propagation"
)

type metadataSupplier struct {
	metadata *metadata.MD
}

// assert that metadataSupplier implements the TextMapCarrier interface.
var _ propagation.TextMapCarrier = &metadataSupplier{}

func (s *metadataSupplier) Get(key string) string {
	values := s.metadata.Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (s *metadataSupplier) Set(key string, value string) {
	s.metadata.Set(key, value)
}

func (s *metadataSupplier) Keys() []string {
	out := make([]string, 0, len(*s.metadata))
	for key := range *s.metadata {
		out = append(out, key)
	}
	return out
}

func injectMetadata(ctx context.Context, cfg *config, md metadata.MD) context.Context {
	cfg.textMapPropagator.Inject(ctx, &metadataSupplier{
		metadata: &md,
	})
	return metadata.NewOutgoingContext(ctx, md)
}

func extractMetadata(ctx context.Context, cfg *config, md metadata.MD) context.Context {
	return cfg.textMapPropagator.Extract(ctx, &metadataSupplier{
		metadata: &md,
	})
}
