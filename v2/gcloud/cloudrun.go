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
	"sort"

	"google.golang.org/api/run/v1"
)

func (c *Client) getRunService(project string) (*run.APIService, error) {
	var err error
	svc := c.services.run

	if svc != nil {
		return svc, nil
	}

	if err := c.ServiceEnable(project, "run.googleapis.com"); err != nil {
		return nil, fmt.Errorf("error activating service for polling: %s", err)
	}

	svc, err = run.NewService(c.ctx, c.opts)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve service: %w", err)
	}

	svc.UserAgent = c.userAgent
	c.services.run = svc

	return svc, nil
}

// RunRegionList will return a list of regions for Cloud Run
func (c *Client) RunRegionList(project string) ([]string, error) {
	resp := []string{}

	svc, err := c.getRunService(project)
	if err != nil {
		return resp, err
	}

	results, err := svc.Projects.Locations.List("projects/" + project).Do()
	if err != nil {
		return resp, err
	}

	for _, v := range results.Locations {
		resp = append(resp, v.LocationId)
	}

	sort.Strings(resp)

	return resp, nil
}
