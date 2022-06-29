// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hello

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/mailgun/mailgun-go/v4"

	cloudbuild "google.golang.org/api/cloudbuild/v1alpha1"
)

var gcloudFuncSourceDir = "."

// PubSubMessage is the payload of a Pub/Sub event.
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// HandleBuild will read in Billing overages and shut down machine's accordingly.
func HandleBuild(ctx context.Context, m PubSubMessage) error {
	data := string(m.Data)

	build := cloudbuild.Build{}
	if err := json.Unmarshal([]byte(data), &build); err != nil {
		return fmt.Errorf("cannot unmarshall Pub/Sub message: %s", err)
	}

	fmt.Printf("Build ID: %s Status: %s\n", build.Id, build.Status)

	if build.Status == "QUEUED" || build.Status == "WORKING" {
		fmt.Printf("Nothing to report. Status: %s\n", build.Status)
		return nil
	}

	if _, err := os.Stat("serverless_function_source_code"); !os.IsNotExist(err) {
		gcloudFuncSourceDir = "serverless_function_source_code"
	}

	// CloudFunctions is doing weird things with compilations and filesystems
	dat, err := os.ReadFile(gcloudFuncSourceDir + "/config.json")
	if err != nil {
		return fmt.Errorf("cannot read config file: %s", err)
	}

	cfg := MailgunConfig{}
	if err := json.Unmarshal(dat, &cfg); err != nil {
		return fmt.Errorf("cannot convert json: %s", err)
	}

	if !strings.Contains(build.Substitutions["TRIGGER_NAME"], "Test-Procedure") {
		fmt.Printf("Don't notify\n")
		return nil
	}

	msg := Message{
		fmt.Sprintf("%s: Project: %s Test %s finished", build.Status, build.ProjectId, build.Id),
		"",
	}

	d := getDuration(build)

	tmpdata := struct {
		URL      string
		Status   string
		Duration string
		Log      template.HTML
	}{
		build.LogUrl,
		build.Status,
		d,
		"",
	}

	var tpl bytes.Buffer

	if build.Status == "SUCCESS" {
		temp := template.Must(template.ParseFiles(gcloudFuncSourceDir + "/template/success.html"))
		if err := temp.Execute(&tpl, tmpdata); err != nil {
			fmt.Printf("error in dealing with template: %s\n", err)
		}
		msg.Body = tpl.String()
		sendMessage(msg, cfg)
		fmt.Printf("%s\n", build.Status)
		return nil
	}

	log, err := getBuildLog(build)
	if err != nil {
		fmt.Printf("could not retrieve log for inclusion: %s", err)
	}

	log = massageLog(log)
	tmpdata.Log = template.HTML(log)

	temp := template.Must(template.ParseFiles(gcloudFuncSourceDir + "/template/unsuccess.html"))
	temp.Execute(&tpl, tmpdata)

	msg.Body += tpl.String()

	sendMessage(msg, cfg)

	fmt.Printf("log url: %s\n", build.LogUrl)
	return nil
}

func massageLog(log string) string {
	log = strings.ReplaceAll(log, "\n", "<br />")
	log = strings.ReplaceAll(log, "[0m", "</span></span>")
	log = strings.ReplaceAll(log, "[0;36m", "<span style=\"color: cyan\">")
	log = strings.ReplaceAll(log, "[0;32m", "<span style=\"color: green\">")
	log = strings.ReplaceAll(log, "[1m", "<span style=\"color: yellow\">")
	log = strings.ReplaceAll(log, "[1;36m", "<span style=\"color: cyan; font-weight: 900\">")
	log = strings.ReplaceAll(log, "[4;36m", "<span style=\"color: cyan; text-decoration: underline\">")
	log = strings.ReplaceAll(log, "[0;31m", "<span style=\"color: red\">")
	log = strings.ReplaceAll(log, "[32m", "<span style=\"color: green\">")
	log = strings.ReplaceAll(log, "[31m", "<span style=\"color: red\">")
	log = strings.ReplaceAll(log, "[4m", "<span style=\"color: cyan\">")

	return log
}

func getBuildLog(build cloudbuild.Build) (string, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()

	// Sets the name for the new bucket.
	bucketName := strings.ReplaceAll(build.LogsBucket, "gs://", "")
	objectName := fmt.Sprintf("log-%s.txt", build.Id)
	// Creates a Bucket instance.
	bucket := client.Bucket(bucketName)

	obj := bucket.Object(objectName)
	rdr, err := obj.NewReader(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create reader: %v", err)
	}
	defer rdr.Close()

	b, err := io.ReadAll(rdr)
	if err != nil {
		return "", fmt.Errorf("failed to create read file: %v", err)
	}

	return string(b), nil
}

type Message struct {
	Subject string
	Body    string
}

type MailgunConfig struct {
	Key    string `json:"MAILGUN_API_KEY"`
	Domain string `json:"MAILGUN_DOMAIN"`
	From   string `json:"MAILGUN_FROM"`
	To     string `json:"MAILGUN_TO"`
}

func getDuration(build cloudbuild.Build) string {
	end, err := time.Parse("2006-01-02T15:04:05.999999Z", build.FinishTime)
	if err != nil {
		fmt.Printf("count not parse time: %s error: %s", build.FinishTime, err)
	}

	begin, err := time.Parse("2006-01-02T15:04:05.999999Z", build.StartTime)
	if err != nil {
		fmt.Printf("count not parse time: %s error: %s", build.FinishTime, err)
	}

	dur := end.Sub(begin)

	d := humanizeDuration(dur)

	return d
}

func sendMessage(m Message, c MailgunConfig) error {
	mg := mailgun.NewMailgun(c.Domain, c.Key)

	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(c.From, m.Subject, m.Body, c.To)
	message.SetHtml(m.Body)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err := mg.Send(ctx, message)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)

	return nil
}

func humanizeDuration(duration time.Duration) string {
	if duration.Seconds() < 60.0 {
		return fmt.Sprintf("%d seconds", int64(duration.Seconds()))
	}
	if duration.Minutes() < 60.0 {
		remainingSeconds := math.Mod(duration.Seconds(), 60)
		return fmt.Sprintf("%d minutes %d seconds", int64(duration.Minutes()), int64(remainingSeconds))
	}
	if duration.Hours() < 24.0 {
		remainingMinutes := math.Mod(duration.Minutes(), 60)
		remainingSeconds := math.Mod(duration.Seconds(), 60)
		return fmt.Sprintf("%d hours %d minutes %d seconds",
			int64(duration.Hours()), int64(remainingMinutes), int64(remainingSeconds))
	}
	remainingHours := math.Mod(duration.Hours(), 24)
	remainingMinutes := math.Mod(duration.Minutes(), 60)
	remainingSeconds := math.Mod(duration.Seconds(), 60)
	return fmt.Sprintf("%d days %d hours %d minutes %d seconds",
		int64(duration.Hours()/24), int64(remainingHours),
		int64(remainingMinutes), int64(remainingSeconds))
}
