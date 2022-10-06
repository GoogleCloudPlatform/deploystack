// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package deploystack provides a series of interfaces for getting Google Cloud
// settings and configurations for use with DeplyStack
package deploystack

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/nyaruka/phonenumbers"
	"google.golang.org/api/option"
	"gopkg.in/yaml.v2"
)

const (
	// TERMCYAN is the terminal code for cyan text
	TERMCYAN = "\033[0;36m"
	// TERMCYANB is the terminal code for bold cyan text
	TERMCYANB = "\033[1;36m"
	// TERMCYANREV is the terminal code for black on cyan text
	TERMCYANREV = "\u001b[46m"
	// TERMRED is the terminal code for red text
	TERMRED = "\033[0;31m"
	// TERMREDB is the terminal code for bold red text
	TERMREDB = "\033[1;31m"
	// TERMREDREV is the terminal code for black on red text
	TERMREDREV = "\033[41m"
	// TERMCLEAR is the terminal code for the clear out color text
	TERMCLEAR = "\033[0m"
	// TERMCLEARSCREEN is the terminal code for clearning the whole screen.
	TERMCLEARSCREEN = "\033[2J"
	// TERMGREY is the terminal code for grey text
	TERMGREY = "\033[1;30m"

	// DefaultRegion is the default compute region used in compute calls.
	DefaultRegion = "us-central1"
	// DefaultMachineType is the default compute machine type used in compute calls.
	DefaultMachineType = "n1-standard"
	// DefaultImageProject is the default project for images used in compute calls.
	DefaultImageProject = "debian-cloud"
	// DefaultImageFamily is the default project for images used in compute calls.
	DefaultImageFamily = "debian-11"
)

// ClearScreen will clear out a terminal screen.
func ClearScreen() {
	fmt.Println(TERMCLEARSCREEN)
}

var (
	// ErrorCustomNotValidPhoneNumber is the error you get when you fail phone
	// number validation.
	ErrorCustomNotValidPhoneNumber = fmt.Errorf("not a valid phone number")
	// ErrorBillingInvalidAccount is the error you get if you pass in a bad
	// Billing Account ID
	ErrorBillingInvalidAccount = fmt.Errorf("not a valid billing account")
	// ErrorBillingNoPermission is the error you get if the user lacks billing
	// related permissions
	ErrorBillingNoPermission = fmt.Errorf("user lacks permission")
	// ErrorProjectCreateTooLong is an error when you try to create a project
	// wuth more than 30 characters
	ErrorProjectCreateTooLong = fmt.Errorf("project_id contains too many characters, limit 30")
	// ErrorProjectInvalidCharacters is an error when you try and pass bad
	// characters into a CreateProjectCall
	ErrorProjectInvalidCharacters = fmt.Errorf("project_id contains invalid characters")
	// ErrorProjectAlreadyExists is an error when you try and create a project
	// That already exists
	ErrorProjectAlreadyExists = fmt.Errorf("project_id already exists")
	// Divider is a text element that draws a horizontal line
	Divider   = ""
	opts      = option.WithCredentialsFile("")
	credspath = ""
)

func init() {
	var err error
	Divider, err = BuildDivider(0)
	if err != nil {
		log.Fatal(err)
	}
}

// BuildDivider captures the size of the terminal screen to build a horizontal
// divider.
func BuildDivider(width int) (string, error) {
	de := 80
	if width == 0 {
		cmd := exec.Command("stty", "size")
		cmd.Stdin = os.Stdin
		out, err := cmd.Output()
		if err != nil {
			width = de
		}

		sl := strings.Split(string(out), " ")

		if len(sl) > 1 {
			width, err = strconv.Atoi(strings.TrimSpace(sl[1]))
			if err != nil {
				width = de
			}
		}
	}

	var sb strings.Builder

	for i := 0; i < width; i++ {
		sb.WriteString("*")
	}
	return sb.String(), nil
}

// Flags is a collection variables that can be passed in from the CLI
type Flags struct {
	Project string            `json:"project"`
	Region  string            `json:"region"`
	Zone    string            `json:"zone"`
	Custom  map[string]string `json:"custom"`
}

// HandleFlags consolidates all of the cli flag login in the package instead of
// relegating that to the calling file. Not super idiomatic, but allows us
// to leave all of this code in one place.
func HandleFlags() Flags {
	f := Flags{}
	m := make(map[string]string)
	projectPtr := flag.String("project", "", "A Google Cloud Project ID")
	regionPtr := flag.String("region", "", "A Google Cloud Region")
	zonePtr := flag.String("zone", "", "A Google Cloud Zone")
	customPtr := flag.String("custom", "", "A list of custom variables that can be passed in")

	flag.Parse()

	f.Project = *projectPtr
	f.Region = *regionPtr
	f.Zone = *zonePtr

	rawString := *customPtr

	cSl := strings.Split(rawString, ",")

	for _, v := range cSl {
		if len(v) == 0 {
			continue
		}
		rawVK := strings.ReplaceAll(v, " ", "=")
		kv := strings.Split(rawVK, "=")
		fmt.Printf("kv %+v\n", kv)
		m[kv[0]] = kv[1]
	}
	f.Custom = m

	return f
}

// Config represents the settings this app will collect from a user. It should
// be in a json file. The idea is minimal programming has to be done to setup
// a DeployStack and export out a tfvars file for terraform part of solution.
type Config struct {
	Title                string            `json:"title" yaml:"title"`
	Description          string            `json:"description" yaml:"description"`
	Duration             int               `json:"duration" yaml:"duration"`
	Project              bool              `json:"collect_project" yaml:"collect_project"`
	ProjectNumber        bool              `json:"collect_project_number" yaml:"collect_project_number"`
	BillingAccount       bool              `json:"collect_billing_account" yaml:"collect_billing_account"`
	Domain               bool              `json:"register_domain" yaml:"register_domain"`
	Region               bool              `json:"collect_region" yaml:"collect_region"`
	RegionType           string            `json:"region_type" yaml:"region_type"`
	RegionDefault        string            `json:"region_default" yaml:"region_default"`
	Zone                 bool              `json:"collect_zone" yaml:"collect_zone"`
	HardSet              map[string]string `json:"hard_settings" yaml:"hard_settings"`
	CustomSettings       []Custom          `json:"custom_settings" yaml:"custom_settings"`
	ConfigureGCEInstance bool              `json:"configure_gce_instance" yaml:"configure_gce_instance"`
	DocumentationLink    string            `json:"documentation_link" yaml:"documentation_link"`
	PathTerraform        string            `json:"path_terraform" yaml:"path_terraform"`
	PathMessages         string            `json:"path_messages" yaml:"path_messages"`
	PathScripts          string            `json:"path_scripts" yaml:"path_scripts"`
}

// Custom represents a custom setting that we would like to collect from a user
// We will collect these settings from the user before continuing.
type Custom struct {
	Name           string   `json:"name"  yaml:"name"`
	Description    string   `json:"description"  yaml:"description"`
	Default        string   `json:"default"  yaml:"default"`
	Value          string   `json:"value"  yaml:"value"`
	Options        []string `json:"options"  yaml:"options"`
	PrependProject bool     `json:"prepend_project"  yaml:"prepend_project"`
	Validation     string   `json:"validation,omitempty"  yaml:"validation,omitempty"`
	project        string
}

// Collect will collect a value for a Custom from a user
func (c *Custom) Collect() error {
	fmt.Printf("%s%s: %s\n", TERMCYANB, c.Description, TERMCLEAR)

	def := c.Default

	if c.PrependProject {
		def = fmt.Sprintf("%s-%s", c.project, c.Default)
	}

	if len(c.Options) > 0 {
		c.Value = listSelect(toLabeledValueSlice(c.Options), def).Value
		return nil
	}

	result := ""

	if len(c.Default) > 0 {
		fmt.Printf("Enter value, or just [enter] for %s%s%s\n", TERMCYANB, c.Default, TERMCLEAR)
	} else {
		fmt.Printf("Enter value: \n")
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil && err.Error() != "EOF" {
			return err
		}
		text = strings.Replace(text, "\n", "", -1)
		result = text

		if len(text) == 0 {
			text = def
		}

		switch c.Validation {

		case "phonenumber":
			num, err := massagePhoneNumber(text)
			if err != nil {
				fmt.Printf("%sThat's not a valid phone number. Please try again.%s\n", TERMRED, TERMCLEAR)
				continue
			}
			result = num
		case "integer":
			_, err := strconv.Atoi(text)
			if err != nil {
				fmt.Printf("%sYour answer '%s' not a valid integer. Please try again.%s\n", TERMRED, text, TERMCLEAR)
				continue
			}
			result = text
		case "yesorno":
			text = strings.TrimSpace(strings.ToLower(text))
			yesList := " yes y "
			noList := " no n "

			if !strings.Contains(yesList+noList, text) {
				fmt.Printf("%sYour answer '%s' is neither 'yes' nor 'no'. Please try again.%s\n", TERMRED, text, TERMCLEAR)
				continue
			}

			if strings.Contains(yesList, text) {
				result = "yes"
			}

			if strings.Contains(noList, text) {
				result = "no"
			}

		default:
			result = text
		}

		c.Value = result
		if len(result) > 0 {
			break
		}

	}

	return nil
}

func massagePhoneNumber(s string) (string, error) {
	num, err := phonenumbers.Parse(s, "US")
	if err != nil {
		return "", ErrorCustomNotValidPhoneNumber
	}
	result := phonenumbers.Format(num, phonenumbers.INTERNATIONAL)
	result = strings.Replace(result, " ", ".", 1)
	result = strings.ReplaceAll(result, "-", "")
	result = strings.ReplaceAll(result, " ", "")

	return result, nil
}

// Customs are a slice of Custom variables.
type Customs []Custom

// Get returns one Custom Variable
func (cs Customs) Get(name string) Custom {
	for _, v := range cs {
		if v.Name == name {
			return v
		}
	}

	return Custom{}
}

// Collect calls the collect method of all of the Custom variables in the
// collection in the order in which they were placed there.
func (cs *Customs) Collect() error {
	for i, v := range *(cs) {
		if err := v.Collect(); err != nil {
			return fmt.Errorf("error getting custom value (%s) from user:  %s", v.Name, err)
		}
		(*cs)[i] = v
	}

	return nil
}

// PrintHeader prints out the header for a DeployStack
func (c Config) PrintHeader() {
	fmt.Printf("%s\n", Divider)
	fmt.Printf("%s%s%s\n", TERMCYANB, c.Title, TERMCLEAR)
	fmt.Printf("%s\n", c.Description)

	timestring := "minute"
	if c.Duration > 1 {
		timestring = "minutes"
	}

	fmt.Printf("It's going to take around %s%d %s%s\n", TERMCYAN, c.Duration, timestring, TERMCLEAR)

	if c.DocumentationLink != "" {
		fmt.Printf("\nIf you would like more information about this stack, please read the \n")
		fmt.Printf("documentation at: \n%s%s%s \n", TERMCYANB, c.DocumentationLink, TERMCLEAR)
	}

	fmt.Printf("%s\n", Divider)
}

// Process runs through all of the options in a config and collects all of the
// necessary data from users.
func (c Config) Process(s *Stack, output string) error {
	Start()
	c.PrintHeader()
	var project, region, zone, projectnumber, billingaccount, projectName string
	var err error

	for i, v := range c.HardSet {
		s.AddSetting(i, v)
	}

	project = s.GetSetting("project_id")
	projectName = s.GetSetting("project_name")
	region = s.GetSetting("region")
	zone = s.GetSetting("zone")

	if c.Project && len(project) == 0 {
		project, projectName, err = ProjectManage()
		if err != nil {
			handleProcessError(fmt.Errorf("error managing project settings: %s", err))
		}
		s.AddSetting("project_id", project)
		s.AddSetting("project_name", projectName)
	}

	if c.ConfigureGCEInstance {
		basename := s.GetSetting("basename")
		config, err := GCEInstanceManage(project, basename)
		if err != nil {
			handleProcessError(fmt.Errorf("error managing compute instance settings: %s", err))
		}

		for i, v := range config {
			s.AddSetting(i, v)
		}

	}

	region = s.GetSetting("region")
	zone = s.GetSetting("zone")

	if c.Region && len(region) == 0 {
		region, err = RegionManage(project, c.RegionType, c.RegionDefault)
		if err != nil {
			handleProcessError(fmt.Errorf("error managing region settings: %s", err))
		}
		s.AddSetting("Region", region)
	}

	if c.Zone && len(zone) == 0 {

		if !c.Region {
			region, err = RegionManage(project, "compute", DefaultRegion)
			if err != nil {
				handleProcessError(fmt.Errorf("error managing region settings: %s", err))
			}
		}

		zone, err = ZoneManage(project, region)
		if err != nil {
			handleProcessError(fmt.Errorf("error managing zone settings: %s", err))
		}
		s.AddSetting("zone", zone)
	}

	if c.ProjectNumber {
		projectnumber, err = projectNumber(project)
		if err != nil {
			handleProcessError(fmt.Errorf("error managing project number settings: %s", err))
		}
		s.AddSetting("project_number", projectnumber)
	}

	if c.Domain {
		domain, err := DomainManage(s)
		if err != nil {
			handleProcessError(fmt.Errorf("error handling domain registration: %s", err))
		}
		s.AddSetting("domain", domain)
	}

	if c.BillingAccount {

		ba, err := BillingAccountManage()
		if err != nil {
			handleProcessError(fmt.Errorf("error managing billing settings: %s", err))
		}
		billingaccount = ba
		s.AddSetting("billing_account", billingaccount)
	}

	for _, v := range c.CustomSettings {
		temp := s.GetSetting(v.Name)

		if len(temp) < 1 {

			v.project = project

			if err := v.Collect(); err != nil {
				handleProcessError(fmt.Errorf("error getting custom value from user: %s", err))
			}
			s.AddSetting(v.Name, v.Value)
		}

	}

	s.PrintSettings()
	s.TerraformFile(output)
	return nil
}

func handleProcessError(err error) {
	fmt.Printf("\n\n%sThere was an issue collecting the information it takes to run this application.                             %s\n\n", TERMREDREV, TERMCLEAR)
	fmt.Printf("%sYou can try again by typing %sdeploystack install%s at the command prompt  %s\n\n", TERMREDB, TERMREDREV, TERMCLEAR+TERMREDB, TERMCLEAR)
	fmt.Printf("%sIf the issue persists, please report at https://github.com/GoogleCloudPlatform/deploystack/issues %s\n\n", TERMREDB, TERMCLEAR)

	fmt.Printf("Extra diagnostic information:\n")

	if strings.Contains(err.Error(), "invalid token JSON from metadata") {
		fmt.Printf("timed out waiting for API activation, you must authorize API use to continue \n")
	}

	fmt.Println(err)
	os.Exit(1)
}

// NewConfigJSON returns a Config object from a file read.
func NewConfigJSON(content []byte) (Config, error) {
	result := Config{}
	if err := json.Unmarshal(content, &result); err != nil {
		return result, fmt.Errorf("unable to convert content to Config: %s", err)
	}

	return result, nil
}

// NewConfigYAML returns a Config object from a file read.
func NewConfigYAML(content []byte) (Config, error) {
	result := Config{}

	if err := yaml.Unmarshal(content, &result); err != nil {
		return result, fmt.Errorf("unable to convert content to Config: %s", err)
	}

	return result, nil
}

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

// ProcessFlags handles adding the contents of the flags to the stack settings
func (s *Stack) ProcessFlags(f Flags) {
	s.AddSetting("project_id", f.Project)
	s.AddSetting("region", f.Region)
	s.AddSetting("zone", f.Zone)

	for i, v := range f.Custom {
		s.AddSetting(i, v)
	}
}

func (s *Stack) findAndReadConfig() (Config, error) {
	flavor := "json"
	config := Config{}

	configPath := ".deploystack/deploystack.json"
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		flavor = "yaml"
		configPath = ".deploystack/deploystack.yaml"
	}

	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		configPath = "deploystack.json"
	}

	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		return config, fmt.Errorf("config file not present, looking for deploystack.json or .deploystack/deploystack.json")
	}

	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, fmt.Errorf("unable to find or read config file: %s", err)
	}

	if flavor == "yaml" {
		config, err = NewConfigYAML(content)
		if err != nil {
			return config, fmt.Errorf("unable to parse config file: %s", err)
		}
		return config, nil
	}

	config, err = NewConfigJSON(content)
	if err != nil {
		return config, fmt.Errorf("unable to parse config file: %s", err)
	}

	return config, nil
}

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
		return fmt.Errorf("unable to locate messages folder: %s", err)
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

// TODO: deprecate and remove
// ReadConfig reads in a Config from a json file.
func (s *Stack) ReadConfig(file, desc string) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("unable to read config file: %s", err)
	}
	config, err := NewConfigJSON(content)
	if err != nil {
		return fmt.Errorf("unable to parse config file: %s", err)
	}

	if len(desc) > 0 {
		description, err := ioutil.ReadFile(desc)
		if err != nil {
			return fmt.Errorf("unable to read description file: %s", err)
		}

		config.Description = string(description)
	}

	s.Config = config

	return nil
}

// Process passes through a process call to the underlying config.
func (s *Stack) Process(output string) error {
	return s.Config.Process(s, output)
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

// PrintSettings prints the settings to the screen
func (s Stack) PrintSettings() {
	keys := []string{}
	for i := range s.Settings {
		keys = append(keys, i)
	}

	longest := longestLength(toLabeledValueSlice(keys))

	fmt.Printf("\n%sProject Details %s \n", TERMCYANREV, TERMCLEAR)

	if s, ok := s.Settings["project_name"]; ok && len(s) > 0 {
		printSetting("project_name", s, longest)
	}

	if s, ok := s.Settings["project_id"]; ok && len(s) > 0 {
		printSetting("project_id", s, longest)
	}

	if s, ok := s.Settings["project_number"]; ok {
		printSetting("project_number", s, longest)
	}

	ordered := []string{}
	for i, v := range s.Settings {
		if i == "project_id" || i == "project_number" || i == "project_name" {
			continue
		}
		if len(v) < 1 {
			continue
		}

		ordered = append(ordered, i)
	}
	sort.Strings(ordered)

	for i := range ordered {
		key := ordered[i]
		printSetting(key, s.Settings[key], longest)
	}
}

func printSetting(name, value string, longest int) {
	sp := buildSpacer(name, longest)
	formatted := strings.Title(strings.ReplaceAll(name, "_", " "))
	fmt.Printf("%s:%s %s%s%s\n", formatted, sp, TERMCYANB, value, TERMCLEAR)
}

// Section allows for division of tasks in a DeployStack
type Section struct {
	Title string
}

// NewSection creates an initialized section
func NewSection(title string) Section {
	return Section{Title: title}
}

// Open prints out the header for a Section.
func (s Section) Open() {
	fmt.Printf("%s\n", Divider)
	fmt.Printf("%s%s%s\n", TERMCYAN, s.Title, TERMCLEAR)
	fmt.Printf("%s\n", Divider)
}

// Close prints out the footer for a Section.
func (s Section) Close() {
	fmt.Printf("%s\n", Divider)
	fmt.Printf("%s%s - %sdone%s\n", TERMCYAN, s.Title, TERMCYANB, TERMCLEAR)
	fmt.Printf("%s\n", Divider)
}

// ProjectID gets the currently set default project
func ProjectID() (string, error) {
	cmd := exec.Command("gcloud", "config", "get-value", "project")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("cannot get project id: %s ", err)
	}

	return strings.TrimSpace(string(out)), nil
}

// ProjectIDSet sets the currently set default project
func ProjectIDSet(project string) error {
	cmd := exec.Command("gcloud", "config", "set", "project", project)
	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("cannot set project id: %s ", err)
	}

	return nil
}

type LabeledValue struct {
	Value string
	Label string
}

type LabeledValues []LabeledValue

func (l LabeledValues) find(value string) LabeledValue {
	for _, v := range l {
		if v.Value == value {
			return v
		}
	}
	return LabeledValue{}
}

func (l *LabeledValues) sort() {
	sort.Slice(*l, func(i, j int) bool {
		return (*l)[i].Label < (*l)[j].Label
	})
}

func toLabeledValueSlice(sl []string) LabeledValues {
	r := LabeledValues{}

	for _, v := range sl {
		val := LabeledValue{v, v}
		if strings.Contains(v, "|") {
			sl := strings.Split(v, "|")
			val = LabeledValue{sl[0], sl[1]}
		}

		r = append(r, val)
	}
	return r
}

// listSelect presents a slice of strings as a list from which
// the user can select. It also highlights and preesnts behvaior for the
// default
func listSelect(sl LabeledValues, def string) LabeledValue {
	itemCount := len(sl)
	halfcount := int(math.Ceil(float64(itemCount / 2)))
	width := longestLength(sl)
	defaultExists := false

	if itemCount < 11 {
		for i, v := range sl {
			if ok := printWithDefault(i+1, width, v.Value, v.Label, def); ok {
				defaultExists = true
			}
			fmt.Printf("\n")
		}
	} else {

		if float64(halfcount) < float64(itemCount)/2 {
			halfcount++
		}

		for i := 0; i < halfcount; i++ {
			v := sl[i]
			if ok := printWithDefault(i+1, width, v.Value, v.Label, def); ok {
				defaultExists = true
			}

			idx := i + halfcount + 1

			if idx > itemCount {
				fmt.Printf("\n")
				break
			}

			v2 := sl[idx-1]
			if ok := printWithDefault(idx, width, v2.Value, v2.Label, def); ok {
				defaultExists = true
			}

			fmt.Printf("\n")
		}
	}

	answer := sl.find(def)
	reader := bufio.NewReader(os.Stdin)
	if defaultExists {
		fmt.Printf("Choose number from list, or just [enter] for %s%s%s\n", TERMCYANB, answer.Label, TERMCLEAR)
	} else {
		fmt.Printf("Choose number from list.\n")
	}

	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		if len(text) == 0 {
			break
		}

		opt, err := strconv.Atoi(text)
		if err != nil || opt > itemCount {
			fmt.Printf("Please enter a numeric between 1 and %d\n", itemCount)
			fmt.Printf("You entered %s\n", text)
			continue
		}

		answer = sl[opt-1]
		break

	}

	return answer
}

func printWithDefault(idx, width int, value, label, def string) bool {
	sp := buildSpacer(cleanTerminalCharsFromString(label), width)

	if value == def {
		fmt.Printf("%s%2d) %s %s%s", TERMCYANB, idx, label, sp, TERMCLEAR)
		return true
	}
	fmt.Printf("%2d) %s %s", idx, label, sp)
	return false
}

func buildSpacer(s string, l int) string {
	sb := strings.Builder{}

	for i := 0; i < l-len(s); i++ {
		sb.WriteString(" ")
	}

	return sb.String()
}

func longestLength(sl []LabeledValue) int {
	longest := 0

	for _, v := range sl {
		if len(v.Label) > longest {
			longest = len(cleanTerminalCharsFromString(v.Label))
		}
	}

	return longest
}

func cleanTerminalCharsFromString(s string) string {
	r := s
	r = strings.ReplaceAll(r, TERMCYAN, "")
	r = strings.ReplaceAll(r, TERMCYANB, "")
	r = strings.ReplaceAll(r, TERMCYANREV, "")
	r = strings.ReplaceAll(r, TERMRED, "")
	r = strings.ReplaceAll(r, TERMREDB, "")
	r = strings.ReplaceAll(r, TERMREDREV, "")
	r = strings.ReplaceAll(r, TERMCLEAR, "")
	r = strings.ReplaceAll(r, TERMCLEARSCREEN, "")
	r = strings.ReplaceAll(r, TERMGREY, "")

	return r
}
