package docker

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types"
	"google.golang.org/protobuf/testing/protocmp"
	cpb "github.com/openconfig/gnoi/containerz"
)

const (
	jsonConfig = `{
  "Args": {
    "Description": "",
    "Name": "",
    "Settable": null,
    "Value": null
  },
  "Description": "%s",
  "Documentation": "",
  "Entrypoint": null,
  "Env": null,
  "Interface": {
    "Socket": "",
    "Types": null
  },
  "IpcHost": false,
  "Linux": {
    "AllowAllDevices": false,
    "Capabilities": null,
    "Devices": null
  },
  "Mounts": null,
  "Network": {
    "Type": ""
  },
  "PidHost": false,
  "PropagatedMount": "",
  "User": {},
  "WorkDir": ""
}`
)

type fakeListPluginStreamer struct {
	msgs []*cpb.ListPluginsResponse
}

func (f *fakeListPluginStreamer) Send(msg *cpb.ListPluginsResponse) error {
	f.msgs = append(f.msgs, msg)
	return nil
}

type fakePluginListingDocker struct {
	fakeDocker
	plugins []*types.Plugin
}

func (f *fakePluginListingDocker) PluginList(ctx context.Context, args filters.Args) (types.PluginsListResponse, error) {
	values := args.Get("name")
	if len(values) == 0 {
		return f.plugins, nil
	}

	for _, plugin := range f.plugins {
		if plugin.Name == values[0] {
			return []*types.Plugin{plugin}, nil
		}
	}
	return nil, nil
}

func TestPluginList(t *testing.T) {
	tests := []struct {
		name        string
		inPlugin    []*types.Plugin
		inInstance  string
		wantPlugins []*cpb.ListPluginsResponse
	}{
		{
			name: "no-plugins",
		},
		{
			name: "one-plugin",
			inPlugin: []*types.Plugin{
				&types.Plugin{
					ID:   "plugin1",
					Name: "plugin1",
					Config: types.PluginConfig{
						Description: "plugin1 config",
					},
				},
			},
			wantPlugins: []*cpb.ListPluginsResponse{
				&cpb.ListPluginsResponse{
					Plugins: []*cpb.Plugin{
						&cpb.Plugin{
							Id:           "plugin1",
							InstanceName: "plugin1",
							Config:       fmt.Sprintf(jsonConfig, "plugin1 config"),
						},
					},
				},
			},
		},
		{
			name: "multiple-plugins",
			inPlugin: []*types.Plugin{
				&types.Plugin{
					ID:   "plugin1",
					Name: "plugin1",
					Config: types.PluginConfig{
						Description: "plugin1 config",
					},
				},
				&types.Plugin{
					ID:   "plugin2",
					Name: "plugin2",
					Config: types.PluginConfig{
						Description: "plugin2 config",
					},
				},
			},
			wantPlugins: []*cpb.ListPluginsResponse{
				&cpb.ListPluginsResponse{
					Plugins: []*cpb.Plugin{
						&cpb.Plugin{
							Id:           "plugin1",
							InstanceName: "plugin1",
							Config:       fmt.Sprintf(jsonConfig, "plugin1 config"),
						},
					},
				},
				&cpb.ListPluginsResponse{
					Plugins: []*cpb.Plugin{
						&cpb.Plugin{
							Id:           "plugin2",
							InstanceName: "plugin2",
							Config:       fmt.Sprintf(jsonConfig, "plugin2 config"),
						},
					},
				},
			},
		},
		{
			name: "multiple-plugins",
			inPlugin: []*types.Plugin{
				&types.Plugin{
					ID:   "plugin1",
					Name: "plugin1",
					Config: types.PluginConfig{
						Description: "plugin1 config",
					},
				},
				&types.Plugin{
					ID:   "plugin2",
					Name: "plugin2",
					Config: types.PluginConfig{
						Description: "plugin2 config",
					},
				},
			},
			inInstance: "plugin1",
			wantPlugins: []*cpb.ListPluginsResponse{
				&cpb.ListPluginsResponse{
					Plugins: []*cpb.Plugin{
						&cpb.Plugin{
							Id:           "plugin1",
							InstanceName: "plugin1",
							Config:       fmt.Sprintf(jsonConfig, "plugin1 config"),
						},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			fsd := &fakePluginListingDocker{
				plugins: tc.inPlugin,
			}
			mgr := New(fsd)

			stream := &fakeListPluginStreamer{}
			if err := mgr.PluginList(ctx, tc.inInstance, stream); err != nil {
				t.Errorf("PluginList(%q) returned error: %v", tc.inInstance, err)
			}

			if diff := cmp.Diff(tc.wantPlugins, stream.msgs, protocmp.Transform()); diff != "" {
				fmt.Println(stream.msgs)
				t.Errorf("PluginList(%q) returned diff(-want, +got):\n%s", tc.inInstance, diff)
			}
		})
	}
}
