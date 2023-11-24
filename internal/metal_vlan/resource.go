package metal_vlan

import (
	"context"
    "fmt"
    "path"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/packethost/packngo"
	"github.com/equinix/terraform-provider-equinix/internal/helper"
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
    "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
    tpfpath "github.com/hashicorp/terraform-plugin-framework/path"
)

type MetalVlanResourceModel struct {
    ID          types.String `tfsdk:"id"`
    ProjectID   types.String `tfsdk:"project_id"`
    Description types.String `tfsdk:"description"`
    Facility    types.String `tfsdk:"facility"`
    Metro       types.String `tfsdk:"metro"`
    Vxlan       types.Int64  `tfsdk:"vxlan"`
}

func (rm *MetalVlanResourceModel) parse(ctx context.Context, vlan *packngo.VirtualNetwork) diag.Diagnostics {
    var diags diag.Diagnostics

    // Assuming 'vlan' is the API response object for a MetalVlan resource
    rm.ProjectID = types.StringValue(vlan.Project.ID)
    rm.Description = types.StringValue(vlan.Description)
    rm.Vxlan = types.Int64Value(int64(vlan.VXLAN))
    rm.Facility = types.StringValue(vlan.FacilityCode)
    rm.Metro = types.StringValue(vlan.MetroCode)

    return diags
}

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "equinix_metal_vlan",
				Schema: &metalVlanResourceSchema,
			},
		),
	}
}

type Resource struct {
	helper.BaseResource
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan MetalVlanResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve the API client
    r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
    client := r.Meta.Metal

    // Translate the plan into an API request
    createRequest := &packngo.VirtualNetworkCreateRequest{
        ProjectID:   plan.ProjectID.ValueString(),
        Description: plan.Description.ValueString(),
        // Include other fields as necessary
    }
    if !plan.Metro.IsNull(){
        createRequest.Metro = plan.Metro.ValueString()
    }
    if !plan.Facility.IsNull(){
        createRequest.Facility = plan.Facility.ValueString()
    }

    // API call to create the MetalVlan resource
    vlan, _, err := client.ProjectVirtualNetworks.Create(createRequest)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error creating MetalVlan",
            "Could not create MetalVlan: " + err.Error(),
        )
        return
    }

    // Parse API response into Terraform state
    diags = plan.parse(ctx, vlan)
    resp.Diagnostics.Append(diags...)
    if diags.HasError() {
        return
    }

    // Set the resource ID and update the state
    resp.State.Set(ctx, &plan)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state MetalVlanResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve the API client
    r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
    client := r.Meta.Metal

    // Extract the ID of the resource from the state
	id := state.ID.ValueString()

    // API call to read the current state of the Metal Vlan
    getOpts := &packngo.GetOptions{Includes: []string{"assigned_to"}}
    vlan, _, err := client.ProjectVirtualNetworks.Get(id, getOpts)
    if err != nil {
        err = helper.FriendlyError(err)
        // Check if the VLAN no longer exists
		if helper.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				"Metal VLAN",
				fmt.Sprintf("[WARN] VLAN (%s) not found, removing from state", id),
			)
			resp.State.RemoveResource(ctx)
			return
		}
        resp.Diagnostics.AddError(
            "Error reading MetalVlan",
            "Could not read MetalVlan with ID " + id + ": " + err.Error(),
        )
        return
    }

    // Parse API response into Terraform state
    diags = state.parse(ctx, vlan)
    resp.Diagnostics.Append(diags...)
    if diags.HasError() {
        return
    }

    // Update the Terraform state
    resp.State.Set(ctx, &state)
}

func (r *Resource) Update(
    ctx context.Context,
    req resource.UpdateRequest,
    resp *resource.UpdateResponse,
) {
	// This resource does not support updates
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var state MetalVlanResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve the API client
    r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
    client := r.Meta.Metal

    // Extract the ID of the resource from the state
	id := state.ID.ValueString()

    // API call to read the current state of the Metal Vlan
    getOpts := &packngo.GetOptions{Includes: []string{
        "instances",
        "instances.network_ports.virtual_networks",
        "internet_gateway",
    }}
    vlan, getResp, err := client.ProjectVirtualNetworks.Get(id, getOpts)
    if helper.IgnoreResponseErrors(helper.HttpForbidden, helper.HttpNotFound)(getResp, err) != nil {
        err = helper.FriendlyError(err)
        resp.Diagnostics.AddError(
            "Error deleting Metal VLAN",
            "Could not retrieve Metal VLAN with ID " + id + ": " + err.Error(),
        )
        return
    } else if err != nil {
		// missing vlans are deleted
		return
	}

    // all device ports must be unassigned before delete
	for _, i := range vlan.Instances {
		for _, p := range i.NetworkPorts {
			for _, a := range p.AttachedVirtualNetworks {
				// a.ID is not set despite including instaces.network_ports.virtual_networks
				// TODO(displague) packngo should offer GetID() that uses ID or Href
				aID := path.Base(a.Href)

				if aID == id {
					_, deleteResp, err := client.Ports.Unassign(p.ID, id)

					if helper.IgnoreResponseErrors(helper.HttpForbidden, helper.HttpNotFound)(deleteResp, err) != nil {
						err = helper.FriendlyError(err)
                        resp.Diagnostics.AddError(
                            "Error deleting Metal VLAN",
                            "Could not unassign Metal VLAN with ID " + id + ": " + err.Error(),
                        )
                        return
					}
				}
			}
		}
	}

    // TODO(displague) do we need to unassign gateway connections before delete?
    err = helper.IgnoreResponseErrors(helper.HttpForbidden, helper.HttpNotFound)(client.ProjectVirtualNetworks.Delete(id))
    if err != nil {
        err = helper.FriendlyError(err)
        resp.Diagnostics.AddError(
            "Error deleting Metal VLAN",
            "Could not delete Metal VLAN with ID " + id + ": " + err.Error(),
        )
        return
    }
}


var metalVlanResourceSchema = schema.Schema{
    Attributes: map[string]schema.Attribute{
        "id": schema.StringAttribute{
            Description: "The unique identifier for the VLAN",
            Computed:    true,
        },
        "project_id": schema.StringAttribute{
            Required:    true,
            Description: "ID of the parent project",
        },
        "description": schema.StringAttribute{
            Optional:    true,
            Description: "Description string",
        },
        "facility": schema.StringAttribute{
            Optional:    true,
            Description: "Facility where to create the VLAN",
            DeprecationMessage: "Use metro instead of facility. For more information, read the migration guide.",
            Validators: []validator.String{
				stringvalidator.ConflictsWith(tpfpath.Expressions{
					tpfpath.MatchRoot("vxlan"),
				}...),
			},
        },
        "metro": schema.StringAttribute{
            Optional:    true,
            Description: "Metro in which to create the VLAN",
            Validators: []validator.String{
				stringvalidator.ExactlyOneOf(tpfpath.Expressions{
					tpfpath.MatchRoot("facility"),
				}...),
			},
        },
        "vxlan": schema.Int64Attribute{
            Optional:    true,
            Computed:    true,
            Description: "VLAN ID, must be unique in metro",
        },
    },
}
