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

// ServiceProfile Service Profile is a software definition for a named provider service and it's network connectivity requirements. This includes the basic marketing information and one or more sets of access points (a set per each access point type) fulfilling the provider service. 
type ServiceProfile struct {
	State *ServiceProfileStateEnum `json:"state,omitempty"`
	// Seller Account for Service Profile.
	Account *SimplifiedAccount `json:"account,omitempty"`
	Project *Project `json:"project,omitempty"`
	// Seller Account for Service Profile.
	ChangeLog *Changelog `json:"changeLog,omitempty"`
	// Service Profile URI response attribute
	Href *string `json:"href,omitempty"`
	Type ServiceProfileTypeEnum `json:"type"`
	// Customer-assigned service profile name
	Name string `json:"name"`
	// Equinix-assigned service profile identifier
	Uuid *string `json:"uuid,omitempty"`
	// User-provided service description
	Description *string `json:"description,omitempty"`
	// Recipients of notifications on service profile change
	Notifications []SimplifiedNotification `json:"notifications,omitempty"`
	Tags []string `json:"tags,omitempty"`
	Visibility ServiceProfileVisibilityEnum `json:"visibility"`
	AllowedEmails []string `json:"allowedEmails,omitempty"`
	AccessPointTypeConfigs []ServiceProfileAccessPointType `json:"accessPointTypeConfigs"`
	CustomFields []CustomField `json:"customFields,omitempty"`
	MarketingInfo *MarketingInfo `json:"marketingInfo,omitempty"`
	Ports []ServiceProfileAccessPointCOLO `json:"ports,omitempty"`
	VirtualDevices []ServiceProfileAccessPointVD `json:"virtualDevices,omitempty"`
	// Derived response attribute.
	Metros []ServiceMetro `json:"metros,omitempty"`
	// response attribute indicates whether the profile belongs to the same organization as the api-invoker.
	SelfProfile *bool `json:"selfProfile,omitempty"`
}

// NewServiceProfile instantiates a new ServiceProfile object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewServiceProfile(type_ ServiceProfileTypeEnum, name string, visibility ServiceProfileVisibilityEnum, accessPointTypeConfigs []ServiceProfileAccessPointType) *ServiceProfile {
	this := ServiceProfile{}
	this.Type = type_
	this.Name = name
	this.Visibility = visibility
	this.AccessPointTypeConfigs = accessPointTypeConfigs
	return &this
}

// NewServiceProfileWithDefaults instantiates a new ServiceProfile object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewServiceProfileWithDefaults() *ServiceProfile {
	this := ServiceProfile{}
	return &this
}

// GetState returns the State field value if set, zero value otherwise.
func (o *ServiceProfile) GetState() ServiceProfileStateEnum {
	if o == nil || o.State == nil {
		var ret ServiceProfileStateEnum
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServiceProfile) GetStateOk() (*ServiceProfileStateEnum, bool) {
	if o == nil || o.State == nil {
		return nil, false
	}
	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *ServiceProfile) HasState() bool {
	if o != nil && o.State != nil {
		return true
	}

	return false
}

// SetState gets a reference to the given ServiceProfileStateEnum and assigns it to the State field.
func (o *ServiceProfile) SetState(v ServiceProfileStateEnum) {
	o.State = &v
}

// GetAccount returns the Account field value if set, zero value otherwise.
func (o *ServiceProfile) GetAccount() SimplifiedAccount {
	if o == nil || o.Account == nil {
		var ret SimplifiedAccount
		return ret
	}
	return *o.Account
}

// GetAccountOk returns a tuple with the Account field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServiceProfile) GetAccountOk() (*SimplifiedAccount, bool) {
	if o == nil || o.Account == nil {
		return nil, false
	}
	return o.Account, true
}

// HasAccount returns a boolean if a field has been set.
func (o *ServiceProfile) HasAccount() bool {
	if o != nil && o.Account != nil {
		return true
	}

	return false
}

// SetAccount gets a reference to the given SimplifiedAccount and assigns it to the Account field.
func (o *ServiceProfile) SetAccount(v SimplifiedAccount) {
	o.Account = &v
}

// GetProject returns the Project field value if set, zero value otherwise.
func (o *ServiceProfile) GetProject() Project {
	if o == nil || o.Project == nil {
		var ret Project
		return ret
	}
	return *o.Project
}

// GetProjectOk returns a tuple with the Project field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServiceProfile) GetProjectOk() (*Project, bool) {
	if o == nil || o.Project == nil {
		return nil, false
	}
	return o.Project, true
}

// HasProject returns a boolean if a field has been set.
func (o *ServiceProfile) HasProject() bool {
	if o != nil && o.Project != nil {
		return true
	}

	return false
}

// SetProject gets a reference to the given Project and assigns it to the Project field.
func (o *ServiceProfile) SetProject(v Project) {
	o.Project = &v
}

// GetChangeLog returns the ChangeLog field value if set, zero value otherwise.
func (o *ServiceProfile) GetChangeLog() Changelog {
	if o == nil || o.ChangeLog == nil {
		var ret Changelog
		return ret
	}
	return *o.ChangeLog
}

// GetChangeLogOk returns a tuple with the ChangeLog field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServiceProfile) GetChangeLogOk() (*Changelog, bool) {
	if o == nil || o.ChangeLog == nil {
		return nil, false
	}
	return o.ChangeLog, true
}

// HasChangeLog returns a boolean if a field has been set.
func (o *ServiceProfile) HasChangeLog() bool {
	if o != nil && o.ChangeLog != nil {
		return true
	}

	return false
}

// SetChangeLog gets a reference to the given Changelog and assigns it to the ChangeLog field.
func (o *ServiceProfile) SetChangeLog(v Changelog) {
	o.ChangeLog = &v
}

// GetHref returns the Href field value if set, zero value otherwise.
func (o *ServiceProfile) GetHref() string {
	if o == nil || o.Href == nil {
		var ret string
		return ret
	}
	return *o.Href
}

// GetHrefOk returns a tuple with the Href field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServiceProfile) GetHrefOk() (*string, bool) {
	if o == nil || o.Href == nil {
		return nil, false
	}
	return o.Href, true
}

// HasHref returns a boolean if a field has been set.
func (o *ServiceProfile) HasHref() bool {
	if o != nil && o.Href != nil {
		return true
	}

	return false
}

// SetHref gets a reference to the given string and assigns it to the Href field.
func (o *ServiceProfile) SetHref(v string) {
	o.Href = &v
}

// GetType returns the Type field value
func (o *ServiceProfile) GetType() ServiceProfileTypeEnum {
	if o == nil {
		var ret ServiceProfileTypeEnum
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *ServiceProfile) GetTypeOk() (*ServiceProfileTypeEnum, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *ServiceProfile) SetType(v ServiceProfileTypeEnum) {
	o.Type = v
}

// GetName returns the Name field value
func (o *ServiceProfile) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *ServiceProfile) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *ServiceProfile) SetName(v string) {
	o.Name = v
}

// GetUuid returns the Uuid field value if set, zero value otherwise.
func (o *ServiceProfile) GetUuid() string {
	if o == nil || o.Uuid == nil {
		var ret string
		return ret
	}
	return *o.Uuid
}

// GetUuidOk returns a tuple with the Uuid field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServiceProfile) GetUuidOk() (*string, bool) {
	if o == nil || o.Uuid == nil {
		return nil, false
	}
	return o.Uuid, true
}

// HasUuid returns a boolean if a field has been set.
func (o *ServiceProfile) HasUuid() bool {
	if o != nil && o.Uuid != nil {
		return true
	}

	return false
}

// SetUuid gets a reference to the given string and assigns it to the Uuid field.
func (o *ServiceProfile) SetUuid(v string) {
	o.Uuid = &v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *ServiceProfile) GetDescription() string {
	if o == nil || o.Description == nil {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServiceProfile) GetDescriptionOk() (*string, bool) {
	if o == nil || o.Description == nil {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *ServiceProfile) HasDescription() bool {
	if o != nil && o.Description != nil {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *ServiceProfile) SetDescription(v string) {
	o.Description = &v
}

// GetNotifications returns the Notifications field value if set, zero value otherwise.
func (o *ServiceProfile) GetNotifications() []SimplifiedNotification {
	if o == nil || o.Notifications == nil {
		var ret []SimplifiedNotification
		return ret
	}
	return o.Notifications
}

// GetNotificationsOk returns a tuple with the Notifications field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServiceProfile) GetNotificationsOk() ([]SimplifiedNotification, bool) {
	if o == nil || o.Notifications == nil {
		return nil, false
	}
	return o.Notifications, true
}

// HasNotifications returns a boolean if a field has been set.
func (o *ServiceProfile) HasNotifications() bool {
	if o != nil && o.Notifications != nil {
		return true
	}

	return false
}

// SetNotifications gets a reference to the given []SimplifiedNotification and assigns it to the Notifications field.
func (o *ServiceProfile) SetNotifications(v []SimplifiedNotification) {
	o.Notifications = v
}

// GetTags returns the Tags field value if set, zero value otherwise.
func (o *ServiceProfile) GetTags() []string {
	if o == nil || o.Tags == nil {
		var ret []string
		return ret
	}
	return o.Tags
}

// GetTagsOk returns a tuple with the Tags field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServiceProfile) GetTagsOk() ([]string, bool) {
	if o == nil || o.Tags == nil {
		return nil, false
	}
	return o.Tags, true
}

// HasTags returns a boolean if a field has been set.
func (o *ServiceProfile) HasTags() bool {
	if o != nil && o.Tags != nil {
		return true
	}

	return false
}

// SetTags gets a reference to the given []string and assigns it to the Tags field.
func (o *ServiceProfile) SetTags(v []string) {
	o.Tags = v
}

// GetVisibility returns the Visibility field value
func (o *ServiceProfile) GetVisibility() ServiceProfileVisibilityEnum {
	if o == nil {
		var ret ServiceProfileVisibilityEnum
		return ret
	}

	return o.Visibility
}

// GetVisibilityOk returns a tuple with the Visibility field value
// and a boolean to check if the value has been set.
func (o *ServiceProfile) GetVisibilityOk() (*ServiceProfileVisibilityEnum, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Visibility, true
}

// SetVisibility sets field value
func (o *ServiceProfile) SetVisibility(v ServiceProfileVisibilityEnum) {
	o.Visibility = v
}

// GetAllowedEmails returns the AllowedEmails field value if set, zero value otherwise.
func (o *ServiceProfile) GetAllowedEmails() []string {
	if o == nil || o.AllowedEmails == nil {
		var ret []string
		return ret
	}
	return o.AllowedEmails
}

// GetAllowedEmailsOk returns a tuple with the AllowedEmails field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServiceProfile) GetAllowedEmailsOk() ([]string, bool) {
	if o == nil || o.AllowedEmails == nil {
		return nil, false
	}
	return o.AllowedEmails, true
}

// HasAllowedEmails returns a boolean if a field has been set.
func (o *ServiceProfile) HasAllowedEmails() bool {
	if o != nil && o.AllowedEmails != nil {
		return true
	}

	return false
}

// SetAllowedEmails gets a reference to the given []string and assigns it to the AllowedEmails field.
func (o *ServiceProfile) SetAllowedEmails(v []string) {
	o.AllowedEmails = v
}

// GetAccessPointTypeConfigs returns the AccessPointTypeConfigs field value
func (o *ServiceProfile) GetAccessPointTypeConfigs() []ServiceProfileAccessPointType {
	if o == nil {
		var ret []ServiceProfileAccessPointType
		return ret
	}

	return o.AccessPointTypeConfigs
}

// GetAccessPointTypeConfigsOk returns a tuple with the AccessPointTypeConfigs field value
// and a boolean to check if the value has been set.
func (o *ServiceProfile) GetAccessPointTypeConfigsOk() ([]ServiceProfileAccessPointType, bool) {
	if o == nil {
		return nil, false
	}
	return o.AccessPointTypeConfigs, true
}

// SetAccessPointTypeConfigs sets field value
func (o *ServiceProfile) SetAccessPointTypeConfigs(v []ServiceProfileAccessPointType) {
	o.AccessPointTypeConfigs = v
}

// GetCustomFields returns the CustomFields field value if set, zero value otherwise.
func (o *ServiceProfile) GetCustomFields() []CustomField {
	if o == nil || o.CustomFields == nil {
		var ret []CustomField
		return ret
	}
	return o.CustomFields
}

// GetCustomFieldsOk returns a tuple with the CustomFields field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServiceProfile) GetCustomFieldsOk() ([]CustomField, bool) {
	if o == nil || o.CustomFields == nil {
		return nil, false
	}
	return o.CustomFields, true
}

// HasCustomFields returns a boolean if a field has been set.
func (o *ServiceProfile) HasCustomFields() bool {
	if o != nil && o.CustomFields != nil {
		return true
	}

	return false
}

// SetCustomFields gets a reference to the given []CustomField and assigns it to the CustomFields field.
func (o *ServiceProfile) SetCustomFields(v []CustomField) {
	o.CustomFields = v
}

// GetMarketingInfo returns the MarketingInfo field value if set, zero value otherwise.
func (o *ServiceProfile) GetMarketingInfo() MarketingInfo {
	if o == nil || o.MarketingInfo == nil {
		var ret MarketingInfo
		return ret
	}
	return *o.MarketingInfo
}

// GetMarketingInfoOk returns a tuple with the MarketingInfo field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServiceProfile) GetMarketingInfoOk() (*MarketingInfo, bool) {
	if o == nil || o.MarketingInfo == nil {
		return nil, false
	}
	return o.MarketingInfo, true
}

// HasMarketingInfo returns a boolean if a field has been set.
func (o *ServiceProfile) HasMarketingInfo() bool {
	if o != nil && o.MarketingInfo != nil {
		return true
	}

	return false
}

// SetMarketingInfo gets a reference to the given MarketingInfo and assigns it to the MarketingInfo field.
func (o *ServiceProfile) SetMarketingInfo(v MarketingInfo) {
	o.MarketingInfo = &v
}

// GetPorts returns the Ports field value if set, zero value otherwise.
func (o *ServiceProfile) GetPorts() []ServiceProfileAccessPointCOLO {
	if o == nil || o.Ports == nil {
		var ret []ServiceProfileAccessPointCOLO
		return ret
	}
	return o.Ports
}

// GetPortsOk returns a tuple with the Ports field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServiceProfile) GetPortsOk() ([]ServiceProfileAccessPointCOLO, bool) {
	if o == nil || o.Ports == nil {
		return nil, false
	}
	return o.Ports, true
}

// HasPorts returns a boolean if a field has been set.
func (o *ServiceProfile) HasPorts() bool {
	if o != nil && o.Ports != nil {
		return true
	}

	return false
}

// SetPorts gets a reference to the given []ServiceProfileAccessPointCOLO and assigns it to the Ports field.
func (o *ServiceProfile) SetPorts(v []ServiceProfileAccessPointCOLO) {
	o.Ports = v
}

// GetVirtualDevices returns the VirtualDevices field value if set, zero value otherwise.
func (o *ServiceProfile) GetVirtualDevices() []ServiceProfileAccessPointVD {
	if o == nil || o.VirtualDevices == nil {
		var ret []ServiceProfileAccessPointVD
		return ret
	}
	return o.VirtualDevices
}

// GetVirtualDevicesOk returns a tuple with the VirtualDevices field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServiceProfile) GetVirtualDevicesOk() ([]ServiceProfileAccessPointVD, bool) {
	if o == nil || o.VirtualDevices == nil {
		return nil, false
	}
	return o.VirtualDevices, true
}

// HasVirtualDevices returns a boolean if a field has been set.
func (o *ServiceProfile) HasVirtualDevices() bool {
	if o != nil && o.VirtualDevices != nil {
		return true
	}

	return false
}

// SetVirtualDevices gets a reference to the given []ServiceProfileAccessPointVD and assigns it to the VirtualDevices field.
func (o *ServiceProfile) SetVirtualDevices(v []ServiceProfileAccessPointVD) {
	o.VirtualDevices = v
}

// GetMetros returns the Metros field value if set, zero value otherwise.
func (o *ServiceProfile) GetMetros() []ServiceMetro {
	if o == nil || o.Metros == nil {
		var ret []ServiceMetro
		return ret
	}
	return o.Metros
}

// GetMetrosOk returns a tuple with the Metros field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServiceProfile) GetMetrosOk() ([]ServiceMetro, bool) {
	if o == nil || o.Metros == nil {
		return nil, false
	}
	return o.Metros, true
}

// HasMetros returns a boolean if a field has been set.
func (o *ServiceProfile) HasMetros() bool {
	if o != nil && o.Metros != nil {
		return true
	}

	return false
}

// SetMetros gets a reference to the given []ServiceMetro and assigns it to the Metros field.
func (o *ServiceProfile) SetMetros(v []ServiceMetro) {
	o.Metros = v
}

// GetSelfProfile returns the SelfProfile field value if set, zero value otherwise.
func (o *ServiceProfile) GetSelfProfile() bool {
	if o == nil || o.SelfProfile == nil {
		var ret bool
		return ret
	}
	return *o.SelfProfile
}

// GetSelfProfileOk returns a tuple with the SelfProfile field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServiceProfile) GetSelfProfileOk() (*bool, bool) {
	if o == nil || o.SelfProfile == nil {
		return nil, false
	}
	return o.SelfProfile, true
}

// HasSelfProfile returns a boolean if a field has been set.
func (o *ServiceProfile) HasSelfProfile() bool {
	if o != nil && o.SelfProfile != nil {
		return true
	}

	return false
}

// SetSelfProfile gets a reference to the given bool and assigns it to the SelfProfile field.
func (o *ServiceProfile) SetSelfProfile(v bool) {
	o.SelfProfile = &v
}

func (o ServiceProfile) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.State != nil {
		toSerialize["state"] = o.State
	}
	if o.Account != nil {
		toSerialize["account"] = o.Account
	}
	if o.Project != nil {
		toSerialize["project"] = o.Project
	}
	if o.ChangeLog != nil {
		toSerialize["changeLog"] = o.ChangeLog
	}
	if o.Href != nil {
		toSerialize["href"] = o.Href
	}
	if true {
		toSerialize["type"] = o.Type
	}
	if true {
		toSerialize["name"] = o.Name
	}
	if o.Uuid != nil {
		toSerialize["uuid"] = o.Uuid
	}
	if o.Description != nil {
		toSerialize["description"] = o.Description
	}
	if o.Notifications != nil {
		toSerialize["notifications"] = o.Notifications
	}
	if o.Tags != nil {
		toSerialize["tags"] = o.Tags
	}
	if true {
		toSerialize["visibility"] = o.Visibility
	}
	if o.AllowedEmails != nil {
		toSerialize["allowedEmails"] = o.AllowedEmails
	}
	if true {
		toSerialize["accessPointTypeConfigs"] = o.AccessPointTypeConfigs
	}
	if o.CustomFields != nil {
		toSerialize["customFields"] = o.CustomFields
	}
	if o.MarketingInfo != nil {
		toSerialize["marketingInfo"] = o.MarketingInfo
	}
	if o.Ports != nil {
		toSerialize["ports"] = o.Ports
	}
	if o.VirtualDevices != nil {
		toSerialize["virtualDevices"] = o.VirtualDevices
	}
	if o.Metros != nil {
		toSerialize["metros"] = o.Metros
	}
	if o.SelfProfile != nil {
		toSerialize["selfProfile"] = o.SelfProfile
	}
	return json.Marshal(toSerialize)
}

type NullableServiceProfile struct {
	value *ServiceProfile
	isSet bool
}

func (v NullableServiceProfile) Get() *ServiceProfile {
	return v.value
}

func (v *NullableServiceProfile) Set(val *ServiceProfile) {
	v.value = val
	v.isSet = true
}

func (v NullableServiceProfile) IsSet() bool {
	return v.isSet
}

func (v *NullableServiceProfile) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableServiceProfile(val *ServiceProfile) *NullableServiceProfile {
	return &NullableServiceProfile{value: val, isSet: true}
}

func (v NullableServiceProfile) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableServiceProfile) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


