# Configuration for using Workload Identity Federation
provider "equinix" {
  # Desired scope of the requested security token. Must be an Access Policy ERN or a string of the form `roleassignments:<organization_id>`
  sts_auth_scope = "roleassignments:<organization_id>"

  # The name of the environment variable containing the STS source token
  # For example, HCP Terraform automatically sets TFC_WORKLOAD_IDENTITY_TOKEN
  sts_source_token_env_var = "TFC_WORKLOAD_IDENTITY_TOKEN"
}