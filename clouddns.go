package deploystack

import (
	"context"
	"fmt"

	"google.golang.org/api/dns/v1"
)

func ZoneCreate(project, name, domain string) error {
	ctx := context.Background()
	svc, err := dns.NewService(ctx)
	if err != nil {
		return err
	}

	zone := *&dns.ManagedZone{
		Name:        name,
		DnsName:     fmt.Sprintf("%s.", domain),
		Description: fmt.Sprintf("A DNS Zone for managing %s", domain),
	}

	if _, err := svc.ManagedZones.Create(project, &zone).Do(); err != nil {
		return err
	}

	return nil
}

func ZoneDelete(project, name string) error {
	ctx := context.Background()
	svc, err := dns.NewService(ctx)
	if err != nil {
		return err
	}

	if err := svc.ManagedZones.Delete(project, name).Do(); err != nil {
		return err
	}

	return nil
}
