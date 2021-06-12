---
page_title: "Equinix Metal: precreated_port"
subcategory: ""
description: |-
  Manipulate device ports
---

# metal_port

Use this resource to set up network ports on an Equnix Metal device. This resource can control both physical and bond ports.

This Terraform resource doesn't create an API resource in Equinix Metal, but rather provides finer control for (Layer 2 networking)[https://metal.equinix.com/developers/docs/layer2-networking/].

The port resource referred is created together with device and accessible either via the device resource or over `/port/<uuid>` API path.

## Network types

`metal_port` can be used to configure the [Network Configuration Types](https://metal.equinix.com/developers/docs/layer2-networking/overview/#network-configuration-types). Network types are bond port configurations.

Following are examples of how the `metal_port` resource can be used to configure various network types, assuming that `bond0_id` is the UUID of the bond interface containing `eth1`.

### Layer 3

```hcl
resource "metal_port" "bond0" {
  port_id = local.bond0_id
  bonded = true
}

resource "metal_port" "eth1" {
  port_id = local.eth1_id
  bonded = true
}
```

### Layer 2 Unbonded

```hcl
resource "metal_port" "bond0" {
  port_id = local.bond0_id
  layer2 = true
  bonded = false
}
```

### Layer 2 Bonded

```hcl
resource "metal_port" "bond0" {
  port_id = local.bond0_id
  layer2 = true
  bonded = true
}
```

### Hybrid Unbonded

```hcl

resource "metal_port" "bond0" {
  port_id = local.bond0_id
  layer2 = false
  bonded = true
  depends_on = [metal_port.eth1]
}

resource "metal_port" "eth1" {
  port_id = local.eth1_id
  bonded = false
}
```

### Hybrid Bonded

```
resource "metal_port" "bond0" {
  port_id = local.bond0_id
  layer2 = false
  bonded = true
  vlan_ids = [metal_vlan.test.id]
}

resource "metal_vlan" "test" {
  description = "test"
  metro = "sv"
  project = metal_project.test.id
}

resource "metal_port" "eth1" {
  port_id = local.eth1_id
  bonded = false
}
```


## Argument Reference

* `port_id` - (Required) ID of the port to read
* `bonded` - (Required) Whether the port should be bonded
* `layer2` - (Optional) Whether to put the port to Layer 2 mode, valid only for bond ports
* `vlan_ids` - (Optional) List off VLAN UUIDs to attach to the port
* `native_vlan_id` - (Optional) UUID of a VLAN to assign as a native VLAN. It must be one of attached VLANs (from `vlan_ids` parameter), valid only for physical (non-bond) ports
* `reset_on_delete` - (Optional) Flag indicating whether to reset port to default settings. For a bond port it means layer3 without VLANs attached, physical ports will be bonded without native VLAN and VLANs attached


## Attributes Reference

* `name` - Name of the port, e.g. `bond0` or `eth0`
* `network_type` - One of layer2-bonded, layer2-individual, layer3, hybrid, hybrid-bonded
* `type` - Type is either "NetworkBondPort" for bond ports or "NetworkPort" for bondable ethernet ports
* `mac` - MAC address of the port
* `bond_id` - UUID of the bond port"
* `bond_name` - Name of the bond port
* `bonded` - Flag indicating whether the port is bonded
* `disbond_supported` - Flag indicating whether the port can be removed from a bond

