package equinix

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var createInvitationRes = &schema.Resource{
	Schema: createInvitationSch(),
}

func createInvitationSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Invitation type",
		},
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
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

var createAccessPointSelectorSimplifiedMetadataEntityRes = &schema.Resource{
	Schema: createAccessPointSelectorSimplifiedMetadataEntitySch(),
}

func createAccessPointSelectorSimplifiedMetadataEntitySch() map[string]*schema.Schema {

	return map[string]*schema.Schema{
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Url to entity",
		},
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Equinix assigned Identifier",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Type of Port",
		},
	}
}

var createServiceTokenLinkProtocolRes = &schema.Resource{
	Schema: createServiceTokenLinkProtocolSch(),
}

func createServiceTokenLinkProtocolSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"UNTAGGED", "DOT1Q", "QINQ", "EVPN_VXLAN"}, true),
			Description:  "Type",
		},
	}
}

var createFabricConnectionServiceTokenAccessPointSelectorRes = &schema.Resource{
	Schema: createFabricConnectionServiceTokenAccessPointSelectorSch(),
}

func createFabricConnectionServiceTokenAccessPointSelectorSch() map[string]*schema.Schema {
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
				Schema: createAccessPointSelectorSimplifiedMetadataEntitySch(),
			},
		},
		"link_protocol": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Link Protocol",
			Elem: &schema.Resource{
				Schema: createServiceTokenLinkProtocolSch(),
			},
		},
	}
}

var createServiceTokenSideRes = &schema.Resource{
	Schema: createServiceTokenSideSch(),
}

func createServiceTokenSideSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"access_point_selectors": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Access Point Selectors",
			Elem: &schema.Resource{
				Schema: createFabricConnectionServiceTokenAccessPointSelectorSch(),
			},
		},
	}
}

var createServiceTokenConnectionSideRes = &schema.Resource{
	Schema: createServiceTokenConnectionSideSch(),
}

func createServiceTokenConnectionSideSch() map[string]*schema.Schema {

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
				Schema: createServiceTokenSideSch(),
			},
		},
		"z_side": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "ZSide",
			Elem: &schema.Resource{
				Schema: createServiceTokenSideSch(),
			},
		},
	}
}

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
				Schema: createServiceTokenConnectionSideSch(),
			},
		},
		"state": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "INACTIVE", "EXPIRED", "DELETED"}, true),
			Description:  "State",
		},
		"notifications": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Service token related notifications",
			Elem: &schema.Resource{
				Schema: createNotificationSch(),
			},
		},
		"account": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "account",
			Elem: &schema.Resource{
				Schema: createAccountSch(),
			},
		},
		"change_log": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Change Log",
			Elem: &schema.Resource{
				Schema: createChangeLogSch(),
			},
		},
	}
}

var createLocationRes = &schema.Resource{
	Schema: createLocationSch(),
}

func createLocationSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Unique Resource Identifier",
		},
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

var createLocationNoIbxRes = &schema.Resource{
	Schema: createLocationNoIbxSch(),
}

func createLocationNoIbxSch() map[string]*schema.Schema {
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

var createServiceProfileAccessPointTypeRes = &schema.Resource{
	Schema: createServiceProfileAccessPointType(),
}

func createServiceProfileAccessPointType() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"type": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"VD", "COLO"}, true),
			Description:  "Type",
		},
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Uuid",
		},
	}
}

var createCustomFieldsRes = &schema.Resource{
	Schema: createCustomFields(),
}

func createCustomFields() map[string]*schema.Schema {
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

var createServiceProfileMarketingProcessStepsRes = &schema.Resource{
	Schema: createServiceProfileMarketingProcessSteps(),
}

func createServiceProfileMarketingProcessSteps() map[string]*schema.Schema {
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

var createServiceProfileMarketingInfoRes = &schema.Resource{
	Schema: createServiceProfileMarketingInfo(),
}

func createServiceProfileMarketingInfo() map[string]*schema.Schema {

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
				Schema: createServiceProfileMarketingProcessSteps(),
			},
		},
	}
}

var createServiceProfileAccessPointColoRes = &schema.Resource{
	Schema: createServiceProfileAccessPointColo(),
}

func createServiceProfileAccessPointColo() map[string]*schema.Schema {
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
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "Type",
			ValidateFunc: validation.StringInSlice([]string{"VD", "COLO"}, true),
		},
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "uuid",
		},
		"location": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Location",
			Elem:        &schema.Resource{Schema: createLocationSch()},
		},
	}
}

var createServiceProfileAccessPointVdRes = &schema.Resource{
	Schema: createServiceProfileAccessPointVd(),
}

func createServiceProfileAccessPointVd() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"VD", "COLO"}, true),
			Description:  "type",
		},
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Uuid",
		},
		"location": {
			Type:        schema.TypeSet,
			Optional:    true,
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

var createServiceProfileMetroSchRes = &schema.Resource{
	Schema: createServiceProfileMetroSch(),
}

func createServiceProfileMetroSch() map[string]*schema.Schema {
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

var createGatewayPackageSchRes = &schema.Resource{
	Schema: createGatewayPackageSch(),
}

func createGatewayPackageSch() map[string]*schema.Schema {
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

var createVirtualGatewaySchRes = &schema.Resource{
	Schema: createVirtualGatewaySch(),
}

func createVirtualGatewaySch() map[string]*schema.Schema {
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
			Elem:        &schema.Resource{Schema: createLocationNoIbxSch()},
		},
		"package": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Package information",
			Elem:        &schema.Resource{Schema: createGatewayPackageSch()},
		},
		"order": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Order Information",
			Elem:        &schema.Resource{Schema: createOrderSch()},
		},
		"project": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Project this gateway created in",
			Elem:        &schema.Resource{Schema: createGatewayProjectSch()},
		},
		"account": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Account info",
			Elem:        &schema.Resource{Schema: createAccountSch()},
		},
		"notifications": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Notifications",
			Elem:        &schema.Resource{Schema: createNotificationSch()},
		},
	}
}

var createServiceProfileSchRes = &schema.Resource{
	Schema: createServiceProfileSch(),
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
		"notifications": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Recipients of notifications on service profile change",
			Elem: &schema.Resource{
				Schema: createNotificationSch(),
			},
		},
		"tags": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Tags",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"visibility": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"PRIVATE", "PUBLIC"}, true),
			Description:  "Visibility of the service profile",
		},
		"allowed_emails": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "User Emails that are allowed to access this service profile",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"access_point_type_configs": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Access Point Type Configs",
			Elem: &schema.Resource{
				Schema: createServiceProfileAccessPointType(),
			},
		},
		"custom_fields": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Custom Fields",
			Elem: &schema.Resource{
				Schema: createCustomFields(),
			},
		},
		"marketing_info": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Marketing Info",
			Elem: &schema.Resource{
				Schema: createServiceProfileMarketingInfo(),
			},
		},
		"ports": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Ports",
			Elem: &schema.Resource{
				Schema: createServiceProfileAccessPointColo(),
			},
		},
		"virtual_devices": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Virtual Devices",
			Elem: &schema.Resource{
				Schema: createServiceProfileAccessPointVd(),
			},
		},
		"metros": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Derived response attribute",
			Elem: &schema.Resource{
				Schema: createServiceProfileMetroSch(),
			},
		},
		"self_profile": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Response attribute indicates whether the profile belongs to the same organization as the api-invoker",
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
			Type:        schema.TypeString,
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

var createAccessPointVirtualDeviceSchRes = &schema.Resource{
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

var createAccessPointInterfaceRes = &schema.Resource{
	Schema: createAccessPointInterface(),
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

var createFabricConnectionRoutingProtocolRes = &schema.Resource{
	Schema: createFabricConnectionRoutingProtocol(),
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

var createFabricConnectionSideOrganizationRes = &schema.Resource{
	Schema: createFabricConnectionSideOrganizationSch(),
}

func createFabricConnectionSideOrganizationSch() map[string]*schema.Schema {
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

var createConnectionSideCompanyProfileRes = &schema.Resource{
	Schema: createConnectionSideCompanyProfileSch(),
}

func createConnectionSideCompanyProfileSch() map[string]*schema.Schema {
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
			Elem:        &schema.Resource{Schema: createFabricConnectionSideOrganizationSch()},
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
			ValidateFunc: validation.StringInSlice([]string{"COLO", "VD", "VG", "SP", "IGW", "IGW", "SUBNET", "GW"}, true),
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
		"invitation": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Invitation based on connection request",
			Elem: &schema.Resource{
				Schema: createInvitationSch(),
			},
		},
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
		"company_profile": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Company Profile",
			Elem: &schema.Resource{
				Schema: createConnectionSideCompanyProfileSch(),
			},
		},
		"nat": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Network Address Translation type",
			Elem: &schema.Resource{
				Schema: createNatSch(),
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

var createNatRes = &schema.Resource{
	Schema: createNatSch(),
}

func createNatSch() map[string]*schema.Schema {

	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Type",
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
		"group": {
			Type:        schema.TypeInt,
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
		"ucm_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "System unique identifier",
		},
		"global_cust_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Global Customer organization identifier",
		},
	}
}

var createNotificationRes = &schema.Resource{
	Schema: createNotificationSch(),
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

var createErrorAdditionalInfoRes = &schema.Resource{
	Schema: createErrorAdditionalInfoSch(),
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

var createOperationalErrorRes = &schema.Resource{
	Schema: createOperationalErrorSch(),
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

var createChangeRes = &schema.Resource{
	Schema: createChangeSch(),
}

func createChangeSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Unique identifier of the change",
		},
		"type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Type of change",
		},
		"status": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Current outcome of the change flow",
		},
		"creation_date_time": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Time change request received",
		},
		"updated_date_time": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Record last updated",
		},
		"information": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Additional information",
		},
		"data": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Change operation data",
			Elem: &schema.Resource{
				Schema: createFabricConnectionChangeDataSch(),
			},
		},
	}
}

var createFabricConnectionChangeDataRes = &schema.Resource{
	Schema: createFabricConnectionChangeDataSch(),
}

func createFabricConnectionChangeDataSch() map[string]*schema.Schema {

	return map[string]*schema.Schema{
		"op": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Operation name",
		},
		"path": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Path inside document leading to updated parameter",
		},
		"value": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "New value for updated parameter",
		},
	}
}

var createAdditionalInfoRes = &schema.Resource{
	Schema: createAdditionalInfoSch(),
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

var createIpv4Res = &schema.Resource{
	Schema: createIpv4Sch(),
}

func createIpv4Sch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"customer_peer_ip": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Customer peering ip",
		},
		"provider_peer_ip": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Provider peering ip",
		},
	}
}

var createRoutingProtocolRes = &schema.Resource{
	Schema: createRoutingProtocolSch(),
}

func createRoutingProtocolSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Routing Protocol Type",
		},
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Routing protocol identifier",
		},
		"customer_asn": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Customer asn",
		},
		"peer_asn": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Peer asn",
		},
		"bgp_auth_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "BGP authorization key",
		},
		"ipv4": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "ip information",
			Elem:        &schema.Resource{Schema: createIpv4Sch()},
		},
		"route_filters": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Route filters values",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
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
		"platform_uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Unique identifier of the connection, internal",
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
			ValidateFunc: validation.StringInSlice([]string{"VG_VC", "EVPL_VC", "EPL_VC", "EC_VC", "GW_VC", "ACCESS_EPL_VC", "NONGENERIC"}, true),
			Description:  "Defines the connection type like VG_VC, EVPL_VC, EPL_VC, EC_VC, GW_VC, ACCESS_EPL_VC, NONGENERIC",
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
		"change": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Represents latest change request and its state information",
			Elem: &schema.Resource{
				Schema: createChangeSch(),
			},
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
		"tags": {
			Type:        schema.TypeList,
			Optional:    true,
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
				Schema: createRoutingProtocolSch(),
			},
		},
	}
}
