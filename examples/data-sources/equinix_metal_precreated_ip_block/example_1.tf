# Create device in your project and then assign /64 subnet from precreated block
# to the new device

# Declare your project ID
locals {
  project_id = "<UUID_of_your_project>"
}

resource "equinix_metal_device" "web1" {
  hostname         = "web1"
  plan             = "c3.small.x86"
  metro            = "sv"
  operating_system = "ubuntu_20_04"
  billing_cycle    = "hourly"
  project_id       = local.project_id

}

data "equinix_metal_precreated_ip_block" "test" {
  metro          = "sv"
  project_id     = local.project_id
  address_family = 6
  public         = true
}

# The precreated IPv6 blocks are /56, so to get /64, we specify 8 more bits for network.
# The cirdsubnet interpolation will pick second /64 subnet from the precreated block.

resource "equinix_metal_ip_attachment" "from_ipv6_block" {
  device_id     = equinix_metal_device.web1.id
  cidr_notation = cidrsubnet(data.equinix_metal_precreated_ip_block.test.cidr_notation, 8, 2)
}
