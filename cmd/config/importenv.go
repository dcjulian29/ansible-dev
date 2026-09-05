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
	"os"

	"github.com/dcjulian29/ansible-dev/internal/settings"
	"github.com/dcjulian29/go-toolbox/textformat"
	"github.com/spf13/cobra"
)

// legacyImport describes one environment variable that used to supply a
// setting, paired with accessors for the configuration field that replaced it.
type legacyImport struct {
	env     string
	key     string
	current func(*settings.Config) string
	set     func(*settings.Config, string)
}

// legacyImports lists the environment variables read by earlier versions of
// ansible-dev, in the order they are reported.
var legacyImports = []legacyImport{
	{
		env:     settings.LegacyRolesEnv,
		key:     "roles_path",
		current: func(c *settings.Config) string { return c.RolesPath },
		set:     func(c *settings.Config, v string) { c.RolesPath = v },
	},
	{
		env:     settings.LegacyRunbooksEnv,
		key:     "runbooks_path",
		current: func(c *settings.Config) string { return c.RunbooksPath },
		set:     func(c *settings.Config, v string) { c.RunbooksPath = v },
	},
}

// importEnvCmd creates "ansible-dev config import-env", a one-time migration
// helper that copies the ANSIBLE_ROLES and ANSIBLE_RUNBOOKS environment
// variables (read by earlier versions) into the configuration file.
//
// A variable that is unset is reported and skipped. A setting that already has
// a value is left alone unless --force is given, so re-running the command is
// safe. The configuration file is written only when something actually changed.
func importEnvCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "import-env",
		Short: "Import the legacy ANSIBLE_ROLES and ANSIBLE_RUNBOOKS variables",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg, err := settings.Load()
			if err != nil {
				return err
			}

			// Tracked so the directory warnings below cover only what this run
			// actually imported, not unrelated settings already in the file.
			var imported []legacyImport

			for _, l := range legacyImports {
				value := os.Getenv(l.env)

				switch {
				case value == "":
					fmt.Println(textformat.Warn(
						fmt.Sprintf("%s is not set; skipping %s", l.env, l.key)))
				case l.current(&cfg) != "" && !force:
					fmt.Println(textformat.Warn(fmt.Sprintf(
						"%s is already set to '%s'; leaving it unchanged (use --force to overwrite)",
						l.key, l.current(&cfg))))
				default:
					l.set(&cfg, value)

					imported = append(imported, l)

					fmt.Println(textformat.Info(
						fmt.Sprintf("%s set to '%s' from %s", l.key, value, l.env)))
				}
			}

			if len(imported) == 0 {
				fmt.Println(textformat.Warn("nothing was imported; the configuration is unchanged"))

				return nil
			}

			if err := settings.Save(&cfg); err != nil {
				return err
			}

			fmt.Println(textformat.Info(fmt.Sprintf(
				"imported %d setting(s); %s and %s are no longer read and can be unset",
				len(imported), settings.LegacyRolesEnv, settings.LegacyRunbooksEnv)))

			// An environment variable left over from an older machine can
			// easily name a directory that is gone, so check what was imported
			// rather than trusting it.
			for _, l := range imported {
				warnAboutDirectory(l.key, l.current(&cfg))
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false,
		"overwrite settings that already have a value")

	return cmd
}
