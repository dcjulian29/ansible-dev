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
package initialize

import (
	"errors"
	"fmt"

	"github.com/dcjulian29/go-toolbox/color"
	"github.com/dcjulian29/go-toolbox/filesystem"
	"github.com/spf13/cobra"
)

var force bool

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "initialize",
		Aliases: []string{"init"},
		Short:   "Initialize an Ansible development vagrant environment",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(color.Yellow("Initializing development environment..."))

			if filesystem.FileExists("ansible.cfg") && !force {
				return errors.New("ansible development environment already exist and force was not provided")
			}

			if err := ansible_cfg(); err != nil {
				return err
			}

			if err := ansible_lint(); err != nil {
				return err
			}

			fmt.Println("  ...  collections/")
			if err := filesystem.EnsureDirectoryExist("collections"); err != nil {
				return err
			}

			if err := group_vars(); err != nil {
				return err
			}

			if err := host_vars(); err != nil {
				return err
			}

			if err := inventory_file(); err != nil {
				return err
			}

			if err := runbook(); err != nil {
				return err
			}

			fmt.Println("  ...  roles/")
			if err := filesystem.EnsureDirectoryExist("roles"); err != nil {
				return err
			}

			if err := yaml_ignore(); err != nil {
				return err
			}

			if err := yaml_lint(); err != nil {
				return err
			}

			if err := vagrant_file(); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "overwrite an existing development environment")

	return cmd
}

func ansible_cfg() error {
	fmt.Println("  ...  ansible.cfg")

	content := []byte(`[defaults]
any_errors_fatal            = true
collections_path            = ./collections
duplicate_dict_key          = error
interpreter_python          = auto_silent
inventory                   = ./hosts.ini
log_path                    = ./ansible.log
roles_path                  = ./roles
callback_result_format      = yaml
verbosity                   = 1
`)

	if err := filesystem.EnsureFileExist("ansible.cfg", content); err != nil {
		return err
	}

	return nil
}

func ansible_lint() error {
	fmt.Println("  ...  .ansible-lint")

	content := []byte(`---
enable_list:
  - args
  - empty-string-compare
  - no-log-password
  - no-same-owner
  - name[prefix]
  - yaml
exclude_paths:
  - collections/
  - roles/
kinds:
  - playbook: "playbooks/*.yml"
profile: production
skip_list:
  - experimental
`)

	if err := filesystem.EnsureFileExist(".ansible-lint", content); err != nil {
		return err
	}

	return nil
}

func group_vars() error {
	fmt.Println("  ...  group_vars/")
	if err := filesystem.EnsureDirectoryExist("group_vars"); err != nil {
		return err
	}

	fmt.Println("  ...    vagrant.yml")

	content := []byte("---\nvarname: value")

	if err := filesystem.EnsureFileExist("group_vars/vagrant.yml", content); err != nil {
		return err
	}

	return nil
}

func host_vars() error {
	fmt.Println("  ...  host_vars/")
	if err := filesystem.EnsureDirectoryExist("host_vars"); err != nil {
		return err
	}

	content := []byte("---\nvarname: value")

	fmt.Println("  ...    debian.yml")

	if err := filesystem.EnsureFileExist("host_vars/debian.yml", content); err != nil {
		return err
	}

	fmt.Println("  ...    alma.yml")

	if err := filesystem.EnsureFileExist("host_vars/alma.yml", content); err != nil {
		return err
	}

	return nil
}

func inventory_file() error {
	fmt.Println("  ...  hosts.ini")

	content := []byte(`[vagrant]
debian ansible_host=192.168.57.5
alma ansible_host=192.168.57.6

[all:vars]
ansible_user=vagrant
ansible_ssh_common_args='-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o CheckHostIP=no'
ansible_port=22
ansible_ssh_private_key_file=~/.ssh/insecure_private_key
`)

	if err := filesystem.EnsureFileExist("hosts.ini", content); err != nil {
		return err
	}

	return nil
}

func runbook() error {
	fmt.Println("  ...  playbooks/")
	if err := filesystem.EnsureDirectoryExist("playbooks"); err != nil {
		return err
	}

	fmt.Println("  ...    runbook.yml")

	content := []byte(`---
- name: Test Ansible Runbook
  hosts: all
  become: true
  any_errors_fatal: true
  gather_facts: true

  handlers:
  #  - name: [Name of Handler]
  #    ...

  tasks:
  #  - name: [Name of Task]
  #    ...

  #  - # Repeat as necessary

  vars:
    name: value
    # variables needed for runbook
`)

	if err := filesystem.EnsureFileExist("playbooks/runbook.yml", content); err != nil {
		return err
	}

	return nil
}

func yaml_ignore() error {
	fmt.Println("  ...  .yamlignore")

	content := []byte(`secrets.yml
`)

	if err := filesystem.EnsureFileExist(".yamlignore", content); err != nil {
		return err
	}

	return nil
}

func yaml_lint() error {
	fmt.Println("  ...  .yamlint")

	content := []byte(`---
extends: default

rules:
  braces:
    min-spaces-inside: 0
    max-spaces-inside: 1
  comments:
    min-spaces-from-content: 1
  comments-indentation: false
  indentation:
    indent-sequences: true
  line-length:
    max: 125
  new-lines: disable
  octal-values:
    forbid-implicit-octal: true
    forbid-explicit-octal: true

ignore-from-file:
  - .gitignore
  - .yamlignore
`)

	if err := filesystem.EnsureFileExist(".yamlint", content); err != nil {
		return err
	}

	return nil
}

func vagrant_file() error {
	fmt.Println("  ...  Vagrantfile")

	content := []byte(`Vagrant.configure("2") do |config|
  config.ssh.insert_key = false
  if Vagrant.has_plugin?("vagrant-vbguest")
    config.vbguest.auto_update = false
  end
    config.vm.boot_timeout = 600
    config.vm.box_check_update = true
  config.vm.provision "ping", type: "shell", inline: "ping -c 1 192.168.57.1", run: "always"
  config.vm.synced_folder ".", "/vagrant", disabled: true
  config.vm.provider "virtualbox" do |vb|
    vb.gui = false
    vb.cpus = 4
    vb.memory = 4096
    vb.check_guest_additions = false
    vb.customize [ "modifyvm", :id, "--uartmode1", "disconnected" ]
    vb.customize [ "modifyvm", :id, "--graphicscontroller", "vmsvga"]
    vb.customize [ "modifyvm", :id, "--ioapic", "on"]
  end
  config.vm.define "debian" do |c|
    c.vm.box = "debian/bookworm64"
    c.vm.hostname = "debian.dev"
    c.vm.network "private_network", ip: "192.168.57.5"
  end
  config.vm.define "alma" do |c|
    c.vm.box = "almalinux/10"
    c.vm.hostname = "alma.dev"
    c.vm.network "private_network", ip: "192.168.57.6"
  end
end
`)

	if err := filesystem.EnsureFileExist("Vagrantfile", content); err != nil {
		return err
	}

	return nil
}
