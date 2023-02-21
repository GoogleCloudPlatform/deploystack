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
	gcloud    = ""
	terraform = ""
)

func init() {
	gcloud = which("gcloud")
	terraform = which("terraform")
}

func which(command string) string {
	cmd := exec.Command("which")
	cmd.Args = append(cmd.Args, command)
	result, _ := cmd.Output()

	return strings.TrimSpace(string(result))
}

// ErrorURLFail is thrown when a http poll fails
var ErrorURLFail = fmt.Errorf("the url did not return 200 after time allotted")

// ErrorCheckFail is thrown when a check fails
var ErrorCheckFail = fmt.Errorf("there was an issue with the poll")

// Terraform is a resource for calling Terraform with a consistent set of
// variables.
type Terraform struct {
	Dir  string            // directory containing .tf files
	Vars map[string]string // collection of vars passed into terraform call
}

func (tf Terraform) exec(command string, opt ...string) (string, error) {
	cmd := tf.cmd(command, opt...)

	dat, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error: '%s'", string(dat))
	}
	return string(dat), nil
}

func (tf Terraform) cmd(command string, opt ...string) *exec.Cmd {
	cmd := exec.Command("terraform")
	cmd.Args = append(cmd.Args, fmt.Sprintf("-chdir=%s", tf.Dir))
	cmd.Args = append(cmd.Args, command)

	if command == "apply" || command == "destroy" {
		cmd.Args = append(cmd.Args, "-auto-approve")
		for i, v := range tf.Vars {
			cmd.Args = append(cmd.Args, "-var")
			cmd.Args = append(cmd.Args, fmt.Sprintf("%s=%s", i, v))
		}
	}

	if command == "output" {
		for _, v := range opt {
			cmd.Args = append(cmd.Args, v)
		}
	}

	return cmd
}

func (tf Terraform) string(command string, opt ...string) string {
	cmd := tf.cmd(command, opt...)
	return cmd.String()
}

// Output extracts a terraform output variable from the terraform state
func (tf Terraform) Output(variable string) (string, error) {
	return tf.exec("output", variable)
}

// Init runs a terraform init command
func (tf Terraform) Init() (string, error) {
	return tf.exec("init")
}

// Apply runs a terraform apply command passing in the variables to the command
func (tf Terraform) Apply() (string, error) {
	return tf.exec("apply")
}

// Destroy runs a terraform destroy command
func (tf Terraform) Destroy() (string, error) {
	return tf.exec("destroy")
}

// InitApplyForTest runs terraform init and apply and can output extra
// information if debug is set to true
func (tf Terraform) InitApplyForTest(test *testing.T, debug bool) {
	out, err := tf.Init()
	if err != nil {
		test.Fatalf("expected no error, got: '%v'", err)
	}

	if debug {
		test.Logf("init: %s\n", out)
	}

	out2, err := tf.Apply()
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
func (tf Terraform) DestroyForTest(test *testing.T, debug bool) {
	out3, err := tf.Destroy()
	if err != nil {
		test.Fatalf("expected no error, got: '%v'", err)
	}

	if debug {
		test.Logf("destroy: %s\n", out3)
	}

	return
}

// Resources is a list of resources
type Resources struct {
	Items   []Resource
	Project string
}

// Init runs through the items in the list and sets some prereqs
func (gs *Resources) Init() {
	for i, v := range gs.Items {
		if v.Project == "" {
			v.Project = gs.Project
			gs.Items[i] = v
		}
	}
}

// Resource represents a resource in Google Cloud that we want to check
// to see if it exists. In most cases a gcloud command will be built using
// this content
type Resource struct {
	Product   string            // Portion of the gcloud command between gcloud and describe
	Name      string            // The name of the resource to describe
	Field     string            // The field to use in format directive. Defaults to 'name'
	Append    string            // A string of content to be added to the end of the command string
	Project   string            // The GCP project to use for the gcloud command
	Expected  string            // The exepcted value that will be checking for. Defaults to vaule of Name
	Arguments map[string]string // A set of key value pairs that will be added to the end of the gcloud command
}

func (g *Resource) desc() *exec.Cmd {
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

	for i, v := range g.Arguments {
		cmd.Args = append(cmd.Args, fmt.Sprintf("--%s", i), v)
	}

	if g.Project != "" {
		cmd.Args = append(cmd.Args, fmt.Sprintf("--project=%s", g.Project))
	}

	if g.Field == "" {
		g.Field = "name"
	}

	cmd.Args = append(cmd.Args, fmt.Sprintf("--format=value(%s)", g.Field))

	return cmd
}

func (g *Resource) delete() *exec.Cmd {
	cmd := exec.Command("gcloud")

	for _, v := range strings.Split(g.Product, " ") {
		cmd.Args = append(cmd.Args, v)
	}

	cmd.Args = append(cmd.Args, "delete", g.Name)

	if len(g.Append) > 0 {
		for _, v := range strings.Split(g.Append, " ") {
			cmd.Args = append(cmd.Args, v)
		}
	}

	for i, v := range g.Arguments {
		cmd.Args = append(cmd.Args, fmt.Sprintf("--%s", i), v)
	}

	if g.Project != "" {
		cmd.Args = append(cmd.Args, fmt.Sprintf("--project=%s", g.Project))
	}
	cmd.Args = append(cmd.Args, "-q")

	return cmd
}

// Exists runs a gcloud describe call to ensure that the resource exists
func (g *Resource) Exists() (string, error) {
	cmd := g.desc()

	dat, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error: '%s'", string(dat))
	}
	out := strings.TrimSpace(string(dat))

	return out, nil
}

// existsString gives the string for  a gcloud describe call to ensure that
// the resource exists
func (g *Resource) existsString() string {
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

func (g *Resource) deleteString() string {
	cmd := g.delete()
	return cmd.String()
}

// TextExistence runs through and tests for the existence of each of the
// GCPResources
func TextExistence(t *testing.T, items []Resource) {
	t.Logf("Testing for existence of GCP resources")
	testsExists := map[string]struct {
		input Resource
		want  string
	}{}
	for _, v := range items {
		if v.Expected == "" {
			v.Expected = v.Name
		}

		testsExists[fmt.Sprintf("Test %s %s exists", v.Product, v.Name)] = struct {
			input Resource
			want  string
		}{v, v.Expected}
	}

	for name, tc := range testsExists {
		t.Run(name, func(t *testing.T) {
			got, err := tc.input.Exists()
			if err != nil {
				debug := strings.ReplaceAll(tc.input.existsString(), gcloud, "gcloud")
				if strings.Contains(err.Error(), "was not found") {
					t.Fatalf("expected item to exist, it did not\n To debug:\n %s", debug)
				}

				t.Fatalf("expected no error, got: '%v' To debug:\n %s", err, debug)
			}

			if !reflect.DeepEqual(tc.want, got) {
				// artifact registry call leaks stuff into stderr
				if strings.Contains(got, "Repository Size") {
					if strings.Contains(got, tc.want) {
						return
					}
				}

				t.Fatalf("expected: '%v', got: '%v'", tc.want, got)
			}
		})
	}
}

// TextNonExistence runs through and tests for the lack of existence of each of
// the GCPResources
func TextNonExistence(t *testing.T, items []Resource) {
	t.Logf("Testing for non-existence of GCP resources")
	testsNotExists := map[string]struct {
		input Resource
	}{}
	for _, v := range items {
		testsNotExists[fmt.Sprintf("Test %s %s does not exist", v.Product, v.Name)] = struct {
			input Resource
		}{v}
	}

	for name, tc := range testsNotExists {
		t.Run(name, func(t *testing.T) {
			_, err := tc.input.Exists()
			if err == nil {
				t.Fatalf("expected error, got no error")
			}
		})
	}
}

// TestOperations Cycles through the operations and runs them.
func TestOperations(t *testing.T, operations Operations, tf Terraform) {
	if len(operations.Items) == 0 {
		return
	}

	t.Logf(operations.Label)
	testsPolls := map[string]struct {
		input Operation
	}{}
	for _, v := range operations.Items {
		testsPolls[fmt.Sprintf("Operation %s %s", operations.Key, v.Type)] = struct {
			input Operation
		}{v}
	}

	for name, tc := range testsPolls {
		t.Run(name, func(t *testing.T) {
			ok, err := tc.input.Do(tf)
			if err != nil {
				t.Fatalf("expected no error, got: '%v'", err)
			}
			if !ok {
				t.Fatalf("operation failed")
			}
		})
	}
}

// httpPoll polls a url attempts number of times with a delay of interval
// between attempts
func httpPoll(url, query string, interval, attempts int) (bool, error) {
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

// customCheck allows for a custom bash command to be run as set in "custom"
func customCheck(command string) (bool, error) {
	sl := strings.Split(command, " ")

	cmd := exec.Command(sl[0])
	cmd.Args = append(cmd.Args, sl[1:]...)

	dat, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("error: '%s'", string(dat))
	}
	return true, nil
}

// sleep is an operation that allows for delays in operations
func sleep(interval int) (bool, error) {
	time.Sleep(time.Duration(interval) * time.Second)
	return true, nil
}

// Operation represents an intersitial check to be run between terraform apply and
// terraform destroy
type Operation struct {
	Output   string
	Type     string
	Attempts int
	Interval int
	Query    string
	Custom   string
}

// Do performs the operation that the check is supposed to run
func (o Operation) Do(tf Terraform) (bool, error) {
	i := o.Interval
	a := o.Attempts

	if a == 0 {
		a = 50
	}

	if i == 0 {
		i = 5
	}

	val, err := tf.Output(o.Output)
	if err != nil {
		return false, err
	}
	switch o.Type {
	case "httpPoll":
		return httpPoll(val, o.Query, a, i)
	case "sleep":
		return sleep(i)
	case "customCheck":
		return customCheck(o.Custom)
	}

	return false, ErrorCheckFail
}

// Operations are a set of operations to perform and certain times in the lifecycle
// of a test
type Operations struct {
	Items []Operation
	Label string
	Key   string
}

// Add an operation to the list of operations
func (os *Operations) Add(o Operation) {
	os.Items = append(os.Items, o)
}

// OperationsSets are the whole collection of all of the pre and post operations
type OperationsSets map[string]Operations

// Add an operation to the underlying set of Operations.
func (os *OperationsSets) Add(target string, o Operation) {
	tmp := (*os)[target]
	tmp.Add(o)
	(*os)[target] = tmp
}

// NewOperationsSet returns the default set of operation sets
func NewOperationsSet() OperationsSets {
	ops := OperationsSets{}

	ops["preTest"] = Operations{
		Key:   "preTest",
		Label: "Operations to be run before any tests",
	}
	ops["preApply"] = Operations{
		Key:   "preApply",
		Label: "Operations to be run after terraform init and before terraform apply",
	}
	ops["postApply"] = Operations{
		Key:   "postApply",
		Label: "Operations to be run after terraform apply",
	}
	ops["preDestroy"] = Operations{
		Key:   "preDestroy",
		Label: "Operations to be run after terraform apply and before terraform destroy",
	}
	ops["postDestroy"] = Operations{
		Key:   "postDestroy",
		Label: "Operations to be run after terraform destroy",
	}
	ops["postTest"] = Operations{
		Key:   "postTest",
		Label: "Operations to be after any tests",
	}

	return ops
}

// TestStack runs the test for an entire Deploystack test given the right inputs
func TestStack(t *testing.T, tf Terraform, resources Resources, ops OperationsSets, debug bool) {
	TestOperations(t, ops["preTest"], tf)
	resources.Init()
	TestOperations(t, ops["preApply"], tf)
	tf.InitApplyForTest(t, debug)
	TestOperations(t, ops["postApply"], tf)
	TextExistence(t, resources.Items)
	TestOperations(t, ops["preDestroy"], tf)
	tf.DestroyForTest(t, debug)
	TestOperations(t, ops["postDestroy"], tf)
	TextNonExistence(t, resources.Items)
	TestOperations(t, ops["postTest"], tf)
}

// DebugCommands will spit out all of the executables that the framework calls
// under the covers for debugging purposes
func DebugCommands(t *testing.T, tf Terraform, resources Resources) {
	fmt.Printf("gcloud describe commands \n")
	for _, v := range resources.Items {
		output := strings.ReplaceAll(v.existsString(), gcloud, "gcloud")
		fmt.Printf("%s\n", output)
	}
	fmt.Println("")

	fmt.Printf("gcloud delete commands \n")
	for _, v := range resources.Items {
		cmd := v.delete().String()
		output := strings.ReplaceAll(cmd, gcloud, "gcloud")
		fmt.Printf("%s\n", output)
	}
	fmt.Println("")

	fmt.Printf("terraform commands \n")
	cmds := []string{"init", "apply", "destroy"}

	for _, v := range cmds {
		output := tf.string(v)
		output = strings.ReplaceAll(output, terraform, "terraform")
		fmt.Printf("%s\n", output)
	}
}

// Clean calls the deletion version of all of the gcloud commands to wipe out
// all of the resources in the project
func Clean(t *testing.T, tf Terraform, resources Resources) {
	for _, v := range resources.Items {
		cmd := v.delete()

		// Big issue is that storage buckets needs to be emptied before they are
		// deleted
		if strings.Contains(cmd.String(), "alpha storage buckets") {
			rm := exec.Command("gcloud")
			rm.Args = append(rm.Args, "alpha", "storage", "rm", "-r")

			for _, v := range cmd.Args {
				if strings.Contains(v, "gs://") {
					rm.Args = append(rm.Args, fmt.Sprintf("%s/**", v))
				}
			}

			dat, err := rm.CombinedOutput()
			if err != nil {
				t.Logf("bucket removal issue: %s", string(dat))
			}

		}

		dat, err := cmd.CombinedOutput()
		if err != nil {
			t.Logf("delete issue: %s", string(dat))
		}

	}
}
