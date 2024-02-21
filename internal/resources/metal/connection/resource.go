package connection

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/packethost/packngo"
	"golang.org/x/exp/slices"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: framework.NewBaseResource(
			framework.BaseResourceConfig{
				Name: "equinix_metal_connection",
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
	resp.Schema = resourceSchema(ctx)
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

	createRequest, diags := generateCreateRequest(ctx, plan)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Create API resource
	var conn *packngo.Connection
	var err error
	projectId := plan.ProjectID.ValueString()

	// Shared connection specific logic
	if plan.Type.ValueString() == string(packngo.ConnectionShared) {
		if projectId == "" {
			resp.Diagnostics.AddError(
				"Missing project_id",
				"project_id is required for 'shared' connection type",
			)
			return
		}

		// Create the shared connection
		var err error
		conn, _, err = client.Connections.ProjectCreate(projectId, createRequest)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating Metal Connection",
				"Could not create a shared Metal Connection: "+err.Error(),
			)
			return
		}
	}

	// Dedicated connection specific logic
	if plan.Type.ValueString() == string(packngo.ConnectionDedicated) {
		organizationId := plan.OrganizationID.ValueString()
		if organizationId == "" {
			proj, _, err := client.Projects.Get(projectId, &packngo.GetOptions{Includes: []string{"organization"}})
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Failed to get Project %s", projectId),
					err.Error(),
				)
				return
			}
			organizationId = proj.Organization.ID
		}

		// Create the dedicated connection
		conn, _, err = client.Connections.OrganizationCreate(organizationId, createRequest)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating Metal Connection",
				"Could not create Metal Connection: "+err.Error(),
			)
			return
		}
	}

	// Use API client to get the current state of the resource
	conn, diags = getConnection(ctx, client, &resp.State, conn.ID)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Parse API response into the Terraform state
	resp.Diagnostics.Append(plan.parse(ctx, conn)...)
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
	var state ResourceModel
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

	// Use API client to get the current state of the resource
	conn, diags := getConnection(ctx, client, &resp.State, id)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(state.parse(ctx, conn)...)
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

	// Prepare update request based on the changes
	updateRequest := &packngo.ConnectionUpdateRequest{}
	// TODO (ocobles) The legacy SDK code includes below code snippet which
	// looks incorrect as "locked" is a device field. Delete it when we are sure it is not necessary
	//
	// if d.HasChange("locked") {
	// 	var action func(string) (*packngo.Response, error)
	// 	if d.Get("locked").(bool) {
	// 		action = client.Devices.Lock
	// 	} else {
	// 		action = client.Devices.Unlock
	// 	}
	// 	if _, err := action(d.Id()); err != nil {
	// 		return friendlyError(err)
	// 	}
	// }
	if !state.Description.Equal(plan.Description) {
		updateRequest.Description = plan.Description.ValueStringPointer()
	}
	if !state.Mode.Equal(plan.Mode) {
		mode := packngo.ConnectionMode(plan.Mode.ValueString())
		updateRequest.Mode = &mode
	}
	if !state.Redundancy.Equal(plan.Redundancy) {
		updateRequest.Redundancy = packngo.ConnectionRedundancy(plan.Redundancy.ValueString())
	}

	// TODO(displague) packngo does not implement ContactEmail for update
	// if !state.ContactEmail.Equal(plan.ContactEmail) { ... }

	if !state.Tags.Equal(plan.Tags) {
		tags := []string{}
		if diags := plan.Tags.ElementsAs(ctx, &tags, false); diags != nil {
			resp.Diagnostics.Append(diags...)
			return
		}
		updateRequest.Tags = tags
	}

	if !reflect.DeepEqual(updateRequest, packngo.ConnectionUpdateRequest{}) {
		if _, _, err := client.Connections.Update(id, updateRequest, nil); err != nil {
			resp.Diagnostics.AddError(
				"Error updating Metal Connection",
				"Could not update Connection with ID "+id+": "+equinix_errors.FriendlyError(err).Error(),
			)
		}
	}

	// Don't update VLANs until _after_ the main ConnectionUpdateRequest has succeeded
	if !state.Vlans.Equal(plan.Vlans) {
		connType := packngo.ConnectionType(plan.Type.ValueString())

		if connType == packngo.ConnectionShared {
			oldVlans := []int{}
			if diags := state.Vlans.ElementsAs(ctx, &oldVlans, false); diags != nil {
				resp.Diagnostics.Append(diags...)
				return
			}

			newVlans := []int{}
			if diags := plan.Vlans.ElementsAs(ctx, &newVlans, false); diags != nil {
				resp.Diagnostics.Append(diags...)
				return
			}

			ports := make([]PortModel, 0, len(plan.Ports.Elements()))
			if diags := plan.Ports.ElementsAs(ctx, &ports, false); diags != nil {
				resp.Diagnostics.Append(diags...)
				return
			}

			for i, oldID := range oldVlans {
				if i < len(newVlans) {
					newID := newVlans[i]

					// If the VNIDs are different
					if oldID != newID {
						// Check if the new VNID is present elsewhere in the oldVlans list
						newIndex := slices.Index(oldVlans, newID)
						if newIndex != -1 {
							// If the new VNID is found in the oldVlans list, unassign the old VNID and assign the new VNID
							if _, _, diags := updateHiddenVirtualCircuitVNID(ctx, client, ports[i], ""); diags.HasError() {
								resp.Diagnostics.Append(diags...)
								return
							}
							if _, _, diags := updateHiddenVirtualCircuitVNID(ctx, client, ports[newIndex], strconv.Itoa(newID)); diags.HasError() {
								resp.Diagnostics.Append(diags...)
								return
							}
						} else {
							// If the new VNID is not found elsewhere in the oldVlans list, assign the new VNID
							if _, _, diags := updateHiddenVirtualCircuitVNID(ctx, client, ports[i], strconv.Itoa(newID)); diags.HasError() {
								resp.Diagnostics.Append(diags...)
								return
							}
						}
					}
				}
			}

			// If newVlans has more VNIDs than oldVlans, assign the remaining new VNIDs
			for i := len(oldVlans); i < len(newVlans); i++ {
				if _, _, diags := updateHiddenVirtualCircuitVNID(ctx, client, ports[i], strconv.Itoa(newVlans[i])); diags.HasError() {
					resp.Diagnostics.Append(diags...)
					return
				}
			}

			// If oldVlans has more VNIDs than newVlans, unassign the removed VNIDs
			for i := len(newVlans); i < len(oldVlans); i++ {
				if _, _, diags := updateHiddenVirtualCircuitVNID(ctx, client, ports[i], ""); diags.HasError() {
					resp.Diagnostics.Append(diags...)
					return
				}
			}
		} else {
			resp.Diagnostics.AddError(
				"Error updating Metal Connection",
				"Could not update Metal Connection with ID "+id+": when you update a 'dedicated' connection, you cannot set vlans",
			)
		}
	}

	// Use API client to get the current state of the resource
	conn, diags := getConnection(ctx, client, &resp.State, id)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(plan.parse(ctx, conn)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the updated state back into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Retrieve the API client
	r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
	client := r.Meta.Metal

	// Retrieve the current state
	var state ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract the ID of the resource from the state
	id := state.ID.ValueString()

	// API call to delete the Metal Connection
	deleteResp, err := client.Connections.Delete(id, true)
	if equinix_errors.IgnoreResponseErrors(equinix_errors.HttpForbidden, equinix_errors.HttpNotFound)(deleteResp, err) != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to delete Metal Connection %s", id),
			equinix_errors.FriendlyError(err).Error(),
		)
	}
}

func generateCreateRequest(ctx context.Context, plan ResourceModel) (*packngo.ConnectionCreateRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	createRequest := &packngo.ConnectionCreateRequest{
		Name:       plan.Name.ValueString(),
		Redundancy: packngo.ConnectionRedundancy(plan.Redundancy.ValueString()),
		Type:       packngo.ConnectionType(plan.Type.ValueString()),
	}
	// missing email is tolerated for user keys (can't be reasonably detected)
	if plan.ContactEmail.ValueString() != "" {
		createRequest.ContactEmail = plan.ContactEmail.ValueString()
	}
	if plan.ServiceTokenType.ValueString() != "" {
		createRequest.ServiceTokenType = packngo.FabricServiceTokenType(plan.ServiceTokenType.ValueString())
	}
	// Handle the speed setting
	// Note: missing speed is tolerated only for shared connections of type z_side
	// https://github.com/equinix/terraform-provider-equinix/issues/276
	if plan.Type.ValueString() == string(packngo.ConnectionDedicated) || plan.ServiceTokenType.ValueString() == "a_side" {
		if plan.Speed.ValueStringPointer() == nil {
			// missing speed
			diags.AddError(
				"Error creating Metal Connection",
				"You must set speed, it's optional only for shared connections of type z_side",
			)
			return nil, diags
		} else {
			speed, err := speedStrToUint(plan.Speed.ValueString())
			if err != nil {
				// wrong speed value
				diags.AddError(
					"Error creating Metal Connection",
					"Could not parse connection speed: "+err.Error(),
				)
				return nil, diags
			}
			createRequest.Speed = speed
		}
	}
	// Add tags if they are set
	if len(plan.Tags.Elements()) > 0 {
		tags := []string{}
		if diags = plan.Tags.ElementsAs(ctx, &tags, false); diags != nil {
			return nil, diags
		}
		createRequest.Tags = tags
	}
	if plan.Metro.ValueString() != "" {
		createRequest.Metro = plan.Metro.ValueString()
	}
	if plan.Facility.ValueString() != "" {
		createRequest.Facility = plan.Facility.ValueString()
	}
	if plan.Description.ValueString() != "" {
		createRequest.Description = plan.Description.ValueStringPointer()
	}
	vlans := []int{}
	if diags := plan.Vlans.ElementsAs(ctx, &vlans, false); diags != nil {
		return nil, diags
	}
	createRequest.VLANs = vlans

	// Shared connection specific logic
	if plan.Type.ValueString() == string(packngo.ConnectionShared) {
		// TODO(ocobles) The "mode" of the interconnection is only relevant to Dedicated Ports.
		// Fabric VCs won't have this field. We should consider add a default "mode" value only
		// when connection plan.Type.ValueString() == packngo.ConnectionDedicated. This validation
		// not needed.
		if packngo.ConnectionMode(plan.Mode.ValueString()) == packngo.ConnectionModeTunnel {
			// wrong mode
			diags.AddError(
				"Wrong mode",
				"tunnel mode is not supported for 'shared' connections",
			)
			return nil, diags
		}
		if createRequest.Redundancy == packngo.ConnectionPrimary && len(vlans) == 2 {
			// wrong number of vlans
			diags.AddError(
				"Wrong number of vlans",
				"when you create a 'shared' connection without redundancy, you must only set max 1 vlan",
			)
			return nil, diags
		}
	}

	// Dedicated connection specific logic
	if plan.Type.ValueString() == string(packngo.ConnectionDedicated) {
		createRequest.Mode = packngo.ConnectionMode(plan.Mode.ValueString())
		if createRequest.ServiceTokenType != "" {
			// must not set service_token_type
			diags.AddError(
				"Failed to create Metal Connection",
				"when you create a 'dedicated' connection, you must not set service_token_type",
			)
			return nil, diags
		}
		if len(createRequest.VLANs) > 0 {
			// must not set vlans
			diags.AddError(
				"Failed to create Metal Connection",
				"when you create a 'dedicated' connection, you must not set vlans",
			)
			return nil, diags
		}
	}
	return createRequest, nil
}

func getConnection(ctx context.Context, client *packngo.Client, state *tfsdk.State, id string) (*packngo.Connection, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Use API client to get the current state of the resource
	getOpts := &packngo.GetOptions{Includes: []string{"service_tokens", "organization", "facility", "metro", "project"}}
	conn, _, err := client.Connections.Get(id, getOpts)
	if err != nil {
		// If the Metal Connection is not found, remove it from the state
		if equinix_errors.IsNotFound(err) {
			diags.AddWarning(
				"Metal Connection",
				fmt.Sprintf("[WARN] Connection (%s) not found, removing from state", id),
			)
			state.RemoveResource(ctx)
			return nil, diags
		}

		diags.AddError(
			"Error reading Metal Connection",
			"Could not read Metal Connection with ID "+id+": "+err.Error(),
		)
		return nil, diags
	}
	return conn, diags
}

func updateHiddenVirtualCircuitVNID(ctx context.Context, client *packngo.Client, port PortModel, newVNID string) (*packngo.VirtualCircuit, *packngo.Response, diag.Diagnostics) {
	// This function is used to update the implicit virtual circuits attached to a shared `metal_connection` resource
	// Do not use this function for a non-shared `metal_connection`
	vcids := make([]types.String, 0, len(port.VirtualCircuitIDs.Elements()))
	diags := port.VirtualCircuitIDs.ElementsAs(ctx, &vcids, false)
	if diags.HasError() {
		return nil, nil, diags
	}
	vcid := vcids[0].ValueString()
	ucr := packngo.VCUpdateRequest{}
	ucr.VirtualNetworkID = &newVNID
	vc, resp, err := client.VirtualCircuits.Update(vcid, &ucr, nil)
	if err != nil {
		diags.AddError(
			"Error Updating Metal Connection",
			"Could not update Metal Connection: "+equinix_errors.FriendlyError(err).Error(),
		)
		return nil, nil, diags
	}
	return vc, resp, nil
}
