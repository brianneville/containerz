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
	"io/fs"
	"os"

	cpb "github.com/openconfig/gnoi/containerz"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RemovePluginTar removes a plugin tar. If the plugin does not exist this operation is a no-op.
func (s *Server) RemovePluginTar(ctx context.Context, request *cpb.RemovePluginTarRequest) (*cpb.RemovePluginTarResponse, error) {
	pluginName := request.GetName()
	if pluginName == "" {
		return nil, status.Error(codes.InvalidArgument, "name field must be set")
	}

	pluginPath := genPluginPath(pluginName)
	stat, err := os.Stat(pluginPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, status.Errorf(codes.NotFound,
				"requested plugin tar %q was not found at %s",
				request.GetName(), pluginPath)
		}
		return nil, err
	}
	if stat.IsDir() { // this probably shouldn't happen but lets check it just in case
		return nil, status.Errorf(codes.Unknown,
			"expected to find plugin tar at %s, found dir",
			pluginPath)
	}
	if err := os.Remove(pluginPath); err != nil {
		return nil, status.Errorf(codes.Unknown,
			"failed to remove plugin tar at %s with error %s",
			pluginPath, err)
	}

	return &cpb.RemovePluginTarResponse{}, nil
}
