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

// Package templates embeds the scaffolding files that are overlaid on top of
// a freshly-created Ansible role or runbook when it is published as a public
// repository. The files carry !!SENTINEL!! placeholders (for example
// !!ROLE_NAME!! and !!ROLE_DESC!!) that callers substitute at render time.
//
// The trees are embedded with the "all:" prefix on purpose: without it Go's
// embed directive silently skips any file or directory whose name begins with
// "." or "_", and every scaffolding tree here is dominated by dotfiles
// (.ansible-lint, .github/, .gitignore, .vscode/, .yamllint). "all:" forces
// them to be included.
package templates

import (
	"embed"
	"io/fs"
)

//go:embed all:role all:runbook
var files embed.FS

// Role returns a filesystem rooted at the role scaffolding tree, so that
// callers walk paths like ".ansible-lint" and "meta/main.yml" rather than
// "role/.ansible-lint". An error is returned only if the embedded tree is
// missing, which would indicate a build performed without the template files.
func Role() (fs.FS, error) {
	return fs.Sub(files, "role")
}

// Runbook returns a filesystem rooted at the runbook scaffolding tree. See
// [Role] for details on rooting and error behavior.
func Runbook() (fs.FS, error) {
	return fs.Sub(files, "runbook")
}
