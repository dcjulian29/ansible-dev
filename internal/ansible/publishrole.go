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

package ansible

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dcjulian29/go-toolbox/execute"
	"github.com/dcjulian29/go-toolbox/filesystem"
)

// BaseRoleName strips the "namespace." prefix from a fully-qualified role name
// (for example "dcjulian29.nginx" becomes "nginx"), matching the bare-name
// convention used for published repositories and requirements.yml sources. A
// name that contains no dot is returned unchanged.
func BaseRoleName(role string) string {
	if i := strings.IndexByte(role, '.'); i >= 0 {
		return role[i+1:]
	}

	return role
}

// PublishRole copies a freshly-scaffolded role from its workspace location
// into the directory named by the ANSIBLE_ROLES environment variable (using
// the role's bare name), initializes a git repository there, and then creates
// and pushes a public GitHub repository named "ansible-role-<name>".
//
// It relies on the "git" and "gh" executables being installed and, in the case
// of gh, already authenticated. An error is returned if ANSIBLE_ROLES is unset,
// the destination already exists, or any external command fails.
func PublishRole(workspaceDir, role, description string) error {
	base := BaseRoleName(role)

	roles := os.Getenv("ANSIBLE_ROLES")
	if len(roles) == 0 {
		return fmt.Errorf("the 'ANSIBLE_ROLES' environment variable is not defined")
	}

	dest := filepath.Join(strings.ReplaceAll(roles, "\\", string(os.PathSeparator)), base)

	if filesystem.DirectoryExist(dest) {
		return fmt.Errorf("published role already exists at '%s'", dest)
	}

	if err := filesystem.EnsureDirectoryExist(filepath.Dir(dest)); err != nil {
		return err
	}

	if err := os.CopyFS(dest, os.DirFS(workspaceDir)); err != nil {
		return err
	}

	commands := [][]string{
		{"-C", dest, "init"},
		{"-C", dest, "branch", "-M", "main"},
		{"-C", dest, "add", "-A"},
		{"-C", dest, "commit", "-m", "Initial commit"},
	}

	for _, c := range commands {
		if err := execute.ExternalProgram("git", c...); err != nil {
			return err
		}
	}

	return execute.ExternalProgram("gh", "repo", "create",
		"dcjulian29/ansible-role-"+base,
		"--source", dest,
		"--remote", "origin",
		"--push",
		"--public",
		"--disable-wiki",
		"--description", "An Ansible role to "+description)
}
