data "equinix_metal_operating_system" "example" {
  distro           = "ubuntu"
  version          = "20.04"
  provisionable_on = "c3.medium.x86"
}

resource "equinix_metal_device" "server" {
  hostname         = "tf.ubuntu"
  plan             = "c3.medium.x86"
  metro            = "ny"
  operating_system = data.equinix_metal_operating_system.example.id
  billing_cycle    = "hourly"
  project_id       = local.project_id
}
