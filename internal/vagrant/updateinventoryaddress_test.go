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

package vagrant

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// inventoryFile writes content to a hosts.ini in a fresh temp directory and
// returns its path.
func inventoryFile(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "hosts.ini")

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	return path
}

func lines(t *testing.T, path string) []string {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	return strings.Split(strings.TrimRight(string(data), "\n"), "\n")
}

func TestUpdateInventoryAddress_RewritesOnlyTheNamedHost(t *testing.T) {
	path := inventoryFile(t, `[vagrant]
debian ansible_host=0.0.0.0
alma   ansible_host=0.0.0.0
`)

	if err := UpdateInventoryAddress(path, "debian", "192.168.1.10"); err != nil {
		t.Fatalf("UpdateInventoryAddress() error = %v, want nil", err)
	}

	got := lines(t, path)

	if got[1] != "debian ansible_host=192.168.1.10" {
		t.Errorf("debian line = %q, want the new address", got[1])
	}

	// The other machine is discovered separately and must not be disturbed.
	if !strings.Contains(got[2], "alma") || !strings.Contains(got[2], "0.0.0.0") {
		t.Errorf("alma line = %q, want it unchanged", got[2])
	}
}

// TestUpdateInventoryAddress_PreservesOtherFields pins that only the
// ansible_host token is replaced; per-host settings on the same line survive.
func TestUpdateInventoryAddress_PreservesOtherFields(t *testing.T) {
	path := inventoryFile(t,
		"[vagrant]\ndebian ansible_host=0.0.0.0 ansible_port=2222 ansible_user=vagrant\n")

	if err := UpdateInventoryAddress(path, "debian", "10.0.0.5"); err != nil {
		t.Fatalf("UpdateInventoryAddress() error = %v, want nil", err)
	}

	want := "debian ansible_host=10.0.0.5 ansible_port=2222 ansible_user=vagrant"
	if got := lines(t, path)[1]; got != want {
		t.Errorf("line = %q, want %q", got, want)
	}
}

func TestUpdateInventoryAddress_KeepsSectionsAndVariables(t *testing.T) {
	path := inventoryFile(t, `[vagrant]
debian ansible_host=0.0.0.0

[all:vars]
ansible_user=vagrant
ansible_port=22
`)

	if err := UpdateInventoryAddress(path, "debian", "172.16.0.2"); err != nil {
		t.Fatalf("UpdateInventoryAddress() error = %v, want nil", err)
	}

	got := strings.Join(lines(t, path), "\n")

	for _, want := range []string{"[vagrant]", "[all:vars]", "ansible_user=vagrant", "ansible_port=22"} {
		if !strings.Contains(got, want) {
			t.Errorf("rewritten file lost %q:\n%s", want, got)
		}
	}
}

// An unknown machine name is not an error: start iterates whatever vagrant
// reports, which need not match every host in the file.
func TestUpdateInventoryAddress_UnknownHostLeavesTheFileUnchanged(t *testing.T) {
	const content = "[vagrant]\ndebian ansible_host=0.0.0.0\n"

	path := inventoryFile(t, content)

	if err := UpdateInventoryAddress(path, "absent", "10.0.0.9"); err != nil {
		t.Fatalf("UpdateInventoryAddress() error = %v, want nil", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != content {
		t.Errorf("file = %q, want it unchanged", string(data))
	}
}

// A host line without an ansible_host token gains nothing: the function only
// substitutes an existing value, it does not add one.
func TestUpdateInventoryAddress_LineWithoutAnAddressIsLeftAlone(t *testing.T) {
	path := inventoryFile(t, "[vagrant]\ndebian ansible_user=vagrant\n")

	if err := UpdateInventoryAddress(path, "debian", "10.0.0.9"); err != nil {
		t.Fatalf("UpdateInventoryAddress() error = %v, want nil", err)
	}

	if got := lines(t, path)[1]; got != "debian ansible_user=vagrant" {
		t.Errorf("line = %q, want it unchanged", got)
	}
}

func TestUpdateInventoryAddress_ErrorsWhenTheFileIsMissing(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "absent.ini")

	if err := UpdateInventoryAddress(missing, "debian", "10.0.0.1"); err == nil {
		t.Error("UpdateInventoryAddress() error = nil, want an error for a missing file")
	}
}

func TestEnsureVagrantfile(t *testing.T) {
	t.Chdir(t.TempDir())

	err := EnsureVagrantfile()
	if err == nil {
		t.Fatal("EnsureVagrantfile() = nil, want an error when the file is missing")
	}

	if !strings.Contains(err.Error(), "Vagrantfile") {
		t.Errorf("error = %q, want it to name the Vagrantfile", err)
	}

	if err := os.WriteFile("Vagrantfile", []byte("Vagrant.configure(\"2\")\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := EnsureVagrantfile(); err != nil {
		t.Errorf("EnsureVagrantfile() = %v, want nil once the file exists", err)
	}
}
