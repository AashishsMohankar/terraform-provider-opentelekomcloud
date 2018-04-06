package opentelekomcloud

import (
	"fmt"
	"log"

	"github.com/huaweicloud/golangsdk/openstack/networking/v2/peerings"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceVpcPeeringConnectionV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVpcPeeringConnectionV2Read,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateName,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"peer_vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"peer_tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceVpcPeeringConnectionV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	peeringClient, err := config.networkingHwV2Client(GetRegion(d, config))

	listOpts := peerings.ListOpts{
		ID:         d.Get("id").(string),
		Name:       d.Get("name").(string),
		Status:     d.Get("status").(string),
		TenantId:   d.Get("peer_tenant_id").(string),
		VpcId:      d.Get("vpc_id").(string),
		Peer_VpcId: d.Get("peer_vpc_id").(string),
	}

	refinedPeering, err := peerings.List(peeringClient, listOpts)
	log.Printf("[DEBUG] Value of allPeeringConnections: %#v", refinedPeering)
	if err != nil {
		return fmt.Errorf("Unable to retrieve vpc peering connections: %s", err)
	}

	if len(refinedPeering) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedPeering) > 1 {
		return fmt.Errorf("Multiple VPC peering connections matched." +
			" Use additional constraints to reduce matches to a single VPC peering connection")
	}

	Peering := refinedPeering[0]

	log.Printf("[DEBUG] Retrieved Vpc peering Connections using given filter %s: %+v", Peering.ID, Peering)
	d.SetId(Peering.ID)

	d.Set("id", Peering.ID)
	d.Set("name", Peering.Name)
	d.Set("status", Peering.Status)
	d.Set("vpc_id", Peering.RequestVpcInfo.VpcId)
	d.Set("peer_vpc_id", Peering.AcceptVpcInfo.VpcId)
	d.Set("peer_tenant_id", Peering.AcceptVpcInfo.TenantId)
	d.Set("region", GetRegion(d, config))

	return nil
}
