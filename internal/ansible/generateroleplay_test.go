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
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

const playPath = ".tmp/play.yml"

func TestGenerateRolePlay_WritesAValidSingleRolePlay(t *testing.T) {
	t.Chdir(t.TempDir())

	if err := GenerateRolePlay("dcjulian29.nginx"); err != nil {
		t.Fatalf("GenerateRolePlay() error = %v, want nil", err)
	}

	data, err := os.ReadFile(filepath.FromSlash(playPath))
	if err != nil {
		t.Fatalf("reading the generated play: %v", err)
	}

	// The file is handed to ansible-playbook, so parsing as YAML is the
	// property that matters, not the exact formatting.
	var plays []struct {
		Hosts  string   `yaml:"hosts"`
		Become bool     `yaml:"become"`
		Roles  []string `yaml:"roles"`
	}

	if err := yaml.Unmarshal(data, &plays); err != nil {
		t.Fatalf("generated play is not valid YAML: %v\n%s", err, data)
	}

	if len(plays) != 1 {
		t.Fatalf("play count = %d, want 1", len(plays))
	}

	if plays[0].Hosts != "all" {
		t.Errorf("hosts = %q, want %q", plays[0].Hosts, "all")
	}

	if !plays[0].Become {
		t.Error("become = false, want true")
	}

	if len(plays[0].Roles) != 1 || plays[0].Roles[0] != "dcjulian29.nginx" {
		t.Errorf("roles = %v, want [dcjulian29.nginx]", plays[0].Roles)
	}
}

func TestGenerateRolePlay_CreatesTheTmpDirectory(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	if err := GenerateRolePlay("nginx"); err != nil {
		t.Fatalf("GenerateRolePlay() error = %v, want nil", err)
	}

	info, err := os.Stat(filepath.Join(dir, ".tmp"))
	if err != nil {
		t.Fatalf("expected a .tmp directory: %v", err)
	}

	if !info.IsDir() {
		t.Error(".tmp is not a directory")
	}
}

// A second run must replace the previous play rather than appending to or
// keeping it, or the wrong role would be applied.
func TestGenerateRolePlay_ReplacesAPreviousPlay(t *testing.T) {
	t.Chdir(t.TempDir())

	if err := GenerateRolePlay("first"); err != nil {
		t.Fatal(err)
	}

	if err := GenerateRolePlay("second"); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.FromSlash(playPath))
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(string(data), "first") {
		t.Errorf("play still references the previous role:\n%s", data)
	}

	if !strings.Contains(string(data), "second") {
		t.Errorf("play does not reference the current role:\n%s", data)
	}
}
