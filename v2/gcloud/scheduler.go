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
	"context"
	"fmt"

	scheduler "cloud.google.com/go/scheduler/apiv1beta1"
	"cloud.google.com/go/scheduler/apiv1beta1/schedulerpb"
)

func (c *Client) getSchedulerService(project string) (*scheduler.CloudSchedulerClient, error) {
	var err error
	svc := c.services.scheduler

	if svc != nil {
		return svc, nil
	}

	if err := c.ServiceEnable(project, "cloudscheduler.googleapis.com"); err != nil {
		return nil, fmt.Errorf("error activating service for polling: %s", err)
	}

	svc, err = scheduler.NewCloudSchedulerClient(c.ctx, c.opts)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve service: %w", err)
	}

	c.services.scheduler = svc

	return svc, nil
}

// JobSchedule creates a Cloud Scheduler Job
func (c *Client) JobSchedule(project, region string, job schedulerpb.Job) error {
	ctx := context.Background()
	svc, err := c.getSchedulerService(project)
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
func (c *Client) JobDelete(project, region, job string) error {
	ctx := context.Background()
	svc, err := c.getSchedulerService(project)
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
