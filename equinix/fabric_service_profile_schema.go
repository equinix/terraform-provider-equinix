package equinix

import (
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func createFabricServiceProfileSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Service Profile URI response attribute",
		},
		"type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Service profile type - L2_PROFILE, L3_PROFILE, ECIA_PROFILE, ECMC_PROFILE",
		},
		"visibility": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Service profile visibility - PUBLIC, PRIVATE",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Customer-assigned service profile name",
		},
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Equinix assigned service profile identifier",
		},
		"description": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "User-provided service description",
		},
		"notifications": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Preferences for notifications on connection configuration or status changes",
			Elem: &schema.Resource{
				Schema: equinix_schema.NotificationSch(),
			},
		},
		"access_point_type_configs": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Access point config information",
			Elem: &schema.Resource{
				Schema: createSPAccessPointTypeConfigSch(),
			},
		},
		"custom_fields": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Custom Fields",
			Elem: &schema.Resource{
				Schema: createCustomFieldSch(),
			},
		},
		"marketing_info": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Marketing Info",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: createMarketingInfoSch(),
			},
		},
		"ports": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Ports",
			Elem: &schema.Resource{
				Schema: createServiceProfileAccessPointColo(),
			},
		},
		"virtual_devices": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Virtual Devices",
			Elem: &schema.Resource{
				Schema: createServiceProfileAccessPointVd(),
			},
		},
		"allowed_emails": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Array of contact emails",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"tags": {
			Type:        schema.TypeList,
			Description: "Tags attached to the connection",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"metros": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Access point config information",
			Elem: &schema.Resource{
				Schema: createServiceMetroSch(),
			},
		},
		"self_profile": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Self Profile",
		},
		"state": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Service profile state - ACTIVE, PENDING_APPROVAL, DELETED, REJECTED",
		},
		"account": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Account",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: createSPAccountSch(),
			},
		},
		"project": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Project information",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: equinix_schema.ProjectSch(),
			},
		},
		"change_log": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Captures connection lifecycle change information",
			Elem: &schema.Resource{
				Schema: equinix_schema.ChangeLogSch(),
			},
		},
	}
}

func createCustomFieldSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"label": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Label",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Description",
		},
		"required": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "Required field",
		},
		"data_type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Data type",
		},
		"options": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Options",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"capture_in_email": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Required field",
		},
	}
}

func createMarketingInfoSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"logo": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Logo",
		},
		"promotion": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Promotion",
		},
		"process_step": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Process Step",
			Elem: &schema.Resource{
				Schema: createProcessStepSch(),
			},
		},
	}
}

func createProcessStepSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"title": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Title",
		},
		"sub_title": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Sub Title",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Description",
		},
	}
}

func createServiceProfileAccessPointColo() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Colo/Port Type",
		},
		"uuid": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Colo/Port Uuid",
		},
		"location": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Colo/Port Location",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: equinix_schema.LocationSch(),
			},
		},
		"seller_region": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Seller Region",
		},
		"seller_region_description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Seller Region details",
		},
		"cross_connect_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Cross Connect Id",
		},
	}
}

func createServiceProfileAccessPointVd() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Virtual Device Type",
		},
		"uuid": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Virtual Device Uuid",
		},
		"location": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Device Location",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: equinix_schema.LocationSch(),
			},
		},
		"interface_uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Device Interface Uuid",
		},
	}
}

func createServiceMetroSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"code": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Metro Code - Example SV",
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Metro Name",
		},
		"ibxs": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "IBX- Equinix International Business Exchange list",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"in_trail": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "In Trail",
		},
		"display_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Display Name",
		},
		"seller_regions": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "Seller Regions",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func createSPAccessPointTypeConfigSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Type of access point type config - VD, COLO",
		},
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Colo/Port Uuid",
		},
		"connection_redundancy_required": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Mandate redundant connections",
		},
		"allow_bandwidth_auto_approval": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Setting to enable or disable the ability of the buyer to change connection bandwidth without approval of the seller",
		},
		"allow_remote_connections": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Setting to allow or prohibit remote connections to the service profile",
		},
		"allow_bandwidth_upgrade": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Availability of a bandwidth upgrade. The default is false",
		},
		"connection_label": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Custom name for Connection",
		},
		"enable_auto_generate_service_key": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable auto generate service key",
		},
		"bandwidth_alert_threshold": {
			Type:        schema.TypeFloat,
			Optional:    true,
			Description: "Percentage of port bandwidth at which an allocation alert is generated",
		},
		"allow_custom_bandwidth": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Setting to enable or disable the ability of the buyer to customize the bandwidth",
		},
		"api_config": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Api configuration details",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: createApiConfigSch(),
			},
		},
		"authentication_key": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Authentication key details",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: createAuthenticationKeySch(),
			},
		},
		"link_protocol_config": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Link protocol configuration details",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: createLinkProtocolConfigSch(),
			},
		},
		"supported_bandwidths": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Supported bandwidths",
			Elem:        &schema.Schema{Type: schema.TypeInt},
		},
	}
}

var createApiConfigSchRes = &schema.Resource{
	Schema: createApiConfigSch(),
}

func createApiConfigSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"api_available": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Setting indicating whether the API is available (true) or not (false)",
		},
		"equinix_managed_vlan": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Setting indicating that the VLAN is managed by Equinix (true) or not (false)",
		},
		"allow_over_subscription": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Setting showing that oversubscription support is available (true) or not (false). The default is false",
		},
		"over_subscription_limit": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "A cap on over subscription",
		},
		"bandwidth_from_api": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Bandwidth from api",
		},
		"integration_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Integration id",
		},
		"equinix_managed_port": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Setting indicating that the port is managed by Equinix (true) or not (false)",
		},
	}
}

var createAuthenticationKeySchRes = &schema.Resource{
	Schema: createAuthenticationKeySch(),
}

func createAuthenticationKeySch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"required": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Required",
		},
		"label": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Label",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Description",
		},
	}
}

func createLinkProtocolConfigSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"encapsulation_strategy": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Encapsulation strategy",
		},
		"reuse_vlan_s_tag": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Reuse vlan sTag",
		},
		"encapsulation": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Port Encapsulation",
		},
	}
}

func createSPAccountSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"account_number": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Account Number",
		},
		"account_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Account Name",
		},
		"org_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Customer organization identifier",
		},
		"organization_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Customer organization name",
		},
		"global_org_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Global organization identifier",
		},
		"global_organization_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Global organization name",
		},
		"global_cust_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Global Customer organization identifier",
		},
		"ucm_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Enterprise datastore id",
		},
	}
}
