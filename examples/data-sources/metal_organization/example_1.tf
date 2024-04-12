# Fetch a organization data and show projects which belong to it
data "equinix_metal_organization" "test" {
  organization_id = local.org_id
}

output "projects_in_the_org" {
  value = data.equinix_metal_organization.test.project_ids
}
