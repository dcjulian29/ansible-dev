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
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Requirements struct {
	Collections []Collection `yaml:"collections"`
	Roles       []Role       `yaml:"roles"`
}

type Collection struct {
	Name    string `yaml:"name"`
	Source  string `yaml:"src"`
	Type    string `yaml:"type"`
	Version string `yaml:"version"`
}

type Role struct {
	Name    string `yaml:"name"`
	Source  string `yaml:"src"`
	Version string `yaml:"version"`
}

func ensureRequirementsFile() {
	if !fileExists("requirements.yml") {
		fmt.Println(Fatal("requirements.yml file is not present!"))
		return
	}
}

func readRequirementsFile() (Requirements, error) {
	ensureRequirementsFile()

	var requirements Requirements

	file, err := os.Open("requirements.yml")
	cobra.CheckErr(err)

	defer file.Close()

	data, err := io.ReadAll(file)
	cobra.CheckErr(err)

	cobra.CheckErr(yaml.Unmarshal(data, &requirements))

	return requirements, nil
}

func writeRequirementsFile(requirements Requirements) error {
	ensureRequirementsFile()

	file, err := os.OpenFile("requirements.yml", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	cobra.CheckErr(err)

	defer file.Close()

	data, err := yaml.Marshal(requirements)
	cobra.CheckErr(err)

	_, err = file.Write(data)
	cobra.CheckErr(err)

	return nil
}
