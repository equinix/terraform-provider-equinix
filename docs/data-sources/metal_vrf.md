---
subcategory: "Metal"
---

~> **Deprecation Notice** This data source will not be available in future versions. The Equinix Metal platform concludes on June 30, 2026. This data source is scheduled for removal in the next major version (5.0.0). For sustained access to Metal services through the sunset date, utilize version 4.x of the Equinix Terraform provider. Consult the documentation at: https://docs.equinix.com/metal/


# equinix_metal_virtual_circuit (Data Source)

Use this data source to retrieve a VRF resource.

See the [Virtual Routing and Forwarding documentation](https://docs.equinix.com/metal/networking/vrf/) for product details and API reference material.

## Example Usage

```terraform
data "equinix_metal_vrf" "example_vrf" {
  vrf_id = "48630899-9ff2-4ce6-a93f-50ff4ebcdf6e"
}
```

## Argument Reference

The following arguments are supported:

* `vrf_id` - (Required) ID of the VRF resource

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `name` - User-supplied name of the VRF, unique to the project
* `metro` - Metro ID or Code where the VRF will be deployed.
* `project_id` - Project ID where the VRF will be deployed.
* `description` - Description of the VRF.
* `local_asn` - The 4-byte ASN set on the VRF.
* `ip_ranges` - All IPv4 and IPv6 Ranges that will be available to BGP Peers. IPv4 addresses must be /8 or smaller with a minimum size of /29. IPv6 must be /56 or smaller with a minimum size of /64. Ranges must not overlap other ranges within the VRF.
