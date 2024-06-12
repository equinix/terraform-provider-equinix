package network

import (
	equinix_fabric_schema "github.com/equinix/terraform-provider-equinix/internal/fabric/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func fabricNetworkChangeSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Absolute URL that returns the details of the given change.\nExample: https://api.equinix.com/fabric/v4/networks/92dc376a-a932-43aa-a6a2-c806dedbd784",
		},
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Asset change request identifier.",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Asset instance change request type.: NETWORK_CREATION, NETWORK_UPDATE, NETWORK_DELETION",
		},
	}
}
func fabricNetworkOperationSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"equinix_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Progress towards provisioning a given asset.",
		},
	}
}
func fabricNetworkProjectSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"project_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Customer project identifier",
		},
	}
}
func fabricNetworkResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Fabric Network URI information",
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(1, 24),
			Description:  "Fabric Network name. An alpha-numeric 24 characters string which can include only hyphens and underscores",
		},
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Equinix-assigned network identifier",
		},
		"state": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Fabric Network overall state",
		},
		"scope": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Fabric Network scope",
		},
		"type": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"IPWAN", "EPLAN", "EVPLAN"}, true),
			Description:  "Supported Network types - EVPLAN, EPLAN, IPWAN",
		},
		"location": {
			Type:        schema.TypeSet,
			Computed:    true,
			Optional:    true,
			Description: "Fabric Network location",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: equinix_fabric_schema.LocationSch(),
			},
		},
		"project": {
			Type:        schema.TypeSet,
			Required:    true,
			Description: "Fabric Network project",
			Elem: &schema.Resource{
				Schema: fabricNetworkProjectSch(),
			},
		},
		"operation": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Network operation information that is associated with this Fabric Network",
			Elem: &schema.Resource{
				Schema: fabricNetworkOperationSch(),
			},
		},
		"change": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Information on asset change operation",
			Elem: &schema.Resource{
				Schema: fabricNetworkChangeSch(),
			},
		},
		"notifications": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Preferences for notifications on Fabric Network configuration or status changes",
			Elem: &schema.Resource{
				Schema: equinix_fabric_schema.NotificationSch(),
			},
		},
		"change_log": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "A permanent record of asset creation, modification, or deletion",
			Elem: &schema.Resource{
				Schema: equinix_fabric_schema.ChangeLogSch(),
			},
		},
		"connections_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Number of connections associated with this network",
		},
	}
}
