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

// VirtualGatewayPostRequest Create Fabric Gateway
type VirtualGatewayPostRequest struct {
	Type string `json:"type"`
	// Customer-provided Fabric Gateway name
	Name string `json:"name"`
	Location SimplifiedLocationWithoutIBX `json:"location"`
	Package VirtualGatewayPackageType `json:"package"`
	Order *Order `json:"order,omitempty"`
	Project *Project `json:"project,omitempty"`
	Account SimplifiedAccount `json:"account"`
	// Preferences for notifications on connection configuration or status changes
	Notifications []SimplifiedNotification `json:"notifications"`
}

// NewVirtualGatewayPostRequest instantiates a new VirtualGatewayPostRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewVirtualGatewayPostRequest(type_ string, name string, location SimplifiedLocationWithoutIBX, package_ VirtualGatewayPackageType, account SimplifiedAccount, notifications []SimplifiedNotification) *VirtualGatewayPostRequest {
	this := VirtualGatewayPostRequest{}
	this.Type = type_
	this.Name = name
	this.Location = location
	this.Package = package_
	this.Account = account
	this.Notifications = notifications
	return &this
}

// NewVirtualGatewayPostRequestWithDefaults instantiates a new VirtualGatewayPostRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewVirtualGatewayPostRequestWithDefaults() *VirtualGatewayPostRequest {
	this := VirtualGatewayPostRequest{}
	return &this
}

// GetType returns the Type field value
func (o *VirtualGatewayPostRequest) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *VirtualGatewayPostRequest) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *VirtualGatewayPostRequest) SetType(v string) {
	o.Type = v
}

// GetName returns the Name field value
func (o *VirtualGatewayPostRequest) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *VirtualGatewayPostRequest) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *VirtualGatewayPostRequest) SetName(v string) {
	o.Name = v
}

// GetLocation returns the Location field value
func (o *VirtualGatewayPostRequest) GetLocation() SimplifiedLocationWithoutIBX {
	if o == nil {
		var ret SimplifiedLocationWithoutIBX
		return ret
	}

	return o.Location
}

// GetLocationOk returns a tuple with the Location field value
// and a boolean to check if the value has been set.
func (o *VirtualGatewayPostRequest) GetLocationOk() (*SimplifiedLocationWithoutIBX, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Location, true
}

// SetLocation sets field value
func (o *VirtualGatewayPostRequest) SetLocation(v SimplifiedLocationWithoutIBX) {
	o.Location = v
}

// GetPackage returns the Package field value
func (o *VirtualGatewayPostRequest) GetPackage() VirtualGatewayPackageType {
	if o == nil {
		var ret VirtualGatewayPackageType
		return ret
	}

	return o.Package
}

// GetPackageOk returns a tuple with the Package field value
// and a boolean to check if the value has been set.
func (o *VirtualGatewayPostRequest) GetPackageOk() (*VirtualGatewayPackageType, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Package, true
}

// SetPackage sets field value
func (o *VirtualGatewayPostRequest) SetPackage(v VirtualGatewayPackageType) {
	o.Package = v
}

// GetOrder returns the Order field value if set, zero value otherwise.
func (o *VirtualGatewayPostRequest) GetOrder() Order {
	if o == nil || o.Order == nil {
		var ret Order
		return ret
	}
	return *o.Order
}

// GetOrderOk returns a tuple with the Order field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VirtualGatewayPostRequest) GetOrderOk() (*Order, bool) {
	if o == nil || o.Order == nil {
		return nil, false
	}
	return o.Order, true
}

// HasOrder returns a boolean if a field has been set.
func (o *VirtualGatewayPostRequest) HasOrder() bool {
	if o != nil && o.Order != nil {
		return true
	}

	return false
}

// SetOrder gets a reference to the given Order and assigns it to the Order field.
func (o *VirtualGatewayPostRequest) SetOrder(v Order) {
	o.Order = &v
}

// GetProject returns the Project field value if set, zero value otherwise.
func (o *VirtualGatewayPostRequest) GetProject() Project {
	if o == nil || o.Project == nil {
		var ret Project
		return ret
	}
	return *o.Project
}

// GetProjectOk returns a tuple with the Project field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VirtualGatewayPostRequest) GetProjectOk() (*Project, bool) {
	if o == nil || o.Project == nil {
		return nil, false
	}
	return o.Project, true
}

// HasProject returns a boolean if a field has been set.
func (o *VirtualGatewayPostRequest) HasProject() bool {
	if o != nil && o.Project != nil {
		return true
	}

	return false
}

// SetProject gets a reference to the given Project and assigns it to the Project field.
func (o *VirtualGatewayPostRequest) SetProject(v Project) {
	o.Project = &v
}

// GetAccount returns the Account field value
func (o *VirtualGatewayPostRequest) GetAccount() SimplifiedAccount {
	if o == nil {
		var ret SimplifiedAccount
		return ret
	}

	return o.Account
}

// GetAccountOk returns a tuple with the Account field value
// and a boolean to check if the value has been set.
func (o *VirtualGatewayPostRequest) GetAccountOk() (*SimplifiedAccount, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Account, true
}

// SetAccount sets field value
func (o *VirtualGatewayPostRequest) SetAccount(v SimplifiedAccount) {
	o.Account = v
}

// GetNotifications returns the Notifications field value
func (o *VirtualGatewayPostRequest) GetNotifications() []SimplifiedNotification {
	if o == nil {
		var ret []SimplifiedNotification
		return ret
	}

	return o.Notifications
}

// GetNotificationsOk returns a tuple with the Notifications field value
// and a boolean to check if the value has been set.
func (o *VirtualGatewayPostRequest) GetNotificationsOk() ([]SimplifiedNotification, bool) {
	if o == nil {
		return nil, false
	}
	return o.Notifications, true
}

// SetNotifications sets field value
func (o *VirtualGatewayPostRequest) SetNotifications(v []SimplifiedNotification) {
	o.Notifications = v
}

func (o VirtualGatewayPostRequest) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["type"] = o.Type
	}
	if true {
		toSerialize["name"] = o.Name
	}
	if true {
		toSerialize["location"] = o.Location
	}
	if true {
		toSerialize["package"] = o.Package
	}
	if o.Order != nil {
		toSerialize["order"] = o.Order
	}
	if o.Project != nil {
		toSerialize["project"] = o.Project
	}
	if true {
		toSerialize["account"] = o.Account
	}
	if true {
		toSerialize["notifications"] = o.Notifications
	}
	return json.Marshal(toSerialize)
}

type NullableVirtualGatewayPostRequest struct {
	value *VirtualGatewayPostRequest
	isSet bool
}

func (v NullableVirtualGatewayPostRequest) Get() *VirtualGatewayPostRequest {
	return v.value
}

func (v *NullableVirtualGatewayPostRequest) Set(val *VirtualGatewayPostRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableVirtualGatewayPostRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableVirtualGatewayPostRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableVirtualGatewayPostRequest(val *VirtualGatewayPostRequest) *NullableVirtualGatewayPostRequest {
	return &NullableVirtualGatewayPostRequest{value: val, isSet: true}
}

func (v NullableVirtualGatewayPostRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableVirtualGatewayPostRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


