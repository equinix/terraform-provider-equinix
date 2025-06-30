resource "equinix_fabric_network" "new_network" {
  name  = "Network-SV"
  type  = "EVPLAN"
  scope = "GLOBAL"
  notifications {
    type   = "ALL"
    emails = ["example@equinix.com","test1@equinix.com"]
  }
  project {
    project_id = "776847000642406"
  }
}
