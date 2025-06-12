# Configuration for using Workload Identity Federation
provider "equinix" {
  # Desired scope of the requested security token. Must be an Access Policy ERN or a string of the form `roleassignments:<organization_id>`
  auth_scope = "roleassignments:<organization_id>"

  # The TFC_WORKLOAD_IDENTITY_TOKEN env variable must contain a workload identity token
}