package gcloud

import (
	"testing"
	"time"
)

func TestServiceAccountCreate(t *testing.T) {
	c := NewClient(ctx, defaultUserAgent, opts)
	tests := map[string]struct {
		project  string
		username string
		err      error
	}{
		"basic": {projectID, "testSA-" + randSeq(5), nil},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			email, err := c.ServiceAccountCreate(tc.project, tc.username, tc.username)
			if err != tc.err {
				t.Fatalf("create: expected: %+v, got: %+v", tc.err, err)
			}
			// Immediately deleting caused intermittent failures
			time.Sleep(time.Second * 2)

			err = c.ServiceAccountDelete(tc.project, email)
			if err != tc.err {
				t.Logf("delete: trying to delete: %s", email)
				t.Fatalf("delete: expected: no error got: %+v", err)
			}
		})
	}
}
