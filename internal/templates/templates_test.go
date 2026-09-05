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

package templates

import (
	"io/fs"
	"strings"
	"testing"
)

// entries lists every file path in the embedded tree.
func entries(t *testing.T, tree fs.FS) []string {
	t.Helper()

	var found []string

	err := fs.WalkDir(tree, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			found = append(found, path)
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	return found
}

func TestRole_IsRootedAtTheRoleTree(t *testing.T) {
	tree, err := Role()
	if err != nil {
		t.Fatalf("Role() error = %v, want nil", err)
	}

	files := entries(t, tree)
	if len(files) == 0 {
		t.Fatal("the embedded role tree is empty")
	}

	// Rooting means callers walk "meta/main.yml", never "role/meta/main.yml".
	for _, f := range files {
		if strings.HasPrefix(f, "role/") {
			t.Errorf("path %q is not rooted at the role tree", f)
		}
	}

	if _, err := fs.Stat(tree, "meta/main.yml"); err != nil {
		t.Errorf("expected meta/main.yml in the role tree: %v", err)
	}
}

func TestRunbook_IsRootedAtTheRunbookTree(t *testing.T) {
	tree, err := Runbook()
	if err != nil {
		t.Fatalf("Runbook() error = %v, want nil", err)
	}

	files := entries(t, tree)
	if len(files) == 0 {
		t.Fatal("the embedded runbook tree is empty")
	}

	for _, f := range files {
		if strings.HasPrefix(f, "runbook/") {
			t.Errorf("path %q is not rooted at the runbook tree", f)
		}
	}

	if _, err := fs.Stat(tree, "galaxy.yml"); err != nil {
		t.Errorf("expected galaxy.yml in the runbook tree: %v", err)
	}
}

// TestTrees_IncludeDotfiles is the regression test for the "all:" embed prefix
// documented on this package. Without it, Go silently drops every path
// beginning with "." or "_" — which is most of the scaffolding — and the
// omission is invisible until a generated role turns out to have no lint or CI
// configuration.
func TestTrees_IncludeDotfiles(t *testing.T) {
	trees := map[string]func() (fs.FS, error){"role": Role, "runbook": Runbook}

	for name, open := range trees {
		t.Run(name, func(t *testing.T) {
			tree, err := open()
			if err != nil {
				t.Fatal(err)
			}

			var dotfiles int

			for _, f := range entries(t, tree) {
				for _, segment := range strings.Split(f, "/") {
					if strings.HasPrefix(segment, ".") || strings.HasPrefix(segment, "_") {
						dotfiles++

						break
					}
				}
			}

			if dotfiles == 0 {
				t.Errorf("no dot-prefixed paths in the %s tree; the 'all:' embed prefix is missing", name)
			}
		})
	}
}

// The two trees are separate roots, so neither should contain the other's
// files. A regression in the fs.Sub call would otherwise go unnoticed.
func TestTrees_AreDistinct(t *testing.T) {
	role, err := Role()
	if err != nil {
		t.Fatal(err)
	}

	runbook, err := Runbook()
	if err != nil {
		t.Fatal(err)
	}

	if _, err := fs.Stat(role, "galaxy.yml"); err == nil {
		t.Error("galaxy.yml (a runbook file) is present in the role tree")
	}

	if _, err := fs.Stat(runbook, "meta/main.yml"); err == nil {
		t.Error("meta/main.yml (a role file) is present in the runbook tree")
	}
}
