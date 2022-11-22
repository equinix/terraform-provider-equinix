---
page_title: "Metal Device Network Types"
---


# Network types

Server network types, such as Layer-2, Layer-3, and Hybrid may be familiar to users of the Equinix Metal Portal. In the Portal, you can toggle the network type with a click of the UI. To take advantage of these features in Terraform, which closely follows the Equinix Metal API, it is important to understand that the network type is a composite string value determined by one or more port bonding, addressing, and VLAN attachment configurations. To change the network type, you must change these underlying properties of the port(s).

For more details, see the Equinix Metal documentation on [Network Configuration Types](https://metal.equinix.com/developers/docs/layer2-networking/overview/#network-configuration-types).

This Terraform provider offers two ways to define the network type.

* [`equinix_metal_port`](#Metal-Port)
* [`equinix_metal_device_network_type`](#Metal-Device-Network-Type)

## Metal Port

The [`equinix_metal_port`](../resources/metal_port.md) resource exposes all of the features needed to affect the network type of a device or port pairing.

Following are examples of how the `equinix_metal_port` resource can be used to configure various network types, assuming that `local.bond0_id` is the UUID of the bond interface containing `eth1` and `local.eth1_id` is the UUID of the `eth1` interface.  These could represent the ports of `equinix_metal_device` resources or data sources.

### Layer 3 Port

Layer-3 (Bonded) is the default port configuration on Equinix Metal devices. The following is provided to illustrate the usage of the `equinix_metal_port` resource. This HCL should not be needed in practice, however it may be useful in some configurations to assert the correct mode is set, port by port, on imported `equinix_metal_device` resources or data-sources.

```hcl
resource "equinix_metal_port" "bond0" {
  port_id = local.bond0_id
  bonded = true
}

resource "equinix_metal_port" "eth1" {
  port_id = local.eth1_id
  bonded = true
}
```

### Layer 2 Unbonded Port

This example configures an Equinix Metal server with a [pure layer 2 unbonded](https://deploy.equinix.com/developers/docs/metal/layer2-networking/layer2-mode/#:~:text=Layer%202%20Unbonded%20Mode) network configuration and adds two VLANs to its `eth1` port; one of them set as the [native VLAN](https://deploy.equinix.com/developers/docs/metal/layer2-networking/native-vlan/). Notice the `depends_on` meta-argument in the `equinix_metal_port.eth1` resource and the `reset_on_delete` attribute in both ports’ configuration. The `reset_on_delete` will set the port to the default settings (layer3 bonded without VLANs attached) before the terraform resource delete/destroy. It is recommended to use the `depends_on` argument here to ensure that the port resources with attached VLANs are reset first, since all VLANs must be detached before re-bonding the ports.

```hcl
resource "equinix_metal_port" "bond0" {
  port_id = local.bond0_id
  layer2 = true
  bonded = false
  reset_on_delete = true
}

resource "equinix_metal_port" "eth1" {
  port_id = local.eth1_id
  bonded  = false
  reset_on_delete = true
  vlan_ids = [equinix_metal_vlan.test1.id, equinix_metal_vlan.test2.id]
  native_vlan_id = equinix_metal_vlan.test1.id
  depends_on = [
	  equinix_metal_port.bond0,
  ]
}

resource "equinix_metal_vlan" "test1" {
  description = "test"
  metro = "sv"
  project = equinix_metal_project.test.id
}

resource "equinix_metal_vlan" "test2" {
  description = "test"
  metro = "sv"
  project = equinix_metal_project.test.id
}
```

### Layer 2 Bonded Port

```hcl
resource "equinix_metal_port" "bond0" {
  port_id = local.bond0_id
  layer2 = true
  bonded = true
}
```

### Hybrid Unbonded Port

```hcl
resource "equinix_metal_port" "eth1" {
  port_id = local.eth1_id
  bonded = false
}
```

### Hybrid Bonded Port

```hcl
resource "equinix_metal_port" "bond0" {
  port_id = local.bond0_id
  layer2 = false
  bonded = true
  vlan_ids = [equinix_metal_vlan.test.id]
}

resource "equinix_metal_vlan" "test" {
  description = "test"
  metro = "sv"
  project = equinix_metal_project.test.id
}
```

### Accessing Port IDs

The port ID value can be obtained from a `equinix_metal_device` using a [`for` expression](https://www.terraform.io/docs/language/expressions/for.html).

Assuming a `equinix_metal_device` exists with the resource name `test`:

```hcl
locals {
  bond0_id = [for p in equinix_metal_device.test.ports: p.id if p.name == "bond0"][0]
   eth1_id = [for p in equinix_metal_device.test.ports: p.id if p.name == "eth1"][0]
}
```

## Metal Device Network Type

The [`equinix_metal_device_network_type`](../resources/metal_device_network_type.md) takes a named network type with any mode required parameters and converts a device to the named network type.  This resource simulated the network type interface for Devices in the Equinix Metal Portal. That interface changed when additional network types were introduced with more diverse port configurations.

When using this resource, keep in mind:

* this resource is not guaranteed to work in devices with more than two ethernet ports
* it may not be able to express all possible port configurations
* subsequent changes to the network configuration may cause this device to detect changes that can not be reconciled without intervention
* `equinix_metal_device_network_type` resources should not be used on devices with ports being controlled with `equinix_metal_port` resources

### Hybrid (Unbonded) Device

This example create one c3.small device and puts it into [hybrid (unbonded) network mode](https://metal.equinix.com/developers/docs/layer2-networking/hybrid-unbonded-mode/).

```hcl
resource "equinix_metal_device" "test" {
  hostname         = "tfacc-device-port-vlan-attachment-test"
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
```

### Hybrid (Unbonded) Device with a VLAN

This example create two devices in [hybrid (unbonded) mode](https://metal.equinix.com/developers/docs/layer2-networking/hybrid-unbonded-mode/) and adds a VLAN to their eth1 ports.

```hcl
locals {
    project_id = "<uuid>"
    device_count = 2
}

resource "equinix_metal_vlan" "test" {
  metro       = "ny"
  project_id  = local.project_id
}


resource "equinix_metal_device" "test" {
  count            = local.device_count
  hostname         = "test${count.index}"
  plan             = "c3.small.x86"
  metro            = "ny"
  operating_system = "ubuntu_20_04"
  billing_cycle    = "hourly"
  project_id       = local.project_id
}

resource "equinix_metal_device_network_type" "test" {
  count     = local.device_count
  device_id = equinix_metal_device.test[count.index].id
  type      = "hybrid"
}


resource "equinix_metal_port_vlan_attachment" "test" {
  count     = local.device_count
  device_id = equinix_metal_device_network_type.test[count.index].id
  port_name = "eth1"
  vlan_vnid = equinix_metal_vlan.test.vxlan
}
```

### Hybrid (Bonded) Device

This example create one c3.small device and puts it into [hybrid-bonded network mode](https://metal.equinix.com/developers/docs/layer2-networking/hybrid-bonded-mode/). Notice, the default network type of `layer3` can be used with VLAN attachments without reconfiguring the device ports.

```hcl
resource "equinix_metal_device" "test" {
  hostname         = "tfacc-device-port-vlan-attachment-test"
  plan             = "c3.small.x86"
  metro            = "ny"
  operating_system = "ubuntu_20_04"
  billing_cycle    = "hourly"
  project_id       = local.project_id
}

resource "equinix_metal_vlan" "test" {
  metro       = "ny"
  project_id  = local.project_id
}

resource "equinix_metal_port_vlan_attachment" "test" {
  count     = local.device_count
  device_id = equinix_metal_device.test.id
  port_name = "bond0"
  vlan_vnid = equinix_metal_vlan.test.vxlan
}
```
