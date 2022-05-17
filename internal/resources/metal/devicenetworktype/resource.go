package devicenetworktype

import (
	"log"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/device"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/packethost/packngo"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Create: resourceMetalDeviceNetworkTypeCreate,
		Read:   resourceMetalDeviceNetworkTypeRead,
		Delete: resourceMetalDeviceNetworkTypeDelete,
		Update: resourceMetalDeviceNetworkTypeUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"device_id": {
				Type:        schema.TypeString,
				Description: "The ID of the device on which the network type should be set",
				Required:    true,
				ForceNew:    true,
			},
			"type": {
				Type:         schema.TypeString,
				Description:  "Network type to set. Must be one of " + device.NetworkTypeListHB,
				Required:     true,
				ValidateFunc: validation.StringInSlice(device.DeviceNetworkTypesHB, false),
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
		_, err := c.DevicePorts.DeviceToNetworkType(devID, targetType)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceMetalDeviceNetworkTypeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*config.Config).MetalClient
	ntype := d.Get("type").(string)

	err := getAndPossiblySetNetworkType(d, client, ntype)
	if err != nil {
		return err
	}
	d.SetId(d.Get("device_id").(string))
	return resourceMetalDeviceNetworkTypeRead(d, meta)
}

func resourceMetalDeviceNetworkTypeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*config.Config).MetalClient

	_, devNType, err := getDevIDandNetworkType(d, client)
	if err != nil {
		err = errors.FriendlyError(err)

		if errors.IsNotFound(err) {
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

	d.Set("type", devNType)

	return nil
}

func resourceMetalDeviceNetworkTypeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*config.Config).MetalClient
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
