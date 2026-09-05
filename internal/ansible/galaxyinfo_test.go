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
	"testing"
)

func TestReadGalaxyInfo_ReadsNamespaceAndName(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "galaxy.yml", `---
namespace: dcjulian29
name: workstation
version: 1.2.3
readme: README.md
`)

	got, err := ReadGalaxyInfo(dir)
	if err != nil {
		t.Fatalf("ReadGalaxyInfo() error = %v, want nil", err)
	}

	if got.Namespace != "dcjulian29" || got.Name != "workstation" {
		t.Errorf("GalaxyInfo = %+v, want namespace dcjulian29 and name workstation", got)
	}
}

func TestReadGalaxyInfo_ErrorsWhenTheFileIsMissing(t *testing.T) {
	if _, err := ReadGalaxyInfo(t.TempDir()); err == nil {
		t.Error("ReadGalaxyInfo() error = nil, want an error without galaxy.yml")
	}
}

func TestReadGalaxyInfo_ErrorsOnMalformedYAML(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "galaxy.yml", "namespace: [unterminated\n")

	if _, err := ReadGalaxyInfo(dir); err == nil {
		t.Error("ReadGalaxyInfo() error = nil, want a parse error")
	}
}

// A galaxy.yml without the two keys parses cleanly but yields empty values.
// "runbook compare" relies on that to build an install path, so pin it: the
// zero value is what the caller sees, not an error.
func TestReadGalaxyInfo_MissingKeysYieldEmptyValues(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "galaxy.yml", "version: 1.0.0\n")

	got, err := ReadGalaxyInfo(dir)
	if err != nil {
		t.Fatalf("ReadGalaxyInfo() error = %v, want nil", err)
	}

	if got.Namespace != "" || got.Name != "" {
		t.Errorf("GalaxyInfo = %+v, want zero values", got)
	}
}
