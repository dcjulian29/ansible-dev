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
package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
)

var collectionPurgeCmd = &cobra.Command{
	Use:   "purge <collection>",
	Short: "Purge all Ansible collection files from the development environment",
	Long:  "Purge all Ansible collection files from the development environment",
	Run: func(cmd *cobra.Command, args []string) {
		purge_collections()
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		ensureAnsibleDirectory()
	},
}

func init() {
	collectionCmd.AddCommand(collectionPurgeCmd)
}

func purge_collections() {
	cfg, err := ini.Load("ansible.cfg")
	cobra.CheckErr(err)

	section, err := cfg.GetSection("defaults")
	cobra.CheckErr(err)

	path, err := section.GetKey("collections_path")
	cobra.CheckErr(err)

	folder := filepath.Join(path.String(), "ansible_collections")

	if !dirExists(folder) {
		return
	}

	removeDir(folder)

	fmt.Println(Info("collections files were purged"))
}
