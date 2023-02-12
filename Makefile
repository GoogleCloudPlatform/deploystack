test:
	go clean -testcache
	go test ./tui/. -cover
	go test ./gcloudtf/.  -cover
	go test ./dsgithub/.  -cover
	go test ./dstester/.  -cover
	go test .  -cover
	go test ./gcloud/.  -cover
	
 