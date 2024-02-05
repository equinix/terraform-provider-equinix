package gateway

import (
	"context"
	"time"
    "fmt"

	"github.com/equinix/terraform-provider-equinix/internal/framework"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/packethost/packngo"
)

var _ resource.ResourceWithModifyPlan = &Resource{}

func NewResource() resource.Resource {
	r := Resource{
		BaseResource: framework.NewBaseResource(
			framework.BaseResourceConfig{
				Name: "equinix_metal_gateway",
			},
		),
	}
	r.SetDefaultDeleteTimeout(20 * time.Minute)

	return &r
}

type Resource struct {
	framework.BaseResource
	framework.WithTimeouts
}

func (r *Resource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
    // Retrieve the current state and plan
    var state, plan ResourceModel
    if !req.State.Raw.IsNull() {
        diags := req.State.Get(ctx, &state)
        resp.Diagnostics.Append(diags...)
        if resp.Diagnostics.HasError() {
            return
        }

        // If the entire plan is null, the resource is planned for destruction.
        if req.Plan.Raw.IsNull() { return }

        diags = req.Plan.Get(ctx, &plan)
        resp.Diagnostics.Append(diags...)
        if resp.Diagnostics.HasError() {
            return
        }

        // NOTE (ocobles) DiffSuppressFunc does not exist in fw
        //
        //DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
        // Suppress diff of IP reservation ID if private_ipv4_subnet_size has been set.
        // When the subnet size is set, the API will create a private subnet and return its ID
        // in this field, which generates a diff (ip_reservation_id is unset in HCL,
        // but the refreshed state shows there's an UUID of the new IPv4 block).
            // 	if d.Get("private_ipv4_subnet_size").(int) != 0 {
            // 		return true
            // 	}
            // 	return false
            // },
        if state.IPReservationID != plan.State {
            if state.PrivateIPv4SubnetSize.ValueInt64() != 0 {
                return
            }
        }
    }
}

func (r *Resource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
    s := resourceSchema(ctx)
    if s.Blocks == nil {
        s.Blocks = make(map[string]schema.Block)
    }
    s.Blocks["timeouts"] = timeouts.Block(ctx, timeouts.Opts{
        Create: true,
    })
    resp.Schema = s
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan ResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve the API client from the provider metadata
    r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
    client := r.Meta.Metal

    // Build the create request based on the plan
    createRequest := packngo.MetalGatewayCreateRequest{
        VirtualNetworkID:      plan.VlanID.ValueString(),
        IPReservationID:       plan.IPReservationID.ValueString(),
        PrivateIPv4SubnetSize: int(plan.PrivateIPv4SubnetSize.ValueInt64()),
    }

    // Call the API to create the resource
    gw, _, err := client.MetalGateways.Create(plan.ProjectID.ValueString(), &createRequest)
    if err != nil {
        resp.Diagnostics.AddError("Error creating MetalGateway", err.Error())
        return
    }

    // Parse API response into the Terraform state
	resp.Diagnostics.Append(plan.parse(gw)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}


func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

    // API call to get the Metal Gateway
    includes := &packngo.GetOptions{Includes: []string{"project", "ip_reservation", "virtual_network", "vrf"}}
    gw, _, err := client.MetalGateways.Get(id, includes)
    if err != nil {
        err = equinix_errors.FriendlyError(err)
        resp.Diagnostics.AddError(
            "Error reading Metal Gateway",
            "Could not read Metal Gateway with ID " + id + ": " + err.Error(),
        )
        return
    }

    // Parse the API response into the Terraform state
    diags = state.parse(gw)
    resp.Diagnostics.Append(diags...)
    if diags.HasError() {
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
	// This resource does not support updates
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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

    // API call to delete the Metal Gateway
    deleteResp, err := client.MetalGateways.Delete(id)
    if err == nil {    
        // Wait for the deletion to be completed
        deleteTimeout := r.DeleteTimeout(ctx, state.Timeouts)
        deleteWaiter := getGatewayStateWaiter(
            client,
            id,
            deleteTimeout,
            []string{string(packngo.MetalGatewayDeleting)},
            []string{},
        )
        _, err = deleteWaiter.WaitForStateContext(ctx)
    }

    if equinix_errors.IgnoreResponseErrors(equinix_errors.HttpForbidden, equinix_errors.HttpNotFound)(deleteResp, err) != nil {
        err = equinix_errors.FriendlyError(err)
        resp.Diagnostics.AddError(
            fmt.Sprintf("Failed to delete Metal Gateway %s", id),
            err.Error(),
        )
    }
}

func getGatewayStateWaiter(client *packngo.Client, id string, timeout time.Duration, pending, target []string) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending: pending,
		Target:  target,
		Refresh: func() (interface{}, string, error) {
			getOpts := &packngo.GetOptions{Includes: []string{"project", "ip_reservation", "virtual_network", "vrf"}}

			gw, _, err := client.MetalGateways.Get(id, getOpts) // TODO: we are not using the returned VRF. Remove the includes?
			if err != nil {
				return 0, "", err
			}
			return gw, string(gw.State), nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}
}