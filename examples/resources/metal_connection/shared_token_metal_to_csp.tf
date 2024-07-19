resource "equinix_metal_connection" "example" {
    name               = "tf-metal-to-azure"
    project_id         = local.project_id
    type               = "shared"
    redundancy         = "redundant"
    metro              = "sv"
    speed              = "1000Mbps"
    service_token_type = "a_side"
    contact_email      = "username@example.com"
}

data "equinix_fabric_sellerprofile" "example" {
  name                     = "Azure ExpressRoute"
  organization_global_name = "Microsoft"
}

resource "equinix_fabric_connection" "example" {
  name              = "tf-metal-to-azure"
  profile_uuid      = data.equinix_fabric_sellerprofile.example.uuid
  speed             = azurerm_express_route_circuit.example.bandwidth_in_mbps
  speed_unit        = "MB"
  notifications     = ["example@equinix.com"]
  service_token     = equinix_metal_connection.example.service_tokens.0.id
  seller_metro_code = "AM"
  authorization_key = azurerm_express_route_circuit.example.service_key
  named_tag         = "PRIVATE"
  secondary_connection {
    name          = "tf-metal-to-azure-sec"
    service_token = equinix_metal_connection.example.service_tokens.1.id
  }
}
