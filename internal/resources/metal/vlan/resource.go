package vlan

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strings"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var (
	vlanDefaultIncludes = []string{"assigned_to", "metro"}
)

type Resource struct {
	framework.BaseResource
	framework.WithTimeouts
}

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

	var data ResourceModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	if data.Facility.IsNull() && data.Metro.IsNull() {
		response.Diagnostics.AddError("Invalid input params",
			"one of facility or metro must be configured")
		return
	}
	if !data.Facility.IsNull() && !data.Vxlan.IsNull() {
		response.Diagnostics.AddError("Invalid input params",
			"you can set vxlan only for metro vlan")
		return
	}

	createRequest := metalv1.VirtualNetworkCreateInput{
		Description: data.Description.ValueStringPointer(),
	}
	if !data.Metro.IsNull() {
		createRequest.Metro = metalv1.PtrString(strings.ToLower(data.Metro.ValueString()))
	}
	if !data.Vxlan.IsNull() {
		createRequest.Vxlan = metalv1.PtrInt32(int32(data.Vxlan.ValueInt64()))
	}
	if !data.Facility.IsNull() {
		createRequest.Facility = data.Facility.ValueStringPointer()
	}
	vlan, _, err := client.VLANsApi.CreateVirtualNetwork(ctx, data.ProjectID.ValueString()).VirtualNetworkCreateInput(createRequest).Execute()
	if err != nil {
		response.Diagnostics.AddError("Error creating Vlan", err.Error())
		return
	}

	// get the current state of newly created vlan with default include fields
	vlan, err = loadVlan(ctx, client, vlan.GetId())
	if err != nil {
		response.Diagnostics.AddError("Error reading Vlan after create", err.Error())
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
	client := r.Meta.NewMetalClientForFramework(ctx, request.ProviderMeta)

	var data ResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	vlan, err := loadVlan(ctx, client, data.ID.ValueString())
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
			err.Error())
		return
	}

	response.Diagnostics.Append(data.parse(vlan)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ResourceModel
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

	var data ResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	vlan, resp, err := client.VLANsApi.GetVirtualNetwork(
		ctx,
		data.ID.ValueString(),
	).Include([]string{"instances", "virtual_networks"}).Execute()
	if err != nil {
		if err := equinix_errors.IgnoreHttpResponseErrors(http.StatusForbidden, http.StatusNotFound)(resp, err); err != nil {
			response.Diagnostics.AddWarning(
				"Equinix Metal Vlan not found during delete",
				err.Error(),
			)
			return
		}
		response.Diagnostics.AddError("Error fetching Vlan using vlanId",
			err.Error())
		return
	}

	// all device ports must be unassigned before delete
	for _, instance := range vlan.Instances {
		for _, port := range instance.NetworkPorts {
			for _, v := range port.GetVirtualNetworks() {
				if v.GetId() == vlan.GetId() {
					_, resp, err = client.PortsApi.UnassignPort(ctx, port.GetId()).PortAssignInput(metalv1.PortAssignInput{
						Vnid: vlan.Id,
					}).Execute()
					if equinix_errors.IgnoreHttpResponseErrors(http.StatusForbidden, http.StatusNotFound)(resp, err) != nil {
						response.Diagnostics.AddError("Error unassign port with Vlan", err.Error())
						return
					}
				}
			}
		}
	}

	resp, err = client.VLANsApi.DeleteVirtualNetwork(ctx, vlan.GetId()).Execute()
	if err := equinix_errors.IgnoreHttpResponseErrors(http.StatusForbidden, http.StatusNotFound)(resp, err); err != nil {
		response.Diagnostics.AddError("Error deleting Vlan",
			err.Error())
		return
	}
}

func loadVlan(ctx context.Context, client *metalv1.APIClient, id string) (*metalv1.VirtualNetwork, error) {
	vlan, _, err := client.VLANsApi.GetVirtualNetwork(ctx, id).Include(vlanDefaultIncludes).Execute()

	// If this is a facility-based VLAN we have to find the facility separately
	// because if we include `facility`, the facility.href attribute becomes
	// unreachable and the API response above won't validate against the VLAN schema
	if err == nil && vlan.Facility != nil {
		facilityId := path.Base(vlan.Facility.GetHref())
		facilities, _, err := client.FacilitiesApi.FindFacilities(ctx).Include([]metalv1.FindFacilitiesIncludeParameterInner{"metro"}).Execute()
		if err == nil {
			for _, facility := range facilities.GetFacilities() {
				if facility.GetId() == facilityId {
					vlan.Metro = (*metalv1.Metro)(facility.Metro)
				}
			}
			if vlan.Metro == nil {
				return vlan, fmt.Errorf("could not find facility %v for VLAN %v", facilityId, vlan.GetId())
			}
		}
	}

	return vlan, err
}
