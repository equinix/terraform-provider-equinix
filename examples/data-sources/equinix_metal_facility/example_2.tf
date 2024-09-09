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
