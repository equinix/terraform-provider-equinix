---
subcategory: "Metal"
---

~> **Deprecation Notice** Equinix Metal will reach end of life on June 30, 2026. All Metal data sources will be removed in version 5.0.0 of this provider. Use version 4.x of this provider for continued use through sunset. See https://docs.equinix.com/metal/ for more information.


# equinix_metal_metro (Data Source)

Provides an Equinix Metal metro datasource.

## Example Usage

```terraform
# Fetch a metro by code and show its ID

data "equinix_metal_metro" "sv" {
  code = "sv"
}

output "id" {
  value = data.equinix_metal_metro.sv.id
}
```

```terraform
# Verify that metro "sv" has capacity for provisioning 2 c3.small.x86 
  devices and 1 c3.medium.x86 device

data "equinix_metal_metro" "test" {
  code = "sv"

  capacity {
    plan = "c3.small.x86"
    quantity = 2
  }

  capacity {
    plan = "c3.medium.x86"
    quantity = 1
  }
}
```

## Argument Reference

The following arguments are supported:

* `code` - (Required) The metro code to search for.
* `capacity` - (Optional) One or more device plans for which the metro must have capacity.
  * `plan` - (Required) Device plan that must be available in selected location.
  * `quantity` - (Optional) Minimum number of devices that must be available in selected location. Default is `1`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the metro.
* `name` - The name of the metro.
* `country` - The country of the metro.
