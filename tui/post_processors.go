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
	"strconv"
	"strings"

	"cloud.google.com/go/domains/apiv1beta1/domainspb"
	"github.com/GoogleCloudPlatform/deploystack/gcloud"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nyaruka/phonenumbers"
)

func processProjectSelection(key string, q *Queue) tea.Cmd {
	return func() tea.Msg {
		if key != "" {

			if errMsg := handleProjectNumber(key, q); errMsg != nil {
				return errMsg
			}

			if err := q.client.ProjectIDSet(key); err != nil {
				return errMsg{err: err}
			}

			q.Save("currentProject", key)

			creator := q.currentKey() + projNewSuffix
			billing := q.currentKey() + billNewSuffix

			q.removeModel(creator)
			q.removeModel(billing)

			return successMsg{}
		}

		return successMsg{}
	}
}

func handleProjectNumber(projectID string, q *Queue) tea.Msg {
	if q.stack.Config.ProjectNumber {
		projectnumber, err := q.client.ProjectNumberGet(projectID)
		if err != nil {
			return errMsg{err: err}
		}
		q.stack.AddSetting("project_number", projectnumber)
	}
	return nil
}

func createProject(project string, q *Queue) tea.Cmd {
	return func() tea.Msg {

		currentProject := q.Get("currentProject").(string)

		if currentProject == "" {
			tmp, err := q.client.ProjectList()
			if err != nil || len(tmp) == 0 || tmp[0].ID == "" {
				return errMsg{err: fmt.Errorf("createProject: could not determine an alternate project for parent detection: %w ", err)}
			}
			currentProject = tmp[0].ID
		}

		parent, err := q.client.ProjectParentGet(currentProject)
		if err != nil {
			return errMsg{err: fmt.Errorf("createProject: could not determine proper parent for project: %w ", err)}
		}

		if err := q.client.ProjectCreate(project, parent.Id, parent.Type); err != nil {
			return errMsg{err: fmt.Errorf("createProject: could not create project: %w", err)}
		}

		if errMsg := handleProjectNumber(project, q); errMsg != nil {
			return errMsg
		}
		if err := q.client.ProjectIDSet(project); err != nil {
			return errMsg{err: err}
		}

		q.Save("currentProject", project)

		return successMsg{}
	}
}

func attachBilling(ba string, q *Queue) tea.Cmd {
	return func() tea.Msg {
		baclean := strings.ReplaceAll(ba, "billingAccounts/", "")
		key := strings.ReplaceAll(q.currentKey(), billNewSuffix, "")
		project := q.stack.GetSetting(key)

		if err := q.client.BillingAccountAttach(project, baclean); err != nil {
			return errMsg{err: fmt.Errorf("attachBilling: could not attach billing to project: %w", err)}
		}

		// If this is one of those billing for project form, let's skip
		// adding it to the stack settings
		if strings.Contains(q.currentKey(), billNewSuffix) {
			return successMsg{unset: true}
		}

		return successMsg{}
	}
}

func validateDomain(domain string, q *Queue) tea.Cmd {
	return func() tea.Msg {
		project := q.Get("currentProject").(string)

		domainInfo, err := q.client.DomainIsAvailable(project, domain)
		if err != nil {
			return errMsg{err: fmt.Errorf("validateDomain: error checking domain availability %w", err)}
		}

		q.Save("domainInfo", domainInfo)
		q.Save("domain", domain)

		if domainInfo.Availability == domainspb.RegisterParameters_UNAVAILABLE {
			isVerified, err := q.client.DomainIsVerified(project, domain)
			if err != nil {
				return errMsg{
					usermsg: "Trying to validate that you own this domain failed due to an error",
					err:     fmt.Errorf("validateDomain: error verifying domain: %s", err),
				}
			}
			if !isVerified {
				return errMsg{
					usermsg: "Domain is owned by someone other than the requestor",
					err:     fmt.Errorf("validateDomain: not owned by requestor: %w", err),
				}
			}

			return successMsg{}
		}
		return successMsg{}
	}
}

func registerDomain(consent string, q *Queue) tea.Cmd {
	return func() tea.Msg {
		userMsg := "There was a problem registering the domain."
		c := strings.ToLower(consent)
		if c != "y" && c != "yes" {
			q.stack.DeleteSetting("domain_consent")
			return errMsg{
				usermsg: userMsg,
				err:     fmt.Errorf("did not consent to being charged"),
				target:  "quit",
			}
		}

		d := gcloud.ContactData{}

		contact := q.Get("contact")

		if contact != nil {
			tmp := contact.(gcloud.ContactData)
			if tmp.AllContacts.Email != "" {
				d = tmp
			}
		}
		if d.AllContacts.Email == "" {
			d = gcloud.ContactData{
				AllContacts: gcloud.DomainRegistrarContact{
					Email: q.stack.GetSetting("domain_email"),
					Phone: q.stack.GetSetting("domain_phone"),
					PostalAddress: gcloud.PostalAddress{
						RegionCode: q.stack.GetSetting(
							"domain_country",
						),
						PostalCode: q.stack.GetSetting(
							"domain_postalcode",
						),
						AdministrativeArea: q.stack.GetSetting("domain_state"),
						Locality:           q.stack.GetSetting("domain_city"),
						AddressLines: []string{
							q.stack.GetSetting("domain_address"),
						},
						Recipients: []string{
							q.stack.GetSetting("domain_name"),
						},
					},
				},
			}
		}

		q.Save("contact", d)

		raw := q.Get("domainInfo")
		domainInfo := raw.(*domainspb.RegisterParameters)

		currentProject := q.Get("currentProject").(string)

		err := q.client.DomainRegister(currentProject, domainInfo, d)
		if err != nil {
			q.stack.AddSetting("domain_consent", "")
			return errMsg{
				usermsg: userMsg,
				err:     fmt.Errorf("registerDomain: error registering domain: %w", err),
				target:  "domain",
			}
		}

		for i := range q.stack.Settings {
			if strings.Contains(i, "domain_") {
				q.stack.DeleteSetting(i)
			}
		}

		return successMsg{}
	}
}

func validateInteger(input string, q *Queue) tea.Cmd {
	return func() tea.Msg {
		_, err := strconv.Atoi(input)
		if err != nil {
			return errMsg{err: fmt.Errorf("Your answer '%s' not a valid integer", input)}
		}
		return successMsg{}
	}
}

func checkYesOrNo(input string) bool {
	text := strings.TrimSpace(strings.ToLower(input))
	yesList := " yes y "
	noList := " no n "

	return strings.Contains(yesList+noList, text)
}

func validateYesOrNo(input string, q *Queue) tea.Cmd {
	return func() tea.Msg {
		text := strings.TrimSpace(strings.ToLower(input))

		if !checkYesOrNo(text) {
			return errMsg{err: fmt.Errorf("Your answer '%s' is neither 'yes' nor 'no'", input)}
		}

		return successMsg{}
	}
}

func validatePhoneNumber(input string, q *Queue) tea.Cmd {
	return func() tea.Msg {
		_, err := massagePhoneNumber(input)
		if err != nil {
			return errMsg{err: fmt.Errorf("Your answer '%s' is not a valid phone number. Please try again", input)}
		}

		return successMsg{}
	}
}

func massagePhoneNumber(s string) (string, error) {
	num, err := phonenumbers.Parse(s, "US")
	if err != nil {
		return "", ErrorCustomNotValidPhoneNumber
	}
	result := phonenumbers.Format(num, phonenumbers.INTERNATIONAL)
	result = strings.Replace(result, " ", ".", 1)
	result = strings.ReplaceAll(result, "-", "")
	result = strings.ReplaceAll(result, " ", "")

	return result, nil
}

// TODO: see if you can test these error conditions
func validateGCEDefault(input string, q *Queue) tea.Cmd {
	return func() tea.Msg {
		text := strings.TrimSpace(strings.ToLower(input))

		if !checkYesOrNo(text) {
			return errMsg{err: fmt.Errorf("Your answer '%s' is neither 'yes' nor 'no'", input)}
		}

		if string(text[0]) == "n" {
			return successMsg{}
		}

		project := q.stack.GetSetting("project-id")
		basename := q.stack.GetSetting("basename")

		defaultImage, err := q.client.ImageLatestGet(project, gcloud.DefaultImageProject, gcloud.DefaultImageFamily)
		if err != nil {
			return errMsg{err: fmt.Errorf("validateGCEDefault: could not get DefaultImage deafult")}
		}

		defaultConfig := map[string]string{
			"instance-image":        defaultImage,
			"instance-disksize":     gcloud.DefaultDiskSize,
			"instance-disktype":     gcloud.DefaultDiskType,
			"instance-tags":         gcloud.HTTPServerTags,
			"instance-name":         fmt.Sprintf("%s-instance", basename),
			"region":                gcloud.DefaultRegion,
			"zone":                  gcloud.DefaultZone,
			"instance-machine-type": gcloud.DefaultInstanceType,
		}

		for i, v := range defaultConfig {
			q.stack.AddSetting(i, v)
		}
		q.removeModel("instance-webserver")
		q.removeModel("instance-image-project")
		q.removeModel("instance-machine-type-family")
		q.removeModel("instance-image")
		q.removeModel("instance-image-type")
		q.removeModel("instance-disksize")
		q.removeModel("instance-disktype")
		q.removeModel("instance-tags")
		q.removeModel("instance-name")
		q.removeModel("instance-machine-type")
		q.removeModel("region")
		q.removeModel("zone")
		q.removeModel("instance-image-family")

		return successMsg{}
	}
}

func validateGCEConfiguration(input string, q *Queue) tea.Cmd {
	return func() tea.Msg {
		q.stack.AddSetting("instance-tags", "")
		instanceWebserver := q.stack.GetSetting("instance-webserver")

		if instanceWebserver == "y" || input == "y" {
			q.stack.AddSetting("instance-tags", gcloud.HTTPServerTags)
		}
		q.stack.DeleteSetting("gce-use-defaults")
		q.stack.DeleteSetting("instance-webserver")
		q.stack.DeleteSetting("instance-image-project")
		q.stack.DeleteSetting("instance-machine-type-family")
		q.stack.DeleteSetting("instance-image-family")
		return successMsg{unset: true}
	}
}
