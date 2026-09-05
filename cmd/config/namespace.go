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

	"github.com/dcjulian29/ansible-dev/internal/ansible"
	"github.com/dcjulian29/ansible-dev/internal/settings"
	"github.com/dcjulian29/go-toolbox/textformat"
	"github.com/spf13/cobra"
)

// namespaceCmd creates "ansible-dev config namespace [name]". With no argument
// it prints the configured namespace; with a name it validates, sets and saves
// it.
//
// Unlike the directory settings this one fails rather than warns on a bad
// value. A directory can reasonably be configured before it exists, but a
// namespace containing "." or "/" is not a value that becomes correct later —
// it would silently produce an unresolvable role name or a malformed
// repository.
func namespaceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "namespace [name]",
		Short: "Show or set the Galaxy/GitHub namespace for new roles",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			cfg, err := settings.Load()
			if err != nil {
				return err
			}

			if len(args) == 0 {
				fmt.Println(cfg.Namespace)

				return nil
			}

			if err := ansible.ValidateNamespace(args[0]); err != nil {
				return err
			}

			cfg.Namespace = args[0]

			if err := settings.Save(&cfg); err != nil {
				return err
			}

			fmt.Println(textformat.Info(fmt.Sprintf("namespace set to '%s'", cfg.Namespace)))

			return nil
		},
	}
}
