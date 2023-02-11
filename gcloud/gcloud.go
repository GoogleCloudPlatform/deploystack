package gcloud

import (
	"context"
	"fmt"

	domains "cloud.google.com/go/domains/apiv1beta1"
	scheduler "cloud.google.com/go/scheduler/apiv1beta1"
	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/deploystack"
	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/cloudbuild/v1"
	"google.golang.org/api/cloudfunctions/v1"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
	"google.golang.org/api/run/v1"
	"google.golang.org/api/secretmanager/v1"
	"google.golang.org/api/serviceusage/v1"
	domainspb "google.golang.org/genproto/googleapis/cloud/domains/v1beta1"
)

// UIClient interface encapsulates all of the calls to gcloud that one needs to
// make the TUI work
type UIClient interface {
	ProjectIDGet() (string, error)
	ProjectList() ([]ProjectWithBilling, error)
	ProjectParentGet(project string) (*cloudresourcemanager.ResourceId, error)
	ProjectCreate(project, parent, parentType string) error
	RegionList(project, product string) ([]string, error)
	ZoneList(project, region string) ([]string, error)
	DomainIsAvailable(project, domain string) (*domainspb.RegisterParameters, error)
	DomainIsVerified(project, domain string) (bool, error)
	DomainRegister(project string, domaininfo *domainspb.RegisterParameters, contact ContactData) error
	ComputeImageLatestGet(project, imageproject, imagefamily string) (string, error)
	ComputeMachineTypeList(project, zone string) (*compute.MachineTypeList, error)
	ComputeMachineTypeFamilyList(imgs *compute.MachineTypeList) deploystack.LabeledValues
	ComputeMachineTypeListByFamily(imgs *compute.MachineTypeList, family string) deploystack.LabeledValues
	ComputeImageList(project, imageproject string) (*compute.ImageList, error)
	ComputeImageTypeListByFamily(imgs *compute.ImageList, project, family string) deploystack.LabeledValues
}

// Client is the tool that will handle all of the communication between gcloud
// and the various product areas
type Client struct {
	ctx             context.Context
	services        services
	userAgent       string
	opts            option.ClientOption
	enabledServices map[string]bool
}

// NewClient initiates a new gcloud Client
func NewClient(ctx context.Context, ua string, opts option.ClientOption) Client {
	c := Client{}
	c.ctx = ctx
	c.userAgent = ua
	c.opts = opts
	c.enabledServices = make(map[string]bool)
	return c
}

type services struct {
	resourceManager *cloudresourcemanager.Service
	billing         *cloudbilling.APIService
	domains         *domains.Client
	serviceUsage    *serviceusage.Service
	computeService  *compute.Service
	functions       *cloudfunctions.Service
	run             *run.APIService
	build           *cloudbuild.Service
	iam             *iam.Service
	scheduler       *scheduler.CloudSchedulerClient
	secretManager   *secretmanager.Service
	storage         *storage.Client
}

// RegionList will return a list of RegionsList depending on product type
func (c *Client) RegionList(project, product string) ([]string, error) {
	switch product {
	case "compute":
		return c.ComputeRegionList(project)
	case "functions":
		return c.FunctionRegionList(project)
	case "run":
		return c.RunRegionList(project)
	}

	return []string{}, fmt.Errorf("invalid product (%s) requested", product)
}
