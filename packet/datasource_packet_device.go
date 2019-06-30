package packet

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/packethost/packngo"
)

func dataSourcePacketDevice() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePacketDeviceRead,
		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePacketDeviceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	hostname := d.Get("hostname").(string)
	projectId := d.Get("project_id").(string)

	ds, _, err := client.Devices.List(projectId, nil)
	if err != nil {
		return err
	}

	dev, err := findDeviceByHostname(ds, hostname)

	if err != nil {
		return err
	}

	d.SetId(dev.ID)
	return nil
}

func findDeviceByHostname(devices []packngo.Device, hostname string) (*packngo.Device, error) {
	results := make([]packngo.Device, 0)
	for _, d := range devices {
		if d.Hostname == hostname {
			results = append(results, d)
		}
	}
	if len(results) == 1 {
		return &results[0], nil
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no device found with hostname %s", hostname)
	}
	return nil, fmt.Errorf("too many devices found with hostname %s (found %d, expected 1)", hostname, len(results))
}
