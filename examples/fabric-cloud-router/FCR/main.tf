provider "equinix" {
  client_id     = var.equinix_client_id
  client_secret = var.equinix_client_secret
}

resource "equinix_fabric_cloud_router" "test"{
  name = var.fcr_name
  type = var.fcr_type
  notifications{
    type=var.notifications_type
    emails=var.notifications_emails
  }
  order {
    purchase_order_number= var.purchase_order_number
  }
  location {
    metro_code=var.fcr_location
  }
  package {
    code=var.fcr_package
  }
  project {
  	project_id = var.fcr_project
  }
  account {
  	account_number = var.fcr_account
  }
}

output "fcr_result" {
  value = equinix_fabric_cloud_router.test.id
}

