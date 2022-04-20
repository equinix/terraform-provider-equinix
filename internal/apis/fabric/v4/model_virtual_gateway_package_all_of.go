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

// VirtualGatewayPackageAllOf struct for VirtualGatewayPackageAllOf
type VirtualGatewayPackageAllOf struct {
	// Gateway package URI
	Href *string `json:"href,omitempty"`
	// Type of Gateway package
	Type *string `json:"type,omitempty"`
	Code *Code `json:"code,omitempty"`
	// Fabric Gateway Package description
	Description *string `json:"description,omitempty"`
	// Gateway package BGP IPv4 routes limit
	TotalIPv4RoutesMax *int32 `json:"totalIPv4RoutesMax,omitempty"`
	// Gateway package BGP IPv6 routes limit
	TotalIPv6RoutesMax *int32 `json:"totalIPv6RoutesMax,omitempty"`
	// Gateway package static IPv4 routes limit
	StaticIPv4RoutesMax *int32 `json:"staticIPv4RoutesMax,omitempty"`
	// Gateway package static IPv6 routes limit
	StaticIPv6RoutesMax *int32 `json:"staticIPv6RoutesMax,omitempty"`
	// Gateway package ACLs limit
	AclMax *int32 `json:"aclMax,omitempty"`
	// Gateway package ACL rules limit
	AclRuleMax *int32 `json:"aclRuleMax,omitempty"`
	// Gateway package high-available configuration support
	IsHaSupported *bool `json:"isHaSupported,omitempty"`
	// Gateway package route filter support
	IsRouteFilterSupported *bool `json:"isRouteFilterSupported,omitempty"`
	Nat *Nat `json:"nat,omitempty"`
	ChangeLog *PackageChangeLog `json:"changeLog,omitempty"`
}

// NewVirtualGatewayPackageAllOf instantiates a new VirtualGatewayPackageAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewVirtualGatewayPackageAllOf() *VirtualGatewayPackageAllOf {
	this := VirtualGatewayPackageAllOf{}
	return &this
}

// NewVirtualGatewayPackageAllOfWithDefaults instantiates a new VirtualGatewayPackageAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewVirtualGatewayPackageAllOfWithDefaults() *VirtualGatewayPackageAllOf {
	this := VirtualGatewayPackageAllOf{}
	return &this
}

// GetHref returns the Href field value if set, zero value otherwise.
func (o *VirtualGatewayPackageAllOf) GetHref() string {
	if o == nil || o.Href == nil {
		var ret string
		return ret
	}
	return *o.Href
}

// GetHrefOk returns a tuple with the Href field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VirtualGatewayPackageAllOf) GetHrefOk() (*string, bool) {
	if o == nil || o.Href == nil {
		return nil, false
	}
	return o.Href, true
}

// HasHref returns a boolean if a field has been set.
func (o *VirtualGatewayPackageAllOf) HasHref() bool {
	if o != nil && o.Href != nil {
		return true
	}

	return false
}

// SetHref gets a reference to the given string and assigns it to the Href field.
func (o *VirtualGatewayPackageAllOf) SetHref(v string) {
	o.Href = &v
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *VirtualGatewayPackageAllOf) GetType() string {
	if o == nil || o.Type == nil {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VirtualGatewayPackageAllOf) GetTypeOk() (*string, bool) {
	if o == nil || o.Type == nil {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *VirtualGatewayPackageAllOf) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *VirtualGatewayPackageAllOf) SetType(v string) {
	o.Type = &v
}

// GetCode returns the Code field value if set, zero value otherwise.
func (o *VirtualGatewayPackageAllOf) GetCode() Code {
	if o == nil || o.Code == nil {
		var ret Code
		return ret
	}
	return *o.Code
}

// GetCodeOk returns a tuple with the Code field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VirtualGatewayPackageAllOf) GetCodeOk() (*Code, bool) {
	if o == nil || o.Code == nil {
		return nil, false
	}
	return o.Code, true
}

// HasCode returns a boolean if a field has been set.
func (o *VirtualGatewayPackageAllOf) HasCode() bool {
	if o != nil && o.Code != nil {
		return true
	}

	return false
}

// SetCode gets a reference to the given Code and assigns it to the Code field.
func (o *VirtualGatewayPackageAllOf) SetCode(v Code) {
	o.Code = &v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *VirtualGatewayPackageAllOf) GetDescription() string {
	if o == nil || o.Description == nil {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VirtualGatewayPackageAllOf) GetDescriptionOk() (*string, bool) {
	if o == nil || o.Description == nil {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *VirtualGatewayPackageAllOf) HasDescription() bool {
	if o != nil && o.Description != nil {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *VirtualGatewayPackageAllOf) SetDescription(v string) {
	o.Description = &v
}

// GetTotalIPv4RoutesMax returns the TotalIPv4RoutesMax field value if set, zero value otherwise.
func (o *VirtualGatewayPackageAllOf) GetTotalIPv4RoutesMax() int32 {
	if o == nil || o.TotalIPv4RoutesMax == nil {
		var ret int32
		return ret
	}
	return *o.TotalIPv4RoutesMax
}

// GetTotalIPv4RoutesMaxOk returns a tuple with the TotalIPv4RoutesMax field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VirtualGatewayPackageAllOf) GetTotalIPv4RoutesMaxOk() (*int32, bool) {
	if o == nil || o.TotalIPv4RoutesMax == nil {
		return nil, false
	}
	return o.TotalIPv4RoutesMax, true
}

// HasTotalIPv4RoutesMax returns a boolean if a field has been set.
func (o *VirtualGatewayPackageAllOf) HasTotalIPv4RoutesMax() bool {
	if o != nil && o.TotalIPv4RoutesMax != nil {
		return true
	}

	return false
}

// SetTotalIPv4RoutesMax gets a reference to the given int32 and assigns it to the TotalIPv4RoutesMax field.
func (o *VirtualGatewayPackageAllOf) SetTotalIPv4RoutesMax(v int32) {
	o.TotalIPv4RoutesMax = &v
}

// GetTotalIPv6RoutesMax returns the TotalIPv6RoutesMax field value if set, zero value otherwise.
func (o *VirtualGatewayPackageAllOf) GetTotalIPv6RoutesMax() int32 {
	if o == nil || o.TotalIPv6RoutesMax == nil {
		var ret int32
		return ret
	}
	return *o.TotalIPv6RoutesMax
}

// GetTotalIPv6RoutesMaxOk returns a tuple with the TotalIPv6RoutesMax field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VirtualGatewayPackageAllOf) GetTotalIPv6RoutesMaxOk() (*int32, bool) {
	if o == nil || o.TotalIPv6RoutesMax == nil {
		return nil, false
	}
	return o.TotalIPv6RoutesMax, true
}

// HasTotalIPv6RoutesMax returns a boolean if a field has been set.
func (o *VirtualGatewayPackageAllOf) HasTotalIPv6RoutesMax() bool {
	if o != nil && o.TotalIPv6RoutesMax != nil {
		return true
	}

	return false
}

// SetTotalIPv6RoutesMax gets a reference to the given int32 and assigns it to the TotalIPv6RoutesMax field.
func (o *VirtualGatewayPackageAllOf) SetTotalIPv6RoutesMax(v int32) {
	o.TotalIPv6RoutesMax = &v
}

// GetStaticIPv4RoutesMax returns the StaticIPv4RoutesMax field value if set, zero value otherwise.
func (o *VirtualGatewayPackageAllOf) GetStaticIPv4RoutesMax() int32 {
	if o == nil || o.StaticIPv4RoutesMax == nil {
		var ret int32
		return ret
	}
	return *o.StaticIPv4RoutesMax
}

// GetStaticIPv4RoutesMaxOk returns a tuple with the StaticIPv4RoutesMax field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VirtualGatewayPackageAllOf) GetStaticIPv4RoutesMaxOk() (*int32, bool) {
	if o == nil || o.StaticIPv4RoutesMax == nil {
		return nil, false
	}
	return o.StaticIPv4RoutesMax, true
}

// HasStaticIPv4RoutesMax returns a boolean if a field has been set.
func (o *VirtualGatewayPackageAllOf) HasStaticIPv4RoutesMax() bool {
	if o != nil && o.StaticIPv4RoutesMax != nil {
		return true
	}

	return false
}

// SetStaticIPv4RoutesMax gets a reference to the given int32 and assigns it to the StaticIPv4RoutesMax field.
func (o *VirtualGatewayPackageAllOf) SetStaticIPv4RoutesMax(v int32) {
	o.StaticIPv4RoutesMax = &v
}

// GetStaticIPv6RoutesMax returns the StaticIPv6RoutesMax field value if set, zero value otherwise.
func (o *VirtualGatewayPackageAllOf) GetStaticIPv6RoutesMax() int32 {
	if o == nil || o.StaticIPv6RoutesMax == nil {
		var ret int32
		return ret
	}
	return *o.StaticIPv6RoutesMax
}

// GetStaticIPv6RoutesMaxOk returns a tuple with the StaticIPv6RoutesMax field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VirtualGatewayPackageAllOf) GetStaticIPv6RoutesMaxOk() (*int32, bool) {
	if o == nil || o.StaticIPv6RoutesMax == nil {
		return nil, false
	}
	return o.StaticIPv6RoutesMax, true
}

// HasStaticIPv6RoutesMax returns a boolean if a field has been set.
func (o *VirtualGatewayPackageAllOf) HasStaticIPv6RoutesMax() bool {
	if o != nil && o.StaticIPv6RoutesMax != nil {
		return true
	}

	return false
}

// SetStaticIPv6RoutesMax gets a reference to the given int32 and assigns it to the StaticIPv6RoutesMax field.
func (o *VirtualGatewayPackageAllOf) SetStaticIPv6RoutesMax(v int32) {
	o.StaticIPv6RoutesMax = &v
}

// GetAclMax returns the AclMax field value if set, zero value otherwise.
func (o *VirtualGatewayPackageAllOf) GetAclMax() int32 {
	if o == nil || o.AclMax == nil {
		var ret int32
		return ret
	}
	return *o.AclMax
}

// GetAclMaxOk returns a tuple with the AclMax field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VirtualGatewayPackageAllOf) GetAclMaxOk() (*int32, bool) {
	if o == nil || o.AclMax == nil {
		return nil, false
	}
	return o.AclMax, true
}

// HasAclMax returns a boolean if a field has been set.
func (o *VirtualGatewayPackageAllOf) HasAclMax() bool {
	if o != nil && o.AclMax != nil {
		return true
	}

	return false
}

// SetAclMax gets a reference to the given int32 and assigns it to the AclMax field.
func (o *VirtualGatewayPackageAllOf) SetAclMax(v int32) {
	o.AclMax = &v
}

// GetAclRuleMax returns the AclRuleMax field value if set, zero value otherwise.
func (o *VirtualGatewayPackageAllOf) GetAclRuleMax() int32 {
	if o == nil || o.AclRuleMax == nil {
		var ret int32
		return ret
	}
	return *o.AclRuleMax
}

// GetAclRuleMaxOk returns a tuple with the AclRuleMax field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VirtualGatewayPackageAllOf) GetAclRuleMaxOk() (*int32, bool) {
	if o == nil || o.AclRuleMax == nil {
		return nil, false
	}
	return o.AclRuleMax, true
}

// HasAclRuleMax returns a boolean if a field has been set.
func (o *VirtualGatewayPackageAllOf) HasAclRuleMax() bool {
	if o != nil && o.AclRuleMax != nil {
		return true
	}

	return false
}

// SetAclRuleMax gets a reference to the given int32 and assigns it to the AclRuleMax field.
func (o *VirtualGatewayPackageAllOf) SetAclRuleMax(v int32) {
	o.AclRuleMax = &v
}

// GetIsHaSupported returns the IsHaSupported field value if set, zero value otherwise.
func (o *VirtualGatewayPackageAllOf) GetIsHaSupported() bool {
	if o == nil || o.IsHaSupported == nil {
		var ret bool
		return ret
	}
	return *o.IsHaSupported
}

// GetIsHaSupportedOk returns a tuple with the IsHaSupported field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VirtualGatewayPackageAllOf) GetIsHaSupportedOk() (*bool, bool) {
	if o == nil || o.IsHaSupported == nil {
		return nil, false
	}
	return o.IsHaSupported, true
}

// HasIsHaSupported returns a boolean if a field has been set.
func (o *VirtualGatewayPackageAllOf) HasIsHaSupported() bool {
	if o != nil && o.IsHaSupported != nil {
		return true
	}

	return false
}

// SetIsHaSupported gets a reference to the given bool and assigns it to the IsHaSupported field.
func (o *VirtualGatewayPackageAllOf) SetIsHaSupported(v bool) {
	o.IsHaSupported = &v
}

// GetIsRouteFilterSupported returns the IsRouteFilterSupported field value if set, zero value otherwise.
func (o *VirtualGatewayPackageAllOf) GetIsRouteFilterSupported() bool {
	if o == nil || o.IsRouteFilterSupported == nil {
		var ret bool
		return ret
	}
	return *o.IsRouteFilterSupported
}

// GetIsRouteFilterSupportedOk returns a tuple with the IsRouteFilterSupported field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VirtualGatewayPackageAllOf) GetIsRouteFilterSupportedOk() (*bool, bool) {
	if o == nil || o.IsRouteFilterSupported == nil {
		return nil, false
	}
	return o.IsRouteFilterSupported, true
}

// HasIsRouteFilterSupported returns a boolean if a field has been set.
func (o *VirtualGatewayPackageAllOf) HasIsRouteFilterSupported() bool {
	if o != nil && o.IsRouteFilterSupported != nil {
		return true
	}

	return false
}

// SetIsRouteFilterSupported gets a reference to the given bool and assigns it to the IsRouteFilterSupported field.
func (o *VirtualGatewayPackageAllOf) SetIsRouteFilterSupported(v bool) {
	o.IsRouteFilterSupported = &v
}

// GetNat returns the Nat field value if set, zero value otherwise.
func (o *VirtualGatewayPackageAllOf) GetNat() Nat {
	if o == nil || o.Nat == nil {
		var ret Nat
		return ret
	}
	return *o.Nat
}

// GetNatOk returns a tuple with the Nat field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VirtualGatewayPackageAllOf) GetNatOk() (*Nat, bool) {
	if o == nil || o.Nat == nil {
		return nil, false
	}
	return o.Nat, true
}

// HasNat returns a boolean if a field has been set.
func (o *VirtualGatewayPackageAllOf) HasNat() bool {
	if o != nil && o.Nat != nil {
		return true
	}

	return false
}

// SetNat gets a reference to the given Nat and assigns it to the Nat field.
func (o *VirtualGatewayPackageAllOf) SetNat(v Nat) {
	o.Nat = &v
}

// GetChangeLog returns the ChangeLog field value if set, zero value otherwise.
func (o *VirtualGatewayPackageAllOf) GetChangeLog() PackageChangeLog {
	if o == nil || o.ChangeLog == nil {
		var ret PackageChangeLog
		return ret
	}
	return *o.ChangeLog
}

// GetChangeLogOk returns a tuple with the ChangeLog field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VirtualGatewayPackageAllOf) GetChangeLogOk() (*PackageChangeLog, bool) {
	if o == nil || o.ChangeLog == nil {
		return nil, false
	}
	return o.ChangeLog, true
}

// HasChangeLog returns a boolean if a field has been set.
func (o *VirtualGatewayPackageAllOf) HasChangeLog() bool {
	if o != nil && o.ChangeLog != nil {
		return true
	}

	return false
}

// SetChangeLog gets a reference to the given PackageChangeLog and assigns it to the ChangeLog field.
func (o *VirtualGatewayPackageAllOf) SetChangeLog(v PackageChangeLog) {
	o.ChangeLog = &v
}

func (o VirtualGatewayPackageAllOf) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Href != nil {
		toSerialize["href"] = o.Href
	}
	if o.Type != nil {
		toSerialize["type"] = o.Type
	}
	if o.Code != nil {
		toSerialize["code"] = o.Code
	}
	if o.Description != nil {
		toSerialize["description"] = o.Description
	}
	if o.TotalIPv4RoutesMax != nil {
		toSerialize["totalIPv4RoutesMax"] = o.TotalIPv4RoutesMax
	}
	if o.TotalIPv6RoutesMax != nil {
		toSerialize["totalIPv6RoutesMax"] = o.TotalIPv6RoutesMax
	}
	if o.StaticIPv4RoutesMax != nil {
		toSerialize["staticIPv4RoutesMax"] = o.StaticIPv4RoutesMax
	}
	if o.StaticIPv6RoutesMax != nil {
		toSerialize["staticIPv6RoutesMax"] = o.StaticIPv6RoutesMax
	}
	if o.AclMax != nil {
		toSerialize["aclMax"] = o.AclMax
	}
	if o.AclRuleMax != nil {
		toSerialize["aclRuleMax"] = o.AclRuleMax
	}
	if o.IsHaSupported != nil {
		toSerialize["isHaSupported"] = o.IsHaSupported
	}
	if o.IsRouteFilterSupported != nil {
		toSerialize["isRouteFilterSupported"] = o.IsRouteFilterSupported
	}
	if o.Nat != nil {
		toSerialize["nat"] = o.Nat
	}
	if o.ChangeLog != nil {
		toSerialize["changeLog"] = o.ChangeLog
	}
	return json.Marshal(toSerialize)
}

type NullableVirtualGatewayPackageAllOf struct {
	value *VirtualGatewayPackageAllOf
	isSet bool
}

func (v NullableVirtualGatewayPackageAllOf) Get() *VirtualGatewayPackageAllOf {
	return v.value
}

func (v *NullableVirtualGatewayPackageAllOf) Set(val *VirtualGatewayPackageAllOf) {
	v.value = val
	v.isSet = true
}

func (v NullableVirtualGatewayPackageAllOf) IsSet() bool {
	return v.isSet
}

func (v *NullableVirtualGatewayPackageAllOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableVirtualGatewayPackageAllOf(val *VirtualGatewayPackageAllOf) *NullableVirtualGatewayPackageAllOf {
	return &NullableVirtualGatewayPackageAllOf{value: val, isSet: true}
}

func (v NullableVirtualGatewayPackageAllOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableVirtualGatewayPackageAllOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


