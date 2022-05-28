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
	cOpts []client.Option
}

func (c *clientSuite) Options() []client.Option {
	return c.cOpts
}

type serverSuite struct {
	sOpts []server.Option
}

func (s *serverSuite) Options() []server.Option {
	return s.sOpts
}

// NewClientSuite client suite for otel with http2 and ttheader meta handler
func NewClientSuite(opts ...Option) *clientSuite {
	clientOpts, cfg := newClientOption(opts...)
	cOpts := []client.Option{
		clientOpts,
		client.WithMiddleware(ClientMiddleware(cfg)),
		client.WithTransportProtocol(transport.TTHeader),
		client.WithMetaHandler(transmeta.ClientHTTP2Handler),
		client.WithMetaHandler(transmeta.ClientTTHeaderHandler),
	}
	return &clientSuite{cOpts}
}

// NewServerSuite server suite for otel with http2 and ttheader meta handler
func NewServerSuite(opts ...Option) *serverSuite {
	serverOpts, cfg := newServerOption(opts...)
	sOpts := []server.Option{
		serverOpts,
		server.WithMiddleware(ServerMiddleware(cfg)),
		server.WithMetaHandler(transmeta.ServerHTTP2Handler),
		server.WithMetaHandler(transmeta.ServerTTHeaderHandler),
	}

	return &serverSuite{sOpts}
}

// Deprecated: Use NewServerSuite instead.
func NewGRPCServerSuite(opts ...Option) *serverSuite {
	serverOpts, cfg := newServerOption(opts...)
	sOpts := []server.Option{
		serverOpts,
		server.WithMiddleware(ServerMiddleware(cfg)),
		server.WithMetaHandler(transmeta.ServerHTTP2Handler),
	}

	return &serverSuite{sOpts}
}

// Deprecated: Use NewClientSuite instead.
func NewGRPCClientSuite(opts ...Option) *clientSuite {
	clientOpts, cfg := newClientOption(opts...)
	cOpts := []client.Option{
		clientOpts,
		client.WithMiddleware(ClientMiddleware(cfg)),
		client.WithTransportProtocol(transport.GRPC),
		client.WithMetaHandler(transmeta.ClientHTTP2Handler),
		client.WithMetaHandler(transmeta.ClientTTHeaderHandler),
	}
	return &clientSuite{cOpts}
}

// Deprecated: Use NewClientSuite instead.
func NewFramedClientSuite(opts ...Option) *clientSuite {
	clientOpts, cfg := newClientOption(opts...)
	cOpts := []client.Option{
		clientOpts,
		client.WithMiddleware(ClientMiddleware(cfg)),
		client.WithTransportProtocol(transport.Framed),
		client.WithMetaHandler(transmeta.ClientTTHeaderHandler),
	}
	return &clientSuite{cOpts}
}
