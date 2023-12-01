package metal_reserved_ip_block

import (
	"context"
	"fmt"
    "time"
    "log"
    "encoding/json"

	"github.com/equinix/terraform-provider-equinix/internal/helper"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/packethost/packngo"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func NewResource(ctx context.Context) resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "equinix_metal_reserved_ip_block",
				Schema: metalReservedIpBlockResourceSchema(ctx),
			},
		),
	}
}

type Resource struct {
	helper.BaseResource
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // Initialize and get values from the plan
    var plan MetalReservedIPBlockResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve the API client from the provider metadata
    r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
    client := r.Meta.Metal

    // Prepare the request for your API using data from the plan
    createRequest := packngo.IPReservationCreateRequest{
        Type:        packngo.IPReservationType(plan.Type.ValueString()),
        Quantity:    int(plan.Quantity.ValueInt64()),
        VRFID:       plan.VrfID.ValueString(),
        Network:     plan.Network.ValueString(),
        CIDR:        int(plan.Cidr.ValueInt64()),
        Description: plan.Description.ValueString(),
        CustomData: plan.CustomData.ValueString(), // NOTE (ocobles) in legacy sdk we were checking if d.HasChange("custom_data") { req.CustomData = d.Get("custom_data")}
    }

    // Implement the conditional logic as per your requirements
    if plan.Type.ValueString() == "global_ipv4" && (plan.Facility.ValueString() != "" || plan.Metro.ValueString() != "") {
        resp.Diagnostics.AddError("Invalid Configuration", "Facility and metro can't be set for global IP block reservation")
        return
    }

    if plan.Type.ValueString() == "public_ipv4" && (plan.Facility.ValueString() == "" && plan.Metro.ValueString() == "") {
        resp.Diagnostics.AddError("Invalid Configuration", "You should set either metro or facility for non-global IP block reservation")
        return
    }

    if plan.Facility.ValueString() != "" {
        createRequest.Facility = plan.Facility.ValueStringPointer()
    }

    if plan.Metro.ValueString() != "" {
        createRequest.Metro = plan.Metro.ValueStringPointer()
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


    // Create the resource using the API
    start := time.Now()
    blockAddr, _, err := client.ProjectIPs.Create(plan.ProjectID.ValueString(), &createRequest)
    if err != nil {
        err = helper.FriendlyError(err)
        resp.Diagnostics.AddError(
            "Error creating Metal Reserved IPBlock",
            fmt.Sprintf("Could not create Metal Reserved IP Block: %s", err),
        )
        return
    }

    // Wait for IP Reservation to reach target state
    wfs := plan.WaitForState.ValueString()
    log.Printf("[DEBUG] Waiting for IP Reservation (%s) to become %s", blockAddr.ID,  wfs)
    target := []string{string(packngo.IPReservationStateCreated)}
    if wfs != string(packngo.IPReservationStateCreated) {
        target = append(target, wfs)
    }
    createTimeout, diags := plan.Timeouts.Create(ctx, 20*time.Minute)
    createTimeout = createTimeout - 30*time.Second - time.Since(start)
    createWaiter := getReservedIpBlockStateWaiter(
        client,
        blockAddr.ID,
        createTimeout,
        []string{string(packngo.IPReservationStatePending)},
        target,
    )

    reservedIPItf, err := createWaiter.WaitForStateContext(ctx)
    if err != nil {
        err = helper.FriendlyError(err)
        resp.Diagnostics.AddError(
            "Error waiting for creationg of IP Reservation",
            fmt.Sprintf("error waiting for IP Reservation (%s) to become %s: %s", blockAddr.ID, wfs, err),
        )
        return
    }

    ip, ok := reservedIPItf.(*packngo.IPAddressReservation)
    if !ok {
        resp.Diagnostics.AddError(
            "Error parsing IP Reservation response",
            "Unexpected response type from API",
        )
        return
    }

    // Map the created resource data back to the Terraform state
    var resourceState MetalReservedIPBlockResourceModel
    resourceState.parse(ctx, ip)
    diags = resp.State.Set(ctx, &resourceState)
    resp.Diagnostics.Append(diags...)
}


func getReservedIpBlockStateWaiter(client *packngo.Client, id string, timeout time.Duration, pending, target []string) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending:    pending,
		Target:     target,
		Refresh:    reservedIPStateRefreshFunc(client, id),
		Timeout:    timeout,
		MinTimeout: 5 * time.Second,
	}
}

func reservedIPStateRefreshFunc(client *packngo.Client, reservedIPId string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		reservedIP, _, err := client.ProjectIPs.Get(reservedIPId, nil)
		if err != nil {
			return nil, "", fmt.Errorf("error retrieving reserved IP block %s: %s", reservedIPId, err)
		}

		return reservedIP, string(reservedIP.State), nil
	}
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state MetalReservedIPBlockResourceModel
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
    getOpts := &packngo.GetOptions{Includes: []string{"facility", "metro", "project", "vrf"}}
	getOpts = getOpts.Filter("types", "public_ipv4,global_ipv4,private_ipv4,public_ipv6,vrf")
    reservedBlock, _, err := client.ProjectIPs.Get(id, getOpts)
    if err != nil {
        err = helper.FriendlyError(err)

        // Check if no longer exists
		if helper.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				"Metal Reserved IP Block",
				fmt.Sprintf("[WARN] IP Block (%s) not found, removing from state", id),
			)
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading Metal Reserved IP Block",
			fmt.Sprintf("Could not read Metal Reserved IP Block with ID %s: %s", id, err),
		)
		return
    }

    // Update the state with the current values of the resource
    diags = state.parse(ctx, reservedBlock)
    resp.Diagnostics.Append(diags...)
    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var plan MetalReservedIPBlockResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    var state MetalReservedIPBlockResourceModel
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
    updateRequest := &packngo.IPAddressUpdateRequest{}
    if !state.Description.Equal(plan.Description) {
        updateRequest.Description = plan.Description.ValueStringPointer()
    }
    if !state.Tags.Equal(plan.Tags) {
        tags := []string{}
        if diags := plan.Tags.ElementsAs(ctx, &tags, false); diags != nil {
            resp.Diagnostics.Append(diags...)
            return 
        }
        updateRequest.Tags = &tags
    }
    if !state.CustomData.Equal(plan.CustomData) {
        var v interface{}
		if err := json.Unmarshal([]byte(plan.CustomData.ValueString()), v); err != nil {
            diags.AddError(
                "Error updating IP Block",
                fmt.Sprintf("Error marshaling custom data to JSON: %s", err.Error()),
            )
            return
		}
		updateRequest.CustomData = v
    }

    // Call your API to update the resource
    updatedBlock, _, err := client.ProjectIPs.Update(id, updateRequest, nil)
    if err != nil {
        err = helper.FriendlyError(err)
        resp.Diagnostics.AddError(
            "Error updating Metal Reserved IP Block",
            fmt.Sprintf("Could not update Metal Reserved IP Block with ID %s: %s", id, err),
        )
        return
    }

    // Update the state with the new values of the resource
    diags = state.parse(ctx, updatedBlock)
    resp.Diagnostics.Append(diags...)
    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var state MetalReservedIPBlockResourceModel
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
    deleteResp, err := client.ProjectIPs.Remove(id)
    if helper.IgnoreResponseErrors(helper.HttpForbidden, helper.HttpNotFound)(deleteResp, err) != nil {
		err = helper.FriendlyError(err)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to delete IP Reservation block %s", id),
			err.Error(),
		)
	}
}
