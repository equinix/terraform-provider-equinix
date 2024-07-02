# Create BGP peering configuration on a existing connection
# between network device and service provider

resource "equinix_network_bgp" "test" {
  connection_id      = "54014acf-9730-4b55-a791-459283d05fb1"
  local_ip_address   = "10.1.1.1/30"
  local_asn          = 12345
  remote_ip_address  = "10.1.1.2"
  remote_asn         = 66123
  authentication_key = "secret"
}
