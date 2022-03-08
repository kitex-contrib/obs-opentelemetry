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
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	"github.com/cloudwego/kitex/transport"
)

var (
	_ client.Suite = (*clientSuite)(nil)
	_ server.Suite = (*serverSuite)(nil)
)

type clientSuite struct {
	opts []Option
}

func NewClientSuite(opts ...Option) *clientSuite {
	return &clientSuite{opts: opts}
}

func (c *clientSuite) Options() []client.Option {
	clientOpts, cfg := newClientOption(c.opts...)
	opts := []client.Option{
		clientOpts,
		client.WithMiddleware(ClientMiddleware(cfg)),
		client.WithTransportProtocol(transport.TTHeader),
		client.WithMetaHandler(transmeta.ClientTTHeaderHandler),
	}
	return opts
}

type serverSuite struct {
	opts []Option
}

func NewServerSuite(opts ...Option) *serverSuite {
	return &serverSuite{opts: opts}
}

func (s *serverSuite) Options() []server.Option {
	serverOpts, cfg := newServerOption(s.opts...)
	opts := []server.Option{
		serverOpts,
		server.WithMiddleware(ServerMiddleware(cfg)),
		server.WithMetaHandler(transmeta.ServerTTHeaderHandler),
	}
	return opts
}
