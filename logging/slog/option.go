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

package slog

import (
	"io"
	"log/slog"
	"os"
)

type Option interface {
	apply(cfg *config)
}

type option func(cfg *config)

func (fn option) apply(cfg *config) {
	fn(cfg)
}

type coreConfig struct {
	opt    slog.HandlerOptions
	writer io.Writer
	level  *slog.LevelVar
}

type config struct {
	coreConfig  coreConfig
	traceConfig *traceConfig
}

func defaultConfig() *config {
	coreConfig := defaultCoreConfig()
	return &config{
		coreConfig: *coreConfig,
		traceConfig: &traceConfig{
			recordStackTraceInSpan: true,
			errorSpanLevel:         slog.LevelError,
		},
	}
}

func defaultCoreConfig() *coreConfig {
	level := new(slog.LevelVar)
	level.Set(slog.LevelInfo)
	return &coreConfig{
		opt: slog.HandlerOptions{
			Level:       level,
			ReplaceAttr: replaceAttr,
		},
		writer: os.Stdout,
		level:  level,
	}
}

// WithHandlerOptions slog handler-options
func WithHandlerOptions(opt *slog.HandlerOptions) Option {
	return option(func(cfg *config) {
		cfg.coreConfig.opt = *opt
	})
}

// WithWriter slog writer
func WithWriter(iow io.Writer) Option {
	return option(func(cfg *config) {
		cfg.coreConfig.writer = iow
	})
}

// WithLevel slog level
func WithLevel(lvl *slog.LevelVar) Option {
	return option(func(cfg *config) {
		cfg.coreConfig.opt.Level = lvl
	})
}

// WithTraceErrorSpanLevel trace error span level option
func WithTraceErrorSpanLevel(level slog.Level) Option {
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

func replaceAttr(_ []string, a slog.Attr) slog.Attr {
	if a.Key == slog.LevelKey {
		level := a.Value.Any().(slog.Level)
		switch level {
		case LevelTrace:
			a.Value = slog.StringValue("Trace")
		case slog.LevelDebug:
			a.Value = slog.StringValue("Debug")
		case slog.LevelInfo:
			a.Value = slog.StringValue("Info")
		case LevelNotice:
			a.Value = slog.StringValue("Notice")
		case slog.LevelWarn:
			a.Value = slog.StringValue("Warn")
		case slog.LevelError:
			a.Value = slog.StringValue("Error")
		case LevelFatal:
			a.Value = slog.StringValue("Fatal")
		default:
			a.Value = slog.StringValue("Warn")
		}
	}
	return a
}
