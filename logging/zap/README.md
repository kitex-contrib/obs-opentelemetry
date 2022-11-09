# Kitex zap (This is a community driven project)

## Introduction

This is a logger library that uses [zap](https://github.com/uber-go/zap) to implement the [Kitex logger interface](https://www.cloudwego.io/docs/kitex/tutorials/basic-feature/logging/), work together with kitex [obs-opentelemetry](https://github.com/kitex-contrib/obs-opentelemetry)

## Usage

Download and install it:

```go
go get github.com/kitex-contrib/obs-opentelemetry/logging/zap
```

Import it in your code:

```go
import kitexzap github.com/kitex-contrib/obs-opentelemetry/logging/zap
```

### Set logger impl

```go
package main

import (
    "github.com/cloudwego/kitex/pkg/klog"
    kitexzap "github.com/kitex-contrib/obs-opentelemetry/logging/zap"
)

func main() {
    klog.SetLogger(kitexzap.NewLogger())
    klog.SetLevel(klog.LevelDebug)
}
```

> We provide some methods to help you customize logger

| Configuration              | Description                                                                   |
| -------------------------- | ----------------------------------------------------------------------------- |
| WithCoreEnc                | [zapcore Encoder](https://pkg.go.dev/go.uber.org/zap/zapcore#Encoder)         |
| WithCoreWs                 | [zapcore WriteSyncer](https://pkg.go.dev/go.uber.org/zap/zapcore#WriteSyncer) |
| WithCoreLevel              | [zap AtomicLevel](https://pkg.go.dev/go.uber.org/zap#AtomicLevel)             |
| WithZapOptions             | [zap Option](https://pkg.go.dev/go.uber.org/zap#Option)                       |
| WithTraceErrorSpanLevel    | trace error span level option                                                 |
| WithRecordStackTraceInSpan | record stack track option                                                     |

### Log with context

```go
// Echo implements the Echo interface.
func (s *EchoImpl) Echo(ctx context.Context, req *api.Request) (resp *api.Response, err error) {
    klog.CtxDebugf(ctx, "echo called: %s", req.GetMessage())
    return &api.Response{Message: req.Message}, nil
}
```

### View log

```json
{
  "level": "info",
  "ts": 1667619647.1459548,
  "msg": "hello world",
  "trace_id": "c77e46c0fb590ee80b6d78ed6682768e",
  "span_id": "b42c96c6dd01ceaf",
  "trace_flags": "01"
}
```

> For some reason, zap will not log extra context info.
