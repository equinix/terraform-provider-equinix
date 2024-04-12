# Creates ACL template and assigns it to the network device
resource "equinix_network_acl_template" "myacl" {
  name        = "test"
  description = "Test ACL template"
  project_id = "a86d7112-d740-4758-9c9c-31e66373746b"
  inbound_rule {
    subnet  = "1.1.1.1/32"
    protocol = "IP"
    src_port = "any"
    dst_port = "any"
    description = "inbound rule description"
  }
  inbound_rule {
    subnet  = "172.16.25.0/24"
    protocol = "UDP"
    src_port = "any"
    dst_port = "53,1045,2041"
  }
}
