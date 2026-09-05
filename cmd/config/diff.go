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
	"runtime"
	"strings"

	"github.com/dcjulian29/ansible-dev/internal/settings"
	"github.com/dcjulian29/go-toolbox/textformat"
	"github.com/spf13/cobra"
)

// diffScalarCmd builds a show/set command for one scalar diff setting of the
// current operating system. With no argument it prints the current value; with
// one it sets and saves it.
func diffScalarCmd(
	use, short string,
	get func(settings.DiffTool) string,
	set func(*settings.Config, string),
) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			cfg, err := settings.Load()
			if err != nil {
				return err
			}

			if len(args) == 0 {
				fmt.Println(get(cfg.CurrentDiff()))

				return nil
			}

			set(&cfg, args[0])

			if err := settings.Save(&cfg); err != nil {
				return err
			}

			fmt.Println(textformat.Info(fmt.Sprintf("%s set for %s", strings.Fields(use)[0], runtime.GOOS)))

			return nil
		},
	}
}

func diffProgramCmd() *cobra.Command {
	return diffScalarCmd(
		"diff-program [path]",
		"Show or set the diff program for this operating system",
		func(d settings.DiffTool) string { return d.Program },
		func(c *settings.Config, v string) { c.SetDiffProgram(v) },
	)
}

func diffRoleFilterCmd() *cobra.Command {
	return diffScalarCmd(
		"diff-role-filter [name]",
		"Show or set the role-compare diff filter for this operating system",
		func(d settings.DiffTool) string { return d.RoleFilter },
		func(c *settings.Config, v string) { c.SetRoleDiffFilter(v) },
	)
}

func diffRunbookFilterCmd() *cobra.Command {
	return diffScalarCmd(
		"diff-runbook-filter [name]",
		"Show or set the runbook-compare diff filter for this operating system",
		func(d settings.DiffTool) string { return d.RunbookFilter },
		func(c *settings.Config, v string) { c.SetRunbookDiffFilter(v) },
	)
}

// diffArgsCmd shows or sets the additional diff arguments for the current OS.
// The argument template may contain {left}, {right}, and {filter}; passing
// multiple arguments replaces the whole list.
func diffArgsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "diff-args [arg ...]",
		Short: "Show or set the additional diff arguments for this operating system",
		RunE: func(_ *cobra.Command, args []string) error {
			cfg, err := settings.Load()
			if err != nil {
				return err
			}

			if len(args) == 0 {
				for _, a := range cfg.CurrentDiff().AdditionalArgs {
					fmt.Println(a)
				}

				return nil
			}

			cfg.SetDiffAdditionalArgs(args)

			if err := settings.Save(&cfg); err != nil {
				return err
			}

			fmt.Println(textformat.Info(fmt.Sprintf("diff-args set for %s", runtime.GOOS)))

			return nil
		},
	}
}
