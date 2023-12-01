package metal_vrf

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
				Name:   "equinix_metal_vrf",
				Schema: &metalVrfResourceSchema,
			},
		),
	}
}

type Resource struct {
	helper.BaseResource
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // Initialize and get values from the plan
    var plan MetalVRFResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve the API client from the provider metadata
    r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
    client := r.Meta.Metal

    // Convert the plan to an API request
    createRequest := &packngo.VRFCreateRequest{
        Name:        plan.Name.ValueString(),
        Description: plan.Description.ValueString(),
        Metro:       plan.Metro.ValueString(),
        LocalASN:    int(plan.LocalASN.ValueInt64()),
    }

    ipRanges := []string{}
    if diags := plan.IPRanges.ElementsAs(ctx, &ipRanges, false); diags != nil {
        resp.Diagnostics.Append(diags...)
        return 
    }
    createRequest.IPRanges = ipRanges

    // API call to create the resource
    vrf, _, err := client.VRFs.Create(plan.ProjectID.ValueString(), createRequest)
    if err != nil {
        err = helper.FriendlyError(err)
        resp.Diagnostics.AddError(
            "Error creating Metal VRF",
            fmt.Sprintf("Could not create Metal VRF: %s", err),
        )
        return
    }

    // Update the Terraform state with the new resource
    var resourceState MetalVRFResourceModel
    resourceState.parse(ctx, vrf)
    diags = resp.State.Set(ctx, &resourceState)
    resp.Diagnostics.Append(diags...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state MetalVRFResourceModel
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
    vrf, _, err := client.VRFs.Get(id, &packngo.GetOptions{})
    if err != nil {
        err = helper.FriendlyError(err)
        // If the VRF was destroyed, mark as gone
		if helper.IsNotFound(err) || helper.IsForbidden(err) {
			resp.Diagnostics.AddWarning(
				"Metal Metal",
				fmt.Sprintf("[WARN] VRF (%s) not accessible, removing from state", id),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading Metal VRF",
			fmt.Sprintf("Could not read Metal VRF with ID %s: %s", id, err),
		)
		return
       
    }

    // Update the state with the current values of the resource
    diags = state.parse(ctx, vrf)
    resp.Diagnostics.Append(diags...)
    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var plan MetalVRFResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    var state MetalVRFResourceModel
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

    // Prepare the update request
    updateRequest := &packngo.VRFUpdateRequest{}
    if !state.Name.Equal(plan.Name) {
        updateRequest.Name = plan.Name.ValueStringPointer()
    }
    if !state.Description.Equal(plan.Description) {
        updateRequest.Description = plan.Description.ValueStringPointer()
    }
    if !state.LocalASN.Equal(plan.LocalASN) {
        asn := int(plan.LocalASN.ValueInt64())
        updateRequest.LocalASN = &asn
    }
    if !state.IPRanges.Equal(plan.IPRanges) {
        ranges := []string{}
        if diags := plan.IPRanges.ElementsAs(ctx, &ranges, false); diags != nil {
            resp.Diagnostics.Append(diags...)
            return 
        }
        updateRequest.IPRanges = &ranges
    }

    // Call your API to update the resource
    updatedVrf, _, err := client.VRFs.Update(id, updateRequest)
    if err != nil {
        err = helper.FriendlyError(err)
        resp.Diagnostics.AddError(
            "Error updating Metal VRF",
            fmt.Sprintf("Could not update Metal VRF with ID %s: %s", id, err),
        )
        return
    }

    // Update the state with the new values of the resource
    diags = state.parse(ctx, updatedVrf)
    resp.Diagnostics.Append(diags...)
    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var state MetalVRFResourceModel
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

    // Call your API to delete the resource
    deleteResp, err := client.VRFs.Delete(id)
    if helper.IgnoreResponseErrors(helper.HttpForbidden, helper.HttpNotFound)(deleteResp, err) != nil {
        err = helper.FriendlyError(err)
        resp.Diagnostics.AddError(
            "Error deleting Metal VRF",
            fmt.Sprintf("Could not delete Metal VRF with ID %s: %s", id, err),
        )
        return
    }
}
