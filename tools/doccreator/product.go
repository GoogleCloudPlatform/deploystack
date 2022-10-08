package main

var pinfo = map[string]product{
	"Compute Engine": {
		Title:         "Compute Engine",
		Youtube:       "https://www.youtube.com/watch?v=1XH0gLlGDdk",
		Walkthrough:   "https://console.cloud.google.com/?tutorial=compute_quickstart",
		Documentation: "https://cloud.google.com/compute/docs/quickstart-linux",
		Description:   "Compute Engine is Google Cloud's Virtual technology. With it you can spin up many different configurations of VM to fit the shape of whatever computing needs you have. ",
		Logo:          "https://deploystack.dev/static/img/compute_engine.svg",
	},
	"Cloud Monitoring": {
		Title:         "Cloud Monitoring",
		Logo:          "https://deploystack.dev/static/img/cloud_monitoring.svg",
		Description:   "Cloud Monitoring exposes visibility into the performance, availability, and health of your applications and infrastructure",
		Youtube:       "https://www.youtube.com/watch?v=wY8cmFY4ua8&t=2s",
		Documentation: "https://cloud.google.com/monitoring/docs/",
	},
	"Cloud Logging": {
		Title:         "Cloud Logging",
		Logo:          "https://deploystack.dev/static/img/cloud_logging.svg",
		Description:   "Cloud Logging provides fully managed, real-time log management with storage, search, analysis and alerting at exabyte scale.",
		Youtube:       "https://www.youtube.com/watch?v=gyDp-Cl_MdA",
		Documentation: "https://cloud.google.com/logging/docs/",
		Walkthrough:   "https://console.cloud.google.com/?tutorial=cloud_ops_logging",
	},
	"Cloud Load Balancing": {
		Title:         "Cloud Load Balancing",
		Logo:          "https://deploystack.dev/static/img/cloud_load_balancing.svg",
		Description:   "Google Cloud Load Balancer allows you to place a load balancer in front of the Storage Bucket - allowing you to use SSL certificates, Logging and Monitoring.",
		Youtube:       "https://www.youtube.com/watch?v=gdeOeu4E7eQ",
		Documentation: "https://cloud.google.com/load-balancing/docs",
		Walkthrough:   "https://console.cloud.google.com/?tutorial=load-balancing__ext-load-balancer-backend-buckets",
	},

	"Google Cloud Pub/Sub": {
		Title:         "Google Cloud Pub/Sub",
		Logo:          "https://deploystack.dev/static/img/cloud_pubsub.svg",
		Description:   "Google Cloud Pub/Sub is a messaging bus for integrating applications on different services on different cloud components into an integrated system.",
		Youtube:       "https://www.youtube.com/watch?v=jLI-84UjZLE",
		Documentation: "https://cloud.google.com/pubsub/docs",
		Walkthrough:   "https://console.cloud.google.com/?tutorial=pubsub_quickstart",
	},

	"Billing Budgets": {
		Title:         "Billing Budgets",
		Logo:          "https://deploystack.dev/static/img/cloud_billing.svg",
		Description:   "Billing Budgets allow you to get notified and take action when your billing surpasses thresholds that you set.",
		Youtube:       "https://www.youtube.com/watch?v=F4omjjMZ54k",
		Documentation: "https://cloud.google.com/billing/docs/how-to/budgets",
		Walkthrough:   "https://console.cloud.google.com/?tutorial=billing_create_budgets",
	},

	"Cloud Functions": {
		Title:         "Cloud Functions",
		Logo:          "https://deploystack.dev/static/img/cloud_functions.svg",
		Description:   "Cloud Functions is a functions a service platform that allows you to listen for Cloud Storage file uploads and run code to create thumbnails of them. ",
		Youtube:       "https://www.youtube.com/watch?v=vM-2O-uKBNQ",
		Documentation: "https://cloud.google.com/functions/docs",
		Walkthrough:   "https://console.cloud.google.com/?tutorial=functions_quickstart",
	},

	"Cloud Run": {
		Title:         "Cloud Run",
		Logo:          "https://deploystack.dev/static/img/cloud_run.svg",
		Description:   "Cloud Run allows you to run application in a container, but in a serverless way, no having to configure number of instances, processors, or memory. Upload a container, get a url.",
		Youtube:       "https://www.youtube.com/watch?v=nhwYc4StHIc",
		Documentation: "https://cloud.google.com/run/docs",
		Walkthrough:   "https://console.cloud.google.com/?tutorial=cloud_run_quickstart_index",
	},

	"Cloud SQL": {
		Title:         "Cloud SQL",
		Logo:          "https://deploystack.dev/static/img/cloud_sql.svg",
		Description:   "Cloud SQL provides managed SQL &em; MySQL, SQL Server, or Postgres for the database layer of your applications.",
		Youtube:       "https://www.youtube.com/watch?v=BQlQ-BTMR1U",
		Documentation: "https://cloud.google.com/sql/docs/mysql/quickstart",
		Walkthrough:   "https://console.cloud.google.com/?tutorial=sql_mysql_quickstart",
	},

	"Cloud Memorystore": {
		Title:         "Cloud Memorystore",
		Logo:          "https://deploystack.dev/static/img/memorystore.svg",
		Description:   "Cloud Memorystore, managed Redis, provides the caching layer for your applications.",
		Youtube:       "https://www.youtube.com/watch?v=sVpoAdbh2nU",
		Documentation: "https://cloud.google.com/memorystore/docs/redis/",
		Walkthrough:   "https://console.cloud.google.com/?tutorial=memorystore__memorystore_redis_gcloud_quickstart",
	},

	"Cloud Build": {
		Title:         "Cloud Build",
		Logo:          "https://deploystack.dev/static/img/cloud_build.svg",
		Description:   "Cloud Build is the tool that packages up the containers and deploys them to be available as Cloud Run services.",
		Youtube:       "https://www.youtube.com/watch?v=w7dMHiEyGAs",
		Documentation: "https://cloud.google.com/build/docs",
		Walkthrough:   "https://console.cloud.google.com/?tutorial=cloud_build_quickstart",
	},

	"Secret Manager": {
		Title:         "Secret Manager",
		Logo:          "https://deploystack.dev/static/img/secret_manager.svg",
		Description:   "Cloud Secret Manager stores sensitive particulars about the application for the build process.",
		Youtube:       "https://www.youtube.com/watch?v=JIE89dneaGo",
		Documentation: "https://cloud.google.com/secret-manager/docs",
		Walkthrough:   "https://console.cloud.google.com/?walkthrough_id=secret-manager__create_secret_secretmanager",
	},

	"Cloud DNS": {
		Title:         "Cloud DNS",
		Logo:          "https://deploystack.dev/static/img/cloud_dns.svg",
		Description:   "Cloud DNS provides DNS management - allowing you to run your own DNS infrastructure using Google's existing DNS infrastructure.",
		Youtube:       "https://www.youtube.com/watch?v=OH_Jw8NhEGU",
		Documentation: "https://cloud.google.com/dns/docs",
		Walkthrough:   "https://console.cloud.google.com/?tutorial=static_dns_quickstart",
	},

	"Cloud Domains": {
		Title:         "Cloud Domains",
		Logo:          "https://deploystack.dev/static/img/cloud_domains.svg",
		Description:   "Cloud Domains allows you to buy domains, acting like a Domain Registrar - but once you do, you can manage them and handle billing through your Google Cloud Account.",
		Youtube:       "https://www.youtube.com/watch?v=971wi9Tt5Ds",
		Documentation: "https://cloud.google.com/domains/docs",
		Walkthrough:   "https://console.cloud.google.com/?walkthrough_id=domains__quickstart-register-domain",
	},

	"Cloud Composer": {
		Title:         "Cloud Composer",
		Logo:          "https://deploystack.dev/static/img/cloud_composer.svg",
		Description:   "A fully managed workflow orchestration service built on Apache Airflow.",
		Youtube:       "https://www.youtube.com/watch?v=bwZOAXnCMf8",
		Documentation: "https://cloud.google.com/composer/docs",
		Walkthrough:   "https://console.cloud.google.com/?walkthrough_id=composer--cloud-composer-create-environment",
	},
	"BigQuery": {
		Title:         "BigQuery",
		Logo:          "https://deploystack.dev/static/img/bigquery.svg",
		Description:   "BigQuery is a serverless, cost-effective and multicloud data warehouse designed to help you turn big data into valuable business insights.",
		Youtube:       "https://www.youtube.com/watch?v=CFw4peH2UwU",
		Documentation: "https://cloud.google.com/bigquery/docs",
		Walkthrough:   "https://console.cloud.google.com/?walkthrough_id=bigquery__bigquery-quickstart-query-public-dataset",
	},
	"Cloud Storage": {
		Title:         "Cloud Storage",
		Logo:          "https://deploystack.dev/static/img/cloud_storage.svg",
		Description:   "Cloud Storage provides file storage and public serving of images over http(s).",
		Youtube:       "https://www.youtube.com/watch?v=zvqJ-pcQChE",
		Documentation: "https://cloud.google.com/storage/docs/quickstart-console",
		Walkthrough:   "https://console.cloud.google.com/?walkthrough_id=storage__quickstart-basic-tasks",
	},
}
