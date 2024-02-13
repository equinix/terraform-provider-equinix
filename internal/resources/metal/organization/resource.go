package organization

import (
	"context"
	"fmt"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/packethost/packngo"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: framework.NewBaseResource(
			framework.BaseResourceConfig{
				Name: "equinix_metal_organization",
			},
		),
	}
}

type Resource struct {
	framework.BaseResource
}

func (r *Resource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = GetResourceSchema(ctx)
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
	client := r.Meta.Metal

	// Retrieve values from plan
	var plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	addressresourcemodel := make([]AddressResourceModel, 1)
	if diags := plan.address.ElementsAs(ctx, &addressresourcemodel, false); diags != nil {
		resp.Diagnostics.AddError(
			"Failed to extract resource data",
			"Unable to process resource data",
		)
		return
	}
	address := packngo.Address{}
	if !addressresourcemodel[0].address.IsNull() {
		address.Address = addressresourcemodel[0].address.ValueString()
	}

	if !addressresourcemodel[0].city.IsNull() {
		address.City = addressresourcemodel[0].city.ValueStringPointer()
	}

	if !addressresourcemodel[0].state.IsNull() {
		address.State = addressresourcemodel[0].state.ValueStringPointer()
	}

	if !addressresourcemodel[0].zipCode.IsNull() {
		address.ZipCode = addressresourcemodel[0].zipCode.ValueString()
	}

	if !addressresourcemodel[0].country.IsNull() {
		address.Country = addressresourcemodel[0].country.ValueString()
	}
	// Generate API request body from plan
	createRequest := &packngo.OrganizationCreateRequest{
		Name:    plan.name.ValueString(),
		Address: address,
	}

	createRequest.Website = plan.website.ValueString()
	createRequest.Description = plan.description.ValueString()
	createRequest.Twitter = plan.twitter.ValueString()
	createRequest.Logo = plan.logo.ValueString()

	org, _, err := client.Organizations.Create(createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create Organizations",
			equinix_errors.FriendlyError(err).Error(),
		)
		return
	}

	// Parse API response into the Terraform state
	plan.parse(ctx, org)

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Retrieve the API client from the provider metadata
	r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
	client := r.Meta.Metal

	// Retrieve values from state
	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// // Extract the ID of the resource from the state
	id := state.id.ValueString()

	org, _, err := client.Organizations.Get(id, &packngo.GetOptions{Includes: []string{"address"}})
	if err != nil {
		err = equinix_errors.FriendlyError(err)

		// If the key is somehow already destroyed, mark as
		// succesfully gone
		if equinix_errors.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				"Equinix Metal Organizations not found during refresh",
				fmt.Sprintf("[WARN] Organization (%s) not found, removing from state", id),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get Organization %s", id),
			err.Error(),
		)
	}

	// Set state to fully populated data
	state.parse(ctx, org)

	// Update the Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
	client := r.Meta.Metal

	// Retrieve values from plan
	var state, plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract the ID of the resource from the state
	id := plan.id.ValueString()

	updateRequest := &packngo.OrganizationUpdateRequest{}
	if !state.name.Equal(plan.name) {
		updateRequest.Name = plan.name.ValueStringPointer()

	}
	if !state.description.Equal(plan.description) {
		updateRequest.Description = plan.description.ValueStringPointer()
	}
	if !state.website.Equal(plan.website) {
		updateRequest.Website = plan.website.ValueStringPointer()
	}

	if !state.twitter.Equal(plan.twitter) {
		updateRequest.Twitter = plan.twitter.ValueStringPointer()
	}

	if !state.address.Equal(plan.address) {

		addressresourcemodel := make([]AddressResourceModel, 1)
		if diags := plan.address.ElementsAs(ctx, &addressresourcemodel, false); diags != nil {
			resp.Diagnostics.AddError(
				"Failed to extract resource data",
				"Unable to process resource data",
			)
			return
		}

		updateRequest.Address.Address = addressresourcemodel[0].address.ValueString()
		updateRequest.Address.City = addressresourcemodel[0].city.ValueStringPointer()
		updateRequest.Address.State = addressresourcemodel[0].state.ValueStringPointer()
		updateRequest.Address.ZipCode = addressresourcemodel[0].zipCode.ValueString()
		updateRequest.Address.Country = addressresourcemodel[0].country.ValueString()
	}

	// Update the resource
	org, _, err := client.Organizations.Update(id, updateRequest)
	if err != nil {
		err = equinix_errors.FriendlyError(err)
		resp.Diagnostics.AddError(
			"Error updating resource",
			"Could not update resource with ID "+id+": "+err.Error(),
		)
		return
	}

	// Set state to fully populated data
	plan.parse(ctx, org)

	// Read the updated state back into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
	client := r.Meta.Metal

	// Retrieve values from plan
	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract the ID of the resource from the state
	id := state.id.ValueString()

	// Use API client to delete the resource
	deleteResp, err := client.Organizations.Delete(id)
	if equinix_errors.IgnoreResponseErrors(equinix_errors.HttpForbidden, equinix_errors.HttpNotFound)(deleteResp, err) != nil {
		err = equinix_errors.FriendlyError(err)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to delete Organizations %s", id),
			err.Error(),
		)
	}
}
