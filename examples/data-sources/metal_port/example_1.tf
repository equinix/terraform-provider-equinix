locals {
  project_id = "<UUID_of_your_project>"
}

resource "equinix_metal_device" "test" {
  hostname         = "tfacc-test-device-port"
  plan             = "c3.medium.x86"
  metro            = "sv"
  operating_system = "ubuntu_20_04"
  billing_cycle    = "hourly"
  project_id       = local.project_id
}

data "equinix_metal_port" "test" {
  device_id = equinix_metal_device.test.id
  name      = "eth0"
}
