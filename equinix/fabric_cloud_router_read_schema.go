package equinix

import (
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func readPackageSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"code": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Fabric Cloud Router package code",
		},
	}
}

func readCloudRouterResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Equinix-assigned Fabric Cloud Router identifier",
		},
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Fabric Cloud Router URI information",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Fabric Cloud Router name. An alpha-numeric 24 characters string which can include only hyphens and underscores",
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Customer-provided Fabric Cloud Router description",
		},
		"state": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Fabric Cloud Router overall state",
		},
		"equinix_asn": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Equinix ASN",
		},
		"bgp_ipv4_routes_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Access point used and maximum number of IPv4 BGP routes",
		},
		"bgp_ipv6_routes_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Access point used and maximum number of IPv6 BGP routes",
		},
		"connections_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Number of connections associated with this Access point",
		},
		"package": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Fabric Cloud Router package information",
			Elem: &schema.Resource{
				Schema: readPackageSch(),
			},
		},
		"change_log": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Captures Fabric Cloud Router lifecycle change information",
			Elem: &schema.Resource{
				Schema: equinix_schema.ChangeLogSch(),
			},
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Defines the Fabric Cloud Router type like XF_GATEWAY",
		},
		"location": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Fabric Cloud Router location",
			Elem: &schema.Resource{
				Schema: equinix_schema.LocationSch(),
			},
		},
		"project": {
			Type:        schema.TypeSet,
			Optional:    true,
			Computed:    true,
			Description: "Project information",
			Elem: &schema.Resource{
				Schema: equinix_schema.ProjectSch(),
			},
		},
		"account": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Customer account information that is associated with this Fabric Cloud Router",
			Elem: &schema.Resource{
				Schema: equinix_schema.AccountSch(),
			},
		},
		"order": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Order information related to this Fabric Cloud Router",
			Elem: &schema.Resource{
				Schema: equinix_schema.OrderSch(),
			},
		},
		"notifications": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Preferences for notifications on Fabric Cloud Router configuration or status changes",
			Elem: &schema.Resource{
				Schema: equinix_schema.NotificationSch(),
			},
		},
	}
}
