package packet

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/packethost/packngo"
)

var (
	computedFields = map[string]schema.ValueType{
		"address":        schema.TypeString,
		"gateway":        schema.TypeString,
		"network":        schema.TypeString,
		"netmask":        schema.TypeString,
		"cidr":           schema.TypeInt,
		"address_family": schema.TypeInt,
		"public":         schema.TypeBool,
		"management":     schema.TypeBool,
		"manageable":     schema.TypeBool,
	}
)

func resourcePacketReservedIPBlock() *schema.Resource {
	reservedBlockSchema := map[string]*schema.Schema{
		"project_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"facility": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"quantity": &schema.Schema{
			Type:     schema.TypeInt,
			Required: true,
			ForceNew: true,
		},
		"cidr_notation": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
	}
	for k, v := range computedFields {
		reservedBlockSchema[k] = &schema.Schema{
			Type:     v,
			Computed: true,
		}
	}

	return &schema.Resource{
		Create: resourcePacketReservedIPBlockCreate,
		Read:   resourcePacketReservedIPBlockRead,
		Delete: resourcePacketReservedIPBlockDelete,

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

	d.Set("cidr_notation", blockAddr.Address)
	d.Set("facility", facility)
	d.Set("quantity", quantity)
	d.Set("project_id", projectID)

	return resourcePacketReservedIPBlockRead(d, meta)
}

func resourcePacketReservedIPBlockRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	id := d.Id()
	cidrNotation := d.Get("cidr_notation").(string)
	projectID := d.Get("project_id").(string)
	var reservedBlock *packngo.IPAddressReservation
	var err error

	if len(id) == 0 {
		if len(cidrNotation) == 0 {
			return fmt.Errorf("can't read device %v", d)
		}
		reservedBlock, _, err = client.ProjectIPs.GetByCIDR(projectID, cidrNotation)
		if err != nil {
			return fmt.Errorf("Error re-reading IP block with CIDR notation %s: %s", cidrNotation, err)
		}
	} else {
		reservedBlock, _, err = client.ProjectIPs.Get(id)
		if err != nil {
			return fmt.Errorf("Error reading IP address block with ID %s: %s", id, err)
		}
	}
	cidrToQuantity := map[int]int{32: 1, 31: 2, 30: 4, 29: 8, 28: 16, 27: 32, 26: 64, 25: 128, 24: 256}

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
	d.Set("quantity", cidrToQuantity[reservedBlock.CIDR])
	d.Set("cidr_notation", fmt.Sprintf("%s/%d", reservedBlock.Address, reservedBlock.CIDR))

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
