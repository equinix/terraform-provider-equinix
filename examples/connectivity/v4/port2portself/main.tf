provider "equinix" {
  client_id     = var.equinix_client_id
  client_secret = var.equinix_client_secret
}

data "equinix_fabric_ports" "aside" {
  filters {
    name = var.aside_port_name
  }
}

data "equinix_fabric_ports" "zside" {
  filters {
    name = var.zside_port_name
  }
}
resource "equinix_fabric_connection" "p2p-qinq" {
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
      port {
        uuid= data.equinix_fabric_ports.aside.data.0.uuid
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
      port {
        uuid= data.equinix_fabric_ports.zside.data.0.uuid
      }
      link_protocol {
        type= var.zside_link_protocol_type
        vlan_s_tag= var.zside_link_protocol_stag
      }
      location {
        metro_code= var.zside_location
      }
    }
  }
}

output "connection_result" {
  value = equinix_fabric_connection.p2p-qinq.id
}
