/*
Equinix Fabric API

Equinix Fabric is an advanced software-defined interconnection solution that enables you to directly, securely and dynamically connect to distributed infrastructure and digital ecosystems on platform Equinix via a single port, Customers can use Fabric to connect to: </br> 1. Cloud Service Providers - Clouds, network and other service providers.  </br> 2. Enterprises - Other Equinix customers, vendors and partners.  </br> 3. Myself - Another customer instance deployed at Equinix. </br>

API version: 4.2
Contact: api-support@equinix.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package v4

import (
	"encoding/json"
)

// BGPConnectionIpv6 struct for BGPConnectionIpv6
type BGPConnectionIpv6 struct {
	// Customer side peering ip
	CustomerPeerIp string `json:"customerPeerIp"`
	// Equinix side peering ip
	EquinixPeerIp string `json:"equinixPeerIp"`
}

// NewBGPConnectionIpv6 instantiates a new BGPConnectionIpv6 object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBGPConnectionIpv6(customerPeerIp string, equinixPeerIp string) *BGPConnectionIpv6 {
	this := BGPConnectionIpv6{}
	this.CustomerPeerIp = customerPeerIp
	this.EquinixPeerIp = equinixPeerIp
	return &this
}

// NewBGPConnectionIpv6WithDefaults instantiates a new BGPConnectionIpv6 object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBGPConnectionIpv6WithDefaults() *BGPConnectionIpv6 {
	this := BGPConnectionIpv6{}
	return &this
}

// GetCustomerPeerIp returns the CustomerPeerIp field value
func (o *BGPConnectionIpv6) GetCustomerPeerIp() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.CustomerPeerIp
}

// GetCustomerPeerIpOk returns a tuple with the CustomerPeerIp field value
// and a boolean to check if the value has been set.
func (o *BGPConnectionIpv6) GetCustomerPeerIpOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CustomerPeerIp, true
}

// SetCustomerPeerIp sets field value
func (o *BGPConnectionIpv6) SetCustomerPeerIp(v string) {
	o.CustomerPeerIp = v
}

// GetEquinixPeerIp returns the EquinixPeerIp field value
func (o *BGPConnectionIpv6) GetEquinixPeerIp() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.EquinixPeerIp
}

// GetEquinixPeerIpOk returns a tuple with the EquinixPeerIp field value
// and a boolean to check if the value has been set.
func (o *BGPConnectionIpv6) GetEquinixPeerIpOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.EquinixPeerIp, true
}

// SetEquinixPeerIp sets field value
func (o *BGPConnectionIpv6) SetEquinixPeerIp(v string) {
	o.EquinixPeerIp = v
}

func (o BGPConnectionIpv6) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["customerPeerIp"] = o.CustomerPeerIp
	}
	if true {
		toSerialize["equinixPeerIp"] = o.EquinixPeerIp
	}
	return json.Marshal(toSerialize)
}

type NullableBGPConnectionIpv6 struct {
	value *BGPConnectionIpv6
	isSet bool
}

func (v NullableBGPConnectionIpv6) Get() *BGPConnectionIpv6 {
	return v.value
}

func (v *NullableBGPConnectionIpv6) Set(val *BGPConnectionIpv6) {
	v.value = val
	v.isSet = true
}

func (v NullableBGPConnectionIpv6) IsSet() bool {
	return v.isSet
}

func (v *NullableBGPConnectionIpv6) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBGPConnectionIpv6(val *BGPConnectionIpv6) *NullableBGPConnectionIpv6 {
	return &NullableBGPConnectionIpv6{value: val, isSet: true}
}

func (v NullableBGPConnectionIpv6) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBGPConnectionIpv6) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


