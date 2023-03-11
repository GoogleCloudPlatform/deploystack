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

	"google.golang.org/api/cloudbuild/v1"
)

func (c *Client) getCloudBuildService(project string) (*cloudbuild.Service, error) {
	var err error
	svc := c.services.build

	if svc != nil {
		return svc, nil
	}

	if err := c.ServiceEnable(project, CloudBuild); err != nil {
		return nil, fmt.Errorf("error activating service for polling: %s", err)
	}

	svc, err = cloudbuild.NewService(c.ctx, c.opts)
	if err != nil {
		return nil, err
	}

	svc.UserAgent = c.userAgent
	c.services.build = svc

	return svc, nil
}

// CloudBuildTriggerCreate creates a build trigger in a given project
func (c *Client) CloudBuildTriggerCreate(project string, trigger cloudbuild.BuildTrigger) (*cloudbuild.BuildTrigger, error) {
	svc, err := c.getCloudBuildService(project)
	if err != nil {
		return nil, err
	}

	req := svc.Projects.Triggers.Create(project, &trigger)
	result, err := req.Do()
	if err != nil {
		return nil, fmt.Errorf("cannot create trigger: %s", err)
	}

	return result, nil
}

// CloudBuildTriggerDelete deletes a build trigger in a given project
func (c *Client) CloudBuildTriggerDelete(project string, triggerid string) error {
	svc, err := c.getCloudBuildService(project)
	if err != nil {
		return err
	}

	if _, err := svc.Projects.Triggers.Delete(project, triggerid).Do(); err != nil {
		return fmt.Errorf("cannot delete trigger: %s", err)
	}

	return nil
}
