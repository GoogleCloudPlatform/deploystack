package gcloud

import (
	"context"
	"fmt"

	domains "cloud.google.com/go/domains/apiv1beta1"
	scheduler "cloud.google.com/go/scheduler/apiv1beta1"
	"cloud.google.com/go/storage"
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
)

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
func NewClient(ctx context.Context, ua string) Client {
	c := Client{}
	c.ctx = ctx
	c.userAgent = ua
	c.opts = option.WithCredentialsFile("")
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
