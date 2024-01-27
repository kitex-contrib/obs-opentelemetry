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
	"errors"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type ExtraKey string

type Option interface {
	apply(cfg *config)
}

type option func(cfg *config)

func (fn option) apply(cfg *config) {
	fn(cfg)
}

type traceConfig struct {
	recordStackTraceInSpan bool
	errorSpanLevel         zerolog.Level
}

type config struct {
	logger      *zerolog.Logger
	traceConfig *traceConfig
}

// defaultConfig default config
func defaultConfig() *config {
	return &config{
		traceConfig: &traceConfig{
			recordStackTraceInSpan: true,
			errorSpanLevel:         zerolog.ErrorLevel,
		},
	}
}

// WithLogger configures logger
func WithLogger(logger *zerolog.Logger) Option {
	return option(func(cfg *config) {
		cfg.logger = logger
	})
}

// WithTraceErrorSpanLevel trace error span level option
func WithTraceErrorSpanLevel(level zerolog.Level) Option {
	return option(func(cfg *config) {
		cfg.traceConfig.errorSpanLevel = level
	})
}

// WithRecordStackTraceInSpan record stack track option
func WithRecordStackTraceInSpan(recordStackTraceInSpan bool) Option {
	return option(func(cfg *config) {
		cfg.traceConfig.recordStackTraceInSpan = recordStackTraceInSpan
	})
}

func (cfg config) defaultZerologHookFn() zerolog.HookFunc {
	return func(e *zerolog.Event, level zerolog.Level, message string) {
		ctx := e.GetCtx()
		span := trace.SpanFromContext(ctx)
		spanCtx := span.SpanContext()

		if !spanCtx.IsValid() {
			return
		}

		e.Any(spanIDKey, spanCtx.SpanID())
		e.Any(traceIDKey, spanCtx.TraceID())
		e.Any(traceFlagsKey, spanCtx.TraceFlags())

		if !span.IsRecording() {
			return
		}

		// set span status
		if level >= cfg.traceConfig.errorSpanLevel {
			span.SetStatus(codes.Error, "")
			span.RecordError(errors.New(message),
				trace.WithStackTrace(cfg.traceConfig.recordStackTraceInSpan))
		}
	}
}
