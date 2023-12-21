package equinix

import (
	"fmt"
	"log"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/mutexkv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
)

func ResourceMetalPortVlanAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceMetalPortVlanAttachmentCreate,
		Read:   resourceMetalPortVlanAttachmentRead,
		Delete: resourceMetalPortVlanAttachmentDelete,
		Update: resourceMetalPortVlanAttachmentUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"force_bond": {
				Type:        schema.TypeBool,
				Description: "Add port back to the bond when this resource is removed. Default is false",
				Optional:    true,
				Default:     false,
				ForceNew:    true,
			},
			"device_id": {
				Type:        schema.TypeString,
				Description: "ID of device to be assigned to the VLAN",
				Required:    true,
				ForceNew:    true,
			},
			"port_name": {
				Type:        schema.TypeString,
				Description: "Name of network port to be assigned to the VLAN",
				Required:    true,
				ForceNew:    true,
			},
			"vlan_vnid": {
				Type:        schema.TypeInt,
				Description: "VXLAN Network Identifier, integer",
				Required:    true,
				ForceNew:    true,
			},
			"vlan_id": {
				Type:        schema.TypeString,
				Description: "UUID of VLAN API resource",
				Computed:    true,
			},
			"port_id": {
				Type:        schema.TypeString,
				Description: "UUID of device port",
				Computed:    true,
			},
			"native": {
				Type:        schema.TypeBool,
				Description: "Mark this VLAN a native VLAN on the port. This can be used only if this assignment assigns second or further VLAN to the port. To ensure that this attachment is not first on a port, you can use depends_on pointing to another equinix_metal_port_vlan_attachment, just like in the layer2-individual example above",
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func resourceMetalPortVlanAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal
	deviceID := d.Get("device_id").(string)
	pName := d.Get("port_name").(string)
	vlanVNID := d.Get("vlan_vnid").(int)

	dev, _, err := client.Devices.Get(deviceID, &packngo.GetOptions{
		Includes: []string{"virtual_networks,project,native_virtual_network"},
	})
	if err != nil {
		return err
	}

	portFound := false
	vlanFound := false
	vlanID := ""
	var port packngo.Port
	for _, p := range dev.NetworkPorts {
		if p.Name == pName {
			portFound = true
			port = p
			for _, n := range p.AttachedVirtualNetworks {
				if vlanVNID == n.VXLAN {
					vlanFound = true
					vlanID = n.ID
					break
				}
			}
			break
		}
	}
	if !portFound {
		return fmt.Errorf("Device %s doesn't have port %s", deviceID, pName)
	}

	par := &packngo.PortAssignRequest{PortID: port.ID}
	if vlanFound {
		log.Printf("Port %s already has VLAN %d assigned", pName, vlanVNID)
		par.VirtualNetworkID = vlanID
	} else {
		projectID := dev.Project.ID
		deviceMetro := dev.Metro.Code
		deviceFacility := dev.Facility.Code
		vlans, _, err := client.ProjectVirtualNetworks.List(projectID, nil)
		if err != nil {
			return err
		}
		for _, n := range vlans.VirtualNetworks {
			// looking up vlan with given vxlan, in the same location as
			// the device - either in the same faclility or metro or both
			vlanMetro := n.MetroCode
			vlanFacility := n.FacilityCode
			if n.VXLAN == vlanVNID {
				facilitiesMatch := deviceFacility == vlanFacility
				metrosMatch := deviceMetro == vlanMetro
				if metrosMatch || facilitiesMatch {
					vlanID = n.ID
					break
				}
			}
		}
		if len(vlanID) == 0 {
			return fmt.Errorf("VLAN with VNID %d doesn't exist in project %s", vlanVNID, projectID)
		}

		par.VirtualNetworkID = vlanID

		// Equinix Metal doesn't allow multiple VLANs to be assigned
		// to the same port at the same time
		lockId := "vlan-attachment-" + port.ID
		mutexkv.Metal.Lock(lockId)
		defer mutexkv.Metal.Unlock(lockId)

		_, _, err = client.DevicePorts.Assign(par)
		if err != nil {
			return err
		}
	}

	d.SetId(port.ID + ":" + vlanID)

	native := d.Get("native").(bool)
	if native {
		_, _, err = client.DevicePorts.AssignNative(par)
		if err != nil {
			return err
		}
	}

	return resourceMetalPortVlanAttachmentRead(d, meta)
}

func resourceMetalPortVlanAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal
	deviceID := d.Get("device_id").(string)
	pName := d.Get("port_name").(string)
	vlanVNID := d.Get("vlan_vnid").(int)

	dev, _, err := client.Devices.Get(deviceID, &packngo.GetOptions{Includes: []string{"virtual_networks,project,native_virtual_network"}})
	if err != nil {
		err = equinix_errors.FriendlyError(err)

		if equinix_errors.IsNotFound(err) {
			log.Printf("[WARN] Device (%s) for Port Vlan Attachment not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}
	portFound := false
	vlanFound := false
	portID := ""
	vlanID := ""
	vlanNative := false
	for _, p := range dev.NetworkPorts {
		if p.Name == pName {
			portFound = true
			portID = p.ID
			for _, n := range p.AttachedVirtualNetworks {
				if vlanVNID == n.VXLAN {
					vlanFound = true
					vlanID = n.ID
					if p.NativeVirtualNetwork != nil {
						vlanNative = vlanID == p.NativeVirtualNetwork.ID
					}
					break
				}
			}
			break
		}
	}
	if !portFound {
		// TODO(displague) should we clear state if the port is unexpectedly
		// gone? Can we treat this like a deletion?
		return fmt.Errorf("Device %s doesn't have port %s", deviceID, pName)
	}
	if !vlanFound {
		d.SetId("")
	}
	d.Set("port_id", portID)
	d.Set("vlan_id", vlanID)
	d.Set("native", vlanNative)
	return nil
}

func resourceMetalPortVlanAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal
	if d.HasChange("native") {
		native := d.Get("native").(bool)
		portID := d.Get("port_id").(string)
		if native {
			vlanID := d.Get("vlan_id").(string)
			par := &packngo.PortAssignRequest{PortID: portID, VirtualNetworkID: vlanID}
			_, _, err := client.DevicePorts.AssignNative(par)
			if err != nil {
				return err
			}
		} else {
			_, _, err := client.DevicePorts.UnassignNative(portID)
			if err != nil {
				return err
			}
		}
	}
	return resourceMetalPortVlanAttachmentRead(d, meta)
}

func resourceMetalPortVlanAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal
	pID := d.Get("port_id").(string)
	vlanID := d.Get("vlan_id").(string)
	native := d.Get("native").(bool)
	if native {
		_, resp, err := client.DevicePorts.UnassignNative(pID)
		if equinix_errors.IgnoreResponseErrors(equinix_errors.HttpForbidden, equinix_errors.HttpNotFound)(resp, err) != nil {
			return err
		}
	}
	par := &packngo.PortAssignRequest{PortID: pID, VirtualNetworkID: vlanID}
	lockId := "vlan-detachment-" + pID
	mutexkv.Metal.Lock(lockId)
	defer mutexkv.Metal.Unlock(lockId)
	portPtr, resp, err := client.DevicePorts.Unassign(par)
	if equinix_errors.IgnoreResponseErrors(equinix_errors.HttpForbidden, equinix_errors.HttpNotFound, equinix_errors.IsNotAssigned)(resp, err) != nil {
		return err
	}
	forceBond := d.Get("force_bond").(bool)
	if forceBond && (len(portPtr.AttachedVirtualNetworks) == 0) {
		deviceID := d.Get("device_id").(string)
		portName := d.Get("port_name").(string)
		port, err := client.DevicePorts.GetPortByName(deviceID, portName)
		if err != nil {
			return equinix_errors.FriendlyError(err)
		}
		_, _, err = client.DevicePorts.Bond(port, false)
		if err != nil {
			return equinix_errors.FriendlyError(err)
		}
	}
	return nil
}
