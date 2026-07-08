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
	"github.com/dcjulian29/go-toolbox/execute"
)

// PublishRunbook initializes a git repository in the already-rendered runbook
// directory dir, then creates and pushes a public GitHub repository named
// "ansible-runbook-<name>".
//
// It relies on the "git" and "gh" executables being installed and, in the case
// of gh, already authenticated. An error is returned if any external command
// fails.
func PublishRunbook(dir, name, description string) error {
	commands := [][]string{
		{"-C", dir, "init"},
		{"-C", dir, "branch", "-M", "main"},
		{"-C", dir, "add", "-A"},
		{"-C", dir, "commit", "-m", "Initial commit"},
	}

	for _, c := range commands {
		if err := execute.ExternalProgram("git", c...); err != nil {
			return err
		}
	}

	return execute.ExternalProgram("gh", "repo", "create",
		"dcjulian29/ansible-runbook-"+name,
		"--source", dir,
		"--remote", "origin",
		"--push",
		"--public",
		"--disable-wiki",
		"--description", "An Ansible runbook that will "+description)
}
