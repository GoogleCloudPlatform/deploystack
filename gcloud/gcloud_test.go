package gcloud

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"google.golang.org/api/option"
)

var (
	projectID        = ""
	creds            map[string]string
	opts             = option.WithCredentialsFile("")
	ctx              = context.Background()
	defaultUserAgent = "deploystack/testing"
)

func TestMain(m *testing.M) {
	var err error
	opts = option.WithCredentialsFile("../creds.json")

	dat, err := os.ReadFile("../creds.json")
	if err != nil {
		log.Fatalf("unable to handle the json config file: %v", err)
	}

	json.Unmarshal(dat, &creds)

	projectID = creds["project_id"]
	if err != nil {
		log.Fatalf("could not get environment project id: %s", err)
	}
	code := m.Run()
	os.Exit(code)
}

func randSeq(n int) string {
	rand.Seed(time.Now().Unix())

	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func interfaceTester() UIClient {
	var r UIClient
	c := NewClient(ctx, defaultUserAgent, opts)
	r = &c
	return r
}
