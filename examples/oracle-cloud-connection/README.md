# ECX Fabric Layer2 Connection to Oracle Cloud

This example shows how create layer 2 connection between ECX Fabric port
and Oracle Cloud FastConnect.
Example covers **provisioning of both sides** of the connection.

## Adjust variables

At minimum, you must set below variables in `terraform.tfvars` file:

* `equinix_client_id` - Equinix client ID (consumer key), obtained after
registering app in the developer platform
* `equinix_client_secret` - Equinix client secret ID (consumer secret),
obtained same way as above
* `equinix_port_name`     - name of ECX Fabric port that you want to connect
to Oracle i.e. *EQUINIX_SVC-FR4-CX-PRI-01*
* `oci_tenancy_ocid` - tenancy's Oracle Cloud Identifier
* `oci_user_ocid` - user's Oracle Cloud Identifier
* `oci_private_key_path` - API Singing private key
* `oci_fingerprint` - API signing private key's fingerprint
* `oci_region` - Oracle Cloud region to connect to
* `oci_compartment_id` - Compartment's Oracle Cloud Identifier

## Initialize

Change directory to example directory and initialize Terraform plugins
by running `terraform init`.

## Deploy template

Apply changes by running `terraform apply`, then **inspect proposed plan**
and approve it.
