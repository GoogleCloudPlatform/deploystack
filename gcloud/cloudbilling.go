package gcloud

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/cloudresourcemanager/v1"
)

var cloudbillingService *cloudbilling.APIService

func (c *Client) getCloudbillingService() (*cloudbilling.APIService, error) {
	var err error
	svc := c.services.cloudbillingService

	if svc != nil {
		return svc, nil
	}

	svc, err = cloudbilling.NewService(context.Background(), c.opts)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve service: %w", err)
	}

	svc.UserAgent = c.userAgent
	c.services.cloudbillingService = svc

	return svc, nil
}

// BillingAccountList gets a list of the billing accounts a user has access to
func (c *Client) BillingAccountList() ([]*cloudbilling.BillingAccount, error) {
	resp := []*cloudbilling.BillingAccount{}
	svc, err := c.getCloudbillingService()
	if err != nil {
		return resp, err
	}

	results, err := svc.BillingAccounts.List().Do()
	if err != nil {
		return resp, err
	}

	return results.BillingAccounts, nil
}

// BillingAccountAttach will enable billing in a given project
func (c *Client) BillingAccountAttach(project, account string) error {
	retries := 10
	svc, err := c.getCloudbillingService()
	if err != nil {
		return err
	}

	ba := fmt.Sprintf("billingAccounts/%s", account)
	proj := fmt.Sprintf("projects/%s", project)

	cfg := cloudbilling.ProjectBillingInfo{
		BillingAccountName: ba,
	}

	var looperr error
	for i := 0; i < retries; i++ {
		_, looperr = svc.Projects.UpdateBillingInfo(proj, &cfg).Do()
		if looperr == nil {
			return nil
		}
		if strings.Contains(looperr.Error(), "User is not authorized to get billing info") {
			continue
		}
	}

	if strings.Contains(looperr.Error(), "Request contains an invalid argument") {
		return ErrorBillingInvalidAccount
	}

	if strings.Contains(looperr.Error(), "Not a valid billing account") {
		return ErrorBillingInvalidAccount
	}

	if strings.Contains(looperr.Error(), "The caller does not have permission") {
		return ErrorBillingNoPermission
	}

	return looperr
}

// ProjectListWithBilling gets a list of projects with their billing information
func (c *Client) ProjectListWithBilling(p []*cloudresourcemanager.Project) ([]ProjectWithBilling, error) {
	res := []ProjectWithBilling{}

	svc, err := c.getCloudbillingService()
	if err != nil {
		return res, err
	}

	projs, _ := c.ProjectListWithBillingEnabled()
	// if err != nil {
	// 	return res, err
	// }

	var wg sync.WaitGroup
	wg.Add(len(p))

	for _, v := range p {
		go func(p *cloudresourcemanager.Project) {
			defer wg.Done()

			if _, ok := projs[p.ProjectId]; ok {
				pwb := ProjectWithBilling{Name: p.Name, ID: p.ProjectId, BillingEnabled: true}
				res = append(res, pwb)
				return
			}

			// Getting random quota errors when somebody had too many projects.
			// sleeping randoming for a second fixed it.
			// I don't think these requests can be fixed by batching.
			sleepRandom()
			if p.LifecycleState == "ACTIVE" && p.Name != "" {
				proj := fmt.Sprintf("projects/%s", p.ProjectId)
				tmp, err := svc.Projects.GetBillingInfo(proj).Do()
				if err != nil {
					if strings.Contains(err.Error(), "The caller does not have permission") {
						// fmt.Printf("project: %+v\n", p)
						return
					}

					fmt.Printf("error getting billing information: %s\n", err)
					return
				}

				pwb := ProjectWithBilling{Name: p.Name, ID: p.ProjectId, BillingEnabled: tmp.BillingEnabled}
				res = append(res, pwb)
				return
			}
		}(v)
	}
	wg.Wait()

	return res, nil
}

// ProjectListWithBillingEnabled queries the billing accounts a user has access to
// to generate a list of projects for each billing account. Will hopefully
// reduce the number of calls made to billing api
func (c *Client) ProjectListWithBillingEnabled() (map[string]bool, error) {
	r := map[string]bool{}
	svc, err := c.getCloudbillingService()
	if err != nil {
		return r, err
	}

	bas, err := c.BillingAccountList()
	if err != nil {
		return r, err
	}

	for _, v := range bas {
		result, err := svc.BillingAccounts.Projects.List(v.Name).Do()
		if err != nil {
			return r, err
		}
		for _, v := range result.ProjectBillingInfo {
			if v.BillingEnabled {
				r[v.ProjectId] = true
			}
		}
	}

	return r, nil
}

func randomInRange(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func sleepRandom() {
	d := time.Second * time.Duration(randomInRange(0, 1))
	time.Sleep(d)
}
