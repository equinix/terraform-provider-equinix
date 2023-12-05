package equinix

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var createNetworkChangeRes = &schema.Resource{
	Schema: createNetworkChangeSch(),
}

func createNetworkChangeSch() map[string]*schema.Schema {
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
			Description: "network type",
		},
	}
}

var createNetworkOperationSchRes = &schema.Resource{
	Schema: createNetworkOperationSch(),
}

func createNetworkOperationSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"equinix_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Network operation status",
		},
	}
}

func createNetworkProjectSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"project_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Project Id",
		},
		"href": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Unique Resource URL",
		},
	}
}

func createNetworkResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Fabric Network URI information",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Fabric Network name. An alpha-numeric 24 characters string which can include only hyphens and underscores",
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
		"equinix_asn": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Equinix ASN",
		},
		"type": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"IPWAN", "EPLAN", "EVPLAN"}, true),
			Description:  "Supported Network types - EVPLAN, EPLAN, IPWAN",
		},
		"location": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Fabric Network location",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: createLocationSch(),
			},
		},
		"project": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Fabric Network project",
			Elem: &schema.Resource{
				Schema: createNetworkProjectSch(),
			},
		},
		"operation": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Network operation information that is associated with this Fabric Network",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: createNetworkOperationSch(),
			},
		},
		"change": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Change information related to this Fabric Network",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: createNetworkChangeSch(),
			},
		},
		"notifications": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Preferences for notifications on Fabric Network configuration or status changes",
			Elem: &schema.Resource{
				Schema: createNotificationSch(),
			},
		},
		"change_log": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Captures Fabric Network lifecycle change information",
			Elem: &schema.Resource{
				Schema: createChangeLogSch(),
			},
		},
	}
}
