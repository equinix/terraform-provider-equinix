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

// Presence Presence
type Presence string

// List of Presence
const (
	MY_PORTS Presence = "MY_PORTS"
)

// All allowed values of Presence enum
var AllowedPresenceEnumValues = []Presence{
	"MY_PORTS",
}

func (v *Presence) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := Presence(value)
	for _, existing := range AllowedPresenceEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid Presence", value)
}

// NewPresenceFromValue returns a pointer to a valid Presence
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewPresenceFromValue(v string) (*Presence, error) {
	ev := Presence(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for Presence: valid values are %v", v, AllowedPresenceEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v Presence) IsValid() bool {
	for _, existing := range AllowedPresenceEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to Presence value
func (v Presence) Ptr() *Presence {
	return &v
}

type NullablePresence struct {
	value *Presence
	isSet bool
}

func (v NullablePresence) Get() *Presence {
	return v.value
}

func (v *NullablePresence) Set(val *Presence) {
	v.value = val
	v.isSet = true
}

func (v NullablePresence) IsSet() bool {
	return v.isSet
}

func (v *NullablePresence) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePresence(val *Presence) *NullablePresence {
	return &NullablePresence{value: val, isSet: true}
}

func (v NullablePresence) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePresence) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

