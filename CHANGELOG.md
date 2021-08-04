## 3.1.0 (Aug 4, 2021)

## Improvements

* Added `mode` to `metal_connection` resource and datasource #129
* Added `tags` to `metal_connection` resource and datasource #134
* Added `project_id` to `metal_connection` datasource #145 
* Added `tags` to `metal_virtual_circuit` resource and datasource #134
* Added `description` to `metal_virtual_circuit` resource and datasource #134
* Added `speed` to `metal_virtual_circuit` resource and datasource #134
* `metal_virtual_circuit` field `nni_vlan` is now optional #129
* Improved acceptance test success rates #170 

## Bug Fixes

* fixed crashing when VLAN is not attached to a virtual circuit #144

## 3.0.0 (Jul 27, 2021)

BREAKING CHANGES:
- Upgraded the Terraform Plugin SDK to v2 (Terraform v0.12+ now required) #113
- `metal_volume` resource has been removed #112
- `metal_volume` datasource has been removed #112
- `metal_volume_attachment` resource has been removed #112

FEATURES:
- New resource and datasource `metal_gateway` #157 
- New resources for API keys: `metal_user_api_key` and `metal_project_api_key` #147 

ENHANCEMENTS:
- `metal_device` `reinstall` options have been added #152 

BUG FIXES:

- Metros will be treated as lower-case #126 / #119 
- Hardware reservation ID properly read in `metal_device` #167
- Crash of facility nil deref in datasource `metal_reserved_ip_block` #163 
- Handling of metro attribute in resource `metal_reserved_ip_block` #169 

IMPROVEMENTS:

- added `reinstall` block to `metal_device` #152 
- added `tags` to `metal_reserved_ip_block` #133
- packngo updated to 0.17.0 #137 / #151
- added `Description` attributes to resource structures #130
- Corrections to guides #111
- Provider example illustrates metro use #146
- CI will run go tests on forked PRs #154


## 2.1.0 (May 20, 2021)

BREAKING CHANGES:
- `metal_spot_market_request` field `locked` is now boolean (not string) ([#78](https://github.com/equinix/terraform-provider-metal/pull/78))

DEPRECATIONS:
- **Deprecated Resource:** `metal_volume` ([#99](https://github.com/equinix/terraform-provider-metal/pull/99))
- **Deprecated Resource:** `metal_volume_attachment` ([#99](https://github.com/equinix/terraform-provider-metal/pull/99))
- **Deprecated Data Source:** `metal_volume` ([#99](https://github.com/equinix/terraform-provider-metal/pull/99))
- removed unintended deprecation of `metal_device`.`facility` ([#82](https://github.com/equinix/terraform-provider-metal/pull/82))

BUG FIXES:
- removed deprecation of `metal_device`.`facility` ([#82](https://github.com/equinix/terraform-provider-metal/pull/82))
- fixed plugin crash handling certain errors ([#89](https://github.com/equinix/terraform-provider-metal/issues/89))

FEATURES:
- **New Resource:** `metal_connection` ([#76](https://github.com/equinix/terraform-provider-metal/pull/76))
- **New Resource:** `metal_virtual_circuit`([#92](https://github.com/equinix/terraform-provider-metal/pull/92))
- **New Data Source:** `metal_vlan` ([#67](https://github.com/equinix/terraform-provider-metal/pull/67))
- **New Data Source:** `metal_reserved_ip_blcok`([#80](https://github.com/equinix/terraform-provider-metal/pull/80))
- **New Data Source:** `metal_port`([#96](https://github.com/equinix/terraform-provider-metal/pull/96))
- **New Data Source:** `metal_hardware_reservation`([#100](https://github.com/equinix/terraform-provider-metal/pull/100))
- **New Guide:** Migrating From the Packet Provider ([#72](https://github.com/equinix/terraform-provider-metal/pull/72))
- **New Guide:** Upgrading Devices from Facilities to Metros ([#103](https://github.com/equinix/terraform-provider-metal/pull/103))

IMPROVEMENTS:
- added `metro` to `metal_connection` data source (`facility` now optional) ([#81](https://github.com/equinix/terraform-provider-metal/pull/81))
- added `devices_min` attribute to `metal_spot_market_request` data source ([#78](https://github.com/equinix/terraform-provider-metal/pull/78))
- added `devices_max` attribute to `metal_spot_market_request` data source ([#78](https://github.com/equinix/terraform-provider-metal/pull/78))
- added `max_bid_price` attribute to `metal_spot_market_request` data source ([#78](https://github.com/equinix/terraform-provider-metal/pull/78))
- added `facilities` attribute to `metal_spot_market_request` data source ([#78](https://github.com/equinix/terraform-provider-metal/pull/78))
- added `metro` attribute to `metal_spot_market_request` data source ([#78](https://github.com/equinix/terraform-provider-metal/pull/78))
- added `project_id` attribute to `metal_spot_market_request` data source ([#78](https://github.com/equinix/terraform-provider-metal/pull/78))
- added `plan` attribute to `metal_spot_market_request` data source ([#78](https://github.com/equinix/terraform-provider-metal/pull/78))
- added `end_dt` attribute to `metal_spot_market_request` data source ([#78](https://github.com/equinix/terraform-provider-metal/pull/78))
- `metal_spot_market_request` resources can now be imported ([#78](https://github.com/equinix/terraform-provider-metal/pull/78))
- `metal_project` `md5_pass` field is now flagged sensitive ([#93](https://github.com/equinix/terraform-provider-metal/issues/93))
- `metal_spot_market_request` `instance_parameters` are documented ([#104](https://github.com/equinix/terraform-provider-metal/issues/104))
- minor `metal_spot_market_request` `max_bid_price` fluctuations (<2%) are ignored to avoid jitter tainting ([#78](https://github.com/equinix/terraform-provider-metal/pull/78))
- updated device, organization, project, reserved IP block, SSH Key, and VLAN documentation to include import instructions ([#78](https://github.com/equinix/terraform-provider-metal/pull/78))
- packngo updated to v0.14.1 ([#106](https://github.com/equinix/terraform-provider-metal/pull/106))

## 2.0.1 (April 15, 2021)

BUG FIXES:

- Fixes `metal_port_vlan_attachment` for Metro VLAN attachment and fixes a bug that caused Facility VLANs to fail attachment when a Metro VLAN was present in the project ([#69](https://github.com/equinix/terraform-provider-metal/pull/69))

IMPROVEMENTS:

- added `metro` to `metal_spot_market_price` data source
- `facility` is now optional in the `metal_spot_market_price` data source
- documentation examples have been updated to use current plans, facilities,
  operating systems, and metros ([#36](https://github.com/equinix/terraform-provider-metal/pull/36))

## 2.0.0 (April 12, 2021)

This release includes support for Equinix Metal Metros. Read more about metros in the
[Equinix Metal Changelog: New Metros Feature Live](https://feedback.equinixmetal.com/changelog/new-metros-feature-live)
post.

Learn more about these changes may enhance and impact your deployments in the [GitHub Discussion](https://github.com/equinix/terraform-provider-metal/discussions/55).

BREAKING CHANGES:

- `facilities` field in `metal_device` is now optional and conflicts with `metro` ([#46](https://github.com/equinix/terraform-provider-metal/pull/46))
- `facilities` field in `metal_vlan` is now optional and conflicts with `metro` ([#53](https://github.com/equinix/terraform-provider-metal/pull/53))
- `facilities` field in `metal_reserved_ip_block` is now optional and conflicts with `metro` ([#56](https://github.com/equinix/terraform-provider-metal/pull/56))
- `facilities` field in `metal_spot_market_request` is now optional and conflicts with `metro` ([#63](https://github.com/equinix/terraform-provider-metal/pull/63))

FEATURES:

- **New Data Source:** `metal_metro` ([#58](https://github.com/equinix/terraform-provider-metal/pull/58))
- **New Data Source:** `metal_connection` ([#41](https://github.com/equinix/terraform-provider-metal/pull/41))

IMPROVEMENTS:

- `metro` field added to `metal_device` resource ([#46](https://github.com/equinix/terraform-provider-metal/pull/46))
- `metro` field added to `metal_vlan` resource ([#53](https://github.com/equinix/terraform-provider-metal/pull/53))
- `metro` field added to `metal_reserved_ip_block` resource ([#56](https://github.com/equinix/terraform-provider-metal/pull/56))
- `metro` field added to `metal_facility` data source ([#58](https://github.com/equinix/terraform-provider-metal/pull/58))
- `metro` field added to `metal_device` data source ([#46](https://github.com/equinix/terraform-provider-metal/pull/46))
- `metro` field added to `metal_precreate_ip_block` data source ([#58](https://github.com/equinix/terraform-provider-metal/pull/58))
- `metro` field added to `metal_spot_market_request` resource ([#63](https://github.com/equinix/terraform-provider-metal/pull/63))
- `facilities` field is computed in `metal_spot_market_request` resource when not specified ([#63](https://github.com/equinix/terraform-provider-metal/pull/63))
- `vxlan` field added to `metal_vlan` resource ([#53](https://github.com/equinix/terraform-provider-metal/pull/53))
- `metro` search parameter added to `metal_ip_block_ranges` data source ([#58](https://github.com/equinix/terraform-provider-metal/pull/58))
- `packngo` update to `v0.13.0` ([#63](https://github.com/equinix/terraform-provider-metal/pull/63))

BUG FIXES:

- virtual_circuit: fixed documentation page title([#59](https://github.com/equinix/terraform-provider-metal/pull/59))
- some deletion 404 and 403s were not treated as successful ([#47](https://github.com/equinix/terraform-provider-metal/pull/47))

## 1.1.0 (March 22, 2021)

BREAKING CHANGES:

- The environment variable `PACKET_AUTH_TOKEN` is deprecated. Use `METAL_AUTH_TOKEN`.
- `CustomData` fields may be parsed differently with packngo v0.6.0+

FIXES:

- `metal_project_ssh_key` `project_id` is reported correctly
- `metal_ssh_key` `project_id` is reported correctly
- Delete operations now treat 404 and 403 HTTP errors as successful

FEATURES:

- **New Data Source** `metal_facility`
- **New Provider Argument** `max_retries`
- **New Provider Argument** `max_retry_wait_seconds`
- Versioned User-Agent is reported by the HTTP Client
- `metal_device_network_type` field `type` now accepts `hybrid-bonded`

IMPROVEMENTS:

- Depends on packngo [v0.7.0](https://github.com/packethost/packngo/releases/tag/v0.7.0)
- Documentation corrections
- E2E testing has less test jitter and greater overall success
- Removed debug logs from metal\_operating\_system data source
- Removed debug logs from metal\_organization data source
- Removed debug logs from metal\_project data source

## 1.0.0 (December 08, 2020)

BREAKING CHANGES:

- [#1](https://github.com/equinix/terraform-provider-metal/issues/1)
  Users migrating from the Packet provider, please follow the instructions at
  the linked issue. In short, all v3.2.0 Packet provider resources have been
  renamed.

FEATURES:

- **New Resource** `metal_bgp_session`
- **New Resource** `metal_device`
- **New Resource** `metal_device_network_type`
- **New Resource** `metal_ip_attachment`
- **New Resource** `metal_organization`
- **New Resource** `metal_port_vlan_attachment`
- **New Resource** `metal_project`
- **New Resource** `metal_project`
- **New Resource** `metal_project_ssh_key`
- **New Resource** `metal_reserved_ip_block`
- **New Resource** `metal_spot_market_request`
- **New Resource** `metal_ssh_key`
- **New Resource** `metal_volume`
- **New Resource** `metal_vlan`
- **New Resource** `metal_volume_attachment`

- **New Data Resource** `metal_device`
- **New Data Resource** `metal_device_bgp_neighbors`
- **New Data Resource** `metal_ip_block_ranges`
- **New Data Resource** `metal_operating_system`
- **New Data Resource** `metal_organization`
- **New Data Resource** `metal_precreated_ip_block`
- **New Data Resource** `metal_project`
- **New Data Resource** `metal_project_ssh_key`
- **New Data Resource** `metal_spot_market_price`
- **New Data Resource** `metal_spot_market_request`
- **New Data Resource** `metal_volume`
