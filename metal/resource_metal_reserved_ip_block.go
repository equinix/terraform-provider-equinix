package metal

import (
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/packethost/packngo"
)

func metalIPComputedFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"address": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"address_family": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Address family as integer (4 or 6)",
		},
		"cidr": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Length of CIDR prefix of the block as integer",
		},
		"gateway": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"netmask": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Mask in decimal notation, e.g. 255.255.255.0",
		},
		"network": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Network IP address portion of the block specification",
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
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Address family as integer (4 or 6)",
	}
	s["public"] = &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Flag indicating whether IP block is addressable from the Internet",
	}
	s["global"] = &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Flag indicating whether IP block is global, i.e. assignable in any location",
	}
	return s
}

func resourceMetalReservedIPBlock() *schema.Resource {
	reservedBlockSchema := metalIPResourceComputedFields()
	reservedBlockSchema["project_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The metal project ID where to allocate the address block",
	}
	reservedBlockSchema["facility"] = &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		ForceNew:      true,
		ConflictsWith: []string{"metro"},
		Description:   "Facility where to allocate the public IP address block, makes sense only for type==public_ipv4, must be empty for type==global_ipv4, conflicts with metro",
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			// suppress diff when unsetting facility
			if len(old) > 0 && new == "" {
				return true
			}
			return old == new
		},
	}
	reservedBlockSchema["metro"] = &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		ForceNew:      true,
		ConflictsWith: []string{"facility"},
		Description:   "Metro where to allocate the public IP address block, makes sense only for type==public_ipv4, must be empty for type==global_ipv4, conflicts with facility",
		DiffSuppressFunc: func(k, fromState, fromHCL string, d *schema.ResourceData) bool {
			_, facOk := d.GetOk("facility")

			// if facility is not in state, treat the diff normally, otherwise do following messy checks:
			if facOk {
				// If metro from HCL is specified, but not present in state, supress the diff.
				// This is legacy, and I think it's here because of migration, so that old
				// facility reservations are not recreated when metro is specified ???)
				if fromHCL != "" && fromState == "" {
					return true
				}
				// If metro is present in state but not present in HCL, supress the diff.
				// This is for "facility-specified" reservation blocks created after ~July 2021.
				// These blocks will have metro "computed" to the TF state, and we don't want to
				// emit a diff if the metro field is empty in HCL.
				if fromHCL == "" && fromState != "" {
					return true
				}
			}
			return fromState == fromHCL
		},
		StateFunc: toLower,
	}
	reservedBlockSchema["description"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "Arbitrary description",
	}
	reservedBlockSchema["quantity"] = &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		ForceNew:    true,
		Description: "The number of allocated /32 addresses, a power of 2",
	}
	reservedBlockSchema["type"] = &schema.Schema{
		Type:         schema.TypeString,
		ForceNew:     true,
		Default:      "public_ipv4",
		Optional:     true,
		Description:  "Either global_ipv4 or public_ipv4, defaults to public_ipv4 for backward compatibility",
		ValidateFunc: validation.StringInSlice([]string{"public_ipv4", "global_ipv4"}, false),
	}
	reservedBlockSchema["cidr_notation"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	reservedBlockSchema["tags"] = &schema.Schema{
		Type:        schema.TypeSet,
		ForceNew:    true,
		Description: "Tags attached to the reserved block",
		Optional:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
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
	facility, facOk := d.GetOk("facility")
	metro, metOk := d.GetOk("metro")

	// no need to guard facOk && metOk, they "ConflictsWith" each-other
	if typ == "global_ipv4" {
		if facOk || metOk {
			return fmt.Errorf("facility and metro can't be set for global IP block reservation")
		}
	} else {
		if !(facOk || metOk) {
			return fmt.Errorf("You should set either metro or facility for non-global IP block reservation")
		}
	}

	if facOk {
		f := facility.(string)
		req.Facility = &f
	}
	if metOk {
		m := metro.(string)
		req.Metro = &m
	}
	desc, ok := d.GetOk("description")
	if ok {
		req.Description = desc.(string)
	}

	if tagsRaw, tagsOk := d.GetOk("tags"); tagsOk {
		for _, tag := range tagsRaw.(*schema.Set).List() {
			req.Tags = append(req.Tags, tag.(string))
		}
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

func getType(r *packngo.IPAddressReservation) (string, error) {
	switch {
	case !r.Public:
		return fmt.Sprintf("private_ipv%d", r.AddressFamily), nil
	case r.Public && !r.Global:
		return fmt.Sprintf("public_ipv%d", r.AddressFamily), nil
	case r.Public && r.Global:
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

	attributeMap := map[string]interface{}{
		"address": reservedBlock.Address,
		"facility": func(d *schema.ResourceData, k string) error {
			if reservedBlock.Facility == nil {
				return nil
			}
			return d.Set(k, reservedBlock.Facility.Code)
		},
		"metro": func(d *schema.ResourceData, k string) error {
			if reservedBlock.Metro == nil {
				return nil
			}
			return d.Set(k, strings.ToLower(reservedBlock.Metro.Code))
		},
		"gateway":        reservedBlock.Gateway,
		"network":        reservedBlock.Network,
		"netmask":        reservedBlock.Netmask,
		"address_family": reservedBlock.AddressFamily,
		"cidr":           reservedBlock.CIDR,
		"type":           typ,
		"tags":           reservedBlock.Tags,
		"public":         reservedBlock.Public,
		"management":     reservedBlock.Management,
		"manageable":     reservedBlock.Manageable,
		"quantity":       quantity,
		"project_id":     path.Base(reservedBlock.Project.Href),
		"cidr_notation":  fmt.Sprintf("%s/%d", reservedBlock.Network, reservedBlock.CIDR),
	}

	// filter out attributes which are not defined in target resource
	for k := range attributeMap {
		if d.Get(k) == nil {
			delete(attributeMap, k)
		}
	}

	err = setMap(d, attributeMap)
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
	d.Set("global", reservedBlock.Global)

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
