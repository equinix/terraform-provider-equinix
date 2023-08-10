provider "equinix" {
  client_id     = var.equinix_client_id
  client_secret = var.equinix_client_secret
}
data "equinix_fabric_ports" "zside" {
  filters {
    name = var.zside_port_name
  }
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

resource "equinix_fabric_connection" "fg2port"{
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
      gateway {
        uuid= equinix_fabric_gateway.test.id
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
        vlan_tag= var.zside_link_protocol_tag
      }
      location {
        metro_code= var.zside_location
      }
    }
  }
}

output "connection_result" {
  value = equinix_fabric_connection.fg2port.id
}

