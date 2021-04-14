---
page_title: "Equinix Metal: metal_vlan"
subcategory: ""
description: |-
  Provides an Equinix Metal Virtual Network datasource. This can be used to read vlans.
---

# metal_vlan

Provides an Equinix Metal Virtual Network datasource. Users can look up VLANs based on either VLAN UUID, or project UUID and vxlan number.

## Example Usage

Fetch a vlan by ID:

```hcl
resource "metal_vlan" "foovlan" {
        project_id = local.project_id
        metro = "sv"
        vxlan = 5
}

data "metal_vlan" "dsvlan" {
        vlan_id = metal_vlan.foovlan.id
}
```

Fetch a vlan by project ID and vxlan

```hcl
resource "metal_vlan" "foovlan" {
        project_id = local.project_id
        metro = "sv"
        vxlan = 5
}

data "metal_vlan" "dsvlan" {
        project_id = local.project_id
        vxlan      = 5
}
```

## Argument Reference

The following arguments are supported:

* `vlan_id` - Metal UUID of the VLAN resource to look up
* `project_id` - UUID of parent project of the VLAN. Use together with the vxland number`
* `vxlan` - vxlan number of the VLAN to look up. Use together with the project_id`


## Attributes Reference

The following attributes are exported:

* `description` - Description text of the VLAN resource
* `facility` - Facility where the VLAN is deployed
* `metro` - Metro where teh VLAN is deployed
* `assigned_devices_ids` - List of device ID to which this VLAN is assinged
