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
	"slices"
	"testing"
)

// Each field of Play must produce its ansible-playbook flag. AskBecomePass
// used to be dropped here, so --ask-become-password never prompted.
func TestPlayArguments_EveryField(t *testing.T) {
	got := playArguments(Play{
		Name:          "acme.nginx",
		Tags:          []string{"install", "configure"},
		AskBecomePass: true,
		AskVaultPass:  true,
		Check:         true,
		FlushCache:    true,
		Step:          true,
		Verbose:       true,
	})

	want := []string{
		"--tags", "install,configure",
		"--flush-cache",
		"--check",
		"--ask-vault-password",
		"--ask-become-pass",
		"-v",
		"--step",
	}

	if !slices.Equal(got, want) {
		t.Errorf("playArguments = %q, want %q", got, want)
	}
}

// A missing playbook used to return nil without running anything. It must
// fail, and fail before the previous run's ansible.log is deleted.
func TestExecutePlay_MissingPlaybook(t *testing.T) {
	t.Chdir(t.TempDir())

	if err := os.WriteFile("ansible.log", []byte("last run"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := ExecutePlay(Play{Name: "acme.nginx"}); err == nil {
		t.Error("ExecutePlay returned nil, want an error for the missing playbook")
	}

	if _, err := os.Stat("ansible.log"); err != nil {
		t.Errorf("ansible.log was removed: %v", err)
	}
}

func TestPlayArguments_NoFields(t *testing.T) {
	if got := playArguments(Play{Name: "acme.nginx"}); len(got) != 0 {
		t.Errorf("playArguments = %q, want no flags", got)
	}
}
