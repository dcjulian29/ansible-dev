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
	"bytes"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/dcjulian29/ansible-dev/internal/templates"
	"github.com/dcjulian29/go-toolbox/filesystem"
)

// ApplyTemplate walks the template filesystem src and writes each file into
// dest, creating parent directories as needed and applying the given
// !!SENTINEL!! replacements to every file's contents. Existing files are
// overwritten. This reproduces the overlay-then-substitute step the legacy
// PowerShell scaffolding performed with the _AddToRole / _AddToRunbook
// directories.
//
// Plain byte substitution is used deliberately rather than text/template:
// several template files (the GitHub Actions workflows in particular) contain
// literal "{{ }}" expressions that would otherwise collide with Go template
// delimiters.
func ApplyTemplate(src fs.FS, dest string, replacements map[string]string) error {
	return fs.WalkDir(src, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		target := filepath.Join(dest, filepath.FromSlash(path))

		if d.IsDir() {
			return filesystem.EnsureDirectoryExist(target)
		}

		content, err := fs.ReadFile(src, path)
		if err != nil {
			return err
		}

		for sentinel, value := range replacements {
			content = bytes.ReplaceAll(content, []byte(sentinel), []byte(value))
		}

		if err := filesystem.EnsureDirectoryExist(filepath.Dir(target)); err != nil {
			return err
		}

		return os.WriteFile(target, content, 0o644)
	})
}

// ApplyRoleTemplate overlays the embedded role scaffolding (LICENSE, README,
// lint configuration, GitHub workflows, meta/main.yml, ...) onto an existing
// role directory, substituting the role's bare name for !!ROLE_NAME!! and the
// supplied description for !!ROLE_DESC!!.
func ApplyRoleTemplate(dir, role, description string) error {
	src, err := templates.Role()
	if err != nil {
		return err
	}

	return ApplyTemplate(src, dir, map[string]string{
		"!!ROLE_NAME!!": BaseRoleName(role),
		"!!ROLE_DESC!!": description,
	})
}
