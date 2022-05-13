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
)

// DomainRegistrarContact represents the data required to register a domain
// with a public registrar.
type DomainRegistrarContact struct {
	Email         string
	Phone         string
	PostalAddress PostalAddress
}

// PostalAddress represents the mail address in a DomainRegistrarContact
type PostalAddress struct {
	RegionCode         string
	PostalCode         string
	AdministrativeArea string
	Locality           string
	AddressLines       []string
	Recipients         []string
}

// YAML outputs the content of this structure into the contact format needed for
// domain registration
func (d DomainRegistrarContact) YAML() (string, error) {
	yaml := `allContacts:
  email: '{{ .Email}}'
  phoneNumber: '{{.Phone}}'
  postalAddress: 
    regionCode: '{{ .PostalAddress.RegionCode}}'
    postalCode: '{{ .PostalAddress.PostalCode}}'
    administrativeArea: '{{ .PostalAddress.AdministrativeArea}}'
    locality: '{{ .PostalAddress.Locality}}'
    addressLines: [{{range $element := .PostalAddress.AddressLines}}'{{$element}}'{{end}}]
    recipients: [{{range $element := .PostalAddress.Recipients}}'{{$element}}'{{end}}]`

	t, err := template.New("yaml").Parse(yaml)
	if err != nil {
		return "", fmt.Errorf("error parsing the yaml template %s", err)
	}
	var tpl bytes.Buffer
	err = t.Execute(&tpl, d)
	if err != nil {
		return "", fmt.Errorf("error executing the yaml template %s", err)
	}

	return tpl.String(), nil
}

func newDomainRegistrarContact() DomainRegistrarContact {
	d := DomainRegistrarContact{}
	d.PostalAddress.AddressLines = []string{}
	d.PostalAddress.Recipients = []string{}
	return d
}

// RegistratContactManage manages collecting domain registraton information
// from the user
func RegistratContactManage(file string) error {
	d := newDomainRegistrarContact()

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

	d.Email = items.Get("email").Value
	d.Phone = items.Get("phone").Value
	d.PostalAddress.RegionCode = items.Get("country").Value
	d.PostalAddress.PostalCode = items.Get("postalcode").Value
	d.PostalAddress.AdministrativeArea = items.Get("state").Value
	d.PostalAddress.Locality = items.Get("city").Value
	d.PostalAddress.AddressLines = append(d.PostalAddress.AddressLines, items.Get("address").Value)
	d.PostalAddress.Recipients = append(d.PostalAddress.Recipients, items.Get("name").Value)

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
