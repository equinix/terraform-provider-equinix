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

// ConnectionIpv4 struct for ConnectionIpv4
type ConnectionIpv4 struct {
	// Customer peering ip
	CustomerPeerIp *string `json:"customerPeerIp,omitempty"`
	// Provider peering ip
	ProviderPeerIp *string `json:"providerPeerIp,omitempty"`
}

// NewConnectionIpv4 instantiates a new ConnectionIpv4 object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConnectionIpv4() *ConnectionIpv4 {
	this := ConnectionIpv4{}
	return &this
}

// NewConnectionIpv4WithDefaults instantiates a new ConnectionIpv4 object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConnectionIpv4WithDefaults() *ConnectionIpv4 {
	this := ConnectionIpv4{}
	return &this
}

// GetCustomerPeerIp returns the CustomerPeerIp field value if set, zero value otherwise.
func (o *ConnectionIpv4) GetCustomerPeerIp() string {
	if o == nil || o.CustomerPeerIp == nil {
		var ret string
		return ret
	}
	return *o.CustomerPeerIp
}

// GetCustomerPeerIpOk returns a tuple with the CustomerPeerIp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConnectionIpv4) GetCustomerPeerIpOk() (*string, bool) {
	if o == nil || o.CustomerPeerIp == nil {
		return nil, false
	}
	return o.CustomerPeerIp, true
}

// HasCustomerPeerIp returns a boolean if a field has been set.
func (o *ConnectionIpv4) HasCustomerPeerIp() bool {
	if o != nil && o.CustomerPeerIp != nil {
		return true
	}

	return false
}

// SetCustomerPeerIp gets a reference to the given string and assigns it to the CustomerPeerIp field.
func (o *ConnectionIpv4) SetCustomerPeerIp(v string) {
	o.CustomerPeerIp = &v
}

// GetProviderPeerIp returns the ProviderPeerIp field value if set, zero value otherwise.
func (o *ConnectionIpv4) GetProviderPeerIp() string {
	if o == nil || o.ProviderPeerIp == nil {
		var ret string
		return ret
	}
	return *o.ProviderPeerIp
}

// GetProviderPeerIpOk returns a tuple with the ProviderPeerIp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConnectionIpv4) GetProviderPeerIpOk() (*string, bool) {
	if o == nil || o.ProviderPeerIp == nil {
		return nil, false
	}
	return o.ProviderPeerIp, true
}

// HasProviderPeerIp returns a boolean if a field has been set.
func (o *ConnectionIpv4) HasProviderPeerIp() bool {
	if o != nil && o.ProviderPeerIp != nil {
		return true
	}

	return false
}

// SetProviderPeerIp gets a reference to the given string and assigns it to the ProviderPeerIp field.
func (o *ConnectionIpv4) SetProviderPeerIp(v string) {
	o.ProviderPeerIp = &v
}

func (o ConnectionIpv4) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.CustomerPeerIp != nil {
		toSerialize["customerPeerIp"] = o.CustomerPeerIp
	}
	if o.ProviderPeerIp != nil {
		toSerialize["providerPeerIp"] = o.ProviderPeerIp
	}
	return json.Marshal(toSerialize)
}

type NullableConnectionIpv4 struct {
	value *ConnectionIpv4
	isSet bool
}

func (v NullableConnectionIpv4) Get() *ConnectionIpv4 {
	return v.value
}

func (v *NullableConnectionIpv4) Set(val *ConnectionIpv4) {
	v.value = val
	v.isSet = true
}

func (v NullableConnectionIpv4) IsSet() bool {
	return v.isSet
}

func (v *NullableConnectionIpv4) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConnectionIpv4(val *ConnectionIpv4) *NullableConnectionIpv4 {
	return &NullableConnectionIpv4{value: val, isSet: true}
}

func (v NullableConnectionIpv4) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConnectionIpv4) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


