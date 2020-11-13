## 1.1.0 (Unreleased)

ENHANCEMENTS:

- resource/equinix_ecx_l2_connection_accepter: creation awaits for desired
connection provider state before succeeding ([#26](https://github.com/equinix/terraform-provider-equinix/issues/26))

BUG FIXES:

- resource/equinix_ecx_l2_connection: z-side port identifier, vlan C-tag and vlan S-tag 
for secondary connection are properly populated with values from the Fabric ([#24](https://github.com/equinix/terraform-provider-equinix/issues/24))

FEATURES:

- **New Resource**: `equinix_network_acl_template` ([#19](https://github.com/equinix/terraform-provider-equinix/issues/19))
- **New Resource**: `equinix_network_bgp` ([#16](https://github.com/equinix/terraform-provider-equinix/issues/16))
- **New Data source**: `equinix_network_account` ([#13](https://github.com/equinix/terraform-provider-equinix/issues/13))
- **New Data source**: `equinix_network_device_type` ([#13](https://github.com/equinix/terraform-provider-equinix/issues/13))
- **New Data source**: `equinix_network_device_software` ([#13](https://github.com/equinix/terraform-provider-equinix/issues/13))
- **New Data source**: `equinix_network_device_platform` ([#13](https://github.com/equinix/terraform-provider-equinix/issues/13))
- **New Resource**: `equinix_network_device` ([#4](https://github.com/equinix/terraform-provider-equinix/issues/4))
- **New Resource**: `equinix_network_ssh_user` ([#4](https://github.com/equinix/terraform-provider-equinix/issues/4))

## 1.0.1 (November 09, 2020)

NOTES:

- this version of module started to use `equinix/rest-go` client
for any REST interactions with Equinix APIs

ENHANCEMENTS:

- resource/equinix_ecx_l2_connection_accepter: added `aws_connection_id` attribute
([#22](https://github.com/equinix/terraform-provider-equinix/issues/22))
- resource/equinix_ecx_l2_connection: removal awaits for desired
connection state before succeeding ([#21](https://github.com/equinix/terraform-provider-equinix/issues/21))
- resource/equinix_ecx_l2_connection: added `device_interface_id` argument ([#18](https://github.com/equinix/terraform-provider-equinix/issues/18))
- resource/equinix_ecx_l2_connection: added `provider_status` and
 `redundancy_type` attributes ([#14](https://github.com/equinix/terraform-provider-equinix/issues/14))
- resource/equinix_ecx_l2_connection: creation awaits for desired
connection state before succeeding ([#15](https://github.com/equinix/terraform-provider-equinix/issues/15))

## 1.0.0 (September 02, 2020)

NOTES:

- first version of official Equinix Terraform provider

FEATURES:

- **New Resource**: `equinix_ecx_l2_connection`
- **New Resource**: `equinix_ecx_l2_connection_accepter`
- **New Resource**: `equinix_ecx_l2_serviceprofile`
- **New Data Source**: `equinix_ecx_port`
- **New Data Source**: `equinix_ecx_l2_sellerprofile`
