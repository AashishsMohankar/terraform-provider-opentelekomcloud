package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/huaweicloud/golangsdk/openstack/cts/v1/tracker"
)

func TestAccCTSTrackerV1_basic(t *testing.T) {
	var tracker []tracker.Tracker

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCTSTrackerV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCTSTrackerV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCTSTrackerV1Exists("opentelekomcloud_cts_tracker_v1.tracker_v1", tracker),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_cts_tracker_v1.tracker_v1", "bucket_name", "tf-test-bucket"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_cts_tracker_v1.tracker_v1", "file_prefix_name", "yO8Q"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_cts_tracker_v1.tracker_v1", "status", "enabled"),
				),
			},
			resource.TestStep{
				Config: testAccCTSTrackerV1_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCTSTrackerV1Exists("opentelekomcloud_cts_tracker_v1.tracker_v1", tracker),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_cts_tracker_v1.tracker_v1", "bucket_name", "tf-test-bucket-update"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_cts_tracker_v1.tracker_v1", "is_support_smn", "false"),
				),
			},
		},
	})
}

func TestAccCTSTrackerV1_timeout(t *testing.T) {
	var tracker []tracker.Tracker

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCTSTrackerV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCTSTrackerV1_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCTSTrackerV1Exists("opentelekomcloud_cts_tracker_v1.tracker_v1", tracker),
				),
			},
		},
	})
}

func testAccCheckCTSTrackerV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	ctsClient, err := config.ctsV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating cts client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_cts_tracker_v1" {
			continue
		}

		_, err := tracker.Get(ctsClient).ExtractTracker()
		if err == nil {
			return fmt.Errorf("cts tracker still exists")
		}
	}

	return nil
}

func testAccCheckCTSTrackerV1Exists(n string, trackers []tracker.Tracker) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		ctsClient, err := config.ctsV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating cts client: %s", err)
		}

		found, err := tracker.Get(ctsClient).ExtractTracker()
		if err != nil {
			return err
		}

		if found[0].TrackerName != rs.Primary.ID {
			return fmt.Errorf("cts tracker not found")
		}

		trackers = found

		return nil
	}
}

var testAccCTSTrackerV1_basic = `
resource "opentelekomcloud_s3_bucket" "bucket" {
  bucket		= "tf-test-bucket"
  acl			= "public-read"
}
resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name			= "tf-test-topic"
  display_name	= "The display name of tf-test-topic"
}

resource "opentelekomcloud_cts_tracker_v1" "tracker_v1" {
  bucket_name		= "${opentelekomcloud_s3_bucket.bucket.bucket}"
  file_prefix_name  = "yO8Q"
  is_support_smn 	= true
  topic_id 			= "${opentelekomcloud_smn_topic_v2.topic_1.id}"
  is_send_all_key_operation = false
  operations 		= ["login"]
  need_notify_user_list = ["user1"]
}
`

var testAccCTSTrackerV1_update = `
resource "opentelekomcloud_s3_bucket" "bucket" {
  bucket = "tf-test-bucket-update"
  acl = "public-read"
}
resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name		  		= "tf-test-topic"
  display_name 		= "The display name of tf-test-topic"
}

resource "opentelekomcloud_cts_tracker_v1" "tracker_v1" {
  bucket_name      		= "${opentelekomcloud_s3_bucket.bucket.bucket}"
  file_prefix_name      = "yO8Q"
  is_support_smn 		= false
  topic_id 				= "${opentelekomcloud_smn_topic_v2.topic_1.id}"
  is_send_all_key_operation = false
  operations 			= ["delete","create","login"]
  need_notify_user_list = ["user1", "user2"]
}
`

var testAccCTSTrackerV1_timeout = `
resource "opentelekomcloud_s3_bucket" "bucket" {
  bucket 		= "tf-test-bucket"
  acl 			= "public-read"
}
resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name		  		= "tf-test-topic"
  display_name    	= "The display name of tf-test-topic"
}

resource "opentelekomcloud_cts_tracker_v1" "tracker_v1" {
  bucket_name      		= "${opentelekomcloud_s3_bucket.bucket.bucket}"
  file_prefix_name      = "yO8Q"
  is_support_smn 		= true
  topic_id 				= "${opentelekomcloud_smn_topic_v2.topic_1.id}"
  is_send_all_key_operation = false
  operations 			= ["login"]
  need_notify_user_list = ["user1"]
}
`
