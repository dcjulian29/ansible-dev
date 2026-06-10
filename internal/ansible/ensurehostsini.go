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
	"github.com/dcjulian29/go-toolbox/filesystem"
)

// EnsureHostsIni verifies that a hosts.ini file exists in the current
// working directory. If the file is missing it is recreated from the
// standard template with placeholder addresses (0.0.0.0) so that
// ansible-dev start can proceed and overwrite the addresses once each
// VM has booted and reported its IP via vagrant ssh-config.
//
// A non-nil error is returned only if the file is missing and cannot
// be recreated.
func EnsureHostsIni() error {
	if filesystem.FileExists("hosts.ini") {
		return nil
	}

	content := []byte(`[vagrant]
debian ansible_host=0.0.0.0
alma   ansible_host=0.0.0.0

[all:vars]
ansible_user=vagrant
ansible_ssh_common_args='-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o CheckHostIP=no'
ansible_port=22
ansible_ssh_private_key_file=~/.ssh/insecure_private_key
`)

	return filesystem.EnsureFileExist("hosts.ini", content)
}
