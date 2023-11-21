package metal_ssh_key

import (
	"context"
	"path"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/packethost/packngo"
	"github.com/equinix/terraform-provider-equinix/internal/helper"
)

type ResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	PublicKey   types.String `tfsdk:"public_key"`
	Fingerprint types.String `tfsdk:"fingerprint"`
	Created     types.String `tfsdk:"created"`
	Updated     types.String `tfsdk:"updated"`
	OwnerID     types.String `tfsdk:"owner_id"`
}

func (rm *ResourceModel) parse(key *packngo.SSHKey) diag.Diagnostics {
	rm.ID = types.StringValue(key.ID)
	rm.Name = types.StringValue(key.Label)
	rm.PublicKey = types.StringValue(key.Key)
	rm.Fingerprint = types.StringValue(key.FingerPrint)
	rm.Created = types.StringValue(key.Created)
	rm.Updated = types.StringValue(key.Updated)
	rm.OwnerID = types.StringValue(path.Base(key.Owner.Href))
	return nil
}

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "equinix_metal_ssh_key",
				Schema: &frameworkResourceSchema,
			},
		),
	}
}

type Resource struct {
	helper.BaseResource
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// r.Meta.AddModuleToMetalUserAgent(d)
	client := r.Meta.Metal

	// Retrieve values from plan
	var plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	createRequest := &packngo.SSHKeyCreateRequest{
		Label: plan.Name.ValueString(),
		Key:   plan.PublicKey.ValueString(),
	}

	// Create API resource
	key, _, err := client.SSHKeys.Create(createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create SSH Key",
			helper.FriendlyError(err).Error(),
		)
		return
	}

	// Parse API response into the Terraform state
	resp.Diagnostics.Append(plan.parse(key)...)
	if resp.Diagnostics.HasError() {
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
	// r.Meta.AddModuleToMetalUserAgent(d)
	client := r.Meta.Metal

	// Retrieve values from plan
	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract the ID of the resource from the state
    id := state.ID.ValueString()

    // Use API client to get the current state of the resource
	key, _, err := client.SSHKeys.Get(id, nil)
	if err != nil {
		err = helper.FriendlyError(err)

		// If the key is somehow already destroyed, mark as
		// succesfully gone
		if helper.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				"SSHKey",
				fmt.Sprintf("[WARN] SSHKey (%s) not found, removing from state", id),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get SSHKey %s", id),
			err.Error(),
		)
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(state.parse(key)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}


func (r *Resource) Update(
    ctx context.Context,
    req resource.UpdateRequest,
    resp *resource.UpdateResponse,
) {
	// r.Meta.AddModuleToMetalUserAgent(d)
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

	updateRequest := &packngo.SSHKeyUpdateRequest{}
	if !state.Name.Equal(plan.Name) {
		updateRequest.Label = plan.Name.ValueStringPointer()
	}
	if !state.PublicKey.Equal(plan.PublicKey) {
		updateRequest.Key = plan.PublicKey.ValueStringPointer()
	}

	// Use your API client to update the resource
	key, _, err := client.SSHKeys.Update(plan.ID.ValueString(), updateRequest)
	if err != nil {
		err = helper.FriendlyError(err)
		resp.Diagnostics.AddError(
			"Error updating resource",
			"Could not update resource with ID " + id + ": " + err.Error(),
		)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(plan.parse(key)...)
	if resp.Diagnostics.HasError() {
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
	// r.Meta.AddModuleToMetalUserAgent(d)
	client := r.Meta.Metal

	// Retrieve values from plan
	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract the ID of the resource from the state
	id := state.ID.ValueString()

	// Use your API client to delete the resource
	deleteResp, err := client.SSHKeys.Delete(id)
	if helper.IgnoreResponseErrors(helper.HttpForbidden, helper.HttpNotFound)(deleteResp, err) != nil {
		err = helper.FriendlyError(err)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to delete SSHKey %s", id),
			err.Error(),
		)
	}
}

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique identifier for this SSH key.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"name": schema.StringAttribute{
			Description: "The name of the SSH key for identification",
			Required:    true,
		},
		"public_key": schema.StringAttribute{
			Description: "The public key",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"fingerprint": schema.StringAttribute{
			Description: "The fingerprint of the SSH key",
			Computed:    true,
		},
		"owner_id": schema.StringAttribute{
			Description: "The UUID of the Equinix Metal API User who owns this key",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"created": schema.StringAttribute{
			Description: "The timestamp for when the SSH key was created",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"updated": schema.StringAttribute{
			Description: "The timestamp for the last time the SSH key was updated",
			Computed:    true,
		},
	},
}
