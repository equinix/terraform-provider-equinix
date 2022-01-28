package equinix

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		Type:        schema.TypeInt,
		Required:    true,
		Description: "4 or 6, depending on which block you are looking for.",
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
	client := meta.(*Config).Client()
	projectID := d.Get("project_id").(string)
	ips, _, err := client.ProjectIPs.List(projectID, nil)
	if err != nil {
		return err
	}
	ipv := d.Get("address_family").(int)
	public := d.Get("public").(bool)
	global := d.Get("global").(bool)

	if !public && global {
		return fmt.Errorf("Private (non-public) global IP address blocks are not supported in Equinix Metal")
	}

	fval, fok := d.GetOk("facility")
	mval, mok := d.GetOk("metro")
	if (fok || mok) && global {
		return fmt.Errorf("You can't specify facility for global IP block - addresses from global blocks can be assigned to devices across several locations")
	}

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
		// lookup of blcok specified with metro
		metro := mval.(string)
		for _, ip := range ips {
			if ip.Metro == nil {
				continue
			}
			if ip.Public == public && ip.AddressFamily == ipv && metro == ip.Metro.Code {
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
	return fmt.Errorf("Could not find matching reserved block, all IPs were %v", ips)

}
