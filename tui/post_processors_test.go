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

package tui

import (
	"fmt"
	"reflect"
	"testing"

	"cloud.google.com/go/domains/apiv1beta1/domainspb"
	"github.com/GoogleCloudPlatform/deploystack/gcloud"
	tea "github.com/charmbracelet/bubbletea"
)

func TestCheckYesOrNo(t *testing.T) {
	tests := map[string]struct {
		in   string
		want bool
	}{
		"y":         {in: "y", want: true},
		"yes":       {in: "yes", want: true},
		"n":         {in: "n", want: true},
		"no":        {in: "no", want: true},
		"Y":         {in: "Y", want: true},
		"YES":       {in: "YES", want: true},
		"N":         {in: "N", want: true},
		"NO":        {in: "NO", want: true},
		"yEs":       {in: "yEs", want: true},
		"noT":       {in: "noT", want: false},
		"sadasdasf": {in: "sadasdasf", want: false},
		"ye":        {in: "ye", want: true},
		"4567":      {in: "4567", want: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := checkYesOrNo(tc.in)

			if tc.want != got {
				t.Fatalf("%s - want '%t' got '%t'", tc.in, tc.want, got)
			}
		})
	}
}

func TestValidateYesOrNo(t *testing.T) {
	tests := map[string]struct {
		in  string
		msg tea.Msg
	}{
		"y":         {in: "y", msg: successMsg{}},
		"yes":       {in: "yes", msg: successMsg{}},
		"n":         {in: "n", msg: successMsg{}},
		"no":        {in: "no", msg: successMsg{}},
		"Y":         {in: "Y", msg: successMsg{}},
		"YES":       {in: "YES", msg: successMsg{}},
		"N":         {in: "N", msg: successMsg{}},
		"NO":        {in: "NO", msg: successMsg{}},
		"yEs":       {in: "yEs", msg: successMsg{}},
		"noT":       {in: "noT", msg: errMsg{err: fmt.Errorf("Your answer '%s' is neither 'yes' nor 'no'", "noT")}},
		"sadasdasf": {in: "sadasdasf", msg: errMsg{err: fmt.Errorf("Your answer '%s' is neither 'yes' nor 'no'", "sadasdasf")}},
		"ye":        {in: "ye", msg: successMsg{}},
		"4567":      {in: "4567", msg: errMsg{err: fmt.Errorf("Your answer '%s' is neither 'yes' nor 'no'", "4567")}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			q := getTestQueue(appTitle, "test")
			cmd := validateYesOrNo(tc.in, &q)

			got := cmd()

			switch tc.msg.(type) {
			case successMsg:
				if tc.msg != got {
					t.Fatalf("%s - want: \n'%+v' \ngot: \n'%+v'", tc.in, tc.msg, got)
				}
			case errMsg:
				gotE := got.(errMsg)
				tcmsgE := tc.msg.(errMsg)

				if tcmsgE.err.Error() != gotE.err.Error() {
					t.Fatalf("want: \n'%+v' \ngot: \n'%+v'", tcmsgE.err.Error(), gotE.err.Error())
				}

			}
		})
	}
}

func TestValidatePhoneNumber(t *testing.T) {
	tests := map[string]struct {
		in  string
		msg tea.Msg
	}{
		"Good":  {"800 555 1234", successMsg{}},
		"Weird": {"d746fd83843", successMsg{}},
		"BAD":   {"dghdhdfuejfhfhfhrghfhfhdhgreh", errMsg{err: fmt.Errorf("Your answer '%s' is not a valid phone number. Please try again", "dghdhdfuejfhfhfhrghfhfhdhgreh")}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			q := getTestQueue(appTitle, "test")
			cmd := validatePhoneNumber(tc.in, &q)

			got := cmd()

			switch tc.msg.(type) {
			case successMsg:
				if tc.msg != got {
					t.Fatalf("%s - want: \n'%+v' \ngot: \n'%+v'", tc.in, tc.msg, got)
				}
			case errMsg:
				gotE := got.(errMsg)
				tcmsgE := tc.msg.(errMsg)

				if tcmsgE.err.Error() != gotE.err.Error() {
					t.Fatalf("want: \n'%+v' \ngot: \n'%+v'", tcmsgE.err.Error(), gotE.err.Error())
				}

			}
		})
	}
}

func TestMassagePhoneNumber(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
		err   error
	}{
		"Good":  {"800 555 1234", "+1.8005551234", nil},
		"Weird": {"d746fd83843", "+1.74683843", nil},
		"BAD":   {"dghdhdfuejfhfhfhrghfhfhdhgreh", "", ErrorCustomNotValidPhoneNumber},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := massagePhoneNumber(tc.input)
			if err != tc.err {
				t.Fatalf("expected: %v, got: %v", tc.err, err)
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestValidateInteger(t *testing.T) {
	tests := map[string]struct {
		in  string
		msg tea.Msg
	}{
		"1":    {in: "1", msg: successMsg{}},
		"12":   {in: "12", msg: successMsg{}},
		"1.4":  {in: "1.4", msg: errMsg{err: fmt.Errorf("Your answer '%s' not a valid integer", "1.4")}},
		"12n":  {in: "12n", msg: errMsg{err: fmt.Errorf("Your answer '%s' not a valid integer", "12n")}},
		"dsds": {in: "dsds", msg: errMsg{err: fmt.Errorf("Your answer '%s' not a valid integer", "dsds")}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			q := getTestQueue(appTitle, "test")
			cmd := validateInteger(tc.in, &q)

			got := cmd()

			switch tc.msg.(type) {
			case successMsg:
				if tc.msg != got {
					t.Fatalf("%s - want: \n'%+v' \ngot: \n'%+v'", tc.in, tc.msg, got)
				}
			case errMsg:
				gotE := got.(errMsg)
				tcmsgE := tc.msg.(errMsg)

				if tcmsgE.err.Error() != gotE.err.Error() {
					t.Fatalf("want: \n'%+v' \ngot: \n'%+v'", tcmsgE.err.Error(), gotE.err.Error())
				}

			}
		})
	}
}

func TestValidateDomain(t *testing.T) {
	tests := map[string]struct {
		in  string
		msg tea.Msg
	}{
		"example.com":  {in: "example.com", msg: errMsg{err: fmt.Errorf("validateDomain: error verifying domain: domain is not verified")}},
		"example2.com": {in: "example2.com", msg: errMsg{err: fmt.Errorf("validateDomain: not owned by requestor: %%!w(<nil>)")}},
		"example3.com": {in: "example3.com", msg: successMsg{}},
		"example4.com": {in: "example4.com", msg: successMsg{}},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			q := getTestQueue(appTitle, "test")
			cmd := validateDomain(tc.in, &q)

			got := cmd()

			switch tc.msg.(type) {
			case successMsg:
				if tc.msg != got {
					t.Fatalf("%s - want: \n'%+v' \ngot: \n'%+v'", tc.in, tc.msg, got)
				}
			case errMsg:
				gotE := got.(errMsg)
				tcmsgE := tc.msg.(errMsg)

				if tcmsgE.err.Error() != gotE.err.Error() {
					t.Fatalf("want: \n'%+v' \ngot: \n'%+v'", tcmsgE.err.Error(), gotE.err.Error())
				}

			}
		})
	}
}

func TestRegisterDomain(t *testing.T) {
	tests := map[string]struct {
		in   string
		msg  tea.Msg
		info *domainspb.RegisterParameters
	}{
		"example.com": {
			in:  "y",
			msg: errMsg{err: fmt.Errorf("registerDomain: error registering domain: domain is already owned. This should have been caught")},
			info: &domainspb.RegisterParameters{
				DomainName: "example.com",
			},
		},
		"example2.com": {
			in:  "y",
			msg: errMsg{err: fmt.Errorf("registerDomain: error registering domain: domain is already owned. This should have been caught")},
			info: &domainspb.RegisterParameters{
				DomainName: "example2.com",
			},
		},
		"example3.com": {
			in:  "y",
			msg: errMsg{err: fmt.Errorf("registerDomain: error registering domain: domain is cursed and cannot be obtained by mortals")},
			info: &domainspb.RegisterParameters{
				DomainName: "example3.com",
			},
		},

		"example4.com": {
			in:  "y",
			msg: successMsg{},
			info: &domainspb.RegisterParameters{
				DomainName: "example4.com",
			},
		},
		"example4.com_no": {
			in:  "n",
			msg: errMsg{err: fmt.Errorf("did not consent to being charged")},
			info: &domainspb.RegisterParameters{
				DomainName: "example4.com",
			},
		},
		"example4.com_garbage": {
			in:  "dasdasdasdsdas",
			msg: errMsg{err: fmt.Errorf("did not consent to being charged")},
			info: &domainspb.RegisterParameters{
				DomainName: "example4.com",
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			q := getTestQueue(appTitle, "test")
			q.Save("domainInfo", tc.info)
			q.stack.AddSetting("domain_consent", "y")
			cmd := registerDomain(tc.in, &q)

			got := cmd()

			switch tc.msg.(type) {
			case successMsg:
				if tc.msg != got {
					t.Fatalf("%s - want: \n'%+v' \ngot: \n'%+v'", tc.in, tc.msg, got)
				}
			case errMsg:
				gotE := got.(errMsg)
				tcmsgE := tc.msg.(errMsg)

				if tcmsgE.err.Error() != gotE.err.Error() {
					t.Fatalf("want: \n'%+v' \ngot: \n'%+v'", tcmsgE.err.Error(), gotE.err.Error())
				}

			}
		})
	}
}

func TestCreateProject(t *testing.T) {
	tests := map[string]struct {
		in  string
		msg tea.Msg
	}{
		"ds-tester-deploystack": {in: "ds-tester-deploystack", msg: errMsg{err: fmt.Errorf("createProject: could not create project: project_id already exists")}},
		"1234":                  {in: "1234", msg: errMsg{err: fmt.Errorf("createProject: could not create project: project_id contains too many characters, limit 30")}},
		"sa1234122132132143145315246754736573568765": {in: "1234", msg: errMsg{err: fmt.Errorf("createProject: could not create project: project_id contains too many characters, limit 30")}},
		"tp-never-used": {in: "tp-never-used", msg: successMsg{}},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			q := getTestQueue(appTitle, "test")
			cmd := createProject(tc.in, &q)

			got := cmd()

			switch tc.msg.(type) {
			case successMsg:
				if tc.msg != got {
					t.Fatalf("%s - want: \n'%+v' \ngot: \n'%+v'", tc.in, tc.msg, got)
				}
			case errMsg:
				gotE := got.(errMsg)
				tcmsgE := tc.msg.(errMsg)

				if tcmsgE.err.Error() != gotE.err.Error() {
					t.Fatalf("want: \n'%+v' \ngot: \n'%+v'", tcmsgE.err.Error(), gotE.err.Error())
				}

			}
		})
	}
}

func TestValidateGCEDefault(t *testing.T) {
	tests := map[string]struct {
		in       string
		msg      tea.Msg
		lenItems int
	}{
		"donotdefault": {in: "n", msg: successMsg{}, lenItems: 12},
		"default":      {in: "y", msg: successMsg{}, lenItems: 1},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			q := getTestQueue(appTitle, "test")
			newGCEInstance(&q)

			cmd := validateGCEDefault(tc.in, &q)

			got := cmd()

			switch tc.msg.(type) {
			case successMsg:
				if tc.msg != got {
					t.Fatalf("%s - want: \n'%+v' \ngot: \n'%+v'", tc.in, tc.msg, got)
				}
			case errMsg:
				gotE := got.(errMsg)
				tcmsgE := tc.msg.(errMsg)

				if tcmsgE.err.Error() != gotE.err.Error() {
					t.Fatalf("want: \n'%+v' \ngot: \n'%+v'", tcmsgE.err.Error(), gotE.err.Error())
				}
			}

			if tc.lenItems != len(q.models) {
				for _, v := range q.models {
					t.Logf("%s %s", v.getKey(), "")
				}

				t.Fatalf("number of models want: '%d' got: '%d'", tc.lenItems, len(q.models))
			}
		})
	}
}

func TestValidateGCEConfiguration(t *testing.T) {
	tests := map[string]struct {
		in    string
		msg   tea.Msg
		value string
	}{
		// "nowebserver":  {in: "n", msg: successMsg{unset: true}, value: ""},
		"yeswebserver": {in: "y", msg: successMsg{unset: true}, value: gcloud.HTTPServerTags},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			q := getTestQueue(appTitle, "test")

			q.stack.AddSetting("instance-webserver", tc.in)
			q.stack.AddSetting("gce-use-defaults", "n")
			q.stack.AddSetting("instance-webserver", "n")
			q.stack.AddSetting("instance-image-project", "n")
			q.stack.AddSetting("instance-machine-type-family", "n")
			q.stack.AddSetting("instance-image-family", "n")

			cmd := validateGCEConfiguration(tc.in, &q)

			got := cmd()

			switch tc.msg.(type) {
			case successMsg:
				if tc.msg != got {
					t.Fatalf("%s - want: \n'%+v' \ngot: \n'%+v'", tc.in, tc.msg, got)
				}
			case errMsg:
				gotE := got.(errMsg)
				tcmsgE := tc.msg.(errMsg)

				if tcmsgE.err.Error() != gotE.err.Error() {
					t.Fatalf("want: \n'%+v' \ngot: \n'%+v'", tcmsgE.err.Error(), gotE.err.Error())
				}
			}

			if tc.value != q.stack.GetSetting("instance-tags") {
				t.Fatalf("tags want: '%s' got: '%s'", tc.value, q.stack.GetSetting("instance-tags"))
			}
		})
	}
}
