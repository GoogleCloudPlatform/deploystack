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

	"google.golang.org/api/iam/v1"
)

func (c *Client) getIAMService(project string) (*iam.Service, error) {
	var err error
	svc := c.services.iam

	if svc != nil {
		return svc, nil
	}

	if err := c.ServiceEnable(project, "domains.googleapis.com"); err != nil {
		return nil, fmt.Errorf("error activating service for polling: %s", err)
	}

	svc, err = iam.NewService(c.ctx, c.opts)
	if err != nil {
		return nil, err
	}

	svc.UserAgent = c.userAgent
	c.services.iam = svc

	return svc, nil
}

// ServiceAccountCreate creates a service account. A little on the nose
func (c *Client) ServiceAccountCreate(project, username, displayName string) (string, error) {
	svc, err := c.getIAMService(project)
	if err != nil {
		return "", err
	}

	req := &iam.CreateServiceAccountRequest{
		AccountId: username,
		ServiceAccount: &iam.ServiceAccount{
			DisplayName: displayName,
		},
	}

	servicaccount, err := svc.Projects.ServiceAccounts.Create(fmt.Sprintf("projects/%s", project), req).Do()
	if err != nil {
		return "", err
	}

	return servicaccount.Email, nil
}

// ServiceAccountDelete deletes a service account. A little on the nose
func (c *Client) ServiceAccountDelete(project, email string) error {
	svc, err := c.getIAMService(project)
	if err != nil {
		return err
	}

	name := fmt.Sprintf("projects/%s/serviceAccounts/%s", project, email)
	_, err = svc.Projects.ServiceAccounts.Delete(name).Do()

	return err
}
