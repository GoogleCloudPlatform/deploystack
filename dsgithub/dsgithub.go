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
	Github      Github          `json:"github" yaml:"github"`
	LocalPath   string          `json:"localpath" yaml:"localpath"`
}

// Github contains the details of a github repo for the purpose of downloading
type Github struct {
	Repo   string `json:"repo" yaml:"repo"`
	Branch string `json:"branch" yaml:"branch"`
}

// NewGithub generates Github from a url that might contain branch information
func NewGithub(repo string) Github {
	result := Github{}
	result.Repo = repo
	result.Branch = "main"

	if strings.Contains(repo, "/tree/") {
		end := strings.Index(repo, "/tree/")
		result.Repo = repo[:end]
		result.Branch = repo[end+6:]
	}

	return result
}

// RepoPath Returns the path that the github content will be cached at
func (g Github) RepoPath(path string) string {
	result := filepath.Base(g.Repo)
	result = strings.ReplaceAll(result, "deploystack-", "")
	result = fmt.Sprintf("%s/repo/%s", path, result)
	return result
}

// Clone performs a git clone to the directory of our choosing
func (g Github) Clone(path string) error {
	localPath := g.RepoPath(path)

	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		fname := filepath.Join(os.TempDir(), "stdout")
		old := os.Stdout            // keep backup of the real stdout
		temp, _ := os.Create(fname) // create temp file
		defer temp.Close()
		os.Stdout = temp
		_, err = git.PlainClone(
			localPath,
			false,
			&git.CloneOptions{
				URL:           g.Repo,
				ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", g.Branch)),
				Progress:      temp,
			})

		if err != nil {
			os.Stdout = old
			out, _ := ioutil.ReadFile(fname)
			fmt.Printf("git response: \n%s\n", string(out))
			return fmt.Errorf("cannot get repo: %s", err)
		}

		os.Stdout = old
	}

	return nil
}

// NewMeta downloads a github repo and parses the DeployStack and Terraform
// information from the stack.
func NewMeta(repo, path string) (Meta, error) {
	g := NewGithub(repo)

	if err := g.Clone(path); err != nil {
		return Meta{}, fmt.Errorf("cannot clone repo: %s", err)
	}

	d, err := NewMetaFromLocal(g.RepoPath(path))
	if err != nil {
		return Meta{}, fmt.Errorf("cannot parse deploystack into: %s", err)
	}
	d.Github = g
	d.LocalPath = g.RepoPath(path)

	return d, nil
}

// NewMetaFromLocal allows project to point at local directories for info
// as well as pulling down from github
func NewMetaFromLocal(path string) (Meta, error) {
	d := Meta{}
	orgpwd, err := os.Getwd()
	if err != nil {
		return d, fmt.Errorf("could not get the wd: %s", err)
	}
	if err := os.Chdir(path); err != nil {
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
	r := filepath.Base(d.Github.Repo)
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
