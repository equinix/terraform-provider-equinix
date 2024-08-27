resource "equinix_metal_organization_member" "member" {
    invitee = "member@example.com"
    roles = ["limited_collaborator"]
    projects_ids = [var.project_id]
    organization_id = var.organization_id
}
