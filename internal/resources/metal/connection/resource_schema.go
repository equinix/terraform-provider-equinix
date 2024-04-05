package connection

import (
	"context"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func resourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttributeDefaultDescription(),
			"name": schema.StringAttribute{
				Description: "Name of the connection resource",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"facility": schema.StringAttribute{
				Description:        "Facility where the connection will be created",
				Optional:           true,
				Computed:           true,
				DeprecationMessage: "Use metro instead of facility. For more information, read the migration guide.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("metro"),
					}...),
				},
			},
			"metro": schema.StringAttribute{
				Description: "Metro where the connection will be created",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"redundancy": schema.StringAttribute{
				Description: "Connection redundancy - redundant or primary",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						string(metalv1.INTERCONNECTIONREDUNDANCY_REDUNDANT),
						string(metalv1.INTERCONNECTIONREDUNDANCY_PRIMARY),
					),
				},
			},
			"contact_email": schema.StringAttribute{
				Description: "The preferred email used for communication and notifications about the Equinix Fabric interconnection",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Description: "Connection type - dedicated, shared or shared_port_vlan",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						string(metalv1.INTERCONNECTIONTYPE_DEDICATED),
						string(metalv1.INTERCONNECTIONTYPE_SHARED),
						string(metalv1.INTERCONNECTIONTYPE_SHARED_PORT_VLAN),
					),
				},
			},
			"project_id": schema.StringAttribute{
				Description: "ID of the project where the connection is scoped to. Required with type \"shared\"",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"speed": schema.StringAttribute{
				Description: "Connection speed -  Values must be in the format '<number>Mbps' or '<number>Gpbs', for example '100Mbps' or '50Gbps'.  Actual supported values will depend on the connection type and whether the connection uses VLANs or VRF.",
				Optional:    true,
				Computed:    true,
			},
			// TODO(ocobles) workaround - API returns "" when description was not provided. Computed=true will
			// ensure backward compatibility since that change was tolerated in SDKv2
			"description": schema.StringAttribute{
				Description: "Description of the connection resource",
				Optional:    true,
				Computed:    true,
			},
			"mode": schema.StringAttribute{
				Description: "Mode for connections in IBX facilities with the dedicated type - standard or tunnel",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(string(metalv1.INTERCONNECTIONMODE_STANDARD)),
				Validators: []validator.String{
					stringvalidator.OneOf(
						string(metalv1.INTERCONNECTIONMODE_STANDARD),
						string(metalv1.INTERCONNECTIONMODE_TUNNEL),
					),
				},
			},
			"tags": schema.ListAttribute{
				Description: "Tags attached to the connection",
				Optional:    true,
				ElementType: types.StringType,
			},
			"vlans": schema.ListAttribute{
				Description: "Only used with shared connection. VLANs to attach. Pass one vlan for Primary/Single connection and two vlans for Redundant connection",
				Optional:    true,
				ElementType: types.Int64Type,
				Validators: []validator.List{
					listvalidator.SizeAtMost(2),
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"vrfs": schema.ListAttribute{
				Description: "Only used with shared connection. VRFs to attach. Pass one VRF for Primary/Single connection and two VRFs for Redundant connection",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.SizeAtMost(2),
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"service_token_type": schema.StringAttribute{
				Description: "Only used with shared connection. Type of service token to use for the connection, a_side or z_side",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("a_side", "z_side"),
				},
			},
			"organization_id": schema.StringAttribute{
				Description: "ID of the organization responsible for the connection. Applicable with type \"dedicated\"",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(path.Expressions{
						path.MatchRoot("project_id"),
					}...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Description: "Status of the connection resource",
				Computed:    true,
			},
			"token": schema.StringAttribute{
				Description: "Only used with shared connection. Fabric Token required to continue the setup process with [equinix_ecx_l2_connection](https://registry.terraform.io/providers/equinix/equinix/latest/docs/resources/equinix_ecx_l2_connection) or from the [Equinix Fabric Portal](https://ecxfabric.equinix.com/dashboard)",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				DeprecationMessage: "If your organization already has connection service tokens enabled, use `service_tokens` instead",
			},
			"ports": schema.ListAttribute{
				Description: "List of connection ports - primary (`ports[0]`) and secondary (`ports[1]`)",
				CustomType:  fwtypes.NewListNestedObjectTypeOf[PortModel](ctx),
				ElementType: fwtypes.NewObjectTypeOf[PortModel](ctx),
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"service_tokens": schema.ListAttribute{
				Description: "Only used with shared connection. List of service tokens required to continue the setup process with [equinix_ecx_l2_connection](https://registry.terraform.io/providers/equinix/equinix/latest/docs/resources/equinix_ecx_l2_connection) or from the [Equinix Fabric Portal](https://ecxfabric.equinix.com/dashboard)",
				CustomType:  fwtypes.NewListNestedObjectTypeOf[ServiceTokenModel](ctx),
				ElementType: fwtypes.NewObjectTypeOf[ServiceTokenModel](ctx),
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"authorization_code": schema.StringAttribute{
				Description: "Only used with Fabric Shared connection. Fabric uses this token to be able to give more detailed information about the Metal end of the network, when viewing resources from within Fabric.",
				Computed:    true,
			},
		},
	}
}
