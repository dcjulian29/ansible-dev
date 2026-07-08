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

package settings

import (
	"errors"
)

// RolesPath returns the configured roles repository directory. It is an error
// when roles_path is not set.
func RolesPath() (string, error) {
	cfg, err := Load()
	if err != nil {
		return "", err
	}

	if cfg.RolesPath == "" {
		return "", errors.New("roles_path is not configured (run 'ansible-dev config roles-path <dir>')")
	}

	return cfg.RolesPath, nil
}

// RunbooksPath returns the configured runbooks repository directory. It is an
// error when runbooks_path is not set.
func RunbooksPath() (string, error) {
	cfg, err := Load()
	if err != nil {
		return "", err
	}

	if cfg.RunbooksPath == "" {
		return "", errors.New("runbooks_path is not configured (run 'ansible-dev config runbooks-path <dir>')")
	}

	return cfg.RunbooksPath, nil
}

// RoleIgnore returns the path substrings excluded by "role compare". An empty
// result means nothing is excluded.
func RoleIgnore() ([]string, error) {
	cfg, err := Load()
	if err != nil {
		return nil, err
	}

	return cfg.RoleIgnore, nil
}

// RunbookIgnore returns the path substrings excluded by "runbook compare". An
// empty result means nothing is excluded.
func RunbookIgnore() ([]string, error) {
	cfg, err := Load()
	if err != nil {
		return nil, err
	}

	return cfg.RunbookIgnore, nil
}
