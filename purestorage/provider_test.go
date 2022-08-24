/*
   Copyright 2018 David Evans

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package purestorage

import (
	"os"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testInitOnce = sync.Once{}

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testInitOnce.Do(
		func() {
			testAccProvider = Provider()
			testAccProviders = map[string]*schema.Provider{
				"purestorage": testAccProvider,
			}
		},
	)

}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	target := os.Getenv("PURE_TARGET")
	username := os.Getenv("PURE_USERNAME")
	password := os.Getenv("PURE_PASSWORD")
	apitoken := os.Getenv("PURE_APITOKEN")
	if target == "" {
		t.Fatalf("PURE_TARGET must be set for acceptance tests")
	}
	if (apitoken == "") && (username == "") && (password == "") {
		t.Fatalf("PURE_USERNAME and PURE_PASSWORD or PURE_APITOKEN must be set for acceptance tests")
	}
	if (username != "") && (password == "") {
		t.Fatalf("PURE_PASSWORD must be set if PURE_USERNAME is set for acceptance tests")
	}
}

func testAccProviderMeta(t *testing.T) (interface{}, error) {
	t.Helper()
	d := schema.TestResourceDataRaw(t, testAccProvider.Schema, make(map[string]interface{}))
	return providerConfigure(d)
}
