resource "equinix_metal_project" "existing_project" {
  name = "The name of the project (if different, will rewrite)"
  bgp_config {
    deployment_type = "local"
    md5             = "C179c28c41a85b"
    asn             = 65000
  }
}
