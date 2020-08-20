provider "equinix" {
  client_id     = "your_client_id"
  client_secret = "your_client_secret"
}

provider "azurerm" {
  version = "=2.22.0"
}

data "equinix_ecx_l2_sellerprofile" "azure" {
  name = "Azure Express Route"
}

data "equinix_ecx_port" "dot1q-1-pri" {
  name = "sit-001-CX-DC5-NL-Dot1q-BO-10G-PRI-JUN-27"
}

data "equinix_ecx_port" "dot1q-1-sec" {
  name = "sit-001-CX-DC6-NL-Dot1q-BO-10G-SEC-JUN-28"
}

resource "azurerm_resource_group" "demo" {
  name     = "TFDemo"
  location = "West Europe"
}

resource "azurerm_express_route_circuit" "demo" {
  name                  = "TFDemoExpressRoute"
  resource_group_name   = azurerm_resource_group.demo.name
  location              = azurerm_resource_group.demo.location
  service_provider_name = "Equinix Test"
  peering_location      = "Area51"
  bandwidth_in_mbps     = 50
  sku {
    tier   = "Premium"
    family = "UnlimitedData"
  }
  allow_classic_operations = false
}

resource "azurerm_express_route_circuit_authorization" "demo" {
  name                       = "TFDemoExpressRouteAuth"
  express_route_circuit_name = azurerm_express_route_circuit.demo.name
  resource_group_name        = azurerm_resource_group.demo.name
}

resource "equinix_ecx_l2_connection" "azure-dot1q-pub" {
  name                  = "tf-azure-dot1q-pub-pri"
  profile_uuid          = data.equinix_ecx_l2_sellerprofile.azure.uuid
  speed                 = azurerm_express_route_circuit.demo.bandwidth_in_mbps
  speed_unit            = "MB"
  notifications         = ["kkolla@equinix.com"]
  purchase_order_number = "1234567890"
  port_uuid             = data.equinix_ecx_port.pri-dot1q.uuid
  vlan_stag             = 101
  seller_metro_code     = "DA"
  authorization_key     = azurerm_express_route_circuit_authorization.demo.authorization_key
  named_tag             = "Public"
  secondary_connection {
    name      = "tf-azure-dot1q-pub-sec"
    port_uuid = data.equinix_ecx_port.sec-dot1q.uuid
    vlan_stag = 130
  }
}
