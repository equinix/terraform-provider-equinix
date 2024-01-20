package utils

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func OrderSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"purchase_order_number": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Purchase order number",
		},
		"billing_tier": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Billing tier for connection bandwidth",
		},
		"order_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Order Identification",
		},
		"order_number": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Order Reference Number",
		},
	}
}

func NotificationSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Notification Type - ALL,CONNECTION_APPROVAL,SALES_REP_NOTIFICATIONS, NOTIFICATIONS",
		},
		"send_interval": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Send interval",
		},
		"emails": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Array of contact emails",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func AccountSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"account_number": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Account Number",
		},
		"account_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Account Name",
		},
		"org_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Customer organization identifier",
		},
		"organization_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Customer organization name",
		},
		"global_org_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Global organization identifier",
		},
		"global_organization_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Global organization name",
		},
		//"ucm_id": {
		//	Type:        schema.TypeString,
		//	Computed:    true,
		//	Description: "Account ucmId",
		//},
		"global_cust_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Global Customer organization identifier",
		},
		//"reseller_account_number": {
		//	Type:        schema.TypeInt,
		//	Computed:    true,
		//	Description: "Reseller account number",
		//},
		//"reseller_account_name": {
		//	Type:        schema.TypeString,
		//	Computed:    true,
		//	Description: "Reseller account name",
		//},
		//"reseller_ucm_id": {
		//	Type:        schema.TypeString,
		//	Computed:    true,
		//	Description: "Reseller account ucmId",
		//},
		//"reseller_org_id": {
		//	Type:        schema.TypeInt,
		//	Computed:    true,
		//	Description: "Reseller customer organization identifier",
		//},
	}
}

func LocationSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"region": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Access point region",
		},
		"metro_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Access point metro name",
		},
		"metro_code": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Access point metro code",
		},
		"ibx": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "IBX Code",
		},
	}
}

func ProjectSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"project_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Project Id",
		},
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Unique Resource URL",
		},
	}
}

func ChangeLogSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"created_by": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Created by User Key",
		},
		"created_by_full_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Created by User Full Name",
		},
		"created_by_email": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Created by User Email Address",
		},
		"created_date_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Created by Date and Time",
		},
		"updated_by": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Updated by User Key",
		},
		"updated_by_full_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Updated by User Full Name",
		},
		"updated_by_email": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Updated by User Email Address",
		},
		"updated_date_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Updated by Date and Time",
		},
		"deleted_by": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Deleted by User Key",
		},
		"deleted_by_full_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Deleted by User Full Name",
		},
		"deleted_by_email": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Deleted by User Email Address",
		},
		"deleted_date_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Deleted by Date and Time",
		},
	}
}

func OperationalErrorSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"error_code": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Error  code",
		},
		"error_message": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Error Message",
		},
		"correlation_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "CorrelationId",
		},
		"details": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Details",
		},
		"help": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Help",
		},
		"additional_info": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Pricing error additional Info",
			Elem: &schema.Resource{
				Schema: ErrorAdditionalInfoSch(),
			},
		},
	}
}

func ErrorAdditionalInfoSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"property": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Property at which the error potentially occurred",
		},
		"reason": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Reason for the error",
		},
	}
}

func RedundancySch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"group": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Redundancy group identifier (UUID of primary connection)",
		},
		"priority": {
			Type:         schema.TypeString,
			Computed:     true,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"PRIMARY", "SECONDARY"}, true),
			Description:  "Connection priority in redundancy group - PRIMARY, SECONDARY",
		},
	}
}

func AccessPointTypeConfigSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Type of access point type config - VD, COLO",
		},
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Equinix-assigned access point type config identifier",
		},
		"connection_redundancy_required": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Mandate redundant connections",
		},
		"allow_bandwidth_auto_approval": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Setting to enable or disable the ability of the buyer to change connection bandwidth without approval of the seller",
		},
		"allow_remote_connections": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Setting to allow or prohibit remote connections to the service profile",
		},
		"allow_bandwidth_upgrade": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Availability of a bandwidth upgrade. The default is false",
		},
		"connection_label": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Custom name for Connection",
		},
		"enable_auto_generate_service_key": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Enable auto generate service key",
		},
		"bandwidth_alert_threshold": {
			Type:        schema.TypeFloat,
			Computed:    true,
			Description: "Percentage of port bandwidth at which an allocation alert is generated",
		},
		"allow_custom_bandwidth": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Setting to enable or disable the ability of the buyer to customize the bandwidth",
		},
		"api_config": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Api configuration details",
			Elem: &schema.Resource{
				Schema: ApiConfigSch(),
			},
		},
		"authentication_key": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Authentication key details",
			Elem: &schema.Resource{
				Schema: AuthenticationKeySch(),
			},
		},
		"link_protocol_config": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Link protocol configuration details",
			Elem: &schema.Resource{
				Schema: LinkProtocolConfigSch(),
			},
		},
		"supported_bandwidths": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Supported bandwidths",
			Elem:        &schema.Schema{Type: schema.TypeInt},
		},
	}
}

func ApiConfigSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"api_available": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Setting indicating whether the API is available (true) or not (false)",
		},
		"equinix_managed_vlan": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Setting indicating that the VLAN is managed by Equinix (true) or not (false)",
		},
		"allow_over_subscription": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Setting showing that oversubscription support is available (true) or not (false). The default is false",
		},
		"over_subscription_limit": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "A cap on over subscription",
		},
		"bandwidth_from_api": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Bandwidth from api",
		},
		"integration_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Integration id",
		},
		"equinix_managed_port": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Setting indicating that the port is managed by Equinix (true) or not (false)",
		},
	}
}

func AuthenticationKeySch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"required": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Required",
		},
		"label": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Label",
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Description",
		},
	}
}

func LinkProtocolConfigSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"encapsulation_strategy": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Encapsulation strategy",
		},
		"reuse_vlan_s_tag": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Reuse vlan sTag",
		},
		"encapsulation": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port Encapsulation",
		},
	}
}
