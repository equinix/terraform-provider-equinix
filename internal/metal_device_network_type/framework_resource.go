package metal_device_network_type

import (
	"context"
	"fmt"

	"github.com/equinix/terraform-provider-equinix/internal/helper"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/packethost/packngo"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "equinix_metal_device_network_type",
				Schema: &metalDeviceNetworkTypeResourceSchema,
			},
		),
	}
}

type Resource struct {
	helper.BaseResource
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan MetalDeviceNetworkTypeResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve the API client from the provider metadata
    r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
    client := r.Meta.Metal

    // Target type
    ntype := plan.Type.ValueString()

    // Making an API call to configure the resource
	device, err := getAndPossiblySetNetworkType(plan, client, ntype)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error creating Metal Device Network Type",
            fmt.Sprintf("Could not configure Metal Device Network Type for device '%s': %s", plan.DeviceID.ValueString(), err),
        )
        return
	}

    // Map the created resource data back to the Terraform state
    var resourceState MetalDeviceNetworkTypeResourceModel
    resourceState.parse(device, device.GetNetworkType())
    diags = resp.State.Set(ctx, &resourceState)
    resp.Diagnostics.Append(diags...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state MetalDeviceNetworkTypeResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve the API client from the provider metadata
    r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
    client := r.Meta.Metal

    // Extract the ID of the resource from the state
    id := state.ID.ValueString()

    // Retrieve the resource from the API
    device, err := getDevIDandNetworkType(state, client)
	if err != nil {
		err = helper.FriendlyError(err)

         // Check if the Device no longer exists
         if helper.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				"Device",
				fmt.Sprintf("[WARN] Device (%s) for Network Type request not found, removing from state", id),
			)
            resp.State.RemoveResource(ctx)
            return
        }

		resp.Diagnostics.AddError(
            "Error reading Device",
            "Could not read Device with ID " + id + ": " + err.Error(),
        )
        return
	}

    // Parse the API response into the Terraform state
    var resourceState MetalDeviceNetworkTypeResourceModel
    resourceState.parse(device, state.Type.ValueString())
    diags = resp.State.Set(ctx, &resourceState)
    resp.Diagnostics.Append(diags...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var plan MetalDeviceNetworkTypeResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    var state MetalDeviceNetworkTypeResourceModel
    diags = req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve the API client from the provider metadata
    r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
    client := r.Meta.Metal

    // Extract the ID of the organization from the state
    id := state.ID.ValueString()

    // Call API to update the resource
    var device *packngo.Device
    var err error
    if !state.Type.Equal(plan.Type) {
        device, err = getAndPossiblySetNetworkType(state, client, plan.Type.ValueString())
    } else {
        device, err = getDevIDandNetworkType(state, client)
    }

    if err != nil {
        resp.Diagnostics.AddError(
            "Error updating Metal Device Network Type",
            fmt.Sprintf("Could not configure Metal Device Network Type for device '%s': %s", id, err),
        )
        return
    }

    // Update the state with the new values of the resource
    var resourceState MetalDeviceNetworkTypeResourceModel
    resourceState.parse(device, state.Type.ValueString())
    diags = resp.State.Set(ctx, &resourceState)
    resp.Diagnostics.Append(diags...)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
   	// This resource does not support delete
}

func getAndPossiblySetNetworkType(data MetalDeviceNetworkTypeResourceModel, c *packngo.Client, targetType string) (*packngo.Device, error) {
    var device *packngo.Device
    var err error

	// "hybrid-bonded" is an alias for "layer3" with VLAN(s) connected. We use
	// other resource for VLAN attachment, so we treat these two as equivalent
	if targetType == "hybrid-bonded" {
		targetType = "layer3"
	}

	device, err = getDevIDandNetworkType(data, c)

	if err == nil && device.GetNetworkType() != targetType {
		device, err = c.DevicePorts.DeviceToNetworkType(device.ID, targetType)
	}
	return device, err
}

func getDevIDandNetworkType(data MetalDeviceNetworkTypeResourceModel, c *packngo.Client) (*packngo.Device, error) {
	deviceID := data.ID
	if deviceID.IsNull() && deviceID.IsUnknown() {
		deviceID = data.DeviceID
	}

	dev, _, err := c.Devices.Get(deviceID.ValueString(), nil)
	return dev, err
}