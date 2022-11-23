package equinix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/packethost/packngo"
)

func dataSourceMetalPreCreatedIPBlock() *schema.Resource {
	s := metalIPComputedFields()
	s["project_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "ID of the project where the searched block should be.",
	}
	s["global"] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Whether to look for global block. Default is false for backward compatibility.",
	}
	s["public"] = &schema.Schema{
		Type:        schema.TypeBool,
		Required:    true,
		Description: "Whether to look for public or private block.",
	}
	s["facility"] = &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		Description:   "Facility of the searched block. (for non-global blocks).",
		ConflictsWith: []string{"metro"},
	}
	s["metro"] = &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		Description:   "Metro of the searched block (for non-global blocks).",
		ConflictsWith: []string{"facility"},
	}
	s["address_family"] = &schema.Schema{
		Type:         schema.TypeInt,
		Required:     true,
		Description:  "4 or 6, depending on which block you are looking for.",
		ValidateFunc: validation.IntInSlice([]int{4, 6}),
	}
	s["cidr_notation"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "CIDR notation of the looked up block.",
	}
	s["quantity"] = &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true,
	}
	s["type"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	return &schema.Resource{
		Read:   dataSourceMetalPreCreatedIPBlockRead,
		Schema: s,
	}
}

func dataSourceMetalPreCreatedIPBlockRead(d *schema.ResourceData, meta interface{}) error {
	var types string
	client := meta.(*Config).metal
	projectID := d.Get("project_id").(string)

	ipv := d.Get("address_family").(int)
	public := d.Get("public").(bool)
	global := d.Get("global").(bool)
	fval, fok := d.GetOk("facility")
	mval, mok := d.GetOk("metro")

	if !public && global {
		return fmt.Errorf("private (non-public) global IP address blocks are not supported in Equinix Metal")
	}

	if (fok || mok) && global {
		return fmt.Errorf("you can't specify facility for global IP block - addresses from global blocks can be assigned to devices across several locations")
	}

	// Public and Address Family are required, prefilter types list based on
	// these values. Global is non-default and exclusive, so we can also filter
	// types on that.
	switch {
	case global:
		types = "global_ipv4"
	case !public:
		types = "private_ipv4,vrf"
	case public:
		switch ipv {
		case 4:
			types = "public_ipv4,global_ipv4"
		case 6:
			types = "public_ipv6"
		}
	}

	getOpts := &packngo.GetOptions{Includes: []string{"facility", "metro", "project", "vrf"}}
	getOpts = getOpts.Filter("types", types)

	ips, _, err := client.ProjectIPs.List(projectID, getOpts)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] filtering ips by (project: %s, facility: %s, metro: %s, family: %d, global: %t)", projectID, fval.(string), mval.(string), ipv, global)

	if fok {
		// lookup of block specified with facility
		facility := fval.(string)
		for _, ip := range ips {
			if ip.Facility == nil {
				continue
			}
			if ip.Public == public && ip.AddressFamily == ipv && facility == ip.Facility.Code {
				return loadBlock(d, &ip)
			}
		}
	} else if mok {
		// lookup of block specified with metro
		metro := mval.(string)
		for _, ip := range ips {
			ipMetro := ip.Metro
			if ip.Metro == nil {
				if ip.Facility.Metro == nil {
					continue
				}
				ipMetro = ip.Facility.Metro
			}
			if ip.Public == public && ip.AddressFamily == ipv && metro == ipMetro.Code {
				return loadBlock(d, &ip)
			}
		}
	} else {
		// lookup of blocks not specified with facility or metro
		for _, ip := range ips {
			if ip.Public == public && ip.AddressFamily == ipv && global == ip.Global {
				return loadBlock(d, &ip)
			}
		}
	}
	log.Printf("[DEBUG] filter not matched in response ips: %v", ips)
	return fmt.Errorf("could not find matching reserved block")
}
