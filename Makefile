test:
	go clean -testcache
	go test .  -cover
	go test ./tui/. -cover
	go test ./gcloudtf/.  -cover
	go test ./dsgithub/.  -cover
	go test ./dstester/.  -cover
	go test ./gcloud/.  -cover
	
update: 
	cd tools/test_files_updater && ./update

# Will only work if you have 
creds.json:
		gcloud secrets versions access latest --secret=creds \
		--project=$(DEPLOYSTACK_TEST_PROJECT) > creds.json