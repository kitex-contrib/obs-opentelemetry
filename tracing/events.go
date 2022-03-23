// Copyright 2022 CloudWeGo Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tracing

import (
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/stats"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var commonEvents = map[string]stats.Event{
	"server_handle_start":  stats.ServerHandleStart,
	"server_handle_finish": stats.ServerHandleFinish,
	"client_conn_start":    stats.ClientConnStart,
	"client_conn_finish":   stats.ClientConnFinish,
	"read_start":           stats.ReadStart,
	"read_finish":          stats.ReadFinish,
	"wait_read_start":      stats.WaitReadStart,
	"wait_read_finish":     stats.WaitReadFinish,
	"write_start":          stats.WriteStart,
	"write_finish":         stats.WriteFinish,
}

func injectStatsEventsToSpan(span trace.Span, st rpcinfo.RPCStats) {
	for name, event := range commonEvents {
		if gotEvent := st.GetEvent(event); gotEvent != nil {
			attrs := []attribute.KeyValue{attribute.Int("event.status", int(gotEvent.Status()))}
			if gotEvent.Info() != "" {
				attrs = append(attrs, attribute.String("event.info", gotEvent.Info()))
			}
			span.AddEvent(name,
				trace.WithTimestamp(gotEvent.Time()),
				trace.WithAttributes(attrs...),
			)
		}
	}
}
