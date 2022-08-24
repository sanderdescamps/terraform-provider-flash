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
	"math/rand"
	"testing"

	"github.com/devans10/pugo/flasharray"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testAccCheckPureProtectiongroupResourceName = "purestorage_protectiongroup.tfprotectiongrouptest"

// Create a protectiongroup
func TestAccResourcePureProtectiongroup_create(t *testing.T) {
	rInt := rand.Int()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPureProtectiongroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPureProtectiongroupConfigBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureProtectiongroupExists(testAccCheckPureProtectiongroupResourceName, true),
				),
			},
		},
	})
}

func TestAccResourcePureProtectiongroup_create_withHosts(t *testing.T) {
	rInt := rand.Int()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPureProtectiongroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPureProtectiongroupConfigWithHosts(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureProtectiongroupExists(testAccCheckPureProtectiongroupResourceName, true),
					testAccCheckPureHostExists("purestorage_host.tfpgrouptesthost", true),
					testAccCheckPureProtectiongroupHosts(testAccCheckPureProtectiongroupResourceName, "tfpgrouptesthost", true),
				),
			},
		},
	})
}

func TestAccResourcePureProtectiongroup_create_withHostgroups(t *testing.T) {
	rInt := rand.Int()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPureProtectiongroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPureProtectiongroupConfigWithHostgroups(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureProtectiongroupExists(testAccCheckPureProtectiongroupResourceName, true),
					testAccCheckPureHostgroupExists("purestorage_hostgroup.tfpgrouptesthgroup", true),
					testAccCheckPureProtectiongroupHostgroups(testAccCheckPureProtectiongroupResourceName, "tfpgrouptesthgroup", true),
				),
			},
		},
	})
}

func TestAccResourcePureProtectiongroup_create_withVolumes(t *testing.T) {
	rInt := rand.Int()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPureProtectiongroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPureProtectiongroupConfigWithVolumes(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureProtectiongroupExists(testAccCheckPureProtectiongroupResourceName, true),
					testAccCheckPureVolumeExists("purestorage_volume.tfpgrouptest-volume", true),
					testAccCheckPureProtectiongroupVolumes(testAccCheckPureProtectiongroupResourceName, fmt.Sprintf("tfpgrouptest-volume-%d", rInt), true),
				),
			},
		},
	})
}

func TestAccResourcePureProtectiongroup_create_withSchedule(t *testing.T) {
	rInt := rand.Int()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPureProtectiongroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPureProtectiongroupConfigWithSchedule(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureProtectiongroupExists(testAccCheckPureProtectiongroupResourceName, true),
				),
			},
		},
	})
}

func TestAccResourcePureProtectiongroup_create_withRetention(t *testing.T) {
	rInt := rand.Int()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPureProtectiongroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPureProtectiongroupConfigWithRetention(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureProtectiongroupExists(testAccCheckPureProtectiongroupResourceName, true),
				),
			},
		},
	})
}

func TestAccResourcePureProtectiongroup_update_withHosts(t *testing.T) {
	rInt := rand.Int()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPureProtectiongroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPureProtectiongroupConfigBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureProtectiongroupExists(testAccCheckPureProtectiongroupResourceName, true),
				),
			},
			{
				Config: testAccCheckPureProtectiongroupConfigWithHosts(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureProtectiongroupExists(testAccCheckPureProtectiongroupResourceName, true),
					testAccCheckPureHostExists("purestorage_host.tfpgrouptesthost", true),
					testAccCheckPureProtectiongroupHosts(testAccCheckPureProtectiongroupResourceName, "tfpgrouptesthost", true),
				),
			},
		},
	})
}

func TestAccResourcePureProtectiongroup_update_withHostgroups(t *testing.T) {
	rInt := rand.Int()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPureProtectiongroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPureProtectiongroupConfigBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureProtectiongroupExists(testAccCheckPureProtectiongroupResourceName, true),
				),
			},
			{
				Config: testAccCheckPureProtectiongroupConfigWithHostgroups(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureProtectiongroupExists(testAccCheckPureProtectiongroupResourceName, true),
					testAccCheckPureHostgroupExists("purestorage_hostgroup.tfpgrouptesthgroup", true),
					testAccCheckPureProtectiongroupHostgroups(testAccCheckPureProtectiongroupResourceName, "tfpgrouptesthgroup", true),
				),
			},
		},
	})
}

func TestAccResourcePureProtectiongroup_update_withVolumes(t *testing.T) {
	rInt := rand.Int()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPureProtectiongroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPureProtectiongroupConfigBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureProtectiongroupExists(testAccCheckPureProtectiongroupResourceName, true),
				),
			},
			{
				Config: testAccCheckPureProtectiongroupConfigWithVolumes(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureProtectiongroupExists(testAccCheckPureProtectiongroupResourceName, true),
					testAccCheckPureVolumeExists("purestorage_volume.tfpgrouptest-volume", true),
					testAccCheckPureProtectiongroupVolumes(testAccCheckPureProtectiongroupResourceName, fmt.Sprintf("tfpgrouptest-volume-%d", rInt), true),
				),
			},
		},
	})
}

func TestAccResourcePureProtectiongroup_update_withSchedule(t *testing.T) {
	rInt := rand.Int()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPureProtectiongroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPureProtectiongroupConfigBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureProtectiongroupExists(testAccCheckPureProtectiongroupResourceName, true),
				),
			},
			{
				Config: testAccCheckPureProtectiongroupConfigWithSchedule(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureProtectiongroupExists(testAccCheckPureProtectiongroupResourceName, true),
				),
			},
		},
	})
}

func TestAccResourcePureProtectiongroup_update_withRetention(t *testing.T) {
	rInt := rand.Int()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPureProtectiongroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPureProtectiongroupConfigBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureProtectiongroupExists(testAccCheckPureProtectiongroupResourceName, true),
				),
			},
			{
				Config: testAccCheckPureProtectiongroupConfigWithRetention(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureProtectiongroupExists(testAccCheckPureProtectiongroupResourceName, true),
				),
			},
		},
	})
}

func testAccCheckPureProtectiongroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*flasharray.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "purestorage_protectiongroup" {
			continue
		}

		_, err := client.Protectiongroups.GetProtectiongroup(rs.Primary.ID, nil)
		if err != nil {
			return nil
		}
		return fmt.Errorf("protectiongroup '%s' stil exists", rs.Primary.ID)
	}

	return nil
}

func testAccCheckPureProtectiongroupExists(n string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*flasharray.Client)
		name := rs.Primary.Attributes["name"]
		_, err := client.Protectiongroups.GetProtectiongroup(name, nil)
		if err != nil {
			if exists {
				return fmt.Errorf("protectiongroup does not exist: %s", n)
			}
			return nil
		}
		return nil
	}
}

func testAccCheckPureProtectiongroupHosts(n string, host string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*flasharray.Client)
		name := rs.Primary.Attributes["name"]
		p, err := client.Protectiongroups.GetProtectiongroup(name, nil)
		if err != nil {
			return fmt.Errorf("protectiongroup does not exist: %s", n)
		}
		if stringInSlice(host, p.Hosts) {
			if exists {
				return nil
			}
			return fmt.Errorf("host %s still connected to Protection Group %s", host, name)
		}
		if exists {
			return fmt.Errorf("host %s not connected to Protection Group %s", host, name)
		}
		return nil
	}
}

func testAccCheckPureProtectiongroupHostgroups(n string, hostgroup string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*flasharray.Client)
		name := rs.Primary.Attributes["name"]
		p, err := client.Protectiongroups.GetProtectiongroup(name, nil)
		if err != nil {
			return fmt.Errorf("protectiongroup does not exist: %s", name)
		}
		if stringInSlice(hostgroup, p.Hgroups) {
			if exists {
				return nil
			}
			return fmt.Errorf("hostgroup %s still connected to Protection Group %s", hostgroup, name)
		}
		if exists {
			return fmt.Errorf("hostgroup %s not connected to Protection Group %s", hostgroup, name)
		}
		return nil
	}
}

func testAccCheckPureProtectiongroupVolumes(n string, volume string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*flasharray.Client)
		name := rs.Primary.Attributes["name"]
		p, err := client.Protectiongroups.GetProtectiongroup(name, nil)
		if err != nil {
			return fmt.Errorf("protectiongroup does not exist: %s", n)
		}
		if stringInSlice(volume, p.Volumes) {
			if exists {
				return nil
			}
			return fmt.Errorf("volume %s still connected to Protection Group", volume)
		}
		if exists {
			return fmt.Errorf("volume %s not connected to Protection Group", volume)
		}
		return nil
	}
}

func testAccCheckPureProtectiongroupConfigBasic(rInt int) string {
	return fmt.Sprintf(`
resource "purestorage_protectiongroup" "tfprotectiongrouptest" {
        name = "tfprotectiongrouptest-%d"
}`, rInt)
}

func testAccCheckPureProtectiongroupConfigWithHosts(rInt int) string {
	return fmt.Sprintf(`
resource "purestorage_host" "tfpgrouptesthost" {
        name = "tfpgrouptesthost"
}

resource "purestorage_protectiongroup" "tfprotectiongrouptest" {
        name = "tfprotectiongrouptest-%d"
        hosts = ["${purestorage_host.tfpgrouptesthost.name}"]
}`, rInt)
}

func testAccCheckPureProtectiongroupConfigWithVolumes(rInt int) string {
	return fmt.Sprintf(`
resource "purestorage_volume" "tfpgrouptest-volume" {
	name = "tfpgrouptest-volume-%d"
	size = 1024000000
}

resource "purestorage_protectiongroup" "tfprotectiongrouptest" {
	name = "tfprotectiongrouptest-%d"
	volumes = ["${purestorage_volume.tfpgrouptest-volume.name}"]
}`, rInt, rInt)
}

func testAccCheckPureProtectiongroupConfigWithHostgroups(rInt int) string {
	return fmt.Sprintf(`
resource "purestorage_hostgroup" "tfpgrouptesthgroup" {
	name = "tfpgrouptesthgroup"
}

resource "purestorage_protectiongroup" "tfprotectiongrouptest" {
        name = "tfprotectiongrouptest-%d"
        hgroups = ["${purestorage_hostgroup.tfpgrouptesthgroup.name}"]
}`, rInt)
}

func testAccCheckPureProtectiongroupConfigWithSchedule(rInt int) string {
	return fmt.Sprintf(`
resource "purestorage_protectiongroup" "tfprotectiongrouptest" {
        name = "tfprotectiongrouptest-%d"
	replicate_enabled = "true"
	replicate_at = "3600"
	replicate_frequency = "86400"
	snap_enabled = "true"
	snap_at = "60"
	snap_frequency = "86400"
}`, rInt)
}

func testAccCheckPureProtectiongroupConfigWithRetention(rInt int) string {
	return fmt.Sprintf(`
resource "purestorage_protectiongroup" "tfprotectiongrouptest" {
	name = "tfprotectiongrouptest-%d"
	all_for = 86400
	days = 8
	per_day = 5
}`, rInt)
}
