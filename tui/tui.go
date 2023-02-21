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

// Package tui provides a BubbleTea powered tui for Deploystack. All rendering
// should happen within this package.
package tui

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/domains/apiv1beta1/domainspb"
	"github.com/GoogleCloudPlatform/deploystack"
	"github.com/GoogleCloudPlatform/deploystack/config"
	"github.com/GoogleCloudPlatform/deploystack/gcloud"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/compute/v1"
)

const (
	explainText           = "DeployStack will walk you through setting some options for the stack this solutions installs. Most questions have a default that you can choose by hitting the Enter key."
	appTitle              = "DeployStack"
	contactfile           = "contact.yaml.tmp"
	validationPhoneNumber = "phonenumber"
	validationYesOrNo     = "yesorno"
	validationInteger     = "integer"
)

var (
	spinnerType = spinner.Line
)

// ErrorCustomNotValidPhoneNumber is the error you get when you fail phone
// number validation.
var ErrorCustomNotValidPhoneNumber = fmt.Errorf("not a valid phone number")

type errMsg struct {
	err     error
	quit    bool
	usermsg string
	target  string
}

func (e errMsg) Error() string { return e.err.Error() }

type successMsg struct {
	msg   string
	unset bool
}

// UIClient interface encapsulates all of the calls to gcloud that one needs to
// make the TUI work
type UIClient interface {
	ProjectIDGet() (string, error)
	ProjectList() ([]gcloud.ProjectWithBilling, error)
	ProjectParentGet(project string) (*cloudresourcemanager.ResourceId, error)
	ProjectCreate(project, parent, parentType string) error
	ProjectNumberGet(id string) (string, error)
	ProjectIDSet(id string) error
	RegionList(project, product string) ([]string, error)
	ZoneList(project, region string) ([]string, error)
	DomainIsAvailable(project, domain string) (*domainspb.RegisterParameters, error)
	DomainIsVerified(project, domain string) (bool, error)
	DomainRegister(project string, domaininfo *domainspb.RegisterParameters, contact gcloud.ContactData) error
	ImageLatestGet(project, imageproject, imagefamily string) (string, error)
	MachineTypeList(project, zone string) (*compute.MachineTypeList, error)
	MachineTypeFamilyList(imgs *compute.MachineTypeList) gcloud.LabeledValues
	MachineTypeListByFamily(imgs *compute.MachineTypeList, family string) gcloud.LabeledValues
	ImageList(project, imageproject string) (*compute.ImageList, error)
	ImageTypeListByFamily(imgs *compute.ImageList, project, family string) gcloud.LabeledValues
	ImageFamilyList(imgs *compute.ImageList) gcloud.LabeledValues
	BillingAccountList() ([]*cloudbilling.BillingAccount, error)
	BillingAccountAttach(project, account string) error
}

// Run takes a deploystack configuration and walks someone through all of the
// input needed to run the eventual terraform
func Run(s *config.Stack, useMock bool) {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	defaultUserAgent := fmt.Sprintf("deploystack/%s", s.Config.Name)

	client := gcloud.NewClient(context.Background(), defaultUserAgent)
	q := NewQueue(s, &client)

	if useMock {
		q = NewQueue(s, GetMock(1))
	}

	q.Save("contact", deploystack.CheckForContact())
	q.InitializeUI()

	p := tea.NewProgram(q.Start(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		Fatal(err)
	}

	if q.Get("halted") != nil {
		Fatal(nil)
	}

	s.TerraformFile("terraform.tfvars")
	deploystack.CacheContact(q.Get("contact"))

	fmt.Print("\n\n")
	fmt.Print(titleStyle.Render("Deploystack"))
	fmt.Print("\n")
	fmt.Print(subTitleStyle.Render(s.Config.Title))
	fmt.Print("\n")
	fmt.Print(strong.Render("Installation will proceed with these settings"))
	fmt.Print(q.getSettings())
}

// Fatal stops processing of Deploystack and halts the calling process. All
// with an eye towards not processing in the shell script of things go wrong.
func Fatal(err error) {
	if err != nil {
		content := `There was an issue collecting the information it takes to run this application.
		You can try again by typing 'deploystack install' at the command prompt 
		If the issue persists, please report at: 
		https://github.com/GoogleCloudPlatform/deploystack/issues
		`

		errmsg := errMsg{
			err:     err,
			usermsg: content,
			quit:    true,
		}

		msg := errorAlert{errmsg}
		fmt.Print("\n\n")
		fmt.Println(titleStyle.Render("DeployStack"))
		fmt.Println(msg.Render())
	}

	os.Exit(1)
}
