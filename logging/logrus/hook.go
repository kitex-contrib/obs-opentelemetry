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

package logrus

import (
	"github.com/cloudwego-contrib/cwgo-pkg/telemetry/instrumentation/otellogrus"

	"github.com/sirupsen/logrus"
)

// Ref to https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/logs/README.md#json-formats
const (
	traceIDKey    = "trace_id"
	spanIDKey     = "span_id"
	traceFlagsKey = "trace_flags"
)

var _ logrus.Hook = (*TraceHook)(nil)

type TraceHookConfig = otellogrus.TraceHookConfig

type TraceHook = otellogrus.TraceHook

func NewTraceHook(cfg *TraceHookConfig) *TraceHook {
	return otellogrus.NewTraceHook(cfg)
}
