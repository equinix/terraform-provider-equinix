resource "equinix_fabric_precision_time" "ntp" {
  type = "NTP"
  name = "tf_ntp_PFCR"
  description = "Equinix Precision Time with NTP Configuration"
  package {
    code = "NTP_STANDARD"
  }
  connections {
    uuid = "30b82c65-ffb4-47d3-ab2b-3cacf46d5b8b"
  }
  ipv4 {
    primary = "192.168.0.2"
    secondary = "192.168.0.3"
    network_mask = "255.255.255.224"
    default_gateway = "192.168.0.1"
  }
}
