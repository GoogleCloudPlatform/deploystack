// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"

	"github.com/GoogleCloudPlatform/deploystack/terraform"
)

func main() {
	b, err := terraform.Extract("../testdata/extracttest")
	if err != nil {
		log.Fatalf("couldn't extract from TF file: %s", err)
	}

	for _, v := range *b {
		fmt.Printf("%s %s %s\n", v.Kind, v.Type, v.Name)
		fmt.Printf("%+v\n", v.Attr)

	}
}
