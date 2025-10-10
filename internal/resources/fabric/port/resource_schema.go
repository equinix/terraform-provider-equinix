package port

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func resourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: `Fabric V4 API compatible resource allows creation and management of Equinix Fabric Ports

Additional Documentation:
* Getting Started: https://docs.equinix.com/en-us/Content/Interconnection/Fabric/ports/fabric-order-port.htm
* API: https://developer.equinix.com/catalog/fabricv4#operation/createPort

~> ** NOTE:** This resource is in beta and is subject to change. Please use with caution. Experimental resource may contain bugs and is not recommended for production use.
* There are no guarantees that a Port Reservation will occur after creating a port order through Terraform
* If a Port Reservation does not occur then the Port Order is not complete and the Terraform resource will not be able to be used as a dependency
* Port Deletions are not a short process and can take 2-5 business days to complete
* Please be advised that a re-run of the Terraform resource with the same settings may not result in an available Port for Reservation even if the previous one was Completed`,
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttributeDefaultDescription(),
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
			"name": schema.StringAttribute{
				Description: "Designated name of the port",
				Optional:    true,
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "Type of the port order request",
				Required:    true,
			},
			"connectivity_source_type": schema.StringAttribute{
				Description: "Connection type that is used from the port after creation",
				Required:    true,
			},
			"location": schema.SingleNestedAttribute{
				Description: "Location details for the port order",
				Required:    true,
				CustomType:  fwtypes.NewObjectTypeOf[locationModel](ctx),
				Attributes: map[string]schema.Attribute{
					"metro_code": schema.StringAttribute{
						Description: "Metro code the port will be created in",
						Required:    true,
					},
				},
			},
			"settings": schema.SingleNestedAttribute{
				Description: "Port order configuration settings",
				Required:    true,
				CustomType:  fwtypes.NewObjectTypeOf[settingsModel](ctx),
				Attributes: map[string]schema.Attribute{
					"package_type": schema.StringAttribute{
						Description: "Billing package for the port being ordered",
						Required:    true,
					},
					"shared_port_type": schema.BoolAttribute{
						Description: "Indicates whether this is a dedicated customer cage or a shared neutral cage",
						Required:    true,
					},
				},
			},
			"encapsulation": schema.SingleNestedAttribute{
				Description: "Port encapsulation settings",
				Required:    true,
				CustomType:  fwtypes.NewObjectTypeOf[encapsulationModel](ctx),
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "Port encapsulation protocol type",
						Required:    true,
					},
					"tag_protocol_id": schema.StringAttribute{
						Description: "Port encapsulation tag protocol identifier",
						Required:    true,
					},
				},
			},
			"account": schema.SingleNestedAttribute{
				Description: "Port order account details",
				Required:    true,
				CustomType:  fwtypes.NewObjectTypeOf[accountModel](ctx),
				Attributes: map[string]schema.Attribute{
					"account_number": schema.Int64Attribute{
						Description: "Account number the port will be created for",
						Required:    true,
					},
					"account_name": schema.StringAttribute{
						Computed:    true,
						Description: "Legal name of the accountholder.",
					},
					//"org_id": schema.StringAttribute{
					//	Computed:    true,
					//	Description: "Equinix-assigned ID of the subscriber's organization.",
					//},
					//"organization_name": schema.StringAttribute{
					//	Computed:    true,
					//	Description: "Equinix-assigned name of the subscriber's organization.",
					//},
					//"global_org_id": schema.StringAttribute{
					//	Computed:    true,
					//	Description: "Equinix-assigned ID of the subscriber's parent organization.",
					//},
					//"global_organization_name": schema.StringAttribute{
					//	Computed:    true,
					//	Description: "Equinix-assigned name of the subscriber's parent organization.",
					//},
					//"global_cust_id": schema.StringAttribute{
					//	Computed:    true,
					//	Description: "Equinix-assigned ID of the subscriber's parent organization.",
					//},
					"ucm_id": schema.StringAttribute{
						Computed:    true,
						Description: "Enterprise datastore id",
					},
				},
			},
			"project": schema.SingleNestedAttribute{
				Description: "Port order project details",
				Required:    true,
				CustomType:  fwtypes.NewObjectTypeOf[projectModel](ctx),
				Attributes: map[string]schema.Attribute{
					"project_id": schema.StringAttribute{
						Description: "Project id the port will be created in",
						Required:    true,
					},
				},
			},
			"redundancy": schema.SingleNestedAttribute{
				Description: "Port redundancy settings",
				Required:    true,
				CustomType:  fwtypes.NewObjectTypeOf[redundancyModel](ctx),
				Attributes: map[string]schema.Attribute{
					"priority": schema.StringAttribute{
						Description: "Port redundancy priority value",
						Required:    true,
					},
				},
			},
			"lag_enabled": schema.BoolAttribute{
				Description: "Boolean value to enable the created port with Link Aggregation Groups",
				Required:    true,
			},
			"device": schema.SingleNestedAttribute{
				Description: "Port device configuration",
				Optional:    true,
				CustomType:  fwtypes.NewObjectTypeOf[deviceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: "Device name for the port",
						Optional:    true,
					},
					"redundancy": schema.SingleNestedAttribute{
						Description: "Device redundancy configuration",
						Optional:    true,

						CustomType: fwtypes.NewObjectTypeOf[deviceRedundancyModel](ctx),
						Attributes: map[string]schema.Attribute{
							"priority": schema.StringAttribute{
								Description: "Redundancy priority (PRIMARY or SECONDARY)",
								Optional:    true,
							},
							"group": schema.StringAttribute{
								Description: "Redundancy group identifier",
								Optional:    true,
							},
						},
					},
				},
			},
			"physical_ports": schema.ListNestedAttribute{
				Description: "Physical ports that will implement this port order",
				Required:    true,
				CustomType:  fwtypes.NewListNestedObjectTypeOf[physicalPortModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Description: "Physical Port type",
							Required:    true,
						},
						"interface": schema.SingleNestedAttribute{
							Description: "Physical port interface configuration",
							Optional:    true,
							CustomType:  fwtypes.NewObjectTypeOf[interfaceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Description: "Interface type for the physical port",
									Optional:    true,
									Computed:    true,
								},
							},
						},
						"demarcation_point": schema.SingleNestedAttribute{
							Description: "Customer physical port",
							Required:    true,
							CustomType:  fwtypes.NewObjectTypeOf[demarcationPointModel](ctx),
							Attributes: map[string]schema.Attribute{
								"ibx": schema.StringAttribute{
									Description: "IBX Metro code for the physical port",
									Required:    true,
								},
								"cage_unique_space_id": schema.StringAttribute{
									Description: "Port cage unique space id",
									Required:    true,
								},
								"cabinet_unique_space_id": schema.StringAttribute{
									Description: "Port cabinet unique space id",
									Required:    true,
								},
								"patch_panel": schema.StringAttribute{
									Description: "Port patch panel",
									Required:    true,
								},
								"connector_type": schema.StringAttribute{
									Description: "Port connector type",
									Required:    true,
								},
							},
						},
					},
				},
			},
			"physical_ports_speed": schema.Int32Attribute{
				Description: "Physical Ports Speed in Mbps",
				Required:    true,
			},
			"physical_ports_type": schema.StringAttribute{
				Description: "Physical Ports Type",
				Required:    true,
			},
			"physical_ports_count": schema.Int32Attribute{
				Description: "Number of physical ports in the Port Order",
				Required:    true,
			},
			"demarcation_point_ibx": schema.StringAttribute{
				Description: "IBX code where the port will be located",
				Required:    true,
			},
			"order": schema.SingleNestedAttribute{
				Description: "Details of the Port Order such as purchaseOrder details and signature",
				Optional:    true,
				Computed:    true,
				CustomType:  fwtypes.NewObjectTypeOf[orderModel](ctx),
				Attributes: map[string]schema.Attribute{
					"purchase_order": schema.SingleNestedAttribute{
						Description: "Purchase order details",
						Optional:    true,
						CustomType:  fwtypes.NewObjectTypeOf[purchaseOrderModel](ctx),
						Attributes: map[string]schema.Attribute{
							"number": schema.StringAttribute{
								Description: "purchase order number",
								Computed:    true,
							},
							"amount": schema.StringAttribute{
								Description: "purchase order amount",
								Computed:    true,
							},
							"attachment_id": schema.StringAttribute{
								Description: "purchase order attachment id",
								Computed:    true,
							},
							"type": schema.StringAttribute{
								Description: "purchase order type",
								Computed:    true,
							},
							"start_date": schema.StringAttribute{
								Description: "purchase order start date",
								Computed:    true,
							},
							"end_date": schema.StringAttribute{
								Description: "purchase order end date",
								Computed:    true,
							},
						},
					},
					"order_number": schema.StringAttribute{
						Description: "Order Reference Number",
						Computed:    true,
					},
					"order_id": schema.StringAttribute{
						Description: "Order Identification",
						Computed:    true,
					},
					"uuid": schema.StringAttribute{
						Description: "Equinix-assigned order identifier, this is a derived response attribute",
						Computed:    true,
					},
					"customer_reference_id": schema.StringAttribute{
						Description: "Customer order reference Id",
						Optional:    true,
					},
					"signature": schema.SingleNestedAttribute{
						Description: "Port order confirmation signature details",
						Optional:    true,
						CustomType:  fwtypes.NewObjectTypeOf[signatureModel](ctx),
						Attributes: map[string]schema.Attribute{
							"signatory": schema.StringAttribute{
								Description: "Port signature Type",
								Required:    true,
							},
							"delegate": schema.SingleNestedAttribute{
								Description: "delegate order details",
								Required:    true,
								CustomType:  fwtypes.NewObjectTypeOf[delegateModel](ctx),
								Attributes: map[string]schema.Attribute{
									"first_name": schema.StringAttribute{
										Description: "First name of the signatory",
										Optional:    true,
									},
									"last_name": schema.StringAttribute{
										Description: "Last name of the signatory",
										Optional:    true,
									},
									"email": schema.StringAttribute{
										Description: "Email of the signatory",
										Required:    true,
									},
								},
							},
						},
					},
				},
			},
			"notifications": schema.ListNestedAttribute{
				Description: "List of notification types and the registered users to receive those notification types",
				Required:    true,
				CustomType:  fwtypes.NewListNestedObjectTypeOf[notificationModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Description: "Notification Type",
							Required:    true,
						},
						"registered_users": schema.ListAttribute{
							Description: "Array of registered users that will receive this notification type on the port",
							Required:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
			"additional_info": schema.ListNestedAttribute{
				Description: "List of key/value objects to provide additional context to the Port order",
				Optional:    true,
				CustomType:  fwtypes.NewListNestedObjectTypeOf[additionalInfoModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Description: "The key name of the key/value pair",
							Required:    true,
						},
						"value": schema.StringAttribute{
							Description: "The value of the key/value pair",
							Required:    true,
						},
					},
				},
			},
			"href": schema.StringAttribute{
				Description: "Equinix assigned URI of the port resource",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"uuid": schema.StringAttribute{
				Description: "Equinix assigned unique identifier of the port resource",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Description: "Value representing provisioning status for the port resource",
				Computed:    true,
			},
			"change_log": schema.SingleNestedAttribute{
				Description: "Details of the last change on the port resource",
				Computed:    true,
				CustomType:  fwtypes.NewObjectTypeOf[changeLogModel](ctx),
				Attributes: map[string]schema.Attribute{
					"created_by": schema.StringAttribute{
						Description: "User name of creator of the port resource",
						Computed:    true,
					},
					"created_by_full_name": schema.StringAttribute{
						Description: "Legal name of creator of the port resource",
						Computed:    true,
					},
					"created_by_email": schema.StringAttribute{
						Description: "Email of creator of the port resource",
						Computed:    true,
					},
					"created_date_time": schema.StringAttribute{
						Description: "Creation time of the port resource",
						Computed:    true,
					},
					"updated_by": schema.StringAttribute{
						Description: "User name of last updater of the port resource",
						Computed:    true,
					},
					"updated_by_full_name": schema.StringAttribute{
						Description: "Legal name of last updater of the port resource",
						Computed:    true,
					},
					"updated_by_email": schema.StringAttribute{
						Description: "Email of last updater of the port resource",
						Computed:    true,
					},
					"updated_date_time": schema.StringAttribute{
						Description: "Last update time of the port resource",
						Computed:    true,
					},
					"deleted_by": schema.StringAttribute{
						Description: "User name of deleter of the port resource",
						Computed:    true,
					},
					"deleted_by_full_name": schema.StringAttribute{
						Description: "Legal name of deleter of the port resource",
						Computed:    true,
					},
					"deleted_by_email": schema.StringAttribute{
						Description: "Email of deleter of the port resource",
						Computed:    true,
					},
					"deleted_date_time": schema.StringAttribute{
						Description: "Deletion time of the port resource",
						Computed:    true,
					},
				},
			},
		},
	}
}
