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

// SimplifiedLocation struct for SimplifiedLocation
type SimplifiedLocation struct {
	// The Canonical URL at which the resource resides.
	Href *string `json:"href,omitempty"`
	Region *string `json:"region,omitempty"`
	MetroName *string `json:"metroName,omitempty"`
	MetroCode *string `json:"metroCode,omitempty"`
	Ibx *string `json:"ibx,omitempty"`
}

// NewSimplifiedLocation instantiates a new SimplifiedLocation object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSimplifiedLocation() *SimplifiedLocation {
	this := SimplifiedLocation{}
	return &this
}

// NewSimplifiedLocationWithDefaults instantiates a new SimplifiedLocation object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSimplifiedLocationWithDefaults() *SimplifiedLocation {
	this := SimplifiedLocation{}
	return &this
}

// GetHref returns the Href field value if set, zero value otherwise.
func (o *SimplifiedLocation) GetHref() string {
	if o == nil || o.Href == nil {
		var ret string
		return ret
	}
	return *o.Href
}

// GetHrefOk returns a tuple with the Href field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SimplifiedLocation) GetHrefOk() (*string, bool) {
	if o == nil || o.Href == nil {
		return nil, false
	}
	return o.Href, true
}

// HasHref returns a boolean if a field has been set.
func (o *SimplifiedLocation) HasHref() bool {
	if o != nil && o.Href != nil {
		return true
	}

	return false
}

// SetHref gets a reference to the given string and assigns it to the Href field.
func (o *SimplifiedLocation) SetHref(v string) {
	o.Href = &v
}

// GetRegion returns the Region field value if set, zero value otherwise.
func (o *SimplifiedLocation) GetRegion() string {
	if o == nil || o.Region == nil {
		var ret string
		return ret
	}
	return *o.Region
}

// GetRegionOk returns a tuple with the Region field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SimplifiedLocation) GetRegionOk() (*string, bool) {
	if o == nil || o.Region == nil {
		return nil, false
	}
	return o.Region, true
}

// HasRegion returns a boolean if a field has been set.
func (o *SimplifiedLocation) HasRegion() bool {
	if o != nil && o.Region != nil {
		return true
	}

	return false
}

// SetRegion gets a reference to the given string and assigns it to the Region field.
func (o *SimplifiedLocation) SetRegion(v string) {
	o.Region = &v
}

// GetMetroName returns the MetroName field value if set, zero value otherwise.
func (o *SimplifiedLocation) GetMetroName() string {
	if o == nil || o.MetroName == nil {
		var ret string
		return ret
	}
	return *o.MetroName
}

// GetMetroNameOk returns a tuple with the MetroName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SimplifiedLocation) GetMetroNameOk() (*string, bool) {
	if o == nil || o.MetroName == nil {
		return nil, false
	}
	return o.MetroName, true
}

// HasMetroName returns a boolean if a field has been set.
func (o *SimplifiedLocation) HasMetroName() bool {
	if o != nil && o.MetroName != nil {
		return true
	}

	return false
}

// SetMetroName gets a reference to the given string and assigns it to the MetroName field.
func (o *SimplifiedLocation) SetMetroName(v string) {
	o.MetroName = &v
}

// GetMetroCode returns the MetroCode field value if set, zero value otherwise.
func (o *SimplifiedLocation) GetMetroCode() string {
	if o == nil || o.MetroCode == nil {
		var ret string
		return ret
	}
	return *o.MetroCode
}

// GetMetroCodeOk returns a tuple with the MetroCode field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SimplifiedLocation) GetMetroCodeOk() (*string, bool) {
	if o == nil || o.MetroCode == nil {
		return nil, false
	}
	return o.MetroCode, true
}

// HasMetroCode returns a boolean if a field has been set.
func (o *SimplifiedLocation) HasMetroCode() bool {
	if o != nil && o.MetroCode != nil {
		return true
	}

	return false
}

// SetMetroCode gets a reference to the given string and assigns it to the MetroCode field.
func (o *SimplifiedLocation) SetMetroCode(v string) {
	o.MetroCode = &v
}

// GetIbx returns the Ibx field value if set, zero value otherwise.
func (o *SimplifiedLocation) GetIbx() string {
	if o == nil || o.Ibx == nil {
		var ret string
		return ret
	}
	return *o.Ibx
}

// GetIbxOk returns a tuple with the Ibx field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SimplifiedLocation) GetIbxOk() (*string, bool) {
	if o == nil || o.Ibx == nil {
		return nil, false
	}
	return o.Ibx, true
}

// HasIbx returns a boolean if a field has been set.
func (o *SimplifiedLocation) HasIbx() bool {
	if o != nil && o.Ibx != nil {
		return true
	}

	return false
}

// SetIbx gets a reference to the given string and assigns it to the Ibx field.
func (o *SimplifiedLocation) SetIbx(v string) {
	o.Ibx = &v
}

func (o SimplifiedLocation) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Href != nil {
		toSerialize["href"] = o.Href
	}
	if o.Region != nil {
		toSerialize["region"] = o.Region
	}
	if o.MetroName != nil {
		toSerialize["metroName"] = o.MetroName
	}
	if o.MetroCode != nil {
		toSerialize["metroCode"] = o.MetroCode
	}
	if o.Ibx != nil {
		toSerialize["ibx"] = o.Ibx
	}
	return json.Marshal(toSerialize)
}

type NullableSimplifiedLocation struct {
	value *SimplifiedLocation
	isSet bool
}

func (v NullableSimplifiedLocation) Get() *SimplifiedLocation {
	return v.value
}

func (v *NullableSimplifiedLocation) Set(val *SimplifiedLocation) {
	v.value = val
	v.isSet = true
}

func (v NullableSimplifiedLocation) IsSet() bool {
	return v.isSet
}

func (v *NullableSimplifiedLocation) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSimplifiedLocation(val *SimplifiedLocation) *NullableSimplifiedLocation {
	return &NullableSimplifiedLocation{value: val, isSet: true}
}

func (v NullableSimplifiedLocation) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSimplifiedLocation) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


