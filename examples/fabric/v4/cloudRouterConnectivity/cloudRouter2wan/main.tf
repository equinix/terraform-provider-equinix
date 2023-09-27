provider "equinix" {
  client_id     = var.equinix_client_id
  client_secret = var.equinix_client_secret
}


resource "equinix_fabric_connection" "fcr2ipwan"{
  name = var.connection_name
  type = var.connection_type
  notifications{
    type=var.notifications_type
    emails=var.notifications_emails
  }
  bandwidth = var.bandwidth

  order {
    purchase_order_number= var.purchase_order_number
  }
  a_side {
    access_point {
      type= var.aside_ap_type
      router {
        uuid= var.fcr_uuid
      }
    }
  }
 z_side {
    access_point {
      type = var.zside_ap_type
      network {
        uuid = var.network_uuid
      }
    }
}
}

output "connection_result" {
  value = equinix_fabric_connection.fcr2ipwan.id
}

