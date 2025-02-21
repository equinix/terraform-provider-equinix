package servicetoken

import (
	equinix_fabric_schema "github.com/equinix/terraform-provider-equinix/internal/fabric/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"VC_TOKEN", "EPL_TOKEN"}, false),
			Description:  "Service Token Type; VC_TOKEN,EPL_TOKEN",
		},
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Equinix-assigned service token identifier",
		},
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "An absolute URL that is the subject of the link's context.",
		},
		"issuer_side": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Information about token side; ASIDE, ZSIDE",
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Name of the Service Token",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Optional Description to the Service Token you will be creating",
		},
		"expiration_date_time": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Expiration date and time of the service token; 2020-11-06T07:00:00Z",
		},
		"service_token_connection": {
			Type:        schema.TypeSet,
			Required:    true,
			Description: "Service Token Connection Type Information",
			Elem:        serviceTokenConnectionSch(),
			Set:         schema.HashResource(serviceTokenConnectionSch()),
		},
		"state": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Service token state; ACTIVE, INACTIVE, EXPIRED, DELETED",
		},
		"notifications": {
			Type:        schema.TypeSet,
			Required:    true,
			Description: "Preferences for notifications on Service Token configuration or status changes",
			Elem: &schema.Resource{
				Schema: equinix_fabric_schema.NotificationSch(),
			},
		},
		"account": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Customer account information that is associated with this service token",
			Elem: &schema.Resource{
				Schema: equinix_fabric_schema.AccountSch(),
			},
		},
		"change_log": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Captures connection lifecycle change information",
			Elem: &schema.Resource{
				Schema: equinix_fabric_schema.ChangeLogSch(),
			},
		},
		"project": {
			Type:        schema.TypeSet,
			Optional:    true,
			Computed:    true,
			Description: "Project information",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: equinix_fabric_schema.ProjectSch(),
			},
		},
	}
}

func serviceTokenConnectionSch() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of Connection supported by Service Token you will create; EVPL_VC, EVPLAN_VC, EPLAN_VC, IPWAN_VC",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Equinix-assigned connection identifier",
			},
			"allow_remote_connection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Authorization to connect remotely",
			},
			"allow_custom_bandwidth": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Allow custom bandwidth value",
			},
			"bandwidth_limit": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 100000),
				Description:  "Connection bandwidth limit in Mbps",
			},
			"supported_bandwidths": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "List of permitted bandwidths'; For Port-based Service Tokens, the maximum allowable bandwidth is 50 Gbps, while for Virtual Device-based Service Tokens, it is limited to 10 Gbps",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"a_side": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "A-Side Connection link protocol,virtual device or network configuration",
				Elem:        serviceTokenAccessPointSch(),
				Set:         schema.HashResource(serviceTokenAccessPointSch()),
			},
			"z_side": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "Z-Side Connection link protocol,virtual device or network configuration",
				Elem:        serviceTokenAccessPointSch(),
				Set:         schema.HashResource(serviceTokenAccessPointSch()),
			},
		},
	}
}

func serviceTokenAccessPointSch() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"access_point_selectors": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of criteria for selecting network access points with optimal efficiency, security, compatibility, and availability",
				Elem:        accessPointSelectorsSch(),
			},
		},
	}
}

func accessPointSelectorsSch() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Type of Access point; COLO, VD, NETWORK",
			},
			"port": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Port Configuration",
				MaxItems:    1,
				Elem:        portSch(),
			},
			"link_protocol": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Link protocol Configuration",
				MaxItems:    1,
				Elem:        linkProtocolSch(),
			},
			"virtual_device": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Virtual Device Configuration",
				MaxItems:    1,
				Elem:        virtualDeviceSch(),
			},
			"interface": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Virtual Device Interface Configuration",
				MaxItems:    1,
				Elem:        interfaceSch(),
			},
			"network": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "Network Configuration",
				MaxItems:    1,
				Elem:        networkSch(),
			},
		},
	}
}

func portSch() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique Resource Identifier",
			},
			"uuid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Equinix-assigned Port identifier",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Type of Port",
			},
			"cvp_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Customer virtual port Id",
			},
			"bandwidth": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Port Bandwidth",
			},
			"port_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Port Name",
			},
			"encapsulation_protocol_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Port Encapsulation",
			},
			"account_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Account Name",
			},
			"priority": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Port Priority",
			},
			"location": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Port Location",
				Elem: &schema.Resource{
					Schema: equinix_fabric_schema.LocationSch(),
				},
			},
		},
	}
}

func virtualDeviceSch() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique Resource Identifier",
			},
			"uuid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Equinix-assigned Virtual Device identifier",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Virtual Device type",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Customer-assigned Virtual Device Name",
			},
			"cluster": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Virtual Device Cluster Information",
			},
		},
	}
}

func linkProtocolSch() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Type of the link protocol - UNTAGGED, DOT1Q, QINQ, EVPN_VXLAN",
				ValidateFunc: validation.StringInSlice([]string{"UNTAGGED", "DOT1Q", "QINQ", "EVPN_VXLAN"}, true),
			},
			"vlan_tag": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Vlan Tag information, vlanTag value specified for DOT1Q connections",
			},
			"vlan_s_tag": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Vlan Provider Tag information, vlanSTag value specified for QINQ connections",
			},
			"vlan_c_tag": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Vlan Customer Tag information, vlanCTag value specified for QINQ connections",
			},
		},
	}
}

func interfaceSch() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Equinix-assigned interface identifier",
			},
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "id",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Interface type",
			},
		},
	}
}

func networkSch() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Equinix-assigned Network identifier",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique Resource Identifier",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of Network",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Network Name",
			},
			"scope": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Scope of Network",
			},
			"location": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Location",
				Elem: &schema.Resource{
					Schema: equinix_fabric_schema.LocationSch(),
				},
			},
		},
	}
}
