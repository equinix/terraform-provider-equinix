package api

//ACLTemplate describes Network Edge device ACL template
type ACLTemplate struct {
	UUID              *string                    `json:"uuid,omitempty"`
	Name              *string                    `json:"name,omitempty"`
	Description       *string                    `json:"description,omitempty"`
	MetroCode         *string                    `json:"metroCode,omitempty"`
	VirtualDeviceUUID *string                    `json:"virtualDeviceUuid,omitempty"`
	DeviceACLStatus   *string                    `json:"deviceAclstatus,omitempty"`
	InboundRules      []ACLTemplateInboundRule   `json:"inboundRules,omitempty"`
	DeviceDetails     []ACLTemplateDeviceDetails `json:"virtualDeviceDetails,omitempty"`
}

//ACLTemplateInboundRule describes inbound ACL rule that is part of
//Network Edge device ACL template
type ACLTemplateInboundRule struct {
	SrcType  *string  `json:"srcType,omitempty"`
	Protocol *string  `json:"protocol,omitempty"`
	SrcPort  *string  `json:"srcPort,omitempty"`
	DstPort  *string  `json:"dstPort,omitempty"`
	Subnets  []string `json:"subnets,omitempty"`
	Subnet   *string  `json:"subnet,omitempty"`
	SeqNO    *int     `json:"seqNo,omitempty"`
}

type ACLTemplateDeviceDetails struct {
	UUID      *string `json:"uuid,omitempty"`
	Name      *string `json:"name,omitempty"`
	ACLStatus *string `json:"aclStatus,omitempty"`
}

//ACLTemplatesResponse describes response for a get ACL template collection request
type ACLTemplatesResponse struct {
	Pagination Pagination    `json:"pagination,omitempty"`
	Data       []ACLTemplate `json:"data,omitempty"`
}
