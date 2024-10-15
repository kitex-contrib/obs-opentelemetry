module github.com/kitex-contrib/obs-opentelemetry/logging/logrus

go 1.21

toolchain go1.21.12

require (
	github.com/cloudwego-contrib/cwgo-pkg/telemetry/instrumentation/otellogrus v0.0.0-20241014044734-80a98dbe0b6a
	github.com/cloudwego/hertz v0.9.3
	github.com/cloudwego/kitex v0.9.1
	github.com/sirupsen/logrus v1.9.2
	go.opentelemetry.io/otel v1.28.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.28.0
	go.opentelemetry.io/otel/sdk v1.28.0
	go.opentelemetry.io/otel/trace v1.28.0
)

require (
	github.com/cloudwego-contrib/cwgo-pkg/log/logging/logrus v0.0.0-20241014044734-80a98dbe0b6a // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	go.opentelemetry.io/otel/metric v1.28.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
)
