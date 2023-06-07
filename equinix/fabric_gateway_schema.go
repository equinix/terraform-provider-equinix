package equinix

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func createPackageSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"code": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Fabric Gateway package code",
		},
	}
}

var createPackageRes = &schema.Resource{
	Schema: createPackageSch(),
}

var createFgAccountRes = &schema.Resource{
	Schema: createFgAccountSch(),
}

func createFgAccountSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"account_number": {
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
			Description: "Account Number",
		},
	}
}

var createFabricGatewayProjectSchRes = &schema.Resource{
	Schema: createGatewayProjectSch(),
}

func createFabricGatewayProjectSch() map[string]*schema.Schema {
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

func createFabricGatewayResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Fabric Gateway URI information",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Fabric Gateway name. An alpha-numeric 24 characters string which can include only hyphens and underscores",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Customer-provided Fabric Gateway description",
		},
		"state": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Fabric Gateway overall state",
		},
		"equinix_asn": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Equinix ASN",
		},
		"package": {
			Type:        schema.TypeSet,
			Required:    true,
			Description: "Fabric Gateway location",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: createPackageSch(),
			},
		},
		"change_log": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Captures Fabric Gateway lifecycle change information",
			Elem: &schema.Resource{
				Schema: createChangeLogSch(),
			},
		},
		"type": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"XF_GATEWAY"}, true),
			Description:  "Defines the FG type like XF_GATEWAY",
		},
		"location": {
			Type:        schema.TypeSet,
			Required:    true,
			Description: "Fabric Gateway location",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: createLocationSch(),
			},
		},
		"project": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Fabric Gateway project",
			Elem: &schema.Resource{
				Schema: createFabricGatewayProjectSch(),
			},
		},
		"account": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Customer account information that is associated with this Fabric Gateway",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: createFgAccountSch(),
			},
		},
		"order": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Order information related to this Fabric Gateway",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: createOrderSch(),
			},
		},
		"notifications": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Preferences for notifications on Fabric Gateway configuration or status changes",
			Elem: &schema.Resource{
				Schema: createNotificationSch(),
			},
		},
	}
}
