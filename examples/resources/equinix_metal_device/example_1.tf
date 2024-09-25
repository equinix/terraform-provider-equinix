resource "equinix_metal_device" "web1" {
  hostname         = "tf.coreos2"
  plan             = "c3.small.x86"
  metro            = "sv"
  operating_system = "ubuntu_20_04"
  billing_cycle    = "hourly"
  project_id       = local.project_id
}
