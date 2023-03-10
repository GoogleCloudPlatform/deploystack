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

package gcloud

import (
	"fmt"
	"strings"
	"time"

	"google.golang.org/api/serviceusage/v1"
)

// ErrorServiceNotExistOrNotAllowed occurs when the user running this code doesn't have
// permission to enable the service in the project or it's a nonexistent service name.
var ErrorServiceNotExistOrNotAllowed = fmt.Errorf("Not found or permission denied for service")

// ErrorProjectRequired communicates that am empty project string has been passed
var ErrorProjectRequired = fmt.Errorf("Project may not be an empty string")

func (c *Client) getServiceUsageService() (*serviceusage.Service, error) {
	var err error
	svc := c.services.serviceUsage

	if svc != nil {
		return svc, nil
	}

	svc, err = serviceusage.NewService(c.ctx, c.opts)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve service: %w", err)
	}

	svc.UserAgent = c.userAgent
	c.services.serviceUsage = svc

	return svc, nil
}

// ServiceEnable enable a service in the selected project so that query calls
// to various lists will work.
func (c *Client) ServiceEnable(project, service string) error {
	if _, ok := c.enabledServices[service]; ok {
		return nil
	}

	svc, err := c.getServiceUsageService()
	if err != nil {
		return fmt.Errorf("could not getServiceUsageService: %s", err)
	}

	enabled, err := c.ServiceIsEnabled(project, service)
	if err != nil {
		return fmt.Errorf("could not confirm if service is already enabled: %w", err)
	}

	if enabled {
		c.enabledServices[service] = true
		return nil
	}

	s := fmt.Sprintf("projects/%s/services/%s", project, service)
	op, err := svc.Services.Enable(s, &serviceusage.EnableServiceRequest{}).Do()
	if err != nil {
		return fmt.Errorf("could not enable service: %s", err)
	}

	if !strings.Contains(string(op.Response), "ENABLED") {
		for i := 0; i < 60; i++ {
			enabled, err = c.ServiceIsEnabled(project, service)
			if err != nil {
				return err
			}
			if enabled {
				c.enabledServices[service] = true
				return nil
			}
			time.Sleep(1 * time.Second)
		}
	}

	c.enabledServices[service] = true
	return nil
}

// ServiceIsEnabled checks to see if the existing service is already enabled
// in the project we are trying to enable it in.
func (c *Client) ServiceIsEnabled(project, service string) (bool, error) {
	svc, err := c.getServiceUsageService()

	if project == "" {
		return false, ErrorProjectRequired
	}

	s := fmt.Sprintf("projects/%s/services/%s", project, service)
	current, err := svc.Services.Get(s).Do()
	if err != nil {
		if strings.Contains(err.Error(), "Not found or permission denied for service") {
			return false, ErrorServiceNotExistOrNotAllowed
		}

		return false, fmt.Errorf("cannot get the service for resource (%s): %w", s, err)
	}

	if current.State == "ENABLED" {
		return true, nil
	}

	return false, nil
}

// ServiceDisable disables a service in the selected project
func (c *Client) ServiceDisable(project, service string) error {
	svc, err := c.getServiceUsageService()
	if err != nil {
		return err
	}
	s := fmt.Sprintf("projects/%s/services/%s", project, service)
	if _, err := svc.Services.Disable(s, &serviceusage.DisableServiceRequest{}).Do(); err != nil {
		if strings.Contains(err.Error(), "Not found or permission denied for service") {
			return ErrorServiceNotExistOrNotAllowed
		}
		return fmt.Errorf("could not disable service: %s", err)
	}

	return nil
}
