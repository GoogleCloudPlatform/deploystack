package deploystack

import (
	"fmt"
	"sort"
	"strings"
)

// ImageManage promps a user to select a disk type.
func ImageManage(project string) (string, error) {
	fmt.Println(Divider)
	fmt.Printf("There are a large number of machine images to choose from. For more infomration, \n")
	fmt.Printf("please refer to the following link for more infomation about machine images.\n")
	fmt.Printf("%shttps://cloud.google.com/compute/docs/images%s\n", TERMCYANB, TERMCLEAR)
	fmt.Println(Divider)

	colorPrintln("Choose an operating system.", TERMCYANB)
	ImageTypeProject := listSelect(DiskProjects, DefaultImageProject)

	fmt.Printf("Polling for %s images...\n", ImageTypeProject.Value)
	images, err := images(project, ImageTypeProject.Value)
	if err != nil {
		return "", err
	}

	families := getListOfImageFamilies(images)

	colorPrintln("Choose a disk family to use for this application.", TERMCYANB)
	family := listSelect(families, families[0].Value)

	imagesByFam := getListOfImageTypesByFamily(images, ImageTypeProject.Value, family.Value)

	colorPrintln("Choose a disk type to use for this application.", TERMCYANB)
	result := listSelect(imagesByFam, imagesByFam[len(imagesByFam)-1].Value)

	return result.Value, nil
}

func colorPrintln(msg, color string) {
	fmt.Printf("%s%s %s\n", color, msg, TERMCLEAR)
}

func MachineTypeManage(project, zone string) (string, error) {
	fmt.Println(Divider)
	fmt.Printf("There are a large number of machine types to choose from. For more infomration, \n")
	fmt.Printf("please refer to the following link for more infomation about Machine types.\n")
	fmt.Printf("%shttps://cloud.google.com/compute/docs/machine-types%s\n", TERMCYANB, TERMCLEAR)
	fmt.Println(Divider)

	fmt.Printf("Polling for machine types...\n")
	types, err := machineTypes(project, zone)
	if err != nil {
		return "", fmt.Errorf("error polling for machine types : %s", err)
	}

	typefamilies := getListOfMachineTypeFamily(types)

	fmt.Printf("Choose an Machine Type Family\n")
	familyProject := listSelect(typefamilies, DefaultMachineType)

	filteredtypes := getListOfMachineTypeByFamily(types, familyProject.Value)

	fmt.Printf("%sChoose a machine type to use for this application. %s\n", TERMCYANB, TERMCLEAR)
	result := listSelect(filteredtypes, filteredtypes[0].Value)

	return result.Value, nil
}

type GCEInstanceConfig map[string]string

func (gce GCEInstanceConfig) Print(title string) {
	keys := []string{}
	for i := range gce {
		keys = append(keys, i)
	}

	longest := longestLength(toLabeledValueSlice(keys))

	colorPrintln(title, TERMCYANREV)
	exclude := []string{}

	if s, ok := gce["instance-name"]; ok && len(s) > 0 {
		printSetting("instance-name", s, longest)
		exclude = append(exclude, "instance-name")
	}

	if s, ok := gce["region"]; ok && len(s) > 0 {
		printSetting("region", s, longest)
		exclude = append(exclude, "region")
	}

	if s, ok := gce["zone"]; ok && len(s) > 0 {
		printSetting("zone", s, longest)
		exclude = append(exclude, "zone")
	}

	ordered := []string{}
	for i, v := range gce {
		if strings.Contains(strings.Join(exclude, " "), i) {
			continue
		}
		if len(v) < 1 {
			continue
		}

		ordered = append(ordered, i)
	}
	sort.Strings(ordered)

	for i := range ordered {
		key := ordered[i]
		printSetting(key, gce[key], longest)
	}
}

func GCEInstanceManage(project, basename string) (GCEInstanceConfig, error) {
	var err error
	configs := make(map[string]string)

	defaultImage, err := getLatestImage(project, DefaultImageProject, DefaultImageFamily)
	if err != nil {
		return configs, err
	}

	defaultConfig := GCEInstanceConfig{
		"instance-image":        defaultImage,
		"instance-disksize":     "200",
		"instance-disktype":     "pd-standard",
		"instance-tags":         "[http-server,https-server]",
		"instance-name":         fmt.Sprintf("%s-instance", basename),
		"region":                DefaultRegion,
		"zone":                  fmt.Sprintf("%s-a", DefaultRegion),
		"instance-machine-type": "n1-standard-1",
	}

	ClearScreen()
	fmt.Println(Divider)
	colorPrintln("Configure a Compute Engine Instance", TERMCYANB)
	fmt.Printf("Let's walk through configuring a Compute Engine Instance (Virtual Machine). \n")
	fmt.Printf("you can either accept a default configuration with settings that work for \n")
	fmt.Printf("trying out most use cases, or hand configure key settings. \n")
	fmt.Println(Divider)

	defaultConfig.Print("Default Configuration")

	chooseDefault := Custom{
		Name:        "choosedefault",
		Description: "Do you want to use the default? ('No' means custom)",
		Default:     "yes",
		Validation:  "yesorno",
	}

	if err := chooseDefault.Collect(); err != nil {
		return configs, err
	}

	if chooseDefault.Value == "yes" {
		return defaultConfig, nil
	}

	nameItem := Custom{
		Name:        "name",
		Description: "Enter the name of the instance",
		Default:     fmt.Sprintf("%s-instance", basename),
	}

	if err := nameItem.Collect(); err != nil {
		return configs, err
	}

	configs["instance-name"] = nameItem.Value

	configs["region"], err = RegionManage(project, "compute", DefaultRegion)
	if err != nil {
		return configs, err
	}

	configs["zone"], err = ZoneManage(project, configs["region"])
	if err != nil {
		return configs, err
	}

	configs["instance-machine-type"], err = MachineTypeManage(project, configs["zone"])
	if err != nil {
		return configs, err
	}
	configs["instance-image"], err = ImageManage(project)
	if err != nil {
		return configs, err
	}

	items := Customs{
		{Name: "instance-disksize", Description: "Enter the size of the boot disk you want in GB", Default: "200", Validation: "integer"},
		{Name: "instance-disktype", Description: "Enter the type of the boot disk you want", Default: "pd-standard", Options: []string{"pd-standard", "pd-balanced", "pd-ssd"}},
		{Name: "webserver", Description: "Do you want this to be a webserver (Expose ports 80 & 443)? ", Default: "no", Validation: "yesorno"},
	}

	if err := items.Collect(); err != nil {
		return configs, err
	}

	for _, v := range items {

		if v.Name == "webserver" {
			configs["instance-tags"] = "[]"
			if v.Value == "yes" {
				configs["instance-tags"] = "[http-server,https-server]"
			}
			continue
		}

		configs[v.Name] = v.Value

	}

	return configs, nil
}
