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
	"log/slog"

	"github.com/cloudwego-contrib/cwgo-pkg/telemetry/instrumentation/otelslog"
	"github.com/cloudwego/kitex/pkg/klog"
)

const (
	LevelTrace  = slog.Level(-8)
	LevelNotice = slog.Level(2)
	LevelFatal  = slog.Level(12)
)

var _ klog.FullLogger = (*Logger)(nil)

type Logger struct {
	*otelslog.KLogger
}

func NewLogger(opts ...Option) *Logger {
	return &Logger{
		otelslog.NewKLogger(opts...),
	}
}
