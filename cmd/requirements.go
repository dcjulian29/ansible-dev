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
	"io"
	"os"

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
		fmt.Println(Fatal("ERROR: Requirements file is not present!"))
		return
	}
}

func readRequirementsFile() (Requirements, error) {
	ensureRequirementsFile()

	var requirements Requirements

	file, err := os.Open("requirements.yml")
	if err != nil {
		fmt.Println(Fatal(err))
		return requirements, err
	}

	defer file.Close()

	data, err := io.ReadAll(file)

	if err != nil {
		fmt.Println(Fatal(err))
		return requirements, err
	}

	if err := yaml.Unmarshal(data, &requirements); err != nil {
		fmt.Println(Fatal(err))
		return requirements, err
	}

	return requirements, nil
}

func writeRequirementsFile(requirements Requirements) error {
	file, err := os.OpenFile("requirements.yml", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(Fatal(err))
		return err
	}

	defer file.Close()

	data, err := yaml.Marshal(requirements)

	if err != nil {
		fmt.Println(Fatal(err))
		return err
	}

	if _, err := file.Write(data); err != nil {
		fmt.Println(Fatal(err))
		return err
	}

	return nil
}
