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

// PriceLocation struct for PriceLocation
type PriceLocation struct {
	MetroCode *string `json:"metroCode,omitempty"`
}

// NewPriceLocation instantiates a new PriceLocation object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPriceLocation() *PriceLocation {
	this := PriceLocation{}
	return &this
}

// NewPriceLocationWithDefaults instantiates a new PriceLocation object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPriceLocationWithDefaults() *PriceLocation {
	this := PriceLocation{}
	return &this
}

// GetMetroCode returns the MetroCode field value if set, zero value otherwise.
func (o *PriceLocation) GetMetroCode() string {
	if o == nil || o.MetroCode == nil {
		var ret string
		return ret
	}
	return *o.MetroCode
}

// GetMetroCodeOk returns a tuple with the MetroCode field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PriceLocation) GetMetroCodeOk() (*string, bool) {
	if o == nil || o.MetroCode == nil {
		return nil, false
	}
	return o.MetroCode, true
}

// HasMetroCode returns a boolean if a field has been set.
func (o *PriceLocation) HasMetroCode() bool {
	if o != nil && o.MetroCode != nil {
		return true
	}

	return false
}

// SetMetroCode gets a reference to the given string and assigns it to the MetroCode field.
func (o *PriceLocation) SetMetroCode(v string) {
	o.MetroCode = &v
}

func (o PriceLocation) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.MetroCode != nil {
		toSerialize["metroCode"] = o.MetroCode
	}
	return json.Marshal(toSerialize)
}

type NullablePriceLocation struct {
	value *PriceLocation
	isSet bool
}

func (v NullablePriceLocation) Get() *PriceLocation {
	return v.value
}

func (v *NullablePriceLocation) Set(val *PriceLocation) {
	v.value = val
	v.isSet = true
}

func (v NullablePriceLocation) IsSet() bool {
	return v.isSet
}

func (v *NullablePriceLocation) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePriceLocation(val *PriceLocation) *NullablePriceLocation {
	return &NullablePriceLocation{value: val, isSet: true}
}

func (v NullablePriceLocation) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePriceLocation) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


