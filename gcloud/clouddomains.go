package gcloud

import (
	"bytes"
	"context"
	"fmt"
	"os"
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

func (c *Client) getDomainsClient(project string) (*domains.Client, error) {
	var err error
	svc := c.services.domains

	if svc != nil {
		return svc, nil
	}

	if err := c.ServiceEnable(project, "domains.googleapis.com"); err != nil {
		return nil, fmt.Errorf("error activating service for polling: %s", err)
	}

	svc, err = domains.NewClient(c.ctx, c.opts)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve service: %w", err)
	}

	c.services.domains = svc

	return svc, nil
}

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

// DomainsSearch checks the Cloud Domain api for the input domain
func (c Client) DomainsSearch(project, domain string) ([]*domainspb.RegisterParameters, error) {
	svc, err := c.getDomainsClient(project)
	if err != nil {
		return nil, err
	}

	req := &domainspb.SearchDomainsRequest{
		Query:    domain,
		Location: fmt.Sprintf("projects/%s/locations/global", project),
	}
	resp, err := svc.SearchDomains(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return resp.RegisterParameters, nil
}

// DomainIsAvailable checks to see if a given domain is available for
// registration
func (c Client) DomainIsAvailable(project, domain string) (*domainspb.RegisterParameters, error) {
	list, err := c.DomainsSearch(project, domain)
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

// DomainIsVerified checks to see if a given domain belongs to this user
func (c Client) DomainIsVerified(project, domain string) (bool, error) {
	svc, err := c.getDomainsClient(project)
	if err != nil {
		return false, fmt.Errorf("cannot get domains client: %s", err)
	}

	req := &domainspb.ListRegistrationsRequest{
		Filter: fmt.Sprintf("domainName=\"%s\"", domain),
		Parent: fmt.Sprintf("projects/%s/locations/global", project),
	}
	it := svc.ListRegistrations(c.ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return false, fmt.Errorf("listing domains failed: %s", err)
		}

		if resp.DomainName == domain {
			return true, nil
		}
	}

	return false, nil
}

// DomainRegister handles registring a domain on behalf of the user.
func (c Client) DomainRegister(project string, domaininfo *domainspb.RegisterParameters, contact ContactData) error {
	parent := fmt.Sprintf("projects/%s/locations/global", project)

	svc, err := c.getDomainsClient(project)
	if err != nil {
		return err
	}

	dnscontact, err := contact.DomainContact()
	if err != nil {
		return err
	}

	req := &domainspb.RegisterDomainRequest{
		DomainNotices: domaininfo.DomainNotices,
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

	if _, err := svc.RegisterDomain(c.ctx, req); err != nil {
		return err
	}

	return nil
}
