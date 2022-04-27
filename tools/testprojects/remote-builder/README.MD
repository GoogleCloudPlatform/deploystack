# DeployStack TestProjects

Since Deploystack spins up completel architectures it's important to have a 
testing project associated with each one. The project's test script can be 
called from test.yaml in a scheduled build. This way we schedule testing for 
each project on creation. Unfortunately we have to add github repos to Build 
Triggers in the we bconsole only. 


Steps to run
* Type `./create [name of stackdeploy]`
* Follow link at bottom of process. 
* Connect this Build setup with Github repo.
* Then return here and type  `./update [name of stackdeploy]`