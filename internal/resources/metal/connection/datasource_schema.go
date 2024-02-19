package connection

import (
	"fmt"

	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/packethost/packngo"
)

var dataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": framework.IDAttributeDefaultDescription(),
		"connection_id": schema.StringAttribute{
			Description: "ID of the connection to lookup",
			Required:    true,
		},
		"name": schema.StringAttribute{
			Description: "Name of the connection resource",
			Computed:    true,
		},
		"facility": schema.StringAttribute{
			Description:        "Facility which the connection is scoped to",
			Computed:           true,
			DeprecationMessage: "Use metro instead of facility. For more information, read the migration guide.",
		},
		"metro": schema.StringAttribute{
			Description: "Metro which the connection is scoped to",
			Computed:    true,
		},
		"redundancy": schema.StringAttribute{
			Description: "Connection redundancy - redundant or primary",
			Computed:    true,
		},
		"contact_email": schema.StringAttribute{
			Description: "The preferred email used for communication and notifications about the Equinix Fabric interconnection",
			Computed:    true,
		},
		"type": schema.StringAttribute{
			Description: "Connection type - dedicated or shared",
			Computed:    true,
		},
		"project_id": schema.StringAttribute{
			Description: "ID of project to which the connection belongs",
			Computed:    true,
		},
		"speed": schema.StringAttribute{
			Description: fmt.Sprintf("Port speed. Possible values are %s", allowedSpeedsString()),
			Computed:    true,
		},
		"description": schema.StringAttribute{
			Description: "Description of the connection resource",
			Computed:    true,
		},
		"mode": schema.StringAttribute{
			Description: fmt.Sprintf("Connection mode - %s or %s",
				string(packngo.ConnectionModeStandard),
				string(packngo.ConnectionModeTunnel),
			),
			Computed: true,
		},
		"tags": schema.ListAttribute{
			Description: "Tags attached to the connection",
			Computed:    true,
			ElementType: types.StringType,
		},
		"vlans": schema.ListAttribute{
			Description: "Attached vlans, only in shared connection",
			Computed:    true,
			ElementType: types.Int64Type,
		},
		"service_token_type": schema.StringAttribute{
			Description: "Only used with shared connection. Type of service token to use for the connection, a_side or z_side",
			Computed:    true,
		},
		"organization_id": schema.StringAttribute{
			Description: "ID of organization to which the connection is scoped to",
			Computed:    true,
		},
		"status": schema.StringAttribute{
			Description: "Status of the connection resource",
			Computed:    true,
		},
		"token": schema.StringAttribute{
			Description:        "Only used with shared connection. Fabric Token required to continue the setup process with [equinix_ecx_l2_connection](https://registry.terraform.io/providers/equinix/equinix/latest/docs/resources/equinix_ecx_l2_connection) or from the [Equinix Fabric Portal](https://ecxfabric.equinix.com/dashboard)",
			Computed:           true,
			DeprecationMessage: "If your organization already has connection service tokens enabled, use `service_tokens` instead",
		},
	},
	Blocks: map[string]schema.Block{
		"service_tokens": schema.ListNestedBlock{
			Description: "Only used with shared connection. List of service tokens required to continue the setup process with [equinix_ecx_l2_connection](https://registry.terraform.io/providers/equinix/equinix/latest/docs/resources/equinix_ecx_l2_connection) or from the [Equinix Fabric Portal](https://ecxfabric.equinix.com/dashboard)",
			NestedObject: schema.NestedBlockObject{
				Attributes: serviceTokensDataSourceNestedAttribute,
			},
		},
		"ports": schema.ListNestedBlock{
			Description: "List of connection ports - primary (`ports[0]`) and secondary (`ports[1]`)",
			NestedObject: schema.NestedBlockObject{
				Attributes: portsDataSourceNestedAttribute,
			},
		},
	},
}

var portsDataSourceNestedAttribute = map[string]schema.Attribute{
	"id": framework.IDAttribute("ID of the connection port resource"),
	"name": schema.StringAttribute{
		Description: "Name of the connection port resource",
		Computed:    true,
	},
	"role": schema.StringAttribute{
		Description: "Role - primary or secondary",
		Computed:    true,
	},
	"speed": schema.Int64Attribute{
		Description: "Port speed in bits per second",
		Computed:    true,
	},
	"status": schema.StringAttribute{
		Description: "Port status",
		Computed:    true,
	},
	"link_status": schema.StringAttribute{
		Description: "Port link status",
		Computed:    true,
	},
	"virtual_circuit_ids": schema.ListAttribute{
		Description: "List of IDs of virtual circuits attached to this port",
		Computed:    true,
		ElementType: types.StringType,
	},
}

var serviceTokensDataSourceNestedAttribute = map[string]schema.Attribute{
	"id": framework.IDAttribute("ID of the service token"),
	"expires_at": schema.StringAttribute{
		Description: "Expiration date of the service token",
		Computed:    true,
	},
	"max_allowed_speed": schema.StringAttribute{
		Description: "Maximum allowed speed for the service token",
		Computed:    true,
	},
	"type": schema.StringAttribute{
		Description: "Type of the service token, a_side or z_side",
		Computed:    true,
	},
	"state": schema.StringAttribute{
		Description: "State of the service token",
		Computed:    true,
	},
	"role": schema.StringAttribute{
		Description: "Role of the service token",
		Computed:    true,
	},
}
