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
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	cloudbuild "cloud.google.com/go/cloudbuild/apiv1"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/iterator"
	option "google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	cloudbuildpb "google.golang.org/genproto/googleapis/devtools/cloudbuild/v1"
)

// PubSubMessage is the payload of a Pub/Sub event.
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

func RecordTest(ctx context.Context, msg PubSubMessage) error {
	project := os.Getenv("PROJECT")

	if project == "" {
		return fmt.Errorf("project not set")
	}

	return collectAndRegister(project)
}

func collectAndRegister(project string) error {
	ctx := context.Background()
	var err error

	// Get Sheets client
	filename := findPath("credentials.json")
	m := make(map[string]string)

	dat, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error opening credentials file: %s", err)
	}

	if err := json.Unmarshal(dat, &m); err != nil {
		return fmt.Errorf("error parsing credentials file: %s", err)
	}

	conf := &jwt.Config{
		Email:        m["client_email"],
		PrivateKey:   []byte(m["private_key"]),
		PrivateKeyID: m["private_key_id"],
		TokenURL:     m["token_uri"],
		Scopes: []string{
			"https://www.googleapis.com/auth/spreadsheets",
		},
	}

	sheetID := m["sheet_id"]

	client := conf.Client(oauth2.NoContext)

	sheetsSVC, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("Unable to retrieve Sheets client: %v", err)
	}

	// Get Build client
	buildSVC, err := cloudbuild.NewClient(ctx)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("failed to create build client: %v", err))
	}
	defer buildSVC.Close()

	// Do the actual things we need todo.
	row, err := querySheet(*sheetsSVC, sheetID, project)
	if err != nil {
		return fmt.Errorf("sheets: failed read from sheet: %v", err)
	}

	r, err := getBuild(ctx, project, *buildSVC)
	if err != nil {
		return fmt.Errorf("cloudbuild: cannot get info from build: %v", err)
	}

	if err = writeStatiiToSheet(r, *sheetsSVC, sheetID, row); err != nil {
		return fmt.Errorf("sheet: can't write to sheet %v", err)
	}

	return nil
}

func getBuild(ctx context.Context, project string, svc cloudbuild.Client) (TestRun, error) {
	build, err := getLastTestResult(ctx, project, svc)
	if err != nil {
		return TestRun{}, fmt.Errorf("failed to get tests: %v", err)
	}

	tr := TestRun{
		ProjectID: project,
		BuildID:   build.Id,
		Status:    build.Status.String(),
		Last:      build.FinishTime.AsTime(),
		Repo:      build.GetSubstitutions()["REPO_NAME"],
	}
	return tr, nil
}

func findPath(name string) string {
	gcloudFuncSourceDir := "./"
	result := ""

	if _, err := os.Stat("serverless_function_source_code"); !os.IsNotExist(err) {
		gcloudFuncSourceDir = "serverless_function_source_code/"
	}

	if _, err := os.Stat("../" + name); !os.IsNotExist(err) {
		gcloudFuncSourceDir = "../"
	}
	result = fmt.Sprintf("%s%s", gcloudFuncSourceDir, name)

	return result
}

type TestRun struct {
	ProjectID string
	BuildID   string
	Status    string
	Last      time.Time
	Repo      string
}

func getLastTestResult(ctx context.Context, project string, svc cloudbuild.Client) (*cloudbuildpb.Build, error) {
	req := &cloudbuildpb.ListBuildsRequest{ProjectId: project}
	fmt.Printf("REQ %+v\n", req)
	result := &cloudbuildpb.Build{}
	it := svc.ListBuilds(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		if resp.Substitutions["TRIGGER_NAME"] == "Test-Procedure" {
			result = resp
			break
		}

		if resp.Substitutions["TRIGGER_NAME"] == "Test-Procedure-push" {
			result = resp
			break
		}

	}

	return result, nil
}

func querySheet(svc sheets.Service, ID string, q string) (int, error) {
	rownumb := -1
	ranges := []string{"Sheet1"}
	includeGridData := true

	resp, err := svc.Spreadsheets.Get(ID).Ranges(ranges...).IncludeGridData(includeGridData).Context(context.Background()).Do()
	if err != nil {
		return rownumb, err
	}

	rownumb = len(resp.Sheets[0].Data[0].RowData)
	for i, row := range resp.Sheets[0].Data[0].RowData {
		for _, cell := range row.Values {
			if cell.FormattedValue == q {
				rownumb = i
			}
		}
	}

	return rownumb, nil
}

func writeStatiiToSheet(t TestRun, svc sheets.Service, ID string, row int) error {
	writeRange := fmt.Sprintf("A%d", row+1)
	var vr sheets.ValueRange

	myval := []interface{}{t.ProjectID, t.Last, t.BuildID, t.Status, t.Repo}
	vr.Values = append(vr.Values, myval)

	if _, err := svc.Spreadsheets.Values.Update(ID, writeRange, &vr).ValueInputOption("RAW").Do(); err != nil {
		return fmt.Errorf("sheets: failed to get old records from spreadsheets %s", err)
	}

	br := sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{{
			SortRange: &sheets.SortRangeRequest{
				Range: &sheets.GridRange{
					SheetId:          0,
					StartRowIndex:    1,
					StartColumnIndex: 0,
					EndColumnIndex:   0,
				},
				SortSpecs: []*sheets.SortSpec{
					{
						SortOrder:      "ASCENDING",
						DimensionIndex: 0,
					},
				},
			},
		}},
	}

	if _, err := svc.Spreadsheets.BatchUpdate(ID, &br).Do(); err != nil {
		return fmt.Errorf("failed get sort spreadsheet: %v", err)
	}

	return nil
}
