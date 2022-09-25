# opentelemetry (这是一个社区驱动的项目)

[English](README.md) | 中文

适用于 [Kitex](https://github.com/cloudwego/kitex) 的 [Opentelemetry](https://opentelemetry.io/) 

## 特性

#### 提供者

- [x] 集成的默认 opentelemetry 程序，达到开箱即用
- [x] 支持设置环境变量

### 遥测工具

#### 链路追踪

- [x] 支持在 kitex 服务端和客户端中的 rpc 链路追踪
- [x] 支持通过元信息自动透明传输对等服务

#### 指标

- [x] 支持 kitex rpc 指标 [Rate, Errors, Duration]
- [x] 支持服务拓扑图度量[服务拓扑图]。
- [x] 支持go runtime 指标

#### 日志

- [x] 在logrus和zap的基础上扩展 kitex 日志工具 
- [x] 实现跟踪自动关联日志



## 通过环境变量来配置

- [Exporter](https://opentelemetry.io/docs/reference/specification/protocol/exporter/)
- [SDK](https://opentelemetry.io/docs/reference/specification/sdk-environment-variables/#general-sdk-configuration)

## 服务端使用示例

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

## 客户端使用示例

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

##  追踪相关日志

## 设置日志

```go
import (
    kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
)

func init()  {
    klog.SetLogger(kitexlogrus.NewLogger())
    klog.SetLevel(klog.LevelDebug)

}
```

## 结合 context 使用日志

```go
// Echo implements the Echo interface.
func (s *EchoImpl) Echo(ctx context.Context, req *api.Request) (resp *api.Response, err error) {
	klog.CtxDebugf(ctx, "echo called: %s", req.GetMessage())
	return &api.Response{Message: req.Message}, nil
}
```

#### 日志格式示例

```log
{"level":"debug","msg":"echo called: my request","span_id":"056e0cf9a8b2cec3","time":"2022-03-09T02:47:28+08:00","trace_flags":"01","trace_id":"33bdd3c81c9eb6cbc0fbb59c57ce088b"}
```



## 示例

[Executable Example](https://github.com/cloudwego/kitex-examples/tree/main/opentelemetry)



## 现已支持的 Mertrics

### RPC Metrics

#### Kitex Server

下面的表格为 RPC server metric 的配置项。

| 名称                  | 指标数据模型 | 单位        | 单位(UCUM) | 描述                  | 状态     | Streaming                                                    |
| --------------------- | ------------ | ----------- | ---------- | --------------------- | -------- | ------------------------------------------------------------ |
| `rpc.server.duration` | Histogram    | millseconds | `ms`       | 测量请求RPC的持续时间 | 推荐使用 | 并不适用， 虽然streaming RPC可能将这个指标记录为*批处理开始到批处理结束*，但在实际使用中很难解释。 |



#### Kitex Client

下面的表格为 RPC server metric 的配置项,这些适用于传统的RPC使用，不支持 streaming RPC

| 名称                  | 指标数据模型 | 单位        | 单位(UCUM) | 描述                  | 状态     | Streaming                                                    |
| --------------------- | ------------ | ----------- | ---------- | --------------------- | -------- | ------------------------------------------------------------ |
| `rpc.server.duration` | Histogram    | millseconds | `ms`       | 测量请求RPC的持续时间 | 推荐使用 | 并不适用， 虽然streaming RPC可能将这个指标记录为*批处理开始到批处理结束*，但在实际使用中很难解释。 |



### R.E.D

R.E.D (Rate, Errors, Duration) 方法定义了你应该为你架构中的每个微服务测量的三个关键指标。我们可以根据`rpc.server.duration`来计算R.E.D。

#### Rate

你的服务每秒钟所提供的请求数。

例如: QPS（Queries Per Second）每秒查询率
```
sum(rate(rpc_server_duration_count{}[5m])) by (service_name, rpc_method)
```

#### Errors
每秒失败的请求数。

例如：错误率
```
sum(rate(rpc_server_duration_count{status_code="Error"}[5m])) by (service_name, rpc_method) / sum(rate(rpc_server_duration_count{}[5m])) by (service_name, rpc_method)
```

#### Duration
每个请求所需时间的分布情况

例如：[P99 Latency](https://stackoverflow.com/questions/12808934/what-is-p99-latency) 

```
histogram_quantile(0.99, sum(rate(rpc_server_duration_bucket{}[5m])) by (le, service_name, rpc_method))
```

### 服务拓扑图

 `rpc.server.duration`将记录对等服务和当前服务维度。基于这个维度，我们可以汇总服务拓扑图

```
sum(rate(rpc_server_duration_count{}[5m])) by (service_name, peer_service)
```

### Runtime Metrics

| 名称                                   | 指标数据模型 | 单位       | 单位(UCUM) | 描述                                |
| -------------------------------------- | ------------ | ---------- | ---------- |-----------------------------------|
| `process.runtime.go.cgo.calls`         | Sum          | -          | -          | 当前进程调用的cgo数量                      |
| `process.runtime.go.gc.count`          | Sum          | -          | -          | 已完成的 gc 周期的数量                     |
| `process.runtime.go.gc.pause_ns`       | Histogram    | nanosecond | `ns`       | 在GC stop-the-world 中暂停的纳秒数量       |
| `process.runtime.go.gc.pause_total_ns` | Histogram    | nanosecond | `ns`       | 自程序启动以来，GC stop-the-world 的累计微秒计数 |
| `process.runtime.go.goroutines`        | Gauge        | -          | -          | 协程数量                              |
| `process.runtime.go.lookups`           | Sum          | -          | -          | 运行时执行的指针查询的数量                     |
| `process.runtime.go.mem.heap_alloc`    | Gauge        | bytes      | `bytes`    | 分配的堆对象的字节数                        |
| `process.runtime.go.mem.heap_idle`     | Gauge        | bytes      | `bytes`    | 空闲（未使用）的堆内存                       |
| `process.runtime.go.mem.heap_inuse`    | Gauge        | bytes      | `bytes`    | 已使用的堆内存                           |
| `process.runtime.go.mem.heap_objects`  | Gauge        | -          | -          | 已分配的堆对象数量                         |
| `process.runtime.go.mem.live_objects`  | Gauge        | -          | -          | 存活对象数量(Mallocs - Frees)           |
| `process.runtime.go.mem.heap_released` | Gauge        | bytes      | `bytes`    | 已交还给操作系统的堆内存                      |
| `process.runtime.go.mem.heap_sys`      | Gauge        | bytes      | `bytes`    | 从操作系统获得的堆内存                       |
| `runtime.uptime`                       | Sum          | ms         | `ms`       | 自应用程序被初始化以来的毫秒数                   |

##  兼容性

OpenTelemetry的 sdk 与1.x opentelemetry-go完全兼容，[详情查看](https://github.com/open-telemetry/opentelemetry-go#compatibility)


维护者: [CoderPoet](https://github.com/CoderPoet)

## 依赖

| **库/框架**                                         | 版本    | 记录   |
| --------------------------------------------------- | ------- | ------ |
| go.opentelemetry.io/otel                            | v1.7.0  | <br /> |
| go.opentelemetry.io/otel/trace                      | v1.7.0  | <br /> |
| go.opentelemetry.io/otel/metric                     | v0.30.0 | <br /> |
| go.opentelemetry.io/otel/semconv                    | v1.7.0  |        |
| go.opentelemetry.io/contrib/instrumentation/runtime | v0.30.0 |        |
| kitex                                               | v0.3.1  |        |
