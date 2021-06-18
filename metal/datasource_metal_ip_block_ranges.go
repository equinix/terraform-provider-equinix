package metal

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
)

func dataSourceMetalIPBlockRanges() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMetalIPBlockRangesRead,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Description: "ID of the project from which to list the blocks",
				Required:    true,
			},
			"facility": {
				Type:        schema.TypeString,
				Description: "Facility code filtering the IP blocks. Global IPv4 blcoks will be listed anyway. If you omit this and metro, all the block from the project will be listed",
				Optional:    true,
			},
			"metro": {
				Type:        schema.TypeString,
				Description: "Metro code filtering the IP blocks. Global IPv4 blcoks will be listed anyway. If you omit this and facility, all the block from the project will be listed",
				Optional:    true,
				StateFunc:   toLower,
			},
			"public_ipv4": {
				Type:        schema.TypeList,
				Description: "List of CIDR expressions for Public IPv4 blocks in the project",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"global_ipv4": {
				Type:        schema.TypeList,
				Description: "List of CIDR expressions for Global IPv4 blocks in the project",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"private_ipv4": {
				Type:        schema.TypeList,
				Description: "List of CIDR expressions for Private IPv4 blocks in the project",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"ipv6": {
				Type:        schema.TypeList,
				Description: "List of CIDR expressions for IPv6 blocks in the project",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
		},
	}
}

func facilityMatch(ref string, facility *packngo.Facility) bool {
	if ref == "" {
		return true
	}
	if facility != nil && ref == facility.Code {
		return true
	}
	return false
}

func metroMatch(ref string, metro *packngo.Metro) bool {
	if ref == "" {
		return true
	}
	if metro != nil && ref == metro.Code {
		return true
	}
	return false
}

func dataSourceMetalIPBlockRangesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	projectID := d.Get("project_id").(string)
	ips, _, err := client.ProjectIPs.List(projectID, nil)
	if err != nil {
		return err
	}

	facility := d.Get("facility").(string)
	metro := d.Get("metro").(string)

	publicIPv4s := []string{}
	globalIPv4s := []string{}
	privateIPv4s := []string{}
	theIPv6s := []string{}
	var targetSlice *[]string

	for _, ip := range ips {
		targetSlice = nil
		cnStr := fmt.Sprintf("%s/%d", ip.Network, ip.CIDR)
		if ip.AddressFamily == 4 {
			if ip.Public {
				if ip.Global {
					globalIPv4s = append(globalIPv4s, cnStr)
				} else {
					targetSlice = &publicIPv4s
				}
			} else {
				targetSlice = &privateIPv4s
			}
		} else {
			targetSlice = &theIPv6s
		}
		if targetSlice != nil && facilityMatch(facility, ip.Facility) && metroMatch(metro, ip.Metro) {
			*targetSlice = append(*targetSlice, cnStr)
		}
	}

	d.Set("public_ipv4", publicIPv4s)
	d.Set("global_ipv4", globalIPv4s)
	d.Set("private_ipv4", privateIPv4s)
	d.Set("ipv6", theIPv6s)

	id := projectID

	// use facility in the index for pre-metro compatibility
	// facility is always returned now, metros is added for future-proofing.
	// facility and metro codes to not clash.
	if facility != "" {
		id = id + "-" + facility
	} else if metro != "" {
		id = id + "-" + metro
	}

	d.SetId(id + "-IPs")
	return nil

}
