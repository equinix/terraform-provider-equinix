// Package PortVlanAttachment provides managment of VLANs on Device Ports.
package port_vlan_attachment

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/mutexkv"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Resource provides the Terraform resource for PortVlanAttachements.
func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMetalPortVlanAttachmentCreate,
		ReadContext:   resourceMetalPortVlanAttachmentRead,
		DeleteContext: resourceMetalPortVlanAttachmentDelete,
		UpdateContext: resourceMetalPortVlanAttachmentUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
func resourceMetalPortVlanAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.Config).NewMetalClientForSDK(d)
	deviceID := d.Get("device_id").(string)
	pName := d.Get("port_name").(string)
	vlanVNID := d.Get("vlan_vnid").(int)

	dev, _, err := client.DevicesApi.FindDeviceById(ctx, deviceID).Include([]string{"virtual_networks,project,native_virtual_network"}).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	portFound := false
	vlanFound := false
	vlanID := ""
	var port metalv1.Port

	for _, p := range dev.NetworkPorts {
		if p.GetName() == pName {
			portFound = true
			port = p

			for _, n := range p.VirtualNetworks {
				if vlanVNID == int(n.GetVxlan()) {
					vlanFound = true
					vlanID = n.GetId()
					break
				}
			}
			break
		}
	}
	if !portFound {
		return diag.FromErr(fmt.Errorf("device %s doesn't have port %s", deviceID, pName))
	}

	if vlanFound {
		log.Printf("Port %s already has VLAN %d assigned", pName, vlanVNID)
	} else {
		projectID := dev.Project.GetId()
		deviceMetroCode := dev.Metro.GetCode()

		vlans, _, err := client.VLANsApi.FindVirtualNetworks(ctx, projectID).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
		for _, n := range vlans.VirtualNetworks {
			// looking up vlan with given vxlan, in the same location as
			// the device - either in the same faclility or metro or both
			vlanMetro := n.GetMetroCode()
			if int(n.GetVxlan()) == vlanVNID {
				if deviceMetroCode == vlanMetro {
					vlanID = n.GetId()
					break
				}
			}
		}
		if len(vlanID) == 0 {
			return diag.FromErr(fmt.Errorf("VLAN with VNID %d doesn't exist in project %s", vlanVNID, projectID))
		}

		// Equinix Metal doesn't allow multiple VLANs to be assigned
		// to the same port at the same time
		lockID := "vlan-attachment-" + port.GetId()
		mutexkv.Metal.Lock(lockID)
		defer mutexkv.Metal.Unlock(lockID)

		assignment := &metalv1.PortAssignInput{}
		assignment.SetVnid(vlanID)

		_, _, err = client.PortsApi.AssignPort(ctx, port.GetId()).PortAssignInput(*assignment).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(port.GetId() + ":" + vlanID)

	native := d.Get("native").(bool)
	if native {
		_, _, err = client.PortsApi.AssignNativeVlan(ctx, port.GetId()).Vnid(vlanID).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceMetalPortVlanAttachmentRead(ctx, d, meta)
}

func resourceMetalPortVlanAttachmentRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.Config).NewMetalClientForSDK(d)
	deviceID := d.Get("device_id").(string)
	pName := d.Get("port_name").(string)
	vlanVNID := d.Get("vlan_vnid").(int)

	dev, _, err := client.DevicesApi.FindDeviceById(ctx, deviceID).Include([]string{"virtual_networks,project,native_virtual_network"}).Execute()
	if err != nil {
		err = equinix_errors.FriendlyError(err)

		if equinix_errors.IsNotFound(err) {
			log.Printf("[WARN] Device (%s) for Port Vlan Attachment not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	portFound := false
	vlanFound := false
	portID := ""
	vlanID := ""
	vlanNative := false
	for _, p := range dev.NetworkPorts {
		if p.GetName() == pName {
			portFound = true
			portID = p.GetId()
			for _, n := range p.VirtualNetworks {
				if vlanVNID == int(n.GetVxlan()) {
					vlanFound = true
					vlanID = n.GetId()
					if p.NativeVirtualNetwork != nil {
						vlanNative = vlanID == p.NativeVirtualNetwork.GetId()
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
		return diag.FromErr(fmt.Errorf("device %s doesn't have port %s", deviceID, pName))
	}
	if !vlanFound {
		d.SetId("")
	}

	if err := d.Set("port_id", portID); err != nil {
		return diag.FromErr(fmt.Errorf("failed to update resource state with new port_id: %w", err))
	}

	if err := d.Set("vlan_id", vlanID); err != nil {
		return diag.FromErr(fmt.Errorf("failed to update resource state with new vlan_id: %w", err))
	}

	if err := d.Set("native", vlanNative); err != nil {
		return diag.FromErr(fmt.Errorf("failed to update resource state with native status: %w", err))
	}

	return nil
}

func resourceMetalPortVlanAttachmentUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.Config).NewMetalClientForSDK(d)
	if d.HasChange("native") {
		native := d.Get("native").(bool)
		portID := d.Get("port_id").(string)
		if native {
			vlanID := d.Get("vlan_id").(string)
			_, _, err := client.PortsApi.AssignNativeVlan(ctx, portID).Vnid(vlanID).Execute()
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			_, _, err := client.PortsApi.UnassignPort(ctx, portID).Execute()
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	return resourceMetalPortVlanAttachmentRead(ctx, d, meta)
}

func resourceMetalPortVlanAttachmentDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.Config).NewMetalClientForSDK(d)
	pID := d.Get("port_id").(string)
	vlanID := d.Get("vlan_id").(string)
	native := d.Get("native").(bool)
	if native {
		_, resp, err := client.PortsApi.DeleteNativeVlan(ctx, pID).Execute()
		switch resp.StatusCode {
		case http.StatusForbidden, http.StatusNotFound:
			// The port has disappeared or device has dissappeared, give up
		default:
			return diag.FromErr(err)
		}
	}

	lockID := "vlan-detachment-" + pID
	mutexkv.Metal.Lock(lockID)
	defer mutexkv.Metal.Unlock(lockID)

	input := metalv1.NewPortAssignInput()
	input.SetVnid(vlanID)

	portPtr, resp, err := client.PortsApi.UnassignPort(ctx, pID).PortAssignInput(*input).Execute()
	switch resp.StatusCode {
	case http.StatusForbidden, http.StatusNotFound:
		// the port or device can't be located, give up
	default:
		return diag.FromErr(err)
	}

	forceBond := d.Get("force_bond").(bool)
	if forceBond && (len(portPtr.GetVirtualNetworks()) == 0) {
		deviceID := d.Get("device_id").(string)
		portName := d.Get("port_name").(string)
		device, _, err := client.DevicesApi.FindDeviceById(ctx, deviceID).Execute()
		if err != nil {
			return diag.FromErr(equinix_errors.FriendlyError(err))
		}
		var port metalv1.Port
		portFound := false
		for _, p := range device.NetworkPorts {
			if p.GetName() == portName {
				port = p
				portFound = true
			}
		}

		if !portFound {
			return diag.FromErr(fmt.Errorf("device %s doesn't have port %s", deviceID, portName))
		}

		_, _, err = client.PortsApi.BondPort(ctx, port.GetId()).BulkEnable(false).Execute()
		if err != nil {
			return diag.FromErr(equinix_errors.FriendlyError(err))
		}
	}
	return nil
}
