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

func purgeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "purge <collection>",
		Short: "Purge all Ansible collection files from the development environment",
		RunE: func(cmd *cobra.Command, args []string) error {
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
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ansible.EnsureAnsibleDirectory()
		},
	}

	return cmd
}
