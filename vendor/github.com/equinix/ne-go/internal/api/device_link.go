package api

type DeviceLinkGroup struct {
	UUID      *string                 `json:"uuid,omitempty"`
	GroupName *string                 `json:"groupName,omitempty"`
	Subnet    *string                 `json:"subnet,omitempty"`
	Status    *string                 `json:"status,omitempty"`
	Devices   []DeviceLinkGroupDevice `json:"linkDevices,omitempty"`
	Links     []DeviceLinkGroupLink   `json:"links"`
}

type DeviceLinkGroupUpdateRequest struct {
	GroupName *string                 `json:"groupName,omitempty"`
	Subnet    *string                 `json:"subnet,omitempty"`
	Devices   []DeviceLinkGroupDevice `json:"linkDevices,omitempty"`
	Links     []DeviceLinkGroupLink   `json:"links,omitempty"`
}

type DeviceLinkGroupDevice struct {
	DeviceUUID  *string `json:"deviceUuid,omitempty"`
	ASN         *int    `json:"asn,omitempty"`
	InterfaceID *int    `json:"interfaceId,omitempty"`
	Status      *string `json:"status,omitempty"`
	IPAddress   *string `json:"ipAssigned,omitempty"`
}

type DeviceLinkGroupLink struct {
	AccountNumber        *string `json:"accountNumber,omitempty"`
	Throughput           *string `json:"throughput,omitempty"`
	ThroughputUnit       *string `json:"throughputUnit,omitempty"`
	SourceMetroCode      *string `json:"sourceMetroCode,omitempty"`
	DestinationMetroCode *string `json:"destinationMetroCode,omitempty"`
	SourceZoneCode       *string `json:"sourceZoneCode,omitempty"`
	DestinationZoneCode  *string `json:"destinationZoneCode,omitempty"`
}

type DeviceLinkGroupCreateResponse struct {
	UUID *string `json:"uuid,omitempty"`
}

type DeviceLinkGroupsGetResponse struct {
	Pagination Pagination        `json:"pagination,omitempty"`
	Data       []DeviceLinkGroup `json:"data,omitempty"`
}
