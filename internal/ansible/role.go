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
package ansible

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dcjulian29/go-toolbox/color"
	"github.com/dcjulian29/go-toolbox/execute"
	"github.com/dcjulian29/go-toolbox/filesystem"
	"gopkg.in/ini.v1"
)

func ApplyRoles(roles []string, tags []string, verbose bool) error {
	play := Play{
		Tags:       tags,
		FlushCache: true,
		Verbose:    verbose,
	}

	if len(roles) > 0 {
		for _, role := range roles {
			fmt.Println(color.Info(fmt.Sprintf("\nApplying the '%s' role...", role)))

			err := GenerateRolePlay(role)
			if err != nil {
				return err
			}

			play.Name = role

			err = ExecutePlay(play)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func NewRole(role string, verbose bool) error {
	path, err := RoleFolder(role)
	if err != nil {
		return err
	}

	param := []string{
		"init",
		role,
		"--init-path",
		strings.ReplaceAll(filepath.Dir(path), "\\", "/"),
	}

	if verbose {
		param = append(param, "--verbose")
	}

	if err := execute.ExternalProgram("ansible-galaxy", param...); err != nil {
		return err
	}

	return nil
}

func RemoveRole(role string) error {
	exists := RoleFolderExists(role)

	if !exists {
		return fmt.Errorf("role '%s' folder not present", role)
	}

	folder, err := RoleFolder(role)
	if err != nil {
		return err
	}

	files, err := filepath.Glob(filepath.Join(folder, "*"))
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return fmt.Errorf("role '%s' files not present", role)
	}

	err = filesystem.RemoveDirectory(folder)
	if err != nil {
		return err
	}

	fmt.Println(color.Info(fmt.Sprintf("role '%s' files were deleted.", role)))

	return nil
}

func RootRoleFolder() (string, error) {
	cfg, err := ini.Load("ansible.cfg")
	if err != nil {
		return "", err
	}

	section, err := cfg.GetSection("defaults")
	if err != nil {
		return "", err
	}

	path, err := section.GetKey("roles_path")
	if err != nil {
		return "", err
	}

	return path.String(), nil
}

func RoleFolder(role string) (string, error) {
	folder, err := RootRoleFolder()
	if err != nil {
		return "", err
	}

	return filepath.Join(folder, role), nil
}

func RoleFolderExists(role string) bool {
	folder, err := RoleFolder(role)
	if err != nil {
		return false
	}

	return filesystem.DirectoryExists(folder)
}
