provider "equinix" {
  client_id     = var.equinix_client_id
  client_secret = var.equinix_client_secret
}

data "equinix_fabric_service_profiles" "sp" {
  filter {
    property = "/name"
    operator = "="
    values   = [var.fabric_sp_name]
  }
}

resource "equinix_fabric_connection" "fcr2profile" {
  name = var.connection_name
  type = var.connection_type
  notifications{
    type   = var.notifications_type
    emails = var.notifications_emails
  }
  bandwidth = var.bandwidth
  redundancy {priority = var.redundancy}
  order {
    purchase_order_number = var.purchase_order_number
  }
  a_side {
    access_point {
      type = var.aside_ap_type
      router {
        uuid = var.fcr_uuid
      }
    }
  }
  z_side {
    access_point {
      type = var.zside_ap_type
      profile {
        type = var.zside_ap_profile_type
        uuid = data.equinix_fabric_service_profiles.sp.id
      }
      location {
        metro_code = var.zside_location
      }
    }
  }
}

output "connection_result" {
  value = equinix_fabric_connection.fcr2profile.id
}
