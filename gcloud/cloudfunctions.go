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

	"google.golang.org/api/cloudfunctions/v1"
)

func (c *Client) getCloudFunctionsService(project string) (*cloudfunctions.Service, error) {
	var err error
	svc := c.services.functions

	if svc != nil {
		return svc, nil
	}

	if err := c.ServiceEnable(project, CloudFunctions); err != nil {
		return nil, fmt.Errorf("error activating service for polling: %s", err)
	}

	svc, err = cloudfunctions.NewService(c.ctx, c.opts)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve service: %w", err)
	}

	svc.UserAgent = c.userAgent
	c.services.functions = svc

	return svc, nil
}

// FunctionRegionList will return a list of regions for Cloud Functions
func (c *Client) FunctionRegionList(project string) ([]string, error) {
	resp := []string{}

	if err := c.ServiceEnable(project, CloudFunctions); err != nil {
		return resp, fmt.Errorf("error activating service for polling: %s", err)
	}

	svc, err := c.getCloudFunctionsService(project)
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

// FunctionDeploy deploys a Cloud Function.
func (c *Client) FunctionDeploy(project, region string, f cloudfunctions.CloudFunction) error {
	svc, err := c.getCloudFunctionsService(project)
	if err != nil {
		return err
	}

	location := fmt.Sprintf("projects/%s/locations/%s", project, region)
	if _, err := svc.Projects.Locations.Functions.Create(location, &f).Do(); err != nil {
		return fmt.Errorf("could not create function: %s", err)
	}

	return nil
}

// FunctionDelete deletes a Cloud Function.
func (c *Client) FunctionDelete(project, region, name string) error {
	svc, err := c.getCloudFunctionsService(project)
	if err != nil {
		return err
	}
	fname := fmt.Sprintf("projects/%s/locations/%s/functions/%s", project, region, name)
	if _, err := svc.Projects.Locations.Functions.Delete(fname).Do(); err != nil {
		return fmt.Errorf("could not create function: %s", err)
	}

	return nil
}

// FunctionGet gets the details of a Cloud Function.
func (c *Client) FunctionGet(project, region, name string) (*cloudfunctions.CloudFunction, error) {
	svc, err := c.getCloudFunctionsService(project)
	if err != nil {
		return nil, err
	}

	fname := fmt.Sprintf("projects/%s/locations/%s/functions/%s", project, region, name)
	result, err := svc.Projects.Locations.Functions.Get(fname).Do()
	if err != nil {
		return nil, fmt.Errorf("could not get function: %s", err)
	}

	return result, nil
}

// FunctionGenerateSignedURL generates a signed url for use with uploading to
// Cloud Storage
func (c *Client) FunctionGenerateSignedURL(project, region string) (string, error) {
	location := fmt.Sprintf("projects/%s/locations/%s", project, region)
	svc, err := c.getCloudFunctionsService(project)
	if err != nil {
		return "", err
	}

	req := &cloudfunctions.GenerateUploadUrlRequest{}

	result, err := svc.Projects.Locations.Functions.GenerateUploadUrl(location, req).Do()
	if err != nil {
		return "", err
	}

	return result.UploadUrl, nil
}
