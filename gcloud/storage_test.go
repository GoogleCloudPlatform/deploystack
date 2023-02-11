package gcloud

import "testing"

func TestBucketCreate(t *testing.T) {
	c := NewClient(ctx, defaultUserAgent)
	tests := map[string]struct {
		project string
		bucket  string
		err     error
	}{
		"basic": {projectID, projectID + "-testing-bucket", nil},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := c.StorageBucketCreate(tc.project, tc.bucket)
			if err != tc.err {
				t.Fatalf("create: expected: %+v, got: %+v", tc.err, err)
			}

			err = c.StorageBucketDelete(tc.project, tc.bucket)
			if err != tc.err {
				t.Fatalf("delete: expected: no error got: %+v", err)
			}
		})
	}
}

func TestObjectCreate(t *testing.T) {
	c := NewClient(ctx, defaultUserAgent)
	tests := map[string]struct {
		project string
		bucket  string
		path    string
		err     error
	}{
		"basic": {projectID, projectID + "-testing-object", "../README.md", nil},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := c.StorageBucketCreate(tc.project, tc.bucket)
			if err != tc.err {
				t.Fatalf("create bucket: expected: %+v, got: %+v", tc.err, err)
			}

			gspath, err := c.StorageObjectCreate(tc.project, tc.bucket, tc.path)
			if err != tc.err {
				t.Fatalf("create object : expected: %+v, got: %+v", tc.err, err)
			}

			err = c.StorageObjectDelete(tc.project, tc.bucket, gspath)
			if err != tc.err {
				t.Fatalf("delete object : expected: %+v, got: %+v", tc.err, err)
			}

			err = c.StorageBucketDelete(tc.project, tc.bucket)
			if err != tc.err {
				t.Fatalf("delete bucket: expected: no error got: %+v", err)
			}
		})
	}
}
