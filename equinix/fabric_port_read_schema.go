package equinix

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func readPortDeviseSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Port Name",
		},
		"redundancy": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Port Device Redundancy",
			Elem: &schema.Resource{
				Schema: readRedundancySch(),
			},
		},
	}
}

func readPortOperationSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"operational_status": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Port Type",
		},
		"connection_count": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Total number of current connections",
		},
		"op_status_changed_at": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Date and time at which port availability changed",
		},
	}
}

func readFabricPortTether() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cross_onnect_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Port cross connect identifier",
		},
		"cabinet_number": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Port cabinet number",
		},
		"system_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Port system name",
		},
		"patch_panel": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Port patch panel",
		},
		"patch_panel_port_a": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Port patch panel port A",
		},
		"patch_panel_port_b": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Port patch panel port B",
		},
	}
}

func readFabricPortResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Port Type",
		},
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port URI information",
		},
		"id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Port ID",
		},
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Equinix-assigned port identifier",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port Name",
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port description",
		},
		"cvpid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Unique ID for a virtual port",
		},
		"state": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port state",
		},
		"operation": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Port specific operational data",
			Elem: &schema.Resource{
				Schema: readPortOperationSch(),
			},
		},
		"bandwidth": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Port bandwidth in Mbps",
		},
		"available_bandwidth": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Port available bandwidth in Mbps",
		},
		"used_bandwidth": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Port used bandwidth in Mbps",
		},
		"service_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port service type",
		},
		"order": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Order related to this Port",
			Elem: &schema.Resource{
				Schema: readOrderSch(),
			},
		},
		"account": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Customer account information that is associated with this port",
			Elem: &schema.Resource{
				Schema: readAccountSch(),
			},
		},
		"change_log": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Captures port lifecycle change information",
			Elem: &schema.Resource{
				Schema: readChangeLogSch(),
			},
		},
		"location": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Port location information",
			Elem: &schema.Resource{
				Schema: readLocationSch(),
			},
		},
		"device": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Port device",
			Elem: &schema.Resource{
				Schema: readPortDeviseSch(),
			},
		},
		"interface": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "", //TODO
			Elem: &schema.Resource{
				//TODO PortInterface
			},
		},
		"redundancy": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Port Redundancy Information",
			Elem: &schema.Resource{
				Schema: readPortRedundancySch(),
			},
		},
		"tether": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "",
			Elem: &schema.Resource{
				Schema: readFabricPortTether(),
			},
		},
		"encapsulation": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "",
			Elem:        &schema.Resource{}, //TODO  PortEncapsulation
		},
		"lag": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "",
			Elem:        &schema.Resource{}, //TODO  PortLag
		},
		"asn": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Port ASN",
		},
		"settings": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "",
			Elem:        &schema.Resource{}, //TODO  PortSettings
		},
		"physical_port_quantity": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Number of physical ports",
		},
		"physical_ports": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Physical ports that implement this port",
			Elem: &schema.Resource{
				Schema: readFabricPortResourceSchema(),
			},
		},
	}
}
