package metal_device_network_type

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/equinix/terraform-provider-equinix/internal/config"
)

var metalDeviceNetworkTypeResourceSchema = schema.Schema{
    Attributes: map[string]schema.Attribute{
        "id": schema.StringAttribute{
            Computed:    true,
            Description: "The unique identifier of the reserved IP block",
            PlanModifiers: []planmodifier.String{
                stringplanmodifier.UseStateForUnknown(),
            },
        },
        "device_id": schema.StringAttribute{
            Required:    true,
            Description: "The ID of the device on which the network type should be set",
        },
        "type": schema.StringAttribute{
            Required:    true,
            Description: "Network type to set. Must be one of " + config.NetworkTypeListHB,
        },
    },
}
