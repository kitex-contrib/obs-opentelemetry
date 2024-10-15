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

package zap

import (
	"github.com/cloudwego-contrib/cwgo-pkg/telemetry/instrumentation/otelzap"

	cwzap "github.com/cloudwego-contrib/cwgo-pkg/log/logging/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ExtraKey = cwzap.ExtraKey

type Option = otelzap.Option

// WithCoreEnc zapcore encoder
func WithCoreEnc(enc zapcore.Encoder) Option {
	return otelzap.WithCoreEnc(enc)
}

// WithCoreWs zapcore write syncer
func WithCoreWs(ws zapcore.WriteSyncer) Option {
	return otelzap.WithCoreWs(ws)
}

// WithCoreLevel zapcore log level
func WithCoreLevel(lvl zap.AtomicLevel) Option {
	return otelzap.WithCoreLevel(lvl)
}

// WithCustomFields record log with the key-value pair.
func WithCustomFields(kv ...interface{}) Option {
	return otelzap.WithCustomFields(kv...)
}

// WithZapOptions add origin zap option
func WithZapOptions(opts ...zap.Option) Option {
	return otelzap.WithZapOptions(opts...)
}

// WithTraceErrorSpanLevel trace error span level option
func WithTraceErrorSpanLevel(level zapcore.Level) Option {
	return otelzap.WithTraceErrorSpanLevel(level)
}

// WithRecordStackTraceInSpan record stack track option
func WithRecordStackTraceInSpan(recordStackTraceInSpan bool) Option {
	return otelzap.WithRecordStackTraceInSpan(recordStackTraceInSpan)
}

/*// WithExtraKeys allow you log extra values from context
func WithExtraKeys(keys []ExtraKey) Option {
	return otelzap.WithE
}*/

// WithExtraKeyAsStr convert extraKey to a string type when retrieving value from context
// Not recommended for use, only for compatibility with certain situations
//
// For more information, refer to the documentation at
// `https://pkg.go.dev/context#WithValue`
/*func WithExtraKeyAsStr() Option {
	return option(func(cfg *config) {
		cfg.extraKeyAsStr = true
	})
}*/
