package equinix

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var readPortDeviceRes = &schema.Resource{
	Schema: readPortDeviceSch(),
}

func readPortDeviceSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port name",
		},
		"redundancy": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Port device redundancy",
			Elem: &schema.Resource{
				Schema: readRedundancySch(),
			},
		},
	}
}

var readPortDeviceRedundancyRes = &schema.Resource{
	Schema: readRedundancySch(),
}

var readPortInterfaceRes = &schema.Resource{
	Schema: readPortInterfaceSch(),
}

func readPortInterfaceSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port interface type",
		},
		"if_index": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port interface index",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port interface name",
		},
	}
}

var readPortEncapsulationRes = &schema.Resource{
	Schema: readFabricPortEncapsulation(),
}

func readFabricPortEncapsulation() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port encapsulation protocol type",
		},
		"tag_protocol_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port encapsulation Tag Protocol Identifier",
		},
	}
}

var readPortLagRes = &schema.Resource{
	Schema: readFabricPortLag(),
}

func readFabricPortLag() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "If LAG enabled",
		},
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Link aggregation group identifier",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Link aggregation group name",
		},
		"member_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "LAG port status",
		},
	}
}

var readPortSettingsRes = &schema.Resource{
	Schema: readFabricPortSettings(),
}

func readFabricPortSettings() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"port_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port type",
		},
	}
}

func readPortOperationSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"operational_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port operation status",
		},
		"connection_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Total number of current connections",
		},
		"op_status_changed_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Date and time at which port availability changed",
		},
	}
}

var readPortsRedundancyRes = &schema.Resource{
	Schema: readPortsRedundancySch(),
}

func readPortsRedundancySch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Access point redundancy",
		},
		"group": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port redundancy group",
		},
		"priority": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Priority type-Primary or Secondary",
		},
	}
}

var readPortTetherRes = &schema.Resource{
	Schema: readFabricPortTether(),
}

func readFabricPortTether() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cross_connect_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port cross connect identifier",
		},
		"cabinet_number": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port cabinet number",
		},
		"system_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port system name",
		},
		"patch_panel": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port patch panel",
		},
		"patch_panel_port_a": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port patch panel port A",
		},
		"patch_panel_port_b": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port patch panel port B",
		},
	}
}

func readFabricPortResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port type",
		},
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port URI information",
		},
		"uuid": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Equinix-assigned port identifier",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port name",
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port description",
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
				Schema: readPortDeviceSch(),
			},
		},
		"redundancy": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Port redundancy information",
			Elem: &schema.Resource{
				Schema: readPortsRedundancySch(),
			},
		},
		"encapsulation": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Port encapsulation protocol",
			Elem: &schema.Resource{
				Schema: readFabricPortEncapsulation(),
			},
		},
		"lag_enabled": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Port Lag",
		},
	}
}

func readGetPortsByNameQueryParamSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Query Parameter to Get Ports By Name",
		},
	}
}

func readFabricPortResourceSchemaUpdated() map[string]*schema.Schema {
	sch := readFabricPortResourceSchema()
	sch["uuid"].Optional = true
	sch["uuid"].Required = false
	return sch
}

func readFabricPortsResponseSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"data": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "List of Ports",
			Elem: &schema.Resource{
				Schema: readFabricPortResourceSchemaUpdated(),
			},
		},
		"filters": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "name",
			Elem: &schema.Resource{
				Schema: readGetPortsByNameQueryParamSch(),
			},
		},
	}
}
