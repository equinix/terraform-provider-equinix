---
subcategory: "Metal"
---

~> **Deprecation Notice** This data source has been deprecated. The Equinix Metal service reaches its end-of-life milestone on June 30, 2026. Scheduled for elimination in provider version 5.0.0, this data source will no longer be available. To sustain Metal operations until the platform concludes, continue with version 4.x of the Equinix Terraform provider. Additional sunset information is available at: https://docs.equinix.com/metal/


# equinix_metal_facility (Data Source)

> **Deprecated** Use `equinix_metal_metro` instead. For more information, refer to the [facility to metro migration guide](https://registry.terraform.io/providers/equinix/equinix/latest/docs/guides/migration_guide_facilities_to_metros_devices).

Provides an Equinix Metal facility datasource.

## Example Usage

```terraform
# Fetch a facility by code and show its ID

data "equinix_metal_facility" "ny5" {
  code = "ny5"
}

output "id" {
  value = data.equinix_metal_facility.ny5.id
}
```

```terraform
# Verify that facility "dc13" has capacity for provisioning 2 c3.small.x86 
  devices and 1 c3.medium.x86 device and has specified features

data "equinix_metal_facility" "test" {
  code = "dc13"

  features_required = ["backend_transfer", "global_ipv4"]

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

* `code` - (Required) The facility code to search for facilities.
* `features_required` - (Optional) Set of feature strings that the facility must have. Some possible values are `baremetal`, `ibx`, `storage`, `global_ipv4`, `backend_transfer`, `layer_2`.
* `capacity` - (Optional) One or more device plans for which the facility must have capacity.
  * `plan` - (Required) Device plan that must be available in selected location.
  * `quantity` - (Optional) Minimun number of devices that must be available in selected location. Default is `1`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the facility.
* `name` - The name of the facility.
* `features` - The features of the facility.
* `metro` - The metro code the facility is part of.
