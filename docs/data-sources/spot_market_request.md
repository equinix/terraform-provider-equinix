---
page_title: "Equinix Metal: metal_spot_market_request"
subcategory: ""
description: |-
  Provides a datasource for existing Spot Market Requests in the Equinix Metal host.
---

# metal_spot_market_request

Provides an Equinix Metal spot_market_request datasource. The datasource will contain list of device IDs created by referenced Spot Market Request.

## Example Usage

```hcl
# Create a Spot Market Request, and print public IPv4 of the created devices, if any.

resource "metal_spot_market_request" "req" {
  project_id       = local.project_id
  max_bid_price    = 0.1
  facilities       = ["ewr1"]
  devices_min      = 2
  devices_max      = 2
  wait_for_devices = true

  instance_parameters {
    hostname         = "testspot"
    billing_cycle    = "hourly"
    operating_system = "ubuntu_16_04"
    plan             = "t1.small.x86"
  }
}

data "metal_spot_market_request" "dreq" {
  request_id = metal_spot_market_request.req.id
}

output "ids" {
  value = data.metal_spot_market_request.dreq.device_ids
}

data "metal_device" "devs" {
  count     = length(data.metal_spot_market_request.dreq.device_ids)
  device_id = data.metal_spot_market_request.dreq.device_ids[count.index]
}

output "ips" {
  value = [for d in data.metal_device.devs : d.access_public_ipv4]
}
```

With the code as `main.tf`, first create the spot market request:

```
terraform apply -target metal_spot_market_request.req
```

When the terraform run ends, run a full apply, and the IPv4 addresses will be printed:

```
$ terraform apply

[...]

ips = [
  "947.85.199.231",
  "947.85.194.181",
]
```

## Argument Reference

The following arguments are supported:

* `request_id` - (Required) The id of the Spot Market Request

## Attributes Reference

The following attributes are exported:

* `device_ids` - List of IDs of devices spawned by the referenced Spot Market Request
