# Create a new Project
resource "equinix_metal_project" "tf_project_1" {
  name = "tftest"
  bgp_config {
    deployment_type = "local"
    md5             = "C179c28c41a85b"
    asn             = 65000
  }
}
