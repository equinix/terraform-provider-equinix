---
page_title: "Equinix Metal: metal_spot_market_request"
subcategory: ""
description: |-
  Provides an Equinix Metal Spot Market Request Resource.
---

# metal\_spot\_market\_request

Provides an Equinix Metal Spot Market Request resource to allow you to
manage spot market requests on your account. For more detail on Spot Market, see [this article in Equinix Metal documentation](https://metal.equinix.com/developers/docs/deploy/spot-market/).

## Example Usage

```hcl
# Create a spot market request
resource "metal_spot_market_request" "req" {
  project_id    = local.project_id
  max_bid_price = 0.03
  facilities    = ["ewr1"]
  devices_min   = 1
  devices_max   = 1

  instance_parameters {
    hostname         = "testspot"
    billing_cycle    = "hourly"
    operating_system = "coreos_stable"
    plan             = "t1.small.x86"
  }
}
```

## Argument Reference

The following arguments are supported:

* `devices_max` - (Required) Maximum number devices to be created
* `devices_min` - (Required) Miniumum number devices to be created
* `max_bid_price` - (Required) Maximum price user is willing to pay per hour per device
* `facilities` - (Required) Facility IDs where devices should be created
* `instance_parameters` - (Required) Device parameters. See device resource for details
* `project_id` - (Required) Project ID
* `wait_for_devices` - (Optional) On resource creation - wait until all desired devices are active, on resource destruction - wait until devices are removed

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 60 mins) Used when creating the Spot Market Request and `wait_for_devices == true`)
* `delete` - (Defaults to 60 mins) Used when destroying the Spot Market Request and `wait for devices == true`

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Spot Market Request
