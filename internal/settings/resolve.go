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
	"fmt"
	"os"

	"github.com/dcjulian29/go-toolbox/filesystem"
)

// The environment variables read by versions of ansible-dev before the
// configuration file existed. They are no longer consulted for a value; they
// are only checked so that [unsetError] can point someone who still has them
// set at "ansible-dev config import-env".
const (
	LegacyRolesEnv    = "ANSIBLE_ROLES"
	LegacyRunbooksEnv = "ANSIBLE_RUNBOOKS"
)

// unsetError builds the error returned when a required path has no value. When
// the environment variable that used to supply it is still set, the message
// points at "config import-env" instead of the individual setter, so an upgrade
// from an older version is a single command rather than a hunt through the
// documentation.
func unsetError(key, command, env string) error {
	if os.Getenv(env) != "" {
		return fmt.Errorf(
			"%s is not configured, but %s is set; that variable is no longer read "+
				"(run 'ansible-dev config import-env' to import it)", key, env)
	}

	return fmt.Errorf("%s is not configured (run 'ansible-dev config %s <dir>')", key, command)
}

// resolvePath returns the usable value of a required directory setting: a
// leading "~" expanded to the current user's home directory, and the result
// confirmed to be a directory that exists.
//
// Expanding here rather than when the value is stored is what keeps one dotfile
// portable: "~/code/ansible/roles" resolves to the right place on each machine,
// including across operating systems, where an absolute path could not.
//
// The existence check is deliberately an error, unlike the warning the config
// setters emit. A setter may reasonably be pointed at a directory that does not
// exist yet, but by the time a command needs the path it has to be real —
// without this, an unexpanded or mistyped path yields no matches and the
// command finishes silently, which reads as "nothing to do" rather than
// "misconfigured".
func resolvePath(value, key, command, env string) (string, error) {
	if value == "" {
		return "", unsetError(key, command, env)
	}

	expanded := filesystem.ExpandHome(value)

	if !filesystem.DirectoryExist(expanded) {
		where := fmt.Sprintf("'%s'", value)
		if expanded != value {
			where = fmt.Sprintf("'%s' (expanded to '%s')", value, expanded)
		}

		return "", fmt.Errorf(
			"%s is set to %s, which is not an existing directory "+
				"(run 'ansible-dev config %s <dir>' to correct it)", key, where, command)
	}

	return expanded, nil
}

// RolesPath returns the configured roles repository directory, with "~"
// expanded. It is an error when roles_path is unset or does not name an
// existing directory.
func RolesPath() (string, error) {
	cfg, err := Load()
	if err != nil {
		return "", err
	}

	return resolvePath(cfg.RolesPath, "roles_path", "roles-path", LegacyRolesEnv)
}

// RunbooksPath returns the configured runbooks repository directory, with "~"
// expanded. It is an error when runbooks_path is unset or does not name an
// existing directory.
func RunbooksPath() (string, error) {
	cfg, err := Load()
	if err != nil {
		return "", err
	}

	return resolvePath(cfg.RunbooksPath, "runbooks_path", "runbooks-path", LegacyRunbooksEnv)
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
