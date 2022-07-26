# ECX Fabric Layer2 Connection to Google Cloud Platform

**NOTE:** There is an
[Equinix Fabric L2 Connection To Google Cloud Interconnect Terraform module](https://registry.terraform.io/modules/equinix-labs/fabric-connection-gcp/equinix/latest)
available with full-fledged examples of connections from Fabric Ports, Network Edge Devices
or Service Tokens.

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
* `gcp_project_name` - GCP project name
* `gcp_region` -  GCP region for a connection
* `gcp_metro_code` -  Equinix metro location with GCP presence for connection's destination.
Given metro location has to correspond with given GCP region

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
