---
subcategory: "Metal"
---

# Data Source: equinix_metal_metro

Provides an Equinix Metal metro datasource.

## Example Usage

```hcl
# Fetch a metro by code and show its ID

data "equinix_metal_metro" "sv" {
  code = "sv"
}

output "id" {
  value = data.equinix_metal_metro.sv.id
}
```


```hcl
# Verify that metro "sv" has capacity for provisioning 2 c3.small.x86 
  devices and 1 c3.medium.x86 device

data "equinix_metal_facility" "test" {
  code = "dc13"
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

* `code` - The metro code

Metros can be looked up by `code`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the metro
* `code` - The code of the metro
* `country` - The country of the metro
* `name` - The name of the metro
* `capacity` - (Optional) Ensure that queried metro has capacity for specified number of given plans
  - `plan` - device plan to check
  - `quantity` - number of device to check
