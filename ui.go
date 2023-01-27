package deploystack

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
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
	ImageTypeProject := DiskProjects.SelectUI()

	fmt.Printf("Polling for %s images...\n", ImageTypeProject.Value)
	images, err := ComputeImageList(project, ImageTypeProject.Value)
	if err != nil {
		return "", err
	}

	families := ComputeImageFamilyList(images)

	colorPrintln("Choose a disk family to use for this application.", TERMCYANB)
	family := families.SelectUI()

	imagesByFam := ComputeImageTypeListByFamily(images, ImageTypeProject.Value, family.Value)

	colorPrintln("Choose a disk type to use for this application.", TERMCYANB)
	result := imagesByFam.SelectUI()

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
	types, err := ComputeMachineTypeList(project, zone)
	if err != nil {
		return "", fmt.Errorf("error polling for machine types : %s", err)
	}

	typefamilies := ComputeMachineTypeFamilyList(types)

	fmt.Printf("Choose an Machine Type Family\n")
	familyProject := typefamilies.SelectUI()

	filteredtypes := ComputeMachineTypeListByFamily(types, familyProject.Value)

	fmt.Printf("%sChoose a machine type to use for this application. %s\n", TERMCYANB, TERMCLEAR)
	result := filteredtypes.SelectUI()

	return result.Value, nil
}

type GCEInstanceConfig map[string]string

func (gce GCEInstanceConfig) Print(title string) {
	keys := []string{}
	for i := range gce {
		keys = append(keys, i)
	}

	list := NewLabeledValues(keys, "")
	longest := list.LongestLen()

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

	defaultImage, err := ComputeImageLatestGet(project, DefaultImageProject, DefaultImageFamily)
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
	accounts, err := BillingAccountList()
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
	list := NewLabeledValues(labeled, labeled[0])
	result := list.SelectUI()

	return extractAccount(result.Value), nil
}

func extractAccount(s string) string {
	sl := strings.Split(s, "(")
	return strings.ReplaceAll(sl[1], ")", "")
}

// projectPrompt manages the interaction of creating a project, including prompts.
func projectPrompt(currentProject string) (string, error) {
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

		if currentProject == "" {
			tmp, err := ProjectList()
			if err != nil || len(tmp) == 0 || tmp[0].ID == "" {
				return "", fmt.Errorf("could not determine an alternate project for parent detection: %s ", err)
			}

			currentProject = tmp[0].ID
		}

		parent, err := ProjectParentGet(currentProject)
		if err != nil {
			return "", fmt.Errorf("could not determine proper parent for project: %s ", err)
		}

		if err := ProjectCreate(text, parent.Id, parent.Type); err != nil {
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

	if err := BillingAccountAttach(result, account); err != nil {
		return "", fmt.Errorf("could not link billing account: %s ", err)
	}
	sec2.Close()
	return result, nil
}

// RegionsList will return a list of RegionsList depending on product type
func RegionsList(project, product string) ([]string, error) {
	switch product {
	case "compute":
		return ComputeRegionList(project)
	case "functions":
		return FunctionRegionList(project)
	case "run":
		return RunRegionsList(project)
	}

	return []string{}, fmt.Errorf("invalid product (%s) requested", product)
}

// RegionManage promps a user to select a region.
func RegionManage(project, product, defaultValue string) (string, error) {
	fmt.Printf("Polling for regions...\n")
	regions, err := RegionsList(project, product)
	if err != nil {
		return "", err
	}

	fmt.Printf("%sChoose a valid region to use for this application. %s\n", TERMCYANB, TERMCLEAR)
	list := NewLabeledValues(regions, defaultValue)
	region := list.SelectUI()

	return region.Value, nil
}

// ZoneManage promps a user to select a zone.
func ZoneManage(project, region string) (string, error) {
	fmt.Printf("Polling for zones...\n")
	zones, err := ComputeZoneList(project, region)
	if err != nil {
		return "", err
	}

	fmt.Printf("%sChoose a valid zone to use for this application. %s\n", TERMCYANB, TERMCLEAR)

	list := NewLabeledValues(zones, zones[0])
	zone := list.SelectUI()

	return zone.Value, nil
}

// Start presents a little documentation screen which also prevents the user
// from timing out the request to activate Cloud Shell
func Start() {
	fmt.Printf(Divider)
	colorPrintln("Deploystack", TERMCYANB)
	fmt.Printf("Deploystack will walk you through setting some options for the\n")
	fmt.Printf("stack this solutions installs.\n")
	fmt.Printf("Most questions have a default that you can choose by hitting the Enter key\n")
	fmt.Printf(Divider)
	colorPrintln("Press the Enter Key to continue", TERMCYANB)
	var input string
	fmt.Scanln(&input)
}

// LabeledValue is a struct that contains a label/value pair
type LabeledValue struct {
	Value     string
	Label     string
	IsDefault bool
}

// NewLabeledValue takes a string and converts it to a LabeledValue. If a |
// delimiter is present it will split into a different label/value
func NewLabeledValue(s string) LabeledValue {
	l := LabeledValue{s, s, false}

	if strings.Contains(s, "|") {
		sl := strings.Split(s, "|")
		l = LabeledValue{sl[0], sl[1], false}
	}

	return l
}

// RenderUI creates a string that will evantually be shown to a user with
// terminal formatting characters and what not.
// extracted render function to make unit testing easier
func (l LabeledValue) RenderUI(index, width int) string {
	stripped := cleanTerminalChars(l.Label)
	offset := len(l.Label) - len(stripped)
	width += offset

	if l.IsDefault {
		return fmt.Sprintf("%s%2d) %-*s %s", TERMCYANB, index, width, l.Label, TERMCLEAR)
	}
	return fmt.Sprintf("%2d) %-*s ", index, width, l.Label)
}

// LabeledValues is collection of LabledValue structs
type LabeledValues []LabeledValue

// SelectUI handles showing a user the list of values an allowing them to select
// one from the list
func (l LabeledValues) SelectUI() LabeledValue {
	itemCount := len(l)
	answer := l.GetDefault()
	defaultExists := answer != LabeledValue{}

	ui := l.RenderListUI()
	fmt.Print(ui)

	if defaultExists {
		fmt.Printf("Choose number from list, or just [enter] for %s%s%s\n", TERMCYANB, answer.Label, TERMCLEAR)
	} else {
		fmt.Printf("Choose number from list.\n")
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		if len(text) == 0 {
			break
		}

		opt, err := strconv.Atoi(text)
		if err != nil || opt > itemCount {
			fmt.Printf("Please enter a numeric between 1 and %d\n", itemCount)
			fmt.Printf("You entered %s\n", text)
			continue
		}

		answer = l[opt-1]
		break

	}

	return answer
}

// Sort orders the LabeledValues by Label
func (l *LabeledValues) Sort() {
	sort.Slice(*l, func(i, j int) bool {
		iStr := strings.ToLower(cleanTerminalChars((*l)[i].Label))
		jStr := strings.ToLower(cleanTerminalChars((*l)[j].Label))
		return iStr < jStr
	})
}

// LongestLen returns the length of longest LABEL in the list
func (l *LabeledValues) LongestLen() int {
	longest := 0

	for _, v := range *l {
		if len(cleanTerminalChars(v.Label)) > longest {
			longest = len(cleanTerminalChars(v.Label))
		}
	}

	return longest
}

// GetDefault returns the deafult value of the LabeledValues list
func (l *LabeledValues) GetDefault() LabeledValue {
	for _, v := range *l {
		if v.IsDefault {
			return v
		}
	}
	return LabeledValue{}
}

// SetDefault sets the default value of the list
func (l *LabeledValues) SetDefault(value string) {
	for i, v := range *l {
		if v.Value == value {
			v.IsDefault = true
			(*l)[i] = v
		}
	}
}

// RenderListUI creates a string that will evantually be shown to a user with
// terminal formatting characters and what not in a multi column list if there
// are enough entries
// extracted render function to make unit testing easier
func (l LabeledValues) RenderListUI() string {
	sb := strings.Builder{}
	width := l.LongestLen()
	itemCount := len(l)

	if itemCount < 11 {
		for i, v := range l {
			sb.WriteString(v.RenderUI(i+1, width))
			sb.WriteString("\n")
		}
	} else {
		halfcount := int(math.Ceil(float64(itemCount / 2)))

		if float64(halfcount) < float64(itemCount)/2 {
			halfcount++
		}

		for i := 0; i < halfcount; i++ {
			sb.WriteString(l[i].RenderUI(i+1, width))

			idx := i + halfcount + 1

			if idx > itemCount {
				sb.WriteString("\n")
				break
			}

			sb.WriteString(l[idx-1].RenderUI(idx, width))

			sb.WriteString("\n")
		}
	}
	return sb.String()
}

// NewLabeledValues takes a slice of strings and returns a list of LabeledValues
func NewLabeledValues(sl []string, defaultValue string) LabeledValues {
	r := LabeledValues{}

	for _, v := range sl {
		val := NewLabeledValue(v)
		if val.Value == defaultValue {
			val.IsDefault = true
		}

		r = append(r, val)
	}
	return r
}

func cleanTerminalChars(s string) string {
	replacements := []string{
		TERMCYAN, "",
		TERMCYANB, "",
		TERMCYANREV, "",
		TERMRED, "",
		TERMREDB, "",
		TERMREDREV, "",
		TERMCLEAR, "",
		TERMCLEARSCREEN, "",
		TERMGREY, "",
	}

	replacer := strings.NewReplacer(replacements...)

	r := replacer.Replace(s)

	return r
}
