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

// JobSchedule creates a Cloud Scheduler Job
func JobSchedule(project, region string, job schedulerpb.Job) error {
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

// JobDelete deletes a Cloud Scheduler Job
func JobDelete(project, region, job string) error {
	ctx := context.Background()
	svc, err := getSchedulerService()
	if err != nil {
		return err
	}
	name := fmt.Sprintf("projects/%s/locations/%s/jobs/%s", project, region, job)

	req := schedulerpb.DeleteJobRequest{
		Name: name,
	}

	if err = svc.DeleteJob(ctx, &req); err != nil {
		return err
	}

	return nil
}
