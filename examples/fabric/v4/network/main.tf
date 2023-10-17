provider "equinix" {
  client_id     = var.equinix_client_id
  client_secret = var.equinix_client_secret
}

resource "equinix_fabric_network" "test"{
  name = var.network_name
  type = var.network_type
  scope = var.network_scope

  notifications{
    type=var.notifications_type
    emails=var.notifications_emails
  }
  order {
    purchase_order_number= var.purchase_order_number
  }
}

output "fcr_result" {
  value = equinix_fabric_network.test.id
}

