## 3.1.0 (Unreleased)

BREAKING CHANGES:
- packngo updated to v0.4.1+, changing the API endpoint from api.packet.net to api.equinix.com/metal/v1

FEATURES:
- [#249](https://github.com/packethost/terraform-provider-packet/pull/249) New datasource `packet_project_ssh_key`

IMPROVEMENTS:
- `packet_device` datasource should query by hostname much faster
- `packet_device_network_type` conversions should be more reliable
- Test sweeper added for SSH keys
- Acceptance testing moved to Github Actions
- Improved logging when resources are not found and removed from state
- Device vlan attachments and network types will be removed from state when their device is removed

## 3.0.1 (August 20, 2020)

BREAKING CHANGES:
- [#246](https://github.com/packethost/terraform-provider-packet/pull/246) Updates URLs to reflect move from terraform-providers/terraform-provider-packet to packethost/terraform-provider-packet

## 3.0.0 (July 21, 2020)

BREAKING CHANGES:
- [#240](https://github.com/packethost/terraform-provider-packet/pull/240) Attribute `packet_device.network_type` is now read-only. Writes (and Reads, when writes are used) should use `packet_device_network_type`.

IMPROVEMENTS:
- [#240](https://github.com/packethost/terraform-provider-packet/pull/240) New resource `packet_device_network_type` for handling network modes of devices.
- removal of deleted `packet_device` attributes: `facility` (deprecated in 2.0.0) and `ip_address_types` (removed in 2.7.5)

## 2.10.1 (July 10, 2020)

BUG FIXES:
- [#239](https://github.com/packethost/terraform-provider-packet/pull/239) Fix conversion from nil during WaitForState loops

## 2.10.0 (July 07, 2020)

FEATURES:
- [#237](https://github.com/packethost/terraform-provider-packet/pull/237) Switch `userdata` attribute of packet_device to ForceNew

## 2.9.0 (May 12, 2020)

FEATURES:
- [#235](https://github.com/packethost/terraform-provider-packet/pull/235) Removed public_ipv4_subnet_size attr of packet_device. It's deprecated and replaced by ip_address attribute
- [#232](https://github.com/packethost/terraform-provider-packet/pull/232) New datasource packet_device_bgp_neighbors

## 2.8.1 (April 22, 2020)

BUG FIXES:
- [#231](https://github.com/packethost/terraform-provider-packet/pull/231) Fix storage param for packet_device 

## 2.8.0 (March 23, 2020)

IMPROVEMENTS:
- [#225](https://github.com/packethost/terraform-provider-packet/pull/225) Bump packngo to version correctly handling hybrid mode for n2.xlarge device plan
- [#226](https://github.com/packethost/terraform-provider-packet/pull/226) Remove deprecated code for ip_address_type attribute of the packet_device resource

## 2.7.5 (February 26, 2020)

FEATURES:
- [#223](https://github.com/packethost/terraform-provider-packet/pull/223) Add new list argument `ip_address` to `packet_device` to allow finer control of assigning subnets from reserved IP blocks to devices
- Deprecate `ip_address_types` for `packet_device`

## 2.7.4 (January 24, 2020)

IMPROVEMENTS:
- Add sweepers for Acceptance tests of Projects, Devices and Volumes
- Prefix names of testing resources with `tfacc-`, to mark them for sweepers
- Fixed links in docs to capture migration of Packet Knowledge Base

BUG FIXES:
 - Fix bug when devices were added to Terraform state
 - Update volumes sizes in Acceptance Tests

## 2.7.3 (December 12, 2019)

IMPROVEMENTS:
- [#210](https://github.com/packethost/terraform-provider-packet/pull/210) Distinguish HTTP 404 NotFound from API and from loadbalancer.

BUG FIXES:
- [#212](https://github.com/packethost/terraform-provider-packet/pull/212) Update retryablehttp dependency to fix #211.

## 2.7.2 (December 06, 2019)

IMPROVEMENTS:

- [#205](https://github.com/packethost/terraform-provider-packet/pull/205) packet_device: allow to force-detach volume on resource removal
- [#208](https://github.com/packethost/terraform-provider-packet/pull/208) packet_device: add caution note to hardware_reservation_id

BUG FIXES:
- [#207](https://github.com/packethost/terraform-provider-packet/pull/205) packet_bgp_session: fix ignored ipv6 value in address_family attribute

## 2.7.1 (December 03, 2019)

IMPROVEMENTS:
- [#202](https://github.com/packethost/terraform-provider-packet/pull/202) packet_volume_attachment documenation: show how to run attach script via Terraform

BUG FIXES:
- [#203](https://github.com/packethost/terraform-provider-packet/pull/203) Fix disappearing API queries in long waits, e.g. when creating ESXi device


## 2.7.0 (November 19, 2019)

FEATURES
- [#201] (https://github.com/packethost/terraform-provider-packet/issues/201) Removed resorce packet_connect

## 2.6.2 (November 16, 2019)

FEATURES
- [#198] (https://github.com/packethost/terraform-provider-packet/issues/198) Deprecated resource: packet_connect

IMPROVEMENTS:
- fix of snapshot_policies example in packet_volume documentation

## 2.6.1 (November 05, 2019)

IMPROVEMENTS:
- Resource imports
  - [#188](https://github.com/packethost/terraform-provider-packet/pull/188) packet_vlan
  - [#189](https://github.com/packethost/terraform-provider-packet/pull/189) packet_ssh_key
  - [#190](https://github.com/packethost/terraform-provider-packet/pull/190) packet_organization
  - [#193](https://github.com/packethost/terraform-provider-packet/pull/193) packet_project_ssh_key
   

BUG FIXES:
- [#194](https://github.com/packethost/terraform-provider-packet/pull/194) HTTP client hotfix - don't retry on HTTP responses with 5xx status code
- [#196](https://github.com/packethost/terraform-provider-packet/pull/196) Add mutext to guard VLAN detachment


## 2.6.0 (October 28, 2019)

IMPROVEMENTS:
- [#187](https://github.com/packethost/terraform-provider-packet/pull/187) Bump Go version to 1.13
- [#186](https://github.com/packethost/terraform-provider-packet/pull/186) Update Packet SDK in order to use retryable HTTP client
- [#162](https://github.com/packethost/terraform-provider-packet/pull/162) Datasource for packet_volume
- [#161](https://github.com/packethost/terraform-provider-packet/pull/161) Datasource for packet_project
- [#160](https://github.com/packethost/terraform-provider-packet/pull/160) Datasource for packet_organization
- [#159](https://github.com/packethost/terraform-provider-packet/pull/159) Datasource for packet_spotmarket_request

BUG FIXES:
- [#180](https://github.com/packethost/terraform-provider-packet/pull/159) Fixes device import 
- [#183](https://github.com/packethost/terraform-provider-packet/pull/184) Added missing error handling when waiting for spot instances to complete

## 2.5.0 (October 18, 2019)

IMPROVEMENTS:
- [#173](https://github.com/packethost/terraform-provider-packet/issues/173) Migrate to TF Plugin SDK
- [#174](https://github.com/packethost/terraform-provider-packet/issues/174) Add timeouts to packet_device
- [#169](https://github.com/packethost/terraform-provider-packet/issues/169) Make userdata not ForceNew
- [#171](https://github.com/packethost/terraform-provider-packet/issues/171) New datasource for selecting preallocated IP blocks in a project

BUG FIXES:
- [#175](https://github.com/packethost/terraform-provider-packet/issues/175) Better error when device fails to provision

## 2.4.0 (September 23, 2019)

IMPROVEMENTS:
- [#168](https://github.com/packethost/terraform-provider-packet/issues/168) Mention that users should use Full HW reservation ID
- [#163](https://github.com/packethost/terraform-provider-packet/issues/163) Add description attribute to packet_reserved_ip_block


## 2.3.0 (August 07, 2019)

BUG FIXES:
- [#156](https://github.com/packethost/terraform-provider-packet/issues/156) Fix filtering logic in packet_operating_system datasource

IMPROVEMENTS:
- [#125](https://github.com/packethost/terraform-provider-packet/issues/125) Add argument `wait_for_reservation_deprovision` to packet_device, in order to wait for proper deprovision of reserved hardware, and avoid errors when attempting to create devices in recently-released hardware reservations.

FEATURES:
- [#158](https://github.com/packethost/terraform-provider-packet/issues/158) New datasource for packet_device

## 2.2.1 (June 05, 2019)

- resource/packet_port_vlan_attachment: Avoid parallel assignment ([#152](https://github.com/packethost/terraform-provider-packet/issues/152))

## 2.2.0 (May 13, 2019)

IMPROVEMENTS:
- [#150](https://github.com/packethost/terraform-provider-packet/pull/150) New ip_address_types attribute for packet_device
- [#147](https://github.com/packethost/terraform-provider-packet/pull/147) Improve error message when trying to create device in nonexistent project


## 2.1.0 (April 30, 2019)

- [#145](https://github.com/packethost/terraform-provider-packet/pull/145) Terraform SDK upgrade with compatibility for Terraform v0.12

## 2.0.0 (April 23, 2019)

- [#144](https://github.com/packethost/terraform-provider-packet/pull/144) Support for Native VLAN in packet_port_vlan_attachment
- [#142](https://github.com/packethost/terraform-provider-packet/pull/142) Remove deprecated facility, fix import of packet_device


## 1.7.2 (April 10, 2019)

BUG FIXES:
- [#140](https://github.com/packethost/terraform-provider-packet/pull/140) Relax the network_type attribute of packet_device, fixing #138
- [#139](https://github.com/packethost/terraform-provider-packet/pull/139) Fix facility json tag for packet_sport_market_request

## 1.7.1 (April 08, 2019)

BUG FIXES:
- [#137](https://github.com/packethost/terraform-provider-packet/pull/137) Remove Disbond call from port-vlan-attachment creation function, in order to fix use-case for layer2-bonded

## 1.7.0 (April 04, 2019)

IMPROVEMENTS:
- [#135](https://github.com/packethost/terraform-provider-packet/pull/135) Add default_policy flag to packet_bgp_session resource

## 1.6.0 (March 29, 2019)

IMPROVEMENTS:
- Documetnation fixes

FEATURES:
- [#132](https://github.com/packethost/terraform-provider-packet/pull/132) New resource /packet_connect: connection to VLANs in other cloud providers

## 1.5.0 (March 20, 2019)

IMPROVEMENTS:
- [#114](https://github.com/packethost/terraform-provider-packet/pull/114) Bump Terraform version to 0.11.11 in order to see JSON from HTTP responses
- Documentation fixes
- Packet Go library updates
- [#122](https://github.com/packethost/terraform-provider-packet/pull/122) backend_transfer attribute in packet_project

FEATURES:
- [#86](https://github.com/packethost/terraform-provider-packet/pull/86) Layer 2 support: network_type attribute in packet_device, new resource packet_port_vlan_attachment

## 1.4.1 (February 21, 2019)

IMPROVEMENTS:

- [#112](https://github.com/packethost/terraform-provider-packet/pull/112) Remove strict validation for facilities, in order to allow non-public private facilities and testing new facilites
  
## 1.4.0 (February 19, 2019)

FEATURES:

- [#101](https://github.com/packethost/terraform-provider-packet/pull/101) Bump Go version to 1.11.5 and switch to Go Modules
- [#96](https://github.com/packethost/terraform-provider-packet/pull/96) New resource/packet_project_ssh_key: Resource for Project SSH Keys
- [#93](https://github.com/packethost/terraform-provider-packet/pull/93) resource/packet_device: Allow list of facilities and "any" facility

IMPROVEMENTS:

- [#99](https://github.com/packethost/terraform-provider-packet/pull/99) resource/packet_reserved_ip_block: extend to allow for Global floating IP blocks

- Various doc improvements

BUG FIXES:

- [#104](https://github.com/packethost/terraform-provider-packet/pull/104) Fix empty error messages on invalid credentials
- [#111](https://github.com/packethost/terraform-provider-packet/pull/111) Fix handling of resources deleted out of Terraform

## 1.3.2 (February 06, 2019)

IMPROVEMENTS:

- [#95](https://github.com/packethost/terraform-provider-packet/pull/95) Hotfix - facility df2 was missing from API lib listing

## 1.3.1 (February 04, 2019)

IMPROVEMENTS:

- [#92](https://github.com/packethost/terraform-provider-packet/pull/92) Hotfix of device network order, back to: 0. Public IPv4, 1. IPv6, 2. Private IPv4

## 1.3.0 (February 01, 2019)

FEATURES:

* [#88](https://github.com/packethost/terraform-provider-packet/pull/88) Support for BGP resources
* [#87](https://github.com/packethost/terraform-provider-packet/pull/87) Upgrade to Go 1.11
* [#85](https://github.com/packethost/terraform-provider-packet/pull/85) resource/packet_vlan: New resource for VLANs

IMPROVEMENTS:

* [#89](https://github.com/packethost/terraform-provider-packet/pull/89) Impose explicit order on network configurations in device resource

## 1.2.5 (September 28, 2018)

FEATURES:

* [#72](https://github.com/packethost/terraform-provider-packet/pull/72) resource/packet_spot_market_request: New resource for Spot Market Request for devices.
* [#71](https://github.com/packethost/terraform-provider-packet/pull/71) datasource/packet_spot_market_price: New datasource for lookup of current hourly spot market price of devices based on location and plan
* [#70](https://github.com/packethost/terraform-provider-packet/pull/70) datasource/packet_operating_system: New datasource for OS lookup

IMPROVEMENTS:

- [#73](https://github.com/packethost/terraform-provider-packet/pull/73) - devices, projects and volumes are now importable (see `terraform import` doc)
- [#69](https://github.com/packethost/terraform-provider-packet/pull/69) - in Device docs, explain how to get OS slugs

BUG FIXES:

- [#74](https://github.com/packethost/terraform-provider-packet/issues/74) fix of broken links in device and ip_attachment docs


## 1.2.4 (May 31, 2018)

BUG FIXES:

- `r/packet_ip_attachment` - handling IP attachments being deleted outside of Terraform ([#68](https://github.com/packethost/terraform-provider-packet/issues/68))

## 1.2.3 (April 27, 2018)

- [#61](https://github.com/packethost/terraform-provider-packet/issues/61), fix volume resource update
- [#63](https://https://github.com/packethost/terraform-provider-packet/pull/63), add Organization resource, add org attirbute to project resource

## 1.2.2 (April 17, 2018)

IMPROVEMENTS:

- [#58](https://github.com/packethost/terraform-provider-packet/issues/58), properly fix resource updates

## 1.2.1 (April 16, 2018)

IMPROVEMENTS:

- [#57](https://github.com/packethost/terraform-provider-packet/issues/57), fix for update of PXE attributes of device resource
- [#52](https://github.com/packethost/terraform-provider-packet/issues/52) fix for project resource update
- [#49](https://github.com/packethost/terraform-provider-packet/issues/49) fix for crash on SSH key update
- [#50](https://github.com/packethost/terraform-provider-packet/issues/50) fix for device update, adds `description` attribute


## 1.2.0 (January 23, 2018)

BACKWARDS INCOMPATIBILITIES / NOTES:

* [#37](https://github.com/packethost/terraform-provider-packet/issues/37), changes computation of packet_reserved_ip_block.cidr_notation. In case you use that attribute down the chain, you will need to refresh and recreate the resource in which you use it.

IMPROVEMENTS:

* datasource/packet_precreated_ip_block: Add Datasource for precreated IP blocks, so that users can assign subnets from those ([#36](https://github.com/packethost/terraform-provider-packet/issues/36))
* mark device userdata as ForceNew, fixing ([#42](https://github.com/packethost/terraform-provider-packet/issues/42))
* Add support for CPR (custom partitioning and RAID)([#35](https://github.com/packethost/terraform-provider-packet/pull/35))

## 1.1.0 (October 09, 2017)

INTERNAL:

* Capture breaking change in the Packet API ((https://github.com/packethost/packngo/pull/47))

## 1.0.0 (September 28, 2017)

INTERNAL:

* provider: Add logging transport for HTTP client (`DEBUG` log now contains all HTTP requests & responses) ([#25](https://github.com/packethost/terraform-provider-packet/issues/25))

IMPROVEMENTS:

* resource/packet_device: Add `public_ipv4_subnet_size` field ([#7](https://github.com/packethost/terraform-provider-packet/issues/7))
* resource/packet_device: Add `hardware_reservation_id` field ([#14](https://github.com/packethost/terraform-provider-packet/issues/14))
* resource/packet_device: Add `root_password` attribute ([#15](https://github.com/packethost/terraform-provider-packet/issues/15))
* resource/packet_volume_attachment: Add resource for attaching volumes ([#9](https://github.com/packethost/terraform-provider-packet/issues/9))
* resource/packet_reserved_ip_block: Add resource for reserving blocks of IP addresses ([#21](https://github.com/packethost/terraform-provider-packet/issues/21))
* resource/packet_ip_attachment: Add resource for attaching floating IP addresses to devices ([#21](https://github.com/packethost/terraform-provider-packet/issues/21))
* resource/packet_devce: Allow to use next-available hardware reservation ([#26](https://github.com/packethost/terraform-provider-packet/issues/26))
* resource/packet_reserved_ip_block: Make reserved IP block resource importable ([#29](https://github.com/packethost/terraform-provider-packet/issues/29))


## 0.1.0 (June 21, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
