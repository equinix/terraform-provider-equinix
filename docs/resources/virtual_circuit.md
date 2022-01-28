---
page_title: "Equinix Metal: virtual_circuit"
subcategory: ""
description: |-
  Create Equinix Fabric Virtual Circuit
---

# metal_virtual_circuit

Use this resource to associate VLAN with a Dedicated Port from [Equinix Fabric - software-defined interconnections](https://metal.equinix.com/developers/docs/networking/fabric/#associating-a-vlan-with-a-dedicated-port).

## Example Usage

Pick an existing Project and Connection, create a VLAN and use `metal_virtual_circuit` to associate it with a Primary Port of the Connection.

```hcl
locals {
	project_id = "52000fb2-ee46-4673-93a8-de2c2bdba33c"
	conn_id = "73f12f29-3e19-43a0-8e90-ae81580db1e0"
}

data "metal_connection" test {
	connection_id = local.conn_id
}

resource "metal_vlan" "test" {
	project_id = local.project_id
	metro      = data.metal_connection.test.metro
}

resource "metal_virtual_circuit" "test" {
	connection_id = local.conn_id
	project_id = local.project_id
	port_id = data.metal_connection.test.ports[0].id
	vlan_id = metal_vlan.test.id
	nni_vlan = 1056
}
```

## Argument Reference

* `connection_id` - (Required) UUID of Connection where the VC is scoped to
* `project_id` - (Required) UUID of the Project where the VC is scoped to
* `port_id` - (Required) UUID of the Connection Port where the VC is scoped to
* `nni_vlan` - (Required) Equinix Metal network-to-network VLAN ID
* `vlan_id` - (Required) UUID of the VLAN to associate
* `name` - (Optional) Name of the Virtual Circuit resource
* `facility` - (Optional) Facility where the connection will be created
* `description` - (Optional) Description for the Virtual Circuit resource
* `tags` - (Optional) Tags for the Virtual Circuit resource
* `speed` - (Optional) Speed of the Virtual Circuit resource

## Attributes Reference

* `status` - Status of the virtal circuit
* `vnid`
* `nni_nvid` - VLAN parameters, see the [documentation for Equinix Fabric](https://metal.equinix.com/developers/docs/networking/fabric/)
