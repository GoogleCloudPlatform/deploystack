package tui

import (
	"cloud.google.com/go/domains/apiv1beta1/domainspb"
	"github.com/GoogleCloudPlatform/deploystack"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/compute/v1"
)

// TODO: move this to package under deploystack deploystack/gcloud?
type client interface {
	ProjectIDGet() (string, error)
	ProjectsList() ([]deploystack.ProjectWithBilling, error)
	ProjectParentGet(project string) (*cloudresourcemanager.ResourceId, error)
	ProjectCreate(project, parent, parentType string) error
	RegionList(project, product string) ([]string, error)
	ZoneList(project, region string) ([]string, error)
	DomainIsAvailable(project, domain string) (*domainspb.RegisterParameters, error)
	DomainIsVerified(project, domain string) (bool, error)
	DomainRegister(project string, domaininfo *domainspb.RegisterParameters, contact deploystack.ContactData) error
	ComputeImageLatestGet(project, imageproject, imagefamily string) (string, error)
	ComputeMachineTypeList(project, zone string) (*compute.MachineTypeList, error)
	ComputeMachineTypeFamilyList(imgs *compute.MachineTypeList) deploystack.LabeledValues
	ComputeMachineTypeListByFamily(imgs *compute.MachineTypeList, family string) deploystack.LabeledValues
	ComputeImageList(project, imageproject string) (*compute.ImageList, error)
	ComputeImageTypeListByFamily(imgs *compute.ImageList, project, family string) deploystack.LabeledValues
}
