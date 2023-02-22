// Copyright 2023 Google LLC
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

package simple_example

import (
	"net/http"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/gcloud"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/tft"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestSimpleExample(t *testing.T) {
	example := tft.NewTFBlueprintTest(t)

	example.DefineVerify(func(assert *assert.Assertions) {
		// default verify reruns plan and asserts no diff
		//TODO: There seems to be a permadiff with module.simple.google_monitoring_dashboard.dashboard
		// example.DefaultVerify(assert)

		// check if cloud run service exists
		projectID := example.GetTFSetupStringOutput("project_id")
		vmName := example.GetStringOutput("xwiki_vm_name")

		// sample assertion to validate VM is running
		opVM := gcloud.Runf(t, "compute instances describe %s --zone us-central1-a --project %s", vmName, projectID)
		assert.Equal("RUNNING", opVM.Get("status").String(), "expected XWiki VM to be running")

		// sample e2e to assert app is working
		wikiURL := example.GetStringOutput("xwiki_url")
		isServing := func() (bool, error) {
			resp, err := http.Get(wikiURL)
			if err != nil || resp.StatusCode != 200 {
				// retry if err or status not 200
				return true, nil
			}
			return false, nil
		}
		utils.Poll(t, isServing, 10, time.Second*10)
	})
	example.Test()
}
