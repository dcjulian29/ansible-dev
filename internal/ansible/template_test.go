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
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/dcjulian29/ansible-dev/internal/templates"
)

func read(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	return string(data)
}

func TestApplyTemplate_SubstitutesAndPreservesLayout(t *testing.T) {
	src := fstest.MapFS{
		"README.md":      {Data: []byte("# !!NAME!!\n\n!!DESC!!\n")},
		"meta/main.yml":  {Data: []byte("role_name: !!NAME!!\n")},
		".gitignore":     {Data: []byte("*.retry\n")},
		"deep/a/b/c.txt": {Data: []byte("!!NAME!! nested\n")},
	}

	dest := t.TempDir()

	err := ApplyTemplate(src, dest, map[string]string{
		"!!NAME!!": "nginx",
		"!!DESC!!": "Ansible role to install nginx",
	})
	if err != nil {
		t.Fatalf("ApplyTemplate() error = %v, want nil", err)
	}

	if got := read(t, filepath.Join(dest, "README.md")); got != "# nginx\n\nAnsible role to install nginx\n" {
		t.Errorf("README.md = %q", got)
	}

	if got := read(t, filepath.Join(dest, "meta", "main.yml")); got != "role_name: nginx\n" {
		t.Errorf("meta/main.yml = %q", got)
	}

	// A file with no sentinels is copied through untouched.
	if got := read(t, filepath.Join(dest, ".gitignore")); got != "*.retry\n" {
		t.Errorf(".gitignore = %q", got)
	}

	if got := read(t, filepath.Join(dest, "deep", "a", "b", "c.txt")); got != "nginx nested\n" {
		t.Errorf("deep/a/b/c.txt = %q", got)
	}
}

// TestApplyTemplate_OverwritesExistingFiles pins the overlay semantics: the
// template is applied on top of a directory ansible-galaxy has already
// populated, so a colliding file must be replaced, not skipped.
func TestApplyTemplate_OverwritesExistingFiles(t *testing.T) {
	dest := t.TempDir()

	if err := os.WriteFile(filepath.Join(dest, "README.md"), []byte("from ansible-galaxy"), 0o644); err != nil {
		t.Fatal(err)
	}

	src := fstest.MapFS{"README.md": {Data: []byte("from the template\n")}}

	if err := ApplyTemplate(src, dest, nil); err != nil {
		t.Fatalf("ApplyTemplate() error = %v, want nil", err)
	}

	if got := read(t, filepath.Join(dest, "README.md")); got != "from the template\n" {
		t.Errorf("README.md = %q, want the template's content", got)
	}
}

// Byte substitution rather than text/template is a deliberate choice, because
// the embedded GitHub workflows contain literal "{{ }}" expressions. This pins
// that they survive untouched.
func TestApplyTemplate_LeavesGoTemplateDelimitersAlone(t *testing.T) {
	src := fstest.MapFS{
		"workflow.yml": {Data: []byte("run: echo ${{ github.ref }} !!NAME!!\n")},
	}

	dest := t.TempDir()

	if err := ApplyTemplate(src, dest, map[string]string{"!!NAME!!": "nginx"}); err != nil {
		t.Fatalf("ApplyTemplate() error = %v, want nil", err)
	}

	want := "run: echo ${{ github.ref }} nginx\n"
	if got := read(t, filepath.Join(dest, "workflow.yml")); got != want {
		t.Errorf("workflow.yml = %q, want %q", got, want)
	}
}

func TestApplyTemplate_UnknownSentinelsAreLeftInPlace(t *testing.T) {
	src := fstest.MapFS{"a.txt": {Data: []byte("!!NAME!! !!UNKNOWN!!\n")}}

	dest := t.TempDir()

	if err := ApplyTemplate(src, dest, map[string]string{"!!NAME!!": "nginx"}); err != nil {
		t.Fatalf("ApplyTemplate() error = %v, want nil", err)
	}

	if got := read(t, filepath.Join(dest, "a.txt")); got != "nginx !!UNKNOWN!!\n" {
		t.Errorf("a.txt = %q", got)
	}
}

// sentinel matches any !!PLACEHOLDER!! token left in a rendered file.
var sentinel = regexp.MustCompile(`!![A-Z_]+!!`)

// remainingSentinels renders the real embedded template and reports every
// placeholder the caller failed to substitute. This is what catches a new
// sentinel being added to a template file without a matching entry in the
// substitution map.
func remainingSentinels(t *testing.T, dir string) []string {
	t.Helper()

	var found []string

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		for _, m := range sentinel.FindAllString(read(t, path), -1) {
			rel, _ := filepath.Rel(dir, path)
			found = append(found, rel+": "+m)
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	return found
}

// TestRemainingSentinels_DetectsThem guards the guard: rendering the real role
// template with no replacements at all must report placeholders. Without this,
// a detector that silently found nothing would let
// TestApplyRoleTemplate_LeavesNoUnsubstitutedSentinels pass vacuously.
func TestRemainingSentinels_DetectsThem(t *testing.T) {
	src, err := templates.Role()
	if err != nil {
		t.Fatal(err)
	}

	dest := t.TempDir()

	if err := ApplyTemplate(src, dest, nil); err != nil {
		t.Fatal(err)
	}

	if left := remainingSentinels(t, dest); len(left) == 0 {
		t.Error("no placeholders found in an unsubstituted template; the detector is broken")
	}
}

func TestApplyRoleTemplate_LeavesNoUnsubstitutedSentinels(t *testing.T) {
	dest := t.TempDir()

	err := ApplyRoleTemplate(dest, "acme", "acme.nginx", "Ansible role to install nginx")
	if err != nil {
		t.Fatalf("ApplyRoleTemplate() error = %v, want nil", err)
	}

	if left := remainingSentinels(t, dest); len(left) > 0 {
		t.Errorf("unsubstituted placeholders remain:\n  %s", strings.Join(left, "\n  "))
	}
}

func TestApplyRoleTemplate_FillsNamespaceNameAndDescription(t *testing.T) {
	dest := t.TempDir()

	err := ApplyRoleTemplate(dest, "acme", "acme.nginx", "Ansible role to install nginx")
	if err != nil {
		t.Fatalf("ApplyRoleTemplate() error = %v, want nil", err)
	}

	meta := read(t, filepath.Join(dest, "meta", "main.yml"))

	for _, want := range []string{
		"namespace: acme",
		"role_name: nginx", // the bare name, with the namespace stripped
		"description: Ansible role to install nginx",
	} {
		if !strings.Contains(meta, want) {
			t.Errorf("meta/main.yml missing %q:\n%s", want, meta)
		}
	}

	readme := read(t, filepath.Join(dest, "README.md"))

	for _, want := range []string{
		"name: acme.nginx",
		"https://github.com/acme/ansible-role-nginx.git",
	} {
		if !strings.Contains(readme, want) {
			t.Errorf("README.md missing %q", want)
		}
	}

	if strings.Contains(readme, "dcjulian29") {
		t.Errorf("README.md still carries the old hard-coded namespace:\n%s", readme)
	}
}

// The role template is dominated by dotfiles, which Go's embed directive skips
// without the "all:" prefix. A regression there would silently produce roles
// missing their lint and CI configuration, so assert a few are present.
func TestApplyRoleTemplate_IncludesDotfiles(t *testing.T) {
	dest := t.TempDir()

	if err := ApplyRoleTemplate(dest, "acme", "nginx", "x"); err != nil {
		t.Fatalf("ApplyRoleTemplate() error = %v, want nil", err)
	}

	for _, name := range []string{".ansible-lint", ".gitignore"} {
		if _, err := os.Stat(filepath.Join(dest, name)); err != nil {
			t.Errorf("expected %s in the rendered role: %v", name, err)
		}
	}

	if _, err := os.Stat(filepath.Join(dest, ".github")); err != nil {
		t.Errorf("expected a .github directory in the rendered role: %v", err)
	}
}
