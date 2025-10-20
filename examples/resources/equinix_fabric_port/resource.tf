resource "equinix_fabric_port" "order" {
  type = "XF_PORT"
  connectivity_source_type = "COLO"
  location = {
    metro_code = "TR"
  }
  settings = {
    package_type = "STANDARD"
    shared_port_type = false
  }
  encapsulation = {
    type = "DOT1Q"
    tag_protocol_id = "0x8100"
  }
  account = {
    account_number = "<account_number>"
  }
  project = {
    project_id = "<project_id>"
  }
  redundancy = {
    priority = "PRIMARY"
  }
  lag_enabled = true
  physical_ports = [
    {
      type = "XF_PHYSICAL_PORT"
      demarcation_point = {
        ibx = "TR2"
        cage_unique_space_id = "TR2:01:002087"
        cabinet_unique_space_id = "Demarc"
        patch_panel = "PP:Demarc:00002087"
        connector_type = "SC"
      }
    }
  ]
  physical_ports_speed = 1000
  physical_ports_type = "1000BASE_LX"
  physical_ports_count = 1
  demarcation_point_ibx = "TR2"
  notifications = [
    {
      type = "TECHNICAL"
      registered_users = [
        "<username>"
      ]
    },
    {
      type = "NOTIFICATION"
      registered_users = [
        "<username>"
      ]
    }
  ]
  additional_info = [
    {
      key = "lagType"
      value = "New"
    }
  ]
}
