package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/GoogleCloudPlatform/deploystack/github"
)

func main() {

}

func getRepo(repo, path string) error {
	candidate := fmt.Sprintf("%s/%s", path, repo)
	dir := getAcceptableDir(candidate)
	return cloneFromRepo(repo, dir)

}

func getAcceptableDir(candidate string) string {
	if _, err := os.Stat(candidate); os.IsNotExist(err) {
		return candidate
	}
	i := 1
	for {
		attempted := fmt.Sprintf("%s_%d", candidate, i)
		if _, err := os.Stat(attempted); os.IsNotExist(err) {
			return attempted
		}
		i++
	}
}

func cloneFromRepo(repo, path string) error {
	host := "github.com"
	owner := "GoogleCloudPlatform"
	fullrepo := fmt.Sprintf("https://%s/%s/%s", host, owner, repo)

	gh := github.NewRepo(fullrepo)
	err := gh.Clone(path)
	if err != nil {
		// This allows using a shortened name of the repo as the label here.
		if !strings.Contains(repo, "deploystack-") {
			fullrepo = fmt.Sprintf("https://%s/%s/deploystack-%s", host, owner, repo)
			gh = github.NewRepo(fullrepo)
			err = gh.Clone(path)
		}
	}

	suffix := strings.ReplaceAll(repo, "deploystack-", "")
	place := fmt.Sprintf("%s/repo/%s", path, suffix)

	// Github code puts the code in a weird place. This moves it
	err = filepath.Walk(place, func(root string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}

		if info.Name() == suffix {
			return nil
		}

		old := filepath.Join(place, info.Name())
		new := filepath.Join(path, info.Name())
		return os.Rename(old, new)
	})

	return err

}
