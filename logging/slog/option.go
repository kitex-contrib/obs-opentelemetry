// Copyright 2023 CloudWeGo Authors.
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

package slog

import (
	"github.com/cloudwego-contrib/cwgo-pkg/telemetry/instrumentation/otelslog"
	"io"
	"log/slog"
)

type Option = otelslog.Option

// WithHandlerOptions slog handler-options
func WithHandlerOptions(opt *slog.HandlerOptions) Option {
	return otelslog.WithHandlerOptions(opt)
}

// WithOutput slog writer
func WithOutput(iow io.Writer) Option {
	return otelslog.WithOutput(iow)
}

// WithLevel slog level
func WithLevel(lvl *slog.LevelVar) Option {
	return otelslog.WithLevel(lvl)
}

// WithTraceErrorSpanLevel trace error span level option
func WithTraceErrorSpanLevel(level slog.Level) Option {
	return otelslog.WithTraceErrorSpanLevel(level)
}

// WithRecordStackTraceInSpan record stack track option
func WithRecordStackTraceInSpan(recordStackTraceInSpan bool) Option {
	return otelslog.WithRecordStackTraceInSpan(recordStackTraceInSpan)
}
