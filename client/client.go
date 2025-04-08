// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package client is a containerz grpc client.
package client

import (
	"context"

	cpb "github.com/openconfig/gnoi/containerz"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	// Dial is the dailer to use to read containerz
	Dial = grpc.DialContext
)

// Client is a grpc containerz client.
type Client struct {
	cli cpb.ContainerzClient
}

// NewClient builds a new containerz client.
func NewClient(ctx context.Context, addr string) (*Client, error) {
	tlsCred := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := Dial(ctx, addr, tlsCred)
	if err != nil {
		return nil, err
	}

	return &Client{
		cli: cpb.NewContainerzClient(conn),
	}, nil
}

// NewClientWithConn creates a client given a ClientConn.
func NewClientWithConn(conn *grpc.ClientConn) *Client {
	return &Client{
		cli: cpb.NewContainerzClient(conn),
	}
}

// NewClientFromStub allows the creation of a client using a client
// obtained via gnoigo.
func NewClientFromStub(c cpb.ContainerzClient) *Client {
	return &Client{
		cli: c,
	}
}
