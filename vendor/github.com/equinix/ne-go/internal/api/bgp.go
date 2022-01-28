package api

//BGPConfiguration describes Network Edge BGP peering configuration
//Used as a body in requests (create,update) and response (get)
type BGPConfiguration struct {
	UUID               *string `json:"uuid,omitempty"`
	ConnectionUUID     *string `json:"connectionUuid,omitempty"`
	VirtualDeviceUUID  *string `json:"virtualDeviceUuid,omitempty"`
	LocalIPAddress     *string `json:"localIpAddress,omitempty"`
	LocalASN           *int    `json:"localAsn,omitempty"`
	RemoteASN          *int    `json:"remoteAsn,omitempty"`
	RemoteIPAddress    *string `json:"remoteIpAddress,omitempty"`
	AuthenticationKey  *string `json:"authenticationKey,omitempty"`
	State              *string `json:"state,omitempty"`
	ProvisioningStatus *string `json:"provisioningStatus,omitempty"`
}

//BGPConfigurationCreateResponse describes response body for
//BGP Configuration create and update request
type BGPConfigurationCreateResponse struct {
	UUID *string `json:"uuid,omitempty"`
}
