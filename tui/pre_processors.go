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

	"github.com/GoogleCloudPlatform/deploystack/gcloud"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func getProjects(q *Queue) tea.Cmd {
	return func() tea.Msg {
		p, err := q.client.ProjectList()
		if err != nil {
			return errMsg{err: err}
		}

		items := []list.Item{}
		for _, v := range p {
			if !v.BillingEnabled {
				label := fmt.Sprintf("%s (Billing Diabled)", v.Name)
				items = append(items, item{value: v.ID, label: billingDisabledStyle.Render(label)})
				continue
			}
			items = append(items, item{
				value: strings.TrimSpace(v.ID),
				label: strings.TrimSpace(v.Name),
			})
		}

		return items
	}
}

func getBillingAccounts(q *Queue) tea.Cmd {
	return func() tea.Msg {
		p, err := q.client.BillingAccountList()
		if err != nil {
			return errMsg{err: err}
		}

		items := []list.Item{}
		for _, v := range p {
			id := strings.ReplaceAll(v.Name, "billingAccounts/", "")
			items = append(items, item{
				value: strings.TrimSpace(id),
				label: strings.TrimSpace(v.DisplayName),
			})
		}

		// If there is only 1 billing account, don't bother the user with
		// setting it up.
		if len(items) == 1 {
			ba := strings.ReplaceAll(p[0].Name, "billingAccounts/", "")

			key := strings.ReplaceAll(q.currentKey(), billNewSuffix, "")
			project := q.stack.GetSetting(key)
			if err := q.client.BillingAccountAttach(project, ba); err != nil {
				return errMsg{err: fmt.Errorf("attachBilling: could not attach billing to project: %w", err)}
			}
			return successMsg{}
		}

		return items
	}
}

func getRegions(q *Queue) tea.Cmd {
	return func() tea.Msg {
		s := q.stack
		project := s.GetSetting("project_id")
		product := s.Config.RegionType

		p, err := q.client.RegionList(project, product)
		if err != nil {
			return errMsg{err: err}
		}

		items := []list.Item{}
		for _, v := range p {
			items = append(items, item{
				value: strings.TrimSpace(v),
				label: strings.TrimSpace(v),
			})
		}

		return items
	}
}

func getZones(q *Queue) tea.Cmd {
	return func() tea.Msg {
		s := q.stack
		project := s.GetSetting("project_id")
		region := s.Settings["region"]

		p, err := q.client.ZoneList(project, region)
		if err != nil {
			return errMsg{err: err}
		}

		items := []list.Item{}
		for _, v := range p {
			items = append(items, item{
				value: strings.TrimSpace(v),
				label: strings.TrimSpace(v),
			})
		}

		return items
	}
}

func getMachineTypeFamilies(q *Queue) tea.Cmd {
	return func() tea.Msg {
		s := q.stack
		project := s.GetSetting("project_id")
		zone := s.GetSetting("zone")

		// TODO: add caching to remove this double request overhead
		types, err := q.client.MachineTypeList(project, zone)
		if err != nil {
			return errMsg{err: err}
		}

		typefamilies := q.client.MachineTypeFamilyList(types)

		items := []list.Item{}
		for _, v := range typefamilies {
			items = append(items, item{
				value: strings.TrimSpace(v.Value),
				label: strings.TrimSpace(v.Label),
			})
		}

		return items
	}
}

func getMachineTypes(q *Queue) tea.Cmd {
	return func() tea.Msg {
		s := q.stack
		project := s.GetSetting("project_id")
		zone := s.GetSetting("zone")
		family := s.GetSetting("instance-machine-type-family")

		// TODO: add caching to remove this double request overhead
		types, err := q.client.MachineTypeList(project, zone)
		if err != nil {
			return errMsg{err: err}
		}

		filteredtypes := q.client.MachineTypeListByFamily(types, family)

		items := []list.Item{}
		for _, v := range filteredtypes {
			items = append(items, item{
				value: strings.TrimSpace(v.Value),
				label: strings.TrimSpace(v.Label),
			})
		}

		return items
	}
}

func getDiskProjects(q *Queue) tea.Cmd {
	return func() tea.Msg {
		diskImages := gcloud.DiskProjects

		items := []list.Item{}
		for _, v := range diskImages {
			items = append(items, item{
				value: strings.TrimSpace(v.Value),
				label: strings.TrimSpace(v.Label),
			})
		}

		return items
	}
}

func getImageFamilies(q *Queue) tea.Cmd {
	return func() tea.Msg {
		s := q.stack
		instanceImageProject := s.GetSetting("instance-image-project")
		project := s.GetSetting("project_id")

		images, err := q.client.ImageList(project, instanceImageProject)
		if err != nil {
			return errMsg{err: err}
		}

		families := q.client.ImageFamilyList(images)

		items := []list.Item{}
		for _, v := range families {
			items = append(items, item{
				value: strings.TrimSpace(v.Value),
				label: strings.TrimSpace(v.Label),
			})
		}

		return items
	}
}

func getImageDisks(q *Queue) tea.Cmd {
	return func() tea.Msg {
		s := q.stack
		instanceImageProject := s.GetSetting("instance-image-project")
		instanceImageFamily := s.GetSetting("instance-image-family")
		project := s.GetSetting("project_id")

		images, err := q.client.ImageList(project, instanceImageProject)
		if err != nil {
			return errMsg{err: err}
		}

		imagesByFam := q.client.ImageTypeListByFamily(images, instanceImageProject, instanceImageFamily)

		items := []list.Item{}
		for _, v := range imagesByFam {
			items = append(items, item{
				value: strings.TrimSpace(v.Value),
				label: strings.TrimSpace(v.Label),
			})
		}

		return items
	}
}

func getDiskTypes(q *Queue) tea.Cmd {
	return func() tea.Msg {
		items := []list.Item{
			item{"Standard", "pd-standard"},
			item{"Balanced", "pd-balanced"},
			item{"SSD", "pd-sdd"},
		}

		return items
	}
}

func getYesOrNo(q *Queue) tea.Cmd {
	return func() tea.Msg {
		items := []list.Item{
			item{"Yes", "y"},
			item{"No", "n"},
		}

		return items
	}
}

func getNoOrYes(q *Queue) tea.Cmd {
	return func() tea.Msg {
		items := []list.Item{
			item{"No", "n"},
			item{"Yes", "y"},
		}

		return items
	}
}

func cleanUp(q *Queue) tea.Cmd {
	return func() tea.Msg {
		// // Don't let these get leaked to terraform
		q.stack.DeleteSetting("domain_consent")

		// Stop autogenerated billing forms from leaking into settings
		for i := range q.stack.Settings {
			if strings.Contains(i, billNewSuffix) {
				q.stack.DeleteSetting(i)
			}
		}

		return ""
	}
}
