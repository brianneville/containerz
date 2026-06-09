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

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var pluginName string

var pluginRemoveTarCmd = &cobra.Command{
	Use:   "remove-tar",
	Short: "Removes a plugin tar by deployed name",
	RunE: func(command *cobra.Command, args []string) error {
		if pluginName == "" {
			fmt.Println("--name must be provided")
		}

		if err := containerzClient.RemovePluginTar(command.Context(), pluginName); err != nil {
			return err
		}

		fmt.Printf("Successfully removed %s\n", instance)
		return nil
	},
}

func init() {
	pluginCmd.AddCommand(pluginRemoveTarCmd)
	pluginRemoveTarCmd.PersistentFlags().StringVar(&pluginName,
		"name", "", "plugin tarball name (as used in Deploy RPC)")
}
