package metal

import (
	"fmt"
	"net"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
)

func dataSourceMetalReservedIPBlock() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMetalReservedIPBlockRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "ID of the block to look up",
				ConflictsWith: []string{"project_id", "address"},
				Computed:      true,
			},
			"project_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Description:   "ID of the project where the searched block should be",
				ConflictsWith: []string{"id"},
			},
			"ip_address": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Find block containing this IP address in given project",
				ConflictsWith: []string{"id"},
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
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Address type, one of public_ipv4, public_ipv6 and private_ipv4",
			},
		},
	}
}

func dataSourceMetalReservedIPBlockRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	blockId, blockIdOk := d.GetOk("id")
	projectId, projectIdOk := d.GetOk("project_id")
	address, addressOk := d.GetOk("ip_address")

	if !(blockIdOk || (projectIdOk && addressOk)) {
		return fmt.Errorf("You must specify either id or project_id and ip_address")
	}
	if blockIdOk {
		block, _, err := client.ProjectIPs.Get(
			blockId.(string),
			&packngo.GetOptions{Includes: []string{"facility", "project", "metro"}},
		)
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

	blocks, _, err := client.ProjectIPs.List(projectId.(string),
		&packngo.GetOptions{Includes: []string{"facility", "project", "metro"}},
	)
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
	return fmt.Errorf("Could not find matching reserved block, all blocks were \n%s", listOfCidrs(blocks))

}

func listOfCidrs(blocks []packngo.IPAddressReservation) string {
	cidrs := []string{}
	for _, b := range blocks {
		cidrs = append(cidrs, fmt.Sprintf("%s/%d", b.Network, b.CIDR))
	}
	return strings.Join(cidrs, "\n")

}
