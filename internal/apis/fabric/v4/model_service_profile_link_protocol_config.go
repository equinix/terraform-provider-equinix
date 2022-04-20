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

// ServiceProfileLinkProtocolConfig Configuration for dot1q to qinq translation support
type ServiceProfileLinkProtocolConfig struct {
	// was tagType - missing on wiki
	EncapsulationStrategy *string `json:"encapsulationStrategy,omitempty"`
	NamedTags []string `json:"namedTags,omitempty"`
	// was ctagLabel
	VlanCTagLabel *string `json:"vlanCTagLabel,omitempty"`
	ReuseVlanSTag *bool `json:"reuseVlanSTag,omitempty"`
	// Port encapsulation - Derived response attribute. Ignored on request payloads.
	Encapsulation *string `json:"encapsulation,omitempty"`
}

// NewServiceProfileLinkProtocolConfig instantiates a new ServiceProfileLinkProtocolConfig object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewServiceProfileLinkProtocolConfig() *ServiceProfileLinkProtocolConfig {
	this := ServiceProfileLinkProtocolConfig{}
	var reuseVlanSTag bool = false
	this.ReuseVlanSTag = &reuseVlanSTag
	return &this
}

// NewServiceProfileLinkProtocolConfigWithDefaults instantiates a new ServiceProfileLinkProtocolConfig object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewServiceProfileLinkProtocolConfigWithDefaults() *ServiceProfileLinkProtocolConfig {
	this := ServiceProfileLinkProtocolConfig{}
	var reuseVlanSTag bool = false
	this.ReuseVlanSTag = &reuseVlanSTag
	return &this
}

// GetEncapsulationStrategy returns the EncapsulationStrategy field value if set, zero value otherwise.
func (o *ServiceProfileLinkProtocolConfig) GetEncapsulationStrategy() string {
	if o == nil || o.EncapsulationStrategy == nil {
		var ret string
		return ret
	}
	return *o.EncapsulationStrategy
}

// GetEncapsulationStrategyOk returns a tuple with the EncapsulationStrategy field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServiceProfileLinkProtocolConfig) GetEncapsulationStrategyOk() (*string, bool) {
	if o == nil || o.EncapsulationStrategy == nil {
		return nil, false
	}
	return o.EncapsulationStrategy, true
}

// HasEncapsulationStrategy returns a boolean if a field has been set.
func (o *ServiceProfileLinkProtocolConfig) HasEncapsulationStrategy() bool {
	if o != nil && o.EncapsulationStrategy != nil {
		return true
	}

	return false
}

// SetEncapsulationStrategy gets a reference to the given string and assigns it to the EncapsulationStrategy field.
func (o *ServiceProfileLinkProtocolConfig) SetEncapsulationStrategy(v string) {
	o.EncapsulationStrategy = &v
}

// GetNamedTags returns the NamedTags field value if set, zero value otherwise.
func (o *ServiceProfileLinkProtocolConfig) GetNamedTags() []string {
	if o == nil || o.NamedTags == nil {
		var ret []string
		return ret
	}
	return o.NamedTags
}

// GetNamedTagsOk returns a tuple with the NamedTags field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServiceProfileLinkProtocolConfig) GetNamedTagsOk() ([]string, bool) {
	if o == nil || o.NamedTags == nil {
		return nil, false
	}
	return o.NamedTags, true
}

// HasNamedTags returns a boolean if a field has been set.
func (o *ServiceProfileLinkProtocolConfig) HasNamedTags() bool {
	if o != nil && o.NamedTags != nil {
		return true
	}

	return false
}

// SetNamedTags gets a reference to the given []string and assigns it to the NamedTags field.
func (o *ServiceProfileLinkProtocolConfig) SetNamedTags(v []string) {
	o.NamedTags = v
}

// GetVlanCTagLabel returns the VlanCTagLabel field value if set, zero value otherwise.
func (o *ServiceProfileLinkProtocolConfig) GetVlanCTagLabel() string {
	if o == nil || o.VlanCTagLabel == nil {
		var ret string
		return ret
	}
	return *o.VlanCTagLabel
}

// GetVlanCTagLabelOk returns a tuple with the VlanCTagLabel field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServiceProfileLinkProtocolConfig) GetVlanCTagLabelOk() (*string, bool) {
	if o == nil || o.VlanCTagLabel == nil {
		return nil, false
	}
	return o.VlanCTagLabel, true
}

// HasVlanCTagLabel returns a boolean if a field has been set.
func (o *ServiceProfileLinkProtocolConfig) HasVlanCTagLabel() bool {
	if o != nil && o.VlanCTagLabel != nil {
		return true
	}

	return false
}

// SetVlanCTagLabel gets a reference to the given string and assigns it to the VlanCTagLabel field.
func (o *ServiceProfileLinkProtocolConfig) SetVlanCTagLabel(v string) {
	o.VlanCTagLabel = &v
}

// GetReuseVlanSTag returns the ReuseVlanSTag field value if set, zero value otherwise.
func (o *ServiceProfileLinkProtocolConfig) GetReuseVlanSTag() bool {
	if o == nil || o.ReuseVlanSTag == nil {
		var ret bool
		return ret
	}
	return *o.ReuseVlanSTag
}

// GetReuseVlanSTagOk returns a tuple with the ReuseVlanSTag field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServiceProfileLinkProtocolConfig) GetReuseVlanSTagOk() (*bool, bool) {
	if o == nil || o.ReuseVlanSTag == nil {
		return nil, false
	}
	return o.ReuseVlanSTag, true
}

// HasReuseVlanSTag returns a boolean if a field has been set.
func (o *ServiceProfileLinkProtocolConfig) HasReuseVlanSTag() bool {
	if o != nil && o.ReuseVlanSTag != nil {
		return true
	}

	return false
}

// SetReuseVlanSTag gets a reference to the given bool and assigns it to the ReuseVlanSTag field.
func (o *ServiceProfileLinkProtocolConfig) SetReuseVlanSTag(v bool) {
	o.ReuseVlanSTag = &v
}

// GetEncapsulation returns the Encapsulation field value if set, zero value otherwise.
func (o *ServiceProfileLinkProtocolConfig) GetEncapsulation() string {
	if o == nil || o.Encapsulation == nil {
		var ret string
		return ret
	}
	return *o.Encapsulation
}

// GetEncapsulationOk returns a tuple with the Encapsulation field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServiceProfileLinkProtocolConfig) GetEncapsulationOk() (*string, bool) {
	if o == nil || o.Encapsulation == nil {
		return nil, false
	}
	return o.Encapsulation, true
}

// HasEncapsulation returns a boolean if a field has been set.
func (o *ServiceProfileLinkProtocolConfig) HasEncapsulation() bool {
	if o != nil && o.Encapsulation != nil {
		return true
	}

	return false
}

// SetEncapsulation gets a reference to the given string and assigns it to the Encapsulation field.
func (o *ServiceProfileLinkProtocolConfig) SetEncapsulation(v string) {
	o.Encapsulation = &v
}

func (o ServiceProfileLinkProtocolConfig) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.EncapsulationStrategy != nil {
		toSerialize["encapsulationStrategy"] = o.EncapsulationStrategy
	}
	if o.NamedTags != nil {
		toSerialize["namedTags"] = o.NamedTags
	}
	if o.VlanCTagLabel != nil {
		toSerialize["vlanCTagLabel"] = o.VlanCTagLabel
	}
	if o.ReuseVlanSTag != nil {
		toSerialize["reuseVlanSTag"] = o.ReuseVlanSTag
	}
	if o.Encapsulation != nil {
		toSerialize["encapsulation"] = o.Encapsulation
	}
	return json.Marshal(toSerialize)
}

type NullableServiceProfileLinkProtocolConfig struct {
	value *ServiceProfileLinkProtocolConfig
	isSet bool
}

func (v NullableServiceProfileLinkProtocolConfig) Get() *ServiceProfileLinkProtocolConfig {
	return v.value
}

func (v *NullableServiceProfileLinkProtocolConfig) Set(val *ServiceProfileLinkProtocolConfig) {
	v.value = val
	v.isSet = true
}

func (v NullableServiceProfileLinkProtocolConfig) IsSet() bool {
	return v.isSet
}

func (v *NullableServiceProfileLinkProtocolConfig) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableServiceProfileLinkProtocolConfig(val *ServiceProfileLinkProtocolConfig) *NullableServiceProfileLinkProtocolConfig {
	return &NullableServiceProfileLinkProtocolConfig{value: val, isSet: true}
}

func (v NullableServiceProfileLinkProtocolConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableServiceProfileLinkProtocolConfig) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


