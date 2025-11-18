package port

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func resourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"port_id": schema.StringAttribute{
				Description: "UUID of the port to lookup",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"bonded": schema.BoolAttribute{
				Description: "Flag indicating whether the port should be bonded",
				Required:    true,
			},
			"layer2": schema.BoolAttribute{
				Description: "Flag indicating whether the port is in layer2 (or layer3) mode. The `layer2` flag can be set only for bond ports.",
				Optional:    true,
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
			},
			"network_type": schema.StringAttribute{
				Description: "One of layer2-bonded, layer2-individual, layer3, hybrid and hybrid-bonded. This attribute is only set on bond ports.",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("layer2-bonded", "layer2-individual", "layer3", "hybrid", "hybrid-bonded"),
				},
			},
			"disbond_supported": schema.BoolAttribute{
				Description: "Flag indicating whether the port can be removed from a bond",
				Computed:    true,
			},
			"bond_name": schema.StringAttribute{
				Description: "Name of the bond port",
				Computed:    true,
			},
			"bond_id": schema.StringAttribute{
				Description: "UUID of the bond port",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "Port type",
				Computed:    true,
			},
			"mac": schema.StringAttribute{
				Description: "MAC address of the port",
				Computed:    true,
			},
		},
	}
}
