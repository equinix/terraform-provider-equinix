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

// ConnectionBulkPostRequest Create bulk connection post request
type ConnectionBulkPostRequest struct {
	Data []Connection `json:"data,omitempty"`
}

// NewConnectionBulkPostRequest instantiates a new ConnectionBulkPostRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConnectionBulkPostRequest() *ConnectionBulkPostRequest {
	this := ConnectionBulkPostRequest{}
	return &this
}

// NewConnectionBulkPostRequestWithDefaults instantiates a new ConnectionBulkPostRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConnectionBulkPostRequestWithDefaults() *ConnectionBulkPostRequest {
	this := ConnectionBulkPostRequest{}
	return &this
}

// GetData returns the Data field value if set, zero value otherwise.
func (o *ConnectionBulkPostRequest) GetData() []Connection {
	if o == nil || o.Data == nil {
		var ret []Connection
		return ret
	}
	return o.Data
}

// GetDataOk returns a tuple with the Data field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConnectionBulkPostRequest) GetDataOk() ([]Connection, bool) {
	if o == nil || o.Data == nil {
		return nil, false
	}
	return o.Data, true
}

// HasData returns a boolean if a field has been set.
func (o *ConnectionBulkPostRequest) HasData() bool {
	if o != nil && o.Data != nil {
		return true
	}

	return false
}

// SetData gets a reference to the given []Connection and assigns it to the Data field.
func (o *ConnectionBulkPostRequest) SetData(v []Connection) {
	o.Data = v
}

func (o ConnectionBulkPostRequest) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Data != nil {
		toSerialize["data"] = o.Data
	}
	return json.Marshal(toSerialize)
}

type NullableConnectionBulkPostRequest struct {
	value *ConnectionBulkPostRequest
	isSet bool
}

func (v NullableConnectionBulkPostRequest) Get() *ConnectionBulkPostRequest {
	return v.value
}

func (v *NullableConnectionBulkPostRequest) Set(val *ConnectionBulkPostRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableConnectionBulkPostRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableConnectionBulkPostRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConnectionBulkPostRequest(val *ConnectionBulkPostRequest) *NullableConnectionBulkPostRequest {
	return &NullableConnectionBulkPostRequest{value: val, isSet: true}
}

func (v NullableConnectionBulkPostRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConnectionBulkPostRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


