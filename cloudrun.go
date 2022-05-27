package deploystack

import (
	"context"
	"fmt"
	"sort"

	"google.golang.org/api/run/v1"
)

var runService *run.APIService

func getRunService(project string) (*run.APIService, error) {
	if runService != nil {
		return runService, nil
	}

	if err := ServiceEnable(project, "run.googleapis.com"); err != nil {
		return nil, fmt.Errorf("error activating service for polling: %s", err)
	}

	ctx := context.Background()
	svc, err := run.NewService(ctx, opts)
	if err != nil {
		return nil, err
	}

	runService = svc

	return svc, nil
}

// regionsRun will return a list of regions for Cloud Run
func regionsRun(project string) ([]string, error) {
	resp := []string{}

	svc, err := getRunService(project)
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
