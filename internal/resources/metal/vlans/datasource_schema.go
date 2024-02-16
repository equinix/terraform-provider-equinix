package vlans

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func dataSourceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for this Metal Vlan",
				Computed:    true,
			},
			"project_id": schema.StringAttribute{
				Description: "ID of parent project of the VLAN. Use together with vxlan and metro or facility",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("vlan_id"),
					}...),
				},
			},
			"vxlan": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Validators: []validator.Int64{
					int64validator.ConflictsWith(path.Expressions{
						path.MatchRoot("vlan_id"),
					}...),
				},
				Description: "VXLAN numner of the VLAN. Unique in a project and facility or metro. Use with project_id",
			},
			"facility": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("vlan_id"),
						path.MatchRoot("metro"),
					}...),
				},
				Description:        "Facility where the VLAN is deployed",
				DeprecationMessage: "Use metro instead of facility.  For more information, read the migration guide: https://registry.terraform.io/providers/equinix/equinix/latest/docs/guides/migration_guide_facilities_to_metros_devices",
			},
			"metro": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("vlan_id"),
						path.MatchRoot("facility"),
					}...),
				},
				Description: "Metro where the VLAN is deployed",
			},
			"vlan_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("project_id"),
						path.MatchRoot("vxlan"),
						path.MatchRoot("metro"),
						path.MatchRoot("facility"),
					}...),
				},
				Description: "Metal UUID of the VLAN resource",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "VLAN description text",
			},
			"assigned_devices_ids": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "List of device IDs to which this VLAN is assigned",
			},
		},
	}
}
