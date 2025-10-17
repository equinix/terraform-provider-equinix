// Package precisiontime for EPT resources and data sources
package precisiontime

import (
	"context"
	"fmt"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func dataSourceAllEptServicesSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: `Fabric V4 API compatible data resource that allow user to fetch Equinix Fabric Precision Time Services with pagination details
Additional Documentation:
* API: https://docs.equinix.com/en-us/Content/KnowledgeCenter/Fabric/API-Reference/API-Precision-Time.htm`,
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttributeDefaultDescription(),
			"data": schema.ListNestedAttribute{
				Description: "Returned list of route aggregation objects",
				Computed:    true,
				CustomType:  fwtypes.NewListNestedObjectTypeOf[basePrecisionTimeModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: getEptServiceSchema(ctx),
				},
			},
			"filters": schema.ListNestedAttribute{
				Description: "List of filters to apply to the stream attachment get request. Maximum of 8. All will be AND'd together with 1 of the 8 being a possible OR group of 3",
				Optional:    true,
				Computed:    true,
				CustomType:  fwtypes.NewListNestedObjectTypeOf[filterModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"property": schema.StringAttribute{
							Description: "Property to apply the filter to",
							Required:    true,
						},
						"operator": schema.StringAttribute{
							Description: "Operation applied to the values of the filter",
							Required:    true,
						},
						"values": schema.ListAttribute{
							Description: "List of values to apply the operation to for the specified property",
							Required:    true,
							ElementType: types.StringType,
						},
						"or": schema.BoolAttribute{
							Description: "Boolean value to specify if this filter is a part of the OR group. Has a maximum of 3 and only counts for 1 of the 8 possible filters",
							Optional:    true,
						},
					},
				},
			},
			"pagination": schema.SingleNestedAttribute{
				Description: "Pagination details for the returned route aggregations list",
				Optional:    true,
				CustomType:  fwtypes.NewObjectTypeOf[paginationModel](ctx),
				Attributes: map[string]schema.Attribute{
					"offset": schema.Int32Attribute{
						Description: "Index of the first item returned in the response. The default is 0",
						Optional:    true,
						Computed:    true,
					},
					"limit": schema.Int32Attribute{
						Description: "Maximum number of search results returned per page. Number must be between 1 and 100, and the default is 20",
						Optional:    true,
						Computed:    true,
					},
				},
			},
			"sort": schema.ListNestedAttribute{
				Description: "Filters for the Data Source Search Request",
				Optional:    true,
				CustomType:  fwtypes.NewListNestedObjectTypeOf[sortModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"direction": schema.StringAttribute{
							Description: "The sorting direction. Can be one of: [DESC, ASC], Defaults to DESC",
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.OneOf("DESC", "ASC"),
							},
						},
						"property": schema.StringAttribute{
							Description: fmt.Sprintf("The property name to use in sorting. One of %v Defaults to /name", fabricv4.AllowedTimeServiceSortByEnumValues),
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceSingleEptServiceSchema(ctx context.Context) schema.Schema {
	baseEptServiceSchema := getEptServiceSchema(ctx)
	baseEptServiceSchema["id"] = framework.IDAttributeDefaultDescription()
	baseEptServiceSchema["ept_service_id"] = schema.StringAttribute{
		Description: "The uuid of the EPT Service this data source should retrieve",
		Required:    true,
	}
	return schema.Schema{
		Description: `Fabric V4 API compatible data resource that allow user to fetch Equinix Precision Time Service by UUID
Additional Documentation:
* API: https://developer.equinix.com/catalog/fabricv4#tag/Precision-Time`,
		Attributes: baseEptServiceSchema,
	}

}
func getEptServiceSchema(ctx context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
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
		"name": schema.StringAttribute{
			Description: "Name of Precision Time Service. Applicable values: Maximum: 24 characters; Allowed characters: alpha-numeric, hyphens ('-') and underscores ('_')",
			Computed:    true,
			Validators: []validator.String{
				stringvalidator.LengthAtMost(24),
			},
		},
		"package": schema.SingleNestedAttribute{
			Description: "Precision Time Service Package Details",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[packageModel](ctx),
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
				},
			},
		},
		"operation": schema.SingleNestedAttribute{
			Description: "Precision Time Service Operation",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[operationModel](ctx),
			Attributes: map[string]schema.Attribute{
				"operational_status": schema.StringAttribute{
					Description: "",
					Optional:    true,
					Computed:    true,
					Validators: []validator.String{
						stringvalidator.OneOf("UP", "DOWN", "DEGRADED"),
					},
				},
			},
		},
		"connections": schema.ListNestedAttribute{
			Description: "An array of objects with unique identifiers of connections.",
			CustomType:  fwtypes.NewListNestedObjectTypeOf[connectionModel](ctx),
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
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
		"ipv4": schema.SingleNestedAttribute{
			Description: "An object that has Network IP Configurations for Timing Master Servers.",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[ipv4Model](ctx),
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
		"ntp_advanced_configuration": schema.ListNestedAttribute{
			Description: "NTP Advanced configuration",
			Optional:    true,
			CustomType:  fwtypes.NewListNestedObjectTypeOf[ntpAdvanceConfigurationModel](ctx),
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "md5 Authentication type",
						Optional:    true,
					},
					"key_number": schema.Int32Attribute{
						Description: "The authentication Key ID",
						Optional:    true,
						Validators: []validator.Int32{
							int32validator.AtLeast(1),
							int32validator.AtMost(65535),
						},
					},
					"key": schema.StringAttribute{
						Description: "The plaintext authentication key. For ASCII type, the key\\\n\\ must contain printable ASCII characters, range 10-20 characters. For\\\n\\ HEX type, range should be 10-40 characters",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(10),
							stringvalidator.LengthAtMost(40),
						},
					},
				},
			},
		},
		"ptp_advanced_configuration": schema.SingleNestedAttribute{
			Description: "PTP Advanced Configuration",
			Optional:    true,
			CustomType:  fwtypes.NewObjectTypeOf[ptpAdvanceConfigurationModel](ctx),
			Attributes: map[string]schema.Attribute{
				"time_scale": schema.StringAttribute{
					Description: "Time Scale value, ARB denotes Arbitrary and PTP denotes Precision Time Protocol",
					Optional:    true,
				},
				"domain": schema.Int32Attribute{
					Description: "The PTP domain value",
					Optional:    true,
					Validators: []validator.Int32{
						int32validator.AtLeast(0),
						int32validator.AtMost(127),
					},
				},
				"priority1": schema.Int32Attribute{
					Description: "The priority1 value determines the best primary clock, Lower value indicates higher priority",
					Optional:    true,
					Validators: []validator.Int32{
						int32validator.AtLeast(0),
						int32validator.AtMost(248),
					},
				},
				"priority2": schema.Int32Attribute{
					Description: "The priority2 value differentiates and prioritizes the primary clock to avoid confusion when priority1-value is the same for different primary clocks in a network",
					Optional:    true,
					Validators: []validator.Int32{
						int32validator.AtLeast(0),
						int32validator.AtMost(248),
					},
				},
				"log_announce_interval": schema.Int32Attribute{
					Description: "Logarithmic value that controls the rate of PTP Announce packets from the PTP time server. Default is 1 (1 packet every 2 seconds), Unit packets/second",
					Optional:    true,
				},
				"log_sync_interval": schema.Int32Attribute{
					Description: "Logarithmic value that controls the rate of PTP Sync packets. Default is -4 (16 packets per second), Unit packets/second..",
					Optional:    true,
				},
				"log_delay_req_interval": schema.Int32Attribute{
					Description: "Logarithmic value that controls the rate of PTP DelayReq packets. Default is -4 (16 packets per second), Unit packets/second..",
					Optional:    true,
				},
				"transport_mode": schema.StringAttribute{
					Description: "ptp transport mode",
					Optional:    true,
				},
				"grant_time": schema.Int32Attribute{
					Description: "Unicast Grant Time in seconds. For Multicast and Hybrid transport modes, grant time defaults to 300 seconds. For Unicast mode, grant time can be between 30 to 7200",
					Optional:    true,
					Validators: []validator.Int32{
						int32validator.AtLeast(30),
						int32validator.AtMost(7200),
					},
				},
			},
		},
		"uuid": schema.StringAttribute{
			Description: "Equinix generated id for the Precision Time Service",
			Computed:    true,
		},
		"href": schema.StringAttribute{
			Description: "Equinix generated Portal link for the created Precision Time Service",
			Computed:    true,
		},
		"state": schema.StringAttribute{
			Description: "Indicator of the state of this Precision Time Service",
			Computed:    true,
		},
		"project": schema.SingleNestedAttribute{
			Description: "Equinix Project attribute object",
			Optional:    true,
			CustomType:  fwtypes.NewObjectTypeOf[projectModel](ctx),
			Attributes: map[string]schema.Attribute{
				"project_id": schema.StringAttribute{
					Description: "Equinix Subscriber-assigned project ID",
					Required:    true,
				},
			},
		},
		"account": schema.SingleNestedAttribute{
			Description: "Equinix User Account associated with Precision Time Service",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[accountModel](ctx),
			Attributes: map[string]schema.Attribute{
				"account_number": schema.Int64Attribute{
					Description: "Equinix Account Number",
					Computed:    true,
				},
				"account_name": schema.StringAttribute{
					Description: "Account Name",
					Computed:    true,
				},
				"org_id": schema.Int64Attribute{
					Description: " Customer organization identifier",
					Computed:    true,
				},
				"organization_name": schema.StringAttribute{
					Description: "Customer organization name",
					Computed:    true,
				},
				"global_org_id": schema.StringAttribute{
					Description: "Customer organization naidentifierme",
					Computed:    true,
				},
				"global_organization_name": schema.StringAttribute{
					Description: "Global organization name",
					Computed:    true,
				},
				"ucm_id": schema.StringAttribute{
					Description: "Global organization name",
					Computed:    true,
				},
				"global_cust_id": schema.StringAttribute{
					Description: "Global Customer Id",
					Computed:    true,
				},
				"reseller_account_number": schema.Int64Attribute{
					Description: "Reseller account number",
					Computed:    true,
				},
				"reseller_account_name": schema.StringAttribute{
					Description: "Reseller account name",
					Computed:    true,
				},
				"reseller_ucm_id": schema.StringAttribute{
					Description: "Reseller account ucmId",
					Computed:    true,
				},
				"reseller_org_id": schema.Int64Attribute{
					Description: "Reseller customer organization identifier",
					Computed:    true,
				},
			},
		},
		"order": schema.SingleNestedAttribute{
			Description: "Precision Time Order",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[orderModel](ctx),
			Attributes: map[string]schema.Attribute{
				"purchase_order_number": schema.StringAttribute{
					Description: "Purchase order number",
					Computed:    true,
				},
				"customer_reference_number": schema.StringAttribute{
					Description: "Customer reference number",
					Computed:    true,
				},
				"order_number": schema.StringAttribute{
					Description: "Order reference number",
					Computed:    true,
				},
			},
		},
		"precision_time_price": schema.SingleNestedAttribute{
			Description: "Precision Time Service Price",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[precisionTimePriceModel](ctx),
			Attributes: map[string]schema.Attribute{
				"currency": schema.StringAttribute{
					Description: "Offering price currency",
					Computed:    true,
				},
				"charges": schema.ListNestedAttribute{
					Description: "offering price charge",
					Computed:    true,
					CustomType:  fwtypes.NewListNestedObjectTypeOf[chargesModel](ctx),
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Description: "Price charge type; MONTHLY_RECURRING, NON_RECURRING",
								Computed:    true,
							},
							"price": schema.Float32Attribute{
								Description: "Offering price",
								Computed:    true,
							},
						},
					},
				},
			},
		},
		"change_log": schema.SingleNestedAttribute{
			Description: "Details of the last change on the route aggregation resource",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[changeLogModel](ctx),
			Attributes: map[string]schema.Attribute{
				"created_by": schema.StringAttribute{
					Description: "User name of creator of the route aggregation resource",
					Computed:    true,
				},
				"created_by_full_name": schema.StringAttribute{
					Description: "Legal name of creator of the route aggregation resource",
					Computed:    true,
				},
				"created_by_email": schema.StringAttribute{
					Description: "Email of creator of the route aggregation resource",
					Computed:    true,
				},
				"created_date_time": schema.StringAttribute{
					Description: "Creation time of the route aggregation resource",
					Computed:    true,
				},
				"updated_by": schema.StringAttribute{
					Description: "User name of last updater of the route aggregation resource",
					Computed:    true,
				},
				"updated_by_full_name": schema.StringAttribute{
					Description: "Legal name of last updater of the route aggregation resource",
					Computed:    true,
				},
				"updated_by_email": schema.StringAttribute{
					Description: "Email of last updater of the route aggregation resource",
					Computed:    true,
				},
				"updated_date_time": schema.StringAttribute{
					Description: "Last update time of the route aggregation resource",
					Computed:    true,
				},
				"deleted_by": schema.StringAttribute{
					Description: "User name of deleter of the route aggregation resource",
					Computed:    true,
				},
				"deleted_by_full_name": schema.StringAttribute{
					Description: "Legal name of deleter of the route aggregation resource",
					Computed:    true,
				},
				"deleted_by_email": schema.StringAttribute{
					Description: "Email of deleter of the route aggregation resource",
					Computed:    true,
				},
				"deleted_date_time": schema.StringAttribute{
					Description: "Deletion time of the route aggregation resource",
					Computed:    true,
				},
			},
		},
	}
}
