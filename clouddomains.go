package deploystack

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"text/template"

	domains "cloud.google.com/go/domains/apiv1beta1"
	"google.golang.org/api/iterator"
	domainspb "google.golang.org/genproto/googleapis/cloud/domains/v1beta1"
	"google.golang.org/genproto/googleapis/type/money"
	"google.golang.org/genproto/googleapis/type/postaladdress"
	"gopkg.in/yaml.v2"
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

// RegistratContactManage manages collecting domain registraton information
// from the user
func RegistratContactManage(file string) error {
	d := newContactData()

	fmt.Printf("%s\n", Divider)
	fmt.Printf("Domain registration requires some contact data. This process only asks for the absolutely mandatory ones. \n")
	fmt.Printf("The domain will be registered with user privacy enabled, so that your contact info will not be public. \n")
	fmt.Printf("This will create a file, so that you never have to do it again. \n")
	fmt.Printf("This file will only exist locally, or in your Cloud Shell environment.  \n")

	items := Customs{
		{Name: "email", Description: "Enter an email address", Default: "person@example.com"},
		{Name: "phone", Description: "Enter a phone number. (Please enter with country code - +1 555 555 5555 for US for example)", Default: "+14155551234"},
		{Name: "country", Description: "Enter a country code", Default: "US"},
		{Name: "postalcode", Description: "Enter a postal code", Default: "94502"},
		{Name: "state", Description: "Enter a state or administrative area", Default: "CA"},
		{Name: "city", Description: "Enter a city", Default: "San Francisco"},
		{Name: "address", Description: "Enter an address", Default: "345 Spear Street"},
		{Name: "name", Description: "Enter name", Default: "Googler"},
	}

	if err := items.Collect(); err != nil {
		return err
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
		return err
	}

	if err := os.WriteFile(file, []byte(yaml), 0o644); err != nil {
		return err
	}

	return nil
}

func DomainsSearch(project, domain string) ([]*domainspb.RegisterParameters, error) {
	ctx := context.Background()

	c, err := domains.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	req := &domainspb.SearchDomainsRequest{
		Query:    domain,
		Location: fmt.Sprintf("projects/%s/locations/global", project),
	}
	resp, err := c.SearchDomains(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.RegisterParameters, nil
}

func DomainIsAvailable(project, domain string) (*domainspb.RegisterParameters, error) {
	list, err := DomainsSearch(project, domain)
	if err != nil {
		return nil, err
	}
	for _, v := range list {
		if v.DomainName == domain {
			if v.Availability.String() == "AVAILABLE" {
				return v, nil
			}
			return nil, err
		}
	}

	return nil, err
}

func DomainsIsVerified(project, domain string) (bool, error) {
	ctx := context.Background()

	c, err := domains.NewClient(ctx)
	if err != nil {
		return false, err
	}
	defer c.Close()

	req := &domainspb.ListRegistrationsRequest{
		Filter: fmt.Sprintf("domainName=\"%s\"", domain),
		Parent: fmt.Sprintf("projects/%s/locations/global", project),
	}
	it := c.ListRegistrations(ctx, req)
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

func DomainRegister(project, domain, dnszone string, cost *money.Money, contact ContactData) error {
	ctx := context.Background()

	parent := fmt.Sprintf("projects/%s/locations/global", project)

	c, err := domains.NewClient(ctx)
	if err != nil {
		return err
	}
	defer c.Close()

	dnscontact, err := contact.DomainContact()
	if err != nil {
		return err
	}

	dnsprovider := domainspb.DnsSettings_CustomDns_{
		CustomDns: &domainspb.DnsSettings_CustomDns{
			NameServers: []string{
				"ns-cloud-e1.googledomains.com",
				"ns-cloud-e2.googledomains.com",
				"ns-cloud-e3.googledomains.com",
				"ns-cloud-e4.googledomains.com",
			},
		},
	}

	dnssettings := domainspb.DnsSettings{
		DnsProvider: &dnsprovider,
	}

	reg := domainspb.Registration{
		Name:            fmt.Sprintf("%s/registrations/%s", parent, domain),
		DomainName:      domain,
		DnsSettings:     &dnssettings,
		ContactSettings: &dnscontact,
	}

	req := &domainspb.RegisterDomainRequest{
		Registration: &reg,
		Parent:       parent,
		YearlyPrice:  cost,
	}

	c.RegisterDomain(ctx, req)

	return nil
}
