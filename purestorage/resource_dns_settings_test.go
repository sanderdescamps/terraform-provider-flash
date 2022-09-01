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
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Configure DNS settings
func TestAccResourcePureDnsSettings_create(t *testing.T) {
	nameservers1 := []string{"10.0.0.1"}
	nameservers2 := []string{"10.0.0.1", "1.0.0.1"}
	domain := "testdrive.local"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPureDnsSettingsConfig(nameservers1, domain),
			},
			{
				Config: testAccCheckPureDnsSettingsConfig(nameservers2, domain),
			},
			{
				Config: testAccCheckPureDnsSettingsConfig(nameservers1, nil),
			},
		},
	})
}

func testAccCheckPureDnsSettingsConfig(nameservers []string, domain interface{}) string {
	if domain == nil {
		return fmt.Sprintf(`
			resource "purefa_dns_settings" "tfdnssettingstest" {
					nameservers = split(";","%s")
			}`, strings.Join(nameservers, ";"))
	} else {
		return fmt.Sprintf(`
			resource "purefa_dns_settings" "tfdnssettingstest" {
					nameservers = split(";","%s")
					domain = "%s"
			}`, strings.Join(nameservers, ";"), domain.(string))
	}

}
