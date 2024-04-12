resource "equinix_metal_device" "pxe1" {
  hostname         = "tf.coreos2-pxe"
  plan             = "c3.small.x86"
  metro            = "sv"
  operating_system = "custom_ipxe"
  billing_cycle    = "hourly"
  project_id       = local.project_id
  ipxe_script_url  = "https://rawgit.com/cloudnativelabs/pxe/master/packet/coreos-stable-metal.ipxe"
  always_pxe       = "false"
  user_data        = data.ignition_config.example.rendered
}
