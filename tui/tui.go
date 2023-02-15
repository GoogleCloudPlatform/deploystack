// Package tui provides a BubbleTea powered tui for Deploystack. All rendering
// should happen within this package.
package tui

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/domains/apiv1beta1/domainspb"
	"github.com/GoogleCloudPlatform/deploystack"
	"github.com/GoogleCloudPlatform/deploystack/gcloud"
	tea "github.com/charmbracelet/bubbletea"
	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/compute/v1"
)

// TODO: put all the text together in a better format
var (
	msgDomainRegisterHeader      = "Manage Domain Registration"
	msgDomainContactHeader       = "Managing Domain contact information"
	msgDomainContactFileWrite    = "Your information was recorded and placed in a file for your future use."
	msgDomainContactFileRead     = "Your information was read from the local contact file."
	msgDomainAvailablityHeader   = "Managing Domain Availability"
	msgDomainAvailablityVerified = "Domain is unavailable for purchase, but records show you are verified as the owner, so it can be used for this application."
	msgDomainPurchase            = "Buying a domain is not reversable, saying 'y' will incur a charge."
	msgDomainRegisterSuccess     = "Domain Registered."
	msgDomainOwnedNotByUser      = "Domain is owned already, by someone other than you. Please pick another domain"
)

const (
	explainText           = "DeployStack will walk you through setting some options for the stack this solutions installs. Most questions have a default that you can choose by hitting the Enter key."
	appTitle              = "DeployStack"
	contactfile           = "contact.yaml.tmp"
	validationPhoneNumber = "phonenumber"
	validationYesOrNo     = "yesorno"
	validationInteger     = "integer"
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

// Start takes a deploystack configuration and walks someone through all of the
// input needed to run the eventual terraform
func Start(s *deploystack.Stack, useMock bool) {
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

func Fatal(err error) {
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
	fmt.Printf("\n\n")
	fmt.Printf(titleStyle.Render("DeployStack"))
	fmt.Printf("\n")
	fmt.Println(msg.Render())
	os.Exit(1)
}
