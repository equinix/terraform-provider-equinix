# Reserve /30 block of max 2 public IPv4 addresses in metro ny for myproject
resource "equinix_metal_reserved_ip_block" "myblock" {
  project_id = local.project_id
  metro      = "ny"
  quantity   = 2
}

# Assign /32 subnet (single address) from reserved block to a device
resource "equinix_metal_ip_attachment" "first_address_assignment" {
  device_id = equinix_metal_device.mydevice.id
  # following expression will result to sth like "147.229.10.152/32"
  cidr_notation = join("/", [cidrhost(metal_reserved_ip_block.myblock.cidr_notation, 0), "32"])
}
