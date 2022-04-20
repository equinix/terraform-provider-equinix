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

// Actions Connection action type
type Actions struct {
}

// NewActions instantiates a new Actions object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewActions() *Actions {
	this := Actions{}
	return &this
}

// NewActionsWithDefaults instantiates a new Actions object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewActionsWithDefaults() *Actions {
	this := Actions{}
	return &this
}

func (o Actions) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	return json.Marshal(toSerialize)
}

type NullableActions struct {
	value *Actions
	isSet bool
}

func (v NullableActions) Get() *Actions {
	return v.value
}

func (v *NullableActions) Set(val *Actions) {
	v.value = val
	v.isSet = true
}

func (v NullableActions) IsSet() bool {
	return v.isSet
}

func (v *NullableActions) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableActions(val *Actions) *NullableActions {
	return &NullableActions{value: val, isSet: true}
}

func (v NullableActions) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableActions) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


