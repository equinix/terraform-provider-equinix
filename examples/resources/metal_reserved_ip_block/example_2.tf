# Allocate /31 block of max 2 public IPv4 addresses in Silicon Valley (sv) metro
resource "equinix_metal_reserved_ip_block" "example" {
  project_id = local.project_id
  metro      = "sv"
  quantity   = 2
}

# Run a device with both public IPv4 from the block assigned

resource "equinix_metal_device" "nodes" {
  project_id       = local.project_id
  metro            = "sv"
  plan             = "c3.small.x86"
  operating_system = "ubuntu_20_04"
  hostname         = "test"
  billing_cycle    = "hourly"

  ip_address {
    type            = "public_ipv4"
    cidr            = 31
    reservation_ids = [equinix_metal_reserved_ip_block.example.id]
  }

  ip_address {
    type = "private_ipv4"
  }
}
