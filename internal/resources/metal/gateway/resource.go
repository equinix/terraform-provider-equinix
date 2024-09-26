package gateway

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

var (
	includes = []string{"project", "ip_reservation", "virtual_network", "vrf"}
)

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
		Delete: true,
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

	client := r.Meta.NewMetalClientForFramework(ctx, req.ProviderMeta)

	// Build the create request based on the plan
	// NOTE: the API spec provides 2 separate schemas for creating a
	// VRF Metal Gateway or a non-VRF Metal Gateway.  Since we can't
	// tell from resource configuration which is being requested, we
	// just use the non-VRF Metal Gateway request object.
	createRequest := metalv1.CreateMetalGatewayRequest{
		MetalGatewayCreateInput: &metalv1.MetalGatewayCreateInput{
			VirtualNetworkId: plan.VlanID.ValueString(),
		},
	}

	if reservationId := plan.IPReservationID.ValueString(); reservationId != "" {
		createRequest.MetalGatewayCreateInput.IpReservationId = &reservationId
	} else {
		// PrivateIpv4SubnetSize is specified as an int32 by the API, but
		// there is currently only an Int64 attribute defined in the plugin
		// framework.  For now we cast to int32; when Int32 attributes are
		// supported we can redefine the schema attribute to match the API
		// Reference: https://github.com/hashicorp/terraform-plugin-framework/pull/1010
		privateSubnetSize := int32(plan.PrivateIPv4SubnetSize.ValueInt64())
		createRequest.MetalGatewayCreateInput.PrivateIpv4SubnetSize = &privateSubnetSize
	}

	// Call the API to create the resource
	gw, _, err := client.MetalGatewaysApi.CreateMetalGateway(ctx, plan.ProjectID.ValueString()).CreateMetalGatewayRequest(createRequest).Include(includes).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error creating MetalGateway", err.Error())
		return
	}

	// API call to get the Metal Gateway
	gwId := ""
	if gw.MetalGateway != nil {
		gwId = gw.MetalGateway.GetId()
	} else {
		gwId = gw.VrfMetalGateway.GetId()
	}

	diags, err = getGatewayAndParse(ctx, client, &plan, gwId)
	resp.Diagnostics.Append(diags...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Metal Gateway",
			"Could not read Metal Gateway with ID "+gwId+": "+err.Error(),
		)
		return
	}

	// Update the Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.NewMetalClientForFramework(ctx, req.ProviderMeta)

	// Extract the ID of the resource from the state
	id := state.ID.ValueString()

	// API call to get the Metal Gateway
	diags, err := getGatewayAndParse(ctx, client, &state, id)
	resp.Diagnostics.Append(diags...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Metal Gateway",
			"Could not read Metal Gateway with ID "+id+": "+err.Error(),
		)
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
	client := r.Meta.NewMetalClientForFramework(ctx, req.ProviderMeta)

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
	// NOTE: we have to send `include` params on the delete request
	// because the delete request returns the gateway JSON and it will
	// only match one of MetalGateway or VrfMetalGateway if the included
	// fields are present in the response
	_, deleteResp, err := client.MetalGatewaysApi.DeleteMetalGateway(ctx, id).Include(includes).Execute()
	if equinix_errors.IgnoreHttpResponseErrors(http.StatusForbidden, http.StatusNotFound)(deleteResp, err) != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to delete Metal Gateway %s", id), err.Error(),
		)
	}

	if err == nil {
		// Wait for the deletion to be completed
		deleteTimeout := r.DeleteTimeout(ctx, state.Timeouts)
		deleteWaiter := getGatewayDeleteWaiter(
			ctx,
			client,
			id,
			deleteTimeout,
			[]string{string(metalv1.METALGATEWAYSTATE_DELETING)},
		)
		_, err = deleteWaiter.WaitForStateContext(ctx)
	}

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to delete Metal Gateway %s", id), err.Error(),
		)
	}
}

func getGatewayAndParse(ctx context.Context, client *metalv1.APIClient, state *ResourceModel, id string) (diags diag.Diagnostics, err error) {
	// API call to get the Metal Gateway
	gw, _, err := client.MetalGatewaysApi.FindMetalGatewayById(ctx, id).Include(includes).Execute()
	if err != nil {
		return diags, err
	}

	// Parse the API response into the Terraform state
	diags = state.parse(gw)
	if diags.HasError() {
		return diags, fmt.Errorf("error parsing Metal Gateway response")
	}

	return diags, nil
}

func getGatewayDeleteWaiter(ctx context.Context, client *metalv1.APIClient, id string, timeout time.Duration, pending []string) *retry.StateChangeConf {
	// deletedMarker is a terraform-provider-only value that is used by the waiter
	// to indicate that the gateway appears to be deleted successfully based on
	// status code
	deletedMarker := "tf-marker-for-deleted-gateway"

	target := []string{deletedMarker}

	return &retry.StateChangeConf{
		Pending: pending,
		Target:  target,
		Refresh: func() (interface{}, string, error) {
			gw, resp, err := client.MetalGatewaysApi.FindMetalGatewayById(ctx, id).Include(includes).Execute()
			if err != nil {
				if equinix_errors.IgnoreHttpResponseErrors(http.StatusForbidden, http.StatusNotFound)(resp, err) != nil {
					return gw, "", err
				} else {
					return gw, deletedMarker, nil
				}
			}
			state := ""
			if gw.MetalGateway != nil {
				state = string(gw.MetalGateway.GetState())
			} else {
				state = string(gw.VrfMetalGateway.GetState())
			}
			return gw, state, nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}
}
