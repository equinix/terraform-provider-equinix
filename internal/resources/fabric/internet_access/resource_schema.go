package internetaccess

import (
	"context"

	"github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func computedString(description string) schema.StringAttribute {
	return schema.StringAttribute{
		Description: description,
		Computed:    true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

func resourcePeeringSchema(ctx context.Context, ipVersion string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: ipVersion + " peering details of the connection",
		Optional:    true,
		Computed:    true,
		CustomType:  fwtypes.NewObjectTypeOf[PeeringModel](ctx),
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.UseStateForUnknown(),
		},
		Attributes: map[string]schema.Attribute{
			"prefix": computedString(ipVersion + " prefix for the peering connection"),
			"prefix_length": schema.Int32Attribute{
				Description:   "Determines the size of the subnet",
				Computed:      true,
				PlanModifiers: []planmodifier.Int32{int32planmodifier.UseStateForUnknown()},
			},
			"equinix_peer_ip": schema.StringAttribute{
				Description:   "Peering IP address on the Equinix side. Assigned by Equinix when not set",
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"customer_peer_ip": computedString("Peering IP address on the customer side"),
			"equinix_vrrp_ip": schema.StringAttribute{
				Description:   "VRRP IP address on the Equinix side. Assigned by Equinix when not set",
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"customer_vrrp_ip": computedString("VRRP IP address on the customer side"),
		},
	}
}

func resourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: `Fabric V4 API compatible resource allows creation and management of Equinix Internet Access (EIA) services.

An EIA service is attached to one or more existing Fabric connections of type IA_VC (two for DUAL_IA).

Additional Documentation:
* API: https://docs.equinix.com/api-catalog/fabricv4/#tag/Internet-Access`,
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttributeDefaultDescription(),
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Read:   true,
				Delete: true,
			}),
			"uuid": computedString("Equinix-assigned unique identifier for the Internet Access service"),
			"href": computedString("Equinix auto generated URI to the Internet Access service"),
			"type": schema.StringAttribute{
				Description: "Type of the Internet Access service. One of: SINGLE_IA, DUAL_IA",
				Required:    true,
				Validators:  []validator.String{stringvalidator.OneOf("SINGLE_IA", "DUAL_IA")},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the Internet Access service",
				Required:    true,
			},
			"bandwidth": schema.Int32Attribute{
				Description: "Bandwidth of the service in Mbps",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"bandwidth_commit": schema.Int32Attribute{
				Description: "Minimum bandwidth commit for the burst billing variant of the service",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Description: "Provisioning state of the service. One of: PROVISIONING, PROVISIONED, FAILED, DEPROVISIONING, DEPROVISIONED",
				Computed:    true,
			},
			"use_case": computedString("Use case of the service. One of: MAIN, MANAGEMENT_ACCESS"),
			"change": schema.SingleNestedAttribute{
				Description: "Current state of the latest change on the service",
				Computed:    true,
				CustomType:  fwtypes.NewObjectTypeOf[ChangeModel](ctx),
				Attributes: map[string]schema.Attribute{
					"href": schema.StringAttribute{
						Description: "URI to the change resource",
						Computed:    true,
					},
				},
			},
			"locations": schema.ListNestedAttribute{
				Description: "Locations associated with the service",
				Computed:    true,
				CustomType:  fwtypes.NewListNestedObjectTypeOf[LocationModel](ctx),
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"metro_href": schema.StringAttribute{Description: "URI to the metro resource", Computed: true},
						"metro_code": schema.StringAttribute{Description: "Two-letter metro code", Computed: true},
						"region":     schema.StringAttribute{Description: "Region of the location. One of: AMER, EMEA, APAC", Computed: true},
						"ibx":        schema.StringAttribute{Description: "IBX data center code", Computed: true},
					},
				},
			},
			"routing_protocol": schema.SingleNestedAttribute{
				Description: "Routing protocol configuration of the service",
				Required:    true,
				CustomType:  fwtypes.NewObjectTypeOf[ResourceRoutingModel](ctx),
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "Routing protocol type. One of: BGP, DIRECT, STATIC",
						Required:    true,
						Validators:  []validator.String{stringvalidator.OneOf("BGP", "DIRECT", "STATIC")},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"customer_routes": schema.ListNestedAttribute{
						Description: "Customer routes announced through the service, given as IP blocks",
						Optional:    true,
						Computed:    true,
						CustomType:  fwtypes.NewListNestedObjectTypeOf[CustomerRouteModel](ctx),
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"ip_block": schema.SingleNestedAttribute{
									Description: "IP block of the route",
									Required:    true,
									CustomType:  fwtypes.NewObjectTypeOf[IpBlockModel](ctx),
									Attributes: map[string]schema.Attribute{
										"uuid": schema.StringAttribute{
											Description: "Unique identifier of the IP block, e.g. of an equinix_fabric_ip_block",
											Required:    true,
										},
										"href": schema.StringAttribute{
											Description: "URI to the IP block",
											Computed:    true,
										},
									},
								},
							},
						},
					},
					"connections": schema.ListNestedAttribute{
						Description: "Fabric IA_VC connections attached to the service",
						Required:    true,
						CustomType:  fwtypes.NewListNestedObjectTypeOf[ResourceConnectionModel](ctx),
						PlanModifiers: []planmodifier.List{
							listplanmodifier.RequiresReplace(),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"uuid": schema.StringAttribute{
									Description: "Unique identifier of the connection",
									Required:    true,
								},
								"href":         computedString("URI to the connection"),
								"peering_ipv4": resourcePeeringSchema(ctx, "IPv4"),
								"peering_ipv6": resourcePeeringSchema(ctx, "IPv6"),
							},
						},
					},
					"export_policy": schema.StringAttribute{
						Description: "BGP only. Export policy. One of: FULL, DEFAULT, FULL_DEFAULT, PARTIAL. Not read back from the API",
						Optional:    true,
						Validators:  []validator.String{stringvalidator.OneOf("FULL", "DEFAULT", "FULL_DEFAULT", "PARTIAL")},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"asn": schema.SingleNestedAttribute{
						Description: "BGP only. Equinix ASN resource used on the customer side. Not read back from the API",
						Optional:    true,
						CustomType:  fwtypes.NewObjectTypeOf[ResourceAsnModel](ctx),
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.RequiresReplace(),
						},
						Attributes: map[string]schema.Attribute{
							"uuid": schema.StringAttribute{
								Description: "Unique identifier of the ASN, e.g. of an equinix_fabric_asn",
								Required:    true,
							},
						},
					},
					"customer_asn": schema.Int64Attribute{
						Description: "BGP only. Customer ASN number. Not read back from the API",
						Optional:    true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.RequiresReplace(),
						},
					},
					"customer_asn_range": schema.StringAttribute{
						Description: "BGP only. Customer ASN range. One of: BITS_16, BITS_32. Not read back from the API",
						Optional:    true,
						Validators:  []validator.String{stringvalidator.OneOf("BITS_16", "BITS_32")},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"bgp_auth_key": schema.StringAttribute{
						Description: "BGP only. BGP authentication key. Not read back from the API",
						Optional:    true,
						Sensitive:   true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
				},
			},
			"billing": schema.SingleNestedAttribute{
				Description: "Billing details of the service",
				Required:    true,
				CustomType:  fwtypes.NewObjectTypeOf[ResourceBillingModel](ctx),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "Billing type. One of: FIXED, USAGE_BASED, BURST_BASED",
						Required:    true,
						Validators:  []validator.String{stringvalidator.OneOf("FIXED", "USAGE_BASED", "BURST_BASED")},
					},
					"enabled": schema.BoolAttribute{
						Description:   "Whether billing is enabled",
						Computed:      true,
						PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
					},
					"start_date": computedString("Start date of the billing period"),
				},
			},
			"account": schema.SingleNestedAttribute{
				Description: "Account the service is billed to",
				Required:    true,
				CustomType:  fwtypes.NewObjectTypeOf[ResourceAccountModel](ctx),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"account_number": schema.StringAttribute{
						Description: "Account number",
						Required:    true,
					},
					"href": computedString("URI to the account"),
				},
			},
			"project": schema.SingleNestedAttribute{
				Description: "Project the service belongs to",
				Required:    true,
				CustomType:  fwtypes.NewObjectTypeOf[ProjectModel](ctx),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"project_id": schema.StringAttribute{
						Description: "Equinix subscriber-assigned project ID",
						Required:    true,
					},
				},
			},
			"order": schema.SingleNestedAttribute{
				Description: "Order details of the service",
				Optional:    true,
				Computed:    true,
				CustomType:  fwtypes.NewObjectTypeOf[ResourceOrderModel](ctx),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"purchase_order_number": schema.StringAttribute{
						Description:   "Purchase order number",
						Optional:      true,
						Computed:      true,
						PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
					},
					"customer_reference_number": computedString("Customer reference number"),
					"billing_tier":              computedString("Billing tier for the connection bandwidth"),
					"order_id":                  computedString("Order identification"),
					"order_number":              computedString("Order reference number"),
					"term_length": schema.Int32Attribute{
						Description:   "Term length in months",
						Computed:      true,
						PlanModifiers: []planmodifier.Int32{int32planmodifier.UseStateForUnknown()},
					},
					"contracted_bandwidth": schema.Int32Attribute{
						Description:   "Contracted bandwidth",
						Computed:      true,
						PlanModifiers: []planmodifier.Int32{int32planmodifier.UseStateForUnknown()},
					},
					"href": computedString("URI to the order"),
				},
			},
			"change_log": schema.SingleNestedAttribute{
				Description: "Audit log of changes to the service",
				Computed:    true,
				CustomType:  fwtypes.NewObjectTypeOf[ChangeLogModel](ctx),
				Attributes: map[string]schema.Attribute{
					"created_date_time": schema.StringAttribute{Description: "Creation time of the service", Computed: true},
					"updated_date_time": schema.StringAttribute{Description: "Last update time of the service", Computed: true},
					"deleted_date_time": schema.StringAttribute{Description: "Deletion time of the service", Computed: true},
				},
			},
		},
	}
}
