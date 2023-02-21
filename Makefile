test:
	go clean -testcache
	go test .  -cover
	go test ./tui/. -cover
	go test ./terraform/.  -cover
	go test ./github/.  -cover
	go test ./dstester/.  -cover
	go test ./gcloud/.  -cover
	
update: 
	gcloud config set project ds-tester-helper	
	cd tools/test_files_updater && ./update

creds.json:
	gcloud secrets versions access latest --secret=creds \
	--project=$(DEPLOYSTACK_TEST_PROJECT) > creds.json