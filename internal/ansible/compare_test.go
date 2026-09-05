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
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// writeFile creates dir/name (including parents) with the given contents.
func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()

	path := filepath.Join(dir, filepath.FromSlash(name))

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	return path
}

// captureStdout redirects os.Stdout for the duration of fn and returns what was
// written. ComparePair reports through fmt.Println rather than a returned
// value, so its --checksum output can only be asserted this way.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	original := os.Stdout

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	os.Stdout = w

	done := make(chan string, 1)

	go func() {
		var sb strings.Builder

		buf := make([]byte, 4096)

		for {
			n, err := r.Read(buf)
			if n > 0 {
				sb.Write(buf[:n])
			}

			if err != nil {
				break
			}
		}

		done <- sb.String()
	}()

	fn()

	w.Close() //nolint:errcheck

	os.Stdout = original

	out := <-done

	r.Close() //nolint:errcheck

	return out
}

// runComparePair calls ComparePair with its progress output captured rather
// than left on the suite's stdout, and hands that output back for the tests
// that assert on it.
func runComparePair(
	t *testing.T,
	primary, secondary string,
	ignore []string,
	checksum bool,
	launch func(left, right string) error,
	home string,
) (bool, string, error) {
	t.Helper()

	var (
		differ bool
		err    error
	)

	out := captureStdout(t, func() {
		differ, err = ComparePair(primary, secondary, ignore, checksum, launch, home)
	})

	return differ, out, err
}

// recorder is a launch func that remembers whether and how it was called.
type recorder struct {
	calls int
	left  string
	right string
	err   error
}

func (rec *recorder) launch(left, right string) error {
	rec.calls++
	rec.left = left
	rec.right = right

	return rec.err
}

func TestComparePair_IdenticalTreesDoNotDiffer(t *testing.T) {
	primary, secondary := t.TempDir(), t.TempDir()

	writeFile(t, primary, "tasks/main.yml", "same")
	writeFile(t, secondary, "tasks/main.yml", "same")

	var rec recorder

	differ, _, err := runComparePair(t, primary, secondary, nil, false, rec.launch, "")
	if err != nil {
		t.Fatalf("ComparePair() error = %v, want nil", err)
	}

	if differ {
		t.Error("differ = true, want false for identical trees")
	}

	if rec.calls != 0 {
		t.Errorf("launch called %d times, want 0 when nothing differs", rec.calls)
	}
}

func TestComparePair_DifferingContentLaunchesDiff(t *testing.T) {
	primary, secondary := t.TempDir(), t.TempDir()

	writeFile(t, primary, "tasks/main.yml", "installed copy")
	writeFile(t, secondary, "tasks/main.yml", "canonical source")

	var rec recorder

	differ, _, err := runComparePair(t, primary, secondary, nil, false, rec.launch, "")
	if err != nil {
		t.Fatalf("ComparePair() error = %v, want nil", err)
	}

	if !differ {
		t.Error("differ = false, want true when file contents differ")
	}

	if rec.calls != 1 {
		t.Fatalf("launch called %d times, want 1", rec.calls)
	}

	// The diff opens with the canonical source on the left and the installed
	// copy on the right, which is the reverse of the comparison order.
	if rec.left != secondary || rec.right != primary {
		t.Errorf("launch(%q, %q), want launch(%q, %q)", rec.left, rec.right, secondary, primary)
	}
}

func TestComparePair_FileCountMismatchDiffers(t *testing.T) {
	tests := []struct {
		name  string
		extra string // "primary" or "secondary"
	}{
		{name: "extra file on the primary side", extra: "primary"},
		{name: "extra file on the secondary side", extra: "secondary"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			primary, secondary := t.TempDir(), t.TempDir()

			writeFile(t, primary, "tasks/main.yml", "same")
			writeFile(t, secondary, "tasks/main.yml", "same")

			if tt.extra == "primary" {
				writeFile(t, primary, "handlers/main.yml", "extra")
			} else {
				writeFile(t, secondary, "handlers/main.yml", "extra")
			}

			differ, _, err := runComparePair(t, primary, secondary, nil, false, nil, "")
			if err != nil {
				t.Fatalf("ComparePair() error = %v, want nil", err)
			}

			if !differ {
				t.Error("differ = false, want true when the file counts disagree")
			}
		})
	}
}

// TestComparePair_IgnoreExcludesBothSides pins that an ignore entry is applied
// to the primary and secondary scans alike. Applying it to only one side would
// leave the file counts unequal and report a spurious difference.
func TestComparePair_IgnoreExcludesBothSides(t *testing.T) {
	primary, secondary := t.TempDir(), t.TempDir()

	writeFile(t, primary, "tasks/main.yml", "same")
	writeFile(t, secondary, "tasks/main.yml", "same")
	writeFile(t, primary, "README.md", "installed readme")
	writeFile(t, secondary, "README.md", "source readme")

	var rec recorder

	differ, _, err := runComparePair(t, primary, secondary, []string{"README.md"}, false, rec.launch, "")
	if err != nil {
		t.Fatalf("ComparePair() error = %v, want nil", err)
	}

	if differ {
		t.Error("differ = true, want false when the only difference is ignored")
	}

	if rec.calls != 0 {
		t.Errorf("launch called %d times, want 0", rec.calls)
	}
}

func TestComparePair_IgnoreSkipsWholeDirectory(t *testing.T) {
	primary, secondary := t.TempDir(), t.TempDir()

	writeFile(t, primary, "tasks/main.yml", "same")
	writeFile(t, secondary, "tasks/main.yml", "same")
	writeFile(t, primary, ".git/config", "installed")

	differ, _, err := runComparePair(t, primary, secondary, []string{".git"}, false, nil, "")
	if err != nil {
		t.Fatalf("ComparePair() error = %v, want nil", err)
	}

	if differ {
		t.Error("differ = true, want false when the extra tree is ignored")
	}
}

// TestComparePair_NilLaunchStillReportsDifference covers the --no-diff path:
// the difference must still be detected and reported, just not opened.
func TestComparePair_NilLaunchStillReportsDifference(t *testing.T) {
	primary, secondary := t.TempDir(), t.TempDir()

	writeFile(t, primary, "tasks/main.yml", "installed copy")
	writeFile(t, secondary, "tasks/main.yml", "canonical source")

	differ, _, err := runComparePair(t, primary, secondary, nil, false, nil, "")
	if err != nil {
		t.Fatalf("ComparePair() error = %v, want nil", err)
	}

	if !differ {
		t.Error("differ = false, want true even when no diff tool is launched")
	}
}

func TestComparePair_LaunchErrorIsReturnedWithTheDifference(t *testing.T) {
	primary, secondary := t.TempDir(), t.TempDir()

	writeFile(t, primary, "tasks/main.yml", "installed copy")
	writeFile(t, secondary, "tasks/main.yml", "canonical source")

	rec := recorder{err: errors.New("diff program missing")}

	differ, _, err := runComparePair(t, primary, secondary, nil, false, rec.launch, "")
	if err == nil {
		t.Fatal("ComparePair() error = nil, want the launch error")
	}

	// The difference was real even though opening it failed, so a caller that
	// reports both is not misled about what was found.
	if !differ {
		t.Error("differ = false, want true alongside the launch error")
	}
}

func TestComparePair_ChecksumReportsMatchesAndDifferences(t *testing.T) {
	primary, secondary := t.TempDir(), t.TempDir()

	writeFile(t, primary, "tasks/main.yml", "same")
	writeFile(t, secondary, "tasks/main.yml", "same")
	writeFile(t, primary, "defaults/main.yml", "installed copy")
	writeFile(t, secondary, "defaults/main.yml", "canonical source")

	differ, out, err := runComparePair(t, primary, secondary, nil, true, nil, "")
	if err != nil {
		t.Fatalf("ComparePair() error = %v, want nil", err)
	}

	if !differ {
		t.Error("differ = false, want true")
	}

	// Checksum lines name the file relative to the primary directory. The
	// header above them still carries the full paths, so only the hash lines
	// are examined here.
	for _, line := range strings.Split(out, "\n") {
		if !strings.Contains(line, "==") && !strings.Contains(line, "!=") {
			continue
		}

		if strings.Contains(line, primary) {
			t.Errorf("checksum line leaks the absolute primary path: %s", line)
		}
	}

	matching := filepath.Join("tasks", "main.yml")
	if !strings.Contains(out, matching+":") || !strings.Contains(out, "==") {
		t.Errorf("output missing an equal line for %s:\n%s", matching, out)
	}

	differing := filepath.Join("defaults", "main.yml")
	if !strings.Contains(out, differing+":") || !strings.Contains(out, "!=") {
		t.Errorf("output missing a not-equal line for %s:\n%s", differing, out)
	}
}

// TestComparePair_HomeFolderAbbreviatesHeader covers the printed header, the
// only use of the homeFolder argument.
func TestComparePair_HomeFolderAbbreviatesHeader(t *testing.T) {
	home := t.TempDir()

	primary := filepath.Join(home, "project", "roles", "nginx")
	secondary := filepath.Join(home, "repo", "nginx")

	writeFile(t, primary, "tasks/main.yml", "same")
	writeFile(t, secondary, "tasks/main.yml", "same")

	_, out, err := runComparePair(t, primary, secondary, nil, false, nil, home)
	if err != nil {
		t.Fatalf("ComparePair() error = %v, want nil", err)
	}

	if strings.Contains(out, home) {
		t.Errorf("header was not abbreviated to '~':\n%s", out)
	}

	if !strings.Contains(out, "~") {
		t.Errorf("header does not contain '~':\n%s", out)
	}
}

func TestComparePair_EmptyHomeFolderLeavesHeaderIntact(t *testing.T) {
	primary, secondary := t.TempDir(), t.TempDir()

	writeFile(t, primary, "tasks/main.yml", "same")
	writeFile(t, secondary, "tasks/main.yml", "same")

	_, out, err := runComparePair(t, primary, secondary, nil, false, nil, "")
	if err != nil {
		t.Fatalf("ComparePair() error = %v, want nil", err)
	}

	if !strings.Contains(out, primary) || !strings.Contains(out, secondary) {
		t.Errorf("header should print both full paths when homeFolder is empty:\n%s", out)
	}
}

func TestHomeFolder_UsesThePlatformVariable(t *testing.T) {
	want := filepath.Join(t.TempDir(), "home")

	if runtime.GOOS == "windows" {
		t.Setenv("USERPROFILE", want)
	} else {
		t.Setenv("HOME", want)
	}

	if got := HomeFolder(); got != want {
		t.Errorf("HomeFolder() = %q, want %q", got, want)
	}
}

// TestHomeFolder_NormalizesSeparators pins the backslash conversion, which is
// what lets the result be compared against paths built with the platform
// separator when the environment variable uses the other one.
func TestHomeFolder_NormalizesSeparators(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("backslash is already the separator on Windows")
	}

	t.Setenv("HOME", `/home/julian\code`)

	if got, want := HomeFolder(), "/home/julian/code"; got != want {
		t.Errorf("HomeFolder() = %q, want %q", got, want)
	}
}
