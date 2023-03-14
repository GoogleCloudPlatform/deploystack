#!/bin/sh
# Don't judge me. I'm sure there's a one liner that could do it.  Sometimes
# you just want it to work and be done with it. 
re1='s/github.com\/GoogleCloudPlatform\///g'
re2='s/require//g'
re3='s/deploystack//g'
cat go.mod| grep deploystack | head -1 | sed  $re1 | sed  $re2 |  sed  $re3 |xargs  > versionDS
date +"%m-%d-%Y %r %Z" | xargs > buildTime