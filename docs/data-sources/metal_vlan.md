---
subcategory: "Metal"
---

# equinix_metal_vlan (Data Source)

Provides an Equinix Metal Virtual Network datasource. VLANs data sources can be searched by VLAN UUID, or project UUID and vxlan number.

## Example Usage

Fetch a vlan by ID:

```terraform
resource "equinix_metal_vlan" "foovlan" {
  project_id = local.project_id
  metro = "sv"
  vxlan = 5
}

data "equinix_metal_vlan" "dsvlan" {
  vlan_id = equinix_metal_vlan.foovlan.id
}
```

Fetch a vlan by project ID, vxlan and metro

```terraform
data "equinix_metal_vlan" "dsvlan" {
  project_id = local.project_id
  vxlan      = 5
  metro      = "sv"
}
```

## Argument Reference

The following arguments are supported:

* `vlan_id` - (Optional) Metal UUID of the VLAN resource to look up.
* `project_id` - (Optional) UUID of parent project of the VLAN. Use together with the vxlan number and metro or facility.
* `vxlan` - (Optional) vxlan number of the VLAN to look up. Use together with the project_id and metro or facility.
* `facility` - (Optional) Facility where the VLAN is deployed. Deprecated, see https://feedback.equinixmetal.com/changelog/bye-facilities-hello-again-metros
* `metro` - (Optional) Metro where the VLAN is deployed.

-> **NOTE:** You must set either `vlan_id` or a combination of `vxlan`, `project_id`, and, `metro` or `facility`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `description` - Description text of the VLAN resource.
* `assigned_devices_ids` - List of device ID to which this VLAN is assigned.
