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

package settings

import (
	"fmt"
	"runtime"
	"strings"
)

// DiffTool is the diff-tool configuration for one operating system: the program
// to run, the filter name for each compare type (empty when the tool has no
// filter concept), and the additional argument template. The template may
// contain the placeholders {left}, {right}, and {filter}, which are substituted
// per comparison by [DiffTool.Command].
type DiffTool struct {
	Program        string   `yaml:"program"`
	RoleFilter     string   `yaml:"role_filter"`
	RunbookFilter  string   `yaml:"runbook_filter"`
	AdditionalArgs []string `yaml:"additional_args"`
}

// Diff returns the diff-tool configuration for the current operating system
// (Config.Diff keyed by runtime.GOOS). It is an error if the current OS has no
// diff program configured. Callers that do not need a visual diff (--no-diff)
// should not call Diff.
func Diff() (DiffTool, error) {
	cfg, err := Load()
	if err != nil {
		return DiffTool{}, err
	}

	diff := cfg.Diff[runtime.GOOS]
	if diff.Program == "" {
		return DiffTool{}, fmt.Errorf(
			"no diff program configured for the current operating system (%s); "+
				"set it with 'ansible-dev config diff-program <path>'",
			runtime.GOOS)
	}

	return diff, nil
}

// CurrentDiff returns the diff configuration for the current operating system
// without requiring that a program be set (used by the config command to
// display current values). The result is a zero DiffTool when the OS has no
// entry.
func (c Config) CurrentDiff() DiffTool {
	return c.Diff[runtime.GOOS]
}

// setDiff mutates (creating if needed) the diff entry for the current operating
// system. Because a new OS is just a new map key, no source change is required
// to configure a platform the tool has not seen before.
func (c *Config) setDiff(mutate func(*DiffTool)) {
	if c.Diff == nil {
		c.Diff = map[string]DiffTool{}
	}

	diff := c.Diff[runtime.GOOS]
	mutate(&diff)
	c.Diff[runtime.GOOS] = diff
}

// SetDiffProgram sets the diff program for the current operating system.
func (c *Config) SetDiffProgram(program string) {
	c.setDiff(func(d *DiffTool) { d.Program = program })
}

// SetRoleDiffFilter sets the role-compare diff filter for the current OS.
func (c *Config) SetRoleDiffFilter(filter string) {
	c.setDiff(func(d *DiffTool) { d.RoleFilter = filter })
}

// SetRunbookDiffFilter sets the runbook-compare diff filter for the current OS.
func (c *Config) SetRunbookDiffFilter(filter string) {
	c.setDiff(func(d *DiffTool) { d.RunbookFilter = filter })
}

// SetDiffAdditionalArgs sets the additional diff arguments for the current OS.
func (c *Config) SetDiffAdditionalArgs(args []string) {
	c.setDiff(func(d *DiffTool) { d.AdditionalArgs = args })
}

// Command builds the program and arguments to launch a diff between left and
// right, expanding the argument template. {left} and {right} become the two
// paths and {filter} becomes filter. When filter is empty, a standalone
// "{filter}" argument is dropped along with the argument immediately before it
// (its flag), so no dangling flag remains; a token that merely contains
// "{filter}" is dropped whole.
func (d DiffTool) Command(filter, left, right string) (string, []string) {
	args := make([]string, 0, len(d.AdditionalArgs))

	for _, a := range d.AdditionalArgs {
		if filter == "" && strings.Contains(a, "{filter}") {
			if a == "{filter}" && len(args) > 0 {
				args = args[:len(args)-1]
			}

			continue
		}

		a = strings.ReplaceAll(a, "{filter}", filter)
		a = strings.ReplaceAll(a, "{left}", left)
		a = strings.ReplaceAll(a, "{right}", right)
		args = append(args, a)
	}

	return d.Program, args
}
