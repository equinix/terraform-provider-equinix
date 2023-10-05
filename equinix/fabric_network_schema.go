package equinix

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var createNetworkAccountRes = &schema.Resource{
	Schema: createNetworkAccountSch(),
}

func createNetworkAccountSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"account_number": {
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
			Description: "Account Number",
		},
	}
}

var createNetworkProjectSchRes = &schema.Resource{
	Schema: createNetworkProjectSch(),
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
		"equinix_asn": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Equinix ASN",
		},
		"change_log": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Captures Fabric Network lifecycle change information",
			Elem: &schema.Resource{
				Schema: createChangeLogSch(),
			},
		},
		"type": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"IP_WAN"}, true),
			Description:  "Defines the Network type like IP_WAN",
		},
		"location": {
			Type:        schema.TypeSet,
			Required:    true,
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
		"account": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Customer account information that is associated with this Fabric Network",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: createNetworkAccountSch(),
			},
		},
		"change": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Order information related to this Fabric Network",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: createChangeSch(),
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
	}
}
