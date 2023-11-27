package metal_connection

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"strconv"

	"github.com/equinix/terraform-provider-equinix/internal/helper"
    "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/packethost/packngo"
	"golang.org/x/exp/slices"
)

var (
	mega          uint64 = 1000 * 1000
	giga          uint64 = 1000 * mega
	allowedSpeeds        = []struct {
		Int uint64
		Str string
	}{
		{50 * mega, "50Mbps"},
		{200 * mega, "200Mbps"},
		{500 * mega, "500Mbps"},
		{1 * giga, "1Gbps"},
		{2 * giga, "2Gbps"},
		{5 * giga, "5Gbps"},
		{10 * giga, "10Gbps"},
	}
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "equinix_metal_connection",
				Schema: &metalConnectionResourceSchema,
			},
		),
	}
}

type Resource struct {
	helper.BaseResource
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan MetalConnectionResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Prepare the API request based on the plan
    createRequest := &packngo.ConnectionCreateRequest{
        Name:       plan.Name.ValueString(),
        Redundancy: packngo.ConnectionRedundancy(plan.Redundancy.ValueString()),
        Type:       packngo.ConnectionType(plan.Type.ValueString()),
    }

    // missing email is tolerated for user keys (can't be reasonably detected)
    if plan.ContactEmail.ValueString() != "" {
        createRequest.ContactEmail = plan.ContactEmail.ValueString()
    }

    var tokenType packngo.FabricServiceTokenType
    if plan.ServiceTokenType.ValueString() != "" {
        tokenType = packngo.FabricServiceTokenType(plan.ServiceTokenType.ValueString())
    }

    // Handle the speed setting
    if plan.Type.ValueString() == string(packngo.ConnectionDedicated) || tokenType == "a_side" {
        if plan.Speed.ValueStringPointer() == nil {
			resp.Diagnostics.AddError(
                "Error creating Metal Connection",
                "You must set speed, it's optional only for shared connections of type z_side",
            )
            return
        }
        speed, err := speedStrToUint(plan.Speed.ValueString())
        if err != nil {
            resp.Diagnostics.AddError(
                "Error creating Metal Connection",
                "Could not parse connection speed: " + err.Error(),
            )
            return
        }
        createRequest.Speed = speed
    }

    // Add tags if they are set
    if len(plan.Tags.Elements()) > 0 {
        tags := []string{}
        if diags := plan.Tags.ElementsAs(ctx, &tags, false); diags != nil {
            resp.Diagnostics.Append(diags...)
            return 
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
        resp.Diagnostics.Append(diags...)
        return 
    }
    
    // Retrieve the API client from the provider metadata
    r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
    client := r.Meta.Metal
    
    // API call to create the dedicated or shared connection
    var conn *packngo.Connection
    var err error
    projectId := plan.ProjectID.ValueString()
    if plan.Type.ValueString() == string(packngo.ConnectionShared) {
        // Shared connection specific logic
        if projectId == "" {
            resp.Diagnostics.AddError(
                "Missing project_id",
                "project_id is required for 'shared' connection type",
            )
            return
        }
        if plan.Mode.ValueString() == string(packngo.ConnectionModeTunnel) {
            resp.Diagnostics.AddError(
                "Wrong mode",
                "tunnel mode is not supported for 'shared' connections",
            )
            return
        }
        if createRequest.Redundancy == packngo.ConnectionPrimary && len(vlans) == 2 {
            resp.Diagnostics.AddError(
                "Wrong number of vlans",
                "when you create a 'shared' connection without redundancy, you must only set max 1 vlan",
            )
            return
        }
        createRequest.VLANs = vlans
        createRequest.ServiceTokenType = tokenType

        // Create the shared connection
        var err error
        conn, _, err = client.Connections.ProjectCreate(projectId, createRequest)
        if err != nil {
            resp.Diagnostics.AddError(
                "Error creating MetalConnection",
                "Could not create MetalConnection: " + err.Error(),
            )
            return
        }
    } else {
        // Dedicated connection specific logic
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
        if plan.ServiceTokenType.ValueString() != "" {
            resp.Diagnostics.AddError(
                "Failed to create Metal Connection",
                "when you create a 'dedicated' connection, you must not set service_token_type",
            )
            return
        }
        if len(vlans) > 0 {
            resp.Diagnostics.AddError(
                "Failed to create Metal Connection",
                "when you create a 'dedicated' connection, you must not set vlans",
            )
            return
        }
		createRequest.Mode = packngo.ConnectionMode(plan.Mode.ValueString())

        // Create the dedicated connection
        conn, _, err = client.Connections.OrganizationCreate(plan.OrganizationID.ValueString(), createRequest)
        if err != nil {
            resp.Diagnostics.AddError(
                "Error creating MetalConnection",
                "Could not create MetalConnection: " + err.Error(),
            )
            return
        }
    }

	// Parse API response into the Terraform state
	resp.Diagnostics.Append(plan.parse(ctx, conn)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state MetalConnectionResourceModel
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

    // Retrieve the Metal Connection from the API
    conn, _, err := client.Connections.Get(id, &packngo.GetOptions{Includes: []string{"service_tokens", "organization", "facility", "metro", "project"}})
    if err != nil {
        // If the Metal Connection is not found, remove it from the state
        if helper.IsNotFound(err) {
            resp.Diagnostics.AddWarning(
				"Metal Connection",
				fmt.Sprintf("[WARN] Connection (%s) not found, removing from state", id),
			)
            resp.State.RemoveResource(ctx)
            return
        }

        resp.Diagnostics.AddError(
            "Error reading Metal Connection",
            "Could not read Metal Connection with ID " + id + ": " + err.Error(),
        )
        return
    }

    // Update the state using the API response
    diags = state.parse(ctx, conn)
    resp.Diagnostics.Append(diags...)
    if diags.HasError() {
        return
    }

    resp.State.Set(ctx, &state)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var state, plan MetalConnectionResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    diags = req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve the API client from the provider metadata
    r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
    client := r.Meta.Metal

    // Extract the ID of the organization from the state
    id := state.ID.ValueString()

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
			err = helper.FriendlyError(err)
            resp.Diagnostics.AddError(
                "Error updating Metal Connection",
                "Could not update Connection with ID " + id + ": " + err.Error(),
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

			maxVlans := int(math.Max(float64(len(oldVlans)), float64(len(newVlans))))

            ports := make([]Port, 0, len(plan.Ports.Elements()))
            if diags := plan.Ports.ElementsAs(ctx, &ports, false); diags != nil {
                resp.Diagnostics.Append(diags...)
                return 
            }

			for i := 0; i < maxVlans; i++ {
                if oldVlans[i] != (newVlans[i]) {
					if i+1 > len(newVlans) {
						// The VNID was removed; unassign the old VNID
						if _, _, diags := updateHiddenVirtualCircuitVNID(ctx, client, ports[i], ""); diags.HasError() {
                            resp.Diagnostics.Append(diags...)
                            return
						}
					} else {
						j := slices.Index(oldVlans, newVlans[i])
						if j > i {
							// The VNID was moved to a different list index; unassign the VNID for the old index so that it is available for reassignment
							if _, _, diags := updateHiddenVirtualCircuitVNID(ctx, client, ports[j], ""); diags.HasError() {
								resp.Diagnostics.Append(diags...)
                                return
							}
						}
						// Assign the VNID (whether it is new or moved) to the correct port
						if _, _, diags := updateHiddenVirtualCircuitVNID(ctx, client, ports[i], strconv.Itoa(newVlans[i])); diags.HasError() {
                            resp.Diagnostics.Append(diags...)
                            return
						}
					}
				}
			}
		}
    } else {
        resp.Diagnostics.AddError(
            "Error updating Metal Connection",
            "Could not update Metal Connection with ID " + id + ": when you update a 'dedicated' connection, you cannot set vlans",
        )
    }

    // Retrieve the Metal Connection from the API
    conn, _, err := client.Connections.Get(id, &packngo.GetOptions{Includes: []string{"service_tokens", "organization", "facility", "metro", "project"}})
    if err != nil {
        // If the Metal Connection is not found, remove it from the state
        if helper.IsNotFound(err) {
            resp.Diagnostics.AddWarning(
				"Metal Connection",
				fmt.Sprintf("[WARN] Connection (%s) not found, removing from state", id),
			)
            resp.State.RemoveResource(ctx)
            return
        }

        resp.Diagnostics.AddError(
            "Error reading Metal Connection",
            "Could not read Metal Connection with ID " + id + ": " + err.Error(),
        )
        return
    }

    // Update the state using the API response
    diags = state.parse(ctx, conn)
    resp.Diagnostics.Append(diags...)
    if diags.HasError() {
        return
    }

    resp.State.Set(ctx, &state)
}


func updateHiddenVirtualCircuitVNID(ctx context.Context, client *packngo.Client, port Port, newVNID string) (*packngo.VirtualCircuit, *packngo.Response, diag.Diagnostics) {
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
        err = helper.FriendlyError(err)
        diags.AddError(
            "Error Updating Metal Connection",
            "Could not update Metal Connection: " + err.Error(),
        )
        return nil, nil, diags
    }
    return vc, resp, nil
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var state MetalConnectionResourceModel
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

    // API call to delete the Metal Connection
    deleteResp, err := client.Connections.Delete(id, true)
    if helper.IgnoreResponseErrors(helper.HttpForbidden, helper.HttpNotFound)(deleteResp, err) != nil {
		err = helper.FriendlyError(err)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to delete Metal Connection %s", id),
			err.Error(),
		)
	}
}
