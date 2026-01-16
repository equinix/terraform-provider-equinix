package port

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func resourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"port_id": schema.StringAttribute{
				Description: "UUID of the port to lookup",
				Required:    true,
			},
			"bonded": schema.BoolAttribute{
				Description: "Flag indicating whether the port should be bonded",
				Required:    true,
			},
			"layer2": schema.BoolAttribute{
				Description: "Flag indicating whether the port is in layer2 (or layer3) mode. The `layer2` flag can be set only for bond ports.",
				Optional:    true,
				Computed:    true,
			},
			"native_vlan_id": schema.StringAttribute{
				Description: "UUID of native VLAN of the port",
				Optional:    true,
			},
			"vxlan_ids": schema.SetAttribute{
				Description: "VLAN VXLAN ids to attach (example: [1000])",
				ElementType: types.Int32Type,
				Optional:    true,
				Computed:    true,
				// PlanModifiers: []planmodifier.Set{
				// 	UsePlanForNewSetValue(),
				// },
				Validators: []validator.Set{
					setvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("vlan_ids"),
					}...),
				},
			},
			"vlan_ids": schema.SetAttribute{
				Description: "UUIDs VLANs to attach. To avoid jitter, use the UUID and not the VXLAN",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				// PlanModifiers: []planmodifier.Set{
				// 	UsePlanForNewSetValue(),
				// },
				Validators: []validator.Set{
					setvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("vxlan_ids"),
					}...),
				},
			},
			"reset_on_delete": schema.BoolAttribute{
				Description: "Behavioral setting to reset the port to default settings (layer3 bonded mode without any vlan attached) before delete/destroy",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the port to look up, e.g. bond0, eth1",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"network_type": schema.StringAttribute{
				Description: "One of layer2-bonded, layer2-individual, layer3, hybrid and hybrid-bonded. This attribute is only set on bond ports.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("layer2-bonded", "layer2-individual", "layer3", "hybrid", "hybrid-bonded"),
				},
			},
			"disbond_supported": schema.BoolAttribute{
				Description: "Flag indicating whether the port can be removed from a bond",
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"bond_name": schema.StringAttribute{
				Description: "Name of the bond port",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"bond_id": schema.StringAttribute{
				Description: "UUID of the bond port",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Description: "Port type",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"mac": schema.StringAttribute{
				Description: "MAC address of the port",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

type updateToPlanSets struct{}

// Description implements planmodifier.Set.
func (v updateToPlanSets) Description(context.Context) string {
	return "Use the plan value and either the API will accept it or provide an error if it can't achieve it."
}

// MarkdownDescription implements planmodifier.Set.
func (v updateToPlanSets) MarkdownDescription(context.Context) string {
	return "Use the plan value and either the API will accept it or provide an error if it can't achieve it."
}

// PlanModifySet implements planmodifier.Set.
func (v updateToPlanSets) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	resp.RequiresReplace = false

	var plan resourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.VLANIDs.IsUnknown() {
		resp.PlanValue = types.SetUnknown(resp.PlanValue.Type(ctx))
		return
	}

	resp.PlanValue = plan.VLANIDs
}

// UsePlanForNewSetValue just changes the plan so that it shows the intended state for computed states.
func UsePlanForNewSetValue() planmodifier.Set {
	return updateToPlanSets{}
}
