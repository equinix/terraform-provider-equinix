resource "equinix_metal_organization_member" "owner" {
    invitee = "admin@example.com"
    roles = ["owner"]
    projects_ids = []
    organization_id = var.organization_id
}
