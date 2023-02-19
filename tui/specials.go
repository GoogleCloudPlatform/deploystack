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
	"strings"

	"cloud.google.com/go/domains/apiv1beta1/domainspb"
	"github.com/GoogleCloudPlatform/deploystack"
	"github.com/GoogleCloudPlatform/deploystack/gcloud"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	projNewSuffix = "_new_project_creator"
	billNewSuffix = "_new_billing_selector"
)

func newProjectCreator(key string) textInput {
	r := newTextInput("Create New Project",
		"",
		key,
		"Checking if project can be created",
	)
	r.addPostProcessor(createProject)

	r.addContent("Project IDs are immutable and can be set only during project ")
	r.addContent("creation. They must start with a lowercase letter and can have ")
	r.addContent("lowercase ASCII letters, digits or hyphens. ")
	r.addContent("Project IDs must be between 6 and 30 characters. ")
	r.addContent("\n\n")
	r.addContent(textInputDefaultStyle.Render("Please enter a new project name to create:"))
	return r
}

func newProjectSelector(key, listLabel, currentProject string, preProcessor tea.Cmd) picker {

	result := newPicker(listLabel, "Retrieving Projects", key, currentProject, preProcessor)
	create := item{"Create New Project", ""}
	result.list.InsertItem(0, create)
	result.addPostProcessor(processProjectSelection)
	return result
}

func newBillingSelector(key string, preProcessor tea.Cmd, postProccessor func(string, *Queue) tea.Cmd) picker {
	result := newPicker("Choose a billing account to use for this project", "Retrieving Billing Accounts", key, "", preProcessor)
	result.postProcessor = postProccessor
	return result
}

func newYesOrNo(q *Queue, listLabel, key string, defaultNo bool, postProcessor func(string, *Queue) tea.Cmd) picker {
	p := newPicker(listLabel, "", key, "", getYesOrNo(q))
	p.list.SetShowStatusBar(false)
	p.list.SetShowFilter(false)
	p.list.SetShowHelp(false)
	p.addPostProcessor(postProcessor)
	p.preProcessor = getNoOrYes(q)

	return p
}

func newCustom(c deploystack.Custom) QueueModel {
	r := newTextInput(c.Description,
		c.Default,
		c.Name,
		"validating",
	)

	switch c.Validation {
	case validationPhoneNumber:
		r.spinnerLabel = "Validating phone number"
		r.addPostProcessor(validatePhoneNumber)
	case validationYesOrNo:
		r.spinnerLabel = "Validating yes or no"
		r.addPostProcessor(validateYesOrNo)
	case validationInteger:
		r.spinnerLabel = "Validating integer"
		r.addPostProcessor(validateInteger)
	}

	return &r
}

func newDomain(q *Queue) {
	contact := gcloud.ContactData{}

	t := newTextInput(
		"Enter a domain you wish to purchase and use for this application",
		"",
		"domain",
		"Checking Domain Availability",
	)
	t.postProcessor = validateDomain
	q.add(&t)

	items := []struct {
		Name         string
		Description  string
		DefaultValue string
		Validator    func(string, *Queue) tea.Cmd
	}{
		{
			Name:         "domain_email",
			Description:  "Enter an email address",
			DefaultValue: "person@example.com",
		},

		{
			Name:         "domain_phone",
			Description:  "Enter a phone number. (Please enter with country code - +1 555 555 5555 for US for example)",
			DefaultValue: "+14155551234",
			Validator:    validatePhoneNumber,
		},

		{
			Name:         "domain_country",
			Description:  "Enter a country code",
			DefaultValue: "US",
		},

		{
			Name:         "domain_postalcode",
			Description:  "Enter a postal code",
			DefaultValue: "94502",
		},

		{
			Name:         "domain_state",
			Description:  "Enter a state or administrative area",
			DefaultValue: "CA",
		},

		{
			Name:         "domain_city",
			Description:  "Enter a city",
			DefaultValue: "San Francisco",
		},

		{
			Name:         "domain_address",
			Description:  "Enter an address",
			DefaultValue: "345 Spear Street",
		},

		{
			Name:         "domain_name",
			Description:  "Enter name",
			DefaultValue: "Googler",
		},
	}

	tmp := q.Get("contact")
	switch v := tmp.(type) {
	case gcloud.ContactData:
		contact = v
	default:
		contact = gcloud.ContactData{}
	}

	if contact.AllContacts.Email == "" {
		for _, v := range items {
			t := newTextInput(v.Description, v.DefaultValue, v.Name, "")
			q.add(&t)
		}
	}

	f := func(q *Queue) {
		domain := q.Get("domain").(string)
		info := q.Get("domainInfo").(*domainspb.RegisterParameters)

		if info.YearlyPrice != nil {
			msg := fmt.Sprintf(
				"Cost for %s will be %s.  %s",
				domain,
				purchaseStyle.Render(
					fmt.Sprintf(
						"%d%s",
						info.YearlyPrice.Units,
						info.YearlyPrice.CurrencyCode,
					),
				),
				textStyle.Render("Continue?"),
			)
			p := q.models[q.current]
			p.clearContent()
			p.addContent(msg)
		}
	}

	dy := newYesOrNo(
		q,
		"Buying a domain is not reversable, saying 'y' will incur a charge.",
		"domain_consent",
		false,
		nil,
	)
	dy.spinnerLabel = "Attempting to register domain"
	dy.addPreView(f)
	dy.addPostProcessor(registerDomain)
	q.add(&dy)
}

func newCustomPages(q *Queue) {
	for _, v := range q.stack.Config.CustomSettings {
		temp := q.stack.GetSetting(v.Name)

		if len(v.Options) > 0 {

			items := []list.Item{}
			for _, opt := range v.Options {
				i := item{value: opt, label: opt}
				if strings.Contains(opt, "|") {
					sl := strings.Split(opt, "|")
					i.label = sl[1]
					i.value = sl[0]
				}

				items = append(items, i)
			}

			f := func(items []list.Item) tea.Cmd {
				return func() tea.Msg {
					return items
				}
			}

			pickerPage := newPicker(v.Description, "", v.Name, v.Default, f(items))
			q.add(&pickerPage)
			continue
		}

		if len(temp) < 1 {
			tiPage := newCustom(v)
			q.add(tiPage)
		}

	}
}

func newGCEInstance(q *Queue) {
	r := newPicker("Do you want to accept the default configuration? (Yes or No)", "", "gce-use-defaults", "", getYesOrNo(q))
	r.omitFromSettings = true
	r.list.SetShowFilter(false)
	r.list.SetShowHelp(false)
	r.list.SetShowStatusBar(false)
	r.addPostProcessor(validateGCEDefault)
	r.addContent(textStyle.Bold(true).Render("Configure a Compute Engine Instance"))
	r.addContent("\n\n")

	m := `Let's walk through configuring a Compute Engine Instance (Virtual Machine).
you can either accept a default configuration with settings that work for
trying out most use cases, or hand configure key settings.
	`

	r.addContent(m)
	q.add(&r)

	basename := q.stack.GetSetting("basename")
	name := newTextInput("Enter the name of the instance",
		fmt.Sprintf("%s-instance", basename),
		"instance-name",
		"",
	)
	q.add(&name)

	newRegion(q)
	newZone(q)
	newMachineTypeManager(q)
	newDiskImageManager(q)

	ds := newTextInput("Enter the size of the boot disk you want in GB",
		"100",
		"instance-disksize",
		"",
	)
	ds.addPostProcessor(validateInteger)
	q.add(&ds)

	dt := newPicker("Pick the type of the boot disk you want", "", "instance-disktype", gcloud.DefaultDiskType, getDiskTypes(q))
	q.add(&dt)

	dy := newYesOrNo(
		q,
		"Do you want this to be a webserver (Expose ports 80 & 443)?",
		"instance-webserver",
		false,
		validateGCEConfiguration,
	)
	q.add(&dy)
}

func newRegion(q *Queue) {
	r := newPicker("Pick a region", "Retrieving regions", "region", q.stack.Config.RegionDefault, getRegions(q))
	q.add(&r)
}

func newZone(q *Queue) {
	z := newPicker("Pick a zone", "Retrieving zones", "zone", gcloud.DefaultZone, getZones(q))
	q.add(&z)
}

func newMachineTypeManager(q *Queue) {
	p := newPicker("Pick a Machine Type Family", "Retrieving machine type familes", "instance-machine-type-family", gcloud.DefaultMachineFamily, getMachineTypeFamilies(q))
	p.addContent(textStyle.Bold(true).Render("Configure a Compute Engine Instance"))
	p.addContent("\n\n")
	p.addContent("There are a large number of machine types to choose from. For more information \n")
	p.addContent("please refer to the following link for more infomation about Machine types: \n")
	p.addContent(url.Render("https://cloud.google.com/compute/docs/machine-types"))
	q.add(&p)

	p2 := newPicker("Pick a Machine Type", "Retrieving machine types", "instance-machine-type", gcloud.DefaultMachineType, getMachineTypes(q))
	p2.addContent(textStyle.Bold(true).Render("Configure a Compute Engine Instance"))
	p2.addContent("\n\n")
	p2.addContent("There are a large number of machine types to choose from. For more information \n")
	p2.addContent("please refer to the following link for more infomation about Machine types: \n")
	p2.addContent(url.Render("https://cloud.google.com/compute/docs/machine-types"))
	q.add(&p2)
}

func newDiskImageManager(q *Queue) {
	p := newPicker("Pick an operating system", "Retrieving operating systems", "instance-image-project", gcloud.DefaultImageProject, getDiskProjects(q))
	p.addContent(textStyle.Bold(true).Render("Configure a Compute Engine Instance"))
	p.addContent("\n\n")
	p.addContent("There are a large number of machine images to choose from. For more information \n")
	p.addContent("please refer to the following link for more infomation about Machine images: \n")
	p.addContent(url.Render("https://cloud.google.com/compute/docs/images"))
	q.add(&p)

	p2 := newPicker("Pick a disk family", "Retrieving disk family", "instance-image-family", gcloud.DefaultImageFamily, getImageFamilies(q))
	p2.addContent(textStyle.Bold(true).Render("Configure a Compute Engine Instance"))
	p2.addContent("\n\n")
	p2.addContent("There are a large number of machine images to choose from. For more information \n")
	p2.addContent("please refer to the following link for more infomation about Machine images: \n")
	p2.addContent(url.Render("https://cloud.google.com/compute/docs/images"))
	q.add(&p2)

	p3 := newPicker("Pick a disk image", "Retrieving disk image", "instance-image", "", getImageDisks(q))
	p3.addContent(textStyle.Bold(true).Render("Configure a Compute Engine Instance"))
	p3.addContent("\n\n")
	p3.addContent("There are a large number of machine images to choose from. For more information \n")
	p3.addContent("please refer to the following link for more infomation about Machine images: \n")
	p3.addContent(url.Render("https://cloud.google.com/compute/docs/images"))
	q.add(&p3)
}
