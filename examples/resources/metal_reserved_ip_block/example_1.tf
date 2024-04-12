# Allocate /31 block of max 2 public IPv4 addresses in Silicon Valley (sv) metro for myproject

resource "equinix_metal_reserved_ip_block" "two_elastic_addresses" {
  project_id = local.project_id
  metro      = "sv"
  quantity   = 2
}

# Allocate 1 floating IP in Silicon Valley (sv) metro

resource "equinix_metal_reserved_ip_block" "test" {
  project_id = local.project_id
  type       = "public_ipv4"
  metro      = "sv"
  quantity   = 1
}

# Allocate 1 global floating IP, which can be assigned to device in any metro

resource "equinix_metal_reserved_ip_block" "test" {
  project_id = local.project_id
  type       = "global_ipv4"
  quantity   = 1
}
