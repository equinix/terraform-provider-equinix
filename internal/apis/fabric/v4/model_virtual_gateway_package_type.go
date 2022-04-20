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

// VirtualGatewayPackageType Fabric Gateway Package Type
type VirtualGatewayPackageType struct {
	// Fabric Gateway URI
	Href *string `json:"href,omitempty"`
	// Gateway package code
	Code string `json:"code"`
}

// NewVirtualGatewayPackageType instantiates a new VirtualGatewayPackageType object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewVirtualGatewayPackageType(code string) *VirtualGatewayPackageType {
	this := VirtualGatewayPackageType{}
	this.Code = code
	return &this
}

// NewVirtualGatewayPackageTypeWithDefaults instantiates a new VirtualGatewayPackageType object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewVirtualGatewayPackageTypeWithDefaults() *VirtualGatewayPackageType {
	this := VirtualGatewayPackageType{}
	return &this
}

// GetHref returns the Href field value if set, zero value otherwise.
func (o *VirtualGatewayPackageType) GetHref() string {
	if o == nil || o.Href == nil {
		var ret string
		return ret
	}
	return *o.Href
}

// GetHrefOk returns a tuple with the Href field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VirtualGatewayPackageType) GetHrefOk() (*string, bool) {
	if o == nil || o.Href == nil {
		return nil, false
	}
	return o.Href, true
}

// HasHref returns a boolean if a field has been set.
func (o *VirtualGatewayPackageType) HasHref() bool {
	if o != nil && o.Href != nil {
		return true
	}

	return false
}

// SetHref gets a reference to the given string and assigns it to the Href field.
func (o *VirtualGatewayPackageType) SetHref(v string) {
	o.Href = &v
}

// GetCode returns the Code field value
func (o *VirtualGatewayPackageType) GetCode() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Code
}

// GetCodeOk returns a tuple with the Code field value
// and a boolean to check if the value has been set.
func (o *VirtualGatewayPackageType) GetCodeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Code, true
}

// SetCode sets field value
func (o *VirtualGatewayPackageType) SetCode(v string) {
	o.Code = v
}

func (o VirtualGatewayPackageType) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Href != nil {
		toSerialize["href"] = o.Href
	}
	if true {
		toSerialize["code"] = o.Code
	}
	return json.Marshal(toSerialize)
}

type NullableVirtualGatewayPackageType struct {
	value *VirtualGatewayPackageType
	isSet bool
}

func (v NullableVirtualGatewayPackageType) Get() *VirtualGatewayPackageType {
	return v.value
}

func (v *NullableVirtualGatewayPackageType) Set(val *VirtualGatewayPackageType) {
	v.value = val
	v.isSet = true
}

func (v NullableVirtualGatewayPackageType) IsSet() bool {
	return v.isSet
}

func (v *NullableVirtualGatewayPackageType) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableVirtualGatewayPackageType(val *VirtualGatewayPackageType) *NullableVirtualGatewayPackageType {
	return &NullableVirtualGatewayPackageType{value: val, isSet: true}
}

func (v NullableVirtualGatewayPackageType) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableVirtualGatewayPackageType) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


