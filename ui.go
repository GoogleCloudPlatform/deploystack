package deploystack

import (
	"bufio"
	"fmt"
	"os"
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
	family := listSelect(families, DefaultImageFamily)

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

// BillingAccountManage either grabs the users only BillingAccount or
// presents a list of BillingAccounts to select from.
func BillingAccountManage() (string, error) {
	accounts, err := billingAccounts()
	if err != nil {
		return "", fmt.Errorf("could not get list of billing accounts: %s", err)
	}

	labeled := []string{}

	for _, v := range accounts {
		labeled = append(labeled, fmt.Sprintf("%s (%s)", v.DisplayName, strings.ReplaceAll(v.Name, "billingAccounts/", "")))
	}

	if len(accounts) == 1 {
		fmt.Printf("\nOnly found one billing account. Using : %s%s%s.\n", TERMCYAN, accounts[0].DisplayName, TERMCLEAR)
		return extractAccount(labeled[0]), nil
	}

	fmt.Printf("\n%sPlease select one of your billing accounts to use with this project%s.\n", TERMCYAN, TERMCLEAR)
	result := listSelect(toLabeledValueSlice(labeled), labeled[0])

	return extractAccount(result.Value), nil
}

func extractAccount(s string) string {
	sl := strings.Split(s, "(")
	return strings.ReplaceAll(sl[1], ")", "")
}

// ProjectManage promps a user to select a project.
func ProjectManage() (string, string, error) {
	createString := "CREATE NEW PROJECT"
	project, err := ProjectID()
	if err != nil {
		return "", "", err
	}

	projects, err := projects()
	if err != nil {
		return "", "", err
	}

	lvs := LabeledValues{}

	for _, v := range projects {
		lv := LabeledValue{Label: v.Name, Value: v.ID}

		if !v.BillingEnabled {
			lv.Label = fmt.Sprintf("%s%s (Billing Disabled)%s", TERMGREY, v.Name, TERMCLEAR)
		}

		lvs = append(lvs, lv)
	}

	lvs = append([]LabeledValue{{createString, createString}}, lvs...)

	fmt.Printf("\n%sChoose a project to use for this application.%s\n\n", TERMCYANB, TERMCLEAR)
	fmt.Printf("%sNOTE:%s This app will make changes to the project. %s\n", TERMCYANREV, TERMCYAN, TERMCLEAR)
	fmt.Printf("While those changes are reverseable, it would be better to put it in a fresh new project. \n")

	lv := listSelect(lvs, project)
	project = lv.Value

	if project == createString {
		project, err = projectPrompt()
		if err != nil {
			return "", "", err
		}
		lv = LabeledValue{project, project}
	}

	if err := ProjectIDSet(project); err != nil {
		return lv.Value, lv.Label, fmt.Errorf("error: unable to set project (%s) in environment: %s", project, err)
	}

	return lv.Value, lv.Label, nil
}

// projectPrompt manages the interaction of creating a project, including prompts.
func projectPrompt() (string, error) {
	result := ""
	sec1 := NewSection("Creating the project")

	sec1.Open()
	fmt.Printf("Project IDs are immutable and can be set only during project \n")
	fmt.Printf("creation. They must start with a lowercase letter and can have \n")
	fmt.Printf("lowercase ASCII letters, digits or hyphens.  \n")
	fmt.Printf("Project IDs must be between 6 and 30 characters. \n")
	fmt.Printf("%sPlease enter a new project name to create: %s\n", TERMCYANB, TERMCLEAR)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		if len(text) == 0 {
			fmt.Printf("%sPlease enter a new project name to create: %s\n", TERMCYANB, TERMCLEAR)
			continue
		}

		if err := projectCreate(text); err != nil {
			fmt.Printf("%sProject name could not be created %s\n", TERMREDREV, TERMCLEAR)
			fmt.Printf("%sReason: %s %s\n", TERMREDB, err, TERMCLEAR)
			fmt.Printf("%sPlease choose another. %s\n", TERMREDREV, TERMCLEAR)
			continue
		}

		fmt.Printf("Project Created\n")
		result = text
		break

	}
	sec1.Close()

	sec2 := NewSection("Activating Billing for the project")
	sec2.Open()
	account, err := BillingAccountManage()
	if err != nil {
		return "", fmt.Errorf("could not determine proper billing account: %s ", err)
	}

	if err := BillingAccountProjectAttach(result, account); err != nil {
		return "", fmt.Errorf("could not link billing account: %s ", err)
	}
	sec2.Close()
	return result, nil
}

// regions will return a list of regions depending on product type
func regions(project, product string) ([]string, error) {
	switch product {
	case "compute":
		return regionsCompute(project)
	case "functions":
		return regionsFunctions(project)
	case "run":
		return regionsRun(project)
	}

	return []string{}, fmt.Errorf("invalid product requested: %s", product)
}

// RegionManage promps a user to select a region.
func RegionManage(project, product, def string) (string, error) {
	fmt.Printf("Polling for regions...\n")
	regions, err := regions(project, product)
	if err != nil {
		return "", err
	}
	fmt.Printf("%sChoose a valid region to use for this application. %s\n", TERMCYANB, TERMCLEAR)
	region := listSelect(toLabeledValueSlice(regions), def)

	return region.Value, nil
}

// ZoneManage promps a user to select a zone.
func ZoneManage(project, region string) (string, error) {
	fmt.Printf("Polling for zones...\n")
	zones, err := zones(project, region)
	if err != nil {
		return "", err
	}

	fmt.Printf("%sChoose a valid zone to use for this application. %s\n", TERMCYANB, TERMCLEAR)
	zone := listSelect(toLabeledValueSlice(zones), zones[0])
	return zone.Value, nil
}

// Start presents a little documentation screen which also prevents the user
// from timing out the request to activate Cloud Shell
func Start() {
	fmt.Printf(Divider)
	colorPrintln("Deploystack", TERMCYANB)
	fmt.Printf("Deploystack will walk you through setting some options for the  \n")
	fmt.Printf("stack this solutions installs. \n")
	fmt.Printf("Most questions have a default that you can choose by hitting the Enter key  \n")
	fmt.Printf(Divider)
	colorPrintln("Press the Enter Key to continue", TERMCYANB)
	var input string
	fmt.Scanln(&input)
}
