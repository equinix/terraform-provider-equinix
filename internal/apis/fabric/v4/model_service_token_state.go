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

// ServiceTokenState Service token state
type ServiceTokenState string

// List of ServiceTokenState
const (
	ACTIVE ServiceTokenState = "ACTIVE"
	INACTIVE ServiceTokenState = "INACTIVE"
	EXPIRED ServiceTokenState = "EXPIRED"
	DELETED ServiceTokenState = "DELETED"
)

// All allowed values of ServiceTokenState enum
var AllowedServiceTokenStateEnumValues = []ServiceTokenState{
	"ACTIVE",
	"INACTIVE",
	"EXPIRED",
	"DELETED",
}

func (v *ServiceTokenState) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := ServiceTokenState(value)
	for _, existing := range AllowedServiceTokenStateEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid ServiceTokenState", value)
}

// NewServiceTokenStateFromValue returns a pointer to a valid ServiceTokenState
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewServiceTokenStateFromValue(v string) (*ServiceTokenState, error) {
	ev := ServiceTokenState(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for ServiceTokenState: valid values are %v", v, AllowedServiceTokenStateEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v ServiceTokenState) IsValid() bool {
	for _, existing := range AllowedServiceTokenStateEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to ServiceTokenState value
func (v ServiceTokenState) Ptr() *ServiceTokenState {
	return &v
}

type NullableServiceTokenState struct {
	value *ServiceTokenState
	isSet bool
}

func (v NullableServiceTokenState) Get() *ServiceTokenState {
	return v.value
}

func (v *NullableServiceTokenState) Set(val *ServiceTokenState) {
	v.value = val
	v.isSet = true
}

func (v NullableServiceTokenState) IsSet() bool {
	return v.isSet
}

func (v *NullableServiceTokenState) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableServiceTokenState(val *ServiceTokenState) *NullableServiceTokenState {
	return &NullableServiceTokenState{value: val, isSet: true}
}

func (v NullableServiceTokenState) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableServiceTokenState) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

