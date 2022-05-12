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
	"context"
	"encoding/json"
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
	"sync"

	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/cloudfunctions/v1"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
	"google.golang.org/api/run/v1"
	"google.golang.org/api/serviceusage/v1"
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
)

// ClearScreen will clear out a terminal screen.
func ClearScreen() {
	fmt.Println(TERMCLEARSCREEN)
}

var (
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
	Title            string            `json:"title"`
	Description      string            `json:"description"`
	Duration         int               `json:"duration"`
	Project          bool              `json:"collect_project"`
	ProjectNumber    bool              `json:"collect_project_number"`
	BillingAccount   bool              `json:"collect_billing_account"`
	Region           bool              `json:"collect_region"`
	RegionType       string            `json:"region_type"`
	RegionDefault    string            `json:"region_default"`
	Zone             bool              `json:"collect_zone"`
	HardSet          map[string]string `json:"hard_settings"`
	CustomSettings   []Custom          `json:"custom_settings"`
	RegistrarContact bool              `json:"collect_registrar_contact"`
}

// Custom represents a custom setting that we would like to collect from a user
// We will collect these settings from the user before continuing.
type Custom struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Default     string   `json:"default"`
	Value       string   `json:"value"`
	Options     []string `json:"options"`
}

// Collect will collect a value for a Custom from a user
func (c *Custom) Collect() error {
	fmt.Printf("%s%s: %s\n", TERMCYANB, c.Description, TERMCLEAR)

	if len(c.Options) > 0 {
		c.Value = listSelect(c.Options, c.Default)
		return nil
	}

	result := ""
	fmt.Printf("Enter value, or just [enter] for %s%s%s\n", TERMCYANB, c.Default, TERMCLEAR)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		result = text
		if len(text) == 0 {
			result = c.Default
		}
		c.Value = result
		break

	}

	return nil
}

type Customs []Custom

func (cs Customs) Get(name string) Custom {
	for _, v := range cs {
		if v.Name == name {
			return v
		}
	}

	return Custom{}
}

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
	fmt.Printf("%s\n", Divider)
}

// Process runs through all of the options in a config and collects all of the
// necessary data from users.
func (c Config) Process(s *Stack, output string) error {
	c.PrintHeader()
	var project, region, zone, projectnumber, billingaccount string
	var err error

	for i, v := range c.HardSet {
		s.AddSetting(i, v)
	}

	project = s.GetSetting("project_id")
	region = s.GetSetting("region")
	zone = s.GetSetting("zone")

	if c.Project && len(project) == 0 {
		project, err = ProjectManage()
		if err != nil {
			log.Fatalf(err.Error())
		}
		s.AddSetting("project_id", project)
	}

	if c.Region && len(region) == 0 {
		region, err = RegionManage(project, c.RegionType, c.RegionDefault)
		if err != nil {
			log.Fatalf(err.Error())
		}
		s.AddSetting("Region", region)
	}

	if c.Zone && len(zone) == 0 {
		zone, err = ZoneManage(project, region)
		if err != nil {
			log.Fatalf(err.Error())
		}
		s.AddSetting("zone", zone)
	}

	if c.ProjectNumber {
		projectnumber, err = ProjectNumber(project)
		if err != nil {
			log.Fatalf(err.Error())
		}
		s.AddSetting("project_number", projectnumber)
	}

	if c.RegistrarContact {
		if err = RegistratContactManage("contact.yaml"); err != nil {
			log.Fatalf(err.Error())
		}
	}

	if c.BillingAccount {

		ba, err := BillingAccountManage()
		if err != nil {
			log.Fatalf(err.Error())
		}
		billingaccount = ba
		s.AddSetting("billing_account", billingaccount)
	}

	for _, v := range c.CustomSettings {
		temp := s.GetSetting(v.Name)

		if len(temp) < 1 {
			if err := v.Collect(); err != nil {
				log.Fatalf("error getting custom value from user:  %s", err)
			}
			s.AddSetting(v.Name, v.Value)
		}

	}

	s.PrintSettings()
	s.TerraformFile(output)
	return nil
}

// NewConfig returns a Config object from a file read.
func NewConfig(content []byte) (Config, error) {
	result := Config{}

	if err := json.Unmarshal(content, &result); err != nil {
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

// ReadConfig reads in a Config from a json file.
func (s *Stack) ReadConfig(file, desc string) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("unable to read config file: %s", err)
	}
	config, err := NewConfig(content)
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
	result := ""

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
		if len(val) < 1 {
			continue
		}
		result = result + fmt.Sprintf("%s=\"%s\"\n", label, val)
	}

	return result
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

	longest := longestLengh(keys)

	fmt.Printf("%sProject Details %s \n", TERMCYANREV, TERMCLEAR)

	if s, ok := s.Settings["project_id"]; ok && len(s) > 0 {
		printSetting("project_id", s, longest)
	}

	if s, ok := s.Settings["project_number"]; ok {
		printSetting("project_number", s, longest)
	}

	ordered := []string{}
	for i, v := range s.Settings {
		if i == "project_id" || i == "project_number" {
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

// ProjectNumber will get the project_number for the input projectid
func ProjectNumber(id string) (string, error) {
	resp := ""
	ctx := context.Background()
	svc, err := cloudresourcemanager.NewService(ctx, opts)
	if err != nil {
		return resp, err
	}

	results, err := svc.Projects.Get(id).Do()
	if err != nil {
		return resp, err
	}

	resp = strconv.Itoa(int(results.ProjectNumber))

	return resp, nil
}

// Projects gets a list of the Projects a user has access to
func Projects() ([]string, error) {
	resp := []string{}
	ctx := context.Background()

	svc, err := cloudresourcemanager.NewService(ctx, opts)
	if err != nil {
		return resp, err
	}

	results, err := svc.Projects.List().Filter("lifecycleState=ACTIVE").Do()
	if err != nil {
		return resp, err
	}
	pwb, err := getBillingForProjects(results.Projects)
	if err != nil {
		return resp, err
	}

	sort.Slice(pwb, func(i, j int) bool {
		return strings.ToLower(pwb[i].Name) < strings.ToLower(pwb[j].Name)
	})

	for _, v := range pwb {
		if v.BillingEnabled {
			resp = append(resp, v.Name)
			continue
		}

		resp = append(resp, fmt.Sprintf("%s (Billing Disabled)", v.Name))

	}

	return resp, nil
}

type projectWithBilling struct {
	Name           string
	BillingEnabled bool
}

func getBillingForProjects(p []*cloudresourcemanager.Project) ([]projectWithBilling, error) {
	res := []projectWithBilling{}

	ctx := context.Background()
	svc, err := cloudbilling.NewService(ctx, opts)
	if err != nil {
		return res, err
	}
	var wg sync.WaitGroup
	wg.Add(len(p))

	for _, v := range p {
		go func(p *cloudresourcemanager.Project) {
			defer wg.Done()
			if p.LifecycleState == "ACTIVE" && p.Name != "" {
				proj := fmt.Sprintf("projects/%s", p.ProjectId)
				tmp, err := svc.Projects.GetBillingInfo(proj).Do()
				if err != nil {
					if strings.Contains(err.Error(), "The caller does not have permission") {
						fmt.Printf("project: %+v\n", p)
						return
					}

					fmt.Printf("error: %s\n", err)
					return
				}

				pwb := projectWithBilling{p.Name, tmp.BillingEnabled}
				res = append(res, pwb)
			}
		}(v)
	}
	wg.Wait()

	return res, nil
}

// billingAccounts gets a list of the billing accounts a user has access to
func billingAccounts() ([]*cloudbilling.BillingAccount, error) {
	resp := []*cloudbilling.BillingAccount{}
	ctx := context.Background()
	svc, err := cloudbilling.NewService(ctx, opts)
	if err != nil {
		return resp, err
	}

	results, err := svc.BillingAccounts.List().Do()
	if err != nil {
		return resp, err
	}

	return results.BillingAccounts, nil
}

// BillingAccountProjectAttach will enable billing in a given project
func BillingAccountProjectAttach(project, account string) error {
	retries := 10
	ctx := context.Background()
	svc, err := cloudbilling.NewService(ctx, opts)
	if err != nil {
		return err
	}

	ba := fmt.Sprintf("billingAccounts/%s", account)
	proj := fmt.Sprintf("projects/%s", project)

	cfg := cloudbilling.ProjectBillingInfo{
		BillingAccountName: ba,
	}

	var looperr error
	for i := 0; i < retries; i++ {
		_, looperr = svc.Projects.UpdateBillingInfo(proj, &cfg).Do()
		if looperr == nil {
			return nil
		}
		if strings.Contains(looperr.Error(), "User is not authorized to get billing info") {
			continue
		}
	}

	if strings.Contains(looperr.Error(), "Request contains an invalid argument") {
		return ErrorBillingInvalidAccount
	}

	if strings.Contains(looperr.Error(), "Not a valid billing account") {
		return ErrorBillingInvalidAccount
	}

	if strings.Contains(looperr.Error(), "The caller does not have permission") {
		return ErrorBillingNoPermission
	}

	return looperr
}

// BillingAccountManage either grabs the users only BillingAccount or
// presents a list of BillingAccounts to select from.
func BillingAccountManage() (string, error) {
	accounts, err := billingAccounts()
	if err != nil {
		return "", fmt.Errorf("could not get list of billing accounts: %s", err)
	}

	labeled := []string{}

	for _, v := range accounts {
		labeled = append(labeled, fmt.Sprintf("%s (%s)", v.DisplayName, strings.ReplaceAll(v.Name, "billingAccounts/", "")))
	}

	if len(accounts) == 1 {
		return extractAccount(labeled[0]), nil
	}

	result := listSelect(labeled, labeled[0])

	return extractAccount(result), nil
}

func extractAccount(s string) string {
	sl := strings.Split(s, "(")
	return strings.ReplaceAll(sl[1], ")", "")
}

// projectCreate does the work of actually creating a new project in your
// GCP account
func projectCreate(project string) error {
	ctx := context.Background()
	svc, err := cloudresourcemanager.NewService(ctx, opts)
	if err != nil {
		return err
	}

	proj := cloudresourcemanager.Project{Name: project, ProjectId: project}

	_, err = svc.Projects.Create(&proj).Do()
	if err != nil {
		if strings.Contains(err.Error(), "project_id must be at most 30 characters long") {
			return ErrorProjectCreateTooLong
		}
		if strings.Contains(err.Error(), "project_id contains invalid characters") {
			return ErrorProjectInvalidCharacters
		}

		return err
	}

	return nil
}

// projectDelete does the work of actually deleting an existing project in
// your GCP account
func projectDelete(project string) error {
	ctx := context.Background()
	svc, err := cloudresourcemanager.NewService(ctx, opts)
	if err != nil {
		return err
	}

	_, err = svc.Projects.Delete(project).Do()
	if err != nil {
		return err
	}

	return nil
}

// ProjectManage promps a user to select a project.
func ProjectManage() (string, error) {
	createString := "CREATE NEW PROJECT"
	project, err := ProjectID()
	if err != nil {
		return "", err
	}

	projects, err := Projects()
	if err != nil {
		return "", err
	}

	projdis := []string{}

	for _, v := range projects {
		if strings.Contains(v, "Billing Disabled") {
			v = fmt.Sprintf("%s%s%s", TERMGREY, v, TERMCLEAR)
		}
		projdis = append(projdis, v)
	}

	projdis = append([]string{createString}, projdis...)

	fmt.Printf("\n%sChoose a project to use for this application.%s\n\n", TERMCYANB, TERMCLEAR)
	fmt.Printf("%sNOTE:%s This app will make changes to the project. %s\n", TERMCYANREV, TERMCYAN, TERMCLEAR)
	fmt.Printf("While those changes are reverseable, it would be better to put it in a fresh new project. \n")

	project = listSelect(projdis, project)

	if project == createString {
		project, err = projectPrompt()
		if err != nil {
			return "", err
		}
	}

	return project, nil
}

// projectPrompt manages the interaction of creating a project, including prompts.
func projectPrompt() (string, error) {
	result := ""
	sec1 := NewSection("Creating the project")

	sec1.Open()
	fmt.Printf("%sPlease enter a new project name to create: %s\n", TERMCYANB, TERMCLEAR)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		if len(text) == 0 {
			fmt.Printf("%sPlease enter a new project name to create: %s\n", TERMCYANB, TERMCLEAR)
			continue
		}

		if err := projectCreate(text); err != nil {
			fmt.Printf("%sProject name could not be created, please choose another,%s\n", TERMREDREV, TERMCLEAR)
			continue
		}

		fmt.Printf("Project Created\n")
		result = text
		break

	}
	sec1.Close()

	sec2 := NewSection("Activating Billing for the project")
	sec2.Open()
	account, err := BillingAccountManage()
	if err != nil {
		return "", fmt.Errorf("could not determine proper billing account: %s ", err)
	}

	if err := BillingAccountProjectAttach(result, account); err != nil {
		return "", fmt.Errorf("could not link billing account: %s ", err)
	}
	sec2.Close()
	return result, nil
}

// regions will return a list of regions depending on product type
func regions(project, product string) ([]string, error) {
	switch product {
	case "compute":
		return regionsCompute(project)
	case "functions":
		return regionsFunctions(project)
	case "run":
		return regionsRun(project)
	}

	return []string{}, fmt.Errorf("invalid product requested: %s", product)
}

// regionsFunctions will return a list of regions for Cloud Functions
func regionsFunctions(project string) ([]string, error) {
	resp := []string{}

	ctx := context.Background()
	svc, err := cloudfunctions.NewService(ctx, opts)
	if err != nil {
		return resp, err
	}

	results, err := svc.Projects.Locations.List("projects/" + project).Do()
	if err != nil {
		return resp, err
	}

	for _, v := range results.Locations {
		resp = append(resp, v.LocationId)
	}

	sort.Strings(resp)

	return resp, nil
}

// regionsRun will return a list of regions for Cloud Run
func regionsRun(project string) ([]string, error) {
	resp := []string{}

	ctx := context.Background()
	svc, err := run.NewService(ctx, opts)
	if err != nil {
		return resp, err
	}

	results, err := svc.Projects.Locations.List("projects/" + project).Do()
	if err != nil {
		return resp, err
	}

	for _, v := range results.Locations {
		resp = append(resp, v.LocationId)
	}

	sort.Strings(resp)

	return resp, nil
}

// regionsCompute will return a list of regions for Compute Engine
func regionsCompute(project string) ([]string, error) {
	resp := []string{}

	ctx := context.Background()
	svc, err := compute.NewService(ctx, opts)
	if err != nil {
		return resp, err
	}

	results, err := svc.Regions.List(project).Do()
	if err != nil {
		return resp, err
	}

	for _, v := range results.Items {
		resp = append(resp, v.Name)
	}

	sort.Strings(resp)

	return resp, nil
}

// RegionManage promps a user to select a region.
func RegionManage(project, product, def string) (string, error) {
	fmt.Printf("Enabling service to poll...\n")
	service := "compute.googleapis.com"
	switch product {
	case "compute":
		service = "compute.googleapis.com"
	case "functions":
		service = "cloudfunctions.googleapis.com"
	case "run":
		service = "run.googleapis.com"
	}

	if err := ServiceEnable(project, service); err != nil {
		return "", fmt.Errorf("error activating service for polling: %s", err)
	}

	fmt.Printf("Polling for regions...\n")
	regions, err := regions(project, product)
	if err != nil {
		return "", err
	}
	fmt.Printf("%sChoose a valid region to use for this application. %s\n", TERMCYANB, TERMCLEAR)
	region := listSelect(regions, def)

	return region, nil
}

// zones will return a list of zones in a given region
func zones(project, region string) ([]string, error) {
	resp := []string{}

	ctx := context.Background()
	svc, err := compute.NewService(ctx, opts)
	if err != nil {
		return resp, err
	}

	filter := fmt.Sprintf("name=%s*", region)

	results, err := svc.Zones.List(project).Filter(filter).Do()
	if err != nil {
		return resp, err
	}

	for _, v := range results.Items {
		resp = append(resp, v.Name)
	}

	sort.Strings(resp)

	return resp, nil
}

// ZoneManage promps a user to select a zone.
func ZoneManage(project, region string) (string, error) {
	fmt.Printf("Enabling service to poll...\n")
	if err := ServiceEnable(project, "compute.googleapis.com"); err != nil {
		return "", fmt.Errorf("error activating service for polling: %s", err)
	}

	fmt.Printf("Polling for zones...\n")
	zones, err := zones(project, region)
	if err != nil {
		return "", err
	}

	fmt.Printf("%sChoose a valid zone to use for this application. %s\n", TERMCYANB, TERMCLEAR)
	zone := listSelect(zones, zones[0])
	return zone, nil
}

// listSelect presents a slice of strings as a list from which
// the user can select. It also highlights and preesnts behvaior for the
// default
func listSelect(sl []string, def string) string {
	itemCount := len(sl)
	halfcount := int(math.Ceil(float64(itemCount / 2)))
	width := longestLengh(sl)
	defaultExists := false

	if itemCount < 11 {
		for i, v := range sl {
			if ok := printWithDefault(i+1, width, v, def); ok {
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
			if ok := printWithDefault(i+1, width, v, def); ok {
				defaultExists = true
			}

			idx := i + halfcount + 1

			if idx > itemCount {
				fmt.Printf("\n")
				break
			}

			v2 := sl[idx-1]
			if ok := printWithDefault(idx, width, v2, def); ok {
				defaultExists = true
			}

			fmt.Printf("\n")
		}
	}

	answer := def
	reader := bufio.NewReader(os.Stdin)
	if defaultExists {
		fmt.Printf("Choose number from list, or just [enter] for %s%s%s\n", TERMCYANB, def, TERMCLEAR)
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

func printWithDefault(idx, width int, value, def string) bool {
	sp := buildSpacer(cleanTerminalCharsFromString(value), width)

	if value == def {
		fmt.Printf("%s%2d) %s %s%s", TERMCYANB, idx, value, sp, TERMCLEAR)
		return true
	}
	fmt.Printf("%2d) %s %s", idx, value, sp)
	return false
}

func buildSpacer(s string, l int) string {
	sb := strings.Builder{}

	for i := 0; i < l-len(s); i++ {
		sb.WriteString(" ")
	}

	return sb.String()
}

func longestLengh(sl []string) int {
	longest := 0

	for _, v := range sl {
		if len(v) > longest {
			longest = len(cleanTerminalCharsFromString(v))
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

// ServiceEnable enable a service in the selected project so that query calls
// to various lists will work.
func ServiceEnable(project, service string) error {
	ctx := context.Background()
	svc, err := serviceusage.NewService(ctx, opts)
	if err != nil {
		return err
	}
	s := fmt.Sprintf("projects/%s/services/%s", project, service)
	op, err := svc.Services.Enable(s, &serviceusage.EnableServiceRequest{}).Do()
	if err != nil {
		return fmt.Errorf("could not enable service: %s", err)
	}

	if !strings.Contains(string(op.Response), "ENABLED") {
		return ServiceEnable(project, service)
	}

	return nil
}

// DomainRegistrarContact represents the data required to register a domain
// with a public registrar.
type DomainRegistrarContact struct {
	Email         string
	Phone         string
	PostalAddress PostalAddress
}

// PostalAddress represents the mail address in a DomainRegistrarContact
type PostalAddress struct {
	RegionCode         string
	PostalCode         string
	AdministrativeArea string
	Locality           string
	AddressLines       []string
	Recipients         []string
}
