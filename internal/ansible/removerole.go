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

	"github.com/dcjulian29/go-toolbox/color"
	"github.com/dcjulian29/go-toolbox/filesystem"
)

// RemoveRole deletes the directory tree for the named role. An error is
// returned if the role folder does not exist or cannot be removed.
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
