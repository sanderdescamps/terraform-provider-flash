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
	"strconv"
	"strings"
	"testing"

	"github.com/devans10/pugo/flasharray"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testAccCheckPureVolumeGroupResourceName = "purefa_volumegroup.tfvolumegrouptest"

// The volumes created in theses tests will not be eradicated.
//
// Create a volume
func TestAccResourcePureVolumeGroup_create(t *testing.T) {
	testID := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPureVolumeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPureVolumeGroupConfig("tfvolumegrouptest", testID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureVolumeGroupExists(testAccCheckPureVolumeGroupResourceName, true),
					resource.TestCheckResourceAttr(testAccCheckPureVolumeGroupResourceName, "name", fmt.Sprintf("tfvolumegrouptest-%d", testID)),
				),
			},
		},
	})
}

func TestAccResourcePureVolumeGroup_update(t *testing.T) {
	testID := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPureVolumeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPureVolumeGroupConfig("tfvolumegrouptest", testID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureVolumeGroupExists(testAccCheckPureVolumeGroupResourceName, true),
					resource.TestCheckResourceAttr(testAccCheckPureVolumeGroupResourceName, "name", fmt.Sprintf("tfvolumegrouptest-%d", testID)),
				),
			},
			{
				Config: testAccCheckPureVolumeGroupConfig("tfvolumegrouptest-rename", testID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureVolumeGroupExists(testAccCheckPureVolumeGroupResourceName, true),
					resource.TestCheckResourceAttr(testAccCheckPureVolumeGroupResourceName, "name", fmt.Sprintf("tfvolumegrouptest-rename-%d", testID)),
				),
			},
		},
	})
}

func TestAccResourcePureVolumeGroup_withVolumes(t *testing.T) {
	testID := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPureVolumeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPureVolumeGroupConfigWithoutVolumes("tfvolumegrouptest", "tfvolumetest", 3, testID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureVolumeGroupExists(testAccCheckPureVolumeGroupResourceName, true),
					testAccCheckPureVolumeCount(strconv.Itoa(testID), 3),
				),
			},
			{
				Config: testAccCheckPureVolumeGroupConfigWithVolumes("tfvolumegrouptest", "tfvolumetest", 3, testID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureVolumeGroupExists(testAccCheckPureVolumeGroupResourceName, true),
					testAccCheckPureVolumeCount(strconv.Itoa(testID), 3),
				),
			},
			{
				Config: testAccCheckPureVolumeGroupConfigMoveVolumes("tfvolumegrouptest", "tfvolumetest", 3, testID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureVolumeGroupCount(strconv.Itoa(testID), 2),
					testAccCheckPureVolumeCount(strconv.Itoa(testID), 3),
				),
			},
		},
	})
}

func testAccCheckPureVolumeGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*flasharray.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "purefa_volumegroup" {
			continue
		}

		_, err := client.Vgroups.GetVgroup(rs.Primary.ID)
		if err != nil {
			return nil
		}
		return fmt.Errorf("volume '%s' stil exists", rs.Primary.ID)
	}

	return nil
}

// The pugo sdk does not support the 'pending_only' parameter. This is a temporary fix.
type Vgroup struct {
	Name    string   `json:"name"`
	Volumes []string `json:"volumes"`

	// response returned with the pending_only=true flag
	TimeRemaining *int `json:"time_remaining,omitempty"`
}

// Checks if resources are still pending for deletion and eredicates if needed
// The pugo sdk does not support listing deleted volumegroups. This method includes a temporary fix.
func testAccCheckPureVolumeGroupEradicate(s *terraform.State) error {
	client := testAccProvider.Meta().(*flasharray.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "purefa_volumegroup" {
			continue
		}

		params := map[string]string{"pending_only": "true"}
		req, _ := client.NewRequest("GET", fmt.Sprintf("vgroup/%s", rs.Primary.ID), params, nil)
		vgroup := &Vgroup{}
		_, err := client.Do(req, &vgroup, false)
		if err != nil {
			return nil
		} else if vgroup != nil && vgroup.TimeRemaining != nil && *vgroup.TimeRemaining > 0 {
			_, err := client.Vgroups.EradicateVgroup(vgroup.Name)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func testAccCheckPureVolumeGroupExists(n string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*flasharray.Client)
		_, err := client.Vgroups.GetVgroup(rs.Primary.ID)
		if err != nil {
			if exists {
				return fmt.Errorf("volume group does not exist: %s", n)
			}
			return nil
		}
		return nil
	}
}

func testAccCheckPureVolumeGroupCount(testID string, count int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rsCount := 0
		vgCount := 0
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "purefa_volumegroup" {
				continue
			}

			if name, ok := rs.Primary.Attributes["name"]; ok {
				if strings.Contains(name, testID) {
					rsCount += 1
				}
			}
		}

		client := testAccProvider.Meta().(*flasharray.Client)
		if vgroups, err := client.Vgroups.ListVgroups(); err == nil {
			for _, vgroup := range vgroups {
				if strings.Contains(vgroup.Name, testID) {
					vgCount += 1
				}
			}
		}

		if rsCount > count {
			return fmt.Errorf("Too many vgroups in state file (has %d expected %d)", rsCount, count)
		} else if rsCount < count {
			return fmt.Errorf("Missing vgroups in state file (has %d expected %d)", rsCount, count)
		} else if vgCount > count {
			return fmt.Errorf("Too many vgroups on storage array (has %d expected %d)", vgCount, count)
		} else if vgCount < count {
			return fmt.Errorf("Missing vgroups on storage array (has %d expected %d)", vgCount, count)
		}
		return nil
	}
}

func testAccCheckPureVolumeGroupConfig(vgroupName string, testID int) string {
	return fmt.Sprintf(`
resource "purefa_volumegroup" "tfvolumegrouptest" {
        name = "%s-%d"
}`, vgroupName, testID)
}

func testAccCheckPureVolumeGroupConfigWithVolumes(vgroupName string, volumeName string, numerOfVolumes int, testID int) string {
	output := ""
	output += fmt.Sprintf(`
		resource "purefa_volumegroup" "tfvolumegrouptest" {
				name = "%s-%d"
		}`, vgroupName, testID)
	output += fmt.Sprintf(`
		resource "purefa_volume" "tfvolumetest" {
			name = "%s-%d-${count.index}"
			size = 1024000000
			volume_group = purefa_volumegroup.tfvolumegrouptest.name
			count = %d
		}`, volumeName, testID, numerOfVolumes)
	return output
}

func testAccCheckPureVolumeGroupConfigWithoutVolumes(vgroupName string, volumeName string, numerOfVolumes int, testID int) string {
	output := ""
	output += fmt.Sprintf(`
		resource "purefa_volumegroup" "tfvolumegrouptest" {
				name = "%s-%d"
		}`, vgroupName, testID)
	output += fmt.Sprintf(`
		resource "purefa_volume" "tfvolumetest" {
			name = "%s-%d-${count.index}"
			size = 1024000000
			count = %d
		}`, volumeName, testID, numerOfVolumes)
	return output
}

func testAccCheckPureVolumeGroupConfigMoveVolumes(vgroupName string, volumeName string, numerOfVolumes int, testID int) string {
	output := ""
	output += fmt.Sprintf(`
		resource "purefa_volumegroup" "tfvolumegrouptest" {
				name = "%s-%d"
		}`, vgroupName, testID)
	output += fmt.Sprintf(`
		resource "purefa_volumegroup" "tfvolumegrouptest2" {
				name = "%s-%d-2"
		}`, vgroupName, testID)
	output += fmt.Sprintf(`
		resource "purefa_volume" "tfvolumetest" {
			name = "%s-%d-${count.index}"
			size = 1024000000
			volume_group = purefa_volumegroup.tfvolumegrouptest2.name
			count = %d
		}`, volumeName, testID, numerOfVolumes)
	return output
}
