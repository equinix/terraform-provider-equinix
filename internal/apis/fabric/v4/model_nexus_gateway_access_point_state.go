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

// NexusGatewayAccessPointState Access point lifecycle state
type NexusGatewayAccessPointState string

// List of NexusGatewayAccessPointState
const (
	PROVISIONED NexusGatewayAccessPointState = "PROVISIONED"
	PROVISIONING NexusGatewayAccessPointState = "PROVISIONING"
	DEPROVISIONED NexusGatewayAccessPointState = "DEPROVISIONED"
)

// All allowed values of NexusGatewayAccessPointState enum
var AllowedNexusGatewayAccessPointStateEnumValues = []NexusGatewayAccessPointState{
	"PROVISIONED",
	"PROVISIONING",
	"DEPROVISIONED",
}

func (v *NexusGatewayAccessPointState) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := NexusGatewayAccessPointState(value)
	for _, existing := range AllowedNexusGatewayAccessPointStateEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid NexusGatewayAccessPointState", value)
}

// NewNexusGatewayAccessPointStateFromValue returns a pointer to a valid NexusGatewayAccessPointState
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewNexusGatewayAccessPointStateFromValue(v string) (*NexusGatewayAccessPointState, error) {
	ev := NexusGatewayAccessPointState(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for NexusGatewayAccessPointState: valid values are %v", v, AllowedNexusGatewayAccessPointStateEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v NexusGatewayAccessPointState) IsValid() bool {
	for _, existing := range AllowedNexusGatewayAccessPointStateEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to NexusGatewayAccessPointState value
func (v NexusGatewayAccessPointState) Ptr() *NexusGatewayAccessPointState {
	return &v
}

type NullableNexusGatewayAccessPointState struct {
	value *NexusGatewayAccessPointState
	isSet bool
}

func (v NullableNexusGatewayAccessPointState) Get() *NexusGatewayAccessPointState {
	return v.value
}

func (v *NullableNexusGatewayAccessPointState) Set(val *NexusGatewayAccessPointState) {
	v.value = val
	v.isSet = true
}

func (v NullableNexusGatewayAccessPointState) IsSet() bool {
	return v.isSet
}

func (v *NullableNexusGatewayAccessPointState) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableNexusGatewayAccessPointState(val *NexusGatewayAccessPointState) *NullableNexusGatewayAccessPointState {
	return &NullableNexusGatewayAccessPointState{value: val, isSet: true}
}

func (v NullableNexusGatewayAccessPointState) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableNexusGatewayAccessPointState) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

