// Package dstester is a collection of tools to make testing Terraform
// resources created for DeployStack easier.
package dstester

import (
	"fmt"
	"net/http"
	"os/exec"
	"reflect"
	"strings"
	"testing"
	"time"
)

var (
	// ErrorURLFail is thrown when a http poll fails
	ErrorURLFail = fmt.Errorf("the url did not return 200 after time allotted")
	// ErrorCheckFail is thrown when a check fails
	ErrorCheckFail = fmt.Errorf("there was an issue with the poll")
)

// Terraform is a resource for calling Terraform with a consistent set of
// variables.
type Terraform struct {
	Dir  string            // directory containing .tf files
	Vars map[string]string // collection of vars passed into terraform call
}

func (t Terraform) exec(command string, opt ...string) (string, error) {
	cmd := exec.Command("terraform")
	cmd.Args = append(cmd.Args, fmt.Sprintf("-chdir=%s", t.Dir))
	cmd.Args = append(cmd.Args, command)

	if command == "apply" || command == "destroy" {
		cmd.Args = append(cmd.Args, "-auto-approve")
		for i, v := range t.Vars {
			cmd.Args = append(cmd.Args, "-var")
			cmd.Args = append(cmd.Args, fmt.Sprintf("%s=%s", i, v))
		}
	}

	if command == "output" {
		for _, v := range opt {
			cmd.Args = append(cmd.Args, v)
		}
	}

	dat, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error: '%s'", string(dat))
	}
	return string(dat), nil
}

// Output extracts a terraform output variable from the terraform state
func (t Terraform) Output(variable string) (string, error) {
	return t.exec("output", variable)
}

// Init runs a terraform init command
func (t Terraform) Init() (string, error) {
	return t.exec("init")
}

// Apply runs a terraform apply command passing in the variables to the command
func (t Terraform) Apply() (string, error) {
	return t.exec("apply")
}

// Destroy runs a terraform destroy command
func (t Terraform) Destroy() (string, error) {
	return t.exec("destroy")
}

// InitApplyForTest runs terraform init and apply and can output extra
// information if debug is set to true
func (t Terraform) InitApplyForTest(test *testing.T, debug bool) {
	out, err := t.Init()
	if err != nil {
		test.Fatalf("expected no error, got: '%v'", err)
	}

	if debug {
		test.Logf("init: %s\n", out)
	}

	out2, err := t.Apply()
	if err != nil {
		test.Fatalf("expected no error, got: '%v'", err)
	}

	if debug {
		test.Logf("apply: %s\n", out2)
	}

	return
}

// DestroyForTest runs terraform destroy and can output extra information if
// debug is set to true
func (t Terraform) DestroyForTest(test *testing.T, debug bool) {
	out3, err := t.Destroy()
	if err != nil {
		test.Fatalf("expected no error, got: '%v'", err)
	}

	if debug {
		test.Logf("destroy: %s\n", out3)
	}

	return
}

// GCPResources is a list of resources
type GCPResources struct {
	Items   []GCPResource
	Project string
}

// Init runs through the items in the list and sets some prereqs
func (gs *GCPResources) Init() {
	for i, v := range gs.Items {
		if v.Project == "" {
			v.Project = gs.Project
			gs.Items[i] = v
		}
	}
}

// GCPResource represents a resource in Google Cloud that we want to check
// to see if it exists
type GCPResource struct {
	Product string
	Name    string
	Field   string
	Append  string
	Project string
}

func (g *GCPResource) desc() *exec.Cmd {
	cmd := exec.Command("gcloud")

	for _, v := range strings.Split(g.Product, " ") {
		cmd.Args = append(cmd.Args, v)
	}

	cmd.Args = append(cmd.Args, "describe", g.Name)

	if len(g.Append) > 0 {
		for _, v := range strings.Split(g.Append, " ") {
			cmd.Args = append(cmd.Args, v)
		}
	}

	if g.Field == "" {
		g.Field = "name"
	}

	cmd.Args = append(cmd.Args, fmt.Sprintf("--format=value(%s)", g.Field))

	if g.Project != "" {
		cmd.Args = append(cmd.Args, fmt.Sprintf("--project=%s", g.Project))
	}

	return cmd
}

// Describe runs a gcloud describe call to ensure that the resource exists
func (g *GCPResource) Describe() (string, error) {
	cmd := g.desc()

	dat, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error: '%s'", string(dat))
	}
	out := strings.TrimSpace(string(dat))

	return out, nil
}

// DescribeCommand gives the string for  a gcloud describe call to ensure that
// the resource exists
func (g *GCPResource) DescribeCommand() string {
	cmd := g.desc()

	for i, v := range cmd.Args {
		if strings.Contains(v, "--format=value") {
			v = fmt.Sprintf("--format=\"value(%s)\"", g.Field)
			cmd.Args[i] = v
			break
		}
	}

	return cmd.String()
}

// TextExistence runs through and tests for the existence of each of the
// GCPResources
func TextExistence(t *testing.T, items []GCPResource) {
	t.Logf("Testing for existence of GCP resources")
	testsExists := map[string]struct {
		input GCPResource
		want  string
	}{}
	for _, v := range items {
		testsExists[fmt.Sprintf("Test %s %s exists", v.Product, v.Name)] = struct {
			input GCPResource
			want  string
		}{v, v.Name}
	}

	for name, tc := range testsExists {
		t.Run(name, func(t *testing.T) {
			got, err := tc.input.Describe()
			if err != nil {

				if strings.Contains(err.Error(), "was not found") {
					debug := strings.ReplaceAll(tc.input.DescribeCommand(), which("gcloud"), "gcloud")
					t.Fatalf("expected item to exist, it did not\n To debug:\n %s", debug)
				}

				t.Fatalf("expected no error, got: '%v'", err)
			}

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: '%v', got: '%v'", tc.want, got)
			}
		})
	}
}

func which(s string) string {
	cmd := exec.Command("which")
	cmd.Args = append(cmd.Args, s)
	result, _ := cmd.Output()
	return strings.TrimSpace(string(result))
}

// TextNonExistence runs through and tests for the lack of existence of each of
// the GCPResources
func TextNonExistence(t *testing.T, items []GCPResource) {
	t.Logf("Testing for non-existence of GCP resources")
	testsNotExists := map[string]struct {
		input GCPResource
	}{}
	for _, v := range items {
		testsNotExists[fmt.Sprintf("Test %s %s does not exist", v.Product, v.Name)] = struct {
			input GCPResource
		}{v}
	}

	for name, tc := range testsNotExists {
		t.Run(name, func(t *testing.T) {
			_, err := tc.input.Describe()
			if err == nil {
				t.Fatalf("expected error, got no error")
			}
		})
	}
}

// TestChecks Cycles through the checks and runs them.
func TestChecks(t *testing.T, polls []Check, tf Terraform) {
	t.Logf("Testing polls to check resource readiness")
	testsPolls := map[string]struct {
		input Check
	}{}
	for _, v := range polls {
		testsPolls[fmt.Sprintf("Test poll %s for %s", v.Output, v.Type)] = struct {
			input Check
		}{v}
	}

	for name, tc := range testsPolls {
		t.Run(name, func(t *testing.T) {
			ok, err := tc.input.Do(tf)
			if err != nil {
				t.Fatalf("expected no error, got: '%v'", err)
			}
			if !ok {
				t.Fatalf("poll failed")
			}
		})
	}
}

// HTTPPoll polls a url attempts number of times with a delay of interval
// between attempts
func HTTPPoll(url, query string, interval, attempts int) (bool, error) {
	urlToUse := strings.ReplaceAll(url, "\"", "")
	urlToUse = strings.TrimSpace(urlToUse)
	client := http.Client{
		Timeout: 2 * time.Second,
	}

	for i := 0; i < attempts; i++ {
		resp, _ := client.Get(urlToUse)

		if resp != nil && resp.StatusCode == http.StatusOK {
			return true, nil
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}

	return false, fmt.Errorf("%s debug: %s", ErrorURLFail, urlToUse)
}

// CustomCheck allows for a custom bash command to be run as set in "custom"
func CustomCheck(command string) (bool, error) {
	sl := strings.Split(command, " ")

	cmd := exec.Command(sl[0])
	cmd.Args = append(cmd.Args, sl[1:]...)

	dat, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("error: '%s'", string(dat))
	}
	return true, nil
}

// Check represents an intersitial check to be run between terraform apply and
// terraform destroy
type Check struct {
	Output   string
	Type     string
	Attempts int
	Interval int
	Query    string
	Custom   string
}

// Do performs the operation that the check is supposed to run
func (c Check) Do(tf Terraform) (bool, error) {
	i := c.Interval
	a := c.Attempts

	if a == 0 {
		a = 50
	}

	if i == 0 {
		i = 5
	}

	val, err := tf.Output(c.Output)
	if err != nil {
		return false, err
	}
	switch c.Type {
	case "httpPoll":
		return HTTPPoll(val, c.Query, a, i)
	case "customCheck":
		return CustomCheck(c.Custom)
	}

	return false, ErrorCheckFail
}
