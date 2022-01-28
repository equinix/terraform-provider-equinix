## 2.0.3 (March 03, 2021)

BUG FIXES:

* creation of redundant connection from a single device is now reflecting
upstream API logic

## 2.0.2 (February 24, 2021)

NOTES:

* upgraded rest-go to v1.3.0 + testify and httpmock modules

## 2.0.1 (February 12, 2021)

BUG FIXES:

* upgraded to rest-go v1.2.2 to solve pagination issues
[equinix/terraform-provider-equinix#41](https://github.com/equinix/terraform-provider-equinix/issues/41)

## 2.0.0 (February 01, 2021)

BREAKING CHANGES:

* General change in create functions: instead of returning copy of input structure,
that might be outdated anyway, functions return pointers to created object's identifiers.
Change affects below functions:
  * `CreateL2Connection`
  * `CreateL2Connection`
  * `CreateL2ServiceProfile`
* `UpdateL2ServiceProfile` doest not return service profile's structure anymore
* General change in L2 Connection and Service Profile models: all basic type fields
are pointers now. Change affects below structures:
  * `L2Connection`
  * `L2ConnectionAdditionalInfo`
  * `L2ConnectionAction`
  * `L2ConnectionActionData`
  * `L2ConnectionToConfirm`
  * `L2ConnectionConfirmation`
  * `L2ServiceProfile`
  * `L2ServiceProfilePort`
  * `L2ServiceProfileSpeedBand`
  * `L2ServiceProfileFeatures`
  * `Port`
  * `L2SellerProfileMetro`
  * `L2SellerProfileAdditionalInfo`

ENHANCEMENTS:

* **L2Connection** redundant connection creation requests maps additional secondary
connection attributes ([equinix/terraform-provider-equinix#17](https://github.com/equinix/terraform-provider-equinix/issues/17)):
  * Speed
  * SpeedUnit
  * ProfileUUID
  * AuthorizationKey
  * SellerRegion
  * SellerMetroCode
  * InterfaceID

## 1.2.0 (January 07, 2021)

NOTES:

* this version of module started to use `equinix/rest-go` client
for any REST interactions with Equinix APIs
* ECX names were removed from descriptions and documentation in favor
of Equinix Fabric name

FEATURES:

* **L2Connection**: `func GetL2OutgoingConnections()` gives possibility to fetch
 all a-side (outgoing) connections for a customer account associated with
authenticated application

ENHANCEMENTS:

* **L2Connection** added additional attributes:
  * *Actions* provide details about pending actions to complete connection provisioning
  * *DeviceInterfaceID* indicates network interface identifier on a network device
  * *ProviderStatus* indicates connection status on a z-side
  * *RedundancyType* indicates whether connection is primary or secondary
  (for redundant connections)

## 1.1.0 (September 22, 2020)

FEATURES:

* **L2Connection** can be created with device identifier (in addition to port identifier)
 to allow interconnections with Network Edge devices

ENHANCEMENTS:

* **L2ServiceProfile** model and fetch logic was enriched with additional data
 useful when fetching seller profiles:
  * additional information that can be provided when creating connection
  * seller's metro locations
  * profile encapsulation
  * global organization and organization names

## 1.0.0 (July 31, 2020)

NOTES:

* first version of Equinix Cloud Exchange Fabric Go client

FEATURES:

* **L2ServiceProfile**: possibility to create, fetch, update (name and bandwidth),
 remove private and public service profiles
* **L2ServiceProfile**: possibility to fetch seller service profiles.
* **L2Connection**: possibility to create, fetch, update, remove ECX Fabric
 layer 2 connections
* **L2Connection**: possibility to approve layer2 connection with provider's
 access and secret keys (AWS use case)
* **UserPort**: possibly to fetch list of user ports
