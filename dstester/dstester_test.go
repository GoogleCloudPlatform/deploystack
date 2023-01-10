package dstester

import (
	"reflect"
	"strings"
	"testing"
)

func TestResourcesInit(t *testing.T) {
	tests := map[string]struct {
		input Resources
		want  string
	}{
		"a": {input: Resources{Project: "test", Items: []Resource{{}}}, want: "test"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.input.Init()
			got := tc.input.Items[0].Project
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestResourceExistsString(t *testing.T) {
	gcloud := which("gcloud")
	tests := map[string]struct {
		input Resource
		want  string
	}{
		"basic": {
			input: Resource{
				Product: "compute instances",
				Name:    "test",
			},
			want: "gcloud compute instances describe test --format=\"value(name)\"",
		},
		"complicated": {
			input: Resource{
				Product:   "compute instances",
				Name:      "test",
				Arguments: map[string]string{"region": "us-central1"},
			},
			want: "gcloud compute instances describe test --region us-central1 --format=\"value(name)\"",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input.existsString()
			got = strings.ReplaceAll(got, gcloud, "gcloud")
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}
