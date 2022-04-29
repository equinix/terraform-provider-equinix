---
subcategory: "Metal"
---

# equinix_metal_virtual_circuit (Resource)

Use this resource to associate VLAN with a Dedicated Port from
[Equinix Fabric - software-defined interconnections](https://metal.equinix.com/developers/docs/networking/fabric/#associating-a-vlan-with-a-dedicated-port).

## Example Usage

Pick an existing Project and Connection, create a VLAN and use `equinix_metal_virtual_circuit`
to associate it with a Primary Port of the Connection.

```hcl
locals {
  project_id = "52000fb2-ee46-4673-93a8-de2c2bdba33c"
  conn_id = "73f12f29-3e19-43a0-8e90-ae81580db1e0"
}

data "equinix_metal_connection" test {
  connection_id = local.conn_id
}

resource "equinix_metal_vlan" "test" {
  project_id = local.project_id
  metro      = data.equinix_metal_connection.test.metro
}

resource "equinix_metal_virtual_circuit" "test" {
  connection_id = local.conn_id
  project_id = local.project_id
  port_id = data.equinix_metal_connection.test.ports[0].id
  vlan_id = equinix_metal_vlan.test.id
  nni_vlan = 1056
}
```

## Argument Reference

The following arguments are supported:

* `connection_id` - (Required) UUID of Connection where the VC is scoped to.
* `project_id` - (Required) UUID of the Project where the VC is scoped to.
* `port_id` - (Required) UUID of the Connection Port where the VC is scoped to.
* `nni_vlan` - (Required) Equinix Metal network-to-network VLAN ID.
* `vlan_id` - (Required) UUID of the VLAN to associate.
* `name` - (Optional) Name of the Virtual Circuit resource.
* `facility` - (Optional) Facility where the connection will be created.
* `description` - (Optional) Description for the Virtual Circuit resource.
* `tags` - (Optional) Tags for the Virtual Circuit resource.
* `speed` - (Optional) Speed of the Virtual Circuit resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `status` - Status of the virtal circuit.
* `vnid` - VNID VLAN parameter, see the [documentation for Equinix Fabric](https://metal.equinix.com/developers/docs/networking/fabric/).
* `nni_vnid` - NNI VLAN parameters, see the [documentation for Equinix Fabric](https://metal.equinix.com/developers/docs/networking/fabric/).

## Import

This resource can be imported using an existing Virtual Circuit ID:

```sh
terraform import equinix_metal_virtual_circuit {existing_id}
```
