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

const testAccCheckPureHostResourceName = "purestorage_host.tfhosttest"

// Create a host
func TestAccResourcePureHost_create(t *testing.T) {
	rInt := rand.Int()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPureHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPureHostConfigBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureHostExists(testAccCheckPureHostResourceName, true),
					resource.TestCheckResourceAttr(testAccCheckPureHostResourceName, "name", fmt.Sprintf("tfhosttest%d", rInt)),
				),
			},
		},
	})
}

func TestAccResourcePureHost_createWithWWN(t *testing.T) {
	rInt := rand.Int()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPureHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPureHostConfigWithWWN(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureHostExists(testAccCheckPureHostResourceName, true),
					resource.TestCheckResourceAttr(testAccCheckPureHostResourceName, "name", fmt.Sprintf("tfhosttest%d", rInt)),
					testAccCheckPureHostWWN(testAccCheckPureHostResourceName, "0000999900009999", true),
				),
			},
		},
	})
}

func TestAccResourcePureHost_createWithVolume(t *testing.T) {
	rInt := rand.Int()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPureHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPureHostConfigWithVolume(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureHostExists(testAccCheckPureHostResourceName, true),
					resource.TestCheckResourceAttr(testAccCheckPureHostResourceName, "name", fmt.Sprintf("tfhosttest%d", rInt)),
					testAccCheckPureHostWWN(testAccCheckPureHostResourceName, "0000999900009999", true),
					testAccCheckPureHostVolumeConnection(testAccCheckPureHostResourceName, fmt.Sprintf("tfhosttest-volume-%d", rInt), true),
				),
			},
		},
	})
}

/*
func TestAccResourcePureHost_createWithCHAP(t *testing.T) {
	rInt := rand.Int()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPureHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPureHostConfig_withCHAP(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureHostExists(testAccCheckPureHostResourceName, true),
					resource.TestCheckResourceAttr(testAccCheckPureHostResourceName, "name", fmt.Sprintf("tfhosttest%d", rInt)),
					resource.TestCheckResourceAttr(testAccCheckPureHostResourceName, "host_user", "myhostuser"),
					resource.TestCheckResourceAttr(testAccCheckPureHostResourceName, "target_user", "mytargetuser"),
					testAccCheckPureHostCHAP(testAccCheckPureHostResourceName, "host_user", "myhostuser", true),
					testAccCheckPureHostCHAP(testAccCheckPureHostResourceName, "target_user", "mytargetuser", true),
				),
			},
		},
	})
}
*/

func TestAccResourcePureHost_createWithPrivateAndSharedVolumes(t *testing.T) {
	rInt := rand.Int()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPureHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPureHostConfigWithPrivateAndSharedVolumes(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureHostExists(testAccCheckPureHostResourceName, true),
					resource.TestCheckResourceAttr(testAccCheckPureHostResourceName, "name", fmt.Sprintf("tfhosttest%d", rInt)),
				),
			},
		},
	})
}

func TestAccResourcePureHost_createWithPersonality(t *testing.T) {
	rInt := rand.Int()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPureHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPureHostConfigWithPersonality(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureHostExists(testAccCheckPureHostResourceName, true),
					resource.TestCheckResourceAttr(testAccCheckPureHostResourceName, "name", fmt.Sprintf("tfhosttest%d", rInt)),
					resource.TestCheckResourceAttr(testAccCheckPureHostResourceName, "personality", "aix"),
					testAccCheckPureHostPersonality(testAccCheckPureHostResourceName, "aix", true),
				),
			},
		},
	})
}

func TestAccResourcePureHost_update(t *testing.T) {
	rInt := rand.Int()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPureHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPureHostConfigBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureHostExists(testAccCheckPureHostResourceName, true),
					resource.TestCheckResourceAttr(testAccCheckPureHostResourceName, "name", fmt.Sprintf("tfhosttest%d", rInt)),
				),
			},
			{
				Config: testAccCheckPureHostConfigWithWWN(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureHostWWN(testAccCheckPureHostResourceName, "0000999900009999", true),
				),
			},
			{
				Config: testAccCheckPureHostConfigRename(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCheckPureHostResourceName, "name", fmt.Sprintf("tfhosttestrename%d", rInt)),
				),
			},
		},
	})
}

func TestAccResourcePureHost_update_AddandRemoveVolume(t *testing.T) {
	rInt := rand.Int()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPureHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPureHostConfigBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureHostExists(testAccCheckPureHostResourceName, true),
					resource.TestCheckResourceAttr(testAccCheckPureHostResourceName, "name", fmt.Sprintf("tfhosttest%d", rInt)),
				),
			},
			{
				Config: testAccCheckPureHostConfigWithVolume(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureHostVolumeConnection(testAccCheckPureHostResourceName, fmt.Sprintf("tfhosttest-volume-%d", rInt), true),
				),
			},
			{
				Config: testAccCheckPureHostConfigWithoutVolume(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureHostVolumeConnection(testAccCheckPureHostResourceName, fmt.Sprintf("tfhosttest-volume-%d", rInt), false),
				),
			},
		},
	})
}

/*
func TestAccResourcePureHost_update_withCHAP(t *testing.T) {
	rInt := rand.Int()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPureHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPureHostConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureHostExists(testAccCheckPureHostResourceName, true),
					resource.TestCheckResourceAttr(testAccCheckPureHostResourceName, "name", fmt.Sprintf("tfhosttest%d", rInt)),
				),
			},
			{
				Config: testAccCheckPureHostConfig_withCHAP(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureHostExists(testAccCheckPureHostResourceName, true),
					resource.TestCheckResourceAttr(testAccCheckPureHostResourceName, "name", fmt.Sprintf("tfhosttest%d", rInt)),
					resource.TestCheckResourceAttr(testAccCheckPureHostResourceName, "host_user", "myhostuser"),
					resource.TestCheckResourceAttr(testAccCheckPureHostResourceName, "target_user", "mytargetuser"),
					testAccCheckPureHostCHAP(testAccCheckPureHostResourceName, "host_user", "myhostuser", true),
					testAccCheckPureHostCHAP(testAccCheckPureHostResourceName, "target_user", "mytargetuser", true),
				),
			},
		},
	})
}
*/

func TestAccResourcePureHost_update_withPersonality(t *testing.T) {
	rInt := rand.Int()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPureHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPureHostConfigBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureHostExists(testAccCheckPureHostResourceName, true),
					resource.TestCheckResourceAttr(testAccCheckPureHostResourceName, "name", fmt.Sprintf("tfhosttest%d", rInt)),
				),
			},
			{
				Config: testAccCheckPureHostConfigWithPersonality(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureHostExists(testAccCheckPureHostResourceName, true),
					resource.TestCheckResourceAttr(testAccCheckPureHostResourceName, "name", fmt.Sprintf("tfhosttest%d", rInt)),
					resource.TestCheckResourceAttr(testAccCheckPureHostResourceName, "personality", "aix"),
					testAccCheckPureHostPersonality(testAccCheckPureHostResourceName, "aix", true),
				),
			},
		},
	})
}

func testAccCheckPureHostDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*flasharray.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "purestorage_host" {
			continue
		}

		_, err := client.Hosts.GetHost(rs.Primary.ID, nil)
		if err != nil {
			return nil
		}
		return fmt.Errorf("host '%s' stil exists", rs.Primary.ID)
	}

	return nil
}

func testAccCheckPureHostExists(n string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*flasharray.Client)
		name, ok := rs.Primary.Attributes["name"]
		_, err := client.Hosts.GetHost(name, nil)
		if err != nil {
			if exists {
				return fmt.Errorf("host does not exist: %s", n)
			}
			return nil
		}
		return nil
	}
}

func testAccCheckPureHostWWN(n string, wwn string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*flasharray.Client)
		name, ok := rs.Primary.Attributes["name"]
		h, err := client.Hosts.GetHost(name, nil)
		if err != nil {
			return fmt.Errorf("host does not exist: %s", n)
		}
		if stringInSlice(wwn, h.Wwn) {
			if exists {
				return nil
			}
			return fmt.Errorf("wwn %s is still set for host", wwn)
		}
		if exists {
			return fmt.Errorf("wwn %s not set for host", wwn)
		}
		return nil
	}
}

func testAccCheckPureHostVolumeConnection(n string, volume string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*flasharray.Client)
		name, ok := rs.Primary.Attributes["name"]
		volumes, err := client.Hosts.ListHostConnections(name, map[string]string{"private": "true"})
		if err != nil {
			return fmt.Errorf("host does not exist: %s", n)
		}
		for _, elem := range volumes {
			if elem.Vol == volume {
				if exists {
					return nil
				}
				return fmt.Errorf("volume %s still connected to host", volume)
			}
		}
		if exists {
			return fmt.Errorf("volume %s not connected to host", volume)
		}
		return nil
	}
}

func testAccCheckPureHostCHAP(n string, param string, value string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*flasharray.Client)
		name, ok := rs.Primary.Attributes["name"]
		h, err := client.Hosts.GetHost(name, map[string]string{"chap": "true"})
		if err != nil {
			return fmt.Errorf("host does not exist: %s", n)
		}

		switch param {
		case "host_password":
			if h.HostPassword == value {
				if exists {
					return nil
				}
				return fmt.Errorf("%s is still set for host", param)
			}
			if exists {
				return fmt.Errorf("%s not set for host", param)
			}
			return nil
		case "host_user":
			if h.HostUser == value {
				if exists {
					return nil
				}
				return fmt.Errorf("%s is still set for host", param)
			}
			if exists {
				return fmt.Errorf("%s not set for host", param)
			}
			return nil
		case "target_password":
			if h.TargetPassword == value {
				if exists {
					return nil
				}
				return fmt.Errorf("%s is still set for host", param)
			}
			if exists {
				return fmt.Errorf("%s not set for host", param)
			}
			return nil
		case "target_user":
			if h.TargetUser == value {
				if exists {
					return nil
				}
				return fmt.Errorf("%s is still set for host", param)
			}
			if exists {
				return fmt.Errorf("%s not set for host", param)
			}
			return nil
		default:
			return fmt.Errorf("%s is not a valid CHAP parameter", param)
		}
	}
}

func testAccCheckPureHostPersonality(n string, personality string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*flasharray.Client)
		name, ok := rs.Primary.Attributes["name"]
		h, err := client.Hosts.GetHost(name, map[string]string{"personality": "true"})
		if err != nil {
			return fmt.Errorf("host does not exist: %s", n)
		}
		if personality == h.Personality {
			if exists {
				return nil
			}
			return fmt.Errorf("personality %s is still set for host", personality)
		}
		if exists {
			return fmt.Errorf("personality %s not set for host", personality)
		}
		return nil
	}
}

func testAccCheckPureHostConfigBasic(rInt int) string {
	return fmt.Sprintf(`
resource "purestorage_host" "tfhosttest" {
        name = "tfhosttest%d"
}`, rInt)
}

func testAccCheckPureHostConfigRename(rInt int) string {
	return fmt.Sprintf(`
resource "purestorage_host" "tfhosttest" {
        name = "tfhosttestrename%d"
	wwn = ["0000999900009999"]
}`, rInt)
}

func testAccCheckPureHostConfigWithWWN(rInt int) string {
	return fmt.Sprintf(`
resource "purestorage_host" "tfhosttest" {
        name = "tfhosttest%d"
	wwn = ["0000999900009999"]
}`, rInt)
}

func testAccCheckPureHostConfigWithVolume(rInt int) string {
	return fmt.Sprintf(`
resource "purestorage_volume" "tfhosttest-volume" {
	name = "tfhosttest-volume-%d"
	size = 1024000000
}
resource "purestorage_host" "tfhosttest" {
        name = "tfhosttest%d"
        wwn = ["0000999900009999"]
	volume {
		vol = "${purestorage_volume.tfhosttest-volume.name}"
		lun = 1
	}
}`, rInt, rInt)
}

func testAccCheckPureHostConfigWithoutVolume(rInt int) string {
	return fmt.Sprintf(`
resource "purestorage_volume" "tfhosttest-volume" {
        name = "tfhosttest-volume-%d"
        size = 1024000000
}
resource "purestorage_host" "tfhosttest" {
        name = "tfhosttest%d"
        wwn = ["0000999900009999"]
}`, rInt, rInt)
}

func testAccCheckPureHostConfigWithCHAP(rInt int) string {
	return fmt.Sprintf(`
resource "purestorage_host" "tfhosttest" {
	name = "tfhosttest%d"
	host_user = "myhostuser"
	target_user = "mytargetuser"
}`, rInt)
}

func testAccCheckPureHostConfigWithPersonality(rInt int) string {
	return fmt.Sprintf(`
resource "purestorage_host" "tfhosttest" {
        name = "tfhosttest%d"
	personality = "aix"
}`, rInt)
}

func testAccCheckPureHostConfigWithPrivateAndSharedVolumes(rInt int) string {
	return fmt.Sprintf(`
resource "purestorage_volume" "tfhosttest-private-volume" {
	name = "tfhosttest-private-volume-%d"
	size = 1024000000
}

resource "purestorage_volume" "tfhosttest-shared-volume" {
        name = "tfhosttest-shared-volume-%d"
        size = 1024000000
}

resource "purestorage_host" "tfhosttest" {
	name = "tfhosttest%d"
	volume {
		vol = "${purestorage_volume.tfhosttest-private-volume.name}"
		lun = 1
	}
}

resource "purestorage_hostgroup" "tfhosttesthostgroup" {
	name = "tfhosthostgroup%d"
	hosts = ["${purestorage_host.tfhosttest.name}"]
	volume {
		vol = "${purestorage_volume.tfhosttest-shared-volume.name}"
		lun = 250
	}
}`, rInt, rInt, rInt, rInt)
}
