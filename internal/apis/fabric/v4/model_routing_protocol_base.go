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

// RoutingProtocolBase struct for RoutingProtocolBase
type RoutingProtocolBase struct {
	RoutingProtocolBGPType *RoutingProtocolBGPType
	RoutingProtocolDirectType *RoutingProtocolDirectType
}

// Unmarshal JSON data into any of the pointers in the struct
func (dst *RoutingProtocolBase) UnmarshalJSON(data []byte) error {
	var err error
	// try to unmarshal JSON data into RoutingProtocolBGPType
	err = json.Unmarshal(data, &dst.RoutingProtocolBGPType);
	if err == nil {
		jsonRoutingProtocolBGPType, _ := json.Marshal(dst.RoutingProtocolBGPType)
		if string(jsonRoutingProtocolBGPType) == "{}" { // empty struct
			dst.RoutingProtocolBGPType = nil
		} else {
			return nil // data stored in dst.RoutingProtocolBGPType, return on the first match
		}
	} else {
		dst.RoutingProtocolBGPType = nil
	}

	// try to unmarshal JSON data into RoutingProtocolDirectType
	err = json.Unmarshal(data, &dst.RoutingProtocolDirectType);
	if err == nil {
		jsonRoutingProtocolDirectType, _ := json.Marshal(dst.RoutingProtocolDirectType)
		if string(jsonRoutingProtocolDirectType) == "{}" { // empty struct
			dst.RoutingProtocolDirectType = nil
		} else {
			return nil // data stored in dst.RoutingProtocolDirectType, return on the first match
		}
	} else {
		dst.RoutingProtocolDirectType = nil
	}

	return fmt.Errorf("Data failed to match schemas in anyOf(RoutingProtocolBase)")
}

// Marshal data from the first non-nil pointers in the struct to JSON
func (src *RoutingProtocolBase) MarshalJSON() ([]byte, error) {
	if src.RoutingProtocolBGPType != nil {
		return json.Marshal(&src.RoutingProtocolBGPType)
	}

	if src.RoutingProtocolDirectType != nil {
		return json.Marshal(&src.RoutingProtocolDirectType)
	}

	return nil, nil // no data in anyOf schemas
}

type NullableRoutingProtocolBase struct {
	value *RoutingProtocolBase
	isSet bool
}

func (v NullableRoutingProtocolBase) Get() *RoutingProtocolBase {
	return v.value
}

func (v *NullableRoutingProtocolBase) Set(val *RoutingProtocolBase) {
	v.value = val
	v.isSet = true
}

func (v NullableRoutingProtocolBase) IsSet() bool {
	return v.isSet
}

func (v *NullableRoutingProtocolBase) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRoutingProtocolBase(val *RoutingProtocolBase) *NullableRoutingProtocolBase {
	return &NullableRoutingProtocolBase{value: val, isSet: true}
}

func (v NullableRoutingProtocolBase) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRoutingProtocolBase) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


