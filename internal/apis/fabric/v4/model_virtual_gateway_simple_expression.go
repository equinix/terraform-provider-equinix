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

// VirtualGatewaySimpleExpression struct for VirtualGatewaySimpleExpression
type VirtualGatewaySimpleExpression struct {
	// Possible field names to use on filters:  * `/project/projectId` - project id (mandatory)  * `/name` - Fabric Gateway name  * `/uuid` - Fabric Gateway uuid  * `/state` - Fabric Gateway status  * `/location/metroCode` - Fabric Gateway metro code  * `/location/metroName` - Fabric Gateway metro name  * `/package/code` - Fabric Gateway package  * `/_*` - all-category search 
	Property *string `json:"property,omitempty"`
	// Possible operators to use on filters:  * `=` - equal  * `!=` - not equal  * `>` - greater than  * `>=` - greater than or equal to  * `<` - less than  * `<=` - less than or equal to  * `[NOT] BETWEEN` - (not) between  * `[NOT] LIKE` - (not) like  * `[NOT] IN` - (not) in  * `~*` - case-insensitive like 
	Operator *string `json:"operator,omitempty"`
	Values []string `json:"values,omitempty"`
}

// NewVirtualGatewaySimpleExpression instantiates a new VirtualGatewaySimpleExpression object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewVirtualGatewaySimpleExpression() *VirtualGatewaySimpleExpression {
	this := VirtualGatewaySimpleExpression{}
	return &this
}

// NewVirtualGatewaySimpleExpressionWithDefaults instantiates a new VirtualGatewaySimpleExpression object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewVirtualGatewaySimpleExpressionWithDefaults() *VirtualGatewaySimpleExpression {
	this := VirtualGatewaySimpleExpression{}
	return &this
}

// GetProperty returns the Property field value if set, zero value otherwise.
func (o *VirtualGatewaySimpleExpression) GetProperty() string {
	if o == nil || o.Property == nil {
		var ret string
		return ret
	}
	return *o.Property
}

// GetPropertyOk returns a tuple with the Property field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VirtualGatewaySimpleExpression) GetPropertyOk() (*string, bool) {
	if o == nil || o.Property == nil {
		return nil, false
	}
	return o.Property, true
}

// HasProperty returns a boolean if a field has been set.
func (o *VirtualGatewaySimpleExpression) HasProperty() bool {
	if o != nil && o.Property != nil {
		return true
	}

	return false
}

// SetProperty gets a reference to the given string and assigns it to the Property field.
func (o *VirtualGatewaySimpleExpression) SetProperty(v string) {
	o.Property = &v
}

// GetOperator returns the Operator field value if set, zero value otherwise.
func (o *VirtualGatewaySimpleExpression) GetOperator() string {
	if o == nil || o.Operator == nil {
		var ret string
		return ret
	}
	return *o.Operator
}

// GetOperatorOk returns a tuple with the Operator field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VirtualGatewaySimpleExpression) GetOperatorOk() (*string, bool) {
	if o == nil || o.Operator == nil {
		return nil, false
	}
	return o.Operator, true
}

// HasOperator returns a boolean if a field has been set.
func (o *VirtualGatewaySimpleExpression) HasOperator() bool {
	if o != nil && o.Operator != nil {
		return true
	}

	return false
}

// SetOperator gets a reference to the given string and assigns it to the Operator field.
func (o *VirtualGatewaySimpleExpression) SetOperator(v string) {
	o.Operator = &v
}

// GetValues returns the Values field value if set, zero value otherwise.
func (o *VirtualGatewaySimpleExpression) GetValues() []string {
	if o == nil || o.Values == nil {
		var ret []string
		return ret
	}
	return o.Values
}

// GetValuesOk returns a tuple with the Values field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VirtualGatewaySimpleExpression) GetValuesOk() ([]string, bool) {
	if o == nil || o.Values == nil {
		return nil, false
	}
	return o.Values, true
}

// HasValues returns a boolean if a field has been set.
func (o *VirtualGatewaySimpleExpression) HasValues() bool {
	if o != nil && o.Values != nil {
		return true
	}

	return false
}

// SetValues gets a reference to the given []string and assigns it to the Values field.
func (o *VirtualGatewaySimpleExpression) SetValues(v []string) {
	o.Values = v
}

func (o VirtualGatewaySimpleExpression) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Property != nil {
		toSerialize["property"] = o.Property
	}
	if o.Operator != nil {
		toSerialize["operator"] = o.Operator
	}
	if o.Values != nil {
		toSerialize["values"] = o.Values
	}
	return json.Marshal(toSerialize)
}

type NullableVirtualGatewaySimpleExpression struct {
	value *VirtualGatewaySimpleExpression
	isSet bool
}

func (v NullableVirtualGatewaySimpleExpression) Get() *VirtualGatewaySimpleExpression {
	return v.value
}

func (v *NullableVirtualGatewaySimpleExpression) Set(val *VirtualGatewaySimpleExpression) {
	v.value = val
	v.isSet = true
}

func (v NullableVirtualGatewaySimpleExpression) IsSet() bool {
	return v.isSet
}

func (v *NullableVirtualGatewaySimpleExpression) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableVirtualGatewaySimpleExpression(val *VirtualGatewaySimpleExpression) *NullableVirtualGatewaySimpleExpression {
	return &NullableVirtualGatewaySimpleExpression{value: val, isSet: true}
}

func (v NullableVirtualGatewaySimpleExpression) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableVirtualGatewaySimpleExpression) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


