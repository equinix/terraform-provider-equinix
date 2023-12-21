package equinix

import (
	"fmt"
	"net"
	"strings"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
)

func DataSourceMetalReservedIPBlock() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMetalReservedIPBlockRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "ID of the block to look up",
				ConflictsWith: []string{"project_id", "ip_address"},
				Computed:      true,
			},
			"project_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Description:   "ID of the project where the searched block should be",
				ConflictsWith: []string{"id"},
				RequiredWith:  []string{"ip_address"},
			},
			"ip_address": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Find block containing this IP address in given project",
				ConflictsWith: []string{"id"},
				RequiredWith:  []string{"project_id"},
			},

			"global": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Addresses from block are attachable in all locations",
			},
			"public": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Addresses from public block are routeable from the Internet",
			},
			"facility": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Facility of the block. (for non-global blocks)",
				Deprecated:  "Use metro instead of facility.  For more information, read the migration guide: https://registry.terraform.io/providers/equinix/equinix/latest/docs/guides/migration_guide_facilities_to_metros_devices",
			},
			"metro": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Metro of the block (for non-global blocks)",
			},
			"address_family": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "4 or 6",
			},
			"address": {
				// I honestly don't know what this "address" is. Maybe next available?
				Type:     schema.TypeString,
				Computed: true,
			},
			"cidr_notation": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CIDR notation of the looked up block",
			},
			"cidr": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Length of CIDR prefix of the block as integer",
			},
			"gateway": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP address of gateway for the block",
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
			"quantity": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vrf_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "VRF ID of the block when type=vrf",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Address type, one of public_ipv4, public_ipv6, private_ipv4, global_ipv4, and vrf",
			},
		},
	}
}

func dataSourceMetalReservedIPBlockRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*config.Config).Metal

	blockId, blockIdOk := d.GetOk("id")
	projectId, projectIdOk := d.GetOk("project_id")
	address, addressOk := d.GetOk("ip_address")
	getOpts := &packngo.GetOptions{Includes: []string{"facility", "metro", "project", "vrf"}}
	getOpts = getOpts.Filter("types", "public_ipv4,global_ipv4,private_ipv4,public_ipv6,vrf")

	if !(blockIdOk || (projectIdOk && addressOk)) {
		return fmt.Errorf("you must specify either id or project_id and ip_address")
	}
	if blockIdOk {
		block, _, err := client.ProjectIPs.Get(
			blockId.(string), getOpts)
		if err != nil {
			return err
		}
		return loadBlock(d, block)
	}
	// we search by project_id and ip_address
	addressStr := address.(string)
	lookupAddress := net.ParseIP(addressStr)
	if lookupAddress == nil {
		return fmt.Errorf("%s is not a valid ip_address", addressStr)
	}

	blocks, _, err := client.ProjectIPs.List(projectId.(string), getOpts)
	if err != nil {
		return err
	}
	for _, b := range blocks {
		cidr := fmt.Sprintf("%s/%d", b.Network, b.CIDR)
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			return fmt.Errorf("CIDR expression of an Equinix Metal IP Block could not be parsed: %s. Please report this in a GitHub issue", cidr)
		}

		if ipNet.Contains(lookupAddress) {
			d.Set("id", b.ID)
			return loadBlock(d, &b)
		}
	}
	return fmt.Errorf("could not find matching reserved block, all blocks were \n%s", listOfCidrs(blocks))
}

func listOfCidrs(blocks []packngo.IPAddressReservation) string {
	cidrs := []string{}
	for _, b := range blocks {
		cidrs = append(cidrs, fmt.Sprintf("%s/%d", b.Network, b.CIDR))
	}
	return strings.Join(cidrs, "\n")
}
