resource "equinix_metal_device" "web1" {
  hostname         = "tf.coreos2"
  plan             = "c3.small.x86"
  metro            = "ny"
  operating_system = "ubuntu_20_04"
  billing_cycle    = "hourly"
  project_id       = local.project_id
  ip_address {
    type = "private_ipv4"
    cidr = 30
  }
}
