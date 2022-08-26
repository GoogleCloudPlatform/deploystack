package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestFindClosingBracket(t *testing.T) {
	tests := map[string]struct {
		start   int
		content string
		want    int
	}{
		"1": {start: 1, content: "", want: 0},
		"2": {start: 4, content: `

		# Enabling services in your GCP project
		variable "gcp_service_list" {
		  description = "The list of apis necessary for the project"
		  type        = list(string)
		  default = [
			"compute.googleapis.com",
		  ]
		}`, want: 9},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := findClosingBracket(tc.start, strings.Split(tc.content, "\n"))
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}
