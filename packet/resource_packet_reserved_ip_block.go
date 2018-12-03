package packet

import (
	"fmt"
	"path"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/packethost/packngo"
)

func packetIPComputedFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"address": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"gateway": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"network": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"netmask": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"cidr": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"address_family": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"public": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"management": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"manageable": {
			Type:     schema.TypeBool,
			Computed: true,
		},
	}
}

func resourcePacketReservedIPBlock() *schema.Resource {
	reservedBlockSchema := packetIPComputedFields()
	reservedBlockSchema["project_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	}
	reservedBlockSchema["facility"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	}
	reservedBlockSchema["quantity"] = &schema.Schema{
		Type:     schema.TypeInt,
		Required: true,
		ForceNew: true,
	}
	reservedBlockSchema["cidr_notation"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	return &schema.Resource{
		Create: resourcePacketReservedIPBlockCreate,
		Read:   resourcePacketReservedIPBlockRead,
		Delete: resourcePacketReservedIPBlockDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: reservedBlockSchema,
	}
}

func resourcePacketReservedIPBlockCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	facility := d.Get("facility").(string)
	quantity := d.Get("quantity").(int)

	req := packngo.IPReservationRequest{
		Type:     "public_ipv4",
		Quantity: quantity,
		Facility: facility,
	}

	projectID := d.Get("project_id").(string)

	blockAddr, _, err := client.ProjectIPs.Request(projectID, &req)
	if err != nil {
		return fmt.Errorf("Error reserving IP address block: %s", err)
	}

	d.Set("facility", facility)
	d.Set("quantity", quantity)
	d.Set("project_id", projectID)
	d.SetId(blockAddr.ID)

	return resourcePacketReservedIPBlockRead(d, meta)
}

func loadBlock(d *schema.ResourceData, reservedBlock *packngo.IPAddressReservation) error {
	ipv4CIDRToQuantity := map[int]int{32: 1, 31: 2, 30: 4, 29: 8, 28: 16, 27: 32, 26: 64, 25: 128, 24: 256}

	d.SetId(reservedBlock.ID)
	d.Set("address", reservedBlock.Address)
	d.Set("facility", reservedBlock.Facility.Code)
	d.Set("gateway", reservedBlock.Gateway)
	d.Set("network", reservedBlock.Network)
	d.Set("netmask", reservedBlock.Netmask)
	d.Set("address_family", reservedBlock.AddressFamily)
	d.Set("cidr", reservedBlock.CIDR)
	d.Set("public", reservedBlock.Public)
	d.Set("management", reservedBlock.Management)
	d.Set("manageable", reservedBlock.Manageable)
	if reservedBlock.AddressFamily == 4 {
		d.Set("quantity", ipv4CIDRToQuantity[reservedBlock.CIDR])
	} else {
		// In Packet, reserved IPv6 block is allocated when device is run in a proejct.
		// It's always /56, and it can't be crated with Terraform, only imported.
		// The longest assignable prefix is /64, making it max 256 subnets per block.
		// The following logic will hold as long as /64 is the smallest assignable subnet size.
		bits := 64 - reservedBlock.CIDR
		if bits > 30 {
			return fmt.Errorf("Strange (too small) CIDR prefix: %d", reservedBlock.CIDR)
		}
		d.Set("quantity", 1<<uint(bits))
	}
	d.Set("project_id", path.Base(reservedBlock.Project.Href))
	d.Set("cidr_notation", fmt.Sprintf("%s/%d", reservedBlock.Network, reservedBlock.CIDR))
	return nil

}

func resourcePacketReservedIPBlockRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	id := d.Id()

	reservedBlock, _, err := client.ProjectIPs.Get(id, nil)
	if err != nil {
		return fmt.Errorf("Error reading IP address block with ID %s: %s", id, err)
	}
	err = loadBlock(d, reservedBlock)
	if err != nil {
		return err
	}

	return nil
}

func resourcePacketReservedIPBlockDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	id := d.Id()

	_, err := client.ProjectIPs.Remove(id)

	if err != nil {
		return fmt.Errorf("Error deleting IP reservation block %s: %s", id, err)
	}

	d.SetId("")
	return nil
}
