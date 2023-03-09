# DeployStack
[![GoDoc](https://godoc.org/github.com/GoogleCloudPlatform/deploystack?status.svg)](https://godoc.org/github.com/GoogleCloudPlatform/deploystack)


[DeployStack](http://deploystack.dev) is a one click solution for running 
Terraform projects for Google Cloud Platfom using [Cloud Shell](https://cloud.google.com/shell) 
It uses [Open in Cloud Shell](https://cloud.google.com/shell/docs/open-in-cloud-shell) 
to guide users from a link to a series or questions to help them install a 
Terraform solution in their own Google Cloud Platform project space - prompting 
them to choose answers to questions like "What [datacenter] zone do you want to 
install in?" And presenting the options to guide them to pic the settings that 
are right for them. 

![DeployStack UX](/assets/demo.gif)

For technical reasons, at this time, it is limited to working with github repos 
owned by [Google Cloud Platform](https://github.com/GoogleCloudPlatform). You 
can see a list of DeployStack projects on [cloud.google.com](https://cloud.google.com/shell/docs/cloud-shell-tutorials/deploystack).

## This Codebase

This project is to centralize all of the tools and processes to get terminal
interfaces for collecting information from users for use with DeployStack. 
Ultimately this codebase creates an executable that runs on Cloud Shell and 
works with other tools to drive the experience.

It's broken up into packages with different responsibilities: 



## Authoring
Authoring information has been moved to [deploystack/AUTHORING.MD](/AUTHORING.MD).

<dl>
  <dt>deploystack</dt>
  <dd>A top level package that ties together user i/o for the executable and 
  passes information to the other packages</dd>
  <dt>deploystack/config</dt>
  <dd>The basic information schema that runs all of the other parts of the 
  executable.</dd>
  <dt>deploystack/gcloud</dt>
  <dd>Communication with Google Cloud via the Go SDK to get things like
  region and zone lists</dd>
  <dt>deploystack/github</dt>
  <dd>Communication with Github for cloning and getting other metadata about 
  projects</dd>
  <dt>deploystack/terraform</dt>
  <dd>Introspection of Terraform files for tooling and other metadata</dd>
  <dt>deploystack/tui</dt>
  <dd>A terminal user interface that is dynamically built based on config files
  rendered using <a href="https://github.com/charmbracelet/bubbletea">BubbleTea</a> , 
  <a href="https://github.com/charmbracelet/lipgloss">LipGloss</a>, and 
  <a href="https://github.com/charmbracelet/bubbles">Bubbles</a> </dd>
</dl>

## Testing this Repo

In order to test the helper app in this repo, we need to do a fair amount of
manipulation of projects and what not. To faciliate that the tests require a
Service Account key json file. To faciliate this there is a script in
`tools/credsfile` that will create a service account, give it the right access
and service enablements, and export out a key file to use with testing.

This is not an offical Google product.
