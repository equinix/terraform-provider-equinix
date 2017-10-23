package packet

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/packethost/packngo"
)

func dataSourcePacketPreCreatedIPBlock() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePacketReservedIPBlockRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"address_family": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"public": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"facility": {
				Type:     schema.TypeString,
				Required: true,
			},
			"quantity": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cidr_notation": {
				Type:     schema.TypeString,
				Computed: true,
			},
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
			"management": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"manageable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourcePacketReservedIPBlockRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	projectID := d.Get("project_id").(string)
	log.Println("[DEBUG] packet_precreated_ip_block - getting list of IPs in a project")
	ips, _, err := client.ProjectIPs.List(projectID)
	if err != nil {
		return err
	}
	ipv := d.Get("address_family").(int)
	public := d.Get("public").(bool)
	facility := d.Get("facility").(string)

	for _, ip := range ips {
		if ip.Public == public && ip.AddressFamily == ipv && facility == ip.Facility.Code {
			loadBlock(d, &ip)
			break
		}
	}
	if d.Get("cidr_notation") == "" {
		return fmt.Errorf("Could not find matching reserved block, all IPs were %v", ips)
	}
	return nil

}
