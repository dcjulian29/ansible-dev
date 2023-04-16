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
	"os"

	"github.com/spf13/cobra"
)

var (
	force            bool
	workingDirectory string

	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize an development environment for Ansible development",
		Long: `Initialize an development environment for Ansible development by creating the folder
	structure and generating the needed files to quickly set up a virtual environment
	ready for development. Vagrant can be used to manage the environment and connect
	to troubleshoot and/or validate.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Initializing development environment...")

			init_env()
		},
	}
)

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().BoolVarP(&force, "force", "f", false, "overwrite an existing development environment")
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func ensureDir(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return err
		}
	}

	return nil
}

func ensureWorkingDirectoryAndExit() {
	if workingDirectory != folderPath {
		if err := os.Chdir(workingDirectory); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	os.Exit(0)
}

func init_env() {
	workingDirectory, _ = os.Getwd()

	if workingDirectory != folderPath {
		if _, err := os.Stat(folderPath); os.IsNotExist(err) {
			fmt.Println("Creating development environment folder...")
			if err := os.MkdirAll(folderPath, 0755); err != nil {
				fmt.Println("Unable to create development environment folder!")
				os.Exit(1)
			}
		}

		if err := os.Chdir(folderPath); err != nil {
			fmt.Println("Unable to access development environment folder!")
			os.Exit(1)
		}
	}

	if fileExists("ansible.cfg") && !force {
		fmt.Println("ERROR: The folder for the ansible development environment already contains an Ansible context and force was not provided.")
		ensureWorkingDirectoryAndExit()
	}

	fmt.Println("    ...   roles/")
	if err := ensureDir("roles"); err != nil {
		fmt.Println(err)
		ensureWorkingDirectoryAndExit()
	}

	fmt.Println("    ...   playbooks/")
	if err := ensureDir("playbooks"); err != nil {
		fmt.Println(err)
		ensureWorkingDirectoryAndExit()
	}

	ansible_cfg()
	ansible_lint()
	inventory_file()
	vagrant_file()

	ensureWorkingDirectoryAndExit()
}

func ansible_cfg() {
	fmt.Println("    ...   ansible.cfg")

	file, err := os.Create("ansible.cfg")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	content := []byte(`[defaults]
any_errors_fatal            = true
collections_path            = ./collections
duplicate_dict_key          = error
error_on_undefined_vars     = true
gathering                   = smart
host_key_checking           = false
inventory                   = ./hosts.ini
log_path                    = ./ansible.log
roles_path                  = ./roles
stdout_callback             = community.general.yaml
verbosity                   = 1

[diff]
always                      = true
`)

	if _, err = file.Write(content); err != nil {
		fmt.Println(err)
		return
	}
}

func inventory_file() {
	fmt.Println("    ...   hosts.ini")
	file, err := os.Create("hosts.ini")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	content := []byte(`[ansibledev]
debian
rocky

[vagrant]
debian
rocky
ubuntu
alma
fedora
`)

	if _, err = file.Write(content); err != nil {
		fmt.Println(err)
		return
	}
}

func vagrant_file() {
	fmt.Println("    ...   Vagrantfile")
	file, err := os.Create("Vagrantfile")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	content := []byte(`Vagrant.configure("2") do |config|
  config.ssh.insert_key = false
  config.ssh.extra_args = "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null"
  config.vm.synced_folder ".", "/vagrant", disabled: true
  config.vm.network :forwarded_port, guest: 22, host: 2220, id: "ssh", disabled: true

  config.vm.provider "virtualbox" do |vb|
    vb.gui = false
    vb.memory = 2048
    vb.cpus = 2
    vb.check_guest_additions = false
    vb.customize [ "modifyvm", :id, "--uartmode1", "disconnected" ]
    vb.customize [ "modifyvm", :id, "--graphicscontroller", "vmsvga"]
  end

  config.vm.define "debian" do |c|
    c.vm.box = "debian/bullseye64"
    c.vm.hostname = "debian.test"
    c.vm.network "private_network", ip: "192.168.57.5"
    c.vm.network :forwarded_port, guest: 22, host: 8005, id: 'ssh'
  end

  config.vm.define "rocky" do |c|
    c.vm.box = "rockylinux/9"
    c.vm.hostname = "rocky.test"
    c.vm.network "private_network", ip: "192.168.57.6"
    c.vm.network :forwarded_port, guest: 22, host: 8006, id: 'ssh'
  end

  config.vm.define "alma" do |c|
    c.vm.box = "generic/alma9"
    c.vm.hostname = "alma.test"
    c.vm.network "private_network", ip: "192.168.57.7"
    c.vm.network :forwarded_port, guest: 22, host: 8007, id: 'ssh'
  end

  config.vm.define "fedora" do |c|
    c.vm.box = "generic/fedora36"
    c.vm.hostname = "fedora.test"
    c.vm.network "private_network", ip: "192.168.57.8"
    c.vm.network :forwarded_port, guest: 22, host: 8008, id: 'ssh'
  end

  config.vm.define "ubuntu" do |c|
    c.vm.box = "generic/ubuntu2204"
    c.vm.hostname = "ubuntu.test"
    c.vm.network "private_network", ip: "192.168.57.9"
    c.vm.network :forwarded_port, guest: 22, host: 8009, id: 'ssh'
  end

end
`)

	if _, err = file.Write(content); err != nil {
		fmt.Println(err)
		return
	}
}

func ansible_lint() {
	fmt.Println("    ...   .ansible-lint")
	file, err := os.Create(".ansible-lint")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	content := []byte(`warn_list:
  - internal-error
  - no-handler
  - experimental

kinds:
  - roles: "roles/"
  - playbook: "playbooks/*.{yml,yaml}"
  - tasks: "**/tasks/*.{yml,yaml}"
  - vars: "**/vars/*.{yml,yaml}"
  - meta: "**/meta/main.yml"
  - yaml: "**/*.{yml,yaml}"

verbosity: 1
`)

	if _, err = file.Write(content); err != nil {
		fmt.Println(err)
		return
	}
}
