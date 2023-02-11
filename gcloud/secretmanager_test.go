package gcloud

import "testing"

func TestSecretCreate(t *testing.T) {
	c := NewClient(ctx, defaultUserAgent)
	tests := map[string]struct {
		project string
		name    string
		payload string
		err     error
	}{
		"basic": {projectID, "testsecret", "secretshhhhhhhhhh", nil},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := c.SecretCreate(tc.project, tc.name, tc.payload)
			if err != tc.err {
				t.Fatalf("expected: %+v, got: %+v", tc.err, err)
			}

			err = c.SecretDelete(tc.project, tc.name)
			if err != tc.err {
				t.Fatalf("expected: no error got: %+v", err)
			}
		})
	}
}
