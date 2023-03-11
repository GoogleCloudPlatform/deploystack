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
	b64 "encoding/base64"
	"fmt"

	"google.golang.org/api/secretmanager/v1"
)

func (c *Client) getSecretManagerService(project string) (*secretmanager.Service, error) {
	var err error
	svc := c.services.secretManager

	if svc != nil {
		return svc, nil
	}

	if err := c.ServiceEnable(project, SecretManager); err != nil {
		return nil, fmt.Errorf("error activating service for polling: %s", err)
	}

	svc, err = secretmanager.NewService(c.ctx, c.opts)
	if err != nil {
		return nil, err
	}

	svc.UserAgent = c.userAgent
	c.services.secretManager = svc

	return svc, nil
}

// SecretCreate creates a secret and populates the lastest version with a payload.
func (c *Client) SecretCreate(project, name, payload string) error {
	svc, err := c.getSecretManagerService(project)
	if err != nil {
		return err
	}

	secret := &secretmanager.Secret{
		Name: fmt.Sprintf("projects/%s/secrets/%s", project, name),
		Replication: &secretmanager.Replication{
			Automatic: &secretmanager.Automatic{},
		},
	}

	parent := fmt.Sprintf("projects/%s", project)

	req := svc.Projects.Secrets.Create(parent, secret)
	req.SecretId(name)

	result, err := req.Do()
	if err != nil {
		return fmt.Errorf("failed to create secret: %s", err)
	}

	version := &secretmanager.AddSecretVersionRequest{
		Payload: &secretmanager.SecretPayload{
			Data: b64.URLEncoding.EncodeToString([]byte(payload)),
		},
	}

	if _, err := svc.Projects.Secrets.AddVersion(result.Name, version).Do(); err != nil {
		return fmt.Errorf("failed to create secret versiopn: %s", err)
	}

	return nil
}

// SecretDelete deletes a secret
func (c *Client) SecretDelete(project, name string) error {
	svc, err := c.getSecretManagerService(project)
	if err != nil {
		return err
	}

	secret := fmt.Sprintf("projects/%s/secrets/%s", project, name)
	if _, err := svc.Projects.Secrets.Delete(secret).Do(); err != nil {
		return fmt.Errorf("could not delete secret (%s) in project (%s)", name, project)
	}

	return nil
}
