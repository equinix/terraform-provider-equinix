## 1.10.0 (Sep 29, 2022)

BUG FIXES:

- Allow to not define service_token_type for shared connections in organizations where service token is not enabled [#251](https://github.com/equinix/terraform-provider-equinix/pull/251)
- Fix some documentation typos with wrong information [#250](https://github.com/equinix/terraform-provider-equinix/pull/250)
- `equinix_metal_precreated_ip_block` was not filtering by project_id [#249](https://github.com/equinix/terraform-provider-equinix/pull/249)

ENHANCEMENTS:

- New resource `equinix_metal_organization_member` [#256](https://github.com/equinix/terraform-provider-equinix/pull/256)

## 1.9.0 (Sep 4, 2022)

BUG FIXES:

- `equinix_metal_device` `reinstall` action was not taken new value for `operating_system` [#247](https://github.com/equinix/terraform-provider-equinix/pull/247)

ENHANCEMENTS:

- Packngo version bumped to 0.26.0

## 1.8.1 (Aug 19, 2022)

BUG FIXES:

- `equinix_metal_device` `operating_system` should only be ForceNew when `reinstall` is not enabled [#244](https://github.com/equinix/terraform-provider-equinix/pull/244)
- clarified use of dedicated connection examples in `equinix_metal_virtual_circuit` docs [#242](https://github.com/equinix/terraform-provider-equinix/pull/242)

## 1.8.0 (Aug 1, 2022)

ENHANCEMENTS:

- `equinix_network_acl_template` added description field for the acl template inbound rule [#236](https://github.com/equinix/terraform-provider-equinix/pull/236)
- Update go version to 1.18.3 in main module and github actions workflows [#219](https://github.com/equinix/terraform-provider-equinix/pull/219)
- Fix missing equinix_metal_ name prefix in `equinix_metal_hardware_reservation` datasource documentation example [#231](https://github.com/equinix/terraform-provider-equinix/pull/231)
- Adds check for forbidden API error in `equinix_metal_device` read function [#235](https://github.com/equinix/terraform-provider-equinix/pull/235)
- Improve service toknes docs and update connectivty examples to link to Fabric terraform modules [#228](https://github.com/equinix/terraform-provider-equinix/pull/228)

## 1.7.0 (Jul 15, 2022)

ENHANCEMENTS:

- `zside_service_token` argument added to `equinix_ecx_l2_connection` to create connections with z-side token [#224](https://github.com/equinix/terraform-provider-equinix/pull/224)
- `vendor_token` attribute added to `equinix_ecx_l2_connection` to populate used a-side/z-side token [#224](https://github.com/equinix/terraform-provider-equinix/pull/224)

## 1.6.1 (Jul 7, 2022)

BUG FIXES:

- Fix client type assertion in equinix_metal_reserved_ip_block [#220](https://github.com/equinix/terraform-provider-equinix/pull/220)

## 1.6.0 (Jul 6, 2022)

FEATURES:

- New data source `equinix_metal_plans` for querying plans using filters [#215](https://github.com/equinix/terraform-provider-equinix/pull/215)
- New resource and data source `equinix_metal_vrf` [#129](https://github.com/equinix/terraform-provider-equinix/pull/129)
- Adds `address` as a datasource attribute and required resource argument to `equinix_metal_organization` [#137](https://github.com/equinix/terraform-provider-equinix/pull/137)
- Adds `vrf` as a datasource attribute to `equinix_metal_gateway` [#129](https://github.com/equinix/terraform-provider-equinix/pull/129)
- Adds `vrf_id` as a datasource attribute to `equinix_metal_gateway` [#129](https://github.com/equinix/terraform-provider-equinix/pull/129)
- Adds `vrf_id`, `peer_asn`, `subnet`, `metal_ip`, `customer_ip`, `md5` as resource arguments and datasource attributes to `equinix_metal_virtual_circuit` [#129](https://github.com/equinix/terraform-provider-equinix/pull/129))
- Adds `vrf_id`, `network`, `cidr` as resource arguments to `equinix_metal_reserved_ip_block` [#129](https://github.com/equinix/terraform-provider-equinix/pull/129)
- Adds `user_ssh_key_ids` as resource argument to `metal_device` [#141](https://github.com/equinix/terraform-provider-equinix/pull/141)

BUG FIXES:

- Change `equinix_network_acl_template` docs subcategory to network edge [#128](https://github.com/equinix/terraform-provider-equinix/pull/128)
- `equinix_network_device` removed hostname validation and fix acl issues in device deletion flow [#126](https://github.com/equinix/terraform-provider-equinix/pull/126)
- Fix provider required credentials [#125](https://github.com/equinix/terraform-provider-equinix/pull/125)
- migration-tool: remove duplicate readme.md [#153](https://github.com/equinix/terraform-provider-equinix/pull/153)

ENHANCEMENTS:

- `mgmt_acl_template_uuid` argument added to `equinix_network_device` [#115](https://github.com/equinix/terraform-provider-equinix/pull/115)
- Improved documentation [#123](https://github.com/equinix/terraform-provider-equinix/pull/123)
- Packngo version bumped to 0.25.0
- update go-getter to 1.5.11 for CWE-532 [#139](https://github.com/equinix/terraform-provider-equinix/pull/139)
- `equinix_metal_gateway` will wait for the Metal Gateway devices to pass through the "deleting" status [#138](https://github.com/equinix/terraform-provider-equinix/pull/138)
- E2E tests use data source `equinix_metal_plans` in all tests with a metal_device to check for available hardware [#215](https://github.com/equinix/terraform-provider-equinix/pull/215)

## 1.5.0 (March 24, 2022)

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

- `equinix_ecx_l2_serviceprofile` detecting diff after refresh [#90](https://github.com/equinix/terraform-provider-equinix/pull/90)
- `equinix_network_device` allow value 0 for additional bandwidth [#91](https://github.com/equinix/terraform-provider-equinix/pull/91)
- `equinix_network_device` hostname max length now match portal limits [#92](https://github.com/equinix/terraform-provider-equinix/pull/92)
- `equinix_ecx_l2_connection` will wait for the secondary connection destroy [#103](https://github.com/equinix/terraform-provider-equinix/pull/103)
- `equinix_ecx_l2_connection` named_tag now is idempotent [#97](https://github.com/equinix/terraform-provider-equinix/issues/97)
- `equinix_ecx_l2_connection` was not storing secondary connection fields [#103](https://github.com/equinix/terraform-provider-equinix/pull/103)

ENHANCEMENTS:

- `service_token` added to `equinix_ecx_l2_connection` [#96](https://github.com/equinix/terraform-provider-equinix/issues/96)
- `service_token` for secondary_connection added to `equinix_ecx_l2_connection` [#111](https://github.com/equinix/terraform-provider-equinix/pull/111)
- update documentation links for timeout parameters  [#101](https://github.com/equinix/terraform-provider-equinix/pull/101)
- `cluster_details` added to `equinix_network_device` [#105](https://github.com/equinix/terraform-provider-equinix/pull/105)

## 1.4.0 (January 14, 2022)

NOTES:

- `equinix_acl_template` argument `metro_code` is now deprecated [#67](https://github.com/equinix/terraform-provider-equinix/pull/67)
- `equinix_acl_template` argument `inbound_rule.#.subnets` is now deprecated [#67](https://github.com/equinix/terraform-provider-equinix/pull/67)
- `equinix_acl_template` attribute `device_id` is now deprecated [#67](https://github.com/equinix/terraform-provider-equinix/pull/67)
- `equinix_ecx_l2_connection_accepter` is now deprecated [#64](https://github.com/equinix/terraform-provider-equinix/pull/64)
- `equinix_network_device_link` argument `device.interface_id` changes taint the resource [#77](https://github.com/equinix/terraform-provider-equinix/pull/77)
- `equinix_network_device_link` attribute `link.src_zone_code` is now deprecated and optional [#77](https://github.com/equinix/terraform-provider-equinix/pull/77)
- `equinix_network_device_link` attribute `link.dest_zone_code` is now deprecated and optional [#77](https://github.com/equinix/terraform-provider-equinix/pull/77)

BUG FIXES:

- `equinix_ecx_l2_connection` will wait for the secondary connection [#87](https://github.com/equinix/terraform-provider-equinix/pull/87)
- `equinix_network_device_link` no longer jitters on zone code fields [#77](https://github.com/equinix/terraform-provider-equinix/pull/77)

ENHANCEMENTS:

- `equinix_acl_template` attribute `device_details` (`uuid`, `name`, `acl_status`) was added [#67](https://github.com/equinix/terraform-provider-equinix/pull/67)
- `equinix_ecx_l2_connection` attribute `actions` was added [#86](https://github.com/equinix/terraform-provider-equinix/pull/86)
- `equinix_network_device_link` argument `device.asn` is now optional [#77](https://github.com/equinix/terraform-provider-equinix/pull/77)
- `equinix_network_device_link` argument `device.subnet` is now optional [#77](https://github.com/equinix/terraform-provider-equinix/pull/77)
- fix connectivity example to establish an Azure connection [#71](https://github.com/equinix/terraform-provider-equinix/pull/71)
- replace Travis CI with GitHub Workflows [#65](https://github.com/equinix/terraform-provider-equinix/pull/65)
- update go modules and update go to 1.17 [#63](https://github.com/equinix/terraform-provider-equinix/pull/63)

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
