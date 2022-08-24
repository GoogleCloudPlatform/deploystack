package main

type product struct {
	TerraformType string
	CommandType   string
	TestType      string
	TestCommand   string
	Suffix        string
	Region        bool
	Zone          bool
	LabelField    string
	Expected      string
	Todo          string
}

// terms := list{"zone", "region", "secret_id", ", "name"}

var prods = map[string]product{
	"google_compute_firewall":                {TestType: "gcloud", TestCommand: "gcloud compute firewalls describe", Zone: true, Suffix: `--format="value(name)"`},
	"google_compute_instance":                {TestType: "gcloud", TestCommand: "gcloud compute instances describe", Zone: true, Suffix: `--format="value(name)"`},
	"google_pubsub_topic":                    {TestType: "gcloud", TestCommand: "gcloud pubsub topics describe", Suffix: `--format="value(name)"`},
	"google_cloud_run_service":               {TestType: "gcloud", TestCommand: "gcloud run services describe", Region: true, Suffix: `--format="value(name)"`},
	"google_storage_bucket":                  {TestType: "gsutil", TestCommand: "gsutil ls | grep -c gs://", Expected: "0"},
	"google_storage_bucket_object":           {TestType: "gsutil", TestCommand: "gsutil ls | grep -c gs://", Expected: "0", Todo: "Make sure you check the bucket details at an actual object"},
	"google_cloudfunctions_function":         {TestType: "gcloud", TestCommand: "gcloud functions describe", Region: true, Suffix: `--format="value(name)"`},
	"google_composer_environment":            {TestType: "gcloud", TestCommand: "gcloud composer environments", Suffix: `--format="value(name)"`},
	"google_bigquery_dataset":                {TestType: "bq", TestCommand: "bq ls | grep -c ", LabelField: "table_id", Expected: "1", Todo: "Double check this set of options for test"},
	"google_bigquery_table":                  {TestType: "bq", TestCommand: "bq ls | grep -c ", LabelField: "dataset_id", Todo: "Double check this set of options for test"},
	"google_compute_snapshot":                {TestType: "gcloud", TestCommand: "gcloud compute snapshots describe", Zone: true, Suffix: `--format="value(name)"`},
	"google_compute_image":                   {TestType: "gcloud", TestCommand: "gcloud compute images describe", Zone: true, Suffix: `--format="value(name)"`},
	"google_compute_instance_template":       {TestType: "gcloud", TestCommand: "gcloud compute instance-templates describe", Zone: true, Suffix: `--format="value(name)"`},
	"google_compute_instance_group_manager":  {TestType: "gcloud", TestCommand: "gcloud compute instance groups managed describe", Zone: true, Suffix: `--format="value(name)"`},
	"google_compute_global_address":          {TestType: "gcloud", TestCommand: "gcloud compute addresses describe", Zone: true, Suffix: `--format="value(name)"`},
	"google_compute_health_check":            {TestType: "gcloud", TestCommand: "gcloud compute health-checks describe", Zone: true, Suffix: `--format="value(name)"`},
	"google_compute_backend_service":         {TestType: "gcloud", TestCommand: "gcloud compute backend-service describe", Zone: true, Suffix: `--format="value(name)"`},
	"google_compute_url_map":                 {TestType: "gcloud", TestCommand: "gcloud compute url-maps describe", Zone: true, Suffix: `--format="value(name)"`},
	"google_compute_target_http_proxy":       {TestType: "gcloud", TestCommand: "gcloud compute target-http-proxies describe", Zone: true, Suffix: `--format="value(name)"`},
	"google_compute_forwarding_rule":         {TestType: "gcloud", TestCommand: "gcloud compute forwarding-rules describe", Zone: true, Suffix: `--format="value(name)"`},
	"google_dns_managed_zone":                {TestType: "gcloud", TestCommand: "gcloud gcloud dns record-sets describe", Suffix: `--format="value(name)"`},
	"google_compute_managed_ssl_certificate": {TestType: "gcloud", TestCommand: "gcloud compute ssl-certificate describe", Zone: true, Suffix: `--format="value(name)"`},
	"google_compute_backend_bucket":          {TestType: "gcloud", TestCommand: "gcloud compute backend-buckets  describe", Zone: true, Suffix: `--format="value(name)"`},
	"google_compute_target_https_proxy":      {TestType: "gcloud", TestCommand: "gcloud compute target-https-proxies describe", Zone: true, Suffix: `--format="value(name)"`},
	"google_dns_record_set":                  {TestType: "gcloud", TestCommand: "gcloud gcloud dns record-sets describe", Suffix: `--format="value(name)"`},
	"google_artifact_registry_repository":    {TestType: "gcloud", TestCommand: "gcloud artifacts repositories describe", Suffix: `--format="value(name)"`, LabelField: "repository_id"},
	"google_compute_network":                 {TestType: "gcloud", TestCommand: "gcloud compute networks describe", Zone: true, Suffix: `--format="value(name)"`},
	"google_vpc_access_connector":            {TestType: "gcloud", TestCommand: "gcloud compute networks vpc-access connectors describe", Zone: true, Suffix: `--format="value(name)"`},
	"google_sql_database_instance":           {TestType: "gcloud", TestCommand: "gcloud sql instances describe", Suffix: `--format="value(name)"`},
	"google_redis_instance":                  {TestType: "gcloud", TestCommand: "gcloud redis instances describe", Suffix: `--format="value(name)"`},
	"google_secret_manager_secret":           {TestType: "gcloud", TestCommand: "gcloud secrets describe", Suffix: `--format="value(name)"`, LabelField: "secret_id"},
	"google_service_account":                 {TestType: "gcloud", TestCommand: "gcloud iam service-accounts describe", Suffix: `--format="value(email)"`, LabelField: "account_id", Todo: "This should be an email and not just an account, so add the @domain.com bit"},
}
