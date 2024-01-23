package metal_connection

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
)

func serviceTokenSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "ID of the service token",
				Computed:    true,
			},
			"expires_at": {
				Type:        schema.TypeString,
				Description: "Expiration date of the service token",
				Computed:    true,
			},
			"max_allowed_speed": {
				Type:        schema.TypeString,
				Description: "Maximum allowed speed for the service token",
				Computed:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "Type of the service token, a_side or z_side",
				Computed:    true,
			},
			"state": {
				Type:        schema.TypeString,
				Description: "State of the service token",
				Computed:    true,
			},
			"role": {
				Type:        schema.TypeString,
				Description: "Role of the service token",
				Computed:    true,
			},
		},
	}
}

func connectionPortSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the connection port resource",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the connection port resource",
			},
			"role": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Role - primary or secondary",
			},
			"speed": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Port speed in bits per second",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Port status",
			},
			"link_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Port link status",
			},
			"virtual_circuit_ids": {
				Computed:    true,
				Type:        schema.TypeList,
				Elem:        schema.TypeString,
				Description: "List of IDs of virtual circuits attached to this port",
			},
		},
	}
}

func DataSource() *schema.Resource {
	speeds := []string{}
	for _, allowedSpeed := range allowedSpeeds {
		speeds = append(speeds, allowedSpeed.Str)
	}
	return &schema.Resource{
		Read: dataSourceMetalConnectionRead,

		Schema: map[string]*schema.Schema{
			"connection_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the connection to lookup",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the connection resource",
			},
			"contact_email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The preferred email used for communication and notifications about the Equinix Fabric interconnection",
			},
			"facility": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Facility which the connection is scoped to",
				Deprecated:  "Use metro instead of facility.  For more information, read the migration guide: https://registry.terraform.io/providers/equinix/equinix/latest/docs/guides/migration_guide_facilities_to_metros_devices",
			},
			"metro": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Metro which the connection is scoped to",
			},
			"redundancy": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Connection redundancy - redundant or primary",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Connection type - dedicated or shared",
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of project to which the connection belongs",
			},
			"speed": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: fmt.Sprintf("Port speed. Possible values are %s", strings.Join(speeds, ", ")),
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the connection resource",
			},
			"mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Connection mode - standard or tunnel",
			},
			"tags": {
				Type:        schema.TypeList,
				Description: "Tags attached to the connection",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"vlans": {
				Type:        schema.TypeList,
				Description: "Attached vlans, only in shared connection",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"service_token_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Only used with shared connection. Type of service token to use for the connection, a_side or z_side",
			},
			"organization_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of organization to which the connection is scoped to",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the connection resource",
			},
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Only used with shared connection. Fabric Token required to continue the setup process with [equinix_ecx_l2_connection](https://registry.terraform.io/providers/equinix/equinix/latest/docs/resources/equinix_ecx_l2_connection) or from the [Equinix Fabric Portal](https://ecxfabric.equinix.com/dashboard)",
				Deprecated:  "If your organization already has connection service tokens enabled, use `service_tokens` instead",
			},
			"ports": {
				Type:        schema.TypeList,
				Elem:        connectionPortSchema(),
				Computed:    true,
				Description: "List of connection ports - primary (`ports[0]`) and secondary (`ports[1]`)",
			},
			"service_tokens": {
				Type:        schema.TypeList,
				Description: "Only used with shared connection. List of service tokens required to continue the setup process with [equinix_ecx_l2_connection](https://registry.terraform.io/providers/equinix/equinix/latest/docs/resources/equinix_ecx_l2_connection) or from the [Equinix Fabric Portal](https://ecxfabric.equinix.com/dashboard)",
				Computed:    true,
				Elem:        serviceTokenSchema(),
			},
		},
	}
}

func getConnectionPorts(cps []packngo.ConnectionPort) []map[string]interface{} {
	ret := make([]map[string]interface{}, len(cps))
	order := map[packngo.ConnectionPortRole]int{
		packngo.ConnectionPortPrimary:   0,
		packngo.ConnectionPortSecondary: 1,
	}

	for _, p := range cps {
		vcIDs := []string{}
		for _, vc := range p.VirtualCircuits {
			vcIDs = append(vcIDs, vc.ID)
		}
		connPort := map[string]interface{}{
			"name":                p.Name,
			"id":                  p.ID,
			"role":                string(p.Role),
			"speed":               p.Speed,
			"status":              p.Status,
			"link_status":         p.LinkStatus,
			"virtual_circuit_ids": vcIDs,
		}
		// sort the ports by role, asserting the API always returns primary for len of 1 responses
		ret[order[p.Role]] = connPort
	}
	return ret
}

func getConnectionVlans(conn *packngo.Connection) []int {
	var ret []int

	if conn.Type == packngo.ConnectionShared {
		order := map[packngo.ConnectionPortRole]int{
			packngo.ConnectionPortPrimary:   0,
			packngo.ConnectionPortSecondary: 1,
		}

		rawVlans := make([]int, len(conn.Ports))
		for _, p := range conn.Ports {
			rawVlans[order[p.Role]] = p.VirtualCircuits[0].VNID
		}

		for _, v := range rawVlans {
			if v > 0 {
				ret = append(ret, v)
			}
		}
	}

	return ret
}

func dataSourceMetalConnectionRead(d *schema.ResourceData, meta interface{}) error {
	connId := d.Get("connection_id").(string)
	d.SetId(connId)
	return resourceMetalConnectionRead(d, meta)
}

func getServiceTokens(tokens []packngo.FabricServiceToken) ([]map[string]interface{}, error) {
	tokenList := []map[string]interface{}{}
	for _, token := range tokens {
		speed, err := speedUintToStr(token.MaxAllowedSpeed)
		if err != nil {
			return nil, err
		}

		rawToken := map[string]interface{}{
			"id":                token.ID,
			"max_allowed_speed": speed,
			"role":              string(token.Role),
			"state":             token.State,
			"type":              string(token.ServiceTokenType),
		}
		if token.ExpiresAt != nil {
			rawToken["expires_at"] = token.ExpiresAt.String()
		}
		tokenList = append(tokenList, rawToken)
	}
	return tokenList, nil
}
