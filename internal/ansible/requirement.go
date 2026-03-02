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
package ansible

import (
	"fmt"
	"io"
	"os"

	"github.com/dcjulian29/go-toolbox/filesystem"
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

func EnsureRequirementsFile() error {
	exist := RequirementsFileExist()
	if !exist {
		return fmt.Errorf("requirements.yml file is not present")
	}

	return nil
}

func ReadRequirements() (Requirements, error) {
	err := EnsureRequirementsFile()
	if err != nil {
		return Requirements{}, err
	}

	var requirements Requirements

	file, err := os.Open("requirements.yml")
	if err != nil {
		return Requirements{}, err
	}

	defer file.Close() //nolint:errcheck

	data, err := io.ReadAll(file)
	if err != nil {
		return Requirements{}, err
	}

	err = yaml.Unmarshal(data, &requirements)
	if err != nil {
		return Requirements{}, err
	}

	return requirements, nil
}

func RequirementsFileExist() bool {
	return filesystem.FileExists("requirements.yml")
}

func SaveRequirements(requirements Requirements) error {
	file, err := os.OpenFile("requirements.yml", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer file.Close() //nolint:errcheck

	data, err := yaml.Marshal(requirements)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
