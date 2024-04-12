---
subcategory: "Metal"
---

# equinix_metal_port_vlan_attachment (Resource)

Provides a resource to attach device ports to VLANs.

Device and VLAN must be in the same metro.

If you need this resource to add the port back to bond on removal, set `force_bond = true`.

To learn more about Layer 2 networking in Equinix Metal, refer to

* https://metal.equinix.com/developers/docs/networking/layer2/
* https://metal.equinix.com/developers/docs/networking/layer2-configs/

## Example Usage

### Hybrid network type

```terraform
resource "equinix_metal_vlan" "test" {
  description = "VLAN in New York"
  metro       = "ny"
  project_id  = local.project_id
}

resource "equinix_metal_device" "test" {
  hostname         = "test"
  plan             = "c3.small.x86"
  metro            = "ny"
  operating_system = "ubuntu_20_04"
  billing_cycle    = "hourly"
  project_id       = local.project_id
}

resource "equinix_metal_device_network_type" "test" {
  device_id = equinix_metal_device.test.id
  type      = "hybrid"
}

resource "equinix_metal_port_vlan_attachment" "test" {
  device_id = equinix_metal_device_network_type.test.id
  port_name = "eth1"
  vlan_vnid = equinix_metal_vlan.test.vxlan
}
```

### Layer 2 network

```terraform
resource "equinix_metal_device" "test" {
  hostname         = "test"
  plan             = "c3.small.x86"
  metro            = "ny"
  operating_system = "ubuntu_20_04"
  billing_cycle    = "hourly"
  project_id       = local.project_id
}

resource "equinix_metal_device_network_type" "test" {
  device_id = equinix_metal_device.test.id
  type      = "layer2-individual"
}

resource "equinix_metal_vlan" "test1" {
  description = "VLAN in New York"
  metro       = "ny"
  project_id  = local.project_id
}

resource "equinix_metal_vlan" "test2" {
  description = "VLAN in New Jersey"
  metro       = "ny"
  project_id  = local.project_id
}

resource "equinix_metal_port_vlan_attachment" "test1" {
  device_id = equinix_metal_device_network_type.test.id
  vlan_vnid = equinix_metal_vlan.test1.vxlan
  port_name = "eth1"
}

resource "equinix_metal_port_vlan_attachment" "test2" {
  device_id  = equinix_metal_device_network_type.test.id
  vlan_vnid  = equinix_metal_vlan.test2.vxlan
  port_name  = "eth1"
  native     = true
  depends_on = ["equinix_metal_port_vlan_attachment.test1"]
}
```

## Argument Reference

The following arguments are supported:

* `device_id` - (Required) ID of device to be assigned to the VLAN.
* `port_name` - (Required) Name of network port to be assigned to the VLAN.
* `vlan_vnid` - (Required) VXLAN Network Identifier.
* `force_bond` - (Optional) Add port back to the bond when this resource is removed. Default is `false`.
* `native` - (Optional) Mark this VLAN a native VLAN on the port. This can be used only if this assignment assigns second or further VLAN to the port. To ensure that this attachment is not first on a port, you can use `depends_on` pointing to another `equinix_metal_port_vlan_attachment`, just like in the layer2-individual example above.

## Attribute Referece

In addition to all arguments above, the following attributes are exported:

* `id` - UUID of device port used in the assignment.
* `vlan_id` - UUID of VLAN API resource.
* `port_id` - UUID of device port.
