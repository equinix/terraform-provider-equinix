## 2.0.1 (April 15, 2021)

BUG FIXES:

- Fixes `metal_port_vlan_attachment` for Metro VLAN attachment and fixes a bug that caused Facility VLANs to fail attachment when a Metro VLAN was present in the project ([#69](https://github.com/equinix/terraform-provider-metal/pull/69))

IMPROVEMENTS:

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
