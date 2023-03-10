#!/bin/sh
re='s/github.com\/GoogleCloudPlatform\///g'
cat go.mod| grep deploystack | head -1 | sed  $re | xargs  > versionDS
date +"%m-%d-%Y %r %Z" | xargs > buildTime