resource "equinix_fabric_cloud_router" "new_cloud_router"{
  name = "Router-SV"
  type = "XF_ROUTER"
  notifications{
    type = "ALL"
    emails = ["example@equinix.com","test1@equinix.com"]
  }
  order {
    purchase_order_number = "1-323292"
  }
  location {
    metro_code = "SV"
  }
  package {
    code = "STANDARD"
  }
  project {
  	project_id = "776847000642406"
  }
  account {
  	account_number = "203612"
  }
}
