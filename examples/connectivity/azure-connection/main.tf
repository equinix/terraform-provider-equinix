provider "equinix" {
  client_id     = var.equinix_client_id
  client_secret = var.equinix_client_secret
}

provider "azurerm" {
  features {}
}

data "equinix_ecx_l2_sellerprofile" "azure" {
  name                     = "Azure ExpressRoute"
  organization_global_name = "Microsoft"
}

data "equinix_ecx_port" "dot1q-1-pri" {
  name = var.equinix_pri_port_name
}

data "equinix_ecx_port" "dot1q-1-sec" {
  name = var.equinix_sec_port_name
}

resource "azurerm_resource_group" "demo" {
  name     = "TFDemo"
  location = var.azure_location
}

resource "azurerm_express_route_circuit" "demo" {
  name                  = "TFDemoExpressRoute"
  resource_group_name   = azurerm_resource_group.demo.name
  location              = azurerm_resource_group.demo.location
  service_provider_name = "Equinix"
  peering_location      = var.azure_peering_location
  bandwidth_in_mbps     = 50
  sku {
    tier   = "Premium"
    family = "UnlimitedData"
  }
  allow_classic_operations = false
}


resource "equinix_ecx_l2_connection" "azure-dot1q-pub" {
  name              = "tf-azure-dot1q-pub-pri"
  profile_uuid      = data.equinix_ecx_l2_sellerprofile.azure.uuid
  speed             = azurerm_express_route_circuit.demo.bandwidth_in_mbps
  speed_unit        = "MB"
  notifications     = ["example@equinix.com"]
  port_uuid         = data.equinix_ecx_port.dot1q-1-pri.uuid
  vlan_stag         = 1010
  seller_metro_code = var.azure_metro_code
  authorization_key = azurerm_express_route_circuit.demo.service_key
  named_tag         = "Public"
  secondary_connection {
    name      = "tf-azure-dot1q-pub-sec"
    port_uuid = data.equinix_ecx_port.dot1q-1-sec.uuid
    vlan_stag = 1300
  }
}
