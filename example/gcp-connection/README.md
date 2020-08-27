# ECX Fabric Layer2 Connection to Google Cloud Platform

This example shows how create layer 2 connection between ECX Fabric port
and Google Cloud Partner Interconnect.
Example covers **provisioning of both sides** of the connection.

## Adjust variables

At minimum, you must set below variables in `terraform.tfvars` file:

* `equinix_client_id` - Equinix client ID (consumer key), obtained after
registering app in the developer platform
* `equinix_client_secret` - Equinix client secret ID (consumer secret),
obtained same way as above
* `equinix_port_name` - name of ECX Fabric port that you want to connect
to GCP i.e. *EQUINIX_SVC-FR4-CX-PRI-01*
* `gcp_region` -  GCP region for connection
* `gcp_project_name` - GCP project name

## GCP login

Log in to GCP with an account that has permission to create the necessary
resources using `gcloud init`.

**NOTE** for purpose of creating partner interconnect, this template creates
network and router as well. Please check `main.tf` for details.

## Initialize

Change directory to example directory and initialize Terraform plugins
by running `terraform init`.

## Deploy template

Apply changes by running `terraform apply`, then **inspect proposed plan**
and approve it.
