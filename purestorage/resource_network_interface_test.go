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
func TestAccResourcePureNetworkInterface_create(t *testing.T) {
	ifname := "vir1"
	address1 := "192.168.6.3"
	address2 := "192.168.6.4"
	gateway := "192.168.6.1"
	netmask := "255.255.255.0"
	resource_name := fmt.Sprintf("purefa_network_interface.tfnetworkinterface%stest", strings.Replace(ifname, ".", "_", -1))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPureNetworkInterfaceConfig(ifname, address1, gateway, netmask, true, 1500),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resource_name, "ifname", ifname),
					resource.TestCheckResourceAttr(resource_name, "address", address1),
					resource.TestCheckResourceAttr(resource_name, "gateway", gateway),
					resource.TestCheckResourceAttr(resource_name, "netmask", netmask),
					resource.TestCheckResourceAttr(resource_name, "enabled", "true"),
					resource.TestCheckResourceAttr(resource_name, "mtu", "1500"),
				),
			},
			{
				Config: testAccCheckPureNetworkInterfaceConfig(ifname, address1, gateway, netmask, false, 1500),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resource_name, "ifname", ifname),
					resource.TestCheckResourceAttr(resource_name, "address", address1),
					resource.TestCheckResourceAttr(resource_name, "gateway", gateway),
					resource.TestCheckResourceAttr(resource_name, "netmask", netmask),
					resource.TestCheckResourceAttr(resource_name, "enabled", "false"),
					resource.TestCheckResourceAttr(resource_name, "mtu", "1500"),
				),
			},
			{
				Config: testAccCheckPureNetworkInterfaceConfig(ifname, address2, gateway, netmask, true, 9000),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resource_name, "ifname", ifname),
					resource.TestCheckResourceAttr(resource_name, "address", address2),
					resource.TestCheckResourceAttr(resource_name, "gateway", gateway),
					resource.TestCheckResourceAttr(resource_name, "netmask", netmask),
					resource.TestCheckResourceAttr(resource_name, "enabled", "true"),
					resource.TestCheckResourceAttr(resource_name, "mtu", "1500"),
				),
			},
		},
	})
}

func testAccCheckPureNetworkInterfaceConfig(ifname string, address string, gateway string, netmask string, enabled bool, mtu int) string {
	return fmt.Sprintf(`
			resource "purefa_network_interface" "tfnetworkinterface%stest" {
				name = "%s"
				address = "%s"
				gateway = "%s"
				netmask = "%s"
				enabled = %t
				mtu = %d
			}`, strings.Replace(ifname, ".", "_", -1), ifname, address, gateway, netmask, enabled, mtu)

}
