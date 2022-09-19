package equinix

import (
	"encoding/json"
	"fmt"
	"log"
	"path"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/packethost/packngo"
)

const (
	ReservedIPCreateTimeout = 10 * time.Minute
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
		"vrf_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func metalIPResourceComputedFields() map[string]*schema.Schema {
	s := metalIPComputedFields()
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
				// If metro from HCL is specified, but not present in state, suppress the diff.
				// This is legacy, and I think it's here because of migration, so that old
				// facility reservations are not recreated when metro is specified ???)
				if fromHCL != "" && fromState == "" {
					return true
				}
				// If metro is present in state but not present in HCL, suppress the diff.
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
		ForceNew:    false,
		Description: "Arbitrary description",
	}
	reservedBlockSchema["quantity"] = &schema.Schema{
		Type:         schema.TypeInt,
		Optional:     true,
		ForceNew:     true,
		Computed:     true,
		ExactlyOneOf: []string{"vrf_id", "quantity"},
		Description:  "The number of allocated /32 addresses, a power of 2",
	}
	reservedBlockSchema["type"] = &schema.Schema{
		Type:         schema.TypeString,
		ForceNew:     true,
		Default:      "public_ipv4",
		Optional:     true,
		Description:  "Either global_ipv4, public_ipv4, or vrf. Defaults to public_ipv4.",
		ValidateFunc: validation.StringInSlice([]string{"public_ipv4", "global_ipv4", "vrf"}, false),
	}
	reservedBlockSchema["cidr_notation"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	reservedBlockSchema["tags"] = &schema.Schema{
		Type:        schema.TypeSet,
		ForceNew:    false,
		Description: "Tags attached to the reserved block",
		Optional:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
	}

	reservedBlockSchema["custom_data"] = &schema.Schema{
		Type:             schema.TypeString,
		Default:          "{}",
		ForceNew:         false,
		Description:      "Custom Data is an arbitrary object (submitted in Terraform as serialized JSON) to assign to the IP Reservation. This may be helpful for self-managed IPAM. The object must be valid JSON.",
		Optional:         true,
		ValidateFunc:     validation.StringIsJSON,
		DiffSuppressFunc: structure.SuppressJsonDiff,
		Elem:             &schema.Schema{Type: schema.TypeString},
	}

	reservedBlockSchema["wait_for_state"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Wait for the IP reservation block to reach a desired state on resource creation. One of: `pending`, `created`. The `created` state is default and recommended if the addresses are needed within the configuration. An error will be returned if a timeout or the `denied` state is encountered.",
		Default:     packngo.IPReservationStateCreated,
		Optional:    true,
		ForceNew:    false,
		ValidateDiagFunc: validation.ToDiagFunc(
			validation.StringInSlice([]string{
				string(packngo.IPReservationStateCreated),
				string(packngo.IPReservationStatePending),
			}, false),
		),
	}

	reservedBlockSchema["vrf_id"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		ForceNew:     true,
		ExactlyOneOf: []string{"vrf_id", "quantity"},
		Description:  "VRF ID for type=vrf reservations",
	}
	reservedBlockSchema["network"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		RequiredWith: []string{"vrf_id"},
		ForceNew:     true,
		Computed:     true,
		Description:  "an unreserved network address from an existing vrf ip_range. `network` can only be specified with vrf_id",
	}
	reservedBlockSchema["cidr"] = &schema.Schema{
		Type:         schema.TypeInt,
		Optional:     true,
		RequiredWith: []string{"vrf_id"},
		ForceNew:     true,
		Computed:     true,
		Description:  "the size of the network to reserve from an existing vrf ip_range. `cidr` can only be specified with `vrf_id`. Minimum range is 22-29, with 30-31 supported and necessary for virtual-circuits",
	}
	// TODO: add comments field, used for reservations that are not automatically approved
	return &schema.Resource{
		Create: resourceMetalReservedIPBlockCreate,
		Read:   resourceMetalReservedIPBlockRead,
		Update: resourceMetalReservedIPBlockUpdate,
		Delete: resourceMetalReservedIPBlockDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: reservedBlockSchema,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(ReservedIPCreateTimeout),
		},
	}
}

func resourceMetalReservedIPBlockCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Config).metal
	quantity := d.Get("quantity").(int)
	typ := d.Get("type").(string)

	req := packngo.IPReservationRequest{
		Type:     packngo.IPReservationType(typ),
		Quantity: quantity,
	}
	facility, facOk := d.GetOk("facility")
	metro, metOk := d.GetOk("metro")

	// no need to guard facOk && metOk, they "ConflictsWith" each-other
	switch typ {
	case "global_ipv4":
		if facOk || metOk {
			return fmt.Errorf("facility and metro can't be set for global IP block reservation")
		}
	case "public_ipv4":
		if !(facOk || metOk) {
			return fmt.Errorf("you should set either metro or facility for non-global IP block reservation")
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

	if d.HasChange("custom_data") {
		req.CustomData = d.Get("custom_data")
	}

	req.VRFID = d.Get("vrf_id").(string)
	req.Network = d.Get("network").(string)
	req.CIDR = d.Get("cidr").(int)

	blockAddr, _, err := client.ProjectIPs.Create(projectID, &req)
	if err != nil {
		return fmt.Errorf("error reserving IP address block: %s", err)
	}
	d.Set("project_id", projectID)
	d.SetId(blockAddr.ID)

	wfs := d.Get("wait_for_state").(string)
	log.Printf("[DEBUG] Waiting for IP Reservation (%s) to become %s", d.Id(), wfs)
	target := []string{string(packngo.IPReservationStateCreated)}
	if wfs != string(packngo.IPReservationStateCreated) {
		target = append(target, wfs)
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{string(packngo.IPReservationStatePending)},
		Target:     target,
		Refresh:    reservedIPStateRefreshFunc(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 15 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for IP Reservation (%s) to become %s: %s", d.Id(), wfs, err)
	}

	return resourceMetalReservedIPBlockRead(d, meta)
}

func resourceMetalReservedIPBlockUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Config).metal
	id := d.Id()
	req := &packngo.IPAddressUpdateRequest{}
	if d.HasChange("tags") {
		tags := []string{}
		if tagsRaw, tagsOk := d.GetOk("tags"); tagsOk {
			for _, tag := range tagsRaw.(*schema.Set).List() {
				tags = append(tags, tag.(string))
			}
		}
		req.Tags = &tags
	}

	if d.HasChange("description") {
		desc := d.Get("description").(string)
		req.Description = &desc
	}

	if d.HasChange("custom_data") {
		var v interface{}
		if err := json.Unmarshal([]byte(d.Get("custom_data").(string)), v); err != nil {
			return fmt.Errorf("error unmarshalling custom_data: %w", err)
		}
		req.CustomData = v
	}

	if _, _, err := client.ProjectIPs.Update(id, req, nil); err != nil {
		return fmt.Errorf("error updating IP reservation: %w", err)
	}

	return resourceMetalReservedIPBlockRead(d, meta)
}

func reservedIPStateRefreshFunc(client *packngo.Client, reservedIPId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		reservedIP, _, err := client.ProjectIPs.Get(reservedIPId, nil)
		if err != nil {
			return nil, "", fmt.Errorf("error retrieving reserved IP block %s: %s", reservedIPId, err)
		}

		return reservedIP, string(reservedIP.State), nil
	}
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
	return "", fmt.Errorf("unknown reservation type %+v", r)
}

func loadBlock(d *schema.ResourceData, reservedBlock *packngo.IPAddressReservation) error {
	d.SetId(reservedBlock.ID)

	quantity := 0
	if reservedBlock.AddressFamily == 4 {
		quantity = 1 << (32 - reservedBlock.CIDR)
	} else {
		// In Equinix Metal, a reserved IPv6 block is allocated when a device is
		// run in a project. It's always /56, and it can't be created with
		// Terraform, only imported. The longest assignable prefix is /64,
		// making it max 256 subnets per block. The following logic will hold as
		// long as /64 is the smallest assignable subnet size.
		bits := 64 - reservedBlock.CIDR
		if bits > 30 {
			return fmt.Errorf("strange (too small) CIDR prefix: %d", reservedBlock.CIDR)
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
		"type":           reservedBlock.Type,
		"tags":           reservedBlock.Tags,
		"public":         reservedBlock.Public,
		"management":     reservedBlock.Management,
		"manageable":     reservedBlock.Manageable,
		"quantity":       quantity,
		"project_id":     path.Base(reservedBlock.Project.Href),
		"cidr_notation":  fmt.Sprintf("%s/%d", reservedBlock.Network, reservedBlock.CIDR),
		"custom_data": func(d *schema.ResourceData, k string) error {
			if reservedBlock.CustomData == nil {
				return nil
			}
			b, err := json.Marshal(reservedBlock.CustomData)
			if err != nil {
				return err
			}
			return d.Set(k, string(b))
		},
		"description": func(d *schema.ResourceData, k string) error {
			if (reservedBlock.Description == nil) || (*(reservedBlock.Description) != "") {
				return nil
			}
			return d.Set(k, *(reservedBlock.Description))
		},
		"global": reservedBlock.Global,
		"vrf_id": func(d *schema.ResourceData, k string) error {
			if reservedBlock.VRF == nil {
				return nil
			}
			return d.Set(k, reservedBlock.VRF.ID)
		},
	}

	// filter out attributes which are not defined in target resource
	for k := range attributeMap {
		if d.Get(k) == nil {
			delete(attributeMap, k)
		}
	}

	return setMap(d, attributeMap)
}

func resourceMetalReservedIPBlockRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Config).metal
	id := d.Id()

	getOpts := &packngo.GetOptions{Includes: []string{"facility", "metro", "project", "vrf"}}
	getOpts = getOpts.Filter("types", "public_ipv4,global_ipv4,private_ipv4,public_ipv6,vrf")

	reservedBlock, _, err := client.ProjectIPs.Get(id, getOpts)
	if err != nil {
		err = friendlyError(err)
		if isNotFound(err) {
			log.Printf("[WARN] Reserved IP Block (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error reading IP address block with ID %s: %s", id, err)
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
	client := meta.(*Config).metal

	id := d.Id()

	resp, err := client.ProjectIPs.Remove(id)

	if ignoreResponseErrors(httpForbidden, httpNotFound)(resp, err) != nil {
		return fmt.Errorf("error deleting IP reservation block %s: %s", id, err)
	}

	d.SetId("")
	return nil
}
