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

const testAccCheckPureAlertRecipientResourceName = "purefa_dns_settings.tfdnssettingstest"

// Create a hostgroup
func TestAccResourcePureAlertRecipient(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: strings.Join([]string{
					testAccCheckPureAlertRecipientConfig("user1@example.com", true),
					testAccCheckPureAlertRecipientConfig("user2@example.com", false)}, "\n"),
			},
			{
				Config: strings.Join([]string{
					testAccCheckPureAlertRecipientConfig("user1@example.org", true),
					testAccCheckPureAlertRecipientConfig("user2@example.com", true)}, "\n"),
			},
			{
				Config: strings.Join([]string{
					testAccCheckPureAlertRecipientConfig("user1@example.org", false),
					testAccCheckPureAlertRecipientConfig("user2@example.com", true)}, "\n"),
			},
		},
	})
}

func testAccCheckPureAlertRecipientConfig(email string, enabled bool) string {
	name := strings.Split(email, "@")[0]
	return fmt.Sprintf(`
			resource "purefa_alert_recipient" "tfalert_%s" {
					email = "%s"
					enabled = %t

			}`, name, email, enabled)
}
