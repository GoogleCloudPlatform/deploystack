# DeployStack

This project is to centralize all of the tools and processes to get terminal
interfaces for collecting information from users for use with DeployStack.

## Authoring

### TLDR;

Quick TLDR omitting testing, which you should totally do but for development and
whatnot you can just do this.

1. Clone this project.
1. Add your `main.tf` file
1. Edit the `deploystack.json` file
1. Edit the `deploystack.txt` file

Running `./deploystack install` should spin up the goapp and collect the config items the
stack needs to run through the deployment

Running `./deploystack uninstall` will destory the whole thing.

### Details

Authors are required to make or edit 4 files.

- `main.tf`
- `deploystack.json`
- `deploystack.txt`
- `test`

The following files need to be included but shoudln't need to be edited at all:

- `deploystack`
- `main.go`
- `test.yaml`

### `main.tf`

This is a standard terraform file with one adjustment for the DeployStack setup.
These files should have several import variables setup with the idea that the
golang helper will get them from the user.

```javascript
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

### `deploystack.json`

This config will be read by the golang helper to prompt the user to create a
tfvars file that will drive the terraform script.

```json
{
  "title": "Basic Title",
  "duration": 5,
  "collect_project": true,
  "collect_region": true,
  "region_type": "functions",
  "region_default": "us-central1",
  "collect_zone": true,
  "hard_settings": {
    "basename": "appprefix"
  },
  "custom_settings": [
    {
      "name": "nodes",
      "description": "Please enter the number of nodes",
      "default": "3"
    }
  ]
}
```

| Name            | Type    | Description                                                                          |
| --------------- | ------- | ------------------------------------------------------------------------------------ |
| title           | string  | You know what a title is                                                             |
| duration        | number  | An estimate as to how long this installation takes                                   |
| collect_project | boolean | Whether or not to walk the user through picking or creating a project.               |
| collect_region  | boolean | Whether or not to walk the user through picking a regions                            |
| region_type     | string  | Which product to select a region for                                                 |
|                 |         | Options: compute, run, functions                                                     |
| region_default  | string  | The highlighted and default choice for region.                                       |
| collect_zone    | string  | Whether or not to walk the user through picking a zone                               |
| hard_settings   |         | Hard Settings are for key value pairs to hardset and not get from the user.          |
|                 |         | `"basename":"appprefix"`                                                             |
| custom_settings |         | Custom Settings are collections of settings that we would like to prompt a user for. |
| name            | string  | The name of the variable                                                             |
| description     | string  | The description of the variable to prompt the user with                              |
| default         | string  | A default value for the variable.                                                    |
| options         | array   | An array of options to turn this into a custom select interface                      |
| prepend_project | bool    | Whether or not to prepend the project id to the default value. Useful for resources like buckets that have to have globally unique names.                       |

### UI Controls

#### Header
```json
  "title":"BASICLB",
  "duration":5,
```

```txt
This process will create the following:

	* Frontend - Cloud Run Service 
	* Middleware - Cloud Run Service
	* Backend - Cloud Sql MySQL instance 
	* Cache - Cloud Memorystore
	* Secrets - Cloud Secret Manager

All of these will spin up configured in a 3 tier application that delievers a
TODO app to show all of these pieces working together.  
```
![UI for Project Selector](assets/ui_header.png)


#### Project Selector
```json
  "collect_project":true
```
![UI for Project Selector](assets/ui_choose_project.png)

#### Region Selector
```json
  "collect_region":true,
  "region_type":"functions",
  "region_default":"us-central1",
```
![UI for Region Selector](assets/ui_change_region.png)

#### Zone Selector
```json
  "collect_zone":true
```
![UI for Zone Selector](assets/ui_choose_zone.png)


#### Custom Settings - no options
```json
"name":"nodes",
"description":"Please enter the number of nodes",
"default": "3"
```
![UI for Custom Settings with no options](assets/ui_custom_no_options.png)

#### Custom Settings - options
```json
"name":"nodes",
"description":"Please enter the number of nodes",
"default": "3"
"options": ["1", "2", "3"]
```
![UI for Custom Settings with options](assets/ui_custom_options.png)


### `deploystack.txt`

This file allows you to add a formatted description to the configuration to
print out to the user. Json files don't do well with newlines.

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


## Testing this Repo

In order to test the helper app in this repo, we need to do a fair amount of
manipulation of projects and what not. To faciliate that the tests require a
Service Account key json file. To faciliate this there is a script in
`tools/credsfile` that will create a service account, give it the right access
and service enablements, and export out a key file to use with testing.

This is not an offical Google product.
