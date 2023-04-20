/*
Copyright © 2023 Julian Easterling <julian@julianscorner.com>

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

var collectionCmd = &cobra.Command{
	Use:     "collection",
	Aliases: []string{"collections"},
	Short:   "Provide management of ansible collections in the development environment",
	Long:    "Provide management of ansible collections in the development environment",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(collectionCmd)
}

func collectionFolder(collection string) (string, error) {
	ensureAnsibleDirectory()

	cfg, err := ini.Load("ansible.cfg")
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	section, err := cfg.GetSection("defaults")
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	path, err := section.GetKey("collections_path")
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	folder := filepath.Join(path.String(), collection)

	return folder, nil
}

func collectionFolderExists(collection string) bool {
	folder, err := collectionFolder(collection)

	if err != nil {
		return false
	}

	return dirExists(folder)
}

func remove_collection(collection string) {
	if !roleFolderExists(collection) {
		fmt.Println(Warn("WARN: Collection '%s' folder not present.", collection))
		return
	}

	folder, err := roleFolder(collection)
	if err != nil {
		fmt.Println(err)
		return
	}

	files, err := filepath.Glob(filepath.Join(folder, "*"))
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(files) == 0 {
		fmt.Println(Warn("WARN: Collection '%s' files not present.", collection))
		return
	}

	removeDir(folder)

	fmt.Println(Info("Collection '%s' files were deleted.", collection))
}
