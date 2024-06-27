package organization

import (
	"context"
	"fmt"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	var plan ResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the API client from the provider metadata
	r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
	client := r.Meta.Metal

	// Generate API request body from plan
	addressresourcemodel, _ := plan.Address.ToSlice(ctx)
	address := packngo.Address{
		Address: addressresourcemodel[0].Address.ValueString(),
		City:    addressresourcemodel[0].City.ValueStringPointer(),
		ZipCode: addressresourcemodel[0].ZipCode.ValueString(),
		Country: addressresourcemodel[0].Country.ValueString(),
	}
	if !addressresourcemodel[0].State.IsNull() {
		address.State = addressresourcemodel[0].State.ValueStringPointer()
	}

	createRequest := &packngo.OrganizationCreateRequest{
		Name:        plan.Name.ValueString(),
		Website:     plan.Website.ValueString(),
		Description: plan.Description.ValueString(),
		Twitter:     plan.Twitter.ValueString(),
		Logo:        plan.Logo.ValueString(),
		Address:     address,
	}

	org, _, err := client.Organizations.Create(createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create Organizations",
			equinix_errors.FriendlyError(err).Error(),
		)
		return
	}

	// API call to get the Metal Organization
	diags, err = getOrganizationAndParse(ctx, client, &plan, org.ID)
	resp.Diagnostics.Append(diags...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Metal Organization",
			"Could not read Metal Organization with ID "+org.ID+": "+err.Error(),
		)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the API client from the provider metadata
	r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
	client := r.Meta.Metal

	// // Extract the ID of the resource from the state
	id := state.ID.ValueString()

	// API call to get the Metal Organization
	diags, err := getOrganizationAndParse(ctx, client, &state, id)
	resp.Diagnostics.Append(diags...)
	if err != nil {
		err = equinix_errors.FriendlyError(err)

		// If the key is somehow already destroyed, mark as
		// succesfully gone
		if equinix_errors.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				"Equinix Metal Organization not found during refresh",
				fmt.Sprintf("[WARN] Organization (%s) not found, removing from state", id),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error updating resource",
			"Could not read Metal Organization with ID "+id+": "+err.Error(),
		)
	}

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
	id := plan.ID.ValueString()

	updateRequest := &packngo.OrganizationUpdateRequest{}

	if !state.Name.Equal(plan.Name) {
		updateRequest.Name = plan.Name.ValueStringPointer()

	}
	if !state.Description.Equal(plan.Description) {
		updateRequest.Description = plan.Description.ValueStringPointer()
	}
	if !state.Website.Equal(plan.Website) {
		updateRequest.Website = plan.Website.ValueStringPointer()
	}

	if !state.Twitter.Equal(plan.Twitter) {
		updateRequest.Twitter = plan.Twitter.ValueStringPointer()
	}

	if !state.Address.Equal(plan.Address) {
		updateRequest.Address = &packngo.Address{}
		addressresourcemodel := make([]AddressResourceModel, 1)
		if diags := plan.Address.ElementsAs(ctx, &addressresourcemodel, false); diags != nil {
			resp.Diagnostics.AddError(
				"Failed to extract resource data",
				"Unable to process resource data",
			)
			return
		}

		if !addressresourcemodel[0].Address.IsNull() {
			updateRequest.Address.Address = *addressresourcemodel[0].Address.ValueStringPointer()
		}

		if !addressresourcemodel[0].City.IsNull() {
			updateRequest.Address.City = addressresourcemodel[0].City.ValueStringPointer()
		}

		if !addressresourcemodel[0].State.IsNull() {
			updateRequest.Address.State = addressresourcemodel[0].State.ValueStringPointer()
		}

		if !addressresourcemodel[0].ZipCode.IsNull() {
			updateRequest.Address.ZipCode = addressresourcemodel[0].ZipCode.ValueString()
		}

		if !addressresourcemodel[0].Country.IsNull() {
			updateRequest.Address.Country = addressresourcemodel[0].Country.ValueString()
		}
	}

	// Update the resource
	_, _, err := client.Organizations.Update(id, updateRequest)
	if err != nil {
		err = equinix_errors.FriendlyError(err)
		resp.Diagnostics.AddError(
			"Error updating resource",
			"Could not update Metal Organization with ID "+id+": "+err.Error(),
		)
		return
	}

	// API call to get the Metal Organization
	diags, err := getOrganizationAndParse(ctx, client, &plan, id)
	resp.Diagnostics.Append(diags...)
	if err != nil {
		err = equinix_errors.FriendlyError(err)
		resp.Diagnostics.AddError(
			"Error updating resource",
			"Could not read Metal Organization with ID "+id+": "+err.Error(),
		)
		return
	}

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
	id := state.ID.ValueString()

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

func getOrganizationAndParse(ctx context.Context, client *packngo.Client, state *ResourceModel, id string) (diags diag.Diagnostics, err error) {
	// API call to get the Metal Organization
	includes := &packngo.GetOptions{Includes: []string{"address"}}
	org, _, err := client.Organizations.Get(id, includes)
	if err != nil {
		return diags, equinix_errors.FriendlyError(err)
	}
	// Parse the API response into the Terraform state
	diags = state.parse(ctx, org)
	if diags.HasError() {
		return diags, fmt.Errorf("error parsing Metal Organization response")
	}

	return diags, nil
}
