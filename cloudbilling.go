package deploystack

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

func getCloudbillingService() (*cloudbilling.APIService, error) {
	if cloudbillingService != nil {
		return cloudbillingService, nil
	}

	ctx := context.Background()
	svc, err := cloudbilling.NewService(ctx, opts)
	if err != nil {
		return nil, err
	}

	cloudbillingService = svc

	return svc, nil
}

// billingAccounts gets a list of the billing accounts a user has access to
func billingAccounts() ([]*cloudbilling.BillingAccount, error) {
	resp := []*cloudbilling.BillingAccount{}
	svc, err := getCloudbillingService()
	if err != nil {
		return resp, err
	}

	results, err := svc.BillingAccounts.List().Do()
	if err != nil {
		return resp, err
	}

	return results.BillingAccounts, nil
}

// BillingAccountProjectAttach will enable billing in a given project
func BillingAccountProjectAttach(project, account string) error {
	retries := 10
	svc, err := getCloudbillingService()
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

func getBillingForProjects(p []*cloudresourcemanager.Project) ([]projectWithBilling, error) {
	res := []projectWithBilling{}

	svc, err := getCloudbillingService()
	if err != nil {
		return res, err
	}
	var wg sync.WaitGroup
	wg.Add(len(p))

	for _, v := range p {
		go func(p *cloudresourcemanager.Project) {
			defer wg.Done()
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

				pwb := projectWithBilling{Name: p.Name, ID: p.ProjectId, BillingEnabled: tmp.BillingEnabled}
				res = append(res, pwb)
			}
		}(v)
	}
	wg.Wait()

	return res, nil
}

func randomInRange(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func sleepRandom() {
	d := time.Second * time.Duration(randomInRange(0, 1))
	time.Sleep(d)
}
