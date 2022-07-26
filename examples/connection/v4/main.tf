provider "equinix" {
  client_id     = var.equinix_client_id
  client_secret = var.equinix_client_secret
}

provider "azurerm" {
  features {}
}

resource "equinix_ecx_l2_connection" "azure-dot1q-pub" {
  name = var.connection_name
  description = var.description
  type = var.connection_type
  notifications{
    type=var.notifications_type
    emails=var.notifications_emails
  }
  bandwidth = var.bandwidth
  redundancy {priority= var.redundancy_pri}
  order {
    purchase_order_number= var.purchase_order_number
  }
  a_side {
    access_point {
      type= var.aside_ap_type
      port {
        uuid= var.aside_port_uuid
      }
      link_protocol {
        type= var.aside_link_protocol_type
        vlan_s_tag= var.aside_link_protocol_stag
      }
    }
  }
  z_side {
    access_point {
      type= var.zside_ap_type
      authentication_key= var.zside_ap_authentication_key
      profile {
        type= var.zside_ap_profile_type
        uuid= var.zside_ap_profile_uuid
      }
      location {
        metro_code= var.zside_location
      }
    }
  }
}
