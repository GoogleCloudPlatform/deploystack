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
	"testing"

	"cloud.google.com/go/scheduler/apiv1beta1/schedulerpb"
	"github.com/stretchr/testify/assert"
)

func TestScheduleJob(t *testing.T) {
	t.Parallel()
	c := NewClient(ctx, defaultUserAgent)
	tests := map[string]struct {
		project string
		region  string
		name    string
		job     schedulerpb.Job
		err     error
	}{
		"basic": {
			projectID,
			"us-central1",
			"testtriggername",
			schedulerpb.Job{
				Name:     fmt.Sprintf("projects/%s/locations/%s/jobs/%s", projectID, "us-central1", "testtriggername"),
				Schedule: "0 6 * * *",
				Target: &schedulerpb.Job_HttpTarget{
					HttpTarget: &schedulerpb.HttpTarget{
						Uri:        "http://example.com",
						HttpMethod: schedulerpb.HttpMethod_GET,
						Headers:    map[string]string{"Content-Type": "application/octet-stream,User-Agent=Google-Cloud-Scheduler"},
					},
				},
			},
			nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := c.JobSchedule(tc.project, tc.region, tc.job)
			if err != tc.err {
				t.Fatalf("expected: %+v, got: %+v", tc.err, err)
			}

			err = c.JobDelete(tc.project, tc.region, tc.name)
			if err != tc.err {
				t.Fatalf("expected: no error got: %+v", err)
			}
		})
	}
}

func TestSchedulerBadProject(t *testing.T) {
	t.Parallel()
	bad := "notavalidprojectnameanditshouldfaildasdas"
	tests := map[string]struct {
		servicefunc func() error
		err         error
	}{
		"JobSchedule": {
			servicefunc: func() error {
				c := NewClient(context.Background(), "testing")
				return c.JobSchedule(bad, "", schedulerpb.Job{})
			},
			err: fmt.Errorf("error activating service for polling"),
		},
		"JobDelete": {
			servicefunc: func() error {
				c := NewClient(context.Background(), "testing")
				return c.JobDelete(bad, "", "")
			},
			err: fmt.Errorf("error activating service for polling"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc := tc
			t.Parallel()
			err := tc.servicefunc()
			assert.ErrorContains(t, err, tc.err.Error())
		})
	}
}
