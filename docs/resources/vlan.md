---
page_title: "Equinix Metal: metal_vlan"
subcategory: ""
description: |-
  Provides a resource for Equinix Metal Virtual Network.
---

# metal_vlan

Provides a resource to allow users to manage Virtual Networks in their projects.

To learn more about Layer 2 networking in Equinix Metal, refer to

* <https://metal.equinix.com/developers/docs/networking/layer2/>
* <https://metal.equinix.com/developers/docs/networking/layer2-configs/>

## Example Usage

```hcl
# Create a new VLAN in datacenter "ewr1"
resource "metal_vlan" "vlan1" {
  description = "VLAN in New Jersey"
  facility    = "ewr1"
  project_id  = local.project_id
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) ID of parent project
* `facility` - (Required) Facility where to create the VLAN
* `description` - Description string

## Attributes Reference

The following attributes are exported:

* `vxlan` - VXLAN segment ID
