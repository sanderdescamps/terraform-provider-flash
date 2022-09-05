package purestorage

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Test if provider can handle the creation and managment of a large amounts of resources.
func TestAccResourcePureLargeTest(t *testing.T) {
	testId := strconv.Itoa(acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeAggregateTestCheckFunc(
			testAccCheckPureHostDestroy,
			testAccCheckPureHostgroupDestroy,
			testAccCheckPureVolumeDestroy,
			testAccCheckPureVolumeGroupDestroy,
			testAccCheckPureVolumeEradicate,
			testAccCheckPureVolumeGroupEradicate,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPureHostConfigFullSetup(100, 100, testId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPureVolumeCount(testId, 100),
				),
			},
		},
	})
}

func testAccCheckPureHostConfigFullSetup(numberOfHosts int, numberOfVolumes int, testID string) string {
	output := ""
	output += fmt.Sprintf(`
		resource "purefa_vgroup" "tfhosttest-volumegroup" {
			name = "tfhosttest-volumegroup-%s"
		}
		`, testID)
	output += fmt.Sprintf(`
		resource "purefa_volume" "tfhosttest-volumes" {
			name = "tfhosttest-volume-%s-${count.index}"
			size = 1024000000
			volume_group = purefa_vgroup.tfhosttest-volumegroup.name
			count = %d
		}
		`, testID, numberOfVolumes)

	output += fmt.Sprintf(`
		resource "purefa_host" "tfhosttest" {
			name = "tfhosttest%s-${count.index}"
			wwn = ["0000999900009${format("%%03s", count.index)}"]
			count = %d
		}`, testID, numberOfHosts)

	output += fmt.Sprintf(`
		resource "purefa_hostgroup" "tfhostgrouptest" {
			name = "tfhosttest%s"
			dynamic "volume" {
				for_each = purefa_volume.tfhosttest-volumes
				content {
					vol = "${volume.value["full_name"]}"
					lun = volume.key + 1
				}
			  }
		}`, testID)
	return output
}
