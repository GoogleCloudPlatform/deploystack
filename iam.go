package deploystack

import (
	"context"
	"fmt"

	"google.golang.org/api/iam/v1"
)

var iamService *iam.Service

func getIAMService() (*iam.Service, error) {
	if iamService != nil {
		return iamService, nil
	}

	ctx := context.Background()
	svc, err := iam.NewService(ctx, opts)
	if err != nil {
		return nil, err
	}

	iamService = svc

	return svc, nil
}

func CreateServiceAccount(project, username, displayName string) (string, error) {
	svc, err := getIAMService()
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
