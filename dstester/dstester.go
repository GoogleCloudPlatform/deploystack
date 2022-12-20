package dstester

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

type Terraform struct {
	Dir  string
	Vars map[string]string
}

func (t Terraform) exec(command string) (string, error) {
	cmd := exec.Command("terraform")
	cmd.Args = append(cmd.Args, fmt.Sprintf("-chdir=%s", t.Dir))
	cmd.Args = append(cmd.Args, command)

	if command == "apply" || command == "destroy" {
		for i, v := range t.Vars {
			cmd.Args = append(cmd.Args, "-auto-approve")
			cmd.Args = append(cmd.Args, "-var")
			cmd.Args = append(cmd.Args, fmt.Sprintf("%s=%s", i, v))
		}
	}

	dat, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error: '%s'", string(dat))
	}
	return string(dat), nil
}

func (t Terraform) Init() (string, error) {
	return t.exec("init")
}

func (t Terraform) Apply() (string, error) {
	return t.exec("apply")
}

func (t Terraform) Destroy() (string, error) {
	return t.exec("destroy")
}

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

type GCloudCMD struct {
	Product string
	Name    string
	Field   string
	Append  string
	Project string
}

func (g *GCloudCMD) Describe() (string, error) {
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

	dat, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error: '%s'", string(dat))
	}
	out := strings.TrimSpace(string(dat))

	return out, nil
}

func (g *GCloudCMD) DescribeCommand() string {
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

	cmd.Args = append(cmd.Args, fmt.Sprintf("--format=\"value(%s)\"", g.Field))

	if g.Project != "" {
		cmd.Args = append(cmd.Args, fmt.Sprintf("--project=%s", g.Project))
	}

	return cmd.String()
}

func log(name, s string) {
	os.WriteFile(name, []byte(s), 0644)
}

func TextExistence(items []GCloudCMD, t *testing.T) {
	testsExists := map[string]struct {
		input GCloudCMD
		want  string
	}{}
	for _, v := range items {
		testsExists[fmt.Sprintf("Test %s %s exists", v.Product, v.Name)] = struct {
			input GCloudCMD
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

func TextNonExistence(items []GCloudCMD, t *testing.T) {
	testsNotExists := map[string]struct {
		input GCloudCMD
	}{}
	for _, v := range items {
		testsNotExists[fmt.Sprintf("Test %s %s does not exist", v.Product, v.Name)] = struct {
			input GCloudCMD
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
