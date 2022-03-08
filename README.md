# opentelemetry

[Opentelemetry](https://opentelemetry.io/) for [Kitex](https://github.com/cloudwego/kitex)

## Server usage
```go
import (
    ...
    "github.com/kitex-contrib/obs-opentelemetry/logging/kitexlogrus"
    "github.com/kitex-contrib/obs-opentelemetry/provider"
    "github.com/kitex-contrib/obs-opentelemetry/tracing"
    "go.opentelemetry.io/otel/attribute"
    semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)


func main()  {
    p := provider.NewOpenTelemetryProvider(
        provider.WithServiceName("echo"),
        provider.WithExportEndpoint("localhost:4317"),
    )
    defer p.Shutdown(context.Background())

    svr := echo.NewServer(
        new(EchoImpl),
        server.WithSuite(tracing.NewServerSuite()),
    )
    if err := svr.Run(); err != nil {
        klog.Fatalf("server stopped with error:", err)
    } 	
}

```

## Client usage
```go
import (
    ...
    "github.com/kitex-contrib/obs-opentelemetry/logging/kitexlogrus"
    "github.com/kitex-contrib/obs-opentelemetry/provider"
    "github.com/kitex-contrib/obs-opentelemetry/tracing"
)

func main(){

    p := provider.NewOpenTelemetryProvider(
        provider.WithServiceName("echo-client"),
        provider.WithExportEndpoint("localhost:4317"),
    )
    defer p.Shutdown(context.Background())
    
    c, err := echo.NewClient(
        "echo",
        client.WithSuite(tracing.NewClientSuite()),
    )
    if err != nil {
        klog.Fatal(err)
    }
	
}

```

## Example

[Executable Example](https://github.com/cloudwego/kitex-examples/tree/main/opentelemetry)

## Supported Metrics

### RPC Metrics

### Kitex Server

Below is a table of RPC server metric instruments.

| Name | Instrument | Unit | Unit ([UCUM](README.md#instrument-units)) | Description | Status | Streaming |
|------|------------|------|-------------------------------------------|-------------|--------|-----------|
| `rpc.server.duration` | Histogram  | milliseconds | `ms` | measures duration of inbound RPC | Recommended | N/A.  While streaming RPCs may record this metric as start-of-batch to end-of-batch, it's hard to interpret in practice. |

### Kitex Client

Below is a table of RPC client metric instruments.  These apply to traditional
RPC usage, not streaming RPCs.

| Name | Instrument | Unit | Unit ([UCUM](README.md#instrument-units)) | Description | Status | Streaming |
|------|------------|------|-------------------------------------------|-------------|--------|-----------|
| `rpc.client.duration` | Histogram | milliseconds | `ms` | measures duration of outbound RPC | Recommended | N/A.  While streaming RPCs may record this metric as start-of-batch to end-of-batch, it's hard to interpret in practice. |

### Runtime Metrics

