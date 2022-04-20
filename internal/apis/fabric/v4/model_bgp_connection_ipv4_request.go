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

// BGPConnectionIpv4Request struct for BGPConnectionIpv4Request
type BGPConnectionIpv4Request struct {
	// Customer side peering ip
	CustomerPeerIp string `json:"customerPeerIp"`
}

// NewBGPConnectionIpv4Request instantiates a new BGPConnectionIpv4Request object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBGPConnectionIpv4Request(customerPeerIp string) *BGPConnectionIpv4Request {
	this := BGPConnectionIpv4Request{}
	this.CustomerPeerIp = customerPeerIp
	return &this
}

// NewBGPConnectionIpv4RequestWithDefaults instantiates a new BGPConnectionIpv4Request object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBGPConnectionIpv4RequestWithDefaults() *BGPConnectionIpv4Request {
	this := BGPConnectionIpv4Request{}
	return &this
}

// GetCustomerPeerIp returns the CustomerPeerIp field value
func (o *BGPConnectionIpv4Request) GetCustomerPeerIp() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.CustomerPeerIp
}

// GetCustomerPeerIpOk returns a tuple with the CustomerPeerIp field value
// and a boolean to check if the value has been set.
func (o *BGPConnectionIpv4Request) GetCustomerPeerIpOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CustomerPeerIp, true
}

// SetCustomerPeerIp sets field value
func (o *BGPConnectionIpv4Request) SetCustomerPeerIp(v string) {
	o.CustomerPeerIp = v
}

func (o BGPConnectionIpv4Request) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["customerPeerIp"] = o.CustomerPeerIp
	}
	return json.Marshal(toSerialize)
}

type NullableBGPConnectionIpv4Request struct {
	value *BGPConnectionIpv4Request
	isSet bool
}

func (v NullableBGPConnectionIpv4Request) Get() *BGPConnectionIpv4Request {
	return v.value
}

func (v *NullableBGPConnectionIpv4Request) Set(val *BGPConnectionIpv4Request) {
	v.value = val
	v.isSet = true
}

func (v NullableBGPConnectionIpv4Request) IsSet() bool {
	return v.isSet
}

func (v *NullableBGPConnectionIpv4Request) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBGPConnectionIpv4Request(val *BGPConnectionIpv4Request) *NullableBGPConnectionIpv4Request {
	return &NullableBGPConnectionIpv4Request{value: val, isSet: true}
}

func (v NullableBGPConnectionIpv4Request) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBGPConnectionIpv4Request) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


