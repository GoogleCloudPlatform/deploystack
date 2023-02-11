package gcloud

import (
	"fmt"
	"testing"

	"cloud.google.com/go/scheduler/apiv1beta1/schedulerpb"
)

func TestScheduleJob(t *testing.T) {
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
