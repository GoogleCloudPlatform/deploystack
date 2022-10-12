package dsgithub

import (
	"os"
	"reflect"
	"testing"
)

func TestShortName(t *testing.T) {
	tests := map[string]struct {
		in   string
		want string
	}{
		"deploystack-repo":     {in: "https://github.com/GoogleCloudPlatform/deploystack-cost-sentry", want: "cost-sentry"},
		"non-deploystack-repo": {in: "https://github.com/tpryan/microservices-demo", want: "microservices-demo"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			m := Meta{}
			m.Github.Repo = tc.in

			got := m.ShortName()
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestShortNameUnderscore(t *testing.T) {
	tests := map[string]struct {
		in   string
		want string
	}{
		"deploystack-repo":     {in: "https://github.com/GoogleCloudPlatform/deploystack-cost-sentry", want: "cost_sentry"},
		"non-deploystack-repo": {in: "https://github.com/tpryan/microservices-demo", want: "microservices_demo"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			m := Meta{}
			m.Github.Repo = tc.in

			got := m.ShortNameUnderscore()
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestNewGithub(t *testing.T) {
	tests := map[string]struct {
		in   string
		want Github
	}{
		"defaultbranch": {in: "https://github.com/GoogleCloudPlatform/deploystack-cost-sentry", want: Github{Repo: "https://github.com/GoogleCloudPlatform/deploystack-cost-sentry", Branch: "main"}},
		"otherbranch":   {in: "https://github.com/tpryan/microservices-demo/tree/deploystack-enable", want: Github{Repo: "https://github.com/tpryan/microservices-demo", Branch: "deploystack-enable"}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := NewGithub(tc.in)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestRepoPath(t *testing.T) {
	tests := map[string]struct {
		in   Github
		path string
		want string
	}{
		"defaultbranch": {
			in: Github{
				Repo:   "https://github.com/GoogleCloudPlatform/deploystack-cost-sentry",
				Branch: "main",
			},
			path: ".",
			want: "./repo/cost-sentry",
		},
		"otherbranch": {
			in: Github{
				Repo:   "https://github.com/tpryan/microservices-demo",
				Branch: "deploystack-enable",
			},
			path: ".",
			want: "./repo/microservices-demo",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.in.RepoPath(tc.path)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestNewMeta(t *testing.T) {
	tests := map[string]struct {
		repo string
		path string
		want Meta
	}{
		"defaultbranch": {
			repo: "https://github.com/GoogleCloudPlatform/deploystack-cost-sentry",
			path: ".",
			want: Meta{Github: Github{Repo: "https://github.com/GoogleCloudPlatform/deploystack-cost-sentry", Branch: "main"}, LocalPath: "./repo/cost-sentry"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			os.RemoveAll(tc.want.LocalPath)
			got, err := NewMeta(tc.repo, tc.path)
			if err != nil {
				t.Fatalf("expected: no error, got: %v", err)
			}
			if !reflect.DeepEqual(tc.want.Github.Repo, got.Github.Repo) {
				t.Fatalf("expected: %v, got: %v", tc.want.Github.Repo, got.Github.Repo)
			}

			if !reflect.DeepEqual(tc.want.Github.Branch, got.Github.Branch) {
				t.Fatalf("expected: %v, got: %v", tc.want.Github.Branch, got.Github.Branch)
			}

			if !reflect.DeepEqual(tc.want.LocalPath, got.LocalPath) {
				t.Fatalf("expected: %v, got: %v", tc.want.LocalPath, got.LocalPath)
			}

			os.RemoveAll(got.LocalPath)
		})
	}
}
