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
	"strings"
	"testing"
)

// hostsProject makes an empty directory current and, when content is non-empty,
// writes it as hosts.ini. Both readers here resolve that name relative to the
// working directory.
func hostsProject(t *testing.T, content string) {
	t.Helper()

	t.Chdir(t.TempDir())

	if content != "" {
		if err := os.WriteFile("hosts.ini", []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestEnsureHostsIni_CreatesAUsableDefault(t *testing.T) {
	hostsProject(t, "")

	if err := EnsureHostsIni(); err != nil {
		t.Fatalf("EnsureHostsIni() error = %v, want nil", err)
	}

	// The generated file must be readable by GetInventory, since the two are
	// used together: start writes placeholder addresses, then overwrites them.
	inventory, err := GetInventory()
	if err != nil {
		t.Fatalf("GetInventory() on the generated file: %v", err)
	}

	if len(inventory) != 2 {
		t.Fatalf("inventory = %+v, want two placeholder hosts", inventory)
	}

	for _, i := range inventory {
		if i.Address != "0.0.0.0" {
			t.Errorf("host %q address = %q, want the 0.0.0.0 placeholder", i.Name, i.Address)
		}
	}
}

// TestEnsureHostsIni_LeavesAnExistingFileAlone matters because the file holds
// addresses discovered from running VMs; regenerating it would discard them.
func TestEnsureHostsIni_LeavesAnExistingFileAlone(t *testing.T) {
	const existing = "[vagrant]\ndebian ansible_host=192.168.1.10\n"

	hostsProject(t, existing)

	if err := EnsureHostsIni(); err != nil {
		t.Fatalf("EnsureHostsIni() error = %v, want nil", err)
	}

	data, err := os.ReadFile("hosts.ini")
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != existing {
		t.Errorf("hosts.ini = %q, want it untouched", string(data))
	}
}

func TestGetInventory_ParsesTheVagrantSection(t *testing.T) {
	hostsProject(t, `[vagrant]
debian ansible_host=192.168.1.10
alma   ansible_host=192.168.1.11

[all:vars]
ansible_user=vagrant
`)

	got, err := GetInventory()
	if err != nil {
		t.Fatalf("GetInventory() error = %v, want nil", err)
	}

	if len(got) != 2 {
		t.Fatalf("inventory = %+v, want two hosts from the [vagrant] section only", got)
	}

	want := map[string]string{"debian": "192.168.1.10", "alma": "192.168.1.11"}

	for _, i := range got {
		if want[i.Name] != i.Address {
			t.Errorf("host %q = %q, want %q", i.Name, i.Address, want[i.Name])
		}
	}
}

func TestGetInventory_ErrorsWithoutTheVagrantSection(t *testing.T) {
	hostsProject(t, "[production]\nweb ansible_host=10.0.0.1\n")

	_, err := GetInventory()
	if err == nil {
		t.Fatal("GetInventory() error = nil, want an error")
	}

	if !strings.Contains(err.Error(), "vagrant") {
		t.Errorf("error = %q, want it to name the missing section", err)
	}
}

func TestGetInventory_ErrorsWhenTheFileIsMissing(t *testing.T) {
	hostsProject(t, "")

	if _, err := GetInventory(); err == nil {
		t.Error("GetInventory() error = nil, want an error when hosts.ini is absent")
	}
}
