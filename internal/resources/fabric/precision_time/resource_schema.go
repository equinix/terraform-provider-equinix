package precision_time

import (
	"context"
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func resourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttributeDefaultDescription(),
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
			"package": schema.SingleNestedAttribute{
				Description: "An object that has the Time Precision Package code and Time Precision Package href",
				Required:    true,
				CustomType:  fwtypes.NewObjectTypeOf[PackageModel](ctx),
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
				},
			},
			"connections": schema.ListAttribute{
				Description: "An array of objects with unique identifiers of connections.",
				Required:    true,
				CustomType:  fwtypes.NewListNestedObjectTypeOf[ConnectionModel](ctx),
				ElementType: fwtypes.NewObjectTypeOf[ConnectionModel](ctx),
			},
			"ipv4": schema.SingleNestedAttribute{
				Description: "An object that has Network IP Configurations for Timing Master Servers.",
				Required:    true,
				CustomType:  fwtypes.NewObjectTypeOf[Ipv4Model](ctx),
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
			"advance_configuration": schema.SingleNestedAttribute{
				Description: "An object that has advanced configuration options.",
				Optional:    true,
				Computed:    true,
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
					},
					"ptp": schema.SingleNestedAttribute{
						Description: "An object that has advanced PTP configuration.",
						Optional:    true,
						Computed:    true,
						CustomType:  fwtypes.NewObjectTypeOf[PTPModel](ctx),
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
							"priority1": schema.Int64Attribute{
								Description: "Specifies the priority level 1 for the clock. The value helps in determining the best clock in the PTP network. Lower values are considered higher priority.",
								Optional:    true,
								Computed:    true,
							},
							"priority2": schema.Int64Attribute{
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
			"project": schema.SingleNestedAttribute{
				Description: "An object that contains the Equinix Fabric project_id used for linking the Time Precision Service to a specific Equinix Fabric Project",
				Optional:    true,
				Computed:    true,
				CustomType:  fwtypes.NewObjectTypeOf[ProjectModel](ctx),
				Attributes: map[string]schema.Attribute{
					"project_id": schema.StringAttribute{
						Description: "Equinix Fabric Project ID",
						Required:    true,
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}
