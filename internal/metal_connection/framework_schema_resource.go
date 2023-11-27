package metal_connection

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/packethost/packngo"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var metalConnectionResourceSchema = schema.Schema{
    Attributes: map[string]schema.Attribute{
        "id": schema.StringAttribute{
            Description: "The unique identifier for this Metal Connection",
            Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
        },
        "name": schema.StringAttribute{
            Required:    true,
            Description: "Name of the connection resource",
            PlanModifiers: []planmodifier.String{
                stringplanmodifier.RequiresReplace(),
            },
        },
        "facility": schema.StringAttribute{
            Optional:    true,
            Computed:    true,
            Description: "Facility where the connection will be created",
            DeprecationMessage: "Use metro instead of facility. For more information, read the migration guide.",
            PlanModifiers: []planmodifier.String{
                stringplanmodifier.RequiresReplace(),
            },
            Validators: []validator.String{
                stringvalidator.ExactlyOneOf(path.Expressions{
					path.MatchRoot("metro"),
				}...),
            },
        },
        "metro": schema.StringAttribute{
            Optional:    true,
            Computed:    true,
            Description: "Metro where the connection will be created",
            PlanModifiers: []planmodifier.String{
                stringplanmodifier.RequiresReplace(),
            },
            // TODO (ocobles)
            // StateFunc: toLower
            // StateFunc doesn't exist in terraform, it requires implementation of bespoke logic before storing state, for instance in resource Create method
        },
        "redundancy": schema.StringAttribute{
            Required:    true,
            Description: "Connection redundancy - redundant or primary",
            Validators: []validator.String{
                stringvalidator.OneOf(
                    string(packngo.ConnectionRedundant),
                    string(packngo.ConnectionPrimary),
                ),
            },
        },
        "contact_email": schema.StringAttribute{
            Optional:    true,
            Computed:    true,
            Description: "The preferred email used for communication and notifications about the Equinix Fabric interconnection",
            PlanModifiers: []planmodifier.String{
                stringplanmodifier.RequiresReplace(),
            },
        },
        "type": schema.StringAttribute{
            Required:    true,
            Description: "Connection type - dedicated or shared",
            PlanModifiers: []planmodifier.String{
                stringplanmodifier.RequiresReplace(),
            },
            Validators: []validator.String{
                stringvalidator.OneOf(
                    string(packngo.ConnectionDedicated),
                    string(packngo.ConnectionShared),
                ),
            },
        },
        "project_id": schema.StringAttribute{
            Optional:    true,
            Description: "ID of the project where the connection is scoped to. Required with type \"shared\"",
            PlanModifiers: []planmodifier.String{
                stringplanmodifier.RequiresReplace(),
                stringplanmodifier.UseStateForUnknown(),
            },
        },
        "speed": schema.StringAttribute{
            Optional:    true,
            Computed:    true,
            Description: "Port speed. Required for a_side connections",
        },
        "description": schema.StringAttribute{
            Optional:    true,
            Description: "Description of the connection resource",
        },
        "mode": schema.StringAttribute{
            Optional:    true,
            Description: "Mode for connections in IBX facilities with the dedicated type - standard or tunnel",
            Default: stringdefault.StaticString(string(packngo.ConnectionModeStandard)),
            Validators: []validator.String{
                stringvalidator.OneOf(
                    string(packngo.ConnectionModeStandard),
                    string(packngo.ConnectionModeTunnel),
                ),
            },
        },
        "tags": schema.ListAttribute{
            Computed: true,
            Description: "Tags attached to the connection",
            ElementType: types.StringType,
        },
        "vlans": schema.ListAttribute{
            Computed: true,
            Description:  "Only used with shared connection. VLANs to attach. Pass one vlan for Primary/Single connection and two vlans for Redundant connection",
            ElementType: types.Int64Type,
            Validators: []validator.List{
                listvalidator.SizeAtMost(2),
            },
        },
        "service_token_type": schema.StringAttribute{
            Optional:    true,
            Description: "Only used with shared connection. Type of service token to use for the connection, a_side or z_side",
            Validators: []validator.String{
				stringvalidator.OneOf("a_side", "z_side"),
			},
        },
        "organization_id": schema.StringAttribute{
            Optional:    true,
            Description: "ID of the organization responsible for the connection. Applicable with type \"dedicated\"",
            Default: stringdefault.StaticString("standard"),
            Validators: []validator.String{
				stringvalidator.AtLeastOneOf(path.Expressions{
					path.MatchRoot("project_id"),
				}...),
			},
        },
        "status": schema.StringAttribute{
            Computed:    true,
            Description: "Status of the connection resource",
        },
        "token": schema.StringAttribute{
            Computed:    true,
            Description: "Only used with shared connection. Fabric Token required to continue the setup process with [equinix_ecx_l2_connection](https://registry.terraform.io/providers/equinix/equinix/latest/docs/resources/equinix_ecx_l2_connection) or from the [Equinix Fabric Portal](https://ecxfabric.equinix.com/dashboard)",
            DeprecationMessage: "If your organization already has connection service tokens enabled, use `service_tokens` instead",
        },
        "service_tokens": schema.ListAttribute{
            Computed: true,
            Description: "List of service tokens required to continue the setup process",
            ElementType: ServiceTokensObjectType,
        },
        "ports": schema.ListAttribute{
            Computed: true,
            Description: "List of connection ports",
            ElementType: PortsObjectType,
        },
    },
}

var PortsObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":                  types.StringType,
		"name":                types.StringType,
		"role":                types.StringType,
		"speed":               types.StringType,
		"status":              types.StringType,
		"link_status":         types.StringType,
		"virtual_circuit_ids": types.ListType{ElemType: types.StringType},
	},
}

var ServiceTokensObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":                types.StringType,
		"max_allowed_speed": types.StringType,
		"role":              types.StringType,
		"state":             types.StringType,
		"type":              types.StringType,
	},
}