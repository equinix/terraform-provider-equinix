data "dns_a_record_set" "www" {
  host = "www.example.com"
}

data "equinix_metal_reserved_ip_block" "www" {
  project_id = local.my_project_id
  address = data.dns_a_record_set.www.addrs[0]
}

resource "equinix_metal_device" "www" {
  project_id = local.my_project_id
  [...]
  ip_address {
    type = "public_ipv4"
    reservation_ids = [data.equinix_metal_reserved_ip_block.www.id]
  }
}
