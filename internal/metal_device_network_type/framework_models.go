package metal_device_network_type

import (
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/packethost/packngo"
)

type MetalDeviceNetworkTypeResourceModel struct {
    ID       types.String `tfsdk:"id"`
    DeviceID types.String `tfsdk:"device_id"`
    Type     types.String `tfsdk:"type"`
}

func (rm *MetalDeviceNetworkTypeResourceModel) parse(device *packngo.Device, currentType string) {
    rm.DeviceID = types.StringValue(device.ID)
    rm.ID = rm.DeviceID

    // if "hybrid-bonded" is set as desired state and current state is "layer3",
	// keep the value in "hybrid-bonded"
    devNType := device.GetNetworkType()
	if currentType == "hybrid-bonded" && device.GetNetworkType() == "layer3" {
		devNType = "hybrid-bonded"
	}
    rm.Type = types.StringValue(devNType)
}