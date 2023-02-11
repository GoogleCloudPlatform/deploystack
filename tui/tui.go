// Package tui provides a BubbleTea powered tui for Deploystack. All rendering
// should happen within this package.
package tui

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
	explainText = "DeployStack will walk you through setting some options for the stack this solutions installs. Most questions have a default that you can choose by hitting the Enter key."
	appTitle    = "DeployStack"
	contactfile = "contact.yaml.tmp"
)

const (
	validationPhoneNumber = "phonenumber"
	validationYesOrNo     = "yesorno"
	validationInteger     = "integer"
)

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
