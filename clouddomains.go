package deploystack

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/template"

	domains "cloud.google.com/go/domains/apiv1beta1"
	"google.golang.org/api/iterator"
	domainspb "google.golang.org/genproto/googleapis/cloud/domains/v1beta1"
	"google.golang.org/genproto/googleapis/type/postaladdress"
	"gopkg.in/yaml.v2"
)

var (
	// ErrorDomainUntenable is returned when a domain isn't available for registration, but
	// is also not owned by the user. It can't be used in this app
	ErrorDomainUntenable = fmt.Errorf("domain is not available, and not owned by attempting user")
	// ErrorDomainUserDeny is returned when an user declines the choice to purchase.
	ErrorDomainUserDeny = fmt.Errorf("user said no to buying the domain")
)

var domainsClient *domains.Client

func getDomainsClient(project string) (*domains.Client, error) {
	if domainsClient != nil {
		return domainsClient, nil
	}

	if err := ServiceEnable(project, "domains.googleapis.com"); err != nil {
		return nil, fmt.Errorf("error activating service for polling: %s", err)
	}

	ctx := context.Background()
	svc, err := domains.NewClient(ctx, opts)
	if err != nil {
		return nil, err
	}

	domainsClient = svc

	return svc, nil
}

var msgDomainRegisterContactExplanation = fmt.Sprintf(`Domain registration requires some contact data. This process only asks for the 
absolutely mandatory ones. The domain will be registered with user privacy 
enabled, so that your contact info will not be public. This will create a file, 
so that you never have to do it again. 
%sThis file will only exist locally, or in your Cloud Shell environment.%s  
`, TERMCYAN, TERMCLEAR)

var (
	msgDomainRegisterHeader      = fmt.Sprintf("%sManage Domain Registration %s", TERMCYANREV, TERMCLEAR)
	msgDomainContactHeader       = fmt.Sprintf("Managing Domain contact information")
	msgDomainContactFileWrite    = fmt.Sprintf("Your information was recorded and placed in a file for your future use.")
	msgDomainContactFileRead     = fmt.Sprintf("Your information was read from the local contact file.")
	msgDomainAvailablityHeader   = fmt.Sprintf("Managing Domain Availability")
	msgDomainAvailablityVerified = fmt.Sprintf("Domain is unavailable for purchase, but records show you are verified as the owner, so it can be used for this application.")
	msgDomainPurchase            = fmt.Sprintf("%sBuying a domain is not reversable, saying 'y' will incur a charge.%s", TERMREDB, TERMCLEAR)
	msgDomainRegisterSuccess     = fmt.Sprintf("%sDomain Registered%s. \n", TERMCYANB, TERMCLEAR)
)

// ContactData represents the structure that we need for Registrar Contact
// Data
type ContactData struct {
	AllContacts DomainRegistrarContact `yaml:"allContacts"`
}

// DomainRegistrarContact represents the data required to register a domain
// with a public registrar.
type DomainRegistrarContact struct {
	Email         string        `yaml:"email"`
	Phone         string        `yaml:"phoneNumber"`
	PostalAddress PostalAddress `yaml:"postalAddress"`
}

// PostalAddress represents the mail address in a DomainRegistrarContact
type PostalAddress struct {
	RegionCode         string   `yaml:"regionCode"`
	PostalCode         string   `yaml:"postalCode"`
	AdministrativeArea string   `yaml:"administrativeArea"`
	Locality           string   `yaml:"locality"`
	AddressLines       []string `yaml:"addressLines"`
	Recipients         []string `yaml:"recipients"`
}

// YAML outputs the content of this structure into the contact format needed for
// domain registration
func (c ContactData) YAML() (string, error) {
	yaml := `allContacts:
  email: '{{ .AllContacts.Email}}'
  phoneNumber: '{{.AllContacts.Phone}}'
  postalAddress: 
    regionCode: '{{ .AllContacts.PostalAddress.RegionCode}}'
    postalCode: '{{ .AllContacts.PostalAddress.PostalCode}}'
    administrativeArea: '{{ .AllContacts.PostalAddress.AdministrativeArea}}'
    locality: '{{ .AllContacts.PostalAddress.Locality}}'
    addressLines: [{{range $element := .AllContacts.PostalAddress.AddressLines}}'{{$element}}'{{end}}]
    recipients: [{{range $element := .AllContacts.PostalAddress.Recipients}}'{{$element}}'{{end}}]`

	t, err := template.New("yaml").Parse(yaml)
	if err != nil {
		return "", fmt.Errorf("error parsing the yaml template %s", err)
	}
	var tpl bytes.Buffer
	err = t.Execute(&tpl, c)
	if err != nil {
		return "", fmt.Errorf("error executing the yaml template %s", err)
	}

	return tpl.String(), nil
}

// DomainContact outputs a varible in the format that Domain Registration
// API needs.
func (c ContactData) DomainContact() (domainspb.ContactSettings, error) {
	dc := domainspb.ContactSettings{}

	pa := postaladdress.PostalAddress{
		RegionCode:         c.AllContacts.PostalAddress.RegionCode,
		PostalCode:         c.AllContacts.PostalAddress.PostalCode,
		AdministrativeArea: c.AllContacts.PostalAddress.AdministrativeArea,
		Locality:           c.AllContacts.PostalAddress.Locality,
		AddressLines:       c.AllContacts.PostalAddress.AddressLines,
		Recipients:         c.AllContacts.PostalAddress.Recipients,
	}

	all := domainspb.ContactSettings_Contact{
		Email:         c.AllContacts.Email,
		PhoneNumber:   c.AllContacts.Phone,
		PostalAddress: &pa,
	}

	dc.AdminContact = &all
	dc.RegistrantContact = &all
	dc.TechnicalContact = &all
	dc.Privacy = domainspb.ContactPrivacy_PRIVATE_CONTACT_DATA

	return dc, nil
}

func newContactData() ContactData {
	c := ContactData{}
	d := DomainRegistrarContact{}
	d.PostalAddress.AddressLines = []string{}
	d.PostalAddress.Recipients = []string{}
	c.AllContacts = d
	return c
}

func newContactDataFromFile(file string) (ContactData, error) {
	c := newContactData()

	dat, err := os.ReadFile(file)
	if err != nil {
		return c, err
	}

	err = yaml.Unmarshal(dat, &c)
	if err != nil {
		return c, err
	}

	return c, nil
}

// DomainManage walks a user through the porocess of collecting contact info and
// registering a domain.
func DomainManage(s *Stack) (string, error) {
	fmt.Println(Divider)
	fmt.Println(msgDomainRegisterHeader)
	fmt.Println(Divider)
	fmt.Println(msgDomainContactHeader)

	contactfile := "contact.yaml"
	contact := ContactData{}
	domain := ""
	project := s.GetSetting("project_id")

	item := Custom{Name: "domain", Description: "Enter a domain you wish to purchase and use for this application"}

	if err := item.Collect(); err != nil {
		return "", fmt.Errorf("trouble getting domain from keyboard: %s", err)
	}

	domain = item.Value

	fmt.Println(Divider)
	fmt.Println(msgDomainAvailablityHeader)

	domainInfo, err := domainIsAvailable(project, domain)
	if err != nil {
		return "", fmt.Errorf("error checking domain %s", err)
	}

	if domainInfo.Availability == domainspb.RegisterParameters_UNAVAILABLE {

		isVerified, err := domainsIsVerified(project, domain)
		if err != nil {
			return "", fmt.Errorf("error verifying domain %s", err)
		}
		// If not, fail and ask the user to repick
		if !isVerified {
			return "", ErrorDomainUntenable
		}

		fmt.Println(msgDomainAvailablityVerified)
		return domain, nil
	}

	if _, err := os.Stat(contactfile); errors.Is(err, os.ErrNotExist) {
		contact, err = RegistrarContactManage(contactfile)
		if err != nil {
			return "", fmt.Errorf("error getting contact data %s", err)
		}
	} else {
		contact, err = newContactDataFromFile(contactfile)
		if err != nil {
			contact, err = RegistrarContactManage(contactfile)
			if err != nil {
				return "", fmt.Errorf("error getting contact data %s", err)
			}
		} else {
			fmt.Println(msgDomainContactFileRead)
		}
	}

	fmt.Println(msgDomainPurchase)
	fmt.Printf("%sCost for %s%s%s%s will be %s%d%s%s%s.  Continue?%s\n",
		TERMREDB, TERMCYAN, domain, TERMCLEAR, TERMREDB,
		TERMCYAN,
		domainInfo.YearlyPrice.Units, domainInfo.YearlyPrice.CurrencyCode,
		TERMCLEAR, TERMREDB, TERMCLEAR)
	choice := Custom{Name: "choice", Description: "y or n?"}

	if err := choice.Collect(); err != nil {
		return "", fmt.Errorf("trouble getting domain from keyboard: %s", err)
	}

	if strings.ToLower(choice.Value) != "y" && strings.ToLower(choice.Value) != "yes" {
		return "", ErrorDomainUserDeny
	}

	if err := domainRegister(project, domainInfo, contact); err != nil {
		return "", fmt.Errorf("error registering domain %s", err)
	}
	fmt.Println(Divider)
	fmt.Println(msgDomainRegisterSuccess)
	return domain, nil
}

// RegistrarContactManage manages collecting domain registraton information
// from the user
func RegistrarContactManage(file string) (ContactData, error) {
	d := newContactData()

	fmt.Println(Divider)
	fmt.Printf(msgDomainRegisterContactExplanation)
	fmt.Println(Divider)

	items := Customs{
		{Name: "email", Description: "Enter an email address", Default: "person@example.com"},
		{Name: "phone", Description: "Enter a phone number. (Please enter with country code - +1 555 555 5555 for US for example)", Default: "+14155551234", Validation: "phonenumber"},
		{Name: "country", Description: "Enter a country code", Default: "US"},
		{Name: "postalcode", Description: "Enter a postal code", Default: "94502"},
		{Name: "state", Description: "Enter a state or administrative area", Default: "CA"},
		{Name: "city", Description: "Enter a city", Default: "San Francisco"},
		{Name: "address", Description: "Enter an address", Default: "345 Spear Street"},
		{Name: "name", Description: "Enter name", Default: "Googler"},
	}

	if err := items.Collect(); err != nil {
		return d, err
	}

	d.AllContacts.Email = items.Get("email").Value
	d.AllContacts.Phone = items.Get("phone").Value
	d.AllContacts.PostalAddress.RegionCode = items.Get("country").Value
	d.AllContacts.PostalAddress.PostalCode = items.Get("postalcode").Value
	d.AllContacts.PostalAddress.AdministrativeArea = items.Get("state").Value
	d.AllContacts.PostalAddress.Locality = items.Get("city").Value
	d.AllContacts.PostalAddress.AddressLines = append(d.AllContacts.PostalAddress.AddressLines, items.Get("address").Value)
	d.AllContacts.PostalAddress.Recipients = append(d.AllContacts.PostalAddress.Recipients, items.Get("name").Value)

	yaml, err := d.YAML()
	if err != nil {
		return d, err
	}

	if err := os.WriteFile(file, []byte(yaml), 0o644); err != nil {
		return d, err
	}
	fmt.Println(msgDomainContactFileWrite)

	return d, nil
}

func domainsSearch(project, domain string) ([]*domainspb.RegisterParameters, error) {
	c, err := getDomainsClient(project)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	req := &domainspb.SearchDomainsRequest{
		Query:    domain,
		Location: fmt.Sprintf("projects/%s/locations/global", project),
	}
	resp, err := c.SearchDomains(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return resp.RegisterParameters, nil
}

func domainIsAvailable(project, domain string) (*domainspb.RegisterParameters, error) {
	list, err := domainsSearch(project, domain)
	if err != nil {
		return nil, err
	}
	for _, v := range list {
		if v.DomainName == domain {
			return v, err
		}
	}

	return nil, err
}

func domainsIsVerified(project, domain string) (bool, error) {
	c, err := getDomainsClient(project)
	if err != nil {
		return false, err
	}
	defer c.Close()

	req := &domainspb.ListRegistrationsRequest{
		Filter: fmt.Sprintf("domainName=\"%s\"", domain),
		Parent: fmt.Sprintf("projects/%s/locations/global", project),
	}
	it := c.ListRegistrations(context.Background(), req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return false, err
		}

		if resp.DomainName == domain {
			return true, nil
		}
	}

	return false, nil
}

func domainRegister(project string, domaininfo *domainspb.RegisterParameters, contact ContactData) error {
	parent := fmt.Sprintf("projects/%s/locations/global", project)

	c, err := getDomainsClient(project)
	if err != nil {
		return err
	}
	defer c.Close()

	dnscontact, err := contact.DomainContact()
	if err != nil {
		return err
	}

	req := &domainspb.RegisterDomainRequest{
		Registration: &domainspb.Registration{
			Name:       fmt.Sprintf("%s/registrations/%s", parent, domaininfo.DomainName),
			DomainName: domaininfo.DomainName,
			DnsSettings: &domainspb.DnsSettings{
				DnsProvider: &domainspb.DnsSettings_CustomDns_{
					CustomDns: &domainspb.DnsSettings_CustomDns{
						NameServers: []string{
							"ns-cloud-e1.googledomains.com",
							"ns-cloud-e2.googledomains.com",
							"ns-cloud-e3.googledomains.com",
							"ns-cloud-e4.googledomains.com",
						},
					},
				},
			},
			ContactSettings: &dnscontact,
		},
		Parent:      parent,
		YearlyPrice: domaininfo.YearlyPrice,
	}

	if _, err := c.RegisterDomain(context.Background(), req); err != nil {
		return err
	}

	return nil
}
