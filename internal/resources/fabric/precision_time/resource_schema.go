package precision_time

import (
	"context"
	"fmt"
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func resourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttributeDefaultDescription(),
			"uuid": schema.StringAttribute{
				Description: "Equinix generated id for the Precision Time Service",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"href": schema.StringAttribute{
				Description: "Equinix generated Portal link for the created Precision Time Service",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Description: "Choose type of Precision Time Service",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						string(fabricv4.PRECISIONTIMESERVICEREQUESTTYPE_NTP),
						string(fabricv4.PRECISIONTIMESERVICEREQUESTTYPE_PTP),
					),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of Precision Time Service. Applicable values: Maximum: 24 characters; Allowed characters: alpha-numeric, hyphens ('-') and underscores ('_')",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(24),
				},
			},
			"description": schema.StringAttribute{
				Description: "Optional description of time service",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Description: fmt.Sprintf("Indicator of the state of this Precision Time Service. One of: [%v]", fabricv4.AllowedPrecisionTimeServiceCreateResponseStateEnumValues),
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"package": schema.SetNestedBlock{
				Description: "Precision Time Service Package Details",
				CustomType:  fwtypes.NewSetNestedObjectTypeOf[PackageModel](ctx),
				Validators: []validator.Set{
					setvalidator.SizeAtMost(1),
				},
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedBlockObject{
					CustomType: fwtypes.NewObjectTypeOf[PackageModel](ctx),
					PlanModifiers: []planmodifier.Object{
						objectplanmodifier.UseStateForUnknown(),
					},
					Attributes: map[string]schema.Attribute{
						"code": schema.StringAttribute{
							Description: "Time Precision Package Code for the desired billing package",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.OneOf(
									string(fabricv4.GETTIMESERVICESPACKAGEBYCODEPACKAGECODEPARAMETER_NTP_STANDARD),
									string(fabricv4.GETTIMESERVICESPACKAGEBYCODEPACKAGECODEPARAMETER_NTP_ENTERPRISE),
									string(fabricv4.GETTIMESERVICESPACKAGEBYCODEPACKAGECODEPARAMETER_PTP_STANDARD),
									string(fabricv4.GETTIMESERVICESPACKAGEBYCODEPACKAGECODEPARAMETER_PTP_ENTERPRISE),
								),
							},
						},
						"href": schema.StringAttribute{
							Description: "Time Precision Package HREF link to corresponding resource in Equinix Portal",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
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
			},
			"ipv4": schema.SetNestedBlock{
				Description: "An object that has Network IP Configurations for Timing Master Servers.",
				CustomType:  fwtypes.NewSetNestedObjectTypeOf[Ipv4Model](ctx),
				Validators: []validator.Set{
					setvalidator.SizeAtMost(1),
				},
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedBlockObject{
					CustomType: fwtypes.NewObjectTypeOf[Ipv4Model](ctx),
					PlanModifiers: []planmodifier.Object{
						objectplanmodifier.UseStateForUnknown(),
					},
					Attributes: map[string]schema.Attribute{
						"primary": schema.StringAttribute{
							Description: "IPv4 address for the Primary Timing Master Server.",
							Required:    true,
						},
						"secondary": schema.StringAttribute{
							Description: "IPv4 address for the Secondary Timing Master Server.",
							Required:    true,
						},
						"network_mask": schema.StringAttribute{
							Description: "IPv4 address that defines the range of consecutive subnets in the network.",
							Required:    true,
						},
						"default_gateway": schema.StringAttribute{
							Description: "IPv4 address that establishes the Routing Interface where traffic is directed. It serves as the next hop in the Network.",
							Required:    true,
						},
					},
				},
			},
			"advance_configuration": schema.SingleNestedBlock{
				Description: "An object that has advanced configuration options.",
				CustomType:  fwtypes.NewObjectTypeOf[AdvanceConfigurationModel](ctx),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"ntp": schema.ListAttribute{
						Description: "Advance Configuration for NTP; a list of MD5 objects",
						Optional:    true,
						Computed:    true,
						CustomType:  fwtypes.NewListNestedObjectTypeOf[MD5Model](ctx),
						ElementType: fwtypes.NewObjectTypeOf[MD5Model](ctx),
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"ptp": schema.SingleNestedBlock{
						Description: "An object that has advanced PTP configuration.",
						CustomType:  fwtypes.NewObjectTypeOf[PTPModel](ctx),
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"time_scale": schema.StringAttribute{
								Description: "Time scale value. ARB denotes Arbitrary, and PTP denotes Precision Time Protocol.",
								Optional:    true,
								Computed:    true,
								Validators: []validator.String{
									stringvalidator.OneOf(
										string(fabricv4.PTPADVANCECONFIGURATIONTIMESCALE_ARB),
										string(fabricv4.PTPADVANCECONFIGURATIONTIMESCALE_PTP),
									),
								},
							},
							"domain": schema.Int64Attribute{
								Description: "Represents the domain number associated with the PTP profile. This is used to differentiate multiple PTP networks within a single physical network.",
								Optional:    true,
								Computed:    true,
							},
							"priority_1": schema.Int64Attribute{
								Description: "Specifies the priority level 1 for the clock. The value helps in determining the best clock in the PTP network. Lower values are considered higher priority.",
								Optional:    true,
								Computed:    true,
							},
							"priority_2": schema.Int64Attribute{
								Description: "Specifies the priority level 2 for the clock. It acts as a tie-breaker if multiple clocks have the same priority 1 value. Lower values are considered higher priority.",
								Optional:    true,
								Computed:    true,
							},
							"log_announce_interval": schema.Int64Attribute{
								Description: "Represents the log2 interval between consecutive PTP announce messages. For example, a value of 0 implies an interval of 2^0 = 1 second.",
								Optional:    true,
								Computed:    true,
							},
							"log_sync_interval": schema.Int64Attribute{
								Description: "Represents the log2 interval between consecutive PTP synchronization messages. A value of 0 implies an interval of 2^0 = 1 second.",
								Optional:    true,
								Computed:    true,
							},
							"log_delay_req_interval": schema.Int64Attribute{
								Description: "Represents the log2 interval between consecutive PTP delay request messages. A value of 0 implies an interval of 2^0 = 1 second.",
								Optional:    true,
								Computed:    true,
							},
							"transport_mode": schema.StringAttribute{
								Description: "Mode of transport for the Time Precision Service.",
								Optional:    true,
								Computed:    true,
								Validators: []validator.String{
									stringvalidator.OneOf(
										string(fabricv4.PTPADVANCECONFIGURATIONTRANSPORTMODE_MULTICAST),
										string(fabricv4.PTPADVANCECONFIGURATIONTRANSPORTMODE_UNICAST),
										string(fabricv4.PTPADVANCECONFIGURATIONTRANSPORTMODE_HYBRID),
									),
								},
							},
							"grant_time": schema.Int64Attribute{
								Description: "Unicast Grant Time in seconds. For Multicast and Hybrid transport modes, grant time defaults to 300 seconds. For Unicast mode, grant time can be between 30 to 7200.",
								Optional:    true,
								Computed:    true,
							},
						},
					},
				},
			},
			"project": schema.SetNestedBlock{
				Description: "An object that contains the Equinix Fabric project_id used for linking the Time Precision Service to a specific Equinix Fabric Project",
				CustomType:  fwtypes.NewSetNestedObjectTypeOf[ProjectModel](ctx),
				Validators: []validator.Set{
					setvalidator.SizeAtMost(1),
				},
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedBlockObject{
					CustomType: fwtypes.NewObjectTypeOf[ProjectModel](ctx),
					PlanModifiers: []planmodifier.Object{
						objectplanmodifier.UseStateForUnknown(),
					},
					Attributes: map[string]schema.Attribute{
						"project_id": schema.StringAttribute{
							Description: "Equinix Fabric Project ID",
							Required:    true,
						},
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
							Required:    true,
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
			"account": schema.SingleNestedBlock{
				Description: "Equinix User Account associated with Precision Time Service",
				CustomType:  fwtypes.NewObjectTypeOf[AccountModel](ctx),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"account_number": schema.Int64Attribute{
						Description: "Equinix User account number",
						Computed:    true,
					},
					"is_reseller_account": schema.BoolAttribute{
						Description: "Equinix User Boolean flag indicating if it is a reseller account",
						Computed:    true,
					},
					"org_id": schema.StringAttribute{
						Description: "Equinix User organization id",
						Computed:    true,
					},
					"global_org_id": schema.StringAttribute{
						Description: "Equinix User global organization id",
						Computed:    true,
					},
				},
			},
		},
	}
}
