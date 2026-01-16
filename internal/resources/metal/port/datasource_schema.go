package port

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func datasourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"port_id": schema.StringAttribute{
				Description: "UUID of the port to lookup",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.Expressions{
							path.MatchRoot("device_id"),
							path.MatchRoot("name"),
						}...,
					),
				},
			},
			"device_id": schema.StringAttribute{
				Description: "Device UUID where to look up the port",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.Expressions{
							path.MatchRoot("port_id"),
						}...,
					),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the port to look up, e.g. bond0, eth1 ",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.Expressions{
							path.MatchRoot("port_id"),
						}...,
					),
				},
			},
			"network_type": schema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"mac": schema.StringAttribute{
				Description: "MAC address of the port",
				Computed:    true,
			},
			"bond_id": schema.StringAttribute{
				Description: "UUID of the bond port",
				Computed:    true,
			},
			"bond_name": schema.StringAttribute{
				Description: "Name of the bond port",
				Computed:    true,
			},
			"bonded": schema.BoolAttribute{
				Description: "Flag indicating whether the port is bonded",
				Computed:    true,
			},
			"disbond_supported": schema.BoolAttribute{
				Description: "Flag indicating whether the port can be removed from a bond",
				Computed:    true,
			},
			"native_vlan_id": schema.StringAttribute{
				Description: "UUID of native VLAN of the port",
				Computed:    true,
			},
			"vlan_ids": schema.SetAttribute{
				Description: "UUIDs of attached VLANs",
				Computed:    true,
				ElementType: types.StringType,
			},
			"vxlan_ids": schema.SetAttribute{
				Description: "VLAN tags of attached VLANs",
				Computed:    true,
				ElementType: types.Int32Type,
			},
			"layer2": schema.BoolAttribute{
				Description: "Flag indicating whether the port is in layer2 (or layer3) mode",
				Computed:    true,
			},
		},
	}
}
