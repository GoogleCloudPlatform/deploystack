package deploystack

import (
	"context"
	"fmt"

	scheduler "cloud.google.com/go/scheduler/apiv1beta1"
	"cloud.google.com/go/scheduler/apiv1beta1/schedulerpb"
)

var schedulerService *scheduler.CloudSchedulerClient

func getSchedulerService() (*scheduler.CloudSchedulerClient, error) {
	if schedulerService != nil {
		return schedulerService, nil
	}

	ctx := context.Background()
	svc, err := scheduler.NewCloudSchedulerClient(ctx, opts)
	if err != nil {
		return nil, err
	}

	schedulerService = svc

	return svc, nil
}

func ScheduleJob(project, region string, job schedulerpb.Job) error {
	ctx := context.Background()
	svc, err := getSchedulerService()
	if err != nil {
		return err
	}
	parent := fmt.Sprintf("projects/%s/locations/%s", project, region)

	req := schedulerpb.CreateJobRequest{
		Parent: parent,
		Job:    &job,
	}

	if _, err = svc.CreateJob(ctx, &req); err != nil {
		return err
	}

	return nil
}
