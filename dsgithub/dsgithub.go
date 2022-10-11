package dsgithub

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/GoogleCloudPlatform/deploystack"
	"github.com/GoogleCloudPlatform/deploystack/gcloudtf"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// Meta is a datastructure that combines the Deploystack, github and Terraform
// bits of metadata about a stack.
type Meta struct {
	DeployStack deploystack.Config
	Terraform   gcloudtf.Blocks `json:"terraform" yaml:"terraform"`
	GitRepo     string          `json:"repo" yaml:"repo"`
	GitBranch   string          `json:"branch" yaml:"branch"`
	LocalPath   string          `json:"localpath" yaml:"localpath"`
}

// NewMeta downloads a github repo and parses the DeployStack and Terraform
// information from the stack.
func NewMeta(repo, path string) (Meta, error) {
	d := Meta{}
	d.GitRepo = repo
	d.GitBranch = "main"

	if strings.Contains(repo, "/tree/") {
		end := strings.Index(repo, "/tree/")
		d.GitRepo = repo[:end]
		d.GitBranch = repo[end+6:]
	}

	repoPath := filepath.Base(d.GitRepo)
	repoPath = strings.ReplaceAll(repoPath, "deploystack-", "")
	repoPath = fmt.Sprintf("%s/repo/%s", path, repoPath)

	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		fname := filepath.Join(os.TempDir(), "stdout")
		old := os.Stdout            // keep backup of the real stdout
		temp, _ := os.Create(fname) // create temp file
		defer temp.Close()
		os.Stdout = temp
		_, err = git.PlainClone(
			repoPath,
			false,
			&git.CloneOptions{
				URL:           d.GitRepo,
				ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", d.GitBranch)),
				Progress:      temp,
			})

		if err != nil {
			os.Stdout = old
			out, _ := ioutil.ReadFile(fname)
			fmt.Printf("git response: \n%s\n", string(out))
			return d, fmt.Errorf("cannot get repo: %s", err)
		}

		os.Stdout = old
	}

	orgpwd, err := os.Getwd()
	if err != nil {
		return d, fmt.Errorf("could not get the wd: %s", err)
	}
	if err := os.Chdir(repoPath); err != nil {
		return d, fmt.Errorf("could not change the wd: %s", err)
	}

	s := deploystack.NewStack()

	if err := s.FindAndReadRequired(); err != nil {
		return d, fmt.Errorf("could not read config file: %s", err)
	}

	b, err := gcloudtf.Extract(s.Config.PathTerraform)
	if err != nil {
		log.Fatalf("couldn't extract from TF file: %s", err)
	}

	d.Terraform = *b
	d.DeployStack = s.Config

	if err := os.Chdir(orgpwd); err != nil {
		return d, fmt.Errorf("could not change the wd back: %s", err)
	}

	return d, nil
}

// ShortName retrieves the shortname of whatever we are calling this stack
func (d Meta) ShortName() string {
	r := filepath.Base(d.GitRepo)
	r = strings.ReplaceAll(r, "deploystack-", "")
	return r
}

// ShortNameUnderscore retrieves the shortname of whatever we are calling
// this stack replacing hyphens with underscores
func (d Meta) ShortNameUnderscore() string {
	r := d.ShortName()
	r = strings.ReplaceAll(r, "-", "_")
	return r
}
