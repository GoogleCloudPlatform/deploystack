# DeployStack - Authoring


### TLDR;

Quick TLDR omitting testing, which you should totally do but for development and
whatnot you can just do this.

1. Create a Folder to hold your stack.
1. Add your `main.tf` file
1. Create a .deploystack folder
1. Add a `deploystack.yaml` file

Running `./deploystack install` should spin up the goapp and collect the config items the
stack needs to run through the deployment

Running `./deploystack uninstall` will destory the whole thing.

### Details

Authors are required to make or edit 4 files.

- `main.tf`
- `.deploystack/deploystack.yaml`
- `.deploystack/test`


### `main.tf`

This is a standard terraform file with one adjustment for the DeployStack setup.
These files should have several import variables setup with the idea that the
golang helper will get them from the user.

```hcl
variable "project_id" {
  type = string
}

variable "project_number" {
  type = string
}

variable "region" {
  type = string
}

variable "zone" {
  type = string
}
```

_DeployStack will also work with the more standard convention of having 
separate main.tf, variables.tf, output.tf, etc files._

### `deploystack.yaml`

This config will be read by the golang helper to prompt the user to create a
tfvars file that will drive the terraform script.

```yaml
title: Basic Title
duration: 5
collect_project: true
collect_region: true
region_type: functions
region_default: us-central1
collect_zone: true
hard_settings:
  basename: appprefix
custom_settings:
- name: nodes
  description: Please enter the number of nodes
  options:
  - roles/reviewer|Project Reviewer
  - roles/owner|Project Owner
  - roles/vison.reader|Cloud Vision Reader
  default: roles/owner|Project Owner
projects:
  allow_duplicates: false
  items:
  - variable_name: project_id
    user_prompt: Choose a project to use for this application
    set_as_default: true
  - variable_name: project_id_2
    user_prompt: Choose a second project to use for this application
    set_as_default: false
```

_JSON is also allowed_


#### DeployStack Config Settings
These are now documented in [deploystack/config](/config).


### `messages/description.txt`

DEPRECATED: This file allows you to add a formatted description to the configuration to
print out to the user. Json files don't do well with newlines. Using Description in
deploystack.yaml is now prefered

### `test`

Test is a shell script that tests the individual pieces of the infrastructure
and tests the desired state at the end of the install.

There are a few functions in the template test file that will help you run one
of these.

- `section_open` - a display function that hellps communicate what is going on.
- `section_close` - paired with section_open
- `evaltest` - take a gcloud command and a desired outcome to make test assertions

```bash
# Setup variables here
source globals
get_project_id PROJECT
get_project_number PROJECT_NUMBER $PROJECT
REGION=us-central1
ZONE=us-central1-a
BASENAME=basiclb
SIZE=3

# Make sure that project is hard set
gcloud config set project ${PROJECT}

# spin up terraform with variables plugged in to build the infrastructure
terraform init
terraform apply -auto-approve -var project_id="${PROJECT}" -var project_number="${PROJECT_NUMBER}" -var region="${REGION}" -var zone="${ZONE}" -var basename="${BASENAME}" -var nodes="${SIZE}"

# You might hace to do some editing here to make these tests work
section_open "Test Managed Instance Group"
    evalTest 'gcloud compute instance-groups managed describe $BASENAME-mig --zone $ZONE --format="value(name)"'  $BASENAME-mig

    COUNT=$(gcloud compute instances list --format="value(name)" | grep $BASENAME-mig | wc -l | xargs)

    if [ $COUNT -ne $SIZE ]
    then
        printf "Halting - error: expected $SIZE instances of GCE got $COUNT  \n"
        exit 1
    else
         printf "number of GCE instances is ok \n"
    fi

section_close

# But in a lot of cases we can just use eval test with a gcloud command and a
# desrired result.
section_open "Test Instance Template"
    evalTest 'gcloud compute instance-templates describe $BASENAME-template --format="value(name)"'  $BASENAME-template
section_close

..

# Now run a destroy operation.
terraform destroy -auto-approve -var project_id="${BASENAME}" -var project_number="${PROJECT_NUMBER}" -var region="${REGION}" -var zone="${ZONE}" -var basename="${BASENAME}" -var nodes="${SIZE}"

# Test all of the parts are destroyed
section_open "Test Managed Instance Group doesn't exist"
    evalTest 'gcloud compute instance-groups managed describe $BASENAME-mig --zone $ZONE --format="value(name)"'  "EXPECTERROR"
section_close

printf "$DIVIDER"
printf "CONGRATS!!!!!!! \n"
printf "You got the end the of your test with everything working. \n"
printf "$DIVIDER"
```