package equinix

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var createServiceTokenRes = &schema.Resource{
	Schema: createServiceTokenSch(),
}

func createServiceTokenSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"VC_TOKEN"}, true),
			Description:  "Token type",
		},
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "An absolute URL that is the subject of the link's context",
		},
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Equinix-assigned service token identifier",
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Service token description",
		},
	}
}

var createLocationRes = &schema.Resource{
	Schema: createLocationSch(),
}

func createLocationSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"region": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Access point region",
		},
		"metro_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Access point metro name",
		},
		"metro_code": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Access point metro code",
		},
		"ibx": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "IBX",
		},
	}
}

//TODO missing uuid in swager generated spec.
func createVirtualGatewaySch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Gateway Type",
		},
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Unique Resource Identifier",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Virtual Device type",
		},
		"project": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Project this gateway created in",
			Elem:        &schema.Resource{Schema: createGatewayProjectSch()},
		},
	}
}

var createGatewayProjectSchRes = &schema.Resource{
	Schema: createGatewayProjectSch(),
}

func createGatewayProjectSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"project_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Project Id",
		},
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Unique Resource URL",
		},
	}
}

func createServiceProfileSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Service Profile URI response attribute",
		},
		"type": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"L2_PROFILE", "L3_PROFILE", "ECIA_PROFILE", "ECMC_PROFILE"}, true),
			Description:  "Service profile type",
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Customer-assigned service profile name",
		},
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Equinix assigned service profile identifier",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "User-provided service description",
		},
	}
}

var createAccessPointLinkProtocolSchRes = &schema.Resource{
	Schema: createAccessPointLinkProtocolSch(),
}

func createAccessPointLinkProtocolSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "Type of the link protocol",
			ValidateFunc: validation.StringInSlice([]string{"UNTAGGED", "DOT1Q", "QINQ", "EVPN_VXLAN"}, true),
		},
		"vlan_tag": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Vlan Tag information, vlanTag value specified for DOT1Q connections",
		},
		"vlan_s_tag": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Vlan Provider Tag information, vlanSTag value specified for QINQ connections",
		},
		"vlan_c_tag": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Vlan Customer Tag information, vlanCTag value specified for QINQ connections",
		},
		"unit": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Unit",
		},
		"vni": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "vni",
		},
		"int_unit": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "int unit",
		},
	}
}

var createAccessPointVirtualDeviceRes = &schema.Resource{
	Schema: createAccessPointVirtualDeviceSch(),
}

func createAccessPointVirtualDeviceSch() map[string]*schema.Schema {

	return map[string]*schema.Schema{
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Unique Resource Identifier",
		},
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Equinix-assigned Virtual Device identifier",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Virtual Device type",
		},
	}
}

func createAccessPointInterface() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Equinix-assigned interface identifier",
		},
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "id",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Interface type",
		},
	}
}

func createFabricConnectionRoutingProtocol() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Routing Protocol type",
		},
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Routing protocol instance identifier",
		},
		"state": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Routing protocol instance state",
		},
	}
}

var createPortRes = &schema.Resource{
	Schema: createPortSch(),
}

func createPortSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Port information",
		},
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Unique Resource Identifier",
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Port name",
		},
		"redundancy": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Redundancy Information",
			Elem: &schema.Resource{
				Schema: createPortRedundancySch(),
			},
		},
	}
}

var createConnectionSideAccessPointRes = &schema.Resource{
	Schema: createConnectionSideAccessPointSch(),
}

func createConnectionSideAccessPointSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"COLO", "VD", "VG", "SP", "IGW", "SUBNET", "GW"}, true),
			Description:  "Access point type",
		},
		"authentication_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Access point type",
		},
		"account": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Account",
			Elem: &schema.Resource{
				Schema: createAccountSch(),
			},
		},
		"location": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Access point location",
			Elem: &schema.Resource{
				Schema: createLocationSch(),
			},
		},
		"port": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Port access point information",
			Elem: &schema.Resource{
				Schema: createPortSch(),
			},
		},
		"profile": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Service Profile",
			Elem: &schema.Resource{
				Schema: createServiceProfileSch(),
			},
		},
		"gateway": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Gateway access point information",
			Elem: &schema.Resource{
				Schema: createVirtualGatewaySch(),
			},
		},
		"link_protocol": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Connection link protocol",
			Elem: &schema.Resource{
				Schema: createAccessPointLinkProtocolSch(),
			},
		},
		"virtual_device": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Virtual device",
			Elem:        &schema.Resource{Schema: createAccessPointVirtualDeviceSch()},
		},
		"interface": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Virtual device interface",
			Elem: &schema.Resource{
				Schema: createAccessPointInterface(),
			},
		},
		"seller_region": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Access point seller region",
		},
		"peering_type": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"PRIVATE", "MICROSOFT", "PUBLIC", "MANUAL"}, true),
			Description:  "Peering Type",
		},
		"routing_protocols": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Access point routing protocols configuration",
			Elem: &schema.Resource{
				Schema: createFabricConnectionRoutingProtocol(),
			},
		},
		"additional_info": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Access point additional Information",
		},
		"provider_connection_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Provider assigned Connection Id",
		},
	}
}

var createFabricConnectionSideRes = &schema.Resource{
	Schema: createFabricConnectionSideSch(),
}

func createFabricConnectionSideSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"service_token": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "For service token based connections, Service tokens authorize users to access protected resources and services. Resource owners can distribute the tokens to trusted partners and vendors, allowing selected third parties to work directly with Equinix network assets",
			Elem: &schema.Resource{
				Schema: createServiceTokenSch(),
			},
		},
		"access_point": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Point of access details",
			Elem: &schema.Resource{
				Schema: createConnectionSideAccessPointSch(),
			},
		},
		"additional_info": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Connection Side additional details",
			Elem: &schema.Resource{
				Schema: createAdditionalInfoSch(),
			},
		},
	}
}

var createRedundancyRes = &schema.Resource{
	Schema: createRedundancySch(),
}

func createRedundancySch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"group": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Redundancy group identifier",
		},
		"priority": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"PRIMARY", "SECONDARY"}, true),
			Description:  "Priority type",
		},
	}
}

var createPortRedundancyRes = &schema.Resource{
	Schema: createPortRedundancySch(),
}

func createPortRedundancySch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"priority": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"PRIMARY", "SECONDARY"}, true),
			Description:  "Priority type",
		},
	}
}

var createChangeLogRes = &schema.Resource{
	Schema: createChangeLogSch(),
}

func createChangeLogSch() map[string]*schema.Schema {
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

var createOrderRes = &schema.Resource{
	Schema: createOrderSch(),
}

func createOrderSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"purchase_order_number": {
			Type:        schema.TypeString,
			Optional:    true,
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

var createAccountRes = &schema.Resource{
	Schema: createAccountSch(),
}

func createAccountSch() map[string]*schema.Schema {
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
	}
}

func createNotificationSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Notification Type",
		},
		"send_interval": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Send interval",
		},
		"emails": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Array of contact emails",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func createErrorAdditionalInfoSch() map[string]*schema.Schema {
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

func createOperationalErrorSch() map[string]*schema.Schema {
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
				Schema: createErrorAdditionalInfoSch(),
			},
		},
	}
}

var createOperationRes = &schema.Resource{
	Schema: createOperationSch(),
}

func createOperationSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"provider_status": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "Connection provider readiness status",
			ValidateFunc: validation.StringInSlice([]string{"AVAILABLE", "DEPROVISIONED", "DEPROVISIONING", "FAILED", "NOT_AVAILABLE", "PENDING_APPROVAL", "PROVISIONED", "PROVISIONING", "REJECTED", "PENDING_BGP", "OUT_OF_BANDWIDTH", "DELETED", "ERROR", "ERRORED", "NOTPROVISIONED", "NOT_PROVISIONED", "ORDERING", "DELETING", "PENDING DELETE", "N/A"}, true),
		},
		"equinix_status": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "Connection status",
			ValidateFunc: validation.StringInSlice([]string{"REJECTED_ACK", "REJECTED", "PENDING_DELETE", "PROVISIONED", "BEING_REPROVISIONED", "BEING_DEPROVISIONED", "BEING_PROVISIONED", "CREATED", "ERRORED", "PENDING_DEPROVISIONING", "APPROVED", "ORDERING", "PENDING_APPROVAL", "NOT_PROVISIONED", "DEPROVISIONING", "NOT_DEPROVISIONED", "PENDING_AUTO_APPROVAL", "PROVISIONING", "PENDING_BGP_PEERING", "PENDING_PROVIDER_VLAN", "DEPROVISIONED", "DELETED", "PENDING_BANDWIDTH_APPROVAL", "AUTO_APPROVAL_FAILED", "UPDATE_PENDING", "DELETED_API", "MODIFIED", "PENDING_PROVIDER_VLAN_ERROR", "DRAFT", "CANCELLED", "PENDING_INTERFACE_CONFIGURATION"}, true),
		},
		"operational_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Connection operational status",
		},
		"errors": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Errors occurred",
			Elem: &schema.Resource{
				Schema: createOperationalErrorSch(),
			},
		},
		"op_status_changed_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "When connection transitioned into current operational status",
		},
	}
}

func createAdditionalInfoSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Additional information key",
		},
		"value": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Additional information value",
		},
	}
}

func createFabricConnectionResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "TBD",
		},
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Connection URI information",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Connection name. An alpha-numeric 24 characters string which can include only hyphens and underscores",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Customer-provided connection description",
		},
		"type": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"VG_VC", "EVPL_VC", "EPL_VC", "EC_VC", "GW_VC", "ACCESS_EPL_VC"}, true),
			Description:  "Defines the connection type like VG_VC, EVPL_VC, EPL_VC, EC_VC, GW_VC, ACCESS_EPL_VC",
		},
		"bandwidth": {
			Type:        schema.TypeInt,
			Required:    true,
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
			Description: "Connection type-specific operational data",
			Elem: &schema.Resource{
				Schema: createOperationSch(),
			},
		},
		"notifications": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Preferences for notifications on connection configuration or status changes",
			Elem: &schema.Resource{
				Schema: createNotificationSch(),
			},
		},
		"additional_info": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Connection additional information",
			Elem: &schema.Resource{
				Schema: createAdditionalInfoSch(),
			},
		},
		"order": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Order related to this connection information",
			Elem: &schema.Resource{
				Schema: createOrderSch(),
			},
		},
		"account": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Customer account information that is associated with this connection",
			Elem: &schema.Resource{
				Schema: createAccountSch(),
			},
		},
		"change_log": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Captures connection lifecycle change information",
			Elem: &schema.Resource{
				Schema: createChangeLogSch(),
			},
		},
		"redundancy": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Redundancy Information",
			Elem: &schema.Resource{
				Schema: createRedundancySch(),
			},
		},
		"direction": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Connection directionality from the requester point of view",
		},
		"a_side": {
			Type:        schema.TypeSet,
			Required:    true,
			Description: "Requester or Customer side connection configuration object of the multi-segment connection",
			Elem: &schema.Resource{
				Schema: createFabricConnectionSideSch(),
			},
		},
		"z_side": {
			Type:        schema.TypeSet,
			Required:    true,
			Description: "Destination or Provider side connection configuration object of the multi-segment connection",
			Elem: &schema.Resource{
				Schema: createFabricConnectionSideSch(),
			},
		},
	}
}
