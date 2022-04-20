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

// RouteTableEntryOrFilter struct for RouteTableEntryOrFilter
type RouteTableEntryOrFilter struct {
	Or []RouteTableEntrySimpleExpression `json:"or,omitempty"`
}

// NewRouteTableEntryOrFilter instantiates a new RouteTableEntryOrFilter object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRouteTableEntryOrFilter() *RouteTableEntryOrFilter {
	this := RouteTableEntryOrFilter{}
	return &this
}

// NewRouteTableEntryOrFilterWithDefaults instantiates a new RouteTableEntryOrFilter object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRouteTableEntryOrFilterWithDefaults() *RouteTableEntryOrFilter {
	this := RouteTableEntryOrFilter{}
	return &this
}

// GetOr returns the Or field value if set, zero value otherwise.
func (o *RouteTableEntryOrFilter) GetOr() []RouteTableEntrySimpleExpression {
	if o == nil || o.Or == nil {
		var ret []RouteTableEntrySimpleExpression
		return ret
	}
	return o.Or
}

// GetOrOk returns a tuple with the Or field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RouteTableEntryOrFilter) GetOrOk() ([]RouteTableEntrySimpleExpression, bool) {
	if o == nil || o.Or == nil {
		return nil, false
	}
	return o.Or, true
}

// HasOr returns a boolean if a field has been set.
func (o *RouteTableEntryOrFilter) HasOr() bool {
	if o != nil && o.Or != nil {
		return true
	}

	return false
}

// SetOr gets a reference to the given []RouteTableEntrySimpleExpression and assigns it to the Or field.
func (o *RouteTableEntryOrFilter) SetOr(v []RouteTableEntrySimpleExpression) {
	o.Or = v
}

func (o RouteTableEntryOrFilter) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Or != nil {
		toSerialize["or"] = o.Or
	}
	return json.Marshal(toSerialize)
}

type NullableRouteTableEntryOrFilter struct {
	value *RouteTableEntryOrFilter
	isSet bool
}

func (v NullableRouteTableEntryOrFilter) Get() *RouteTableEntryOrFilter {
	return v.value
}

func (v *NullableRouteTableEntryOrFilter) Set(val *RouteTableEntryOrFilter) {
	v.value = val
	v.isSet = true
}

func (v NullableRouteTableEntryOrFilter) IsSet() bool {
	return v.isSet
}

func (v *NullableRouteTableEntryOrFilter) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRouteTableEntryOrFilter(val *RouteTableEntryOrFilter) *NullableRouteTableEntryOrFilter {
	return &NullableRouteTableEntryOrFilter{value: val, isSet: true}
}

func (v NullableRouteTableEntryOrFilter) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRouteTableEntryOrFilter) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


