---
subcategory: "Metal"
---

# Resource: equinix\_metal\_spot\_market\_request

Provides an Equinix Metal Spot Market Request resource to allow you to
manage spot market requests on your account. For more detail on Spot Market, see [this article in Equinix Metal documentation](https://metal.equinix.com/developers/docs/deploy/spot-market/).

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

* `devices_max` - (Required) Maximum number devices to be created
* `devices_min` - (Required) Miniumum number devices to be created
* `max_bid_price` - (Required) Maximum price user is willing to pay per hour per device
* `project_id` - (Required) Project ID
* `wait_for_devices` - (Optional) On resource creation - wait until all desired devices are active, on resource destruction - wait until devices are removed
* `facilities` - (Optional) Facility IDs where devices should be created
* `metro` - (Optional) Metro where devices should be created
* `locked` - (Optional) Blocks deletion of the SpotMarketRequest device until the lock is disabled
* `instance_parameters` - (Required) Parameters for devices provisioned from this request. You can find the parameter description from the [equinix_metal_device doc](metal_device.md).
  * `billing_cycle`
  * `plan`
  * `operating_system`
  * `hostname`
  * `termintation_time`
  * `always_pxe`
  * `description`
  * `features`
  * `locked`
  * `project_ssh_keys`
  * `user_ssh_keys`
  * `userdata`
  * `customdata`
  * `ipxe_script_url`
  * `tags`

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 60 mins) Used when creating the Spot Market Request and `wait_for_devices == true`)
* `delete` - (Defaults to 60 mins) Used when destroying the Spot Market Request and `wait for devices == true`

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Spot Market Request
* `facilities` - The facilities where the Spot Market Request is applied. This is computed when `metro` is set or no specific location was requested.

## Import

This resource can be imported using an existing spot market request ID:

```sh
terraform import equinix_metal_spot_market_request {existing_spot_market_request_id}
```
