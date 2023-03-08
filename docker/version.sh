#!/bin/sh
re='s/github.com\/GoogleCloudPlatform\///g'
cat go.mod| grep deploystack | head -1 | sed  $re | xargs  > versionDS
cat go.mod| grep deploystack/tui  |  sed  $re | xargs > versionTUI
cat go.mod| grep deploystack/gcloud  |  sed  $re | xargs > versionGcloud
cat go.mod| grep deploystack/config  |  sed  $re | xargs > versionConfig
cat go.mod| grep deploystack/terraform  |  sed  $re | xargs > versionTerraform
date +"%m-%d-%Y %r %Z" | xargs > buildTime