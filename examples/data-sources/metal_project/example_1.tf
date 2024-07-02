# Get Project by name and print UUIDs of its users
data "equinix_metal_project" "tf_project_1" {
  name = "Terraform Fun"
}

output "users_of_Terraform_Fun" {
  value = data.equinix_metal_project.tf_project_1.user_ids
}
