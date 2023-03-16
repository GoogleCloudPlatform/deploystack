// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Stack represents the input config and output settings for this DeployStack
type Stack struct {
	Settings Settings
	Config   Config
}

// NewStack returns an initialized Stack
func NewStack() Stack {
	s := Stack{}
	s.Settings = Settings{}
	return s
}

func (s *Stack) findAndReadConfig(path string) (Config, error) {
	config := Config{}

	candidates := []string{
		".deploystack/deploystack.yaml",
		".deploystack/deploystack.json",
		"deploystack.json",
	}

	configPath := ""
	for _, v := range candidates {
		candidate := filepath.Join(path, v)
		if _, err := os.Stat(candidate); err == nil {
			configPath = candidate
			break
		}

	}

	if configPath == "" {
		return config, ErrConfigNotExist
	}

	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, fmt.Errorf("unable to find or read config (%s) file: %s", configPath, err)
	}

	switch filepath.Ext(configPath) {
	case ".yaml":
		config, err = NewConfigYAML(content)
		if err != nil {
			return config, fmt.Errorf("unable to parse config file: %s", err)
		}
		return config, nil
	default:
		config, err = NewConfigJSON(content)
		if err != nil {
			return config, fmt.Errorf("unable to parse config file: %s", err)
		}

	}
	return config, nil
}

// ErrConfigNotExist is what happens when a config file either does not exist
// or exists but is not readable.
var ErrConfigNotExist = fmt.Errorf("could not find and parse a config file")

func (s *Stack) findDSFolder(path, folder string) (string, error) {
	switch folder {
	case "messages":
		if s.Config.PathMessages != "" {
			return s.Config.PathMessages, nil
		}
	case "scripts":
		if s.Config.PathScripts != "" {
			return s.Config.PathScripts, nil
		}
	}

	dsPath := filepath.Join(path, folder)

	if _, err := os.Stat(dsPath); err == nil {
		return dsPath, nil
	}

	dsPath = filepath.Join(path, ".deploystack", folder)

	if _, err := os.Stat(dsPath); err == nil {
		return dsPath, nil
	}

	return fmt.Sprintf("./%s", folder), fmt.Errorf("requirement (%s) was not found either in the root, or in .deploystack folder nor was it set in deploystack.json", folder)
}

func (s *Stack) findTFFolder(path string) (string, error) {
	if s.Config.PathTerraform != "" {
		return s.Config.PathTerraform, nil
	}

	mains := []string{}

	err := filepath.Walk(path, func(walkpath string, info os.FileInfo, err error) error {

		if info == nil {
			return fmt.Errorf("info is nil: walkpath: %s err: %s", walkpath, err)
		}

		if info.Name() == "main.tf" {
			dir := filepath.Dir(walkpath)
			mains = append(mains, dir)
			return err
		}
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("findTFFolder: could not find a terraform folder:, %s", err)
	}

	// I want the top most main file here. And that should be the shortest
	sort.Slice(mains, func(i, j int) bool {
		return len(mains[i]) < len(mains[j])
	})

	if len(mains) > 0 {
		return filepath.Rel(path, mains[0])
	}

	return "", nil
}

// FindAndRead figures out a default config, or reads it if it is there
// has option to insure various things and folders exist
func (s *Stack) FindAndRead(path string, required bool) error {
	errs := []error{}

	config, err := s.findAndReadConfig(path)
	s.Config = config
	errs = append(errs, err)

	tfPath, err := s.findTFFolder(path)
	s.Config.PathTerraform = tfPath
	errs = append(errs, err)

	scriptPath, err := s.findDSFolder(path, "scripts")
	s.Config.PathScripts, _ = filepath.Rel(path, scriptPath)
	errs = append(errs, err)

	messagePath, err := s.findDSFolder(path, "messages")
	s.Config.PathMessages, _ = filepath.Rel(path, messagePath)
	errs = append(errs, err)

	if config.Description == "" {
		descText := fmt.Sprintf("%s/description.txt", messagePath)
		description, err := ioutil.ReadFile(descText)
		s.Config.Description = string(description)
		errs = append(errs, err)
	}

	s.Config.convertHardset()
	s.Config.defaultAuthorSettings()

	if required && len(errs) > 0 {
		return errs[0]
	}

	return nil
}

// FindAndReadRequired finds and reads in a Config from a json file.
func (s *Stack) FindAndReadRequired(path string) error {
	return s.FindAndRead(path, true)
}

// AddSetting stores a setting key/value pair.
func (s *Stack) AddSetting(key, value string) {
	s.Settings.Add(key, value)
}

// AddSettingComplete passes a completely intact setting to the underlying
// setting structure
func (s *Stack) AddSettingComplete(set Setting) {
	s.Settings.AddComplete(set)
}

// GetSetting returns a setting value.
func (s *Stack) GetSetting(key string) string {
	set := s.Settings.Find(key)

	if set != nil {
		return set.Value
	}

	return ""
}

// DeleteSetting removes a setting value.
func (s *Stack) DeleteSetting(key string) {
	for i, v := range s.Settings {
		if v.Name == key {
			s.Settings = append(s.Settings[:i], s.Settings[i+1:]...)
		}
	}

}

// Terraform returns all of the settings as a Terraform variables format.
func (s Stack) Terraform() string {
	result := strings.Builder{}

	s.Settings.Sort()

	for _, v := range s.Settings {
		if v.Name == "" {
			continue
		}
		label := v.TFvarsName()

		if label == "project_name" {
			continue
		}

		if label == "stack_name" {
			continue
		}

		if len(v.Value) == 0 && len(v.List) == 0 && v.Map == nil {
			continue
		}

		result.WriteString(v.TFVars())

	}

	return result.String()
}

// TerraformFile exports TFVars format to input file.
func (s Stack) TerraformFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.WriteString(s.Terraform()); err != nil {
		return err
	}

	return nil
}
