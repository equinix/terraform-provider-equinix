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

// BGPConnectionIpv4 struct for BGPConnectionIpv4
type BGPConnectionIpv4 struct {
	// Customer side peering ip
	CustomerPeerIp string `json:"customerPeerIp"`
	// Equinix side peering ip
	EquinixPeerIp string `json:"equinixPeerIp"`
}

// NewBGPConnectionIpv4 instantiates a new BGPConnectionIpv4 object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBGPConnectionIpv4(customerPeerIp string, equinixPeerIp string) *BGPConnectionIpv4 {
	this := BGPConnectionIpv4{}
	this.CustomerPeerIp = customerPeerIp
	this.EquinixPeerIp = equinixPeerIp
	return &this
}

// NewBGPConnectionIpv4WithDefaults instantiates a new BGPConnectionIpv4 object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBGPConnectionIpv4WithDefaults() *BGPConnectionIpv4 {
	this := BGPConnectionIpv4{}
	return &this
}

// GetCustomerPeerIp returns the CustomerPeerIp field value
func (o *BGPConnectionIpv4) GetCustomerPeerIp() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.CustomerPeerIp
}

// GetCustomerPeerIpOk returns a tuple with the CustomerPeerIp field value
// and a boolean to check if the value has been set.
func (o *BGPConnectionIpv4) GetCustomerPeerIpOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CustomerPeerIp, true
}

// SetCustomerPeerIp sets field value
func (o *BGPConnectionIpv4) SetCustomerPeerIp(v string) {
	o.CustomerPeerIp = v
}

// GetEquinixPeerIp returns the EquinixPeerIp field value
func (o *BGPConnectionIpv4) GetEquinixPeerIp() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.EquinixPeerIp
}

// GetEquinixPeerIpOk returns a tuple with the EquinixPeerIp field value
// and a boolean to check if the value has been set.
func (o *BGPConnectionIpv4) GetEquinixPeerIpOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.EquinixPeerIp, true
}

// SetEquinixPeerIp sets field value
func (o *BGPConnectionIpv4) SetEquinixPeerIp(v string) {
	o.EquinixPeerIp = v
}

func (o BGPConnectionIpv4) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["customerPeerIp"] = o.CustomerPeerIp
	}
	if true {
		toSerialize["equinixPeerIp"] = o.EquinixPeerIp
	}
	return json.Marshal(toSerialize)
}

type NullableBGPConnectionIpv4 struct {
	value *BGPConnectionIpv4
	isSet bool
}

func (v NullableBGPConnectionIpv4) Get() *BGPConnectionIpv4 {
	return v.value
}

func (v *NullableBGPConnectionIpv4) Set(val *BGPConnectionIpv4) {
	v.value = val
	v.isSet = true
}

func (v NullableBGPConnectionIpv4) IsSet() bool {
	return v.isSet
}

func (v *NullableBGPConnectionIpv4) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBGPConnectionIpv4(val *BGPConnectionIpv4) *NullableBGPConnectionIpv4 {
	return &NullableBGPConnectionIpv4{value: val, isSet: true}
}

func (v NullableBGPConnectionIpv4) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBGPConnectionIpv4) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


