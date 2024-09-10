package equinix

import (
	"log"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/network"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/packethost/packngo"
)

func resourceMetalDeviceNetworkType() *schema.Resource {
	return &schema.Resource{
		Create: resourceMetalDeviceNetworkTypeCreate,
		Read:   resourceMetalDeviceNetworkTypeRead,
		Delete: resourceMetalDeviceNetworkTypeDelete,
		Update: resourceMetalDeviceNetworkTypeUpdate,
		Importer: &schema.ResourceImporter{
			//nolint
			State: schema.ImportStatePassthrough,
		},
		DeprecationMessage: "The metal_device_network_type resource is deprecated.  Please use metal_port instead.  See the [Metal Device Network Types guide](https://registry.terraform.io/providers/equinix/equinix/latest/docs/guides/network_types) for more info",
		Description: `This resource controls network type of Equinix Metal devices.

To learn more about Layer 2 networking in Equinix Metal, refer to

* https://metal.equinix.com/developers/docs/networking/layer2/
* https://metal.equinix.com/developers/docs/networking/layer2-configs/

If you are attaching VLAN to a device (i.e. using equinix_metal_port_vlan_attachment), link the device ID from this resource, in order to make the port attachment implicitly dependent on the state of the network type. If you link the device ID from the equinix_metal_device resource, Terraform will not wait for the network type change. See examples in [equinix_metal_port_vlan_attachment](port_vlan_attachment).

-> **NOTE:** This resource takes a named network type with any mode required parameters and converts a device to the named network type. This resource simulated the network type interface for Devices in the Equinix Metal Portal. That interface changed when additional network types were introduced with more diverse port configurations and it is not guaranteed to work in devices with more than two ethernet ports. See the [Network Types Guide](../guides/network_types.md) for examples of this resource and to learn about the recommended ` + "`equinix_metal_port`" + ` alternative.`,
		Schema: map[string]*schema.Schema{
			"device_id": {
				Type:        schema.TypeString,
				Description: "The ID of the device on which the network type should be set",
				Required:    true,
				ForceNew:    true,
			},
			"type": {
				Type:         schema.TypeString,
				Description:  "Network type to set. Must be one of " + network.NetworkTypeListHB,
				Required:     true,
				ValidateFunc: validation.StringInSlice(network.DeviceNetworkTypesHB, false),
			},
		},
	}
}

func getDevIDandNetworkType(d *schema.ResourceData, c *packngo.Client) (string, string, error) {
	deviceID := d.Id()
	if len(deviceID) == 0 {
		deviceID = d.Get("device_id").(string)
	}

	dev, _, err := c.Devices.Get(deviceID, nil)
	if err != nil {
		return "", "", err
	}
	devType := dev.GetNetworkType()

	return dev.ID, devType, nil
}

func getAndPossiblySetNetworkType(d *schema.ResourceData, c *packngo.Client, targetType string) error {
	// "hybrid-bonded" is an alias for "layer3" with VLAN(s) connected. We use
	// other resource for VLAN attachment, so we treat these two as equivalent
	if targetType == "hybrid-bonded" {
		targetType = "layer3"
	}
	devID, devType, err := getDevIDandNetworkType(d, c)
	if err != nil {
		return err
	}

	if devType != targetType {
		//nolint
		_, err := c.DevicePorts.DeviceToNetworkType(devID, targetType)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceMetalDeviceNetworkTypeCreate(d *schema.ResourceData, meta interface{}) error {
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal

	ntype := d.Get("type").(string)
	err := getAndPossiblySetNetworkType(d, client, ntype)
	if err != nil {
		return err
	}
	d.SetId(d.Get("device_id").(string))
	return resourceMetalDeviceNetworkTypeRead(d, meta)
}

func resourceMetalDeviceNetworkTypeRead(d *schema.ResourceData, meta interface{}) error {
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal

	_, devNType, err := getDevIDandNetworkType(d, client)
	if err != nil {
		err = equinix_errors.FriendlyError(err)

		if equinix_errors.IsNotFound(err) {
			log.Printf("[WARN] Device (%s) for Network Type request not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	// if "hybrid-bonded" is set as desired state and current state is "layer3",
	// keep the value in "hybrid-bonded"
	currentType := d.Get("type").(string)
	if currentType == "hybrid-bonded" && devNType == "layer3" {
		devNType = "hybrid-bonded"
	}

	err = d.Set("type", devNType)

	return err
}

func resourceMetalDeviceNetworkTypeUpdate(d *schema.ResourceData, meta interface{}) error {
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal

	ntype := d.Get("type").(string)
	if d.HasChange("type") {
		err := getAndPossiblySetNetworkType(d, client, ntype)
		if err != nil {
			return err
		}
	}
	return resourceMetalDeviceNetworkTypeRead(d, meta)
}

func resourceMetalDeviceNetworkTypeDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
