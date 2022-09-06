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
	"strings"
	"testing"

	"github.com/devans10/pugo/flasharray"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testAccCheckPureVolumeResourceName = "purefa_volume.tfvolumetest"
const testAccCheckPureVolumeCloneResourceName = "purefa_volume.tfclonevolumetest"

// The volumes created in theses tests will not be eradicated.
//
// Create a volume
func TestAccResourcePureVolume_create(t *testing.T) {
	rInt := rand.Int()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPureVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPureVolumeConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureVolumeExists(testAccCheckPureVolumeResourceName, true),
					resource.TestCheckResourceAttr(testAccCheckPureVolumeResourceName, "name", fmt.Sprintf("tfvolumetest-%d", rInt)),
					resource.TestCheckResourceAttr(testAccCheckPureVolumeResourceName, "size", "1024000000"),
					resource.TestCheckResourceAttrSet(testAccCheckPureVolumeResourceName, "serial"),
				),
			},
		},
	})
}
func TestAccResourcePureVolume_clone(t *testing.T) {
	rInt := rand.Int()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPureVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPureVolumeConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureVolumeExists(testAccCheckPureVolumeResourceName, true),
					resource.TestCheckResourceAttr(testAccCheckPureVolumeResourceName, "name", fmt.Sprintf("tfvolumetest-%d", rInt)),
					resource.TestCheckResourceAttr(testAccCheckPureVolumeResourceName, "size", "1024000000"),
					resource.TestCheckResourceAttrSet(testAccCheckPureVolumeResourceName, "serial"),
				),
			},
			{
				Config: testAccCheckPureVolumeConfigClone(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureVolumeExists(testAccCheckPureVolumeCloneResourceName, true),
					resource.TestCheckResourceAttr(testAccCheckPureVolumeCloneResourceName, "source", fmt.Sprintf("tfvolumetest-%d", rInt)),
				),
			},
		},
	})
}

func TestAccResourcePureVolume_update(t *testing.T) {
	rInt := rand.Int()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPureVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPureVolumeConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureVolumeExists(testAccCheckPureVolumeResourceName, true),
					resource.TestCheckResourceAttr(testAccCheckPureVolumeResourceName, "name", fmt.Sprintf("tfvolumetest-%d", rInt)),
					resource.TestCheckResourceAttr(testAccCheckPureVolumeResourceName, "size", "1024000000"),
					resource.TestCheckResourceAttrSet(testAccCheckPureVolumeResourceName, "serial"),
				),
			},
			{
				Config: testAccCheckPureVolumeConfigResize(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureVolumeExists(testAccCheckPureVolumeResourceName, true),
					resource.TestCheckResourceAttr(testAccCheckPureVolumeResourceName, "name", fmt.Sprintf("tfvolumetest-%d", rInt)),
					resource.TestCheckResourceAttr(testAccCheckPureVolumeResourceName, "size", "2048000000"),
					resource.TestCheckResourceAttrSet(testAccCheckPureVolumeResourceName, "serial"),
				),
			},
			{
				Config: testAccCheckPureVolumeConfigRename(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureVolumeExists(testAccCheckPureVolumeResourceName, true),
					resource.TestCheckResourceAttr(testAccCheckPureVolumeResourceName, "name", fmt.Sprintf("tfvolumetest-rename-%d", rInt)),
					resource.TestCheckResourceAttr(testAccCheckPureVolumeResourceName, "size", "2048000000"),
					resource.TestCheckResourceAttrSet(testAccCheckPureVolumeResourceName, "serial"),
				),
			},
		},
	})
}

func testAccCheckPureVolumeDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*flasharray.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "purefa_volume" {
			continue
		}

		_, err := client.Volumes.GetVolume(rs.Primary.ID, nil)
		if err != nil {
			return nil
		}
		return fmt.Errorf("volume '%s' stil exists", rs.Primary.ID)
	}

	return nil
}

// The pugo sdk does not support the 'pending_only' parameter. This is a temporary fix.
type Volume struct {
	Name    string `json:"name,omitempty"`
	Source  string `json:"source,omitempty"`
	Serial  string `json:"serial,omitempty"`
	Size    int    `json:"size,omitempty"`
	Created string `json:"created,omitempty"`

	// response returned with the pending_only=true flag
	TimeRemaining            *int    `json:"time_remaining,omitempty"`
	PromotionStatus          *string `json:"promotion_status,omitempty"`
	RequestedPromotionStatus *string `json:"requested_promotion_status,omitempty"`
}

// Checks if resources are still pending for deletion and eredicates if needed
// The pugo sdk does not support listing deleted volumegroups. This method includes a temporary fix.
func testAccCheckPureVolumeEradicate(s *terraform.State) error {
	client := testAccProvider.Meta().(*flasharray.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "purefa_volume" {
			continue
		}

		params := map[string]string{"pending_only": "true"}
		req, _ := client.NewRequest("GET", fmt.Sprintf("volume/%s", rs.Primary.ID), params, nil)
		volume := &Volume{}
		_, err := client.Do(req, &volume, false)
		if err != nil {
			return nil
		} else if volume != nil && volume.TimeRemaining != nil && *volume.TimeRemaining > 0 {
			_, err := client.Volumes.EradicateVolume(volume.Name)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func testAccCheckPureVolumeExists(n string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*flasharray.Client)
		_, err := client.Volumes.GetVolume(rs.Primary.ID, nil)
		if err != nil {
			if exists {
				return fmt.Errorf("volume does not exist: %s", n)
			}
			return nil
		}
		return nil
	}
}

func testAccCheckPureVolumeCount(testID string, count int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rsCount := 0
		volCount := 0
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "purefa_volume" {
				continue
			}

			if name, ok := rs.Primary.Attributes["name"]; ok {
				if strings.Contains(name, testID) {
					rsCount += 1
				}
			}
		}

		client := testAccProvider.Meta().(*flasharray.Client)
		if volumes, err := client.Volumes.ListVolumes(nil); err == nil {
			for _, volume := range volumes {
				if strings.Contains(volume.Name, testID) {
					volCount += 1
				}
			}
		}

		if rsCount > count {
			return fmt.Errorf("Too many volumes in state file (has %d expected %d)", rsCount, count)
		} else if rsCount < count {
			return fmt.Errorf("Missing volumes in state file (has %d expected %d)", rsCount, count)
		} else if volCount > count {
			return fmt.Errorf("Too many volumes on storage array (has %d expected %d)", volCount, count)
		} else if volCount < count {
			return fmt.Errorf("Missing volumes on storage array (has %d expected %d)", volCount, count)
		}
		return nil
	}
}

func testAccCheckPureVolumeConfig(rInt int) string {
	return fmt.Sprintf(`
resource "purefa_volume" "tfvolumetest" {
        name = "tfvolumetest-%d"
        size = 1024000000
		allow_destroy = true
}`, rInt)
}

func testAccCheckPureVolumeConfigClone(rInt int) string {
	return fmt.Sprintf(`
resource "purefa_volume" "tfvolumetest" {
        name = "tfvolumetest-%d"
        size = 1024000000
		allow_destroy = true
}

resource "purefa_volume" "tfclonevolumetest" {
        name = "tfclonevolumetest-%d"
        source = "${purefa_volume.tfvolumetest.name}"
		allow_destroy = true
}`, rInt, rInt)
}

func testAccCheckPureVolumeConfigResize(rInt int) string {
	return fmt.Sprintf(`
resource "purefa_volume" "tfvolumetest" {
	name = "tfvolumetest-%d"
	size = 2048000000
	allow_destroy = true
}`, rInt)
}

func testAccCheckPureVolumeConfigRename(rInt int) string {
	return fmt.Sprintf(`
resource "purefa_volume" "tfvolumetest" {
        name = "tfvolumetest-rename-%d"
        size = 2048000000
		allow_destroy = true
}`, rInt)
}
