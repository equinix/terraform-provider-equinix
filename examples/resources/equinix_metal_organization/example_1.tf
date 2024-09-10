# Create a new Organization
resource "equinix_metal_organization" "tf_organization_1" {
  name        = "foobar"
  description = "quux"
}
