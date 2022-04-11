package deploystack

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/cloudfunctions/v1"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/run/v1"
	"google.golang.org/api/serviceusage/v1"
)

const (
	// TERMCYAN is the terminal code for cyan text
	TERMCYAN = "\033[0;36m"
	// TERMCYANB is the terminal code for bold cyan text
	TERMCYANB = "\033[1;36m"
	// TERMCYANREV is the terminal code for black on cyan text
	TERMCYANREV = "\033[36m"
	// TERMRED is the terminal code for red text
	TERMRED = "\033[0;31m"
	// TERMREDB is the terminal code for bold red text
	TERMREDB = "\033[1;31m"
	// TERMREDREV is the terminal code for black on red text
	TERMREDREV = "\033[41m"
	// TERMCLEAR is the terminal code for the clear out color text
	TERMCLEAR = "\033[0m"
)

var divider = ""

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
)

func init() {
	var err error
	divider, err = buildDivider()
	if err != nil {
		log.Fatal(err)
	}
}

func buildDivider() (string, error) {
	de := 80
	width := de
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

	var sb strings.Builder

	for i := 0; i < width; i++ {
		sb.WriteString("*")
	}
	return sb.String(), nil
}

// Config represents the settings this app will collect from a user. It should
// be in a json file. The idea is minimal programming has to be done to setup
// a DeployStack and export out a tfvars file for terraform part of solution.
type Config struct {
	Title          string            `json:"title"`
	Description    string            `json:"description"`
	Duration       int               `json:"duration"`
	Project        bool              `json:"collect_project"`
	ProjectNumber  bool              `json:"collect_project_number"`
	BillingAccount bool              `json:"collect_billing_account"`
	Region         bool              `json:"collect_region"`
	RegionType     string            `json:"region_type"`
	RegionDefault  string            `json:"region_default"`
	Zone           bool              `json:"collect_zone"`
	HardSet        map[string]string `json:"hard_settings"`
	CustomSettings []Custom          `json:"custom_settings"`
}

// Custom represents a custom setting that we would like to collect from a user
// We will collect these settings from the user before continuing.
type Custom struct {
	Name        string
	Description string
	Default     string
	Value       string
}

// Collect will collect a value for a Custom from a user
func (c *Custom) Collect() error {
	result := ""
	fmt.Printf("%s%s: %s\n", TERMCYANB, c.Description, TERMCLEAR)
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

// PrintHeader prints out the header for a DeployStack
func (c Config) PrintHeader() {
	fmt.Printf("%s\n", divider)
	fmt.Printf("%s%s%s\n", TERMCYANB, c.Title, TERMCLEAR)
	fmt.Printf("%s\n", c.Description)

	timestring := "minute"
	if c.Duration > 1 {
		timestring = "minutes"
	}

	fmt.Printf("It's going to take around %s%d %s%s\n", TERMCYAN, c.Duration, timestring, TERMCLEAR)
	fmt.Printf("%s\n", divider)
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

	if c.BillingAccount {

		ba, err := BillingAccountManage()
		if err != nil {
			log.Fatalf(err.Error())
		}
		billingaccount = ba
		s.AddSetting("billing_account", billingaccount)
	}

	for _, v := range c.CustomSettings {
		if err := v.Collect(); err != nil {
			log.Fatalf("error getting custom value from user:  %s", err)
		}
		s.AddSetting(v.Name, v.Value)
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

	if s, ok := s.Settings["project_id"]; ok {
		printSetting("project_id", s, longest)
	}

	if s, ok := s.Settings["project_number"]; ok {
		printSetting("project_number", s, longest)
	}

	for i, v := range s.Settings {
		if i == "project_id" || i == "project_number" {
			continue
		}
		if len(v) < 1 {
			continue
		}
		printSetting(i, v, longest)
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
	fmt.Printf("%s\n", divider)
	fmt.Printf("%s%s%s\n", TERMCYAN, s.Title, TERMCLEAR)
	fmt.Printf("%s\n", divider)
}

// Close prints out the footer for a Section.
func (s Section) Close() {
	fmt.Printf("%s\n", divider)
	fmt.Printf("%s%s - %sdone%s\n", TERMCYAN, s.Title, TERMCYANB, TERMCLEAR)
	fmt.Printf("%s\n", divider)
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
	svc, err := cloudresourcemanager.NewService(ctx)
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

// Projects gets a list of the projects a user has access to
func Projects() ([]string, error) {
	resp := []string{}
	ctx := context.Background()
	svc, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		return resp, err
	}

	results, err := svc.Projects.List().Filter("lifecycleState=ACTIVE").Do()
	if err != nil {
		return resp, err
	}

	for _, v := range results.Projects {
		if v.LifecycleState == "ACTIVE" {
			resp = append(resp, v.ProjectId)
		}
	}

	sort.Strings(resp)

	return resp, nil
}

// BillingAccounts gets a list of the billing accounts a user has access to
func BillingAccounts() ([]string, error) {
	resp := []string{}
	ctx := context.Background()
	svc, err := cloudbilling.NewService(ctx)
	if err != nil {
		return resp, err
	}

	results, err := svc.BillingAccounts.List().Do()
	if err != nil {
		return resp, err
	}

	for _, v := range results.BillingAccounts {
		resp = append(resp, strings.Replace(v.Name, "billingAccounts/", "", -1))
	}

	sort.Strings(resp)

	return resp, nil
}

// BillingAccountProjectAttach will enable billing in a given project
func BillingAccountProjectAttach(project, account string) error {
	retries := 10
	ctx := context.Background()
	svc, err := cloudbilling.NewService(ctx)
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

	fmt.Printf("LoopErr: %s\n", err)

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
	accounts, err := BillingAccounts()
	if err != nil {
		return "", fmt.Errorf("could not get list of billing accounts: %s", err)
	}

	if len(accounts) == 1 {
		return accounts[0], nil
	}

	return ListSelect(accounts, accounts[0]), nil
}

// ProjectCreate does the work of actually creating a new project in your
// GCP account
func ProjectCreate(project string) error {
	ctx := context.Background()
	svc, err := cloudresourcemanager.NewService(ctx)
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

// ProjectDelete does the work of actually deleting an existing project in
// your GCP account
func ProjectDelete(project string) error {
	ctx := context.Background()
	svc, err := cloudresourcemanager.NewService(ctx)
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

	projects = append([]string{createString}, projects...)

	fmt.Printf("\n%sChoose a project to use for this application.%s\n\n", TERMCYANB, TERMCLEAR)
	fmt.Printf("%sNOTE:%s This app will make changes to the project. %s\n", TERMCYANREV, TERMCYAN, TERMCLEAR)
	fmt.Printf("While those changes are reverseable, it would be better to put it in a fresh new project. \n")

	project = ListSelect(projects, project)

	if project == createString {
		project, err = ProjectPrompt()
		if err != nil {
			return "", err
		}
	}

	return project, nil
}

// ProjectPrompt manages the interaction of creating a project, including prompts.
func ProjectPrompt() (string, error) {
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

		if err := ProjectCreate(text); err != nil {
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

// Regions will return a list of regions depending on product type
func Regions(project, product string) ([]string, error) {
	switch product {
	case "compute":
		return RegionsCompute(project)
	case "functions":
		return RegionsFunctions(project)
	case "run":
		return RegionsRun(project)
	}

	return []string{}, fmt.Errorf("invalid product requested: %s", product)
}

// RegionsFunctions will return a list of regions for Cloud Functions
func RegionsFunctions(project string) ([]string, error) {
	resp := []string{}

	ctx := context.Background()
	svc, err := cloudfunctions.NewService(ctx)
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

// RegionsRun will return a list of regions for Cloud Run
func RegionsRun(project string) ([]string, error) {
	resp := []string{}

	ctx := context.Background()
	svc, err := run.NewService(ctx)
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

// RegionsCompute will return a list of regions for Compute Engine
func RegionsCompute(project string) ([]string, error) {
	resp := []string{}

	ctx := context.Background()
	svc, err := compute.NewService(ctx)
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
	regions, err := Regions(project, product)
	if err != nil {
		return "", err
	}
	fmt.Printf("%sChoose a valid region to use for this application. %s\n", TERMCYANB, TERMCLEAR)
	region := ListSelect(regions, def)

	return region, nil
}

// Zones will return a list of zones in a given region
func Zones(project, region string) ([]string, error) {
	resp := []string{}

	ctx := context.Background()
	svc, err := compute.NewService(ctx)
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
	zones, err := Zones(project, region)
	if err != nil {
		return "", err
	}

	fmt.Printf("%sChoose a valid zone to use for this application. %s\n", TERMCYANB, TERMCLEAR)
	zone := ListSelect(zones, zones[0])
	return zone, nil
}

// ListSelect presents a slice of strings as a list from which
// the user can select. It also highlights and preesnts behvaior for the
// default
func ListSelect(sl []string, def string) string {
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
	sp := buildSpacer(value, width)

	if value == def {
		fmt.Printf("%s%2d) %s %s %s", TERMCYANB, idx, value, sp, TERMCLEAR)
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
			longest = len(v)
		}
	}

	return longest
}

// ServiceEnable enable a service in the selected project so that query calls
// to various lists will work.
func ServiceEnable(project, service string) error {
	ctx := context.Background()
	svc, err := serviceusage.NewService(ctx)
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
