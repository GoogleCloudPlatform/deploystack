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

package deploystack

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Stack represents the input config and output settings for this DeployStack
type Stack struct {
	Settings map[string]string
	Config   Config
}

// NewStack returns an initialized Stack
func NewStack() Stack {
	s := Stack{}
	s.Settings = make(map[string]string)
	return s
}

func (s *Stack) findAndReadConfig() (Config, error) {
	config := Config{}

	candidates := []string{
		".deploystack/deploystack.yaml",
		".deploystack/deploystack.json",
		"deploystack.json",
	}

	configPath := ""
	wd, _ := os.Getwd()
	for _, v := range candidates {
		file := fmt.Sprintf("%s/%s", wd, v)

		if _, err := os.Stat(file); err == nil {
			configPath = file
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

var ErrConfigNotExist = fmt.Errorf("could not file and parse a config file")

func (s *Stack) findDSFolder(c Config, folder string) (string, error) {
	switch folder {
	case "messages":
		if c.PathMessages != "" {
			return c.PathMessages, nil
		}
	case "scripts":
		if c.PathScripts != "" {
			return c.PathScripts, nil
		}
	}

	path := fmt.Sprintf(".deploystack/%s", folder)

	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	path = folder

	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	return fmt.Sprintf("./%s", folder), fmt.Errorf("requirement (%s) was not found either in the root, or in .deploystack folder nor was it set in deploystack.json", folder)
}

func (s *Stack) findTFFolder(c Config) (string, error) {
	if c.PathTerraform != "" {
		return c.PathTerraform, nil
	}

	path := "terraform"

	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	return ".", nil
}

// FindAndReadRequired finds and reads in a Config from a json file.
func (s *Stack) FindAndReadRequired() error {
	config, err := s.findAndReadConfig()
	if err != nil {
		return fmt.Errorf("unable to parse config file: %s", err)
	}

	tfPath, err := s.findTFFolder(config)
	if err != nil {
		return fmt.Errorf("unable to locate terraform folder: %s", err)
	}
	config.PathTerraform = tfPath

	scriptPath, _ := s.findDSFolder(config, "scripts")
	if err != nil {
		log.Printf("WARNING - unable to locate scripts folder, folder not required, : %s", err)
	}
	config.PathScripts = scriptPath

	messagePath, err := s.findDSFolder(config, "messages")
	if err != nil {
		log.Printf("WARNING - unable to locate messages folder, folder not required, : %s", err)
	}
	config.PathMessages = messagePath

	descText := fmt.Sprintf("%s/description.txt", messagePath)
	if _, err := os.Stat(descText); err == nil {
		description, err := ioutil.ReadFile(descText)
		if err != nil {
			return fmt.Errorf("unable to read description file: %s", err)
		}

		config.Description = string(description)
	}

	s.Config = config

	return nil
}

// AddSetting stores a setting key/value pair.
func (s Stack) AddSetting(key, value string) {
	k := strings.ToLower(key)
	s.Settings[k] = value
}

// GetSetting returns a setting value.
func (s Stack) GetSetting(key string) string {
	return s.Settings[key]
}

// DeleteSetting removes a setting value.
func (s Stack) DeleteSetting(key string) {
	delete(s.Settings, key)
}

// Terraform returns all of the settings as a Terraform variables format.
func (s Stack) Terraform() string {
	result := strings.Builder{}

	keys := []string{}
	for i := range s.Settings {
		keys = append(keys, i)
	}

	sort.Strings(keys)

	for _, v := range keys {
		if len(v) < 1 {
			continue
		}
		label := strings.ToLower(strings.ReplaceAll(v, " ", "_"))
		val := s.Settings[v]

		if label == "project_name" {
			continue
		}

		if label == "stack_name" {
			continue
		}

		if len(val) < 1 {
			continue
		}

		if val[0:1] == "[" {
			sb := strings.Builder{}
			sb.WriteString("[")
			tmp := strings.ReplaceAll(val, "[", "")
			tmp = strings.ReplaceAll(tmp, "]", "")
			sl := strings.Split(tmp, ",")

			for i, v := range sl {
				sl[i] = fmt.Sprintf("\"%s\"", v)
			}

			delimtext := strings.Join(sl, ",")

			sb.WriteString(delimtext)
			sb.WriteString("]")
			set := sb.String()
			set = strings.ReplaceAll(set, "\"\"", "")

			result.WriteString(fmt.Sprintf("%s=%s\n", label, set))
			continue
		}

		result.WriteString(fmt.Sprintf("%s=\"%s\"\n", label, val))

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
