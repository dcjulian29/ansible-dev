/*
Copyright © 2026 Julian Easterling

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"fmt"

	"github.com/dcjulian29/ansible-dev/internal/settings"
	"github.com/spf13/cobra"
)

// pathCmd creates "ansible-dev config path", which prints the absolute path of
// the configuration file. The path is printed whether or not the file exists
// yet, so it can be piped to an editor to open or create it.
func pathCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "Print the absolute path of the configuration file",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			path, err := settings.Path()
			if err != nil {
				return err
			}

			fmt.Println(path)

			return nil
		},
	}
}
