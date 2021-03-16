package metal

import (
	"fmt"
	"log"
	"path"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/packethost/packngo"
)

func metalIPComputedFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"address": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"address_family": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"cidr": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"gateway": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"netmask": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"network": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"manageable": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"management": {
			Type:     schema.TypeBool,
			Computed: true,
		},
	}
}

func metalIPResourceComputedFields() map[string]*schema.Schema {
	s := metalIPComputedFields()
	s["address_family"] = &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true,
	}
	s["public"] = &schema.Schema{
		Type:     schema.TypeBool,
		Computed: true,
	}
	s["global"] = &schema.Schema{
		Type:     schema.TypeBool,
		Computed: true,
	}
	return s
}

func resourceMetalReservedIPBlock() *schema.Resource {
	reservedBlockSchema := metalIPResourceComputedFields()
	reservedBlockSchema["project_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	}
	reservedBlockSchema["facility"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
	}
	reservedBlockSchema["description"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
	}
	reservedBlockSchema["quantity"] = &schema.Schema{
		Type:     schema.TypeInt,
		Required: true,
		ForceNew: true,
	}
	reservedBlockSchema["type"] = &schema.Schema{
		Type:         schema.TypeString,
		ForceNew:     true,
		Default:      "public_ipv4",
		Optional:     true,
		ValidateFunc: validation.StringInSlice([]string{"public_ipv4", "global_ipv4"}, false),
	}
	reservedBlockSchema["cidr_notation"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	return &schema.Resource{
		Create: resourceMetalReservedIPBlockCreate,
		Read:   resourceMetalReservedIPBlockRead,
		Delete: resourceMetalReservedIPBlockDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: reservedBlockSchema,
	}
}

func resourceMetalReservedIPBlockCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	quantity := d.Get("quantity").(int)
	typ := d.Get("type").(string)

	req := packngo.IPReservationRequest{
		Type:     typ,
		Quantity: quantity,
	}
	f, ok := d.GetOk("facility")

	if ok && typ == "global_ipv4" {
		return fmt.Errorf("Facility can not be set for type == global_ipv4")
	}
	fs := f.(string)
	if typ == "public_ipv4" {
		req.Facility = &fs
	}
	desc, ok := d.GetOk("description")
	if ok {
		req.Description = desc.(string)
	}

	projectID := d.Get("project_id").(string)

	blockAddr, _, err := client.ProjectIPs.Request(projectID, &req)
	if err != nil {
		return fmt.Errorf("Error reserving IP address block: %s", err)
	}

	d.Set("project_id", projectID)
	d.SetId(blockAddr.ID)

	return resourceMetalReservedIPBlockRead(d, meta)
}

func getGlobalBool(r *packngo.IPAddressReservation) bool {
	if r.Global != nil {
		return *(r.Global)
	}
	return false
}

func getType(r *packngo.IPAddressReservation) (string, error) {
	globalBool := getGlobalBool(r)
	switch {
	case !r.Public:
		return fmt.Sprintf("private_ipv%d", r.AddressFamily), nil
	case r.Public && !globalBool:
		return fmt.Sprintf("public_ipv%d", r.AddressFamily), nil
	case r.Public && globalBool:
		return fmt.Sprintf("global_ipv%d", r.AddressFamily), nil
	}
	return "", fmt.Errorf("Unknown reservation type %+v", r)
}

func loadBlock(d *schema.ResourceData, reservedBlock *packngo.IPAddressReservation) error {
	ipv4CIDRToQuantity := map[int]int{32: 1, 31: 2, 30: 4, 29: 8, 28: 16, 27: 32, 26: 64, 25: 128, 24: 256}

	d.SetId(reservedBlock.ID)

	typ, err := getType(reservedBlock)
	if err != nil {
		return err
	}
	quantity := 0
	if reservedBlock.AddressFamily == 4 {
		quantity = ipv4CIDRToQuantity[reservedBlock.CIDR]
	} else {
		// In Equinix Metal, a reserved IPv6 block is allocated when a device is
		// run in a project. It's always /56, and it can't be created with
		// Terraform, only imported. The longest assignable prefix is /64,
		// making it max 256 subnets per block. The following logic will hold as
		// long as /64 is the smallest assignable subnet size.
		bits := 64 - reservedBlock.CIDR
		if bits > 30 {
			return fmt.Errorf("Strange (too small) CIDR prefix: %d", reservedBlock.CIDR)
		}
		quantity = 1 << uint(bits)
	}

	err = setMap(d, map[string]interface{}{
		"address": reservedBlock.Address,
		"facility": func(d *schema.ResourceData, k string) error {
			if reservedBlock.Facility == nil {
				return nil
			}
			return d.Set(k, reservedBlock.Facility.Code)
		},
		"gateway":        reservedBlock.Gateway,
		"network":        reservedBlock.Network,
		"netmask":        reservedBlock.Netmask,
		"address_family": reservedBlock.AddressFamily,
		"cidr":           reservedBlock.CIDR,
		"type":           typ,
		"public":         reservedBlock.Public,
		"management":     reservedBlock.Management,
		"manageable":     reservedBlock.Manageable,
		"quantity":       quantity,
		"project_id":     path.Base(reservedBlock.Project.Href),
		"cidr_notation":  fmt.Sprintf("%s/%d", reservedBlock.Network, reservedBlock.CIDR),
	})
	return err
}

func resourceMetalReservedIPBlockRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	id := d.Id()

	reservedBlock, _, err := client.ProjectIPs.Get(id, nil)
	if err != nil {
		err = friendlyError(err)
		if isNotFound(err) {
			log.Printf("[WARN] Reserved IP Block (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading IP address block with ID %s: %s", id, err)
	}
	err = loadBlock(d, reservedBlock)
	if err != nil {
		return err
	}

	if (reservedBlock.Description != nil) && (*(reservedBlock.Description) != "") {
		d.Set("description", *(reservedBlock.Description))
	}
	d.Set("global", getGlobalBool(reservedBlock))

	return nil
}

func resourceMetalReservedIPBlockDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	id := d.Id()

	resp, err := client.ProjectIPs.Remove(id)

	if ignoreResponseErrors(httpForbidden, httpNotFound)(resp, err) != nil {
		return fmt.Errorf("Error deleting IP reservation block %s: %s", id, err)
	}

	d.SetId("")
	return nil
}
