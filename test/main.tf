terraform {
  required_providers {
    equinix = {
      source = "equinix/equinix"
      version = "1.25.1"
    }
  }
}

# Credentials for all Equinix resources
provider "equinix" {}

# Get Project by name and print UUIDs of its users
data "equinix_metal_project" "tf_project_1" {
  name = "ocm-test"
}

resource "equinix_metal_project_ssh_key" "foobar" {
	name = "ocm-test"
	public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDM/unxJeFqxsTJcu6mhqsMHSaVlpu+Jj/P+44zrm6X/MAoHSX3X9oLgujEjjZ74yLfdfe0bJrbL2YgJzNaEkIQQ1VPMHB5EhTKUBGnzlPP0hHTnxsjAm9qDHgUPgvgFDQSAMzdJRJ0Cexo16Ph9VxCoLh3dxiE7s2gaM2FdVg7P8aSxKypsxAhYV3D0AwqzoOyT6WWhBoQ0xZ85XevOTnJCpImSemEGs6nVGEsWcEc1d1YvdxFjAK4SdsKUMkj4Dsy/leKsdi/DEAf356vbMT1UHsXXvy5TlHu/Pa6qF53v32Enz+nhKy7/8W2Yt2yWx8HnQcT2rug9lvCXagJO6oauqRTO77C4QZn13ZLMZgLT66S/tNh2EX0gi6vmIs5dth8uF+K6nxIyKJXbcA4ASg7F1OJrHKFZdTc5v1cPeq6PcbqGgc+8SrPYQmzvQqLoMBuxyos2hUkYOmw3aeWJj9nFa8Wu5WaN89mUeOqSkU4S5cgUzWUOmKey56B/j/s1sVys9rMhZapVs0wL4L9GBBM48N5jAQZnnpo85A8KsZq5ME22bTLqnxsDXqDYZvS7PSI6Dxi7eleOFE/NYYDkrgDLHTQri8ucDMVeVWHgoMY2bPXdn7KKy5jW5jKsf8EPARXg77A4gRYmgKrcwIKqJEUPqyxJBe0CPoGTqgXPRsUiQ== tomk@hp2"
	project_id = data.equinix_metal_project.tf_project_1.id
}

data "equinix_metal_project_ssh_key" "foobar" {
	search = equinix_metal_project_ssh_key.foobar.fingerprint
	project_id = data.equinix_metal_project.tf_project_1.id
}

output "users_of_Terraform_Fun" {
  value = data.equinix_metal_project_ssh_key.foobar.fingerprint
}
