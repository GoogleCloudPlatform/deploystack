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
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"cloud.google.com/go/domains/apiv1beta1/domainspb"
	"github.com/GoogleCloudPlatform/deploystack/gcloud"
	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/compute/v1"
	"google.golang.org/genproto/googleapis/type/money"
)

// GetMock returns a mock gcloud.Client
// from github.com/GoogleCloudPlatform/deploystack/gcloud
//
//revive:disable:unexported-return
func GetMock(delay int) mock {
	return mock{d: delay}
}

//revive:enable:unexported-return

type mock struct {
	d        int
	forceErr bool
	cache    map[string]interface{}
}

func (m mock) delay() {
	time.Sleep(time.Second * time.Duration(m.d))
}

func (m mock) ProjectIDGet() (string, error) {
	if m.forceErr {
		return "", errForced
	}
	return "ds-tester-singlevm", nil
}

func (m mock) ProjectIDSet(id string) error {
	if m.forceErr {
		return errForced
	}
	return nil
}

func (m mock) ProjectList() ([]gcloud.ProjectWithBilling, error) {
	m.delay()
	if m.forceErr {
		return nil, errForced
	}
	r := []gcloud.ProjectWithBilling{
		{ID: "ds-test-ms-ua2jjt3u", Name: "ds-test-ms-ua2jjt3u", BillingEnabled: true},
		{ID: "ds-test-ms-bpbfnumc", Name: "ds-test-ms-bpbfnumc", BillingEnabled: true},
		{ID: "ds-test-ms-smx0mkqq", Name: "ds-test-ms-smx0mkqq", BillingEnabled: true},
		{ID: "ds-test-ms-8vt3qrfj", Name: "ds-test-ms-8vt3qrfj", BillingEnabled: true},
		{ID: "ds-test-ms-q5gttybe", Name: "ds-test-ms-q5gttybe", BillingEnabled: true},
		{ID: "load-balanced-vms-4ead", Name: "load-balanced-vms", BillingEnabled: true},
		{ID: "ds-test-ms-0y27r3x2", Name: "ds-test-ms-0y27r3x2", BillingEnabled: true},
		{ID: "ds-test-ms-850k1gfj", Name: "ds-test-ms-850k1gfj", BillingEnabled: true},
		{ID: "ds-test-ms-cq98yj8x", Name: "ds-test-ms-cq98yj8x", BillingEnabled: true},
		{ID: "ds-test-ms-cfjsrtp3", Name: "ds-test-ms-cfjsrtp3", BillingEnabled: true},
		{ID: "ds-test-ms-t9dkufjm", Name: "ds-test-ms-t9dkufjm", BillingEnabled: true},
		{ID: "ds-test-ms-cdoafhqa", Name: "ds-test-ms-cdoafhqa", BillingEnabled: true},
		{ID: "ds-test-ms-it7w27el", Name: "ds-test-ms-it7w27el", BillingEnabled: true},
		{ID: "ds-test-ms-cfuk8b5v", Name: "ds-test-ms-cfuk8b5v", BillingEnabled: true},
		{ID: "ds-test-ms-qpridp11", Name: "ds-test-ms-qpridp11", BillingEnabled: true},
		{ID: "ds-test-ms-qhmn8elm", Name: "ds-test-ms-qhmn8elm", BillingEnabled: true},
		{ID: "ds-test-ms-4p9szpjt", Name: "ds-test-ms-4p9szpjt", BillingEnabled: true},
		{ID: "ds-test-ms-czhyncv2", Name: "ds-test-ms-czhyncv2", BillingEnabled: true},
		{ID: "ds-test-ms-gisjij3o", Name: "ds-test-ms-gisjij3o", BillingEnabled: true},
		{ID: "ds-test-ms-wtfosvv2", Name: "ds-test-ms-wtfosvv2", BillingEnabled: true},
		{ID: "ds-test-ms-t6g7l7el", Name: "ds-test-ms-t6g7l7el", BillingEnabled: true},
		{ID: "ds-test-ms-odrfhxu1", Name: "ds-test-ms-odrfhxu1", BillingEnabled: true},
		{ID: "ds-test-ms-dcq41vmo", Name: "ds-test-ms-dcq41vmo", BillingEnabled: true},
		{ID: "ds-test-ms-jnsq6zr4", Name: "ds-test-ms-jnsq6zr4", BillingEnabled: true},
		{ID: "ds-test-ms-ikvy5obn", Name: "ds-test-ms-ikvy5obn", BillingEnabled: true},
		{ID: "ds-test-ms-pdmymgst", Name: "ds-test-ms-pdmymgst", BillingEnabled: true},
		{ID: "ds-test-ms-1hkja8o9", Name: "ds-test-ms-1hkja8o9", BillingEnabled: true},
		{ID: "ds-test-ms-f3nimk87", Name: "ds-test-ms-f3nimk87", BillingEnabled: true},
		{ID: "ds-test-ms-xh1isutj", Name: "ds-test-ms-xh1isutj", BillingEnabled: true},
		{ID: "ds-test-ms-mkso9apf", Name: "ds-test-ms-mkso9apf", BillingEnabled: true},
		{ID: "ds-tester-glb-and-armor", Name: "ds-tester-glb-and-armor", BillingEnabled: true},
		{ID: "ds-tester-auditlogs", Name: "ds-tester-auditlogs", BillingEnabled: true},
		{ID: "ds-tester-wordpress-run", Name: "ds-tester-wordpress-run", BillingEnabled: true},
		{ID: "ds-tester-cloudsql-multiregion", Name: "ds-tester-cloudsql-multiregion", BillingEnabled: true},
		{ID: "ds-tester-gcs-to-bq", Name: "ds-tester-gcs-to-bq", BillingEnabled: true},
		{ID: "sic-deleteme-3ta-373719", Name: "sic-deleteme-3ta", BillingEnabled: true},
		{ID: "tpryan-test-project", Name: "tpryan-test-project", BillingEnabled: true},
		{ID: "tf-contributions-tpryan", Name: "tf-contributions-tpryan", BillingEnabled: true},
		{ID: "ds-tooling-app", Name: "ds-tooling-app", BillingEnabled: true},
		{ID: "coldfusion-demo-2", Name: "coldfusion-demo-2", BillingEnabled: true},
		{ID: "coldfusion-demo", Name: "coldfusion-demo", BillingEnabled: true},
		{ID: "ds-tester-microservices-demo", Name: "ds-tester-microservices-demo", BillingEnabled: true},
		{ID: "sic-tester", Name: "sic-tester", BillingEnabled: true},
		{ID: "ds-tester-e2e-new", Name: "ds-tester-e2e-new", BillingEnabled: true},
		{ID: "ds-break-things", Name: "DS-BREAK-THINGS", BillingEnabled: false},
		{ID: "coltsays-360004", Name: "coltsays", BillingEnabled: true},
		{ID: "ds-tester-etl-pipeline", Name: "ds-tester-etl-pipeline", BillingEnabled: true},
		{ID: "sic-container-repo", Name: "sic-container-repo", BillingEnabled: true},
		{ID: "ds-opsagent", Name: "ds-opsagent", BillingEnabled: true},
		{ID: "ds-tester-nosql-client-server", Name: "ds-tester-nosql-client-server", BillingEnabled: true},
		{ID: "neos-tester", Name: "neos-tester", BillingEnabled: true},
		{ID: "ds-artifacts-cloudshell", Name: "ds-artifacts-cloudshell", BillingEnabled: true},
		{ID: "summit-walkthrough", Name: "summit-walkthrough", BillingEnabled: true},
		{ID: "ds-tester-todo-fixed", Name: "ds-tester-todo-fixed", BillingEnabled: true},
		{ID: "ds-tester-opsagent", Name: "ds-tester-opsagent", BillingEnabled: true},
		{ID: "ds-tester-singlevm", Name: "ds-tester-singlevm", BillingEnabled: true},
		{ID: "run-integrations-test", Name: "run-integrations-test", BillingEnabled: true},
		{ID: "ds-tester-deploystack", Name: "ds-tester-deploystack", BillingEnabled: true},
		{ID: "ds-test-no-billing", Name: "ds-test-no-billing", BillingEnabled: false},
		{ID: "ds-tester-helper", Name: "ds-tester-helper", BillingEnabled: true},
		{ID: "ds-tester-basiclb", Name: "ds-tester-basiclb", BillingEnabled: true},
		{ID: "ds-tester-yesornosite", Name: "ds-tester-yesornosite", BillingEnabled: true},
		{ID: "ds-tester-scaler", Name: "ds-tester-scaler", BillingEnabled: true},
		{ID: "ds-tester-costsentry", Name: "ds-tester-costsentry", BillingEnabled: true},
		{ID: "deploystack-terraform-2", Name: "deploystack-terraform-2", BillingEnabled: true},
		{ID: "deploystack-terraform", Name: "deploystack-terraform", BillingEnabled: true},
		{ID: "deploy-terraform", Name: "deploy-terraform", BillingEnabled: true},
		{ID: "stack-terraform", Name: "stack-terraform", BillingEnabled: true},
		{ID: "microsites-deploystack", Name: "microsites-deploystack", BillingEnabled: true},
		{ID: "microsites-stackables", Name: "microsites-stackables", BillingEnabled: true},
		{ID: "vertexaitester", Name: "vertexaitester", BillingEnabled: true},
		{ID: "stackinaboxtester", Name: "stackinaboxtester", BillingEnabled: true},
		{ID: "stackinabox", Name: "stackinabox", BillingEnabled: true},
		{ID: "aiab-test-project", Name: "aiab-test-project", BillingEnabled: true},
		{ID: "cost-sentry-experiments", Name: "cost-sentry-experiments", BillingEnabled: true},
		{ID: "appinabox-yesornosite-demo", Name: "deploystack-yesornosite-demo", BillingEnabled: true},
		{ID: "bucketsite-test", Name: "bucketsite-test", BillingEnabled: true},
		{ID: "microsites-appinabox", Name: "microsites-appinabox", BillingEnabled: true},
		{ID: "scaler-microsite", Name: "scaler-microsite", BillingEnabled: true},
		{ID: "todo-microsite", Name: "todo-microsite", BillingEnabled: true},
		{ID: "neosregional", Name: "NeosRegional", BillingEnabled: true},
		{ID: "cloudicons", Name: "cloudicons", BillingEnabled: true},
		{ID: "sustained-racer-323200", Name: "GoogleCloudCheatSheet", BillingEnabled: true},
		{ID: "cloud-logging-generator", Name: "cloud-logging", BillingEnabled: true},
		{ID: "neos-log-test", Name: "neos-log-test", BillingEnabled: true},
		{ID: "neos-test-304321", Name: "neos-test", BillingEnabled: true},
	}

	sort.Slice(r, func(i, j int) bool {
		return strings.ToLower(r[i].Name) < strings.ToLower(r[j].Name)
	})

	return r, nil
}

func (m mock) RegionList(project, product string) ([]string, error) {
	m.delay()
	if m.forceErr {
		return nil, errForced
	}
	r := []string{
		"asia-east1",
		"asia-east2",
		"asia-northeast1",
		"asia-northeast2",
		"asia-northeast3",
		"asia-south1",
		"asia-south2",
		"asia-southeast1",
		"asia-southeast2",
		"australia-southeast1",
		"australia-southeast2",
		"europe-central2",
		"europe-north1",
		"europe-southwest1",
		"europe-west1",
		"europe-west2",
		"europe-west3",
		"europe-west4",
		"europe-west6",
		"europe-west8",
		"europe-west9",
		"me-west1",
		"northamerica-northeast1",
		"northamerica-northeast2",
		"southamerica-east1",
		"southamerica-west1",
		"us-central1",
		"us-east1",
		"us-east4",
		"us-east5",
		"us-south1",
		"us-west1",
		"us-west2",
		"us-west3",
		"us-west4",
	}

	return r, nil
}

func (m mock) ZoneList(project, region string) ([]string, error) {
	m.delay()
	if m.forceErr {
		return nil, errForced
	}
	z := []string{
		"us-east1-b",
		"us-east1-c",
		"us-east1-d",
		"us-east4-c",
		"us-east4-b",
		"us-east4-a",
		"us-central1-c",
		"us-central1-a",
		"us-central1-f",
		"us-central1-b",
		"us-west1-b",
		"us-west1-c",
		"us-west1-a",
		"europe-west4-a",
		"europe-west4-b",
		"europe-west4-c",
		"europe-west1-b",
		"europe-west1-d",
		"europe-west1-c",
		"europe-west3-c",
		"europe-west3-a",
		"europe-west3-b",
		"europe-west2-c",
		"europe-west2-b",
		"europe-west2-a",
		"asia-east1-b",
		"asia-east1-a",
		"asia-east1-c",
		"asia-southeast1-b",
		"asia-southeast1-a",
		"asia-southeast1-c",
		"asia-northeast1-b",
		"asia-northeast1-c",
		"asia-northeast1-a",
		"asia-south1-c",
		"asia-south1-b",
		"asia-south1-a",
		"australia-southeast1-b",
		"australia-southeast1-c",
		"australia-southeast1-a",
		"southamerica-east1-b",
		"southamerica-east1-c",
		"southamerica-east1-a",
		"asia-east2-a",
		"asia-east2-b",
		"asia-east2-c",
		"asia-northeast2-a",
		"asia-northeast2-b",
		"asia-northeast2-c",
		"asia-northeast3-a",
		"asia-northeast3-b",
		"asia-northeast3-c",
		"asia-south2-a",
		"asia-south2-b",
		"asia-south2-c",
		"asia-southeast2-a",
		"asia-southeast2-b",
		"asia-southeast2-c",
		"australia-southeast2-a",
		"australia-southeast2-b",
		"australia-southeast2-c",
		"europe-central2-a",
		"europe-central2-b",
		"europe-central2-c",
		"europe-north1-a",
		"europe-north1-b",
		"europe-north1-c",
		"europe-southwest1-a",
		"europe-southwest1-b",
		"europe-southwest1-c",
		"europe-west6-a",
		"europe-west6-b",
		"europe-west6-c",
		"europe-west8-a",
		"europe-west8-b",
		"europe-west8-c",
		"europe-west9-a",
		"europe-west9-b",
		"europe-west9-c",
		"me-west1-a",
		"me-west1-b",
		"me-west1-c",
		"northamerica-northeast1-a",
		"northamerica-northeast1-b",
		"northamerica-northeast1-c",
		"northamerica-northeast2-a",
		"northamerica-northeast2-b",
		"northamerica-northeast2-c",
		"southamerica-west1-a",
		"southamerica-west1-b",
		"southamerica-west1-c",
		"us-east5-a",
		"us-east5-b",
		"us-east5-c",
		"us-south1-a",
		"us-south1-b",
		"us-south1-c",
		"us-west2-a",
		"us-west2-b",
		"us-west2-c",
		"us-west3-a",
		"us-west3-b",
		"us-west3-c",
		"us-west4-a",
		"us-west4-b",
		"us-west4-c",
	}

	r := []string{}

	for _, v := range z {
		if strings.Contains(v, region) {
			r = append(r, v)
		}
	}

	sort.Strings(r)

	return r, nil
}

func (m mock) ProjectParentGet(project string) (*cloudresourcemanager.ResourceId, error) {
	m.delay()
	if m.forceErr {
		return nil, errForced
	}
	r := &cloudresourcemanager.ResourceId{}

	r.Id = "298490623289"
	r.Type = "organization"
	m.delay()
	return r, nil
}

func (m mock) ProjectCreate(project, parent, parentType string) error {
	m.delay()
	if m.forceErr {
		return errForced
	}
	if len(project) > 32 {
		return gcloud.ErrorProjectCreateTooLong
	}

	if len(project) < 6 {
		return gcloud.ErrorProjectCreateTooLong
	}

	if strings.Contains(project, "!") {
		return gcloud.ErrorProjectInvalidCharacters
	}

	list, _ := m.ProjectList()

	for _, v := range list {
		if v.ID == project {
			return gcloud.ErrorProjectAlreadyExists
		}
		if v.Name == project {
			return gcloud.ErrorProjectAlreadyExists
		}
	}

	return nil
}

func (m mock) DomainIsAvailable(project, domain string) (*domainspb.RegisterParameters, error) {
	m.delay()
	if m.forceErr {
		return nil, errForced
	}
	r := &domainspb.RegisterParameters{}

	if domain == "example.com" {
		r.Availability = domainspb.RegisterParameters_UNAVAILABLE
		return r, nil
	}

	if domain == "example2.com" {
		r.Availability = domainspb.RegisterParameters_UNAVAILABLE
		return r, nil
	}

	r.DomainName = domain
	r.Availability = domainspb.RegisterParameters_AVAILABLE
	r.YearlyPrice = &money.Money{
		Units:        12,
		CurrencyCode: "USD",
	}

	return r, nil
}

func (m mock) DomainIsVerified(project, domain string) (bool, error) {
	m.delay()
	if m.forceErr {
		return false, errForced
	}
	if domain == "example2.com" {
		return false, nil
	}
	if domain == "example.com" {
		return false, fmt.Errorf("domain is not verified")
	}

	return true, nil
}

func (m mock) DomainRegister(project string, domaininfo *domainspb.RegisterParameters, contact gcloud.ContactData) error {
	m.delay()
	if m.forceErr {
		return errForced
	}
	if domaininfo.DomainName == "example3.com" {
		return fmt.Errorf("domain is cursed and cannot be obtained by mortals")
	}
	if domaininfo.DomainName == "example2.com" {
		return fmt.Errorf("domain is already owned. This should have been caught")
	}
	if domaininfo.DomainName == "example.com" {
		return fmt.Errorf("domain is already owned. This should have been caught")
	}

	return nil
}

func (m mock) ImageLatestGet(project, imageproject, imagefamily string) (string, error) {
	if m.forceErr {
		return "", errForced
	}
	return "debian-cloud/debian-11-bullseye-v20230202", nil
}

func (m mock) MachineTypeList(project, zone string) (*compute.MachineTypeList, error) {
	if m.forceErr {
		return nil, errForced
	}
	r := compute.MachineTypeList{
		Items: []*compute.MachineType{
			{GuestCpus: 12, MemoryMb: 87040, Name: "a2-highgpu-1g"},
			{GuestCpus: 24, MemoryMb: 174080, Name: "a2-highgpu-2g"},
			{GuestCpus: 48, MemoryMb: 348160, Name: "a2-highgpu-4g"},
			{GuestCpus: 96, MemoryMb: 696320, Name: "a2-highgpu-8g"},
			{GuestCpus: 96, MemoryMb: 1392640, Name: "a2-megagpu-16g"},
			{GuestCpus: 16, MemoryMb: 65536, Name: "c2-standard-16"},
			{GuestCpus: 30, MemoryMb: 122880, Name: "c2-standard-30"},
			{GuestCpus: 4, MemoryMb: 16384, Name: "c2-standard-4"},
			{GuestCpus: 60, MemoryMb: 245760, Name: "c2-standard-60"},
			{GuestCpus: 8, MemoryMb: 32768, Name: "c2-standard-8"},
			{GuestCpus: 112, MemoryMb: 229376, Name: "c2d-highcpu-112"},
			{GuestCpus: 16, MemoryMb: 32768, Name: "c2d-highcpu-16"},
			{GuestCpus: 2, MemoryMb: 4096, Name: "c2d-highcpu-2"},
			{GuestCpus: 32, MemoryMb: 65536, Name: "c2d-highcpu-32"},
			{GuestCpus: 4, MemoryMb: 8192, Name: "c2d-highcpu-4"},
			{GuestCpus: 56, MemoryMb: 114688, Name: "c2d-highcpu-56"},
			{GuestCpus: 8, MemoryMb: 16384, Name: "c2d-highcpu-8"},
			{GuestCpus: 112, MemoryMb: 917504, Name: "c2d-highmem-112"},
			{GuestCpus: 16, MemoryMb: 131072, Name: "c2d-highmem-16"},
			{GuestCpus: 2, MemoryMb: 16384, Name: "c2d-highmem-2"},
			{GuestCpus: 32, MemoryMb: 262144, Name: "c2d-highmem-32"},
			{GuestCpus: 4, MemoryMb: 32768, Name: "c2d-highmem-4"},
			{GuestCpus: 56, MemoryMb: 458752, Name: "c2d-highmem-56"},
			{GuestCpus: 8, MemoryMb: 65536, Name: "c2d-highmem-8"},
			{GuestCpus: 112, MemoryMb: 458752, Name: "c2d-standard-112"},
			{GuestCpus: 16, MemoryMb: 65536, Name: "c2d-standard-16"},
			{GuestCpus: 2, MemoryMb: 8192, Name: "c2d-standard-2"},
			{GuestCpus: 32, MemoryMb: 131072, Name: "c2d-standard-32"},
			{GuestCpus: 4, MemoryMb: 16384, Name: "c2d-standard-4"},
			{GuestCpus: 56, MemoryMb: 229376, Name: "c2d-standard-56"},
			{GuestCpus: 8, MemoryMb: 32768, Name: "c2d-standard-8"},
			{GuestCpus: 16, MemoryMb: 16384, Name: "e2-highcpu-16"},
			{GuestCpus: 2, MemoryMb: 2048, Name: "e2-highcpu-2"},
			{GuestCpus: 32, MemoryMb: 32768, Name: "e2-highcpu-32"},
			{GuestCpus: 4, MemoryMb: 4096, Name: "e2-highcpu-4"},
			{GuestCpus: 8, MemoryMb: 8192, Name: "e2-highcpu-8"},
			{GuestCpus: 16, MemoryMb: 131072, Name: "e2-highmem-16"},
			{GuestCpus: 2, MemoryMb: 16384, Name: "e2-highmem-2"},
			{GuestCpus: 4, MemoryMb: 32768, Name: "e2-highmem-4"},
			{GuestCpus: 8, MemoryMb: 65536, Name: "e2-highmem-8"},
			{GuestCpus: 2, MemoryMb: 4096, Name: "e2-medium"},
			{GuestCpus: 2, MemoryMb: 1024, Name: "e2-micro"},
			{GuestCpus: 2, MemoryMb: 2048, Name: "e2-small"},
			{GuestCpus: 16, MemoryMb: 65536, Name: "e2-standard-16"},
			{GuestCpus: 2, MemoryMb: 8192, Name: "e2-standard-2"},
			{GuestCpus: 32, MemoryMb: 131072, Name: "e2-standard-32"},
			{GuestCpus: 4, MemoryMb: 16384, Name: "e2-standard-4"},
			{GuestCpus: 8, MemoryMb: 32768, Name: "e2-standard-8"},
			{GuestCpus: 1, MemoryMb: 614, Name: "f1-micro"},
			{GuestCpus: 1, MemoryMb: 1740, Name: "g1-small"},
			{GuestCpus: 96, MemoryMb: 1468006, Name: "m1-megamem-96"},
			{GuestCpus: 160, MemoryMb: 3936256, Name: "m1-ultramem-160"},
			{GuestCpus: 40, MemoryMb: 984064, Name: "m1-ultramem-40"},
			{GuestCpus: 80, MemoryMb: 1968128, Name: "m1-ultramem-80"},
			{GuestCpus: 416, MemoryMb: 9043968, Name: "m2-hypermem-416"},
			{GuestCpus: 416, MemoryMb: 6029312, Name: "m2-megamem-416"},
			{GuestCpus: 208, MemoryMb: 6029312, Name: "m2-ultramem-208"},
			{GuestCpus: 416, MemoryMb: 12058624, Name: "m2-ultramem-416"},
			{GuestCpus: 128, MemoryMb: 1998848, Name: "m3-megamem-128"},
			{GuestCpus: 64, MemoryMb: 999424, Name: "m3-megamem-64"},
			{GuestCpus: 128, MemoryMb: 3997696, Name: "m3-ultramem-128"},
			{GuestCpus: 32, MemoryMb: 999424, Name: "m3-ultramem-32"},
			{GuestCpus: 64, MemoryMb: 1998848, Name: "m3-ultramem-64"},
			{GuestCpus: 16, MemoryMb: 14746, Name: "n1-highcpu-16"},
			{GuestCpus: 2, MemoryMb: 1843, Name: "n1-highcpu-2"},
			{GuestCpus: 32, MemoryMb: 29491, Name: "n1-highcpu-32"},
			{GuestCpus: 4, MemoryMb: 3686, Name: "n1-highcpu-4"},
			{GuestCpus: 64, MemoryMb: 58982, Name: "n1-highcpu-64"},
			{GuestCpus: 8, MemoryMb: 7373, Name: "n1-highcpu-8"},
			{GuestCpus: 96, MemoryMb: 88474, Name: "n1-highcpu-96"},
			{GuestCpus: 16, MemoryMb: 106496, Name: "n1-highmem-16"},
			{GuestCpus: 2, MemoryMb: 13312, Name: "n1-highmem-2"},
			{GuestCpus: 32, MemoryMb: 212992, Name: "n1-highmem-32"},
			{GuestCpus: 4, MemoryMb: 26624, Name: "n1-highmem-4"},
			{GuestCpus: 64, MemoryMb: 425984, Name: "n1-highmem-64"},
			{GuestCpus: 8, MemoryMb: 53248, Name: "n1-highmem-8"},
			{GuestCpus: 96, MemoryMb: 638976, Name: "n1-highmem-96"},
			{GuestCpus: 96, MemoryMb: 1468006, Name: "n1-megamem-96"},
			{GuestCpus: 1, MemoryMb: 3840, Name: "n1-standard-1"},
			{GuestCpus: 16, MemoryMb: 61440, Name: "n1-standard-16"},
			{GuestCpus: 2, MemoryMb: 7680, Name: "n1-standard-2"},
			{GuestCpus: 32, MemoryMb: 122880, Name: "n1-standard-32"},
			{GuestCpus: 4, MemoryMb: 15360, Name: "n1-standard-4"},
			{GuestCpus: 64, MemoryMb: 245760, Name: "n1-standard-64"},
			{GuestCpus: 8, MemoryMb: 30720, Name: "n1-standard-8"},
			{GuestCpus: 96, MemoryMb: 368640, Name: "n1-standard-96"},
			{GuestCpus: 160, MemoryMb: 3936256, Name: "n1-ultramem-160"},
			{GuestCpus: 40, MemoryMb: 984064, Name: "n1-ultramem-40"},
			{GuestCpus: 80, MemoryMb: 1968128, Name: "n1-ultramem-80"},
			{GuestCpus: 16, MemoryMb: 16384, Name: "n2-highcpu-16"},
			{GuestCpus: 2, MemoryMb: 2048, Name: "n2-highcpu-2"},
			{GuestCpus: 32, MemoryMb: 32768, Name: "n2-highcpu-32"},
			{GuestCpus: 4, MemoryMb: 4096, Name: "n2-highcpu-4"},
			{GuestCpus: 48, MemoryMb: 49152, Name: "n2-highcpu-48"},
			{GuestCpus: 64, MemoryMb: 65536, Name: "n2-highcpu-64"},
			{GuestCpus: 8, MemoryMb: 8192, Name: "n2-highcpu-8"},
			{GuestCpus: 80, MemoryMb: 81920, Name: "n2-highcpu-80"},
			{GuestCpus: 96, MemoryMb: 98304, Name: "n2-highcpu-96"},
			{GuestCpus: 128, MemoryMb: 884736, Name: "n2-highmem-128"},
			{GuestCpus: 16, MemoryMb: 131072, Name: "n2-highmem-16"},
			{GuestCpus: 2, MemoryMb: 16384, Name: "n2-highmem-2"},
			{GuestCpus: 32, MemoryMb: 262144, Name: "n2-highmem-32"},
			{GuestCpus: 4, MemoryMb: 32768, Name: "n2-highmem-4"},
			{GuestCpus: 48, MemoryMb: 393216, Name: "n2-highmem-48"},
			{GuestCpus: 64, MemoryMb: 524288, Name: "n2-highmem-64"},
			{GuestCpus: 8, MemoryMb: 65536, Name: "n2-highmem-8"},
			{GuestCpus: 80, MemoryMb: 655360, Name: "n2-highmem-80"},
			{GuestCpus: 96, MemoryMb: 786432, Name: "n2-highmem-96"},
			{GuestCpus: 128, MemoryMb: 524288, Name: "n2-standard-128"},
			{GuestCpus: 16, MemoryMb: 65536, Name: "n2-standard-16"},
			{GuestCpus: 2, MemoryMb: 8192, Name: "n2-standard-2"},
			{GuestCpus: 32, MemoryMb: 131072, Name: "n2-standard-32"},
			{GuestCpus: 4, MemoryMb: 16384, Name: "n2-standard-4"},
			{GuestCpus: 48, MemoryMb: 196608, Name: "n2-standard-48"},
			{GuestCpus: 64, MemoryMb: 262144, Name: "n2-standard-64"},
			{GuestCpus: 8, MemoryMb: 32768, Name: "n2-standard-8"},
			{GuestCpus: 80, MemoryMb: 327680, Name: "n2-standard-80"},
			{GuestCpus: 96, MemoryMb: 393216, Name: "n2-standard-96"},
			{GuestCpus: 128, MemoryMb: 131072, Name: "n2d-highcpu-128"},
			{GuestCpus: 16, MemoryMb: 16384, Name: "n2d-highcpu-16"},
			{GuestCpus: 2, MemoryMb: 2048, Name: "n2d-highcpu-2"},
			{GuestCpus: 224, MemoryMb: 229376, Name: "n2d-highcpu-224"},
			{GuestCpus: 32, MemoryMb: 32768, Name: "n2d-highcpu-32"},
			{GuestCpus: 4, MemoryMb: 4096, Name: "n2d-highcpu-4"},
			{GuestCpus: 48, MemoryMb: 49152, Name: "n2d-highcpu-48"},
			{GuestCpus: 64, MemoryMb: 65536, Name: "n2d-highcpu-64"},
			{GuestCpus: 8, MemoryMb: 8192, Name: "n2d-highcpu-8"},
			{GuestCpus: 80, MemoryMb: 81920, Name: "n2d-highcpu-80"},
			{GuestCpus: 96, MemoryMb: 98304, Name: "n2d-highcpu-96"},
			{GuestCpus: 16, MemoryMb: 131072, Name: "n2d-highmem-16"},
			{GuestCpus: 2, MemoryMb: 16384, Name: "n2d-highmem-2"},
			{GuestCpus: 32, MemoryMb: 262144, Name: "n2d-highmem-32"},
			{GuestCpus: 4, MemoryMb: 32768, Name: "n2d-highmem-4"},
			{GuestCpus: 48, MemoryMb: 393216, Name: "n2d-highmem-48"},
			{GuestCpus: 64, MemoryMb: 524288, Name: "n2d-highmem-64"},
			{GuestCpus: 8, MemoryMb: 65536, Name: "n2d-highmem-8"},
			{GuestCpus: 80, MemoryMb: 655360, Name: "n2d-highmem-80"},
			{GuestCpus: 96, MemoryMb: 786432, Name: "n2d-highmem-96"},
			{GuestCpus: 128, MemoryMb: 524288, Name: "n2d-standard-128"},
			{GuestCpus: 16, MemoryMb: 65536, Name: "n2d-standard-16"},
			{GuestCpus: 2, MemoryMb: 8192, Name: "n2d-standard-2"},
			{GuestCpus: 224, MemoryMb: 917504, Name: "n2d-standard-224"},
			{GuestCpus: 32, MemoryMb: 131072, Name: "n2d-standard-32"},
			{GuestCpus: 4, MemoryMb: 16384, Name: "n2d-standard-4"},
			{GuestCpus: 48, MemoryMb: 196608, Name: "n2d-standard-48"},
			{GuestCpus: 64, MemoryMb: 262144, Name: "n2d-standard-64"},
			{GuestCpus: 8, MemoryMb: 32768, Name: "n2d-standard-8"},
			{GuestCpus: 80, MemoryMb: 327680, Name: "n2d-standard-80"},
			{GuestCpus: 96, MemoryMb: 393216, Name: "n2d-standard-96"},
			{GuestCpus: 1, MemoryMb: 4096, Name: "t2a-standard-1"},
			{GuestCpus: 16, MemoryMb: 65536, Name: "t2a-standard-16"},
			{GuestCpus: 2, MemoryMb: 8192, Name: "t2a-standard-2"},
			{GuestCpus: 32, MemoryMb: 131072, Name: "t2a-standard-32"},
			{GuestCpus: 4, MemoryMb: 16384, Name: "t2a-standard-4"},
			{GuestCpus: 48, MemoryMb: 196608, Name: "t2a-standard-48"},
			{GuestCpus: 8, MemoryMb: 32768, Name: "t2a-standard-8"},
			{GuestCpus: 1, MemoryMb: 4096, Name: "t2d-standard-1"},
			{GuestCpus: 16, MemoryMb: 65536, Name: "t2d-standard-16"},
			{GuestCpus: 2, MemoryMb: 8192, Name: "t2d-standard-2"},
			{GuestCpus: 32, MemoryMb: 131072, Name: "t2d-standard-32"},
			{GuestCpus: 4, MemoryMb: 16384, Name: "t2d-standard-4"},
			{GuestCpus: 48, MemoryMb: 196608, Name: "t2d-standard-48"},
			{GuestCpus: 60, MemoryMb: 245760, Name: "t2d-standard-60"},
			{GuestCpus: 8, MemoryMb: 32768, Name: "t2d-standard-8"},
		},
	}
	m.delay()
	return &r, nil
}

func (m mock) MachineTypeFamilyList(imgs *compute.MachineTypeList) gcloud.LabeledValues {
	client := gcloud.NewClient(context.Background(), "deploystack/test")

	return client.MachineTypeFamilyList(imgs)
}

func (m mock) MachineTypeListByFamily(imgs *compute.MachineTypeList, family string) gcloud.LabeledValues {
	client := gcloud.NewClient(context.Background(), "deploystack/test")
	return client.MachineTypeListByFamily(imgs, family)
}

func (m mock) ImageList(project, imageproject string) (*compute.ImageList, error) {
	if m.forceErr {
		return nil, errForced
	}
	imageList := &compute.ImageList{
		Items: []*compute.Image{
			{Name: "centos-7-v20230203 ", Kind: "centos-cloud", Family: "centos-7"},
			{Name: "centos-stream-8-v20230203 ", Kind: "centos-cloud", Family: "centos-stream-8"},
			{Name: "centos-stream-9-v20230203 ", Kind: "centos-cloud", Family: "centos-stream-9"},
			{Name: "cos-101-17162-40-56 ", Kind: "cos-cloud", Family: "cos-101-lts"},
			{Name: "cos-89-16108-798-7", Kind: "cos-cloud", Family: "cos-89-lts"},
			{Name: "cos-93-16623-341-8", Kind: "cos-cloud", Family: "cos-93-lts"},
			{Name: "cos-97-16919-235-9", Kind: "cos-cloud", Family: "cos-97-lts"},
			{Name: "cos-arm64-101-17162-40-56 ", Kind: "cos-cloud", Family: "cos-arm64-101-lts"},
			{Name: "cos-arm64-beta-101-17162-40-56", Kind: "cos-cloud", Family: "cos-arm64-beta"},
			{Name: "cos-arm64-dev-105-17400-0-0 ", Kind: "cos-cloud", Family: "cos-arm64-dev"},
			{Name: "cos-arm64-stable-101-17162-40-56", Kind: "cos-cloud", Family: "cos-arm64-stable"},
			{Name: "cos-beta-101-17162-40-56", Kind: "cos-cloud", Family: "cos-beta"},
			{Name: "debian-10-buster-v20221206", Kind: "debian-cloud ", Family: "debian-10"},
			{Name: "debian-11-bullseye-arm64-v20221102", Kind: "debian-cloud ", Family: "debian-11-arm64"},
			{Name: "debian-11-bullseye-v20221206", Kind: "debian-cloud ", Family: "debian-11"},
			{Name: "fedora-cloud-base-gcp-34-1-2-x86-64 ", Kind: "fedora-cloud ", Family: "fedora-cloud-34"},
			{Name: "fedora-cloud-base-gcp-35-1-2-x86-64 ", Kind: "fedora-cloud ", Family: "fedora-cloud-35"},
			{Name: "fedora-cloud-base-gcp-36-20220506-n-0-x86-64", Kind: "fedora-cloud ", Family: "fedora-cloud-36"},
			{Name: "fedora-cloud-base-gcp-37-beta-1-5-x86-64", Kind: "fedora-cloud ", Family: "fedora-cloud-37"},
			{Name: "opensuse-leap-15-4-v20221201-arm64", Kind: "opensuse-cloud", Family: "opensuse-leap-arm64"},
			{Name: "opensuse-leap-15-4-v20221201-x86-64 ", Kind: "opensuse-cloud", Family: "opensuse-leap"},
			{Name: "rhel-7-v20230203", Kind: "rhel-cloud ", Family: "rhel-7"},
			{Name: "rhel-8-v20230202", Kind: "rhel-cloud ", Family: "rhel-8"},
			{Name: "rhel-9-arm64-v20230203", Kind: "rhel-cloud ", Family: "rhel-9-arm64"},
			{Name: "rhel-9-v20230203", Kind: "rhel-cloud ", Family: "rhel-9"},
			{Name: "rhel-7-7-sap-v20230203", Kind: "rhel-sap-cloud ", Family: "rhel-7-7-sap-ha"},
			{Name: "rhel-7-9-sap-v20230203", Kind: "rhel-sap-cloud ", Family: "rhel-7-9-sap-ha"},
			{Name: "rhel-8-1-sap-v20230202", Kind: "rhel-sap-cloud ", Family: "rhel-8-1-sap-ha"},
			{Name: "rhel-8-2-sap-v20230202", Kind: "rhel-sap-cloud ", Family: "rhel-8-2-sap-ha"},
			{Name: "rhel-8-4-sap-v20230202", Kind: "rhel-sap-cloud ", Family: "rhel-8-4-sap-ha"},
			{Name: "rhel-8-6-sap-v20230202", Kind: "rhel-sap-cloud ", Family: "rhel-8-6-sap-ha"},
			{Name: "rocky-linux-8-optimized-gcp-arm64-v20230202 ", Kind: "rocky-linux-cloud", Family: "rocky-linux-8-optimized-gcp-arm64"},
			{Name: "rocky-linux-8-optimized-gcp-v20230202 ", Kind: "rocky-linux-cloud", Family: "rocky-linux-8-optimized-gcp"},
			{Name: "rocky-linux-8-v20230202 ", Kind: "rocky-linux-cloud", Family: "rocky-linux-8"},
			{Name: "rocky-linux-9-arm64-v20230203 ", Kind: "rocky-linux-cloud", Family: "rocky-linux-9-arm64"},
			{Name: "rocky-linux-9-optimized-gcp-arm64-v20230203 ", Kind: "rocky-linux-cloud", Family: "rocky-linux-9-optimized-gcp-arm64"},
			{Name: "rocky-linux-9-optimized-gcp-v20230203 ", Kind: "rocky-linux-cloud", Family: "rocky-linux-9-optimized-gcp"},
			{Name: "rocky-linux-9-v20230203 ", Kind: "rocky-linux-cloud", Family: "rocky-linux-9"},
			{Name: "sles-12-sp5-v20221104-x86-64", Kind: "suse-cloud", Family: "sles-12"},
			{Name: "sles-15-sp4-v20221104-arm64 ", Kind: "suse-cloud", Family: "sles-15-arm64"},
			{Name: "sles-15-sp4-v20221104-x86-64", Kind: "suse-cloud", Family: "sles-15"},
			{Name: "sles-12-sp5-sap-v20230116-x86-64", Kind: "suse-sap-cloud ", Family: "sles-12-sp5-sap"},
			{Name: "sles-15-sp1-sap-v20221108-x86-64", Kind: "suse-sap-cloud ", Family: "sles-15-sp1-sap"},
			{Name: "sles-15-sp2-sap-v20221108-x86-64", Kind: "suse-sap-cloud ", Family: "sles-15-sp2-sap"},
			{Name: "sles-15-sp3-sap-v20221108-x86-64", Kind: "suse-sap-cloud ", Family: "sles-15-sp3-sap"},
			{Name: "sles-15-sp4-sap-v20221104-x86-64", Kind: "suse-sap-cloud ", Family: "sles-15-sp4-sap"},
			{Name: "ubuntu-1804-bionic-arm64-v20230131", Kind: "ubuntu-os-cloud", Family: "ubuntu-1804-lts-arm64"},
			{Name: "ubuntu-pro-1604-xenial-v20221201", Kind: "ubuntu-os-pro-cloud", Family: "ubuntu-pro-1604-lts"},
			{Name: "ubuntu-pro-1804-bionic-v20230124", Kind: "ubuntu-os-pro-cloud", Family: "ubuntu-pro-1804-lts"},
			{Name: "ubuntu-pro-2004-focal-v20230126 ", Kind: "ubuntu-os-pro-cloud", Family: "ubuntu-pro-2004-lts"},
			{Name: "ubuntu-pro-2204-jammy-v20230114 ", Kind: "ubuntu-os-pro-cloud", Family: "ubuntu-pro-2204-lts"},
			{Name: "ubuntu-pro-fips-1804-bionic-v20230131 ", Kind: "ubuntu-os-pro-cloud", Family: "ubuntu-pro-fips-1804-lts"},
			{Name: "ubuntu-pro-fips-2004-focal-v20230126", Kind: "ubuntu-os-pro-cloud", Family: "ubuntu-pro-fips-2004-lts"},
			{Name: "windows-server-2012-r2-dc-core-v20230113", Kind: "windows-cloud", Family: "windows-2012-r2-core"},
			{Name: "windows-server-2012-r2-dc-v20230112 ", Kind: "windows-cloud", Family: "windows-2012-r2"},
			{Name: "sql-2014-enterprise-windows-2012-r2-dc-v20230112", Kind: "windows-sql-cloud", Family: "sql-ent-2014-win-2012-r2"},
			{Name: "sql-2014-enterprise-windows-2016-dc-v20230111 ", Kind: "windows-sql-cloud", Family: "sql-ent-2014-win-2016"},
			{Name: "sql-2014-standard-windows-2012-r2-dc-v20230112", Kind: "windows-sql-cloud", Family: "sql-std-2014-win-2012-r2"},
			{Name: "cos-dev-105-17400-0-0 ", Kind: "cos-cloud", Family: "cos-dev"},
			{Name: "cos-stable-101-17162-40-56", Kind: "cos-cloud", Family: "cos-stable"},
			{Name: "ubuntu-1804-bionic-v20230131", Kind: "ubuntu-os-cloud", Family: "ubuntu-1804-lts"},
			{Name: "ubuntu-2004-focal-arm64-v20230125 ", Kind: "ubuntu-os-cloud", Family: "ubuntu-2004-lts-arm64"},
			{Name: "ubuntu-2004-focal-v20230125 ", Kind: "ubuntu-os-cloud", Family: "ubuntu-2004-lts"},
			{Name: "ubuntu-2204-jammy-arm64-v20230114 ", Kind: "ubuntu-os-cloud", Family: "ubuntu-2204-lts-arm64"},
			{Name: "ubuntu-2204-jammy-v20230114 ", Kind: "ubuntu-os-cloud", Family: "ubuntu-2204-lts"},
			{Name: "ubuntu-2210-kinetic-amd64-v20230125 ", Kind: "ubuntu-os-cloud", Family: "ubuntu-2210-amd64"},
			{Name: "ubuntu-2210-kinetic-arm64-v20230125 ", Kind: "ubuntu-os-cloud", Family: "ubuntu-2210-arm64"},
			{Name: "ubuntu-minimal-1804-bionic-arm64-v20230125", Kind: "ubuntu-os-cloud", Family: "ubuntu-minimal-1804-lts-arm64"},
			{Name: "ubuntu-minimal-1804-bionic-v20230125", Kind: "ubuntu-os-cloud", Family: "ubuntu-minimal-1804-lts"},
			{Name: "windows-server-2016-dc-core-v20230111 ", Kind: "windows-cloud", Family: "windows-2016-core"},
			{Name: "windows-server-2016-dc-v20230111", Kind: "windows-cloud", Family: "windows-2016"},
			{Name: "windows-server-2019-dc-core-for-containers-v20230113", Kind: "windows-cloud", Family: "windows-2019-core-for-containers"},
			{Name: "windows-server-2019-dc-core-v20230111 ", Kind: "windows-cloud", Family: "windows-2019-core"},
			{Name: "windows-server-2019-dc-for-containers-v20230113 ", Kind: "windows-cloud", Family: "windows-2019-for-containers"},
			{Name: "windows-server-2019-dc-v20230111", Kind: "windows-cloud", Family: "windows-2019"},
			{Name: "windows-server-2022-dc-core-v20230111 ", Kind: "windows-cloud", Family: "windows-2022-core"},
			{Name: "windows-server-2022-dc-v20230111", Kind: "windows-cloud", Family: "windows-2022"},
			{Name: "sql-2014-web-windows-2012-r2-dc-v20230112 ", Kind: "windows-sql-cloud", Family: "sql-web-2014-win-2012-r2"},
			{Name: "sql-2016-enterprise-windows-2012-r2-dc-v20230112", Kind: "windows-sql-cloud", Family: "sql-ent-2016-win-2012-r2"},
			{Name: "sql-2016-enterprise-windows-2016-dc-v20230111 ", Kind: "windows-sql-cloud", Family: "sql-ent-2016-win-2016"},
			{Name: "sql-2016-enterprise-windows-2019-dc-v20230111 ", Kind: "windows-sql-cloud", Family: "sql-ent-2016-win-2019"},
			{Name: "sql-2016-standard-windows-2012-r2-dc-v20230112", Kind: "windows-sql-cloud", Family: "sql-std-2016-win-2012-r2"},
			{Name: "sql-2016-standard-windows-2016-dc-v20230111 ", Kind: "windows-sql-cloud", Family: "sql-std-2016-win-2016"},
			{Name: "sql-2016-standard-windows-2019-dc-v20230111 ", Kind: "windows-sql-cloud", Family: "sql-std-2016-win-2019"},
			{Name: "ubuntu-minimal-2004-focal-arm64-v20230126 ", Kind: "ubuntu-os-cloud", Family: "ubuntu-minimal-2004-lts-arm64"},
			{Name: "ubuntu-minimal-2004-focal-v20230126 ", Kind: "ubuntu-os-cloud", Family: "ubuntu-minimal-2004-lts"},
			{Name: "ubuntu-minimal-2204-jammy-arm64-v20230124 ", Kind: "ubuntu-os-cloud", Family: "ubuntu-minimal-2204-lts-arm64"},
			{Name: "ubuntu-minimal-2204-jammy-v20230124 ", Kind: "ubuntu-os-cloud", Family: "ubuntu-minimal-2204-lts"},
			{Name: "ubuntu-minimal-2210-kinetic-amd64-v20230126 ", Kind: "ubuntu-os-cloud", Family: "ubuntu-minimal-2210-amd64"},
			{Name: "ubuntu-minimal-2210-kinetic-arm64-v20230126 ", Kind: "ubuntu-os-cloud", Family: "ubuntu-minimal-2210-arm64"},
			{Name: "sql-2016-web-windows-2012-r2-dc-v20230112 ", Kind: "windows-sql-cloud", Family: "sql-web-2016-win-2012-r2"},
			{Name: "sql-2016-web-windows-2016-dc-v20230111", Kind: "windows-sql-cloud", Family: "sql-web-2016-win-2016"},
			{Name: "sql-2016-web-windows-2019-dc-v20230111", Kind: "windows-sql-cloud", Family: "sql-web-2016-win-2019"},
			{Name: "sql-2017-enterprise-windows-2016-dc-v20230111 ", Kind: "windows-sql-cloud", Family: "sql-ent-2017-win-2016"},
			{Name: "sql-2017-enterprise-windows-2019-dc-v20230111 ", Kind: "windows-sql-cloud", Family: "sql-ent-2017-win-2019"},
			{Name: "sql-2017-enterprise-windows-2022-dc-v20230111 ", Kind: "windows-sql-cloud", Family: "sql-ent-2017-win-2022"},
			{Name: "sql-2017-express-windows-2012-r2-dc-v20230112 ", Kind: "windows-sql-cloud", Family: "sql-exp-2017-win-2012-r2"},
			{Name: "sql-2017-express-windows-2016-dc-v20230111", Kind: "windows-sql-cloud", Family: "sql-exp-2017-win-2016"},
			{Name: "sql-2017-express-windows-2019-dc-v20230111", Kind: "windows-sql-cloud", Family: "sql-exp-2017-win-2019"},
			{Name: "fedora-coreos-37-20230110-3-1-gcp-x86-64", Kind: "fedora-coreos-cloud", Family: "fedora-coreos-stable"},
			{Name: "fedora-coreos-37-20230122-1-1-gcp-x86-64", Kind: "fedora-coreos-cloud", Family: "fedora-coreos-next"},
			{Name: "fedora-coreos-37-20230122-2-0-gcp-x86-64", Kind: "fedora-coreos-cloud", Family: "fedora-coreos-testing"},
			{Name: "sql-2017-standard-windows-2016-dc-v20230111 ", Kind: "windows-sql-cloud", Family: "sql-std-2017-win-2016"},
			{Name: "sql-2017-standard-windows-2019-dc-v20230111 ", Kind: "windows-sql-cloud", Family: "sql-std-2017-win-2019"},
			{Name: "sql-2017-standard-windows-2022-dc-v20230111 ", Kind: "windows-sql-cloud", Family: "sql-std-2017-win-2022"},
			{Name: "sql-2017-web-windows-2016-dc-v20230111", Kind: "windows-sql-cloud", Family: "sql-web-2017-win-2016"},
			{Name: "sql-2017-web-windows-2019-dc-v20230111", Kind: "windows-sql-cloud", Family: "sql-web-2017-win-2019"},
			{Name: "sql-2017-web-windows-2022-dc-v20230111", Kind: "windows-sql-cloud", Family: "sql-web-2017-win-2022"},
			{Name: "sql-2019-enterprise-windows-2019-dc-v20230111 ", Kind: "windows-sql-cloud", Family: "sql-ent-2019-win-2019"},
			{Name: "sql-2019-enterprise-windows-2022-dc-v20230111 ", Kind: "windows-sql-cloud", Family: "sql-ent-2019-win-2022"},
			{Name: "sql-2019-standard-windows-2019-dc-v20230111 ", Kind: "windows-sql-cloud", Family: "sql-std-2019-win-2019"},
			{Name: "sql-2019-standard-windows-2022-dc-v20230111 ", Kind: "windows-sql-cloud", Family: "sql-std-2019-win-2022"},
			{Name: "sql-2019-web-windows-2019-dc-v20230111", Kind: "windows-sql-cloud", Family: "sql-web-2019-win-2019"},
			{Name: "sql-2019-web-windows-2022-dc-v20230111", Kind: "windows-sql-cloud", Family: "sql-web-2019-win-2022"},
			{Name: "sql-2022-enterprise-windows-2019-dc-v20230112 ", Kind: "windows-sql-cloud", Family: "sql-ent-2022-win-2019"},
			{Name: "sql-2022-enterprise-windows-2022-dc-v20230112 ", Kind: "windows-sql-cloud", Family: "sql-ent-2022-win-2022"},
			{Name: "sql-2022-standard-windows-2019-dc-v20230112 ", Kind: "windows-sql-cloud", Family: "sql-std-2022-win-2019"},
			{Name: "sql-2022-standard-windows-2022-dc-v20230112 ", Kind: "windows-sql-cloud", Family: "sql-std-2022-win-2022"},
			{Name: "sql-2022-web-windows-2019-dc-v20230112", Kind: "windows-sql-cloud", Family: "sql-web-2022-win-2019"},
			{Name: "sql-2022-web-windows-2022-dc-v20230112", Kind: "windows-sql-cloud", Family: "sql-web-2022-win-2022"},
		},
	}

	resp := &compute.ImageList{}

	for _, v := range imageList.Items {
		if strings.TrimSpace(v.Kind) == strings.TrimSpace(imageproject) {
			resp.Items = append(resp.Items, v)
		}
	}

	return resp, nil
}

func (m mock) ImageTypeListByFamily(imgs *compute.ImageList, project, family string) gcloud.LabeledValues {
	lb := gcloud.LabeledValues{}

	for _, v := range imgs.Items {
		if v.Family == family {
			value := fmt.Sprintf("%s/%s", project, v.Name)
			lb = append(lb, gcloud.LabeledValue{Value: value, Label: v.Name, IsDefault: false})
		}
	}

	last := lb[len(lb)-1]
	last.Label = fmt.Sprintf("%s (Latest)", last.Label)
	lb[len(lb)-1] = last
	lb.Sort()
	lb.SetDefault(last.Value)

	return lb
}

func (m mock) ProjectNumberGet(id string) (string, error) {
	if m.forceErr {
		return "", errForced
	}
	return "123234567755", nil
}

func (m mock) ImageFamilyList(imgs *compute.ImageList) gcloud.LabeledValues {
	fam := make(map[string]bool)
	lb := gcloud.LabeledValues{}

	for _, v := range imgs.Items {
		fam[v.Family] = false
	}

	for i := range fam {
		if i == "" {
			continue
		}
		lb = append(lb, gcloud.LabeledValue{
			Value:     i,
			Label:     i,
			IsDefault: false,
		})
	}
	lb.SetDefault(gcloud.DefaultImageFamily)
	lb.Sort()
	return lb
}

func (m *mock) save(key string, value interface{}) {
	if m.cache == nil {
		m.cache = make(map[string]interface{})
	}

	m.cache[key] = value
}

func (m *mock) get(key string) interface{} {
	return m.cache[key]
}

func (m mock) BillingAccountList() ([]*cloudbilling.BillingAccount, error) {
	if m.forceErr {
		return nil, errForced
	}

	result := []*cloudbilling.BillingAccount{
		{
			DisplayName: "Very Limted Funds",
			Name:        "billingAccounts/000000-000000-00000Y",
		},
		{
			DisplayName: "Unlimted Funds",
			Name:        "billingAccounts/000000-000000-00000X",
		},
	}

	i := m.get("BillingAccountList")
	switch val := i.(type) {
	case []*cloudbilling.BillingAccount:
		return val, nil
	}

	return result, nil
}

var errForced = fmt.Errorf("this is a forced error for mocking")

func (m mock) BillingAccountAttach(project, account string) error {
	if m.forceErr {
		return errForced
	}
	return nil
}
