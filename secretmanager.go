package deploystack

import (
	"context"
	b64 "encoding/base64"
	"fmt"

	"google.golang.org/api/secretmanager/v1"
)

var secretManagerService *secretmanager.Service

func getSecretManagerService() (*secretmanager.Service, error) {
	if secretManagerService != nil {
		return secretManagerService, nil
	}

	ctx := context.Background()
	svc, err := secretmanager.NewService(ctx)
	if err != nil {
		return nil, err
	}

	secretManagerService = svc

	return svc, nil
}

// SecretCreate creates a secret and populates the lastest version with a payload.
func SecretCreate(project, name, payload string) error {
	svc, err := getSecretManagerService()
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
func SecretDelete(project, name string) error {
	svc, err := getSecretManagerService()
	if err != nil {
		return err
	}

	secret := fmt.Sprintf("projects/%s/secrets/%s", project, name)
	if _, err := svc.Projects.Secrets.Delete(secret).Do(); err != nil {
		return fmt.Errorf("could not delete secret (%s) in project (%s)", name, project)
	}

	return nil
}
