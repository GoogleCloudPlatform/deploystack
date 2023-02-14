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