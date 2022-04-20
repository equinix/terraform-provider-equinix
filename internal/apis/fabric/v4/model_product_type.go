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
	"fmt"
)

// ProductType Product type
type ProductType string

// List of ProductType
const (
	VIRTUAL_CONNECTION_PRODUCT ProductType = "VIRTUAL_CONNECTION_PRODUCT"
	IP_BLOCK_PRODUCT ProductType = "IP_BLOCK_PRODUCT"
	VIRTUAL_PORT_PRODUCT ProductType = "VIRTUAL_PORT_PRODUCT"
	FABRIC_GATEWAY_PRODUCT ProductType = "FABRIC_GATEWAY_PRODUCT"
)

// All allowed values of ProductType enum
var AllowedProductTypeEnumValues = []ProductType{
	"VIRTUAL_CONNECTION_PRODUCT",
	"IP_BLOCK_PRODUCT",
	"VIRTUAL_PORT_PRODUCT",
	"FABRIC_GATEWAY_PRODUCT",
}

func (v *ProductType) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := ProductType(value)
	for _, existing := range AllowedProductTypeEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid ProductType", value)
}

// NewProductTypeFromValue returns a pointer to a valid ProductType
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewProductTypeFromValue(v string) (*ProductType, error) {
	ev := ProductType(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for ProductType: valid values are %v", v, AllowedProductTypeEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v ProductType) IsValid() bool {
	for _, existing := range AllowedProductTypeEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to ProductType value
func (v ProductType) Ptr() *ProductType {
	return &v
}

type NullableProductType struct {
	value *ProductType
	isSet bool
}

func (v NullableProductType) Get() *ProductType {
	return v.value
}

func (v *NullableProductType) Set(val *ProductType) {
	v.value = val
	v.isSet = true
}

func (v NullableProductType) IsSet() bool {
	return v.isSet
}

func (v *NullableProductType) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableProductType(val *ProductType) *NullableProductType {
	return &NullableProductType{value: val, isSet: true}
}

func (v NullableProductType) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableProductType) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

