# Create a Spot Market Request, and print public IPv4 of the created devices, if any.

resource "equinix_metal_spot_market_request" "req" {
  project_id       = local.project_id
  max_bid_price    = 0.1
  metro            = "ny"
  devices_min      = 2
  devices_max      = 2
  wait_for_devices = true

  instance_parameters {
    hostname         = "testspot"
    billing_cycle    = "hourly"
    operating_system = "ubuntu_20_04"
    plan             = "c3.small.x86"
  }
}

data "equinix_metal_spot_market_request" "dreq" {
  request_id = equinix_metal_spot_market_request.req.id
}

output "ids" {
  value = data.equinix_metal_spot_market_request.dreq.device_ids
}

data "equinix_metal_device" "devs" {
  count     = length(data.equinix_metal_spot_market_request.dreq.device_ids)
  device_id = data.equinix_metal_spot_market_request.dreq.device_ids[count.index]
}

output "ips" {
  value = [for d in data.equinix_metal_device.devs : d.access_public_ipv4]
}
