package vlan

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/packethost/packngo"
)

var (
	vlanDefaultIncludes = []string{"assigned_to", "facility", "metro"}
)

// Resource defines the Terraform resource implementation for managing VLANs.
// It embeds framework.BaseResource to inherit core resource behavior and
// framework.WithTimeouts to support customizable operation timeouts.
type Resource struct {
	framework.BaseResource
	framework.WithTimeouts
}

// NewResource creates and returns a new instance of the VLAN Terraform resource.
// It initializes the resource with a base configuration, including its name,
// and returns it as a framework-compatible resource.Resource interface
func NewResource() resource.Resource {
	r := Resource{
		BaseResource: framework.NewBaseResource(
			framework.BaseResourceConfig{
				Name: "equinix_metal_vlan",
			},
		),
	}

	return &r
}

// Schema defines the Terraform schema for the equinix_metal_vlan resource.
// It retrieves the base schema using resourceSchema, ensures the Blocks map is initialized,
// and assigns the resulting schema to the response. This method is called by the Terraform
// framework during provider initialization to understand the structure of the resource.
func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	s := resourceSchema(ctx)
	if s.Blocks == nil {
		s.Blocks = make(map[string]schema.Block)
	}

	resp.Schema = s
}

// Create provisions a new VLAN resource in Equinix Metal using the Terraform framework.
// It performs the following steps:
//    1. Adds the Terraform framework module to the Equinix Metal user agent.
//    2. Parses and validates the input configuration from the Terraform plan.
//    3. Validates that either a facility or metro is specified, and that VXLAN is only set for metro VLANs.
//    4. Constructs a VirtualNetworkCreateRequest and sends it to the Equinix Metal API.
//    5. Retrieves the newly created VLAN with default include fields to ensure full state population.
//    6. Parses the API response into the Terraform state model.
//    7. Sets the final state for Terraform to track.
//
// Any errors encountered during these steps are added to the diagnostics response to inform the user.
func (r *Resource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	r.Meta.AddFwModuleToMetalUserAgent(ctx, request.ProviderMeta)
	client := r.Meta.Metal

	var data ResourceModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	if data.Facility.IsNull() && data.Metro.IsNull() {
		response.Diagnostics.AddError("Invalid input params",
			equinix_errors.FriendlyError(errors.New("one of facility or metro must be configured")).Error())
		return
	}
	if !data.Facility.IsNull() && !data.Vxlan.IsNull() {
		response.Diagnostics.AddError("Invalid input params",
			equinix_errors.FriendlyError(errors.New("you can set vxlan only for metro vlan")).Error())
		return
	}

	createRequest := &packngo.VirtualNetworkCreateRequest{
		ProjectID:   data.ProjectID.ValueString(),
		Description: data.Description.ValueString(),
	}
	if !data.Metro.IsNull() {
		createRequest.Metro = strings.ToLower(data.Metro.ValueString())
		createRequest.VXLAN = int(data.Vxlan.ValueInt64())
	}
	if !data.Facility.IsNull() {
		createRequest.Facility = data.Facility.ValueString()
	}
	vlan, _, err := client.ProjectVirtualNetworks.Create(createRequest)
	if err != nil {
		response.Diagnostics.AddError("Error creating Vlan", equinix_errors.FriendlyError(err).Error())
		return
	}

	// get the current state of newly created vlan with default include fields
	vlan, _, err = client.ProjectVirtualNetworks.Get(vlan.ID, &packngo.GetOptions{Includes: vlanDefaultIncludes})
	if err != nil {
		response.Diagnostics.AddError("Error reading Vlan after create", equinix_errors.FriendlyError(err).Error())
		return
	}

	// Parse API response into the Terraform state
	response.Diagnostics.Append(data.parse(vlan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *Resource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	r.Meta.AddFwModuleToMetalUserAgent(ctx, request.ProviderMeta)
	client := r.Meta.Metal

	var data ResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	vlan, _, err := client.ProjectVirtualNetworks.Get(
		data.ID.ValueString(),
		&packngo.GetOptions{Includes: vlanDefaultIncludes},
	)
	if err != nil {
		if equinix_errors.IsNotFound(err) {
			response.Diagnostics.AddWarning(
				"Equinix Metal Vlan not found during refresh",
				fmt.Sprintf("[WARN] Vlan (%s) not found, removing from state", data.ID.ValueString()),
			)
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError("Error fetching Vlan using vlanId",
			equinix_errors.FriendlyError(err).Error())
		return
	}

	response.Diagnostics.Append(data.parse(vlan)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

// Update modifies an existing VLAN resource in Equinix Metal based on the desired state
// provided in the Terraform plan. It performs the following steps:
//    1. Initializes a new Metal client using the provider metadata.
//    2. Retrieves the current and planned state from the Terraform request.
//    3. Compares relevant fields (currently only Description) and constructs an update request.
//    4. Sends the update request to the Equinix Metal API.
//    5. Parses the updated VLAN response into the Terraform state model.
//    6. Updates the Terraform state with the new data.
//
// Any errors encountered during the update process are added to the diagnostics response.
func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	client := r.Meta.NewMetalClientForFramework(ctx, req.ProviderMeta)

	// Retrieve values from plan
	var state, plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract the ID of the resource from the state
	id := plan.ID.ValueString()

	updateRequest := &metalv1.VirtualNetworkUpdateInput{}
	if !state.Description.Equal(plan.Description) {
		updateRequest.Description = plan.Description.ValueStringPointer()
	}

	// Update the resource
	vlan, _, err := client.VLANsApi.UpdateVirtualNetwork(ctx, id).VirtualNetworkUpdateInput(*updateRequest).Include([]string{"assigned_to"}).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating resource",
			"Could not update resource with ID "+id+": "+err.Error(),
		)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(plan.parseMetalV1(vlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the updated state back into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete removes a VLAN resource from Equinix Metal using the Terraform framework.
// It performs the following steps:
//    1. Adds the Terraform framework module to the Equinix Metal user agent.
//    2. Retrieves the current state of the resource from Terraform.
//    3. Fetches the VLAN from the Equinix Metal API, including related instances and ports.
//    4. Iterates through all attached devices and unassigns the VLAN from their ports.
//    5. Deletes the VLAN from the Equinix Metal project.
//
// If the VLAN is not found or access is forbidden, the method exits gracefully with a warning.
// Any other errors encountered during unassignment or deletion are added to the diagnostics response.
func (r *Resource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	r.Meta.AddFwModuleToMetalUserAgent(ctx, request.ProviderMeta)
	client := r.Meta.Metal

	var data ResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	vlan, resp, err := client.ProjectVirtualNetworks.Get(
		data.ID.ValueString(),
		&packngo.GetOptions{Includes: []string{"instances", "virtual_networks", "meta_gateway"}},
	)
	if err != nil {
		if equinix_errors.IgnoreResponseErrors(equinix_errors.HttpForbidden, equinix_errors.HttpNotFound)(resp, err) != nil {
			response.Diagnostics.AddWarning(
				"Equinix Metal Vlan not found during delete",
				equinix_errors.FriendlyError(err).Error(),
			)
			return
		}
		response.Diagnostics.AddError("Error fetching Vlan using vlanId",
			equinix_errors.FriendlyError(err).Error())
		return
	}

	// all device ports must be unassigned before delete
	for _, instance := range vlan.Instances {
		for _, port := range instance.NetworkPorts {
			for _, v := range port.AttachedVirtualNetworks {
				if path.Base(v.Href) == vlan.ID {
					_, resp, err = client.Ports.Unassign(port.ID, vlan.ID)
					if equinix_errors.IgnoreResponseErrors(equinix_errors.HttpForbidden, equinix_errors.HttpNotFound)(resp, err) != nil {
						response.Diagnostics.AddError("Error unassign port with Vlan",
							equinix_errors.FriendlyError(err).Error())
						return
					}
				}
			}
		}
	}

	if err := equinix_errors.IgnoreResponseErrors(equinix_errors.HttpForbidden, equinix_errors.HttpNotFound)(client.ProjectVirtualNetworks.Delete(vlan.ID)); err != nil {
		response.Diagnostics.AddError("Error deleting Vlan",
			equinix_errors.FriendlyError(err).Error())
		return
	}
}
