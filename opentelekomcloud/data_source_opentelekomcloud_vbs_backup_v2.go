package opentelekomcloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/vbs/v2/backups"
)

func dataSourceVBSBackupV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVBSBackupV2Read,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"fail_reason": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"object_count": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"container": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_metadata": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceVBSBackupV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vbsClient, err := config.vbsV2Client(GetRegion(d, config))

	listOpts := backups.ListOpts{
		Id:         d.Get("id").(string),
		Name:       d.Get("name").(string),
		Status:     d.Get("status").(string),
		VolumeId:   d.Get("volume_id").(string),
		SnapshotId: d.Get("snapshot_id").(string),
	}

	refinedBackups, err := backups.List(vbsClient, listOpts)
	if err != nil {
		return fmt.Errorf("Unable to retrieve backups: %s", err)
	}

	if len(refinedBackups) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedBackups) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	Backup := refinedBackups[0]

	log.Printf("[INFO] Retrieved Backup using given filter %s: %+v", Backup.Id, Backup)
	d.SetId(Backup.Id)

	d.Set("name", Backup.Name)
	d.Set("description", Backup.Description)
	d.Set("status", Backup.Status)
	d.Set("availability_zone", Backup.AZ)
	d.Set("snapshot_id", Backup.SnapshotId)
	d.Set("service_metadata", Backup.ServiceMetadata)
	d.Set("size", Backup.Size)
	d.Set("fail_reason", Backup.FailReason)
	d.Set("container", Backup.Container)
	d.Set("tenant_id", Backup.TenantId)
	d.Set("object_count", Backup.ObjectCount)
	d.Set("volume_id", Backup.VolumeId)
	d.Set("region", GetRegion(d, config))

	return nil
}