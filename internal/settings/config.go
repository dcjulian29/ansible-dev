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

// Package settings holds the ansible-dev user configuration, persisted to
// ~/.config/ansible-dev.yml via the go-toolbox configuration package. It also
// provides resolver helpers that return the effective value for a setting:
// the roles/runbooks paths (required), the role/runbook ignore lists (empty
// means exclude nothing), and the diff-tool settings for the current OS.
package settings

import (
	"github.com/dcjulian29/go-toolbox/configuration"
)

// Config is the ansible-dev configuration persisted to
// ~/.config/ansible-dev.yml.
//
// Diff holds the external diff-tool settings keyed by operating system
// (runtime.GOOS), so a single dotfile can serve several machines and a new
// platform is just a new map entry. Only the entry for the OS you run on is
// consulted; the compare commands error when the current OS has no diff program
// configured (unless run with --no-diff).
type Config struct {
	// RolesPath is the directory holding the published role repositories
	// (bare names). Required; role commands error when it is unset.
	RolesPath string `yaml:"roles_path"`

	// RunbooksPath is the directory holding the published runbook repositories
	// (bare names). Required; runbook commands error when it is unset.
	RunbooksPath string `yaml:"runbooks_path"`

	// RoleIgnore is the set of path substrings excluded by "role compare".
	// When empty, nothing is excluded (every file is compared).
	RoleIgnore []string `yaml:"role_ignore"`

	// RunbookIgnore is the set of path substrings excluded by "runbook
	// compare". When empty, nothing is excluded (every file is compared).
	RunbookIgnore []string `yaml:"runbook_ignore"`

	// Diff maps an operating system (runtime.GOOS: "windows", "linux",
	// "darwin", ...) to its diff-tool configuration.
	Diff map[string]DiffTool `yaml:"diff"`
}

// file is the process-wide singleton backing ~/.config/ansible-dev.yml.
var file = configuration.New[Config]("ansible-dev.yml")

// Load returns the persisted configuration (singleton-cached on first read).
func Load() (Config, error) {
	return file.Load()
}

// Save persists cfg and refreshes the cached configuration.
func Save(cfg *Config) error {
	return file.Save(cfg)
}

// Show returns the persisted configuration rendered as YAML.
func Show() (string, error) {
	return file.Show()
}

// Path returns the absolute path of the configuration file.
func Path() (string, error) {
	return file.Path()
}
