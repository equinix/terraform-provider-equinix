package equinix

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func readServiceTokenSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Token type - VC_TOKEN",
		},
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "An absolute URL that is the subject of the link's context",
		},
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Equinix-assigned service token identifier",
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Service token description",
		},
	}
}

func readLocationSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"region": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Access point region",
		},
		"metro_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Access point metro name",
		},
		"metro_code": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Access point metro code",
		},
		"ibx": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "IBX Code",
		},
	}
}

func readVirtualGatewaySch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Gateway unique identifier",
		},
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Unique Resource Identifier",
		},
		"state": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Virtual Gateway State",
		},
	}
}

func readServiceProfileSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Service Profile URI response attribute",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Service profile type- LAYER_2_PROFILE, LAYER_3_PROFILE",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Customer-assigned service profile name",
		},
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Equinix assigned service profile identifier",
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "User-provided service profile description",
		},
		"access_point_type_configs": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Access point config information",
			Elem: &schema.Resource{
				Schema: readAccessPointTypeConfigSch(),
			},
		},
	}
}

func readAccessPointTypeConfigSch() map[string]*schema.Schema {
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
	}
}

func readAccessPointLinkProtocolSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Type of the link protocol - DOT1Q, QINQ, UNTAGGED",
		},
		"vlan_tag": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Vlan Tag information, vlanTag value specified for DOT1Q connections",
		},
		"vlan_s_tag": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Vlan Provider Tag information, vlanSTag value specified for QINQ connections",
		},
		"vlan_c_tag": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Vlan Customer Tag information, vlanCTag value specified for QINQ connections",
		},
	}
}

func readAccessPointVirtualDeviceSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Unique Resource Identifier",
		},
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Equinix-assigned Virtual Device identifier",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Virtual Device type",
		},
	}
}

func readAccessPointInterface() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Equinix-assigned interface identifier",
		},
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Access Point Interface id",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Interface type- CSP",
		},
	}
}

func readPortSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Equinix-assigned port identifier",
		},
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Unique Resource Identifier",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port name",
		},
		"redundancy": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Redundancy Information",
			Elem: &schema.Resource{
				Schema: readPortRedundancySch(),
			},
		},
	}
}

func readConnectionSideAccessPointSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Access point type - VD, COLO",
		},
		"authentication_key": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Authentication key for provider based connections",
		},
		"account": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Account",
			Elem: &schema.Resource{
				Schema: readAccountSch(),
			},
		},
		"location": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Access point location",
			Elem: &schema.Resource{
				Schema: readLocationSch(),
			},
		},
		"port": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Port access point information",
			Elem: &schema.Resource{
				Schema: readPortSch(),
			},
		},
		"profile": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Service Profile",
			Elem: &schema.Resource{
				Schema: readServiceProfileSch(),
			},
		},
		"gateway": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Gateway access point information",
			Elem: &schema.Resource{
				Schema: readVirtualGatewaySch(),
			},
		},
		"link_protocol": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Connection link protocol",
			Elem: &schema.Resource{
				Schema: readAccessPointLinkProtocolSch(),
			},
		},
		"virtual_device": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Virtual device",
			Elem:        &schema.Resource{Schema: readAccessPointVirtualDeviceSch()},
		},
		"interface": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Virtual device interface",
			Elem: &schema.Resource{
				Schema: readAccessPointInterface(),
			},
		},
		"seller_region": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Access point seller region",
		},
		"peering_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Peering Type - for Azure - Private or Public",
		},
		"provider_connection_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Provider assigned Connection Id",
		},
	}
}

func readFabricConnectionSideSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"service_token": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "For service token based connections, Service tokens authorize users to access protected resources and services. Resource owners can distribute the tokens to trusted partners and vendors, allowing selected third parties to work directly with Equinix network assets",
			Elem: &schema.Resource{
				Schema: readServiceTokenSch(),
			},
		},
		"access_point": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Point of access details",
			Elem: &schema.Resource{
				Schema: readConnectionSideAccessPointSch(),
			},
		},
	}
}

func readRedundancySch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"group": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Redundancy group identifier",
		},
		"priority": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Priority type - Primary or Secondary",
		},
	}
}

func readPortRedundancySch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"priority": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Priority type-Primary or Secondary",
		},
	}
}

func readChangeLogSch() map[string]*schema.Schema {
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
			Description: "Created on Date and Time",
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
			Description: "Updated on Date and Time",
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
			Description: "Deleted on Date and Time",
		},
	}
}

func readOrderSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"purchase_order_number": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Purchase order number",
		},
		"billing_tier": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Billing tier for connection bandwidth",
		},
		"order_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Order Identification",
		},
		"order_number": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Order Reference Number",
		},
	}
}

func readAccountSch() map[string]*schema.Schema {
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
		"global_cust_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Global Customer organization identifier",
		},
	}
}

func readNotificationSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Notification Type- ALL,CONNECTION_APPROVAL,SALES_REP_NOTIFICATIONS, NOTIFICATIONS",
		},
		"send_interval": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Send interval",
		},
		"emails": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Array of contact emails",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func readErrorAdditionalInfoSch() map[string]*schema.Schema {
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

func readOperationalErrorSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"error_code": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Error code",
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
				Schema: readErrorAdditionalInfoSch(),
			},
		},
	}
}

func readOperationSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"provider_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Connection provider readiness status - AVAILABLE, DEPROVISIONED, DEPROVISIONING ...",
		},
		"equinix_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Connection status - REJECTED_ACK, REJECTED, PENDING_DELETE, PROVISIONED ...",
		},
		"errors": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Errors occurred",
			Elem: &schema.Resource{
				Schema: readOperationalErrorSch(),
			},
		},
	}
}

func readAdditionalInfoSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Additional information key",
		},
		"value": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Additional information value",
		},
	}
}

func readFabricConnectionResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Equinix-assigned connection identifier",
		},
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Connection URI information",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Connection name. An alpha-numeric 24 characters string which can include only hyphens and underscores",
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Customer-provided connection description",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Defines the connection type like VG_VC, EVPL_VC, EPL_VC, EC_VC, GW_VC, ACCESS_EPL_VC",
		},
		"bandwidth": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Connection bandwidth in Mbps",
		},
		"is_remote": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Connection property derived from access point locations",
		},
		"state": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Connection overall state",
		},
		"operation": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Connection specific operational data",
			Elem: &schema.Resource{
				Schema: readOperationSch(),
			},
		},
		"notifications": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Preferences for notifications on connection configuration or status changes",
			Elem: &schema.Resource{
				Schema: readNotificationSch(),
			},
		},
		"additional_info": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Connection additional information",
			Elem: &schema.Resource{
				Schema: readAdditionalInfoSch(),
			},
		},
		"order": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Order related to this connection information",
			Elem: &schema.Resource{
				Schema: readOrderSch(),
			},
		},
		"account": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Customer account information that is associated with this connection",
			Elem: &schema.Resource{
				Schema: readAccountSch(),
			},
		},
		"change_log": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Captures connection lifecycle change information",
			Elem: &schema.Resource{
				Schema: readChangeLogSch(),
			},
		},
		"redundancy": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Redundancy Information",
			Elem: &schema.Resource{
				Schema: readRedundancySch(),
			},
		},
		"direction": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Connection directionality from the requester point of view",
		},
		"a_side": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Requester or Customer side connection configuration object of the multi-segment connection",
			Elem: &schema.Resource{
				Schema: readFabricConnectionSideSch(),
			},
		},
		"z_side": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Destination or Provider side connection configuration object of the multi-segment connection",
			Elem: &schema.Resource{
				Schema: readFabricConnectionSideSch(),
			},
		},
		"project": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Project information",
			Elem: &schema.Resource{
				Schema: createGatewayProjectSch(),
			},
		},
	}
}
