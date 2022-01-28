package ne

import (
	"net/http"
	"net/url"

	"github.com/equinix/ne-go/internal/api"
	"github.com/equinix/rest-go"
)

type restDeviceLinkUpdateRequest struct {
	uuid      string
	groupName *string
	subnet    *string
	devices   []DeviceLinkGroupDevice
	links     []DeviceLinkGroupLink
	c         RestClient
}

//GetDeviceLinkGroups retrieves list of existing device link groups
//(along with their details)
func (c RestClient) GetDeviceLinkGroups() ([]DeviceLinkGroup, error) {
	path := "/ne/v1/links"
	content, err := c.GetOffsetPaginated(path, &api.DeviceLinkGroupsGetResponse{},
		rest.DefaultOffsetPagingConfig())
	if err != nil {
		return nil, err
	}
	transformed := make([]DeviceLinkGroup, len(content))
	for i := range content {
		transformed[i] = *mapDeviceLinkGroupAPIToDomain(content[i].(api.DeviceLinkGroup))

	}
	return transformed, nil
}

//GetDeviceLinkGroups retrieves details of a device link group
//with a given identifier
func (c RestClient) GetDeviceLinkGroup(uuid string) (*DeviceLinkGroup, error) {
	path := "/ne/v1/links/" + url.PathEscape(uuid)
	result := api.DeviceLinkGroup{}
	request := c.R().SetResult(&result)
	if err := c.Execute(request, http.MethodGet, path); err != nil {
		return nil, err
	}
	return mapDeviceLinkGroupAPIToDomain(result), nil
}

//CreateDeviceLinkGroup creates given device link group and returns
//its identifier upon successful creation
func (c RestClient) CreateDeviceLinkGroup(linkGroup DeviceLinkGroup) (*string, error) {
	path := "/ne/v1/links"
	reqBody := mapDeviceLinkGroupDomainToAPI(linkGroup)
	respBody := api.DeviceLinkGroupCreateResponse{}
	req := c.R().SetBody(&reqBody).SetResult(&respBody)
	if err := c.Execute(req, http.MethodPost, path); err != nil {
		return nil, err
	}
	return respBody.UUID, nil
}

//NewDeviceLinkGroupUpdateRequest creates new update request for a device link
//group with a given identifier
func (c RestClient) NewDeviceLinkGroupUpdateRequest(uuid string) DeviceLinkUpdateRequest {
	return &restDeviceLinkUpdateRequest{uuid: uuid, c: c}
}

//DeleteDeviceLinkGroup removes device link group with a given identifier
func (c RestClient) DeleteDeviceLinkGroup(uuid string) error {
	path := "/ne/v1/links/" + url.PathEscape(uuid)
	if err := c.Execute(c.R(), http.MethodDelete, path); err != nil {
		return err
	}
	return nil
}

func (req *restDeviceLinkUpdateRequest) WithGroupName(name string) DeviceLinkUpdateRequest {
	req.groupName = &name
	return req
}

func (req *restDeviceLinkUpdateRequest) WithSubnet(subnet string) DeviceLinkUpdateRequest {
	req.subnet = &subnet
	return req
}

func (req *restDeviceLinkUpdateRequest) WithDevices(devices []DeviceLinkGroupDevice) DeviceLinkUpdateRequest {
	req.devices = devices
	return req
}

func (req *restDeviceLinkUpdateRequest) WithLinks(links []DeviceLinkGroupLink) DeviceLinkUpdateRequest {
	req.links = links
	return req
}

func (req *restDeviceLinkUpdateRequest) Execute() error {
	reqBody := api.DeviceLinkGroupUpdateRequest{}
	if StringValue(req.groupName) != "" {
		reqBody.GroupName = req.groupName
	}
	if StringValue(req.subnet) != "" {
		reqBody.Subnet = req.subnet
	}
	reqBody.Links = make([]api.DeviceLinkGroupLink, len(req.links))
	for i := range req.links {
		reqBody.Links[i] = mapDeviceLinkGroupLinkDomainToAPI(req.links[i])
	}
	reqBody.Devices = make([]api.DeviceLinkGroupDevice, len(req.devices))
	for i := range req.devices {
		reqBody.Devices[i] = mapDeviceLinkGroupDeviceDomainToAPI(req.devices[i])
	}
	path := "/ne/v1/links/" + url.PathEscape(req.uuid)
	httpReq := req.c.R().SetBody(&reqBody)
	if err := req.c.Execute(httpReq, http.MethodPatch, path); err != nil {
		return err
	}
	return nil
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Unexported package methods
//_______________________________________________________________________

func mapDeviceLinkGroupAPIToDomain(apiLinkGroup api.DeviceLinkGroup) *DeviceLinkGroup {
	linkGroup := DeviceLinkGroup{}
	linkGroup.UUID = apiLinkGroup.UUID
	linkGroup.Name = apiLinkGroup.GroupName
	linkGroup.Subnet = apiLinkGroup.Subnet
	linkGroup.Status = apiLinkGroup.Status
	linkGroup.Devices = make([]DeviceLinkGroupDevice, len(apiLinkGroup.Devices))
	for i := range apiLinkGroup.Devices {
		linkGroup.Devices[i] = mapDeviceLinkGroupDeviceAPIToDomain(apiLinkGroup.Devices[i])
	}
	linkGroup.Links = make([]DeviceLinkGroupLink, len(apiLinkGroup.Links))
	for i := range apiLinkGroup.Links {
		linkGroup.Links[i] = mapDeviceLinkGroupLinkAPIToDomain(apiLinkGroup.Links[i])
	}
	return &linkGroup
}

func mapDeviceLinkGroupDeviceAPIToDomain(apiLinkGroupDevice api.DeviceLinkGroupDevice) DeviceLinkGroupDevice {
	return DeviceLinkGroupDevice{
		DeviceID:    apiLinkGroupDevice.DeviceUUID,
		ASN:         apiLinkGroupDevice.ASN,
		InterfaceID: apiLinkGroupDevice.InterfaceID,
		Status:      apiLinkGroupDevice.Status,
		IPAddress:   apiLinkGroupDevice.IPAddress,
	}
}

func mapDeviceLinkGroupLinkAPIToDomain(apiLinkGroupLink api.DeviceLinkGroupLink) DeviceLinkGroupLink {
	return DeviceLinkGroupLink{
		AccountNumber:        apiLinkGroupLink.AccountNumber,
		Throughput:           apiLinkGroupLink.Throughput,
		ThroughputUnit:       apiLinkGroupLink.ThroughputUnit,
		SourceMetroCode:      apiLinkGroupLink.SourceMetroCode,
		DestinationMetroCode: apiLinkGroupLink.DestinationMetroCode,
		SourceZoneCode:       apiLinkGroupLink.SourceZoneCode,
		DestinationZoneCode:  apiLinkGroupLink.DestinationZoneCode,
	}
}

func mapDeviceLinkGroupDomainToAPI(linkGroup DeviceLinkGroup) api.DeviceLinkGroup {
	apiLinkGroup := api.DeviceLinkGroup{}
	apiLinkGroup.GroupName = linkGroup.Name
	apiLinkGroup.Subnet = linkGroup.Subnet
	apiLinkGroup.Devices = make([]api.DeviceLinkGroupDevice, len(linkGroup.Devices))
	for i := range linkGroup.Devices {
		apiLinkGroup.Devices[i] = mapDeviceLinkGroupDeviceDomainToAPI(linkGroup.Devices[i])
	}
	apiLinkGroup.Links = make([]api.DeviceLinkGroupLink, len(linkGroup.Links))
	for i := range linkGroup.Links {
		apiLinkGroup.Links[i] = mapDeviceLinkGroupLinkDomainToAPI(linkGroup.Links[i])
	}
	return apiLinkGroup
}

func mapDeviceLinkGroupDeviceDomainToAPI(linkGroupDevice DeviceLinkGroupDevice) api.DeviceLinkGroupDevice {
	return api.DeviceLinkGroupDevice{
		DeviceUUID:  linkGroupDevice.DeviceID,
		ASN:         linkGroupDevice.ASN,
		InterfaceID: linkGroupDevice.InterfaceID,
		Status:      linkGroupDevice.Status,
		IPAddress:   linkGroupDevice.IPAddress,
	}
}

func mapDeviceLinkGroupLinkDomainToAPI(linkGroupLink DeviceLinkGroupLink) api.DeviceLinkGroupLink {
	return api.DeviceLinkGroupLink{
		AccountNumber:        linkGroupLink.AccountNumber,
		Throughput:           linkGroupLink.Throughput,
		ThroughputUnit:       linkGroupLink.ThroughputUnit,
		SourceMetroCode:      linkGroupLink.SourceMetroCode,
		DestinationMetroCode: linkGroupLink.DestinationMetroCode,
		SourceZoneCode:       linkGroupLink.SourceZoneCode,
		DestinationZoneCode:  linkGroupLink.DestinationZoneCode,
	}
}
