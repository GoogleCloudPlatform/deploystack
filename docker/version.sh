#!/bin/sh
cat go.mod| grep github.com/GoogleCloudPlatform/deploystack | head -1 | sed  's/github.com\/GoogleCloudPlatform\///g' | xargs  > versionDS
cat go.mod| grep github.com/GoogleCloudPlatform/deploystack/tui  |  sed  's/github.com\/GoogleCloudPlatform\///g' | xargs > versionTUI
cat go.mod| grep github.com/GoogleCloudPlatform/deploystack/gcloud  |  sed  's/github.com\/GoogleCloudPlatform\///g' | xargs > versionGcloud
cat go.mod| grep github.com/GoogleCloudPlatform/deploystack/config  |  sed  's/github.com\/GoogleCloudPlatform\///g' | xargs > versionConfig
date +"%m-%d-%Y %r %Z" | xargs > buildTime