package deploystack

import (
	"context"
	"fmt"
	"sort"

	"google.golang.org/api/compute/v1"
)

// DiskProjects are the list of projects for disk images for Compute Engine
var DiskProjects = labeledValues{
	labeledValue{label: "CentOS", value: "centos-cloud"},
	labeledValue{label: "Container-Optimized OS (COS)", value: "cos-cloud"},
	labeledValue{label: "Debian", value: "debian-cloud"},
	labeledValue{label: "Fedora CoreOS", value: "fedora-coreos-cloud"},
	labeledValue{label: "Red Hat Enterprise Linux (RHEL)", value: "rhel-cloud"},
	labeledValue{label: "Red Hat Enterprise Linux (RHEL) for SAP", value: "rhel-sap-cloud"},
	labeledValue{label: "Rocky Linux", value: "rocky-linux-cloud"},
	labeledValue{label: "SQL Server", value: "windows-sql-cloud"},
	labeledValue{label: "SUSE Linux Enterprise Server (SLES)", value: "suse-cloud"},
	labeledValue{label: "SUSE Linux Enterprise Server (SLES) for SAP", value: "suse-cloud"},
	labeledValue{label: "SUSE Linux Enterprise Server (SLES) BYOS", value: "suse-byos-cloud"},
	labeledValue{label: "Ubuntu LTS", value: "ubuntu-os-cloud"},
	labeledValue{label: "Ubuntu Pro", value: "ubuntu-os-pro-cloud"},
	labeledValue{label: "Windows Server", value: "windows-cloud"},
}

// regionsCompute will return a list of regions for Compute Engine
func regionsCompute(project string) ([]string, error) {
	resp := []string{}

	ctx := context.Background()
	svc, err := compute.NewService(ctx, opts)
	if err != nil {
		return resp, err
	}

	results, err := svc.Regions.List(project).Do()
	if err != nil {
		return resp, err
	}

	for _, v := range results.Items {
		resp = append(resp, v.Name)
	}

	sort.Strings(resp)

	return resp, nil
}

// zones will return a list of zones in a given region
func zones(project, region string) ([]string, error) {
	resp := []string{}

	ctx := context.Background()
	svc, err := compute.NewService(ctx, opts)
	if err != nil {
		return resp, err
	}

	filter := fmt.Sprintf("name=%s*", region)

	results, err := svc.Zones.List(project).Filter(filter).Do()
	if err != nil {
		return resp, err
	}

	for _, v := range results.Items {
		resp = append(resp, v.Name)
	}

	sort.Strings(resp)

	return resp, nil
}

func machineTypes(project, zone string) (labeledValues, error) {
	resp := labeledValues{}
	ctx := context.Background()
	svc, err := compute.NewService(ctx, opts)
	if err != nil {
		return resp, err
	}

	results, err := svc.MachineTypes.List(project, zone).Do()
	if err != nil {
		return resp, err
	}

	for _, v := range results.Items {
		lb := labeledValue{}
		lb.value = v.Name
		mb := formatMBToGB(v.MemoryMb)
		lb.label = fmt.Sprintf("%s CPUs: %d Mem: %s", v.Name, v.GuestCpus, mb)

		resp = append(resp, lb)
	}

	return resp, nil
}

func formatMBToGB(i int64) string {
	return fmt.Sprintf("%d GB", i/1024)
}

func diskTypes(project string) (*compute.ImageList, error) {
	resp := &compute.ImageList{}
	ctx := context.Background()
	svc, err := compute.NewService(ctx, opts)
	if err != nil {
		return resp, err
	}

	results, err := svc.Images.List(project).Do()
	if err != nil {
		return resp, err
	}

	return results, nil
}

func getListOfDiskFamilies(imgs *compute.ImageList) labeledValues {
	fam := make(map[string]bool)
	lb := labeledValues{}

	for _, v := range imgs.Items {
		fam[v.Family] = false
	}

	for i := range fam {
		if i == "" {
			continue
		}
		lb = append(lb, labeledValue{i, i})
	}
	lb.sort()
	return lb
}

func getListOfDiskTypes(imgs *compute.ImageList, family string) labeledValues {
	lb := labeledValues{}

	for _, v := range imgs.Items {
		if v.Family == family {
			lb = append(lb, labeledValue{v.Name, v.Name})
		}
	}

	last := lb[len(lb)-1]
	last.label = fmt.Sprintf("%s (Latest)", last.value)
	lb[len(lb)-1] = last
	lb.sort()

	return lb
}
