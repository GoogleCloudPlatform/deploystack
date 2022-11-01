package deploystack

import (
	"context"
	"fmt"

	"google.golang.org/api/cloudbuild/v1"
)

var cloudBuildService *cloudbuild.Service

func getCloudBuildService() (*cloudbuild.Service, error) {
	if cloudBuildService != nil {
		return cloudBuildService, nil
	}

	ctx := context.Background()
	svc, err := cloudbuild.NewService(ctx, opts)
	if err != nil {
		return nil, err
	}

	cloudBuildService = svc

	return svc, nil
}

// CreateCloudBuildTrigger creates a build trigger in a given project
func CreateCloudBuildTrigger(project string, trigger cloudbuild.BuildTrigger) (*cloudbuild.BuildTrigger, error) {
	svc, err := getCloudBuildService()
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

// DeleteCloudBuildTrigger deletes a build trigger in a given project
func DeleteCloudBuildTrigger(project string, triggerid string) error {
	svc, err := getCloudBuildService()
	if err != nil {
		return err
	}

	if _, err := svc.Projects.Triggers.Delete(project, triggerid).Do(); err != nil {
		return fmt.Errorf("cannot delete trigger: %s", err)
	}

	return nil
}
