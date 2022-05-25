package deploystack

import "fmt"

// DiskTypeManage promps a user to select a disk type.
func DiskTypeManage(project string) (string, error) {
	fmt.Printf("Enabling service to poll...\n")
	if err := ServiceEnable(project, "compute.googleapis.com"); err != nil {
		return "", fmt.Errorf("error activating service for polling: %s", err)
	}

	fmt.Printf("Choose an operating system\n")
	familyProject := listSelect(DiskProjects, DiskProjects[0].value)

	fmt.Printf("Polling for %s disk images...\n", familyProject.value)
	types, err := diskTypes(familyProject.value)
	if err != nil {
		return "", err
	}

	diskTypes := labeledValues{}

	for _, v := range types.Items {
		lv := labeledValue{}
		lv.label = v.Name
		lv.value = v.Name
		diskTypes = append(diskTypes, lv)
	}

	famtypes := getListOfDiskFamilies(types)

	fmt.Printf("%sChoose a disk family to use for this application. %s\n", TERMCYANB, TERMCLEAR)
	diskfamily := listSelect(famtypes, famtypes[0].value)

	typesbyfam := getListOfDiskTypes(types, diskfamily.value)

	fmt.Printf("%sChoose a disk type to use for this application. %s\n", TERMCYANB, TERMCLEAR)
	result := listSelect(typesbyfam, typesbyfam[len(typesbyfam)-1].value)

	return result.value, nil
}
