package deploystack

import (
	"context"
	"fmt"
	"sort"

	"google.golang.org/api/cloudfunctions/v1"
)

var cloudfunctionsService *cloudfunctions.Service

func getCloudFunctionsService(project string) (*cloudfunctions.Service, error) {
	if cloudfunctionsService != nil {
		return cloudfunctionsService, nil
	}

	if err := ServiceEnable(project, "cloudfunctions.googleapis.com"); err != nil {
		return nil, fmt.Errorf("error activating service for polling: %s", err)
	}

	ctx := context.Background()
	svc, err := cloudfunctions.NewService(ctx, opts)
	if err != nil {
		return nil, err
	}

	cloudfunctionsService = svc

	return svc, nil
}

// regionsFunctions will return a list of regions for Cloud Functions
func regionsFunctions(project string) ([]string, error) {
	resp := []string{}

	if err := ServiceEnable(project, "cloudfunctions.googleapis.com"); err != nil {
		return resp, fmt.Errorf("error activating service for polling: %s", err)
	}

	svc, err := getCloudFunctionsService(project)
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
