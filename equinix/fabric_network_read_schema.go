package equinix

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func readChangeSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "href",
		},
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "UUID of Network",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "network change type",
		},
	}
}

func readNetworkResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Equinix-assigned Fabric Network identifier",
		},
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Fabric Network URI information",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Fabric Network name. An alpha-numeric 24 characters string which can include only hyphens and underscores",
		},
		"state": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Fabric Network overall state",
		},
		"scope": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Network scope",
		},
		"connections_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Connections count",
		},
		"change_log": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Captures Fabric Network lifecycle change information",
			Elem: &schema.Resource{
				Schema: readChangeLogSch(),
			},
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Defines the Fabric Network type like IPWAN",
		},
		"location": {
			Type:        schema.TypeSet,
			Computed:    true,
			Optional:    true,
			Description: "Fabric Network location",
			Elem: &schema.Resource{
				Schema: readLocationSch(),
			},
		},
		"project": {
			Type:        schema.TypeSet,
			Optional:    true,
			Computed:    true,
			Description: "Project information",
			Elem: &schema.Resource{
				Schema: createGatewayProjectSch(),
			},
		},
		"change": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Order information related to this Fabric Network",
			Elem: &schema.Resource{
				Schema: readChangeSch(),
			},
		},
		"notifications": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Preferences for notifications on Fabric Network configuration or status changes",
			Elem: &schema.Resource{
				Schema: readNotificationSch(),
			},
		},
	}
}
