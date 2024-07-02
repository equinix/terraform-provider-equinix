# Create a spot market request
resource "equinix_metal_spot_market_request" "req" {
  project_id    = local.project_id
  max_bid_price = 0.03
  metro         = "ny"
  devices_min   = 1
  devices_max   = 1

  instance_parameters {
    hostname         = "testspot"
    billing_cycle    = "hourly"
    operating_system = "ubuntu_20_04"
    plan             = "c3.small.x86"
  }
}
