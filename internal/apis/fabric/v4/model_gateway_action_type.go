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

// GatewayActionType Gateway action type
type GatewayActionType string

// List of GatewayActionType
const (
	BGP_SESSION_STATUS_UPDATE GatewayActionType = "BGP_SESSION_STATUS_UPDATE"
	ROUTE_TABLE_ENTRY_UPDATE GatewayActionType = "ROUTE_TABLE_ENTRY_UPDATE"
)

// All allowed values of GatewayActionType enum
var AllowedGatewayActionTypeEnumValues = []GatewayActionType{
	"BGP_SESSION_STATUS_UPDATE",
	"ROUTE_TABLE_ENTRY_UPDATE",
}

func (v *GatewayActionType) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := GatewayActionType(value)
	for _, existing := range AllowedGatewayActionTypeEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid GatewayActionType", value)
}

// NewGatewayActionTypeFromValue returns a pointer to a valid GatewayActionType
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewGatewayActionTypeFromValue(v string) (*GatewayActionType, error) {
	ev := GatewayActionType(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for GatewayActionType: valid values are %v", v, AllowedGatewayActionTypeEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v GatewayActionType) IsValid() bool {
	for _, existing := range AllowedGatewayActionTypeEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to GatewayActionType value
func (v GatewayActionType) Ptr() *GatewayActionType {
	return &v
}

type NullableGatewayActionType struct {
	value *GatewayActionType
	isSet bool
}

func (v NullableGatewayActionType) Get() *GatewayActionType {
	return v.value
}

func (v *NullableGatewayActionType) Set(val *GatewayActionType) {
	v.value = val
	v.isSet = true
}

func (v NullableGatewayActionType) IsSet() bool {
	return v.isSet
}

func (v *NullableGatewayActionType) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGatewayActionType(val *GatewayActionType) *NullableGatewayActionType {
	return &NullableGatewayActionType{value: val, isSet: true}
}

func (v NullableGatewayActionType) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGatewayActionType) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

