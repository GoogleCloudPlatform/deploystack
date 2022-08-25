# Test Creator

Will stub out a Deploystack test based on targeted terraform directory.

usage

```shell
cd tools/testcreator
./clean && go run *.go -folder [root folder of a deplpoystack project]

```

It will output a test file to ./out. Move that to your deploystack project folder. 

For example here is the output from deploystack-load-balanced-vms looks like:

``` bash
# Boilerplate ommitted

# TODO: Script is pretty rudimentary, so make sure all of these variables are actually needed for test
PROJECT_ID=[SET VALUE HERE]
PROJECT_NUMBER=[SET VALUE HERE]
REGION=[SET VALUE HERE]
ZONE=[SET VALUE HERE]
BASENAME=[SET VALUE HERE]
NODES=[SET VALUE HERE]
 
gcloud config set project ${PROJECT}

terraform init
# TODO: Make sure all of these variables are correct and actually need to be passed in.:
terraform apply  \
   -var project_id=$PROJECT_ID \
   -var project_number=$PROJECT_NUMBER \
   -var region=$REGION \
   -var zone=$ZONE \
   -var basename=$BASENAME \
   -var nodes=$NODES \
   -auto-approve


section_open "Test google_compute_instance exists"
    evalTest 'gcloud compute instances describe $BASENAME-exemplar --zone $ZONE --format="value(name)"' "$BASENAME-exemplar"
section_close

section_open "Test google_compute_snapshot exists"
    evalTest 'gcloud compute snapshots describe $BASENAME-snapshot --zone $ZONE --format="value(name)"' "$BASENAME-snapshot"
section_close

section_open "Test google_compute_image exists"
    evalTest 'gcloud compute images describe $BASENAME-latest --zone $ZONE --format="value(name)"' "$BASENAME-latest"
section_close

section_open "Test google_compute_instance_template exists"
    evalTest 'gcloud compute instance-templates describe $BASENAME-template --zone $ZONE --format="value(name)"' "$BASENAME-template"
section_close

section_open "Test google_compute_instance_group_manager exists"
    evalTest 'gcloud compute instance groups managed describe $BASENAME-mig --zone $ZONE --format="value(name)"' "$BASENAME-mig"
section_close

section_open "Test google_compute_global_address exists"
    evalTest 'gcloud compute addresses describe $BASENAME-ip --zone $ZONE --format="value(name)"' "$BASENAME-ip"
section_close

section_open "Test google_compute_health_check exists"
    evalTest 'gcloud compute health-checks describe $BASENAME-health-chk --zone $ZONE --format="value(name)"' "$BASENAME-health-chk"
section_close

section_open "Test google_compute_firewall exists"
    evalTest 'gcloud compute firewalls describe allow-health-check --zone $ZONE --format="value(name)"' "allow-health-check"
section_close

section_open "Test google_compute_backend_service exists"
    evalTest 'gcloud compute backend-service describe $BASENAME-service --zone $ZONE --format="value(name)"' "$BASENAME-service"
section_close

section_open "Test google_compute_url_map exists"
    evalTest 'gcloud compute url-maps describe $BASENAME-lb --zone $ZONE --format="value(name)"' "$BASENAME-lb"
section_close

section_open "Test google_compute_target_http_proxy exists"
    evalTest 'gcloud compute target-http-proxies describe $BASENAME-lb-proxy --zone $ZONE --format="value(name)"' "$BASENAME-lb-proxy"
section_close

section_open "Test google_compute_forwarding_rule exists"
    evalTest 'gcloud compute forwarding-rules describe $BASENAME-http-lb-forwarding-rule --zone $ZONE --format="value(name)"' "$BASENAME-http-lb-forwarding-rule"
section_close


# TODO: Make sure all of these variables are correct and actually need to be passed in.:
terraform destroy  \
   -var project_id=$PROJECT_ID \
   -var project_number=$PROJECT_NUMBER \
   -var region=$REGION \
   -var zone=$ZONE \
   -var basename=$BASENAME \
   -var nodes=$NODES \
   -var gcp_service_list=$GCP_SERVICE_LIST \
   -auto-approve


section_open "Test google_compute_instance does not exists"
    evalTest 'gcloud compute instances describe $BASENAME-exemplar --zone $ZONE --format="value(name)"' "EXPECTERROR"
section_close

section_open "Test google_compute_snapshot does not exists"
    evalTest 'gcloud compute snapshots describe $BASENAME-snapshot --zone $ZONE --format="value(name)"' "EXPECTERROR"
section_close

section_open "Test google_compute_image does not exists"
    evalTest 'gcloud compute images describe $BASENAME-latest --zone $ZONE --format="value(name)"' "EXPECTERROR"
section_close

section_open "Test google_compute_instance_template does not exists"
    evalTest 'gcloud compute instance-templates describe $BASENAME-template --zone $ZONE --format="value(name)"' "EXPECTERROR"
section_close

section_open "Test google_compute_instance_group_manager does not exists"
    evalTest 'gcloud compute instance groups managed describe $BASENAME-mig --zone $ZONE --format="value(name)"' "EXPECTERROR"
section_close

section_open "Test google_compute_global_address does not exists"
    evalTest 'gcloud compute addresses describe $BASENAME-ip --zone $ZONE --format="value(name)"' "EXPECTERROR"
section_close

section_open "Test google_compute_health_check does not exists"
    evalTest 'gcloud compute health-checks describe $BASENAME-health-chk --zone $ZONE --format="value(name)"' "EXPECTERROR"
section_close

section_open "Test google_compute_firewall does not exists"
    evalTest 'gcloud compute firewalls describe allow-health-check --zone $ZONE --format="value(name)"' "EXPECTERROR"
section_close

section_open "Test google_compute_backend_service does not exists"
    evalTest 'gcloud compute backend-service describe $BASENAME-service --zone $ZONE --format="value(name)"' "EXPECTERROR"
section_close

section_open "Test google_compute_url_map does not exists"
    evalTest 'gcloud compute url-maps describe $BASENAME-lb --zone $ZONE --format="value(name)"' "EXPECTERROR"
section_close

section_open "Test google_compute_target_http_proxy does not exists"
    evalTest 'gcloud compute target-http-proxies describe $BASENAME-lb-proxy --zone $ZONE --format="value(name)"' "EXPECTERROR"
section_close

section_open "Test google_compute_forwarding_rule does not exists"
    evalTest 'gcloud compute forwarding-rules describe $BASENAME-http-lb-forwarding-rule --zone $ZONE --format="value(name)"' "EXPECTERROR"
section_close

printf "$DIVIDER"
printf "CONGRATS!!!!!!! \n"
printf "You got the end the of your test with everything working. \n"
printf "$DIVIDER"    


```