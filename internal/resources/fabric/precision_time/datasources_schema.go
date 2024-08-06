package precision_time

import (
	"context"
	"fmt"
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func dataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttributeDefaultDescription(),
			"uuid": schema.StringAttribute{
				Description: "Uuid of Precision Time Service resource; used for lookup",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "Choose type of Precision Time Service",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						string(fabricv4.PRECISIONTIMESERVICEREQUESTTYPE_NTP),
						string(fabricv4.PRECISIONTIMESERVICEREQUESTTYPE_PTP),
					),
				},
			},
			"href": schema.StringAttribute{
				Description: "Equinix generated Portal link for the created Precision Time Service",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of Precision Time Service. Applicable values: Maximum: 24 characters; Allowed characters: alpha-numeric, hyphens ('-') and underscores ('_')",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(24),
				},
			},
			"description": schema.StringAttribute{
				Description: "Optional description of time service",
				Computed:    true,
			},
			"state": schema.StringAttribute{
				Description: fmt.Sprintf("Indicator of the state of this Precision Time Service. One of: [%v]", fabricv4.AllowedPrecisionTimeServiceCreateResponseStateEnumValues),
				Computed:    true,
			},
			"project_id": schema.StringAttribute{
				Description: "Equinix Fabric Project ID",
				Computed:    true,
			},
			"advance_configuration": schema.ListAttribute{
				Description: "An object that has advanced configuration options.",
				CustomType:  fwtypes.NewListNestedObjectTypeOf[AdvanceConfigurationModel](ctx),
				ElementType: fwtypes.NewObjectTypeOf[AdvanceConfigurationModel](ctx),
				Computed:    true,
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"account": schema.ListAttribute{
				Description: "Equinix User Account associated with Precision Time Service",
				CustomType:  fwtypes.NewListNestedObjectTypeOf[AccountModel](ctx),
				ElementType: fwtypes.NewObjectTypeOf[AccountModel](ctx),
				Computed:    true,
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"package": schema.SingleNestedBlock{
				Description: "Precision Time Service Package Details",
				CustomType:  fwtypes.NewObjectTypeOf[PackageModel](ctx),
				Attributes: map[string]schema.Attribute{
					"code": schema.StringAttribute{
						Description: "Time Precision Package Code for the desired billing package",
						Computed:    true,
					},
					"href": schema.StringAttribute{
						Description: "Time Precision Package HREF link to corresponding resource in Equinix Portal",
						Computed:    true,
					},
					"type": schema.StringAttribute{
						Description: "Type of the Precision Time Service Package",
						Computed:    true,
					},
					"bandwidth": schema.Int64Attribute{
						Description: "Bandwidth of the Precision Time Service",
						Computed:    true,
					},
					"clients_per_second_max": schema.Int64Attribute{
						Description: "Maximum clients available per second for the Precision Time Service",
						Computed:    true,
					},
					"redundancy_supported": schema.BoolAttribute{
						Description: "Boolean flag indicating if this Precision Time Service supports redundancy",
						Computed:    true,
					},
					"multi_subnet_supported": schema.BoolAttribute{
						Description: "Boolean flag indicating if this Precision Time Service supports multi subnetting",
						Computed:    true,
					},
					"accuracy_unit": schema.StringAttribute{
						Description: "Time unit of accuracy for the Precision Time Service; e.g. microseconds",
						Computed:    true,
					},
					"accuracy_sla": schema.Int64Attribute{
						Description: "SLA for the accuracy provided by the Precision Time Service",
						Computed:    true,
					},
					"accuracy_avg_min": schema.Int64Attribute{
						Description: "Average minimum accuracy provided by the Precision Time Service",
						Computed:    true,
					},
					"accuracy_avg_max": schema.Int64Attribute{
						Description: "Average maximum accuracy provided by the Precision Time Service",
						Computed:    true,
					},
				},
			},
			"ipv4": schema.SingleNestedBlock{
				Description: "An object that has Network IP Configurations for Timing Master Servers.",
				CustomType:  fwtypes.NewObjectTypeOf[Ipv4Model](ctx),
				Attributes: map[string]schema.Attribute{
					"primary": schema.StringAttribute{
						Description: "IPv4 address for the Primary Timing Master Server.",
						Computed:    true,
					},
					"secondary": schema.StringAttribute{
						Description: "IPv4 address for the Secondary Timing Master Server.",
						Computed:    true,
					},
					"network_mask": schema.StringAttribute{
						Description: "IPv4 address that defines the range of consecutive subnets in the network.",
						Computed:    true,
					},
					"default_gateway": schema.StringAttribute{
						Description: "IPv4 address that establishes the Routing Interface where traffic is directed. It serves as the next hop in the Network.",
						Computed:    true,
					},
				},
			},
			"connections": schema.ListNestedBlock{
				Description: "An array of objects with unique identifiers of connections.",
				CustomType:  fwtypes.NewListNestedObjectTypeOf[ConnectionModel](ctx),
				NestedObject: schema.NestedBlockObject{
					CustomType: fwtypes.NewObjectTypeOf[ConnectionModel](ctx),
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Description: "Equinix Fabric Connection UUID; Precision Time Service will be connected with it",
							Computed:    true,
						},
						"href": schema.StringAttribute{
							Description: "Link to the Equinix Fabric Connection associated with the Precision Time Service",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "Type of the Equinix Fabric Connection associated with the Precision Time Service",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}
