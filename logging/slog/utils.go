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
	"fmt"
	"log/slog"
	"strings"

	"github.com/cloudwego/kitex/pkg/klog"
)

// get format msg
func getMessage(template string, fmtArgs []interface{}) string {
	if len(fmtArgs) == 0 {
		return template
	}

	if template != "" {
		return fmt.Sprintf(template, fmtArgs...)
	}

	if len(fmtArgs) == 1 {
		if str, ok := fmtArgs[0].(string); ok {
			return str
		}
	}
	return fmt.Sprint(fmtArgs...)
}

// OtelSeverityText convert slog level to otel severityText
// ref to https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/logs/data-model.md#severity-fields
func OtelSeverityText(lv slog.Level) string {
	s := lv.String()
	if s == "warning" {
		s = "warn"
	}
	return strings.ToUpper(s)
}

// Adapt klog level to slog level
func tranSLevel(level klog.Level) (lvl slog.Level) {
	switch level {
	case klog.LevelTrace:
		lvl = LevelTrace
	case klog.LevelDebug:
		lvl = slog.LevelDebug
	case klog.LevelInfo:
		lvl = slog.LevelInfo
	case klog.LevelWarn:
		lvl = slog.LevelWarn
	case klog.LevelNotice:
		lvl = LevelNotice
	case klog.LevelError:
		lvl = slog.LevelError
	case klog.LevelFatal:
		lvl = LevelFatal
	default:
		lvl = slog.LevelWarn
	}
	return
}
