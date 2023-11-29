package metal_ip_attachment

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
				Name:   "equinix_metal_ip_attachment",
				Schema: &metalIPAttachmentResourceSchema,
			},
		),
	}
}

type Resource struct {
	helper.BaseResource
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // Initialize and get values from the plan
    var plan MetalIPAttachmentResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve the API client from the provider metadata
    r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
    client := r.Meta.Metal

    // Prepare the request body
    createRequest := packngo.AddressStruct{
        Address: plan.CIDRNotation.ValueString(),
    }

    // API call to create the IP Attachment
    assignment, _, err := client.DeviceIPs.Assign(plan.DeviceID.ValueString(), &createRequest)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error assigning IP to device",
            fmt.Sprintf("Could not assign IP address %s to device %s: %s", plan.CIDRNotation.ValueString(), plan.DeviceID.ValueString(), err),
        )
        return
    }

    // Parse API response into the Terraform state
    stateDiags := plan.parse(assignment)
    resp.Diagnostics.Append(stateDiags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Set the state
    diags = resp.State.Set(ctx, &plan)
    resp.Diagnostics.Append(diags...)
}


func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    // Initialize and get current state
    var state MetalIPAttachmentResourceModel
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

    // API call to get the current state of the IP Attachment
    assignment, _, err := client.DeviceIPs.Get(id, nil)
    if err != nil {
        err = helper.FriendlyError(err)

        // If the IP Attachment is not found, mark as successfully removed
        if helper.IsNotFound(err) {
            resp.Diagnostics.AddWarning(
				"Metal IP Attachment",
				fmt.Sprintf("[WARN] IP Attachment (%s) not found, removing from state", id),
			)
            resp.State.RemoveResource(ctx)
            return
        }
        resp.Diagnostics.AddError(
            "Error reading Metal IP Attachment",
            "Could not read IP Attachment with ID " + id + ": " + err.Error(),
        )
        return
    }

    // Update the state using the API response
    diags = state.parse(assignment)
    resp.Diagnostics.Append(diags...)
    if diags.HasError() {
        return
    }

    resp.State.Set(ctx, &state)
}

func (r *Resource) Update(
    ctx context.Context,
    req resource.UpdateRequest,
    resp *resource.UpdateResponse,
) {
	// This resource does not support updates
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    // Initialize and get current state
    var state MetalIPAttachmentResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve the API client from the provider metadata
    r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
    client := r.Meta.Metal

    // Extract the ID of the organization from the state
    id := state.ID.ValueString()

    // API call to delete the IP Attachment
    deleteResp, err := client.DeviceIPs.Unassign(id)
    if helper.IgnoreResponseErrors(helper.HttpForbidden, helper.HttpNotFound)(deleteResp, err) != nil {
		err = helper.FriendlyError(err)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to delete IP Attachment %s", id),
			err.Error(),
		)
	}
}
