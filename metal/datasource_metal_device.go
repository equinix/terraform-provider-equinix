package metal

import (
	"encoding/json"
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/packethost/packngo"
)

func dataSourceMetalDevice() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMetalDeviceRead,
		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:          schema.TypeString,
				Description:   "The device name",
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"device_id"},
			},
			"project_id": {
				Type:          schema.TypeString,
				Description:   "The id of the project in which the devices exists",
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"device_id"},
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description string for the device",
				Computed:    true,
			},
			"device_id": {
				Type:          schema.TypeString,
				Description:   "Device ID",
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"project_id", "hostname"},
			},
			"facility": {
				Type:        schema.TypeString,
				Description: "The facility where the device is deployed",
				Computed:    true,
			},
			"metro": {
				Type:        schema.TypeString,
				Description: "The metro where the device is deployed",
				Computed:    true,
			},
			"plan": {
				Type:        schema.TypeString,
				Description: "The hardware config of the device",
				Computed:    true,
			},
			"operating_system": {
				Type:        schema.TypeString,
				Description: "The operating system running on the device",
				Computed:    true,
			},
			"state": {
				Type:        schema.TypeString,
				Description: "The state of the device",
				Computed:    true,
			},
			"billing_cycle": {
				Type:        schema.TypeString,
				Description: "The billing cycle of the device (monthly or hourly)",
				Computed:    true,
			},
			"access_public_ipv6": {
				Type:        schema.TypeString,
				Description: "The ipv6 management IP assigned to the device",
				Computed:    true,
			},

			"access_public_ipv4": {
				Type:        schema.TypeString,
				Description: "The ipv4 management IP assigned to the device",
				Computed:    true,
			},
			"access_private_ipv4": {
				Type:        schema.TypeString,
				Description: "The ipv4 private IP assigned to the device",
				Computed:    true,
			},
			"tags": {
				Type:        schema.TypeList,
				Description: "Tags attached to the device",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"ssh_key_ids": {
				Type:        schema.TypeList,
				Description: "List of IDs of SSH keys deployed in the device, can be both user or project SSH keys",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"network_type": {
				Type:        schema.TypeString,
				Description: "L2 network type of the device, one of" + NetworkTypeList,
				Computed:    true,
			},
			"hardware_reservation_id": {
				Type:        schema.TypeString,
				Description: "The id of hardware reservation which this device occupies",
				Computed:    true,
			},
			"storage": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"root_password": {
				Type:        schema.TypeString,
				Description: "Root password to the server (if still available)",
				Computed:    true,
				Sensitive:   true,
			},
			"always_pxe": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ipxe_script_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network": {
				Type:        schema.TypeList,
				Description: "The device's private and public IP (v4 and v6) network details. When a device is run without any special network configuration, it will have 3 networks: ublic IPv4 at metal_device.name.network.0, IPv6 at metal_device.name.network.1 and private IPv4 at metal_device.name.network.2. Elastic addresses then stack by type - an assigned public IPv4 will go after the management public IPv4 (to index 1), and will then shift the indices of the IPv6 and private IPv4. Assigned private IPv4 will go after the management private IPv4 (to the end of the network list).",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:        schema.TypeString,
							Description: "IPv4 or IPv6 address string",
							Computed:    true,
						},
						"gateway": {
							Type:        schema.TypeString,
							Description: "Address of router",
							Computed:    true,
						},
						"family": {
							Type:        schema.TypeInt,
							Description: "IP version - \"4\" or \"6\"",
							Computed:    true,
						},
						"cidr": {
							Type:        schema.TypeInt,
							Description: "Bit length of the network mask of the address",
							Computed:    true,
						},
						"public": {
							Type:        schema.TypeBool,
							Description: "Whether the address is routable from the Internet",
							Computed:    true,
						},
					},
				},
			},
			"ports": {
				Type:        schema.TypeList,
				Description: "Ports assigned to the device",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "Name of the port (e.g. eth0, or bond0)",
							Computed:    true,
						},
						"id": {
							Type:        schema.TypeString,
							Description: "The ID of the device",
							Computed:    true,
						},
						"type": {
							Type:        schema.TypeString,
							Description: "Type of the port (e.g. NetworkPort or NetworkBondPort)",
							Computed:    true,
						},
						"mac": {
							Type:        schema.TypeString,
							Description: "MAC address assigned to the port",
							Computed:    true,
						},
						"bonded": {
							Type:        schema.TypeBool,
							Description: "Whether this port is part of a bond in bonded network setup",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceMetalDeviceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	hostnameRaw, hostnameOK := d.GetOk("hostname")
	projectIdRaw, projectIdOK := d.GetOk("project_id")
	deviceIdRaw, deviceIdOK := d.GetOk("device_id")

	if !deviceIdOK && !hostnameOK {
		return fmt.Errorf("You must supply device_id or hostname")
	}
	var device *packngo.Device
	if hostnameOK {
		if !projectIdOK {
			return fmt.Errorf("If you lookup via hostname, you must supply project_id")
		}
		hostname := hostnameRaw.(string)
		projectId := projectIdRaw.(string)

		ds, _, err := client.Devices.List(projectId, &packngo.ListOptions{Search: hostname})
		if err != nil {
			return err
		}

		device, err = findDeviceByHostname(ds, hostname)
		if err != nil {
			return err
		}
	} else {
		deviceId := deviceIdRaw.(string)
		var err error
		device, _, err = client.Devices.Get(deviceId, nil)
		if err != nil {
			return err
		}
	}

	d.Set("hostname", device.Hostname)
	d.Set("project_id", device.Project.ID)
	d.Set("device_id", device.ID)
	d.Set("plan", device.Plan.Slug)
	d.Set("facility", device.Facility.Code)
	if device.Metro != nil {
		d.Set("metro", strings.ToLower(device.Metro.Code))
	}
	d.Set("operating_system", device.OS.Slug)
	d.Set("state", device.State)
	d.Set("billing_cycle", device.BillingCycle)
	d.Set("ipxe_script_url", device.IPXEScriptURL)
	d.Set("always_pxe", device.AlwaysPXE)
	d.Set("root_password", device.RootPassword)
	if device.Storage != nil {
		rawStorageBytes, err := json.Marshal(device.Storage)
		if err != nil {
			return fmt.Errorf("[ERR] Error getting storage JSON string for device (%s): %s", d.Id(), err)
		}

		storageString, err := structure.NormalizeJsonString(string(rawStorageBytes))
		if err != nil {
			return fmt.Errorf("[ERR] Error normalizing storage JSON string for device (%s): %s", d.Id(), err)
		}
		d.Set("storage", storageString)
	}

	if device.HardwareReservation != nil {
		d.Set("hardware_reservation_id", device.HardwareReservation.ID)
	}
	networkType := device.GetNetworkType()

	d.Set("network_type", networkType)

	d.Set("tags", device.Tags)

	keyIDs := []string{}
	for _, k := range device.SSHKeys {
		keyIDs = append(keyIDs, path.Base(k.URL))
	}
	d.Set("ssh_key_ids", keyIDs)
	networkInfo := getNetworkInfo(device.Network)

	sort.SliceStable(networkInfo.Networks, func(i, j int) bool {
		famI := networkInfo.Networks[i]["family"].(int)
		famJ := networkInfo.Networks[j]["family"].(int)
		pubI := networkInfo.Networks[i]["public"].(bool)
		pubJ := networkInfo.Networks[j]["public"].(bool)
		return getNetworkRank(famI, pubI) < getNetworkRank(famJ, pubJ)
	})

	d.Set("network", networkInfo.Networks)
	d.Set("access_public_ipv4", networkInfo.PublicIPv4)
	d.Set("access_private_ipv4", networkInfo.PrivateIPv4)
	d.Set("access_public_ipv6", networkInfo.PublicIPv6)

	ports := getPorts(device.NetworkPorts)
	d.Set("ports", ports)

	d.SetId(device.ID)
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
