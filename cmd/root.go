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
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.szostok.io/version/extension"
)

var cfgFile string

var rootCmd = &cobra.Command{
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

func Execute() {
	rootCmd.AddCommand(
		extension.NewVersionCobraCmd(
			extension.WithUpgradeNotice("dcjulian29", "ansible-dev"),
		),
	)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "specify configuration file")
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
