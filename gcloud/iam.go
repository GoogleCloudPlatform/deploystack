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
