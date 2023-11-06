provider "equinix" {
  client_id     = var.equinix_client_id
  client_secret = var.equinix_client_secret
}

data "equinix_fabric_service_profiles" "ibm" {
  filter {
    property = "/name"
    operator = "="
    values   = [var.fabric_sp_name]
  }
}

data "equinix_fabric_ports" "port" {
  filters {
    name = var.equinix_port_name
  }
}

resource "equinix_fabric_connection" "ibm1" {
  name = var.connection_name
  type = var.connection_type

  notifications {
    type   = var.notifications_type
    emails = var.notifications_emails
  }

  bandwidth = var.bandwidth

  additional_info  = [{ key = "ASN", value = var.seller_asn }]

  redundancy {
    priority  = var.redundancy
  }
  order {
    purchase_order_number = var.purchase_order_number
  }
  a_side {
    access_point {
      type = var.aside_ap_type
      port {
        uuid = data.equinix_fabric_ports.port.id
      }
      link_protocol {
        type     = var.aside_link_protocol_type
        vlan_tag = var.aside_link_protocol_tag
      }
    }
  }
  z_side {
    access_point {
      type                = var.zside_ap_type
      authentication_key  = var.zside_ap_authentication_key
      seller_region       = var.seller_region
      profile {
        type  = var.zside_ap_profile_type
        uuid  = data.equinix_fabric_service_profiles.ibm.id
      }
      location {
        metro_code  = var.zside_location
      }
    }
  }
}

output "connection_result" {
  value = equinix_fabric_connection.ibm1.id
}
