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

// ignoreCmd builds a command group (used for both "role-ignore" and
// "runbook-ignore") that manages one ignore list on the configuration. The get
// and set closures read and write the specific slice field so the same
// list/add/remove/clear plumbing serves both lists.
func ignoreCmd(
	use, short string,
	get func(*settings.Config) []string,
	set func(*settings.Config, []string),
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
	}

	cmd.AddCommand(ignoreListCmd(get))
	cmd.AddCommand(ignoreAddCmd(get, set))
	cmd.AddCommand(ignoreRemoveCmd(get, set))
	cmd.AddCommand(ignoreClearCmd(set))

	return cmd
}

func ignoreListCmd(get func(*settings.Config) []string) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List the ignored path substrings",
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg, err := settings.Load()
			if err != nil {
				return err
			}

			for _, v := range get(&cfg) {
				fmt.Println(v)
			}

			return nil
		},
	}
}

func ignoreAddCmd(
	get func(*settings.Config) []string,
	set func(*settings.Config, []string),
) *cobra.Command {
	return &cobra.Command{
		Use:   "add <substring>",
		Short: "Add a path substring to ignore",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			cfg, err := settings.Load()
			if err != nil {
				return err
			}

			list := get(&cfg)
			for _, v := range list {
				if v == args[0] {
					return nil // already present; nothing to do
				}
			}

			set(&cfg, append(list, args[0]))

			if err := settings.Save(&cfg); err != nil {
				return err
			}

			fmt.Println(textformat.Info(fmt.Sprintf("added '%s'", args[0])))

			return nil
		},
	}
}

func ignoreRemoveCmd(
	get func(*settings.Config) []string,
	set func(*settings.Config, []string),
) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <substring>",
		Short: "Remove a path substring from the ignore list",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			cfg, err := settings.Load()
			if err != nil {
				return err
			}

			kept := make([]string, 0, len(get(&cfg)))
			for _, v := range get(&cfg) {
				if v != args[0] {
					kept = append(kept, v)
				}
			}

			set(&cfg, kept)

			if err := settings.Save(&cfg); err != nil {
				return err
			}

			fmt.Println(textformat.Info(fmt.Sprintf("removed '%s'", args[0])))

			return nil
		},
	}
}

func ignoreClearCmd(set func(*settings.Config, []string)) *cobra.Command {
	return &cobra.Command{
		Use:   "clear",
		Short: "Remove all entries from the ignore list",
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg, err := settings.Load()
			if err != nil {
				return err
			}

			set(&cfg, nil)

			if err := settings.Save(&cfg); err != nil {
				return err
			}

			fmt.Println(textformat.Info("ignore list cleared"))

			return nil
		},
	}
}
