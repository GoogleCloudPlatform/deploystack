package deploystack

import (
	"fmt"
)

// DiskTypeManage promps a user to select a disk type.
func DiskTypeManage(project string) (string, error) {
	fmt.Printf("Enabling service to poll...\n")
	if err := ServiceEnable(project, "compute.googleapis.com"); err != nil {
		return "", fmt.Errorf("error activating service for polling: %s", err)
	}

	fmt.Printf("Choose an operating system\n")
	familyProject := listSelect(DiskProjects, DefaultImageProject)

	fmt.Printf("Polling for %s disk images...\n", familyProject.Value)
	types, err := diskTypes(familyProject.Value)
	if err != nil {
		return "", err
	}

	diskTypes := LabeledValues{}

	for _, v := range types.Items {
		lv := LabeledValue{}
		lv.Label = v.Name
		lv.Value = v.Name
		diskTypes = append(diskTypes, lv)
	}

	famtypes := getListOfImageFamilies(types)

	fmt.Printf("%sChoose a disk family to use for this application. %s\n", TERMCYANB, TERMCLEAR)
	diskfamily := listSelect(famtypes, famtypes[0].Value)

	typesbyfam := getListOfImageTypesByFamily(types, familyProject.Value, diskfamily.Value)

	fmt.Printf("%sChoose a disk type to use for this application. %s\n", TERMCYANB, TERMCLEAR)
	result := listSelect(typesbyfam, typesbyfam[len(typesbyfam)-1].Value)

	return result.Value, nil
}

func MachineTypeManage(project, zone string) (string, error) {
	fmt.Println(Divider)
	fmt.Printf("There are a large number of machine types to choose from. For more infomration, \n")
	fmt.Printf("please refer to the following link for more infomation about Machine types.\n")
	fmt.Printf("%shttps://cloud.google.com/compute/docs/machine-types%s\n", TERMCYANB, TERMCLEAR)
	fmt.Println(Divider)

	fmt.Printf("Enabling service to poll...\n")
	if err := ServiceEnable(project, "compute.googleapis.com"); err != nil {
		return "", fmt.Errorf("error activating service for polling: %s", err)
	}

	fmt.Printf("Polling for machine types...\n")
	types, err := MachineTypes(project, zone)
	if err != nil {
		return "", fmt.Errorf("error polling for machine types : %s", err)
	}

	typefamilies := GetListOfMachineTypeFamily(types)

	fmt.Printf("Choose an Machine Type Family\n")
	familyProject := listSelect(typefamilies, DefaultMachineType)

	filteredtypes := GetListOfMachineTypeByFamily(types, familyProject.Value)

	fmt.Printf("%sChoose a machine type to use for this application. %s\n", TERMCYANB, TERMCLEAR)
	result := listSelect(filteredtypes, filteredtypes[0].Value)

	return result.Value, nil
}

func GCEInstanceManage(project string) (map[string]string, error) {
	// err := errors.New("")
	configs := make(map[string]string)

	// configs["region"], err = RegionManage(project, "compute", DefaultRegion)
	// if err != nil {
	// 	return configs, err
	// }

	// configs["zone"], err = ZoneManage(project, configs["region"])
	// if err != nil {
	// 	return configs, err
	// }

	// configs["machine-type"], err = MachineTypeManage(project, configs["zone"])
	// if err != nil {
	// 	return configs, err
	// }
	// configs["image"], err = DiskTypeManage(project)
	// if err != nil {
	// 	return configs, err
	// }

	items := Customs{
		{Name: "disksize", Description: "Enter the size of the boot disk you want in GB", Default: "200", Validation: "integer"},
		{Name: "disktype", Description: "Enter the type of the boot disk you want", Default: "pd-standard", Options: []string{"pd-standard", "pd-balanced", "pd-ssd"}},
		{Name: "webserver", Description: "Do you want this to be a webserver? ", Default: "no", Validation: "yesorno"},
	}

	if err := items.Collect(); err != nil {
		return configs, err
	}

	for _, v := range items {
		configs[v.Name] = v.Value
	}

	return configs, nil
}
