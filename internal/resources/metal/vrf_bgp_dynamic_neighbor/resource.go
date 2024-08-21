package vrfbgpdynamicneighbor

import (
	"context"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/equinix/terraform-provider-equinix/internal/framework"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var (
	bgpNeighborIncludes = []string{"metal_gateway"}
	// `created_by` is specified as a UserLimited.  To avoid an error
	// due to missing UserLimited.id field, have to either exclude
	// or include `created_by`.  Since we're including `metal_gateway`
	// we also have to either exclude or include `ip_reservation` to
	// avoid a deserialization error due to the required, enumerated
	// `type` property on VRF IP reservations
	bgpNeighborExcludes = []string{"created_by", "ip_reservation"}
)

type Resource struct {
	framework.BaseResource
	framework.WithTimeouts
}

func NewResource() resource.Resource {
	r := Resource{
		BaseResource: framework.NewBaseResource(
			framework.BaseResourceConfig{
				Name: "equinix_metal_vrf_bgp_dynamic_neighbor",
			},
		),
	}

	return &r
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

	resp.Schema = s
}

func (r *Resource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	client := r.Meta.NewMetalClientForFramework(ctx, request.ProviderMeta)

	var plan Model
	response.Diagnostics.Append(request.Config.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	createRequest := metalv1.BgpDynamicNeighborCreateInput{
		BgpNeighborRange: plan.Range.ValueString(),
		BgpNeighborAsn:   plan.ASN.ValueInt64(),
	}

	response.Diagnostics.Append(getPlanTags(ctx, plan, &createRequest.Tags)...)
	if response.Diagnostics.HasError() {
		return
	}

	neighbor, _, err := client.MetalGatewaysApi.CreateBgpDynamicNeighbor(ctx, plan.GatewayID.ValueString()).
		BgpDynamicNeighborCreateInput(createRequest).
		Exclude(bgpNeighborExcludes).
		Include(bgpNeighborIncludes).
		Execute()

	if err != nil {
		response.Diagnostics.AddError(
			"Error creating VRF BGP dynamic neighbor range",
			"Could not create VRF BGP dynamic neighbor range: "+err.Error(),
		)
	}

	// Parse API response into the Terraform state
	response.Diagnostics.Append(plan.parse(ctx, neighbor)...)
	if response.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	client := r.Meta.NewMetalClientForFramework(ctx, request.ProviderMeta)

	var data Model
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	neighbor, _, err := client.VRFsApi.BgpDynamicNeighborsIdGet(ctx, data.ID.ValueString()).
		Exclude(bgpNeighborExcludes).
		Include(bgpNeighborIncludes).
		Execute()

	if err != nil {
		response.Diagnostics.AddError(
			"Error reading VRF BGP dynamic neighbor range",
			"Could not read VRF BGP dynamic neighbor with ID "+data.ID.ValueString()+": "+err.Error(),
		)
	}

	response.Diagnostics.Append(data.parse(ctx, neighbor)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// TODO: ideally it should be possible to update tags, but the API doesn't have an
	// update endpoint for BGP dynamic neighbors, so for now update is a no-op and
	// tag changes force resource recreation
	var data Model
	if diag := req.Plan.Get(ctx, &data); diag.HasError() {
		resp.Diagnostics.Append(diag...)
		return
	}

	if diag := resp.State.Set(ctx, &data); diag.HasError() {
		resp.Diagnostics.Append(diag...)
		return
	}
}

func (r *Resource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	client := r.Meta.NewMetalClientForFramework(ctx, request.ProviderMeta)

	var data Model
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	// TODO: should we do something with the neighbor object returned here?
	// For example: do we need to poll the API until neighbor.GetState() has
	// as particular value?
	_, _, err := client.VRFsApi.DeleteBgpDynamicNeighborById(ctx, data.ID.ValueString()).
		Exclude(bgpNeighborExcludes).
		Execute()

	if err != nil {
		response.Diagnostics.AddError(
			"Error deleting VRF BGP dynamic neighbor range",
			"Could not delete VRF BGP dynamic neighbor with ID "+data.ID.ValueString()+": "+err.Error(),
		)
	}
}

func getPlanTags(ctx context.Context, plan Model, tags *[]string) diag.Diagnostics {
	if len(plan.Tags.Elements()) != 0 {
		return plan.Tags.ElementsAs(ctx, tags, false)
	}
	return diag.Diagnostics{}
}
