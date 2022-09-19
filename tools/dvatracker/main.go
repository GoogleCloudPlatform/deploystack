package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/GoogleCloudPlatform/deploystack"
	"github.com/GoogleCloudPlatform/deploystack/gcloudtf"
	"gopkg.in/src-d/go-git.v4"
)

var required = []string{
	"/deploystack.json",
	"/messages/description.txt",
	"/main.tf",
}

var versionStrings = []string{"v1", "alpha", "beta", "v1beta1", "v1beta2", "v1beta3", "v1beta4"}

var repos = []string{
	"https://github.com/GoogleCloudPlatform/deploystack-cost-sentry",
	"https://github.com/GoogleCloudPlatform/deploystack-etl-pipeline",
	"https://github.com/GoogleCloudPlatform/deploystack-load-balanced-vms",
	"https://github.com/GoogleCloudPlatform/deploystack-nosql-client-server",
	"https://github.com/GoogleCloudPlatform/deploystack-ops-agent",
	"https://github.com/GoogleCloudPlatform/deploystack-serverless-e2e-photo-sharing-app",
	"https://github.com/GoogleCloudPlatform/deploystack-single-vm",
	"https://github.com/GoogleCloudPlatform/deploystack-static-hosting-with-domain",
	"https://github.com/GoogleCloudPlatform/deploystack-storage-event-function-app",
	"https://github.com/GoogleCloudPlatform/deploystack-three-tier-app",
}

func main() {
	dvasGlobal := DVAs{}

	for _, repo := range repos {
		DSMeta, err := NewDSMeta(repo)
		if err != nil {
			log.Fatalf("error loading project %s", err)
		}

		dvas := GetDVAs(DSMeta.GetShortName(), DSMeta.Blocks)

		dvasGlobal = append(dvasGlobal, dvas...)
	}
	if err := dvasGlobal.ToCSV("out/ref.csv"); err != nil {
		log.Fatalf("error writing csv: %s", err)
	}
}

func GetDVAs(shortname string, b gcloudtf.Blocks) []DVA {
	result := []DVA{}
	for _, block := range b {
		if block.Kind == "managed" {
			apis, ok := prods[block.Type]
			if !ok {
				fmt.Printf("Needs an entry: %s\n", block.Type)
				continue
			}
			if len(apis) == 0 {
				continue
			}
			for i := range apis {
				for _, version := range versionStrings {
					d := DVA{shortname, strings.Replace(i, "[version]", version, 1)}
					result = append(result, d)
				}
			}

		}
	}

	return result
}

type DVAs []DVA

func (d DVAs) ToCSV(file string) error {
	f, err := os.Create(file)
	defer f.Close()

	if err != nil {
		return fmt.Errorf("could not create file %s: %s", file, err)
	}
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{"appname", "api_method"}); err != nil {
		log.Fatalln("error writing header to file", err)
	}

	for _, record := range d {
		if err := w.Write(record.toSlice()); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}

	return nil
}

type DVA struct {
	Stack string
	API   string
}

func (d DVA) toSlice() []string {
	return []string{d.Stack, d.API}
}

type DSMeta struct {
	DeployStack deploystack.Config
	Blocks      gcloudtf.Blocks
	GitRepo     string
}

func NewDSMeta(repo string) (DSMeta, error) {
	d := DSMeta{}
	d.GitRepo = repo

	repoPath := filepath.Base(repo)
	repoPath = strings.ReplaceAll(repoPath, "deploystack-", "")
	repoPath = fmt.Sprintf("./repo/%s", repoPath)

	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		fname := filepath.Join(os.TempDir(), "stdout")
		old := os.Stdout            // keep backup of the real stdout
		temp, _ := os.Create(fname) // create temp file
		defer temp.Close()
		os.Stdout = temp
		_, err = git.PlainClone(repoPath, false, &git.CloneOptions{
			URL:      repo,
			Progress: temp,
		})
		if err != nil {
			os.Stdout = old
			out, _ := ioutil.ReadFile(fname)
			fmt.Printf("git response: \n%s\n", string(out))
			return d, fmt.Errorf("cannot get repo: %s", err)
		}

		os.Stdout = old
	}

	for _, v := range required {

		v = fmt.Sprintf("%s/%s", repoPath, v)

		if _, err := os.Stat(v); os.IsNotExist(err) {
			return d, fmt.Errorf("required file does not exist: %s", v)
		}
	}

	b, err := gcloudtf.Extract(repoPath)
	if err != nil {
		log.Fatalf("couldn't extract from TF file: %s", err)
	}

	s := deploystack.NewStack()
	if err := s.ReadConfig(fmt.Sprintf("%s/deploystack.json", repoPath), fmt.Sprintf("%s/messages/description.txt", repoPath)); err != nil {
		return d, fmt.Errorf("could not read config file: %s", err)
	}

	d.Blocks = *b
	d.DeployStack = s.Config

	return d, nil
}

func (d DSMeta) GetShortName() string {
	r := filepath.Base(d.GitRepo)
	r = strings.ReplaceAll(r, "deploystack-", "")
	return r
}

func (d DSMeta) GetShortNameUnderScore() string {
	r := d.GetShortName()
	r = strings.ReplaceAll(r, "-", "_")
	return r
}

var prods = map[string]map[string]bool{
	"google_artifact_registry_repository": {"google.devtools.artifactregistry.[version].ArtifactRegistry.CreateRepository": true},
	"google_bigquery_dataset":             {"google.cloud.bigquery.[version].DatasetService.InsertDataset": true},
	"google_bigquery_table": {
		"google.cloud.bigquery.[version].TableService.InsertTable": true,
		"google.cloud.bigquery.[version].TableService.UpdateTable": true,
		"google.cloud.bigquery.[version].TableService.PatchTable":  true,
	},
	"google_cloud_run_service":            {"google.cloud.run.[version].Services.CreateService": true},
	"google_cloud_run_service_iam_policy": {"google.cloud.run.[version].Services.SetIamPolicy": true},
	"google_cloudfunctions_function":      {"google.cloud.functions.[version].CloudFunctionsService.CreateFunction": true},
	"google_composer_environment":         {"google.cloud.orchestration.airflow.service.[version].Environments.CreateEnvironment": true},
	"google_compute_backend_bucket":       {"compute.[version].BackendBucketsService.Insert": true},
	"google_compute_backend_service":      {"compute.[version].BackendServicesService.Insert": true},
	"google_compute_firewall":             {"compute.[version].FirewallsService.Insert": true},
	"google_compute_forwarding_rule":      {"compute.[version].GlobalForwardingRulesService.Insert": true},
	"google_compute_global_address":       {"compute.[version].GlobalAddressesService.Insert": true},
	"google_compute_health_check":         {"compute.[version].HealthChecksService.Insert": true},
	"google_compute_image":                {"compute.[version].ImagesService.Insert": true},
	"google_compute_instance": {
		"compute.[version].InstancesService.Insert":      true,
		"compute.[version].InstancesService.SetMetadata": true,
	},
	"google_compute_instance_group_manager":        {"compute.[version].InstanceGroupManagersService.Insert": true},
	"google_compute_instance_template":             {"compute.[version].InstanceTemplatesService.Insert": true},
	"google_compute_managed_ssl_certificate":       {"compute.[version].SslCertificatesService.Insert": true},
	"google_compute_network":                       {"compute.[version].NetworksService.Insert": true},
	"google_compute_snapshot":                      {"compute.[version].DisksService.CreateSnapshot": true},
	"google_compute_url_map":                       {"compute.[version].UrlMapsService.Insert": true},
	"google_dns_managed_zone":                      {"cloud.dns.api.[version].ChangesService.Create": true},
	"google_dns_record_set":                        {"cloud.dns.api.[version].ManagedZonesService.Create": true},
	"google_pubsub_topic":                          {"google.pubsub.[version].Publisher.CreateTopic": true},
	"google_redis_instance":                        {"google.cloud.redis.[version].CloudRedis.CreateInstance": true},
	"google_secret_manager_secret":                 {"google.cloud.secretmanager.[version].SecretManagerService.CreateSecret": true},
	"google_secret_manager_secret_version":         {"google.cloud.secretmanager.[version].SecretManagerService.AddSecretVersion": true},
	"google_service_networking_connection":         {"google.cloud.servicenetworking.[version].ServicePeeringManager.UpdateConnection": true},
	"google_sql_database_instance":                 {"google.cloud.sql.[version].SqlDatabasesService.Insert": true},
	"google_storage_bucket":                        {"storage.buckets.insert": true},
	"google_storage_bucket_iam_binding":            {"storage.iam.update": true},
	"google_storage_bucket_object":                 {"storage.objects.insert": true, "storage.objects.update": true},
	"google_vpc_access_connector":                  {"google.cloud.vpcaccess.[version].VpcAccessService.CreateConnector": true},
	"google_project_service":                       {},
	"null_resource":                                {},
	"random_password":                              {},
	"random_id":                                    {},
	"random_string":                                {},
	"time_sleep":                                   {},
	"google_service_account":                       {"google.iam.admin.[version].IAM.CreateServiceAccount": true},
	"google_secret_manager_secret_iam_binding":     {"google.cloud.secretmanager.[version].SecretManagerService.SetIamPolicy": true},
	"google_container_registry":                    {},
	"google_sql_user":                              {"google.cloud.sql.[version].SqlUsersService.Insert": true},
	"google_sql_database":                          {"google.cloud.sql.[version].SqlDatabasesService.Insert": true},
	"google_storage_bucket_iam_member":             {"storage.iam.update": true},
	"google_project_iam_member":                    {"google.iam.admin.[version].IAM.UpdateRole": true},
	"google_cloud_run_service_iam_member":          {"google.cloud.run.[version].Services.SetIamPolicy": true},
	"google_service_account_iam_binding":           {"google.iam.admin.[version].IAM.SetIamPolicy": true},
	"google_compute_target_http_proxy":             {"compute.[version].RegionTargetHttpProxiesService.Insert": true, "compute.[version].TargetHttpProxiesService.Insert": true},
	"google_compute_region_network_endpoint_group": {"compute.[version].RegionNetworkEndpointGroupsService.Insert": true},
	"google_compute_target_https_proxy":            {"compute.[version].RegionTargetHttpsProxiesService.Insert": true, "compute.[version].TargetHttpsProxiesService.Insert": true},
}
