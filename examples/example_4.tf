# Configuration for using Workload Identity Federation
provider "equinix" {
  # Desired scope of the requested security token. Must be an Access Policy ERN or a string of the form `roleassignments:<organization_id>`
  auth_scope = "roleassignments:<organization_id>"

  # An OIDC ID token issued by a trusted OIDC provider to a trusted client.
  sts_source_token = "some_workload_identity_token"
}