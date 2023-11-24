package metal_port

import (
	"context"
	"fmt"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/helper"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/packethost/packngo"
)

var (
	l2Types = []string{"layer2-individual", "layer2-bonded"}
	l3Types = []string{"layer3", "hybrid", "hybrid-bonded"}
)

type MetalPortResourceModel struct {
	ID               types.String `tfsdk:"id"`
    PortID           types.String `tfsdk:"port_id"`
    DeviceID         types.String `tfsdk:"device_id"`
    Bonded           types.Bool   `tfsdk:"bonded"`
    Layer2           types.Bool   `tfsdk:"layer2"`
    NativeVLANID     types.String `tfsdk:"native_vlan_id"`
    VxlanIDs         types.Set    `tfsdk:"vxlan_ids"`
    VlanIDs          types.Set    `tfsdk:"vlan_ids"`
    ResetOnDelete    types.Bool   `tfsdk:"reset_on_delete"`
    Name             types.String `tfsdk:"name"`
    NetworkType      types.String `tfsdk:"network_type"`
    DisbondSupported types.Bool   `tfsdk:"disbond_supported"`
    BondName         types.String `tfsdk:"bond_name"`
    BondID           types.String `tfsdk:"bond_id"`
    Type             types.String `tfsdk:"type"`
    Mac              types.String `tfsdk:"mac"`
}

func (rm *MetalPortResourceModel) parse(ctx context.Context, port *packngo.Port) diag.Diagnostics {
    var diags diag.Diagnostics

    // Assuming 'port' is the API response object for a MetalPort resource
    rm.PortID = types.StringValue(port.ID)
    rm.Bonded = types.BoolValue(port.Data.Bonded)
	// Layer2 is only true if the network type is not in l3Types and is in l2Types
    rm.Layer2 = types.BoolValue(
		!helper.Contains(l3Types, port.NetworkType) && helper.Contains(l2Types, port.NetworkType),
	)
    rm.NativeVLANID = types.StringValue("")
    if port.NativeVirtualNetwork != nil {
        rm.NativeVLANID = types.StringValue(port.NativeVirtualNetwork.ID)
    }
    
    // Convert VXLAN IDs and VLAN IDs to types.Set
    portVxlanIDs := make([]int64, len(port.AttachedVirtualNetworks))
	portVlanIDs := make([]string, len(port.AttachedVirtualNetworks))
	for i, v := range port.AttachedVirtualNetworks {
		portVxlanIDs[i] = int64(v.VXLAN)
		portVlanIDs[i] = v.ID
	}
	vxlanIDs, diags := types.SetValueFrom(ctx, types.Int64Type, portVxlanIDs)
	if diags.HasError() {
		return diags
	}
	rm.VxlanIDs = vxlanIDs
	vlanIDs, diags := types.SetValueFrom(ctx, types.StringType, portVlanIDs)
	if diags.HasError() {
		return diags
	}
    rm.VlanIDs = vlanIDs

    rm.Name = types.StringValue(port.Name)
    rm.NetworkType = types.StringValue(port.NetworkType)
    rm.DisbondSupported = types.BoolValue(port.DisbondOperationSupported)
    
    if port.Bond != nil {
        rm.BondName = types.StringValue(port.Bond.Name)
        rm.BondID = types.StringValue(port.Bond.ID)
    }

    rm.Type = types.StringValue(port.Type)
    rm.Mac = types.StringValue(port.Data.MAC)

    return diags
}

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "equinix_metal_port",
				Schema: &metalPortResourceSchema,
			},
		),
	}
}

type Resource struct {
	helper.BaseResource
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan MetalPortResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

	// Retrieve the API client from the provider metadata
	r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
	client := r.Meta.Metal

    // API call to create/update the Metal Port resource
    err := updatePort(ctx, client, plan)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error creating Metal Port resource",
            "Could not create Metal Port resource: " + err.Error(),
        )
        return
    }

	// Retrieve updated Metal Port from API
	port, err := getPortByResourceData(plan, client)
	if err != nil {
		err = helper.FriendlyError(err)
		// If the org was destroyed, mark as gone
		if helper.IsNotFound(err) || helper.IsForbidden(err) {
			resp.Diagnostics.AddWarning(
				"Metal Port",
				fmt.Sprintf("[WARN] Port (%s) not accessible, removing from state", port.ID),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading Metal Port",
			"Could not read port: " + err.Error(),
		)
		return
	}

    // Parse API response into Terraform state
    stateDiags := plan.parse(ctx, port)
    resp.Diagnostics.Append(stateDiags...)
    if stateDiags.HasError() {
        return
    }

    // Set the resource ID
    resp.State.Set(ctx, &plan)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    // Retrieve the current state
    var state MetalPortResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve the API client from the provider meta
	r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
    client := r.Meta.Metal

    // Retrieve updated Metal Port from API
    port, err := getPortByResourceData(state, client)
    if err != nil {
        err = helper.FriendlyError(err)
        if helper.IsNotFound(err) || helper.IsForbidden(err) {
            resp.Diagnostics.AddWarning(
                "Metal Port",
                fmt.Sprintf("[WARN] Port (%s) not accessible, removing from state", state.PortID.ValueString()),
            )
            resp.State.RemoveResource(ctx)
            return
        }
        resp.Diagnostics.AddError(
            "Error reading Metal Port",
            "Could not read port: " + err.Error(),
        )
        return
    }

    // Parse the API response into the Terraform state
    parseDiags := state.parse(ctx, port)
    resp.Diagnostics.Append(parseDiags...)
    if parseDiags.HasError() {
        return
    }

    // Update the Terraform state
    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
}


func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan MetalPortResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

	// Retrieve the API client from the provider metadata
	r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
	client := r.Meta.Metal

    // API call to create/update the Metal Port resource
    err := updatePort(ctx, client, plan)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error creating Metal Port resource",
            "Could not create Metal Port resource: " + err.Error(),
        )
        return
    }

	// Retrieve updated Metal Port from API
	port, err := getPortByResourceData(plan, client)
	if err != nil {
		err = helper.FriendlyError(err)
		// If the org was destroyed, mark as gone
		if helper.IsNotFound(err) || helper.IsForbidden(err) {
			resp.Diagnostics.AddWarning(
				"Metal Port",
				fmt.Sprintf("[WARN] Port (%s) not accessible, removing from state", port.ID),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading Metal Port",
			"Could not read port: " + err.Error(),
		)
		return
	}

    // Parse API response into Terraform state
    stateDiags := plan.parse(ctx, port)
    resp.Diagnostics.Append(stateDiags...)
    if stateDiags.HasError() {
        return
    }

    // Set the resource ID
    resp.State.Set(ctx, &plan)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    // Retrieve the current state
    var state MetalPortResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve the API client from the provider meta
	r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
    client := r.Meta.Metal

	 // Check if the port should be reset to default settings before deletion
	 if state.ResetOnDelete.ValueBool() {

		start := time.Now()
		// update cpd.Data with state
		cpd, cResp, err := getPortData(client, state)
		if helper.IgnoreResponseErrors(helper.HttpForbidden, helper.HttpNotFound)(cResp, err) != nil {
			err = helper.FriendlyError(err)
			resp.Diagnostics.AddError(
                "Error retrieving  Metal Port",
                "Could not retrieve Metal Port: " + err.Error(),
            )
			return
		}

		// to reset the port to defaults we iterate through helpers (used in
		// create/update), some of which rely on resource state. reuse those helpers by
		// setting ephemeral state.
		cpd.Data.Layer2 = types.BoolValue(false)
		cpd.Data.Bonded = types.BoolValue(true)
		cpd.Data.NativeVLANID = types.StringNull()
		vlanIDs, diags := types.SetValueFrom(ctx, types.StringType, []string{})
		if diags.HasError() {
			return
		}
		cpd.Data.VlanIDs = vlanIDs
		cpd.Data.VxlanIDs = types.SetNull(types.Int64Null().Type(ctx))

		for _, f := range [](func(*ClientPortData) error){
			batchVlans(ctx, start, true),
			makeBond,
			convertToL3,
		} {
			if err := f(cpd); err != nil {
				resp.Diagnostics.AddError(
					"Error resetting Metal Port",
					"Could not reset Metal Port to default settings: " + err.Error(),
				)
				return
			}
		}

		// TODO(displague) error or warn?
		if warn := portProperlyDestroyed(cpd.Port); warn != nil {
			resp.Diagnostics.AddWarning(
				"Metal Port",
				fmt.Sprintf("[WARN] %s\n", warn),
			)
		}
    }
}

var metalPortResourceSchema = schema.Schema{
    Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
            Description: "UUID of the port",
            Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
        },
        "bonded": schema.BoolAttribute{
            Required:    true,
            Description: "Flag indicating whether the port should be bonded",
        },
		"port_id": schema.StringAttribute{
            Optional:    true,
            Description: "UUID of the port to lookup. You must specify either port_id or (device_id and name)",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
				stringplanmodifier.UseStateForUnknown(),
			},
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(path.Expressions{
					path.MatchRoot("device_id"),
				}...),
			},
        },
		"device_id": schema.StringAttribute{
            Optional:    true,
            Description: "UUID of the device to lookup",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
				stringplanmodifier.UseStateForUnknown(),
			},
			Validators: []validator.String{
				stringvalidator.AlsoRequires(path.Expressions{
					path.MatchRoot("name"),
				}...),
			},
        },
		"name": schema.StringAttribute{
			Optional:    true,
			Description: "Name of the port to look up, e.g., bond0, eth1",
            Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
				stringplanmodifier.UseStateForUnknown(),
			},
        },
        "layer2": schema.BoolAttribute{
            Optional:    true,
            Description: "Flag indicating whether the port is in layer2 (or layer3) mode. The `layer2` flag can be set only for bond ports.",
        },
        "native_vlan_id": schema.StringAttribute{
            Optional:    true,
            Description: "UUID of native VLAN of the port",
        },
        "vxlan_ids": schema.SetAttribute{
            Optional:      true,
            Computed:      true,
            Description:   "VLAN VXLAN ids to attach (example: [1000])",
            ElementType:   types.Int64Type,
			Validators: []validator.Set{
				setvalidator.ConflictsWith(path.Expressions{
					path.MatchRoot("vlan_ids"),
				}...),
			},
        },
        "vlan_ids": schema.SetAttribute{
            Optional:      true,
            Computed:      true,
            Description:   "UUIDs VLANs to attach. To avoid jitter, use the UUID and not the VXLAN",
            ElementType:      types.StringType,
            Validators: []validator.Set{
				setvalidator.ConflictsWith(path.Expressions{
					path.MatchRoot("vxlan_ids"),
				}...),
			},
        },
        "reset_on_delete": schema.BoolAttribute{
            Optional:    true,
            Description: "Behavioral setting to reset the port to default settings (layer3 bonded mode without any vlan attached) before delete/destroy",
        },
        "network_type": schema.StringAttribute{
            Computed:    true,
            Description: "One of layer2-bonded, layer2-individual, layer3, hybrid, and hybrid-bonded. This attribute is only set on bond ports.",
        },
        "disbond_supported": schema.BoolAttribute{
            Computed:    true,
            Description: "Flag indicating whether the port can be removed from a bond",
        },
        "bond_name": schema.StringAttribute{
            Computed:    true,
            Description: "Name of the bond port",
        },
        "bond_id": schema.StringAttribute{
            Computed:    true,
            Description: "UUID of the bond port",
        },
        "type": schema.StringAttribute{
            Computed:    true,
            Description: "Port type",
        },
        "mac": schema.StringAttribute{
            Computed:    true,
            Description: "MAC address of the port",
        },
    },
}
