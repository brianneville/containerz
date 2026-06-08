// Copyright 2026 Google LLC
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

package client

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	cpb "github.com/openconfig/gnoi/containerz"
	"google.golang.org/protobuf/testing/protocmp"
)

type fakeRemovePluginTarServer struct {
	fakeContainerzServer

	recvMsg *cpb.RemovePluginTarRequest
}

func (f *fakeRemovePluginTarServer) RemovePluginTar(ctx context.Context, req *cpb.RemovePluginTarRequest) (*cpb.RemovePluginTarResponse, error) {
	f.recvMsg = req
	return nil, nil
}

func TestRemovePluginTar(t *testing.T) {
	tests := []struct {
		name       string
		pluginName string
		wantReq    *cpb.RemovePluginTarRequest
	}{
		{
			name:       "remove-plugin",
			pluginName: "test",
			wantReq: &cpb.RemovePluginTarRequest{
				Name: "test",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			fcm := &fakeRemovePluginTarServer{}
			addr, stop := newServer(t, fcm)
			defer stop()
			cli, err := NewClient(ctx, addr)
			if err != nil {
				t.Fatalf("NewClient(%v) returned an unexpected error: %v",
					addr, err)
			}

			if err := cli.RemovePluginTar(ctx, tc.pluginName); err != nil {
				t.Fatalf("RemovePluginTar(%v) returned an unexpected error: %v",
					tc.pluginName, err)
			}
			if diff := cmp.Diff(tc.wantReq, fcm.recvMsg, protocmp.Transform()); diff != "" {
				t.Errorf("RemovePlugin(%v) returned an unexpected diff (-want +got):\n%s",
					tc.pluginName, diff)
			}
		})
	}
}
