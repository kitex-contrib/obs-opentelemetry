# opentelemetry (This is a community driven project)

[Opentelemetry](https://opentelemetry.io/) for [Kitex](https://github.com/cloudwego/kitex)

## Feature
#### Provider
- [x] Out-of-the-box  default opentelemetry provider
- [x] Support setting via environment variables

### Instrumentation

#### Tracing
- [x] Support server and client kitex rpc tracing
- [x] Support automatic transparent transmission of peer service through meta info

#### Metrics
- [x] Support kitex rpc metrics [R.E.D]
- [x] Support service topology map metrics [Service Topology Map]
- [x] Support go runtime metrics

#### Logging
- [x] Extend kitex logger based on logrus
- [x] Implement tracing auto associated logs

## Configuration via environment variables
- [Exporter](https://opentelemetry.io/docs/reference/specification/protocol/exporter/)
- [SDK](https://opentelemetry.io/docs/reference/specification/sdk-environment-variables/#general-sdk-configuration)

## Server usage
```go
import (
    ...
    "github.com/kitex-contrib/obs-opentelemetry/provider"
    "github.com/kitex-contrib/obs-opentelemetry/tracing"
)


func main()  {
    serviceName := "echo"
	
    p := provider.NewOpenTelemetryProvider(
        provider.WithServiceName(serviceName),
        provider.WithExportEndpoint("localhost:4317"),
        provider.WithInsecure(),
    )
    defer p.Shutdown(context.Background())

    svr := echo.NewServer(
        new(EchoImpl),
        server.WithSuite(tracing.NewServerSuite()),
        // Please keep the same as provider.WithServiceName
        server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: serviceName}),
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
    serviceName := "echo-client"
	
    p := provider.NewOpenTelemetryProvider(
        provider.WithServiceName(serviceName),
        provider.WithExportEndpoint("localhost:4317"),
        provider.WithInsecure(),
    )
    defer p.Shutdown(context.Background())
    
    c, err := echo.NewClient(
        "echo",
        client.WithSuite(tracing.NewClientSuite()),
        // Please keep the same as provider.WithServiceName
        client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: serviceName}),
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
server.WithSuite(tracing.NewGRPCServerSuite())
// client
client.WithSuite(tracing.NewGRPCClientSuite())
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

| Name | Instrument | Unit | Unit (UCUM) | Description | Status | Streaming |
|------|------------|------|-------------------------------------------|-------------|--------|-----------|
| `rpc.server.duration` | Histogram  | milliseconds | `ms` | measures duration of inbound RPC | Recommended | N/A.  While streaming RPCs may record this metric as start-of-batch to end-of-batch, it's hard to interpret in practice. |

#### Kitex Client

Below is a table of RPC client metric instruments.  These apply to traditional
RPC usage, not streaming RPCs.

| Name | Instrument | Unit | Unit (UCUM) | Description | Status | Streaming |
|------|------------|------|-------------------------------------------|-------------|--------|-----------|
| `rpc.client.duration` | Histogram | milliseconds | `ms` | measures duration of outbound RPC | Recommended | N/A.  While streaming RPCs may record this metric as start-of-batch to end-of-batch, it's hard to interpret in practice. |


### R.E.D
The RED Method defines the three key metrics you should measure for every microservice in your architecture. We can calculate RED based on `rpc.server.duration`.

#### Rate
the number of requests, per second, you services are serving.

eg: QPS
```
sum(rate(rpc_server_duration_count{}[5m])) by (service_name, rpc_method)
```

#### Errors
the number of failed requests per second.

eg: Error ratio
```
sum(rate(rpc_server_duration_count{status_code="Error"}[5m])) by (service_name, rpc_method) / sum(rate(rpc_server_duration_count{}[5m])) by (service_name, rpc_method)
```

#### Duration
distributions of the amount of time each request takes

eg: P99 Latency
```
histogram_quantile(0.99, sum(rate(rpc_server_duration_bucket{}[5m])) by (le, service_name, rpc_method))
```

### Service Topology Map
The `rpc.server.duration` will record the peer service and the current service dimension. Based on this dimension, we can aggregate the service topology map
```
sum(rate(rpc_server_duration_count{}[5m])) by (service_name, peer_service)
```

### Runtime Metrics
| Name | Instrument | Unit | Unit (UCUM)) | Description |
|------|------------|------|-------------------------------------------|-------------|
| `runtime.go.cgo.calls` | Sum | - | - | Number of cgo calls made by the current process. |
| `runtime.go.gc.count` | Sum | - | - | Number of completed garbage collection cycles. |
| `runtime.go.gc.pause_ns` | Histogram | nanosecond | `ns` | Amount of nanoseconds in GC stop-the-world pauses. |
| `runtime.go.gc.pause_total_ns` | Histogram | nanosecond | `ns` | Cumulative nanoseconds in GC stop-the-world pauses since the program started. |
| `runtime.go.goroutines` | Gauge | - | - | measures duration of outbound RPC. | 
| `runtime.go.lookups` | Sum | - | - | Number of pointer lookups performed by the runtime. |
| `runtime.go.mem.heap_alloc` | Gauge | bytes | `bytes` | Bytes of allocated heap objects. |
| `runtime.go.mem.heap_idle` | Gauge | bytes | `bytes` | Bytes in idle (unused) spans. |
| `runtime.go.mem.heap_inuse` | Gauge | bytes | `bytes` | Bytes in in-use spans. |
| `runtime.go.mem.heap_objects` | Gauge | - | - | Number of allocated heap objects. |
| `runtime.go.mem.live_objects` | Gauge | - | - | Number of live objects is the number of cumulative Mallocs - Frees. |
| `runtime.go.mem.heap_released` | Gauge | bytes | `bytes` | Bytes of idle spans whose physical memory has been returned to the OS. |
| `runtime.go.mem.heap_sys` | Gauge | bytes | `bytes` | Bytes of idle spans whose physical memory has been returned to the OS. |
| `runtime.uptime` | Sum | ms | `ms` | Milliseconds since application was initialized. |


## Compatibility
The sdk of OpenTelemetry is fully compatible with 1.X opentelemetry-go. [see](https://github.com/open-telemetry/opentelemetry-go#compatibility)


maintained by: [CoderPoet](https://github.com/CoderPoet)


