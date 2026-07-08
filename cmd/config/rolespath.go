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
	"github.com/dcjulian29/go-toolbox/textformat"
	"github.com/spf13/cobra"
)

// rolesPathCmd creates "ansible-dev config roles-path [directory]". With no
// argument it prints the configured roles path; with a directory it sets and
// saves it.
func rolesPathCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "roles-path [directory]",
		Short: "Show or set the roles repository directory",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			cfg, err := settings.Load()
			if err != nil {
				return err
			}

			if len(args) == 0 {
				fmt.Println(cfg.RolesPath)

				return nil
			}

			cfg.RolesPath = args[0]

			if err := settings.Save(&cfg); err != nil {
				return err
			}

			fmt.Println(textformat.Info(fmt.Sprintf("roles_path set to '%s'", cfg.RolesPath)))

			return nil
		},
	}
}
