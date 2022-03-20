# opentelemetry

[Opentelemetry](https://opentelemetry.io/) for [Kitex](https://github.com/cloudwego/kitex)

## Feature
#### Provider
- [x] Out-of-the-box  default opentelemetry provider

### Instrumentation

#### Tracing
- [x] Support server and client kitex rpc tracing
- [x] Support automatic transparent transmission of peer service through baggage

#### Metrics
- [x] Support kitex rpc metrics [RED]
- [x] Support peer service dimension in rpc metrics
- [x] Support go runtime metrics

#### Logging
- [x] Extend kitex logger based on logrus
- [x] Implement tracing auto associated logs

## Server usage
```go
import (
    ...
    "github.com/kitex-contrib/obs-opentelemetry/provider"
    "github.com/kitex-contrib/obs-opentelemetry/tracing"
)


func main()  {
    p := provider.NewOpenTelemetryProvider(
        provider.WithServiceName("echo"),
        provider.WithExportEndpoint("localhost:4317"),
        provider.WithInsecure(),
    )
    defer p.Shutdown(context.Background())

    svr := echo.NewServer(
        new(EchoImpl),
        server.WithSuite(tracing.NewServerSuite()),
        server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "echo"}),
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
    "github.com/kitex-contrib/obs-opentelemetry/provider"
    "github.com/kitex-contrib/obs-opentelemetry/tracing"
)

func main(){

    p := provider.NewOpenTelemetryProvider(
        provider.WithServiceName("echo-client"),
        provider.WithExportEndpoint("localhost:4317"),
        provider.WithInsecure(),
    )
    defer p.Shutdown(context.Background())
    
    c, err := echo.NewClient(
        "echo",
        client.WithSuite(tracing.NewClientSuite()),
        client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "echo-client"}),
    )
    if err != nil {
        klog.Fatal(err)
    }
	
}

```


## Kitex Multi Protocol
> Kitex supports multiple protocols, we can use different suites to initialize opentelemetry observability suites.
#### Thrift + TTHeader
```go
// server
server.WithSuite(tracing.NewServerSuite())
// client
client.WithSuite(tracing.NewClientSuite())
```

#### GRPC + HTTP2
```go
// server
server.WithSuite(tracing.NewGrpcServerSuite())
// client
client.WithSuite(tracing.NewGrpcClientSuite())
```



## Tracing associated Logs

#### set logger impl
```go
import (
    kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
)

func init()  {
    klog.SetLogger(kitexlogrus.NewLogger())
    klog.SetLevel(klog.LevelDebug)

}
```

#### log with context

```go
// Echo implements the Echo interface.
func (s *EchoImpl) Echo(ctx context.Context, req *api.Request) (resp *api.Response, err error) {
	klog.CtxDebugf(ctx, "echo called: %s", req.GetMessage())
	return &api.Response{Message: req.Message}, nil
}
```

#### view log

```log
{"level":"debug","msg":"echo called: my request","span_id":"056e0cf9a8b2cec3","time":"2022-03-09T02:47:28+08:00","trace_flags":"01","trace_id":"33bdd3c81c9eb6cbc0fbb59c57ce088b"}
```


## Example

[Executable Example](https://github.com/cloudwego/kitex-examples/tree/main/opentelemetry)

## Supported Metrics

### RPC Metrics

#### Kitex Server

Below is a table of RPC server metric instruments.

| Name | Instrument | Unit | Unit ([UCUM](README.md#instrument-units)) | Description | Status | Streaming |
|------|------------|------|-------------------------------------------|-------------|--------|-----------|
| `rpc.server.duration` | Histogram  | milliseconds | `ms` | measures duration of inbound RPC | Recommended | N/A.  While streaming RPCs may record this metric as start-of-batch to end-of-batch, it's hard to interpret in practice. |

#### Kitex Client

Below is a table of RPC client metric instruments.  These apply to traditional
RPC usage, not streaming RPCs.

| Name | Instrument | Unit | Unit ([UCUM](README.md#instrument-units)) | Description | Status | Streaming |
|------|------------|------|-------------------------------------------|-------------|--------|-----------|
| `rpc.client.duration` | Histogram | milliseconds | `ms` | measures duration of outbound RPC | Recommended | N/A.  While streaming RPCs may record this metric as start-of-batch to end-of-batch, it's hard to interpret in practice. |

#### Runtime Metrics

