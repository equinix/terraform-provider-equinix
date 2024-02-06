package gateway

import (
	"context"
	"fmt"
	"time"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/packethost/packngo"
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

	// API call to get the Metal Gateway
	diags, err = getGatewayAndParse(client, &plan, gw.ID)
	resp.Diagnostics.Append(diags...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Metal Gateway",
			"Could not read Metal Gateway with ID "+gw.ID+": "+err.Error(),
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

	// Retrieve the API client from the provider metadata
	r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
	client := r.Meta.Metal

	// Extract the ID of the resource from the state
	id := state.ID.ValueString()

	// API call to get the Metal Gateway
	diags, err := getGatewayAndParse(client, &state, id)
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
		deleteResp = nil
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
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to delete Metal Gateway %s", id),
			equinix_errors.FriendlyError(err).Error(),
		)
	}
}

func getGatewayAndParse(client *packngo.Client, state *ResourceModel, id string) (diags diag.Diagnostics, err error) {
	// API call to get the Metal Gateway
	includes := &packngo.GetOptions{Includes: []string{"project", "ip_reservation", "virtual_network", "vrf"}}
	gw, _, err := client.MetalGateways.Get(id, includes)
	if err != nil {
		return diags, equinix_errors.FriendlyError(err)
	}

	// Parse the API response into the Terraform state
	diags = state.parse(gw)
	if diags.HasError() {
		return diags, fmt.Errorf("error parsing Metal Gateway response")
	}

	return diags, nil
}

func getGatewayStateWaiter(client *packngo.Client, id string, timeout time.Duration, pending, target []string) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending: pending,
		Target:  target,
		Refresh: func() (interface{}, string, error) {
			getOpts := &packngo.GetOptions{Includes: []string{"project", "ip_reservation", "virtual_network", "vrf"}}

			gw, _, err := client.MetalGateways.Get(id, getOpts) // TODO: we are not using the returned gw. Remove the includes?
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
