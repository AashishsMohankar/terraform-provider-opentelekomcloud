package opentelekomcloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/peerings"
	"log"
	"time"
)

func resourceVpcPeeringConnectionV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceVPCPeeringV2Create, //providers.go
		Read:   resourceVPCPeeringV2Read,
		Update: resourceVPCPeeringV2Update,
		Delete: resourceVPCPeeringV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{ //request and response parameters
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validateName,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: false,
				Computed: true,
			},
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"peer_vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"peer_tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceVPCPeeringV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	peeringClient, err := config.networkingHwV2Client(GetRegion(d, config))

	log.Printf("[DEBUG] Value of peeringClient: %#v", peeringClient)

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud Vpc Peering Connection Client: %s", err)
	}

	requestvpcinfo := peerings.VpcInfo{
		VpcId: d.Get("vpc_id").(string),
	}

	acceptvpcinfo := peerings.VpcInfo{
		VpcId:    d.Get("peer_vpc_id").(string),
		TenantId: d.Get("peer_tenant_id").(string),
	}

	createOpts := peerings.CreateOpts{
		Name:           d.Get("name").(string),
		RequestVpcInfo: requestvpcinfo,
		AcceptVpcInfo:  acceptvpcinfo,
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	n, err := peerings.Create(peeringClient, createOpts).Extract()

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud Vpc Peering Connection: %s", err)
	}
	d.SetId(n.ID)

	log.Printf("[INFO] Vpc Peering Connection ID: %s", n.ID)

	log.Printf("[DEBUG] Waiting for OpenTelekomCloud Vpc Peering Connection(%s) to become available", n.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATING"},
		Target:     []string{"PENDING_ACCEPTANCE", "ACTIVE"},
		Refresh:    waitForVpcPeeringActive(peeringClient, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	d.SetId(n.ID)

	return resourceVPCPeeringV2Read(d, meta)

}

func resourceVPCPeeringV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	peeringClient, err := config.networkingHwV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud   Vpc Peering Connection Client: %s", err)
	}

	n, err := peerings.Get(peeringClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving OpenTelekomCloud Vpc Peering Connection: %s", err)
	}

	log.Printf("[DEBUG] Retrieved Vpc Peering Connection %s: %+v", d.Id(), n)

	d.Set("id", n.ID)
	d.Set("name", n.Name)
	d.Set("status", n.Status)
	d.Set("vpc_id", n.RequestVpcInfo.VpcId)
	d.Set("peer_vpc_id", n.AcceptVpcInfo.VpcId)
	d.Set("peer_tenant_id", n.AcceptVpcInfo.TenantId)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceVPCPeeringV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	peeringClient, err := config.networkingHwV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud  Vpc Peering Connection Client: %s", err)
	}

	var update bool
	var updateOpts peerings.UpdateOpts

	if d.HasChange("name") {
		update = true
		updateOpts.Name = d.Get("name").(string)
	}

	log.Printf("[DEBUG] Updating Vpc Peering Connection %s with options: %+v", d.Id(), updateOpts)

	if update {
		log.Printf("[DEBUG] Updating Vpc Peering Connection %s with options: %#v", d.Id(), updateOpts)
		_, err = peerings.Update(peeringClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating OpenTelekomCloud Vpc Peering Connection: %s", err)
		}
	}

	return resourceVPCPeeringV2Read(d, meta)
}

func resourceVPCPeeringV2Delete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Destroy vpc peering connection: %s", d.Id())

	config := meta.(*Config)
	peeringClient, err := config.networkingHwV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud  Vpc Peering Connection Client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForVpcPeeringDelete(peeringClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting OpenTelekomCloud Vpc Peering Connection: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForVpcPeeringActive(peeringClient *golangsdk.ServiceClient, peeringId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := peerings.Get(peeringClient, peeringId).Extract()
		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] OpenTelekomCloud Peering Client: %+v", n)
		if n.Status == "PENDING_ACCEPTANCE" || n.Status == "ACTIVE" {
			return n, n.Status, nil
		}

		return n, "CREATING", nil
	}
}

func waitForVpcPeeringDelete(peeringClient *golangsdk.ServiceClient, peeringId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete OpenTelekomCloud vpc peering connection %s.\n", peeringId)

		r, err := peerings.Get(peeringClient, peeringId).Extract()
		log.Printf("[DEBUG] Value after extract: %#v", r)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted OpenTelekomCloud vpc peering connection %s", peeringId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		err = peerings.Delete(peeringClient, peeringId).ExtractErr()
		log.Printf("[DEBUG] Value if error: %#v", err)

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted OpenTelekomCloud vpc peering connection %s", peeringId)
				return r, "DELETED", nil
			}
			if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
				if errCode.Actual == 409 {
					return r, "ACTIVE", nil
				}
			}
			return r, "ACTIVE", err
		}

		log.Printf("[DEBUG] OpenTelekomCloud vpc peering connection %s still active.\n", peeringId)
		return r, "ACTIVE", nil
	}
}
