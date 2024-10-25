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

package zerolog

import (
	"github.com/cloudwego-contrib/cwgo-pkg/telemetry/instrumentation/otelzerolog"

	"github.com/rs/zerolog"
)

type ExtraKey string

type Option = otelzerolog.Option

// WithLogger configures logger
func WithLogger(logger *zerolog.Logger) Option {
	return otelzerolog.WithZeroLogger(logger)
}

// WithTraceErrorSpanLevel trace error span level option
func WithTraceErrorSpanLevel(level zerolog.Level) Option {
	return otelzerolog.WithTraceErrorSpanLevel(level)
}

// WithRecordStackTraceInSpan record stack track option
func WithRecordStackTraceInSpan(recordStackTraceInSpan bool) Option {
	return otelzerolog.WithRecordStackTraceInSpan(recordStackTraceInSpan)
}
