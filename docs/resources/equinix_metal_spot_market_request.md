---
subcategory: "Metal"
---

# equinix_metal_spot_market_request (Resource)

Provides an Equinix Metal Spot Market Request resource to allow you to
manage spot market requests on your account. For more detail on Spot Market,
see [this article in Equinix Metal documentation](https://metal.equinix.com/developers/docs/deploy/spot-market/).

## Example Usage

```hcl
# Create a spot market request
resource "equinix_metal_spot_market_request" "req" {
  project_id    = local.project_id
  max_bid_price = 0.03
  facilities    = ["ny5"]
  devices_min   = 1
  devices_max   = 1

  instance_parameters {
    hostname         = "testspot"
    billing_cycle    = "hourly"
    operating_system = "ubuntu_20_04"
    plan             = "c3.small.x86"
  }
}
```

## Argument Reference

The following arguments are supported:

* `devices_max` - (Required) Maximum number devices to be created.
* `devices_min` - (Required) Miniumum number devices to be created.
* `max_bid_price` - (Required) Maximum price user is willing to pay per hour per device.
* `project_id` - (Required) Project ID.
* `wait_for_devices` - (Optional) On resource creation wait until all desired devices are active.
On resource destruction wait until devices are removed.
* `facilities` - (Optional) Facility IDs where devices should be created.
* `end_at` - (Optional) Deadline for the request. For example \"2021-09-03T16:32:00+03:00\". If you don't supply timezone info, timestamp is assumed to be in UTC."
* `metro` - (Optional) Metro where devices should be created.
* `locked` - (Optional) Blocks deletion of the SpotMarketRequest device until the lock is disabled.
* `instance_parameters` - (Required) Key/Value pairs of parameters for devices provisioned from
this request. Valid keys are: `billing_cycle`, `plan`, `operating_system`, `hostname`, `always_pxe` (Defaults to false), `description`, `features`, `locked`, `project_ssh_keys`,
`user_ssh_keys`, `userdata`, `customdata`, `ipxe_script_url`, `tags`. You can find each parameter
description in [equinix_metal_device](equinix_metal_device.md) docs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Spot Market Request.

The `instance_parameters` block additionally exports:

* `termination_time` - Timestamp for device termination shown once the device has been outbid.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/configuration/resources#operation-timeouts) for certain actions:

* `create` - (Defaults to 60 mins) Used when creating the Spot Market Request and `wait_for_devices`
is set to `true`.
* `delete` - (Defaults to 60 mins) Used when destroying the Spot Market Request and `wait_for_devices`
is set to `true`.

## Import

This resource can be imported using an existing spot market request ID:

```sh
terraform import equinix_metal_spot_market_request {existing_spot_market_request_id}
```
