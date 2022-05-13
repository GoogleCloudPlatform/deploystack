package deploystack

import (
	"strings"
	"testing"
)

func TestZoneCreate(t *testing.T) {
	tests := map[string]struct {
		project string
		name    string
		domain  string
		errMsg  string
	}{
		"ShouldWork":  {projectID, "testing-zone", "yesornositetester.com", ""},
		"BadDomain":   {projectID, "testing-zone2", "example.com", "may be reserved or registered already"},
		"BadZoneName": {projectID, "TestingDNSZone_examplecom", "examplecom", "Invalid value for 'entity.managedZone.name'"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := ZoneCreate(tc.project, tc.name, tc.domain)
			if !errorContains(err, tc.errMsg) {
				t.Logf("Project: %s", projectID)
				t.Fatalf("expected: error(%v) got: error(%v)", tc.errMsg, err)
			}
		})
	}

	for _, tc := range tests {
		if tc.errMsg == "" {
			ZoneDelete(tc.project, tc.name)
		}
	}
}

func errorContains(out error, want string) bool {
	if out == nil {
		return want == ""
	}
	if want == "" {
		return false
	}
	return strings.Contains(out.Error(), want)
}
