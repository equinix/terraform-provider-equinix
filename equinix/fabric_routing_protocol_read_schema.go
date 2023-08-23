package equinix

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/maps"
)

func readRoutingProtocolBaseSch() map[string]*schema.Schema {
	rpDataSchema := map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Routing Protocol configuration type",
		},
	}
	maps.Copy(rpDataSchema, readRoutingProtocolDirectTypeSch())
	maps.Copy(rpDataSchema, readRoutingProtocolBgpTypeSch())

	return rpDataSchema
}

func readRoutingProtocolDirectTypeSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Defines the routing protocol type as DIRECT",
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Routing Protocol name. An alpha-numeric 24 characters string which can include only hyphens and underscores",
		},
		"direct_ipv4": {
			Type:        schema.TypeSet,
			Required:    true,
			Description: "Routing Protocol Direct IPv4",
			Elem: &schema.Resource{
				Schema: readDirectConnectionIpv4Sch(),
			},
		},
		"direct_ipv6": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Routing Protocol Direct IPv6",
			Elem: &schema.Resource{
				Schema: readDirectConnectionIpv6Sch(),
			},
		},
	}
}

func readDirectConnectionIpv4Sch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"equinix_iface_ip": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Equinix side Interface IP address",
		},
	}
}

func readDirectConnectionIpv6Sch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"equinix_iface_ip": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Equinix side Interface IP address\n\n",
		},
	}
}

func readRoutingProtocolBgpTypeSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Defines the routing protocol type as DIRECT",
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Routing Protocol name. An alpha-numeric 24 characters string which can include only hyphens and underscores",
		},
		"bgp_ipv4": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Routing Protocol BGP IPv4",
			Elem: &schema.Resource{
				Schema: readBgpConnectionIpv4Sch(),
			},
		},
		"bgp_ipv6": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Routing Protocol BGP IPv6",
			Elem: &schema.Resource{
				Schema: readBgpConnectionIpv6Sch(),
			},
		},
		"customer_asn": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Customer-provided ASN",
		},
		"equinix_asn": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Equinix ASN",
		},
		"bgp_auth_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "BGP authorization key",
		},
		"bfd": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Bidirectional Forwarding Detection",
			Elem: &schema.Resource{
				Schema: readRoutingProtocolBfdSch(),
			},
		},
	}
}

func readBgpConnectionIpv4Sch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"customer_peer_ip": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Customer side peering ip",
		},
		"equinix_peer_ip": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Equinix side peering ip",
		},
		"enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Admin status for the BGP session",
		},
	}
}

func readBgpConnectionIpv6Sch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"customer_peer_ip": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Customer side peering ip",
		},
		"equinix_peer_ip": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Equinix side peering ip",
		},
		"enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Admin status for the BGP session",
		},
	}
}

func readRoutingProtocolBfdSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "Bidirectional Forwarding Detection enablement",
		},
		"interval": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     100,
			Description: "Interval range between the received BFD control packets",
		},
	}
}

func readRoutingProtocolOperationSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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

func readRoutingProtocolChangeSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		//"description": {
		//	Type: schema.TypeString,
		//	Computed: true,
		//	Description: "Details of latest Routing Protocol change",
		//},
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Uniquely identifies a change",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Type of change",
		},
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Routing Protocol Change URI",
		},
	}
}

func readFabricRoutingProtocolResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"connection_uuid": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Connection URI associated with Routing Protocol",
		},
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Routing Protocol URI information",
		},
		"type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Defines the routing protocol type like BGP or DIRECT",
		},
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Equinix-assigned routing protocol identifier",
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Routing Protocol name. An alpha-numeric 24 characters string which can include only hyphens and underscores",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Customer-provided Fabric Routing Protocol description",
		},
		"state": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Routing Protocol overall state",
		},
		"operation": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Routing Protocol type-specific operational data",
			Elem: &schema.Resource{
				Schema: readRoutingProtocolOperationSch(),
			},
		},
		"change": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Routing Protocol configuration Changes",
			Elem: &schema.Resource{
				Schema: readRoutingProtocolChangeSch(),
			},
		},
		"direct_ipv4": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Routing Protocol Direct IPv4",
			Elem: &schema.Resource{
				Schema: readDirectConnectionIpv4Sch(),
			},
		},
		"direct_ipv6": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Routing Protocol Direct IPv6",
			Elem: &schema.Resource{
				Schema: readDirectConnectionIpv6Sch(),
			},
		},
		"bgp_ipv4": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Routing Protocol BGP IPv4",
			Elem: &schema.Resource{
				Schema: readBgpConnectionIpv4Sch(),
			},
		},
		"bgp_ipv6": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Routing Protocol BGP IPv6",
			Elem: &schema.Resource{
				Schema: readBgpConnectionIpv6Sch(),
			},
		},
		"customer_asn": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Customer-provided ASN",
		},
		"equinix_asn": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Equinix ASN",
		},
		"bgp_auth_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "BGP authorization key",
		},
		"bfd": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Bidirectional Forwarding Detection",
			Elem: &schema.Resource{
				Schema: readRoutingProtocolBfdSch(),
			},
		},
		"change_log": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Captures Routing Protocol lifecycle change information",
			Elem: &schema.Resource{
				Schema: readChangeLogSch(),
			},
		},
	}
}
