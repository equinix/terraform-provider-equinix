package metal_gateway
import (
	"context"
	"time"
    "fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/packethost/packngo"
	"github.com/equinix/terraform-provider-equinix/internal/helper"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)


func NewResource(ctx context.Context) resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "equinix_metal_gateway",
				Schema: metalGatewayResourceSchema(ctx),
			},
		),
	}
}

type Resource struct {
	helper.BaseResource
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan MetalGatewayResourceModel
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
    result, _, err := client.MetalGateways.Create(plan.ProjectID.ValueString(), &createRequest)
    if err != nil {
        resp.Diagnostics.AddError("Error creating MetalGateway", err.Error())
        return
    }

    // Update the Terraform state with the new ID
    diags = resp.State.Set(ctx, &MetalGatewayResourceModel{
        ID: types.StringValue(result.ID),
        // Set other fields as necessary
    })
    resp.Diagnostics.Append(diags...)
}


func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state MetalGatewayResourceModel
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
    mg, _, err := client.MetalGateways.Get(id, includes)
    if err != nil {
        err = helper.FriendlyError(err)
        resp.Diagnostics.AddError(
            "Error reading Metal Gateway",
            "Could not read Metal Gateway with ID " + id + ": " + err.Error(),
        )
        return
    }

    // Parse the API response into the Terraform state
    diags = state.parse(ctx, mg)
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
    var state MetalGatewayResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Extract the ID of the resource from the state
    id := state.ID.ValueString()

    // API call to delete the Metal Gateway
    deleteResp, err := client.MetalGateways.Delete(id)
    if helper.IgnoreResponseErrors(helper.HttpForbidden, helper.HttpNotFound)(deleteResp, err) != nil {
		err = helper.FriendlyError(err)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to delete Metal Gateway %s", id),
			err.Error(),
		)
	}

    // Wait for the deletion to be completed
    //
    // NOTE (ocobles) WaitForStateorRetryContext doesn't exist in terraform framework
    // using sdk library https://discuss.hashicorp.com/t/terraform-plugin-framework-what-is-the-replacement-for-waitforstate-or-retrycontext/45538
    deleteTimeout, diags := state.Timeouts.Delete(ctx, 20*time.Minute)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
    deleteWaiter := getGatewayStateWaiter(
        client,
        id,
        deleteTimeout,
        []string{string(packngo.MetalGatewayDeleting)},
        []string{},
    )

    _, err = deleteWaiter.WaitForStateContext(ctx)
    if err != nil {
        err = helper.FriendlyError(err)
        resp.Diagnostics.AddError(
            "Error waiting for deletion of MetalGateway",
            "Failed to delete MetalGateway with ID " + id + ": " + err.Error(),
        )
        return
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
