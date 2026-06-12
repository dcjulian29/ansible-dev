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

// Package initialize implements the "ansible-dev initialize" (aliased as
// "init") command, which scaffolds a complete Ansible development
// environment in the current directory. The generated layout includes
// configuration files, linter settings, inventory, variable directories,
// a sample runbook playbook, and a multi-VM Vagrantfile.
package initialize

import (
	"errors"
	"fmt"

	"github.com/dcjulian29/ansible-dev/internal/ansible"
	"github.com/dcjulian29/go-toolbox/color"
	"github.com/dcjulian29/go-toolbox/filesystem"
	"github.com/spf13/cobra"
)

var force bool

// NewCommand creates and returns the Cobra command for
// "ansible-dev initialize". The command is also aliased as "init".
//
// When executed it scaffolds the following project structure in the current
// working directory:
//
//	ansible.cfg            – Ansible configuration (roles path, inventory, logging, etc.)
//	.ansible-lint          – ansible-lint profile and rule configuration
//	collections/           – empty directory for Galaxy collections
//	group_vars/vagrant.yml – group variables for the [vagrant] inventory group
//	host_vars/debian.yml   – host variables for the "debian" VM
//	host_vars/alma.yml     – host variables for the "alma" VM
//	hosts.ini              – INI inventory with [vagrant] and [all:vars] sections
//	playbooks/runbook.yml  – skeleton runbook playbook
//	roles/                 – empty directory for Ansible roles
//	.yamlignore            – files excluded from YAML linting (secrets.yml)
//	.yamlint               – yamllint configuration
//	Vagrantfile            – multi-VM Vagrant config (Debian + AlmaLinux, VirtualBox)
//
// If ansible.cfg already exists and --force is not set, the command returns
// an error without modifying any files. When --force is set, every file is
// re-created (or left as-is by the underlying EnsureFileExist helper if
// unchanged).
//
// Flags:
//   - --force, -f: overwrite an existing development environment
//     (default false).
//
// Execution stops and returns at the first file or directory creation error.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "initialize",
		Aliases: []string{"init"},
		Short:   "Initialize an Ansible development vagrant environment",
		RunE: func(_ *cobra.Command, _ []string) error {
			fmt.Println(color.Yellow("Initializing development environment..."))

			if filesystem.FileExists("ansible.cfg") && !force {
				return errors.New("ansible development environment already exist and force was not provided")
			}

			if err := ansibleConfig(); err != nil {
				return err
			}

			if err := ansibleLint(); err != nil {
				return err
			}

			fmt.Println("  ...  collections/")
			if err := filesystem.EnsureDirectoryExist("collections"); err != nil {
				return err
			}

			if err := gitIgnore(); err != nil {
				return err
			}

			if err := groupVariables(); err != nil {
				return err
			}

			if err := hostVariables(); err != nil {
				return err
			}

			if err := inventoryFile(); err != nil {
				return err
			}

			if err := runbook(); err != nil {
				return err
			}

			fmt.Println("  ...  roles/")
			if err := filesystem.EnsureDirectoryExist("roles"); err != nil {
				return err
			}

			if err := yamlIgnore(); err != nil {
				return err
			}

			if err := yamlLint(); err != nil {
				return err
			}

			if err := vagrantFile(); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "overwrite an existing development environment")

	return cmd
}

// ansibleConfig creates the "ansible.cfg" file with default Ansible
// configuration values. Notable settings include:
//   - roles_path:       ./roles
//   - collections_path: ./collections
//   - inventory:        ./hosts.ini
//   - log_path:         ./ansible.log
//   - verbosity:        1
//
// An error is returned if the file cannot be created or written.
func ansibleConfig() error {
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

// ansibleLint creates the ".ansible-lint" configuration file with a
// production-grade profile. The configuration enables specific rules
// (args, empty-string-compare, no-log-password, no-same-owner,
// name[prefix], yaml), excludes the collections/ and roles/ directories
// from linting, and skips experimental rules.
//
// An error is returned if the file cannot be created or written.
func ansibleLint() error {
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

// gitIgnore creates the ".gitignore" file for the development environment.
// It excludes runtime-modified and generated files that should not be
// committed: hosts.ini (rewritten on every vagrant up by ansible-dev),
// ansible.log, .vagrant/, .tmp/, and other build-time Ansible directories.
func gitIgnore() error {
	fmt.Println(" ... .gitignore")

	content := []byte(`.vagrant/
.ansible/
.tmp/
.vagrant/
ansible.log
collections/
hosts.ini
roles/
`)

	if err := filesystem.EnsureFileExist(".gitignore", content); err != nil {
		return err
	}

	return nil
}

// groupVariables creates the "group_vars/" directory and a starter
// "group_vars/vagrant.yml" file containing a placeholder variable. The
// file applies to all hosts in the [vagrant] inventory group.
//
// An error is returned if the directory or file cannot be created.
func groupVariables() error {
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

// hostVariables creates the "host_vars/" directory and starter variable
// files for each VM defined in the default Vagrantfile:
//   - host_vars/debian.yml – variables for the "debian" VM.
//   - host_vars/alma.yml   – variables for the "alma" VM.
//
// Each file contains a single placeholder variable. An error is returned
// if the directory or any file cannot be created.
func hostVariables() error {
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

// inventoryFile creates the "hosts.ini" Ansible inventory file. The
// generated inventory defines:
//   - A [vagrant] group with two hosts: "debian" (192.168.57.5) and
//     "alma" (192.168.57.6).
//   - An [all:vars] section with SSH connection parameters configured for
//     Vagrant (insecure private key, disabled host-key checking, port 22).
//
// An error is returned if the file cannot be created or written.
func inventoryFile() error {
	fmt.Println("  ...  hosts.ini")

	return ansible.EnsureHostsIni()
}

// runbook creates the "playbooks/" directory and a skeleton runbook
// playbook at "playbooks/runbook.yml". The generated playbook targets all
// hosts with become enabled, fatal error escalation, and fact gathering.
// Task, handler, and variable sections are present as commented-out
// placeholders.
//
// An error is returned if the directory or file cannot be created.
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

// yamlIgnore creates the ".yamlignore" file, which tells yamllint to skip
// specified files. By default only "secrets.yml" is excluded.
//
// An error is returned if the file cannot be created or written.
func yamlIgnore() error {
	fmt.Println("  ...  .yamlignore")

	content := []byte(`secrets.yml
`)

	if err := filesystem.EnsureFileExist(".yamlignore", content); err != nil {
		return err
	}

	return nil
}

// yamlLint creates the ".yamlint" configuration file for yamllint. The
// configuration extends the default ruleset with customized settings for
// braces, comments, indentation, line length (max 125), and octal value
// handling. It also references .gitignore and .yamlignore for file
// exclusion.
//
// An error is returned if the file cannot be created or written.
func yamlLint() error {
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

// vagrantFile creates the "Vagrantfile" for the development environment.
// The generated configuration defines a VirtualBox-backed, multi-VM setup
// with the following characteristics:
//
//   - Shared settings: insecure SSH key insertion disabled, synced folders
//     disabled, 600-second boot timeout, and a connectivity ping to the
//     host at 192.168.57.1 on every provision.
//   - VirtualBox provider: 4 CPUs, 4 GB RAM, headless mode, VMSVGA
//     graphics, IOAPIC enabled, and guest additions checking disabled.
//   - "debian" VM: debian/bookworm64 box at 192.168.57.5.
//   - "alma" VM: almalinux/10 box at 192.168.57.6.
//
// An error is returned if the file cannot be created or written.
func vagrantFile() error {
	fmt.Println("  ...  Vagrantfile")

	content := []byte(`VM_CPUS   = 2
VM_MEMORY = 4096

host = Vagrant::Util::Platform.platform

if host =~ /mswin|mingw|cygwin/
  hyperv_available = File.exist?("C:/Windows/System32/vmms.exe")
  ENV["VAGRANT_DEFAULT_PROVIDER"] ||= (hyperv_available ? "hyperv" : "virtualbox")
elsif host =~ /linux/
  libvirt_ok = system("virsh --version >/dev/null 2>&1")
  ENV["VAGRANT_DEFAULT_PROVIDER"] ||= (libvirt_ok ? "libvirt" : "virtualbox")
else
  ENV["VAGRANT_DEFAULT_PROVIDER"] ||= "virtualbox"
end

Vagrant.configure("2") do |config|
  config.ssh.insert_key = false
  if Vagrant.has_plugin?("vagrant-vbguest")
    config.vbguest.auto_update = false
  end
  config.vm.boot_timeout = 3600
  config.vm.box_check_update = true
  config.vm.synced_folder ".", "/vagrant", disabled: true

  if ENV["VAGRANT_DEFAULT_PROVIDER"] == "hyperv"
    config.vm.network "public_network", bridge: "Default Switch"
  end

  config.vm.provider "virtualbox" do |vb|
    vb.gui    = false
    vb.cpus   = VM_CPUS
    vb.memory = VM_MEMORY
    vb.check_guest_additions = false
    vb.customize [ "modifyvm", :id, "--uartmode1", "disconnected" ]
    vb.customize [ "modifyvm", :id, "--graphicscontroller", "vmsvga" ]
    vb.customize [ "modifyvm", :id, "--ioapic", "on" ]
  end

  config.vm.provider "hyperv" do |hv|
    hv.cpus   = VM_CPUS
    hv.memory = VM_MEMORY
    hv.enable_enhanced_session_mode = false
    hv.auto_start_action = "Nothing"
  end

  config.vm.provider "libvirt" do |lv|
    lv.cpus   = VM_CPUS
    lv.memory = VM_MEMORY
    lv.driver = "kvm"
  end

  config.vm.define "debian" do |c|
    c.vm.box      = "dcjulian29/debian-13"
    c.vm.hostname = "debian.dev"
  end

  config.vm.define "alma" do |c|
    c.vm.box      = "dcjulian29/almalinux-10"
    c.vm.hostname = "alma.dev"
  end
end
`)

	if err := filesystem.EnsureFileExist("Vagrantfile", content); err != nil {
		return err
	}

	return nil
}
