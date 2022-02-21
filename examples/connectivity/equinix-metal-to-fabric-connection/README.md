# ECX Fabric Layer2 Connection from Equinix Metal to Service Provider

This example shows how to create a layer 2 connection between Equinix Metal and a service provider, using an Equinix Metal a-side service token

## Adjust variables

At minimum, you must override below variables in `terraform.tfvars` file:

* `equinix_client_id` - Equinix client ID (consumer key), obtained after registering app in the developer platform
* `equinix_client_secret` - Equinix client secret ID (consumer secret), obtained same way as above
* `metal_auth_token` - This is your Equinix Metal API Auth token
* `metal_project_name` - Name of an existing metal project
* `seller_profile_name`- Name of the service provider to connect with, i.e. 'AWS Direct Connect'
* `seller_region`- The region code in which the seller port resides, i.e. 'us-west-1'
* `seller_authorization_key`- Text field used to authorize connection on the provider side. Value depends on a provider service profile used for connection
* `connection_notification_users` - List of email addresses used for sending connection update notifications

## Initialize

Change directory to example directory and initialize Terraform plugins
by running `terraform init`.

## Deploy template

Apply changes by running `terraform apply`, then **inspect proposed plan**
and approve it.
