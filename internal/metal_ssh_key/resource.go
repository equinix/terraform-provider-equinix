package metal_ssh_key

import (
	"context"
	"path"

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
	Name        types.String `tfsdk:"name,omitempty"`
	PublicKey   types.String `tfsdk:"public_key,omitempty"`
	ProjectID   types.String `tfsdk:"project_id,omitempty"`
	Fingerprint types.String `tfsdk:"fingerprint,omitempty"`
	Updated     types.String `tfsdk:"updated,omitempty"`
	OwnerID     types.String `tfsdk:"owner_id,omitempty"`
}

func (rm *ResourceModel) parse(key *packngo.SSHKey) {
	rm.ID = types.StringValue(key.ID)
	rm.Name = types.StringValue(key.Label)
	rm.PublicKey = types.StringValue(key.Key)
	rm.ProjectID = types.StringValue(path.Base(key.Owner.Href))
	rm.Fingerprint = types.StringValue(key.FingerPrint)
	rm.Updated = types.StringValue(key.Updated)
	rm.OwnerID = types.StringValue(path.Base(key.Owner.Href))
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
    var rm ResourceModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &rm)...)
    if resp.Diagnostics.HasError() {
        return
    }

	// Generate API request body from plan
	createRequest := &packngo.SSHKeyCreateRequest{
		Label: rm.Name.ValueString(),
		Key:   rm.PublicKey.ValueString(),
	}

	if rm.ProjectID.ValueString() != "" {
		createRequest.ProjectID = rm.ProjectID.ValueString()
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

	// Set state to fully populated data
	rm.parse(key)
	resp.Diagnostics.Append(resp.State.Set(ctx, &rm)...)
}

func (r *Resource) Read(
    ctx context.Context,
    req resource.ReadRequest,
    resp *resource.ReadResponse,
) {
    // client := req.ProviderMeta.(*Config).metal
    
	// var rm ResourceModel

	// resp.Diagnostics.Append(req.State.Get(ctx, &rm)...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }

    // key, _, err := client.SSHKeys.Get(rm.ID.ValueString(), nil)
    // if err != nil {
    //     resp.Error = helper.FriendlyError(err)
    //     return
    // }

    // // Set the resource's state with the populated ResourceModel
	// rm.parse(key)
	// resp.Diagnostics.Append(resp.State.Set(ctx, &rm)...)
}


func (r *Resource) Update(
    ctx context.Context,
    req resource.UpdateRequest,
    resp *resource.UpdateResponse,
) {
    // client := r.Meta.Metal
    // id := req.ID

    // // Check if any attributes have changed
    // if req.HasChange("name") || req.HasChange("public_key") {
    //     updateRequest := &packngo.SSHKeyUpdateRequest{}

    //     if req.HasChange("name") {
    //         name := req.New.State.(*ResourceModel).Name
    //         updateRequest.Label = string(name)
    //     }

    //     if req.HasChange("public_key") {
    //         publicKey := req.New.State.(*ResourceModel).PublicKey
    //         updateRequest.Key = string(publicKey)
    //     }

    //     _, _, err := client.SSHKeys.Update(id, updateRequest)
    //     if err != nil {
    //         resp.Error = friendlyError(err)
    //         return
    //     }
    // }
}

func (r *Resource) Delete(
    ctx context.Context,
    req resource.DeleteRequest,
    resp *resource.DeleteResponse,
) {
    // client := req.Meta.Metal
    // id := req.ID

    // _, err := client.SSHKeys.Delete(id)
    // if err != nil && !isNotFound(err) {
    //     resp.Error = friendlyError(err)
    //     return
    // }

    // // Set the resource's ID to an empty string to mark it as deleted
    // resp.State = nil
}


var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
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
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
	},
}
