## 1.5.0 (UNRELEASED)

FEATURES:

- **New Guide:** Migrating From the Packet Provider
- **New Guide:** Upgrading Devices from Facilities to Metros

- **New Provider Argument** `auth_token`
- **New Provider Argument** `max_retries`
- **New Provider Argument** `max_retry_wait_seconds`
- **New Resource:** `equinix_metal_port`
- **New Resource:** `equinix_metal_user_api_key`
- **New Resource:** `equinix_metal_project_api_key`
- **New Resource:** `equinix_metal_gateway`
- **New Resource:** `equinix_metal_connection`
- **New Resource:** `equinix_metal_virtual_circuit`
- **New Resource** `equinix_metal_bgp_session`
- **New Resource** `equinix_metal_device`
- **New Resource** `equinix_metal_device_network_type`
- **New Resource** `equinix_metal_ip_attachment`
- **New Resource** `equinix_metal_organization`
- **New Resource** `equinix_metal_port_vlan_attachment`
- **New Resource** `equinix_metal_project`
- **New Resource** `equinix_metal_project`
- **New Resource** `equinix_metal_project_ssh_key`
- **New Resource** `equinix_metal_reserved_ip_block`
- **New Resource** `equinix_metal_spot_market_request`
- **New Resource** `equinix_metal_ssh_key`
- **New Resource** `equinix_metal_vlan`
- **New Data Resource:** `equinix_metal_vlan`
- **New Data Resource:** `equinix_metal_reserved_ip_block`
- **New Data Resource:** `equinix_metal_port`
- **New Data Resource:** `equinix_metal_hardware_reservation`
- **New Data Resource:** `equinix_metal_metro`
- **New Data Resource:** `equinix_metal_connection`
- **New Data Resource** `equinix_metal_facility`
- **New Data Resource** `equinix_metal_device`
- **New Data Resource** `equinix_metal_device_bgp_neighbors`
- **New Data Resource** `equinix_metal_ip_block_ranges`
- **New Data Resource** `equinix_metal_operating_system`
- **New Data Resource** `equinix_metal_organization`
- **New Data Resource** `equinix_metal_precreated_ip_block`
- **New Data Resource** `equinix_metal_project`
- **New Data Resource** `equinix_metal_project_ssh_key`
- **New Data Resource** `equinix_metal_spot_market_price`
- **New Data Resource** `equinix_metal_spot_market_request`

BUG FIXES:

* `equinix_ecx_l2_serviceprofile` detecting diff after refresh [#90](https://github.com/equinix/terraform-provider-equinix/pull/90)
* `equinix_network_device` allow value 0 for additional bandwidth [#91](https://github.com/equinix/terraform-provider-equinix/pull/91)
* `equinix_ecx_l2_connection` will wait for the secondary connection destroy [#103](https://github.com/equinix/terraform-provider-equinix/pull/103)
* `equinix_ecx_l2_connection` named_tag now is idempotent [#97](https://github.com/equinix/terraform-provider-equinix/issues/97)

ENHANCEMENTS:

- `service_token` added to `equinix_ecx_l2_connection` [#96](https://github.com/equinix/terraform-provider-equinix/issues/96)
- update documentation links for timeout parameters  [#101](https://github.com/equinix/terraform-provider-equinix/pull/101)

## 1.4.0 (January 14, 2022)

NOTES:

* `equinix_acl_template` argument `metro_code` is now deprecated [#67](https://github.com/equinix/terraform-provider-equinix/pull/67)
* `equinix_acl_template` argument `inbound_rule.#.subnets` is now deprecated [#67](https://github.com/equinix/terraform-provider-equinix/pull/67)
* `equinix_acl_template` attribute `device_id` is now deprecated [#67](https://github.com/equinix/terraform-provider-equinix/pull/67)
* `equinix_ecx_l2_connection_accepter` is now deprecated [#64](https://github.com/equinix/terraform-provider-equinix/pull/64)
* `equinix_network_device_link` argument `device.interface_id` changes taint the resource [#77](https://github.com/equinix/terraform-provider-equinix/pull/77)
* `equinix_network_device_link` attribute `link.src_zone_code` is now deprecated and optional [#77](https://github.com/equinix/terraform-provider-equinix/pull/77)
* `equinix_network_device_link` attribute `link.dest_zone_code` is now deprecated and optional [#77](https://github.com/equinix/terraform-provider-equinix/pull/77)

BUG FIXES:

* `equinix_ecx_l2_connection` will wait for the secondary connection [#87](https://github.com/equinix/terraform-provider-equinix/pull/87)
* `equinix_network_device_link` no longer jitters on zone code fields [#77](https://github.com/equinix/terraform-provider-equinix/pull/77)

ENHANCEMENTS:

* `equinix_acl_template` attribute `device_details` (`uuid`, `name`, `acl_status`) was added [#67](https://github.com/equinix/terraform-provider-equinix/pull/67)
* `equinix_ecx_l2_connection` attribute `actions` was added [#86](https://github.com/equinix/terraform-provider-equinix/pull/86)
* `equinix_network_device_link` argument `device.asn` is now optional [#77](https://github.com/equinix/terraform-provider-equinix/pull/77)
* `equinix_network_device_link` argument `device.subnet` is now optional [#77](https://github.com/equinix/terraform-provider-equinix/pull/77)
* fix connectivity example to establish an Azure connection [#71](https://github.com/equinix/terraform-provider-equinix/pull/71)
* replace Travis CI with GitHub Workflows [#65](https://github.com/equinix/terraform-provider-equinix/pull/65)
* update go modules and update go to 1.17 [#63](https://github.com/equinix/terraform-provider-equinix/pull/63)

## 1.3.0 (November 18, 2021)

BUG FIXES:

- `equinix_network_device` no longer loses `license_token` values after updating the resource ([#59](https://github.com/equinix/terraform-provider-equinix/issues/59))

ENHANCEMENTS:

- `wan_interface_id` added to `equinix_network_device` ([#59](https://github.com/equinix/terraform-provider-equinix/issues/59))
- `equinix_ecx_l2_connection` resources can now be imported ([#49](https://github.com/equinix/terraform-provider-equinix/issues/49))
- `equinix_ecx_l2_connection_accepter` resources can now be imported ([#49](https://github.com/equinix/terraform-provider-equinix/issues/49))
- `equinix_ecx_l2_serviceprofile` resources can now be imported ([#50](https://github.com/equinix/terraform-provider-equinix/issues/50))
- `equinix_network_acl_template` resources can now be imported ([#50](https://github.com/equinix/terraform-provider-equinix/issues/50))
- `equinix_network_bgp` resources can now be imported ([#50](https://github.com/equinix/terraform-provider-equinix/issues/50))
- `equinix_network_device_link` resources can now be imported ([#50](https://github.com/equinix/terraform-provider-equinix/issues/50))
- `equinix_network_ssh_key` resources can now be imported ([#50](https://github.com/equinix/terraform-provider-equinix/issues/50))
- `equinix_network_ssh_user` resources can now be imported ([#50](https://github.com/equinix/terraform-provider-equinix/issues/50))
- darwin/arm64 binaries are now published ([#51](https://github.com/equinix/terraform-provider-equinix/issues/51))
- Go version used for builds was updated to 1.17 ([#52](https://github.com/equinix/terraform-provider-equinix/issues/52))

## 1.2.0 (April 27, 2021)

FEATURES:

- **New Resource**: `equinix_network_device_link` ([#43](https://github.com/equinix/terraform-provider-equinix/issues/43))

## 1.1.0 (April 09, 2021)

BUG FIXES:

- creation of Equinix Fabric layer2 redundant connection from a single device
is now possible by specifying same `deviceUUID` argument for both primary and
secondary connection. API logic of Fabric is reflected accordingly in client module

FEATURES:

- **New Data source**: `equinix_ecx_l2_sellerprofiles`: ([#40](https://github.com/equinix/terraform-provider-equinix/issues/40))
- **New Resource**: `equinix_network_ssh_key` ([#25](https://github.com/equinix/terraform-provider-equinix/issues/25))
- **New Resource**: `equinix_network_acl_template` ([#19](https://github.com/equinix/terraform-provider-equinix/issues/19))
- **New Resource**: `equinix_network_bgp` ([#16](https://github.com/equinix/terraform-provider-equinix/issues/16))
- **New Data source**: `equinix_network_account` ([#13](https://github.com/equinix/terraform-provider-equinix/issues/13))
- **New Data source**: `equinix_network_device_type` ([#13](https://github.com/equinix/terraform-provider-equinix/issues/13))
- **New Data source**: `equinix_network_device_software` ([#13](https://github.com/equinix/terraform-provider-equinix/issues/13))
- **New Data source**: `equinix_network_device_platform` ([#13](https://github.com/equinix/terraform-provider-equinix/issues/13))
- **New Resource**: `equinix_network_device` ([#4](https://github.com/equinix/terraform-provider-equinix/issues/4))
- **New Resource**: `equinix_network_ssh_user` ([#4](https://github.com/equinix/terraform-provider-equinix/issues/4))

ENHANCEMENTS:

- Equinix provider: setting `TF_LOG` to `TRACE` enables logging of Equinix REST
API requests and responses
- resource/equinix_ecx_l2_connection: internal representation of secondary connection
block has changed from Set to List. This enables plan to better communicate secondary
connection changes and allows using `Optional` + `Computed` schema options
([#39](https://github.com/equinix/terraform-provider-equinix/issues/39))
- resource/equinix_ecx_l2_connection: added additional arguments for `secondary_connection`
([#18](https://github.com/equinix/terraform-provider-equinix/issues/18)):
  - `speed`
  - `speed_unit`
  - `profile_uuid`
  - `authorization_key`
  - `seller_metro_code`
  - `seller_region`

## 1.0.3 (January 07, 2021)

ENHANCEMENTS:

- resource/equinix_ecx_l2_connection_accepter: AWS credentials can be provided
using additional ways: environmental variables and shared configuration files
- resource/equinix_ecx_l2_service_profile: introduced schema validations,
updated acceptance tests and resource documentation

BUG FIXES:

- resource/equinix_ecx_l2_connection_accepter: creation waits for PROVISIONED provider
status of the connection before succeeding
([#37](https://github.com/equinix/terraform-provider-equinix/issues/37))

## 1.0.2 (November 17, 2020)

ENHANCEMENTS:

- resource/equinix_ecx_l2_connection_accepter: creation awaits for desired
connection provider state before succeeding ([#26](https://github.com/equinix/terraform-provider-equinix/issues/26))

BUG FIXES:

- resource/equinix_ecx_l2_connection: z-side port identifier, vlan C-tag and vlan
S-tag for secondary connection are properly populated with values from the Fabric
([#24](https://github.com/equinix/terraform-provider-equinix/issues/24))

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
