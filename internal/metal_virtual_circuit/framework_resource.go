package metal_virtual_circuit

import (
	"context"
	"fmt"
    "time"
    "reflect"

	"github.com/equinix/terraform-provider-equinix/internal/helper"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/packethost/packngo"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func NewResource(ctx context.Context) resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "equinix_metal_virtual_circuit",
				Schema: metalVirtualCircuitResourceSchema(ctx),
			},
		),
	}
}

type Resource struct {
	helper.BaseResource
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan MetalVirtualCircuitResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve the API client from the provider metadata
    r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
    client := r.Meta.Metal

    // Check connection status
    connId := plan.ConnectionID.ValueString()
    conn, _, err := client.Connections.Get(connId, nil)
	if err != nil {
        err = helper.FriendlyError(err)
        resp.Diagnostics.AddError(
            "Error Creating Metal Virtual Circuit",
            fmt.Sprintf("Could not read Connection with ID %s: %s", connId, err),
        )
        return
    }
	if conn.Status == string(packngo.VCStatusPending) {
        resp.Diagnostics.AddError(
            "Error Creating Metal Virtual Circuit",
            fmt.Sprintf("Connection request with name %s and ID %s wasn't approved yet", conn.Name, connId),
        )
        return
	}

    // Convert the plan to an API request format
    createRequest := packngo.VCCreateRequest{
        VirtualNetworkID: plan.VlanID.ValueString(),
        Name:             plan.Name.ValueString(),
        Description:      plan.Description.ValueString(),
        Speed:            plan.Speed.ValueString(),
        VRFID:            plan.VrfID.ValueString(),
        PeerASN:          int(plan.PeerASN.ValueInt64()),
        Subnet:           plan.Subnet.ValueString(),
        MetalIP:          plan.MetalIP.ValueString(),
        CustomerIP:       plan.CustomerIP.ValueString(),
        MD5:              plan.MD5.ValueString(),
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

    if !plan.NniVLAN.IsNull() && !plan.NniVLAN.IsUnknown(){
        createRequest.NniVLAN = int(plan.NniVLAN.ValueInt64())
    }

    // API call to create the resource
    vc, _, err := client.VirtualCircuits.Create(plan.ProjectID.ValueString(), connId, plan.PortID.ValueString(), &createRequest, nil)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error creating Metal Virtual Circuit",
            fmt.Sprintf("Could not create Metal Virtual Circuit: %s", err),
        )
        return
    }

    // Wait for VC to reach target state
    // TODO: offer to wait while VCStatusPending
    targetState := string(packngo.VCStatusActive)
    createTimeout, diags := plan.Timeouts.Create(ctx, 20*time.Minute)
    if diags.HasError() {
        resp.Diagnostics.Append(diags...)
        return
    }
    createWaiter := getVCStateWaiter(
        client,
        vc.ID,
        createTimeout,
        []string{string(packngo.VCStatusActivating)},
        []string{targetState},
    )

    vcItf, err := createWaiter.WaitForStateContext(ctx)
    if err != nil {
        err = helper.FriendlyError(err)
        resp.Diagnostics.AddError(
            "Error waiting for creation of Metal Virtual Circuit",
            fmt.Sprintf("error waiting for Virtual Circuit (%s) to become %s: %s", vc.ID, targetState, err),
        )
        return
    }

    vc, ok := vcItf.(*packngo.VirtualCircuit)
    if !ok {
        resp.Diagnostics.AddError(
            "Error parsing Virtual Circuit response",
            "Unexpected response type from API",
        )
        return
    }

    // Update the Terraform state with the new resource
    var resourceState MetalVirtualCircuitResourceModel
    resourceState.parse(ctx, vc)
    diags = resp.State.Set(ctx, &resourceState)
    resp.Diagnostics.Append(diags...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state MetalVirtualCircuitResourceModel
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
    getOpts := &packngo.GetOptions{Includes: []string{"project", "virtual_network", "vrf"}}
    vc, _, err := client.VirtualCircuits.Get(id, getOpts)
    if err != nil {
        err = helper.FriendlyError(err)
        resp.Diagnostics.AddError(
            "Error reading Metal Virtual Circuit",
            fmt.Sprintf("Could not read Metal Virtual Circuit with ID %s: %s", id, err),
        )
        return
    }

    // Update the state with the current values of the resource
    resourceState := MetalVirtualCircuitResourceModel{}
    resourceState.parse(ctx, vc)
    diags = resp.State.Set(ctx, &resourceState)
    resp.Diagnostics.Append(diags...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var plan MetalVirtualCircuitResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    var state MetalVirtualCircuitResourceModel
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
    updateRequest := &packngo.VCUpdateRequest{}
    if !state.VlanID.Equal(plan.VlanID) {
        updateRequest.VirtualNetworkID = plan.VlanID.ValueStringPointer()
    }
    if !state.Name.Equal(plan.Name) {
        updateRequest.Name = plan.Name.ValueStringPointer()
    }
    if !state.Description.Equal(plan.Description) {
        updateRequest.Description = plan.Description.ValueStringPointer()
    }
    if !state.Speed.Equal(plan.Speed) {
        updateRequest.Speed = plan.Speed.ValueString()
    }
    if !state.Tags.Equal(plan.Tags) {
        tags := []string{}
        if diags := plan.Tags.ElementsAs(ctx, &tags, false); diags != nil {
            resp.Diagnostics.Append(diags...)
            return 
        }
        updateRequest.Tags = &tags
    }
    
    if !reflect.DeepEqual(updateRequest, packngo.VCUpdateRequest{}) {
        var updatedVC *packngo.VirtualCircuit
        var err error
		if updatedVC, _, err = client.VirtualCircuits.Update(id, updateRequest, nil); err != nil {
			err = helper.FriendlyError(err)
            resp.Diagnostics.AddError(
                "Error updating Metal Virtual Circuit",
                fmt.Sprintf("Could not update Virtual Circuit with ID %s: %s", id, err),
            )
            return
		}
        // Update the state with the new values of the resource
        diags = state.parse(ctx, updatedVC)
        resp.Diagnostics.Append(diags...)
	}

    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var state MetalVirtualCircuitResourceModel
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
    deleteResp, err := client.VirtualCircuits.Delete(id)
    if helper.IgnoreResponseErrors(helper.HttpForbidden, helper.HttpNotFound)(deleteResp, err) != nil {
		err = helper.FriendlyError(err)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to delete Metal Virtual Circuit %s", id),
			err.Error(),
		)
	}

    // Wait for VC to be deleted
    deleteTimeout, diags := state.Timeouts.Delete(ctx, 20*time.Minute)
    if diags.HasError() {
        resp.Diagnostics.Append(diags...)
        return
    }
    deleteTimeout = deleteTimeout - 30*time.Second
    deleteWaiter := getVCStateWaiter(
		client,
		id,
		deleteTimeout,
		[]string{string(packngo.VCStatusDeleting)},
		[]string{},
	)

    _, err = deleteWaiter.WaitForStateContext(ctx)
    if helper.IgnoreResponseErrors(helper.HttpForbidden, helper.HttpNotFound)(deleteResp, err) != nil {
		err = helper.FriendlyError(err)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to delete Metal Virtual Circuit %s", id),
            fmt.Sprintf("error waiting for Virtual Circuit (%s) to be deleted: %s", id, err),
		)
	}
}

func getVCStateWaiter(client *packngo.Client, id string, timeout time.Duration, pending, target []string) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending: pending,
		Target:  target,
		Refresh: func() (interface{}, string, error) {
			vc, _, err := client.VirtualCircuits.Get(id, nil)
			if err != nil {
				return 0, "", err
			}
			return vc, string(vc.Status), nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}
}
