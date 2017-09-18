package huaweicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

// DEPRECATED
func TestAccBlockStorageV1Volume_importBasic(t *testing.T) {
	resourceName := "huaweicloud_blockstorage_volume_v1.volume_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBlockStorageV1VolumeDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccBlockStorageV1Volume_basic,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
