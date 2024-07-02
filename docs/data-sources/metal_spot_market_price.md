---
subcategory: "Metal"
---

# equinix_metal_operating_system (Data Source)

Use this data source to get Equinix Metal Spot Market Price for a plan.

## Example Usage

Lookup by metro:

```terraform
data "equinix_metal_spot_market_price" "example" {
  metro    = "sv"
  plan     = "c3.small.x86"
}
```

## Argument Reference

The following arguments are supported:

* `plan` - (Required) Name of the plan.
* `facility` - (**Deprecated**) Name of the facility. Use metro instead; read the [facility to metro migration guide](https://registry.terraform.io/providers/equinix/equinix/latest/docs/guides/migration_guide_facilities_to_metros_devices)
* `metro` - (Optional) Name of the metro.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `price` - Current spot market price for given plan in given facility.
