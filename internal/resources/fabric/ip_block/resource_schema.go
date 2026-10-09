package ipblock

import (
	"context"

	"github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func resourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: `Fabric V4 API compatible resource allows creation and management of Equinix Fabric IP Blocks.

Additional Documentation:
* API: https://docs.equinix.com/api-catalog/fabricv4/#tag/IP-Blocks`,
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttributeDefaultDescription(),
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Delete: true,
			}),
			"type": schema.StringAttribute{
				Description: "IP block type; IPV4_IP_BLOCK or IPV6_IP_BLOCK",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project": schema.SingleNestedAttribute{
				Description: "Project this IP block belongs to",
				Required:    true,
				CustomType:  fwtypes.NewObjectTypeOf[ResourceProjectModel](ctx),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"project_id": schema.StringAttribute{
						Description: "Equinix-assigned project ID",
						Required:    true,
					},
				},
			},
			"location": schema.SingleNestedAttribute{
				Description: "Metro location for the IP block",
				Optional:    true,
				Computed:    true,
				CustomType:  fwtypes.NewObjectTypeOf[LocationModel](ctx),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"metro_code": schema.StringAttribute{
						Description: "Two-letter metro code (e.g. SV, DC, AM)",
						Required:    true,
					},
					"metro_href": schema.StringAttribute{
						Description: "Metro resource URL",
						Computed:    true,
					},
				},
			},
			"prefix_length": schema.Int32Attribute{
				Description: "CIDR prefix length (e.g. 28 for a /28 block)",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"prefix": schema.StringAttribute{
				Description: "CIDR prefix (e.g. 192.0.2.0/28)",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"account": schema.SingleNestedAttribute{
				Description: "Billing account for the IP block",
				Optional:    true,
				Computed:    true,
				CustomType:  fwtypes.NewObjectTypeOf[AccountModel](ctx),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"account_number": schema.StringAttribute{
						Description: "Account number",
						Optional:    true,
						Computed:    true,
					},
				},
			},
			"uuid": schema.StringAttribute{
				Description: "Equinix-assigned UUID of the IP block",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"href": schema.StringAttribute{
				Description: "URL of the IP block resource",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Description: "Provisioning state of the IP block (PENDING, ACTIVE, DELETING, DELETED, FAILED)",
				Computed:    true,
			},
			"ownership": schema.StringAttribute{
				Description: "Ownership type; EQUINIX or CUSTOMER",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"assets": schema.ListNestedAttribute{
				Description: "Resources currently using this IP block",
				Computed:    true,
				CustomType:  fwtypes.NewListNestedObjectTypeOf[AssetModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Description: "Asset type",
							Computed:    true,
						},
						"uuid": schema.StringAttribute{
							Description: "Asset UUID",
							Computed:    true,
						},
						"href": schema.StringAttribute{
							Description: "Asset URL",
							Computed:    true,
						},
					},
				},
			},
			"change": schema.SingleNestedAttribute{
				Description: "Current pending change on this IP block",
				Computed:    true,
				CustomType:  fwtypes.NewObjectTypeOf[ChangeModel](ctx),
				Attributes: map[string]schema.Attribute{
					"href": schema.StringAttribute{
						Description: "Change URL",
						Computed:    true,
					},
				},
			},
			"order": schema.SingleNestedAttribute{
				Description: "Order details for this IP block",
				Computed:    true,
				CustomType:  fwtypes.NewObjectTypeOf[OrderModel](ctx),
				Attributes: map[string]schema.Attribute{
					"href": schema.StringAttribute{
						Description: "Order URL",
						Computed:    true,
					},
					"order_number": schema.StringAttribute{
						Description: "Order number",
						Computed:    true,
					},
				},
			},
			"change_log": schema.SingleNestedAttribute{
				Description: "Timestamps tracking resource lifecycle",
				Computed:    true,
				CustomType:  fwtypes.NewObjectTypeOf[ChangeLogModel](ctx),
				Attributes: map[string]schema.Attribute{
					"created_date_time": schema.StringAttribute{
						Description: "Creation timestamp",
						Computed:    true,
					},
					"updated_date_time": schema.StringAttribute{
						Description: "Last update timestamp",
						Computed:    true,
					},
					"deleted_date_time": schema.StringAttribute{
						Description: "Deletion timestamp",
						Computed:    true,
					},
				},
			},
		},
	}
}
