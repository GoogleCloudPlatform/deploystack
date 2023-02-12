package gcloud

import (
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	"google.golang.org/api/cloudresourcemanager/v1"
)

func (c *Client) getCloudResourceManagerService() (*cloudresourcemanager.Service, error) {
	var err error
	svc := c.services.resourceManager

	if svc != nil {
		return svc, nil
	}

	svc, err = cloudresourcemanager.NewService(c.ctx, c.opts)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve service: %w", err)
	}

	svc.UserAgent = c.userAgent
	c.services.resourceManager = svc

	return svc, nil
}

// ProjectNumberGet will get the project_number for the input projectid
func (c *Client) ProjectNumberGet(id string) (string, error) {
	resp := ""
	svc, err := c.getCloudResourceManagerService()
	if err != nil {
		return resp, err
	}

	results, err := svc.Projects.Get(id).Do()
	if err != nil {
		return resp, err
	}

	resp = strconv.Itoa(int(results.ProjectNumber))

	return resp, nil
}

// ProjectParentGet returns the parent of an input project
func (c *Client) ProjectParentGet(id string) (*cloudresourcemanager.ResourceId, error) {
	svc, err := c.getCloudResourceManagerService()
	if err != nil {
		return nil, err
	}

	results, err := svc.Projects.Get(id).Do()
	if err != nil {
		return nil, err
	}

	return results.Parent, nil
}

// ProjectList gets a list of the ProjectList a user has access to
func (c *Client) ProjectList() ([]ProjectWithBilling, error) {
	resp := []ProjectWithBilling{}

	svc, err := c.getCloudResourceManagerService()
	if err != nil {
		return resp, err
	}

	results, err := svc.Projects.List().Filter("lifecycleState=ACTIVE").Do()
	if err != nil {
		return resp, err
	}

	pwb, err := c.ProjectListWithBilling(results.Projects)
	if err != nil {
		return resp, err
	}

	sort.Slice(pwb, func(i, j int) bool {
		return strings.ToLower(pwb[i].Name) < strings.ToLower(pwb[j].Name)
	})

	return pwb, nil
}

// ProjectWithBilling is a project with it's billing status
type ProjectWithBilling struct {
	Name           string
	ID             string
	BillingEnabled bool
}

// ProjectCreate does the work of actually creating a new project in your
// GCP account
func (c *Client) ProjectCreate(project, parent, parentType string) error {
	svc, err := c.getCloudResourceManagerService()
	if err != nil {
		return err
	}

	par := &cloudresourcemanager.ResourceId{}
	if parent != "" && parentType != "" {
		par.Id = parent
		par.Type = parentType
	}

	proj := cloudresourcemanager.Project{
		Name:      project,
		ProjectId: project,
		Parent:    par,
	}

	result, err := svc.Projects.Create(&proj).Do()
	if err != nil {
		if strings.Contains(err.Error(), "project_id must be at most 30 characters long") {
			return ErrorProjectCreateTooLong
		}
		if strings.Contains(err.Error(), "project_id contains invalid characters") {
			return ErrorProjectInvalidCharacters
		}
		if strings.Contains(err.Error(), "requested entity already exists") {
			return ErrorProjectAlreadyExists
		}

		return err
	}

	for i := 0; i < 20; i++ {
		op, err := svc.Operations.Get(result.Name).Do()
		if err != nil {
			return fmt.Errorf("could not poll for project completion: %s", err)
		}
		if op.Done {
			if op.Error != nil {
				return fmt.Errorf("project creation was unsuccessful, reason: %s ", op.Error.Message)
			}
			return nil
		}
		time.Sleep(2 * time.Second)
	}

	return ErrorProjectCreateTooLong
}

// ProjectDelete does the work of actually deleting an existing project in
// your GCP account
func (c *Client) ProjectDelete(project string) error {
	svc, err := c.getCloudResourceManagerService()
	if err != nil {
		return err
	}

	_, err = svc.Projects.Delete(project).Do()
	if err != nil {
		return err
	}

	return nil
}

// ProjectGrantIAMRole grants a given principal a given role in a given project
func (c *Client) ProjectGrantIAMRole(project, role, principal string) error {
	svc, err := c.getCloudResourceManagerService()
	if err != nil {
		return err
	}
	getReq := cloudresourcemanager.GetIamPolicyRequest{}

	policy, err := svc.Projects.GetIamPolicy(project, &getReq).Do()
	if err != nil {
		return fmt.Errorf("cannot get iam policy for project (%s): %s", project, err)
	}

	b := cloudresourcemanager.Binding{}
	b.Role = role
	b.Members = append(b.Members, principal)
	policy.Bindings = append(policy.Bindings, &b)

	setReq := cloudresourcemanager.SetIamPolicyRequest{}
	setReq.Policy = policy

	if _, err = svc.Projects.SetIamPolicy(project, &setReq).Do(); err != nil {
		return fmt.Errorf("cannot set iam policy role (%s) for project (%s): %s", role, project, err)
	}

	return nil
}

// ProjectIDGet gets the currently set default project
func (c Client) ProjectIDGet() (string, error) {
	cmd := exec.Command("gcloud", "config", "get-value", "project")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("cannot get project id: %s ", err)
	}

	return strings.TrimSpace(string(out)), nil
}

// ProjectIDSet sets the currently set default project
func (c *Client) ProjectIDSet(project string) error {
	cmd := exec.Command("gcloud", "config", "set", "project", project)
	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("cannot set project id: %s ", err)
	}

	return nil
}

// ProjectExists confirms that a project actually exists
func (c *Client) ProjectExists(project string) bool {
	svc, err := c.getCloudResourceManagerService()
	if err != nil {
		return false
	}

	_, err = svc.Projects.Get(project).Do()
	if err != nil {
		return false
	}

	return true
}
