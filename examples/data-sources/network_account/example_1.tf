# Retrieve details of an account in Active status in DC metro
data "equinix_network_account" "dc" {
  metro_code = "DC"
  status     = "Active"
  project_id = "a86d7112-d740-4758-9c9c-31e66373746b" 
}

output "number" {
  value = data.equinix_network_account.dc.number
}
