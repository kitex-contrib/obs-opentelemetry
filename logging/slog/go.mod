module github.com/kitex-contrib/obs-opentelemetry/logging/slog

go 1.21

require (
	github.com/cloudwego-contrib/cwgo-pkg/telemetry/instrumentation/otelslog v0.0.0
	github.com/cloudwego/kitex v0.9.1
	github.com/stretchr/testify v1.9.0
	go.opentelemetry.io/otel v1.28.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.28.0
	go.opentelemetry.io/otel/sdk v1.28.0
	go.opentelemetry.io/otel/trace v1.28.0
)

require (
	github.com/cloudwego-contrib/cwgo-pkg/log/logging/slog v0.0.0-20241014044734-80a98dbe0b6a // indirect
	github.com/cloudwego/hertz v0.9.3 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/otel/metric v1.28.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/cloudwego-contrib/cwgo-pkg/telemetry/instrumentation/otelslog => github.com/smx-Morgan/cwgo-pkg/telemetry/instrumentation/otelslog v0.0.0-20241019002536-84cf43046703
