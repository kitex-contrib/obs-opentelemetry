module github.com/kitex-contrib/obs-opentelemetry

go 1.16

require (
	github.com/bytedance/gopkg v0.0.0-20210910103821-e4efae9c17c3
	github.com/cloudwego/kitex v0.2.0
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/contrib/instrumentation/runtime v0.29.0
	go.opentelemetry.io/contrib/propagators/b3 v1.4.0
	go.opentelemetry.io/contrib/propagators/jaeger v1.4.0
	go.opentelemetry.io/contrib/propagators/opencensus v0.29.0
	go.opentelemetry.io/contrib/propagators/ot v1.4.0
	go.opentelemetry.io/otel v1.4.1
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric v0.27.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.27.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.4.1
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.4.1
	go.opentelemetry.io/otel/metric v0.27.0
	go.opentelemetry.io/otel/sdk v1.4.1
	go.opentelemetry.io/otel/sdk/metric v0.27.0
	go.opentelemetry.io/otel/trace v1.4.1
)
