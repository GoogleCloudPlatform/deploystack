package deploystack

import (
	"fmt"
	"testing"

	"google.golang.org/api/cloudbuild/v1"
)

func TestTriggerCreate(t *testing.T) {
	tests := map[string]struct {
		project string
		trigger cloudbuild.BuildTrigger
		err     error
	}{
		"basic": {
			projectID,
			cloudbuild.BuildTrigger{
				Name:      "testtrigger",
				EventType: "MANUAL",
				GitFileSource: &cloudbuild.GitFileSource{
					RepoType: "GITHUB",
					Uri:      "http://github.com/googlecloudplatform/deploystack-single-vm",
					Path:     "test.yaml",
					Revision: fmt.Sprintf("refs/heads/%s", "main"),
				},
				SourceToBuild: &cloudbuild.GitRepoSource{
					Uri:      "http://github.com/googlecloudplatform/deploystack-single-vm",
					Ref:      fmt.Sprintf("refs/heads/%s", "main"),
					RepoType: "GITHUB",
				},
			},
			nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := CloudBuildTriggerCreate(tc.project, tc.trigger)
			if err != tc.err {
				t.Fatalf("expected: %+v, got: %+v", tc.err, err)
			}

			err = CloudBuildTriggerDelete(tc.project, got.Id)
			if err != tc.err {
				t.Fatalf("expected: no error got: %+v", err)
			}
		})
	}
}
