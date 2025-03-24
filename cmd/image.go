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

package cmd

import (
	"github.com/openconfig/containerz/client"
	"github.com/spf13/cobra"

	cpb "github.com/openconfig/gnoi/containerz"
)

var (
	containerzClient *client.Client
	cli              cpb.ContainerzClient
	image, tag       string
)

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "General image operations",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		containerzClient, err = client.NewClient(cmd.Context(), addr)
		return err
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	RootCmd.AddCommand(imageCmd)
	imageCmd.PersistentFlags().StringVar(&image, "image", "", "container image name")
	imageCmd.PersistentFlags().StringVar(&tag, "tag", "latest", "container image tag")
}
