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

package server

import (
	"context"
	"errors"
	"os"
	"testing"

	cpb "github.com/openconfig/gnoi/containerz"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRemovePluginTar(t *testing.T) {
	defer func(oldPluginLocation string) {
		pluginLocation = oldPluginLocation
	}(pluginLocation)
	pluginLocation = t.TempDir()

	for _, tc := range []struct {
		name        string
		pluginName  string
		tarExists   bool
		tarIsDir    bool
		expectedErr error
	}{{
		name:        "name_unset",
		expectedErr: status.Error(codes.InvalidArgument, "name field must be set"),
	}, {
		name:       "tar_not_exist",
		pluginName: "something-else",
		expectedErr: status.Errorf(codes.NotFound,
			"requested plugin tar %q was not found at %s/%[1]s.tar",
			"something-else", pluginLocation),
	}, {
		name:       "dir_exist",
		pluginName: "the-dir",
		tarIsDir:   true,
		expectedErr: status.Errorf(codes.Unknown,
			"expected to find plugin tar at %s/%s.tar, found dir",
			pluginLocation, "the-dir"),
	}, {
		name:       "tar_exist",
		pluginName: "the-tar",
		tarExists:  true,
	}} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.tarIsDir {
				if err := os.Mkdir(genPluginPath(tc.pluginName),
					os.ModePerm|os.ModeDir); err != nil {
					t.Fatal(err)
				}
			} else if tc.tarExists {
				if err := os.WriteFile(
					genPluginPath(tc.pluginName), []byte("this"),
					os.ModePerm); err != nil {
					t.Fatal(err)
				}
			}
			ctx := context.Background()
			cli, s := startServerAndReturnClient(ctx, t, &fakeContainerManager{}, nil)
			defer s.Halt(ctx)
			resp, err := cli.RemovePluginTar(ctx, &cpb.RemovePluginTarRequest{
				Name: tc.pluginName,
			})
			if err != nil && resp != nil {
				t.Errorf("expected no response when error is non-nil, got %#v", resp)
			}
			if !errors.Is(err, tc.expectedErr) {
				t.Fatalf("want error %v, got error %v", tc.expectedErr, err)
			}
		})
	}
}
