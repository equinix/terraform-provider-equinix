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

// NatStaticMapping Network Address Translation configuration - STATIC_MAPPING
type NatStaticMapping struct {
	Ipv4 NatIpv4Ipv6Configuration `json:"ipv4"`
	Ipv6 NatIpv4Ipv6Configuration `json:"ipv6"`
}

// NewNatStaticMapping instantiates a new NatStaticMapping object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewNatStaticMapping(ipv4 NatIpv4Ipv6Configuration, ipv6 NatIpv4Ipv6Configuration, type_ string) *NatStaticMapping {
	this := NatStaticMapping{}
	this.Type = type_
	return &this
}

// NewNatStaticMappingWithDefaults instantiates a new NatStaticMapping object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewNatStaticMappingWithDefaults() *NatStaticMapping {
	this := NatStaticMapping{}
	return &this
}

// GetIpv4 returns the Ipv4 field value
func (o *NatStaticMapping) GetIpv4() NatIpv4Ipv6Configuration {
	if o == nil {
		var ret NatIpv4Ipv6Configuration
		return ret
	}

	return o.Ipv4
}

// GetIpv4Ok returns a tuple with the Ipv4 field value
// and a boolean to check if the value has been set.
func (o *NatStaticMapping) GetIpv4Ok() (*NatIpv4Ipv6Configuration, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Ipv4, true
}

// SetIpv4 sets field value
func (o *NatStaticMapping) SetIpv4(v NatIpv4Ipv6Configuration) {
	o.Ipv4 = v
}

// GetIpv6 returns the Ipv6 field value
func (o *NatStaticMapping) GetIpv6() NatIpv4Ipv6Configuration {
	if o == nil {
		var ret NatIpv4Ipv6Configuration
		return ret
	}

	return o.Ipv6
}

// GetIpv6Ok returns a tuple with the Ipv6 field value
// and a boolean to check if the value has been set.
func (o *NatStaticMapping) GetIpv6Ok() (*NatIpv4Ipv6Configuration, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Ipv6, true
}

// SetIpv6 sets field value
func (o *NatStaticMapping) SetIpv6(v NatIpv4Ipv6Configuration) {
	o.Ipv6 = v
}

func (o NatStaticMapping) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["ipv4"] = o.Ipv4
	}
	if true {
		toSerialize["ipv6"] = o.Ipv6
	}
	return json.Marshal(toSerialize)
}

type NullableNatStaticMapping struct {
	value *NatStaticMapping
	isSet bool
}

func (v NullableNatStaticMapping) Get() *NatStaticMapping {
	return v.value
}

func (v *NullableNatStaticMapping) Set(val *NatStaticMapping) {
	v.value = val
	v.isSet = true
}

func (v NullableNatStaticMapping) IsSet() bool {
	return v.isSet
}

func (v *NullableNatStaticMapping) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableNatStaticMapping(val *NatStaticMapping) *NullableNatStaticMapping {
	return &NullableNatStaticMapping{value: val, isSet: true}
}

func (v NullableNatStaticMapping) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableNatStaticMapping) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


