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

// NatIpv4Ipv6Configuration Network address translation rules
type NatIpv4Ipv6Configuration struct {
	// IPv4 or IPv6 specific configuration
	Rules []NatRule `json:"rules"`
}

// NewNatIpv4Ipv6Configuration instantiates a new NatIpv4Ipv6Configuration object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewNatIpv4Ipv6Configuration(rules []NatRule) *NatIpv4Ipv6Configuration {
	this := NatIpv4Ipv6Configuration{}
	this.Rules = rules
	return &this
}

// NewNatIpv4Ipv6ConfigurationWithDefaults instantiates a new NatIpv4Ipv6Configuration object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewNatIpv4Ipv6ConfigurationWithDefaults() *NatIpv4Ipv6Configuration {
	this := NatIpv4Ipv6Configuration{}
	return &this
}

// GetRules returns the Rules field value
func (o *NatIpv4Ipv6Configuration) GetRules() []NatRule {
	if o == nil {
		var ret []NatRule
		return ret
	}

	return o.Rules
}

// GetRulesOk returns a tuple with the Rules field value
// and a boolean to check if the value has been set.
func (o *NatIpv4Ipv6Configuration) GetRulesOk() ([]NatRule, bool) {
	if o == nil {
		return nil, false
	}
	return o.Rules, true
}

// SetRules sets field value
func (o *NatIpv4Ipv6Configuration) SetRules(v []NatRule) {
	o.Rules = v
}

func (o NatIpv4Ipv6Configuration) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["rules"] = o.Rules
	}
	return json.Marshal(toSerialize)
}

type NullableNatIpv4Ipv6Configuration struct {
	value *NatIpv4Ipv6Configuration
	isSet bool
}

func (v NullableNatIpv4Ipv6Configuration) Get() *NatIpv4Ipv6Configuration {
	return v.value
}

func (v *NullableNatIpv4Ipv6Configuration) Set(val *NatIpv4Ipv6Configuration) {
	v.value = val
	v.isSet = true
}

func (v NullableNatIpv4Ipv6Configuration) IsSet() bool {
	return v.isSet
}

func (v *NullableNatIpv4Ipv6Configuration) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableNatIpv4Ipv6Configuration(val *NatIpv4Ipv6Configuration) *NullableNatIpv4Ipv6Configuration {
	return &NullableNatIpv4Ipv6Configuration{value: val, isSet: true}
}

func (v NullableNatIpv4Ipv6Configuration) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableNatIpv4Ipv6Configuration) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


