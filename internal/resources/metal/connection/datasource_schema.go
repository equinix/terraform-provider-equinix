package connection

import (
	"context"
	"fmt"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func dataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
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
				Description: "Connection speed - Values will be in the format '<number>Mbps' or '<number>Gpbs', for example '100Mbps`, '50Gbps', etc.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the connection resource",
				Computed:    true,
			},
			"mode": schema.StringAttribute{
				Description: fmt.Sprintf("Connection mode - %s or %s",
					string(metalv1.INTERCONNECTIONMODE_STANDARD),
					string(metalv1.INTERCONNECTIONMODE_TUNNEL),
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
			"vrfs": schema.ListAttribute{
				Description: "Attached vrfs, only in shared connection",
				Computed:    true,
				ElementType: types.StringType,
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
				Description:        "Only used with shared connection. Fabric Token required to continue the setup process with [equinix_fabric_connection](https://registry.terraform.io/providers/equinix/equinix/latest/docs/resources/fabric_connection) or from the [Equinix Fabric Portal](https://fabric.equinix.com/dashboard)",
				Computed:           true,
				DeprecationMessage: "If your organization already has connection service tokens enabled, use `service_tokens` instead",
			},
			"ports": schema.ListAttribute{
				Description: "List of connection ports - primary (`ports[0]`) and secondary (`ports[1]`)",
				CustomType:  fwtypes.NewListNestedObjectTypeOf[PortModel](ctx),
				ElementType: fwtypes.NewObjectTypeOf[PortModel](ctx),
				Computed:    true,
			},
			"service_tokens": schema.ListAttribute{
				Description: "Only used with shared connection. List of service tokens required to continue the setup process with [equinix_fabric_connection](https://registry.terraform.io/providers/equinix/equinix/latest/docs/resources/fabric_connection) or from the [Equinix Fabric Portal](https://fabric.equinix.com/dashboard)",
				CustomType:  fwtypes.NewListNestedObjectTypeOf[ServiceTokenModel](ctx),
				ElementType: fwtypes.NewObjectTypeOf[ServiceTokenModel](ctx),
				Computed:    true,
			},
			"authorization_code": schema.StringAttribute{
				Description: "Only used with Fabric Shared connection. Fabric uses this token to be able to give more detailed information about the Metal end of the network, when viewing resources from within Fabric.",
				Computed:    true,
			},
		},
	}
}
