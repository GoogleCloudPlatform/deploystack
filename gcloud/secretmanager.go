package gcloud

import (
	b64 "encoding/base64"
	"fmt"

	"google.golang.org/api/secretmanager/v1"
)

func (c *Client) getSecretManagerService(project string) (*secretmanager.Service, error) {
	var err error
	svc := c.services.secretManagerService

	if svc != nil {
		return svc, nil
	}

	if err := c.ServiceEnable(project, "secretmanager.googleapis.com"); err != nil {
		return nil, fmt.Errorf("error activating service for polling: %s", err)
	}

	svc, err = secretmanager.NewService(c.ctx, c.opts)
	if err != nil {
		return nil, err
	}

	svc.UserAgent = c.userAgent
	c.services.secretManagerService = svc

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
