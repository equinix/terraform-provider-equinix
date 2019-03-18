package packet

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/packethost/packngo"
)

func resourcePacketPortVlanAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourcePacketPortVlanAttachmentCreate,
		Read:   resourcePacketPortVlanAttachmentRead,
		Delete: resourcePacketPortVlanAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"force_bond": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"device_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vlan_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePacketPortVlanAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	dID := d.Get("device_id").(string)
	pName := d.Get("port_name").(string)
	vID := d.Get("vlan_id").(string)
	log.Printf("[DEBUG] Attaching Port (%s) to VLAN (%s)\n", pName, vID)
	dev, _, err := client.Devices.Get(dID, &packngo.GetOptions{Includes: []string{"network_ports"}})
	if err != nil {
		return err
	}

	portFound := false
	vlanFound := false
	var port packngo.Port
	for _, p := range dev.NetworkPorts {
		if p.Name == pName {
			portFound = true
			port = p
			for _, n := range p.AttachedVirtualNetworks {
				if vID == (filepath.Base(n.Href)) {
					vlanFound = true
					break
				}
			}
			break
		}
	}
	if !portFound {
		return fmt.Errorf("Device %s doesn't have port %s", dID, pName)
	}
	if vlanFound {
		log.Printf("Port %s already has VLAN %s assigned", pName, vID)
		return nil
	}

	if port.Data.Bonded {
		_, _, err := client.DevicePorts.Disbond(&packngo.DisbondRequest{PortID: port.ID, BulkDisable: false})
		if err != nil {
			return friendlyError(err)
		}
	}

	par := &packngo.PortAssignRequest{PortID: port.ID, VirtualNetworkID: vID}

	_, _, err = client.DevicePorts.Assign(par)
	if err != nil {
		return err
	}

	d.SetId(port.ID + ":" + vID)
	return resourcePacketPortVlanAttachmentRead(d, meta)
}

func resourcePacketPortVlanAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	dID := d.Get("device_id").(string)
	pName := d.Get("port_name").(string)
	vID := d.Get("vlan_id").(string)

	dev, _, err := client.Devices.Get(dID, &packngo.GetOptions{Includes: []string{"network_ports"}})
	if err != nil {
		return err
	}
	portFound := false
	vlanFound := false
	portID := ""
	for _, p := range dev.NetworkPorts {
		if p.Name == pName {
			portFound = true
			portID = p.ID
			for _, n := range p.AttachedVirtualNetworks {
				if vID == (filepath.Base(n.Href)) {
					vlanFound = true
					break
				}
			}
			break
		}
	}
	d.Set("port_id", portID)
	if !portFound {
		return fmt.Errorf("Device %s doesn't have port %s", dID, pName)
	}
	if !vlanFound {
		d.SetId(portID)
	}
	return nil
}

func resourcePacketPortVlanAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	pID := d.Get("port_id").(string)
	vID := d.Get("vlan_id").(string)
	par := &packngo.PortAssignRequest{PortID: pID, VirtualNetworkID: vID}
	client := meta.(*packngo.Client)
	portPtr, _, err := client.DevicePorts.Unassign(par)
	if err != nil {
		return err
	}
	forceBond := d.Get("force_bond").(bool)
	if forceBond && (len(portPtr.AttachedVirtualNetworks) == 0) {
		_, _, err = client.DevicePorts.Bond(&packngo.BondRequest{PortID: pID, BulkEnable: false})
		if err != nil {
			return friendlyError(err)
		}
	}
	return nil
}
