module github.com/kitex-contrib/obs-opentelemetry/logging/logrus

go 1.19

require (
	github.com/cloudwego/kitex v0.5.2
	github.com/sirupsen/logrus v1.9.0
	go.opentelemetry.io/otel v1.15.1
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.15.1
	go.opentelemetry.io/otel/sdk v1.15.1
	go.opentelemetry.io/otel/trace v1.15.1
)

require (
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	golang.org/x/sys v0.8.0 // indirect
)
