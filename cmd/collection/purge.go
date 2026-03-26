/*
Copyright © 2026 Julian Easterling <julian@julianscorner.com>

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

package collection

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/dcjulian29/ansible-dev/internal/ansible"
	"github.com/dcjulian29/go-toolbox/color"
	"github.com/dcjulian29/go-toolbox/filesystem"
	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
)

// purgeCmd creates the Cobra command for "ansible-dev collection purge",
// which deletes all installed Ansible collection files from the development
// environment.
//
// The command reads the "collections_path" key from the [defaults] section
// of ansible.cfg, appends the standard "ansible_collections" subdirectory,
// and recursively removes that directory tree. A confirmation message is
// printed to stdout on success.
//
// Note: despite the Use string showing "<collection>", this command does
// not accept a positional argument — it unconditionally purges the entire
// ansible_collections directory.
//
// An error is returned if:
//   - ansible.cfg cannot be loaded or is missing the [defaults] section.
//   - The "collections_path" key is not defined.
//   - The ansible_collections directory does not exist.
//   - The directory cannot be removed.
//
// A PreRunE hook calls [ansible.EnsureAnsibleDirectory] to verify the
// current directory is a valid Ansible project.
func purgeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "purge <collection>",
		Short: "Purge all Ansible collection files from the development environment",
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg, err := ini.Load("ansible.cfg")
			if err != nil {
				return err
			}

			section, err := cfg.GetSection("defaults")
			if err != nil {
				return err
			}

			path, err := section.GetKey("collections_path")
			if err != nil {
				return err
			}

			folder := filepath.Join(path.String(), "ansible_collections")

			if !filesystem.DirectoryExists(folder) {
				return errors.New("collections files do not exists")
			}

			if err := filesystem.RemoveDirectory(folder); err != nil {
				return err
			}

			fmt.Println(color.Info("collections files were purged"))

			return nil
		},
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return ansible.EnsureAnsibleDirectory()
		},
	}

	return cmd
}
