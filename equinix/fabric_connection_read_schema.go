package equinix

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func readInvitationSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Invitation type",
		},
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Equinix-assigned invitation identifier",
		},
		"state": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Invitation status as it is today",
		},
		"message": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Invitation message",
		},
		"email": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Invitation recipient",
		},
		"expiry": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Invitation expiry time",
		},
	}
}

func readAccessPointSelectorSimplifiedMetadataEntitySch() map[string]*schema.Schema {

	return map[string]*schema.Schema{
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Url to entity",
		},
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Equinix assigned Identifier",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Type of Port",
		},
	}
}

func readServiceTokenLinkProtocolSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Type",
		},
	}
}

func readFabricConnectionServiceTokenAccessPointSelectorSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Selector Type",
		},
		"port": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "port",
			Elem: &schema.Resource{
				Schema: readAccessPointSelectorSimplifiedMetadataEntitySch(),
			},
		},
		"link_protocol": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Link Protocol",
			Elem: &schema.Resource{
				Schema: readServiceTokenLinkProtocolSch(),
			},
		},
	}
}

func readServiceTokenSideSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"access_point_selectors": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Access Point Selectors",
			Elem: &schema.Resource{
				Schema: readFabricConnectionServiceTokenAccessPointSelectorSch(),
			},
		},
	}
}

func readServiceTokenConnectionSideSch() map[string]*schema.Schema {

	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Type of Connection",
		},
		"allow_remote_connection": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Authorization to connect remotely",
		},
		"bandwidth_limit": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Connection bandwidth limit in Mbps",
		},
		"supported_bandwidths": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "List of permitted bandwidths",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"a_side": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "ASide",
			Elem: &schema.Resource{
				Schema: readServiceTokenSideSch(),
			},
		},
		"z_side": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "ZSide",
			Elem: &schema.Resource{
				Schema: readServiceTokenSideSch(),
			},
		},
	}
}

func readServiceTokenSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Token type",
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
		"expiry": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Lifespan (in days) of the service token",
		},
		"expiration_date_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Expiration date and time of the service token",
		},
		"connection": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Connection",
			Elem: &schema.Resource{
				Schema: readServiceTokenConnectionSideSch(),
			},
		},
		"state": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "State",
		},
		"notifications": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Service token related notifications",
			Elem: &schema.Resource{
				Schema: readNotificationSch(),
			},
		},
		"account": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "account",
			Elem: &schema.Resource{
				Schema: readAccountSch(),
			},
		},
		"change_log": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Change Log",
			Elem: &schema.Resource{
				Schema: readChangeLogSch(),
			},
		},
	}
}

func readLocationSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Unique Resource Identifier",
		},
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
			Description: "IBX",
		},
	}
}

func readLocationNoIbxSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Unique Resource Identifier",
		},
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
	}
}

func readServiceProfileAccessPointType() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Type",
		},
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Uuid",
		},
	}
}

func readCustomFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
		"required": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Required",
		},
		"data_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Data Type",
		},
		"options": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Options",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"capture_in_email": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Capture this field as a part of email notification",
		},
	}
}

func readServiceProfileMarketingProcessSteps() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"title": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Service profile custom step title",
		},
		"sub_title": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Service profile custom step sub title",
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Service profile custom step description",
		},
	}
}

func readServiceProfileMarketingInfo() map[string]*schema.Schema {

	return map[string]*schema.Schema{
		"logo": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Logo file name",
		},
		"promotion": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Profile promotion on marketplace",
		},
		"process_steps": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Process Steps",
			Elem: &schema.Resource{
				Schema: readServiceProfileMarketingProcessSteps(),
			},
		},
	}
}

func readServiceProfileAccessPointColo() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"seller_region": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Seller Region",
		},
		"seller_region_description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Seller Region Description",
		},
		"cross_connect_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Cross Connect Id",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Type",
		},
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "uuid",
		},
		"location": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Location",
			Elem:        &schema.Resource{Schema: createLocationSch()},
		},
	}
}

func readServiceProfileAccessPointVd() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "type",
		},
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Uuid",
		},
		"location": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Location",
			Elem:        &schema.Resource{Schema: createLocationSch()},
		},
		"interface_uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "InterfaceUuid",
		},
	}
}

func readServiceProfileMetroSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"code": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Metro code",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Metro name",
		},
		"ibxs": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Allowed ibxes in the metro",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"in_trail": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "InTrail",
		},
		"display_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Service metro display name",
		},
		"seller_regions": {
			Type:        schema.TypeMap,
			Computed:    true,
			Description: "Seller Regions",
		},
	}
}

func readGatewayProjectSch() map[string]*schema.Schema {
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

func readGatewayPackageSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Unique Resource Identifier",
		},
		"code": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Gateway package code",
		},
	}
}

func readVirtualGatewaySch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Gateway Type",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Gateway Name",
		},
		"location": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Gateway Location",
			Elem:        &schema.Resource{Schema: readLocationNoIbxSch()},
		},
		"package": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Package information",
			Elem:        &schema.Resource{Schema: readGatewayPackageSch()},
		},
		"order": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Order Information",
			Elem:        &schema.Resource{Schema: readOrderSch()},
		},
		"project": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Project this gateway created in",
			Elem:        &schema.Resource{Schema: readGatewayProjectSch()},
		},
		"account": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Account info",
			Elem:        &schema.Resource{Schema: readAccountSch()},
		},
		"notifications": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Notifications",
			Elem:        &schema.Resource{Schema: readNotificationSch()},
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
			Description: "Service profile type",
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
			Description: "User-provided service description",
		},
		"notifications": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Recipients of notifications on service profile change",
			Elem: &schema.Resource{
				Schema: readNotificationSch(),
			},
		},
		"tags": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Tags",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"visibility": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Visibility of the service profile",
		},
		"allowed_emails": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "User Emails that are allowed to access this service profile",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"access_point_type_configs": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Access Point Type Configs",
			Elem: &schema.Resource{
				Schema: readServiceProfileAccessPointType(),
			},
		},
		"custom_fields": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Custom Fields",
			Elem: &schema.Resource{
				Schema: readCustomFields(),
			},
		},
		"marketing_info": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Marketing Info",
			Elem: &schema.Resource{
				Schema: readServiceProfileMarketingInfo(),
			},
		},
		"ports": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Ports",
			Elem: &schema.Resource{
				Schema: readServiceProfileAccessPointColo(),
			},
		},
		"virtual_devices": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Virtual Devices",
			Elem: &schema.Resource{
				Schema: readServiceProfileAccessPointVd(),
			},
		},
		"metros": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Derived response attribute",
			Elem: &schema.Resource{
				Schema: readServiceProfileMetroSch(),
			},
		},
		"self_profile": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Response attribute indicates whether the profile belongs to the same organization as the api-invoker",
		},
	}
}

func readAccessPointLinkProtocolSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Type of the link protocol",
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
		"unit": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Unit",
		},
		"vni": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "vni",
		},
		"int_unit": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "int unit",
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
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Customer-assigned Virtual Device name",
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
			Description: "id",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Interface type",
		},
	}
}

func readFabricConnectionRoutingProtocol() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Routing Protocol type",
		},
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Routing protocol instance identifier",
		},
		"state": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Routing protocol instance state",
		},
	}
}

func readFabricConnectionSideOrganizationSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Organization id -customer organization id",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Organization id -customer organization id",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Organization name -customer organization name",
		},
	}
}

func readConnectionSideCompanyProfileSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Common company profile id",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Common company profile name",
		},
		"organization": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Organization details",
			Elem:        &schema.Resource{Schema: readFabricConnectionSideOrganizationSch()},
		},
	}
}

func readPortSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port information",
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

func readPortRedundancySch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"group": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Redundancy group identifier",
		},
		"priority": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Priority type",
		},
	}
}

func readConnectionSideAccessPointSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Access point type",
		},
		"authentication_key": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Access point type",
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
			Description: "Peering Type",
		},
		"routing_protocols": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Access point routing protocols configuration",
			Elem: &schema.Resource{
				Schema: readFabricConnectionRoutingProtocol(),
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

func readFabricConnectionSideSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"invitation": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Invitation based on connection request",
			Elem: &schema.Resource{
				Schema: readInvitationSch(),
			},
		},
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
		"company_profile": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Company Profile",
			Elem: &schema.Resource{
				Schema: readConnectionSideCompanyProfileSch(),
			},
		},
		"nat": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Network Address Translation type",
			Elem: &schema.Resource{
				Schema: readNatSch(),
			},
		},
		"additional_info": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Connection Side additional details",
			Elem: &schema.Resource{
				Schema: readAdditionalInfoSch(),
			},
		},
	}
}

func readNatSch() map[string]*schema.Schema {

	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Type",
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
			Description: "Priority type",
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
		"ucm_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "System unique identifier",
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
			Description: "Notification Type",
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
			Description: "Connection provider readiness status",
		},
		"equinix_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Connection status",
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
				Schema: readOperationalErrorSch(),
			},
		},
		"op_status_changed_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "When connection transitioned into current operational status",
		},
	}
}

func readChangeSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Unique identifier of the change",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Type of change",
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Current outcome of the change flow",
		},
		"creation_date_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Time change request received",
		},
		"updated_date_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Record last updated",
		},
		"information": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Additional information",
		},
		"data": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Change operation data",
			Elem: &schema.Resource{
				Schema: readFabricConnectionChangeDataSch(),
			},
		},
	}
}

func readFabricConnectionChangeDataSch() map[string]*schema.Schema {

	return map[string]*schema.Schema{
		"op": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Operation name",
		},
		"path": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Path inside document leading to updated parameter",
		},
		"value": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "New value for updated parameter",
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

func readIpv4Sch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"customer_peer_ip": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Customer peering ip",
		},
		"provider_peer_ip": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Provider peering ip",
		},
	}
}

func readRoutingProtocolSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Routing Protocol Type",
		},
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Routing protocol identifier",
		},
		"customer_asn": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Customer asn",
		},
		"peer_asn": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Peer asn",
		},
		"bgp_auth_key": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "BGP authorization key",
		},
		"ipv4": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "ip information",
			Elem:        &schema.Resource{Schema: readIpv4Sch()},
		},
		"route_filters": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Route filters values",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func readFabricConnectionResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"uuid": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "TBD",
		},
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Connection URI information",
		},
		"platform_uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Unique identifier of the connection, internal",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			ForceNew:    true,
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
			Description: "Defines the connection type like VG_VC, EVPL_VC, EPL_VC, EC_VC, GW_VC, ACCESS_EPL_VC, NONGENERIC",
		},
		"bandwidth": {
			Type:        schema.TypeInt,
			Computed:    true,
			ForceNew:    true,
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
			ForceNew:    true,
			Description: "Connection overall state",
		},
		"change": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Represents latest change request and its state information",
			Elem: &schema.Resource{
				Schema: readChangeSch(),
			},
		},
		"operation": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Connection type-specific operational data",
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
			ForceNew:    true,
			Description: "Order related to this connection information",
			Elem: &schema.Resource{
				Schema: readOrderSch(),
			},
		},
		"account": {
			Type:        schema.TypeSet,
			Computed:    true,
			ForceNew:    true,
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
			ForceNew:    true,
			Description: "Requester or Customer side connection configuration object of the multi-segment connection",
			Elem: &schema.Resource{
				Schema: readFabricConnectionSideSch(),
			},
		},
		"z_side": {
			Type:        schema.TypeSet,
			Computed:    true,
			ForceNew:    true,
			Description: "Destination or Provider side connection configuration object of the multi-segment connection",
			Elem: &schema.Resource{
				Schema: readFabricConnectionSideSch(),
			},
		},
		"tags": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "User provided tags for the connection",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"routing_protocols": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Configured Routing protocol for the connection",
			Elem: &schema.Resource{
				Schema: readRoutingProtocolSch(),
			},
		},
	}
}
