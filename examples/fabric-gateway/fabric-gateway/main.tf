provider "equinix" {
  client_id     = var.equinix_client_id
  client_secret = var.equinix_client_secret
}

resource "equinix_fabric_gateway" "test"{
  name = var.fg_name
  type = var.fg_type
  notifications{
    type=var.notifications_type
    emails=var.notifications_emails
  }
  order {
    purchase_order_number= var.purchase_order_number
  }
  location {
    metro_code=var.fg_location
  }
  package {
    code=var.fg_package
  }
  project {
  	project_id = var.fg_project
  }
  account {
  	account_number = var.fg_account
  }
}

output "fg_result" {
  value = equinix_fabric_gateway.test.id
}

