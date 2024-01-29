# Kitex zerolog (This is a community driven project)

## Introduction

This is a logger library that uses [zerolog](https://github.com/rs/zerolog) to implement the [Kitex logger interface](https://www.cloudwego.io/docs/kitex/tutorials/basic-feature/logging/), work together with kitex [obs-opentelemetry](https://github.com/kitex-contrib/obs-opentelemetry)

## Usage

Download and install it:

```go
go get github.com/kitex-contrib/obs-opentelemetry/logging/zerolog
```

Import it in your code:

```go
import kitexzerolog github.com/kitex-contrib/obs-opentelemetry/logging/zerolog
```

### Set logger impl

```go
package main

import (
    "github.com/rs/zerolog/log"
    "github.com/cloudwego/kitex/pkg/klog"
    kitexzerolog "github.com/kitex-contrib/obs-opentelemetry/logging/zerolog"
)

func main() {
    logger := kitexzerolog.NewLogger()
    klog.SetLogger(logger)
    klog.SetLevel(klog.LevelDebug)

    // OR / AND using global logger
    log.Logger = *logger.Logger()
}
```

> We provide some methods to help you customize logger

| Configuration              | Description                                                                   |
| -------------------------- | ----------------------------------------------------------------------------- |
| WithLogger                 | [Logger](https://pkg.go.dev/github.com/rs/zerolog#Logger)                       |
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

### Log with global logger

In case your code base has already used global logger from zerolog

```go
import (
    "github.com/rs/zerolog/log"
)

// Echo implements the Echo interface.
func (s *EchoImpl) Echo(ctx context.Context, req *api.Request) (resp *api.Response, err error) {
    log.Debug().Ctx(ctx).Msgf("echo called: %s", req.GetMessage())
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

> For some reason, zerolog will not log extra context info.
