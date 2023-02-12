package gcloud

import (
	"context"
	"fmt"
	"sort"
	"strings"

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

var (
	// DefaultRegion is the default compute region used in compute calls.
	DefaultRegion = "us-central1"
	// DefaultMachineType is the default compute machine type used in compute calls.
	DefaultMachineType = "n1-standard"
	// DefaultImageProject is the default project for images used in compute calls.
	DefaultImageProject = "debian-cloud"
	// DefaultImageFamily is the default project for images used in compute calls.
	DefaultImageFamily = "debian-11"
	// DefaultDiskSize is the default size for making disks for Compute Engine
	DefaultDiskSize = "200"
	// DefaultDiskType is the default style of disk
	DefaultDiskType = "pd-standard"
	// DefaultInstanceType is the default machine type of compute engine
	DefaultInstanceType = "n1-standard-1"
	// HTTPServerTags are the instance tags to open up the instance to be a
	// http server
	HTTPServerTags = "[http-server,https-server]"
	// DefaultZone is the default zone used in compute calls.
	DefaultZone = "us-central1-a"
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

// LabeledValue is a struct that contains a label/value pair
type LabeledValue struct {
	Value     string
	Label     string
	IsDefault bool
}

// NewLabeledValue takes a string and converts it to a LabeledValue. If a |
// delimiter is present it will split into a different label/value
func NewLabeledValue(s string) LabeledValue {
	l := LabeledValue{s, s, false}

	if strings.Contains(s, "|") {
		sl := strings.Split(s, "|")
		l = LabeledValue{sl[0], sl[1], false}
	}

	return l
}

// LabeledValues is collection of LabledValue structs
type LabeledValues []LabeledValue

// Sort orders the LabeledValues by Label
func (l *LabeledValues) Sort() {
	sort.Slice(*l, func(i, j int) bool {
		iStr := strings.ToLower((*l)[i].Label)
		jStr := strings.ToLower((*l)[j].Label)
		return iStr < jStr
	})
}

// LongestLen returns the length of longest LABEL in the list
func (l *LabeledValues) LongestLen() int {
	longest := 0

	for _, v := range *l {
		if len(v.Label) > longest {
			longest = len(v.Label)
		}
	}

	return longest
}

// GetDefault returns the deafult value of the LabeledValues list
func (l *LabeledValues) GetDefault() LabeledValue {
	for _, v := range *l {
		if v.IsDefault {
			return v
		}
	}
	return LabeledValue{}
}

// SetDefault sets the default value of the list
func (l *LabeledValues) SetDefault(value string) {
	for i, v := range *l {
		if v.Value == value {
			v.IsDefault = true
			(*l)[i] = v
		}
	}
}

// NewLabeledValues takes a slice of strings and returns a list of LabeledValues
func NewLabeledValues(sl []string, defaultValue string) LabeledValues {
	r := LabeledValues{}

	for _, v := range sl {
		val := NewLabeledValue(v)
		if val.Value == defaultValue {
			val.IsDefault = true
		}

		r = append(r, val)
	}
	return r
}
