/*
Copyright © 2023 Julian Easterling julian@julianscorner.com

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
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.szostok.io/version/extension"
)

var (
	cfgFile          string
	folderPath       string
	workingDirectory string
	Info             = Teal
	Warn             = Yellow
	Fatal            = Red
	Black            = Color("\033[1;30m%s\033[0m")
	Red              = Color("\033[1;31m%s\033[0m")
	Green            = Color("\033[1;32m%s\033[0m")
	Yellow           = Color("\033[1;33m%s\033[0m")
	Purple           = Color("\033[1;34m%s\033[0m")
	Magenta          = Color("\033[1;35m%s\033[0m")
	Teal             = Color("\033[1;36m%s\033[0m")
	White            = Color("\033[1;37m%s\033[0m")

	rootCmd = &cobra.Command{
		Use:   "ansible-dev",
		Short: "ansible-dev enables development of Ansible playbooks, roles, and modules.",
		Long: `ansible-dev integrates with Vagrant to enable users to define, develop, and test Ansible
playbooks, roles, and modules. It allows users to define and manage infrastructure resources and
uses the providers automation engine to provision and run plays.

By utilizing Ansible's playbooks, roles, and modules, developers can automate the deployment of
software applications across multiple hosting providers, reducing the time and effort required to
manage complex infrastructure environments. These playbooks, roles, and modules enable developers to
create, manage, and provision infrastructure resources like virtual machines, load balancers, and
databases.`,
	}
)

func Execute() {
	workingDirectory, _ = os.Getwd()

	rootCmd.AddCommand(
		extension.NewVersionCobraCmd(
			extension.WithUpgradeNotice("dcjulian29", "ansible-dev"),
		),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	pwd, _ := os.Getwd()

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "specify configuration file")
	rootCmd.PersistentFlags().StringVar(&folderPath, "path", pwd, "path to development folder")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		pwd, err := os.Getwd()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.AddConfigPath(pwd)
		viper.SetConfigType("yml")
		viper.SetConfigName(".ansible-dev")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}

func ensureDir(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return err
		}
	}

	return nil
}

func ensureAnsibleDirectory() {
	if workingDirectory != folderPath {
		if err := os.Chdir(folderPath); err != nil {
			fmt.Println("Unable to access development environment folder!")
			os.Exit(1)
		}
	}
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

func executeExternalProgram(program string, params ...string) {
	cmd := exec.Command(program, params...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		ensureWorkingDirectoryAndExit()
	}
}

func removeDir(dirPath string) {
	if dirExists(dirPath) {
		files, err := filepath.Glob(filepath.Join(dirPath, "*"))
		if err != nil {
			fmt.Println(err)
			ensureWorkingDirectoryAndExit()
		}

		for _, file := range files {
			if err := os.RemoveAll(file); err != nil {
				fmt.Println(err)
				ensureWorkingDirectoryAndExit()
			}
		}

		if err := os.Remove(dirPath); err != nil {
			fmt.Println(err)
			ensureWorkingDirectoryAndExit()
		}
	}
}

func removeFile(filePath string) {
	if fileExists(filePath) {
		if err := os.Remove(filePath); err != nil {
			fmt.Println(err)
			ensureWorkingDirectoryAndExit()
		}
	}
}

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}

	return sprint
}

func ensureVagrantfile() {
	if !fileExists("Vagrantfile") {
		fmt.Println(Fatal("ERROR: Can't find the Vagrantfile!"))
		ensureWorkingDirectoryAndExit()
	}
}

func fileHash(filePath string) string {
	hash := sha256.New()

	sourceFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(Fatal("error opening file:", err))
	}

	defer sourceFile.Close()

	if _, err := io.Copy(hash, sourceFile); err != nil {
		fmt.Println(Fatal("error calculating hash:", err))
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func scanDirectory(dir_path string, ignore []string) ([]string, []string) {
	folders := []string{}
	files := []string{}

	filepath.Walk(dir_path, func(path string, f os.FileInfo, err error) error {
		_continue := false

		for _, i := range ignore {
			if strings.Contains(path, i) {
				_continue = true
			}
		}

		if !_continue {
			f, err = os.Stat(path)
			if err != nil {
				fmt.Println(Fatal("ERROR: Scanning '" + dir_path + "!"))
			}

			f_mode := f.Mode()

			if f_mode.IsDir() {
				folders = append(folders, path)
			} else if f_mode.IsRegular() {
				files = append(files, path)
			}
		}

		return nil
	})

	return folders, files
}
