package zap

import (
	"fmt"

	"go.uber.org/zap/zapcore"
)

// getMessage format with Sprint, Sprintf, or neither.
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

// OtelSeverityText convert zapcore level to otel severityText
// ref to https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/logs/data-model.md#severity-fields
func OtelSeverityText(lv zapcore.Level) string {
	s := lv.CapitalString()
	if s == "DPANIC" || s == "PANIC" {
		s = "FATAL"
	}
	return s
}
