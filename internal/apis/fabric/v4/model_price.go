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

// Price struct for Price
type Price struct {
	// An absolute URL that returns specified pricing data
	Href *string `json:"href,omitempty"`
	Type *ProductType `json:"type,omitempty"`
	// Equinix-assigned product code
	Code *string `json:"code,omitempty"`
	// Full product name
	Name *string `json:"name,omitempty"`
	// Product description
	Description *string `json:"description,omitempty"`
	Account *SimplifiedAccount `json:"account,omitempty"`
	Charges []PriceCharge `json:"charges,omitempty"`
	// Product offering price currency
	Currency *string `json:"currency,omitempty"`
	// In months. No value means unlimited
	TermLength *int32 `json:"termLength,omitempty"`
	Catgory *PriceCategory `json:"catgory,omitempty"`
	Connection *VirtualConnectionPrice `json:"connection,omitempty"`
	IpBlock *IpBlockPrice `json:"ipBlock,omitempty"`
	Gateway *FabricGatewayPrice `json:"gateway,omitempty"`
	Port *VirtualPortPrice `json:"port,omitempty"`
}

// NewPrice instantiates a new Price object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPrice() *Price {
	this := Price{}
	return &this
}

// NewPriceWithDefaults instantiates a new Price object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPriceWithDefaults() *Price {
	this := Price{}
	return &this
}

// GetHref returns the Href field value if set, zero value otherwise.
func (o *Price) GetHref() string {
	if o == nil || o.Href == nil {
		var ret string
		return ret
	}
	return *o.Href
}

// GetHrefOk returns a tuple with the Href field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Price) GetHrefOk() (*string, bool) {
	if o == nil || o.Href == nil {
		return nil, false
	}
	return o.Href, true
}

// HasHref returns a boolean if a field has been set.
func (o *Price) HasHref() bool {
	if o != nil && o.Href != nil {
		return true
	}

	return false
}

// SetHref gets a reference to the given string and assigns it to the Href field.
func (o *Price) SetHref(v string) {
	o.Href = &v
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *Price) GetType() ProductType {
	if o == nil || o.Type == nil {
		var ret ProductType
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Price) GetTypeOk() (*ProductType, bool) {
	if o == nil || o.Type == nil {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *Price) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}

// SetType gets a reference to the given ProductType and assigns it to the Type field.
func (o *Price) SetType(v ProductType) {
	o.Type = &v
}

// GetCode returns the Code field value if set, zero value otherwise.
func (o *Price) GetCode() string {
	if o == nil || o.Code == nil {
		var ret string
		return ret
	}
	return *o.Code
}

// GetCodeOk returns a tuple with the Code field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Price) GetCodeOk() (*string, bool) {
	if o == nil || o.Code == nil {
		return nil, false
	}
	return o.Code, true
}

// HasCode returns a boolean if a field has been set.
func (o *Price) HasCode() bool {
	if o != nil && o.Code != nil {
		return true
	}

	return false
}

// SetCode gets a reference to the given string and assigns it to the Code field.
func (o *Price) SetCode(v string) {
	o.Code = &v
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *Price) GetName() string {
	if o == nil || o.Name == nil {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Price) GetNameOk() (*string, bool) {
	if o == nil || o.Name == nil {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *Price) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *Price) SetName(v string) {
	o.Name = &v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *Price) GetDescription() string {
	if o == nil || o.Description == nil {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Price) GetDescriptionOk() (*string, bool) {
	if o == nil || o.Description == nil {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *Price) HasDescription() bool {
	if o != nil && o.Description != nil {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *Price) SetDescription(v string) {
	o.Description = &v
}

// GetAccount returns the Account field value if set, zero value otherwise.
func (o *Price) GetAccount() SimplifiedAccount {
	if o == nil || o.Account == nil {
		var ret SimplifiedAccount
		return ret
	}
	return *o.Account
}

// GetAccountOk returns a tuple with the Account field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Price) GetAccountOk() (*SimplifiedAccount, bool) {
	if o == nil || o.Account == nil {
		return nil, false
	}
	return o.Account, true
}

// HasAccount returns a boolean if a field has been set.
func (o *Price) HasAccount() bool {
	if o != nil && o.Account != nil {
		return true
	}

	return false
}

// SetAccount gets a reference to the given SimplifiedAccount and assigns it to the Account field.
func (o *Price) SetAccount(v SimplifiedAccount) {
	o.Account = &v
}

// GetCharges returns the Charges field value if set, zero value otherwise.
func (o *Price) GetCharges() []PriceCharge {
	if o == nil || o.Charges == nil {
		var ret []PriceCharge
		return ret
	}
	return o.Charges
}

// GetChargesOk returns a tuple with the Charges field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Price) GetChargesOk() ([]PriceCharge, bool) {
	if o == nil || o.Charges == nil {
		return nil, false
	}
	return o.Charges, true
}

// HasCharges returns a boolean if a field has been set.
func (o *Price) HasCharges() bool {
	if o != nil && o.Charges != nil {
		return true
	}

	return false
}

// SetCharges gets a reference to the given []PriceCharge and assigns it to the Charges field.
func (o *Price) SetCharges(v []PriceCharge) {
	o.Charges = v
}

// GetCurrency returns the Currency field value if set, zero value otherwise.
func (o *Price) GetCurrency() string {
	if o == nil || o.Currency == nil {
		var ret string
		return ret
	}
	return *o.Currency
}

// GetCurrencyOk returns a tuple with the Currency field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Price) GetCurrencyOk() (*string, bool) {
	if o == nil || o.Currency == nil {
		return nil, false
	}
	return o.Currency, true
}

// HasCurrency returns a boolean if a field has been set.
func (o *Price) HasCurrency() bool {
	if o != nil && o.Currency != nil {
		return true
	}

	return false
}

// SetCurrency gets a reference to the given string and assigns it to the Currency field.
func (o *Price) SetCurrency(v string) {
	o.Currency = &v
}

// GetTermLength returns the TermLength field value if set, zero value otherwise.
func (o *Price) GetTermLength() int32 {
	if o == nil || o.TermLength == nil {
		var ret int32
		return ret
	}
	return *o.TermLength
}

// GetTermLengthOk returns a tuple with the TermLength field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Price) GetTermLengthOk() (*int32, bool) {
	if o == nil || o.TermLength == nil {
		return nil, false
	}
	return o.TermLength, true
}

// HasTermLength returns a boolean if a field has been set.
func (o *Price) HasTermLength() bool {
	if o != nil && o.TermLength != nil {
		return true
	}

	return false
}

// SetTermLength gets a reference to the given int32 and assigns it to the TermLength field.
func (o *Price) SetTermLength(v int32) {
	o.TermLength = &v
}

// GetCatgory returns the Catgory field value if set, zero value otherwise.
func (o *Price) GetCatgory() PriceCategory {
	if o == nil || o.Catgory == nil {
		var ret PriceCategory
		return ret
	}
	return *o.Catgory
}

// GetCatgoryOk returns a tuple with the Catgory field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Price) GetCatgoryOk() (*PriceCategory, bool) {
	if o == nil || o.Catgory == nil {
		return nil, false
	}
	return o.Catgory, true
}

// HasCatgory returns a boolean if a field has been set.
func (o *Price) HasCatgory() bool {
	if o != nil && o.Catgory != nil {
		return true
	}

	return false
}

// SetCatgory gets a reference to the given PriceCategory and assigns it to the Catgory field.
func (o *Price) SetCatgory(v PriceCategory) {
	o.Catgory = &v
}

// GetConnection returns the Connection field value if set, zero value otherwise.
func (o *Price) GetConnection() VirtualConnectionPrice {
	if o == nil || o.Connection == nil {
		var ret VirtualConnectionPrice
		return ret
	}
	return *o.Connection
}

// GetConnectionOk returns a tuple with the Connection field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Price) GetConnectionOk() (*VirtualConnectionPrice, bool) {
	if o == nil || o.Connection == nil {
		return nil, false
	}
	return o.Connection, true
}

// HasConnection returns a boolean if a field has been set.
func (o *Price) HasConnection() bool {
	if o != nil && o.Connection != nil {
		return true
	}

	return false
}

// SetConnection gets a reference to the given VirtualConnectionPrice and assigns it to the Connection field.
func (o *Price) SetConnection(v VirtualConnectionPrice) {
	o.Connection = &v
}

// GetIpBlock returns the IpBlock field value if set, zero value otherwise.
func (o *Price) GetIpBlock() IpBlockPrice {
	if o == nil || o.IpBlock == nil {
		var ret IpBlockPrice
		return ret
	}
	return *o.IpBlock
}

// GetIpBlockOk returns a tuple with the IpBlock field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Price) GetIpBlockOk() (*IpBlockPrice, bool) {
	if o == nil || o.IpBlock == nil {
		return nil, false
	}
	return o.IpBlock, true
}

// HasIpBlock returns a boolean if a field has been set.
func (o *Price) HasIpBlock() bool {
	if o != nil && o.IpBlock != nil {
		return true
	}

	return false
}

// SetIpBlock gets a reference to the given IpBlockPrice and assigns it to the IpBlock field.
func (o *Price) SetIpBlock(v IpBlockPrice) {
	o.IpBlock = &v
}

// GetGateway returns the Gateway field value if set, zero value otherwise.
func (o *Price) GetGateway() FabricGatewayPrice {
	if o == nil || o.Gateway == nil {
		var ret FabricGatewayPrice
		return ret
	}
	return *o.Gateway
}

// GetGatewayOk returns a tuple with the Gateway field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Price) GetGatewayOk() (*FabricGatewayPrice, bool) {
	if o == nil || o.Gateway == nil {
		return nil, false
	}
	return o.Gateway, true
}

// HasGateway returns a boolean if a field has been set.
func (o *Price) HasGateway() bool {
	if o != nil && o.Gateway != nil {
		return true
	}

	return false
}

// SetGateway gets a reference to the given FabricGatewayPrice and assigns it to the Gateway field.
func (o *Price) SetGateway(v FabricGatewayPrice) {
	o.Gateway = &v
}

// GetPort returns the Port field value if set, zero value otherwise.
func (o *Price) GetPort() VirtualPortPrice {
	if o == nil || o.Port == nil {
		var ret VirtualPortPrice
		return ret
	}
	return *o.Port
}

// GetPortOk returns a tuple with the Port field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Price) GetPortOk() (*VirtualPortPrice, bool) {
	if o == nil || o.Port == nil {
		return nil, false
	}
	return o.Port, true
}

// HasPort returns a boolean if a field has been set.
func (o *Price) HasPort() bool {
	if o != nil && o.Port != nil {
		return true
	}

	return false
}

// SetPort gets a reference to the given VirtualPortPrice and assigns it to the Port field.
func (o *Price) SetPort(v VirtualPortPrice) {
	o.Port = &v
}

func (o Price) MarshalJSON() ([]byte, error) {
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
	if o.Name != nil {
		toSerialize["name"] = o.Name
	}
	if o.Description != nil {
		toSerialize["description"] = o.Description
	}
	if o.Account != nil {
		toSerialize["account"] = o.Account
	}
	if o.Charges != nil {
		toSerialize["charges"] = o.Charges
	}
	if o.Currency != nil {
		toSerialize["currency"] = o.Currency
	}
	if o.TermLength != nil {
		toSerialize["termLength"] = o.TermLength
	}
	if o.Catgory != nil {
		toSerialize["catgory"] = o.Catgory
	}
	if o.Connection != nil {
		toSerialize["connection"] = o.Connection
	}
	if o.IpBlock != nil {
		toSerialize["ipBlock"] = o.IpBlock
	}
	if o.Gateway != nil {
		toSerialize["gateway"] = o.Gateway
	}
	if o.Port != nil {
		toSerialize["port"] = o.Port
	}
	return json.Marshal(toSerialize)
}

type NullablePrice struct {
	value *Price
	isSet bool
}

func (v NullablePrice) Get() *Price {
	return v.value
}

func (v *NullablePrice) Set(val *Price) {
	v.value = val
	v.isSet = true
}

func (v NullablePrice) IsSet() bool {
	return v.isSet
}

func (v *NullablePrice) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePrice(val *Price) *NullablePrice {
	return &NullablePrice{value: val, isSet: true}
}

func (v NullablePrice) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePrice) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


