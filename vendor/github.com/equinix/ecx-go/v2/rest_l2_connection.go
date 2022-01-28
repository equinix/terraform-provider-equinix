package ecx

import (
	"net/http"
	"net/url"

	"github.com/equinix/ecx-go/v2/internal/api"
	"github.com/equinix/rest-go"
)

type restL2ConnectionUpdateRequest struct {
	uuid      string
	name      *string
	speed     *int
	speedUnit *string
	c         RestClient
}

//GetL2OutgoingConnections retrieves list of all originating (a-side) layer 2 connections
//for a customer account associated with authenticated application
func (c RestClient) GetL2OutgoingConnections(statuses []string) ([]L2Connection, error) {
	path := "/ecx/v3/l2/buyer/connections"
	pagingConfig := rest.DefaultPagingConfig().
		SetSizeParamName("pageSize").
		SetPageParamName("pageNumber").
		SetFirstPageNumber(0)
	if len(statuses) > 0 {
		pagingConfig.SetAdditionalParams(map[string]string{"status": buildQueryParamValueString(statuses)})
	}
	content, err := c.GetPaginated(path, &api.L2BuyerConnectionsResponse{}, pagingConfig)
	if err != nil {
		return nil, err
	}
	transformed := make([]L2Connection, len(content))
	for i := range content {
		transformed[i] = *mapGETToL2Connection(content[i].(api.L2ConnectionResponse))
	}
	return transformed, nil
}

//GetL2Connection operation retrieves layer 2 connection with a given UUID
func (c RestClient) GetL2Connection(uuid string) (*L2Connection, error) {
	path := "/ecx/v3/l2/connections/" + url.PathEscape(uuid)
	respBody := api.L2ConnectionResponse{}
	req := c.R().SetResult(&respBody)
	if err := c.Execute(req, http.MethodGet, path); err != nil {
		return nil, err
	}
	return mapGETToL2Connection(respBody), nil
}

//CreateL2Connection operation creates non-redundant layer 2 connection with a given connection structure.
//Upon successful creation, connection structure, enriched with assigned UUID, will be returned
func (c RestClient) CreateL2Connection(l2connection L2Connection) (*string, error) {
	path := "/ecx/v3/l2/connections"
	if StringValue(l2connection.DeviceUUID) != "" {
		path = "/ne/v1/l2/connections"
	}
	reqBody := createL2ConnectionRequest(l2connection)
	respBody := api.CreateL2ConnectionResponse{}
	req := c.R().SetBody(&reqBody).SetResult(&respBody)
	if err := c.Execute(req, http.MethodPost, path); err != nil {
		return nil, err
	}
	return respBody.PrimaryConnectionID, nil
}

//CreateL2RedundantConnection operation creates redundant layer2 connection with
//given connection structures.
//Primary connection structure is used as a baseline for underlaying API call,
//whereas secondary connection structure provices supplementary information only.
//Upon successful creation, primary connection structure, enriched with assigned UUID
//and redundant connection UUID, will be returned
func (c RestClient) CreateL2RedundantConnection(primary L2Connection, secondary L2Connection) (*string, *string, error) {
	path := "/ecx/v3/l2/connections"
	if StringValue(primary.DeviceUUID) != "" {
		path = "/ne/v1/l2/connections"
	}
	reqBody := createL2RedundantConnectionRequest(primary, secondary)
	respBody := api.CreateL2ConnectionResponse{}
	req := c.R().SetBody(&reqBody).SetResult(&respBody)
	if err := c.Execute(req, http.MethodPost, path); err != nil {
		return nil, nil, err
	}
	return respBody.PrimaryConnectionID, respBody.SecondaryConnectionID, nil
}

//DeleteL2Connection deletes layer 2 connection with a given UUID
func (c RestClient) DeleteL2Connection(uuid string) error {
	path := "/ecx/v3/l2/connections/" + url.PathEscape(uuid)
	respBody := api.DeleteL2ConnectionResponse{}
	req := c.R().SetResult(&respBody)
	if err := c.Execute(req, http.MethodDelete, path); err != nil {
		return err
	}
	return nil
}

//NewL2ConnectionUpdateRequest creates new composite update request
//for a connection with a given UUID
func (c RestClient) NewL2ConnectionUpdateRequest(uuid string) L2ConnectionUpdateRequest {
	return &restL2ConnectionUpdateRequest{
		uuid: uuid,
		c:    c,
	}
}

//WithName sets new connection name in a composite connection update request
func (req *restL2ConnectionUpdateRequest) WithName(name string) L2ConnectionUpdateRequest {
	req.name = &name
	return req
}

//WithBandwidth sets new connection bandwidth in a composite connection update request
func (req *restL2ConnectionUpdateRequest) WithBandwidth(speed int, speedUnit string) L2ConnectionUpdateRequest {
	req.speed = &speed
	req.speedUnit = &speedUnit
	return req
}

//WithSpeed sets new connection speed in a composite connection update request
func (req *restL2ConnectionUpdateRequest) WithSpeed(speed int) L2ConnectionUpdateRequest {
	req.speed = &speed
	return req
}

//WithSpeedUnit sets new connection speed unit in a composite connection update request
func (req *restL2ConnectionUpdateRequest) WithSpeedUnit(speedUnit string) L2ConnectionUpdateRequest {
	req.speedUnit = &speedUnit
	return req
}

//Execute attempts to update connection according new data set in composite update request.
//This is not atomic operation and if any update will fail, other changes won't be reverted.
//UpdateError will be returned if any of requested data failed to update
func (req *restL2ConnectionUpdateRequest) Execute() error {
	path := "/ecx/v3/l2/connections/" + url.PathEscape(req.uuid)
	reqBody := api.L2ConnectionUpdateRequest{
		Name:      req.name,
		Speed:     req.speed,
		SpeedUnit: req.speedUnit,
	}
	if StringValue(req.name) != "" || (IntValue(req.speed) > 0 && StringValue(req.speedUnit) != "") {
		restReq := req.c.R().SetQueryParam("action", "update").SetBody(&reqBody)
		if err := req.c.Execute(restReq, http.MethodPatch, path); err != nil {
			return err
		}
	}
	return nil
}

func mapGETToL2Connection(getResponse api.L2ConnectionResponse) *L2Connection {
	return &L2Connection{
		UUID:                getResponse.UUID,
		Name:                getResponse.Name,
		ProfileUUID:         getResponse.SellerServiceUUID,
		Speed:               getResponse.Speed,
		SpeedUnit:           getResponse.SpeedUnit,
		Status:              getResponse.Status,
		ProviderStatus:      getResponse.ProviderStatus,
		Notifications:       getResponse.Notifications,
		PurchaseOrderNumber: getResponse.PurchaseOrderNumber,
		PortUUID:            getResponse.PortUUID,
		DeviceUUID:          getResponse.VirtualDeviceUUID,
		VlanSTag:            getResponse.VlanSTag,
		VlanCTag:            getResponse.VlanCTag,
		NamedTag:            getResponse.NamedTag,
		AdditionalInfo:      mapAdditionalInfoAPIToDomain(getResponse.AdditionalInfo),
		ZSidePortUUID:       getResponse.ZSidePortUUID,
		ZSideVlanSTag:       getResponse.ZSideVlanSTag,
		ZSideVlanCTag:       getResponse.ZSideVlanCTag,
		SellerRegion:        getResponse.SellerRegion,
		SellerMetroCode:     getResponse.SellerMetroCode,
		AuthorizationKey:    getResponse.AuthorizationKey,
		RedundantUUID:       getResponse.RedundantUUID,
		RedundancyType:      getResponse.RedundancyType,
		Actions:             mapL2ConnectionActionsAPIToDomain(getResponse.ActionDetails),
	}
}

func createL2ConnectionRequest(l2connection L2Connection) api.L2ConnectionRequest {
	return api.L2ConnectionRequest{
		PrimaryName:          l2connection.Name,
		ProfileUUID:          l2connection.ProfileUUID,
		Speed:                l2connection.Speed,
		SpeedUnit:            l2connection.SpeedUnit,
		Notifications:        l2connection.Notifications,
		PurchaseOrderNumber:  l2connection.PurchaseOrderNumber,
		PrimaryPortUUID:      l2connection.PortUUID,
		VirtualDeviceUUID:    l2connection.DeviceUUID,
		InterfaceID:          l2connection.DeviceInterfaceID,
		PrimaryVlanSTag:      l2connection.VlanSTag,
		PrimaryVlanCTag:      l2connection.VlanCTag,
		NamedTag:             l2connection.NamedTag,
		AdditionalInfo:       mapAdditionalInfoDomainToAPI(l2connection.AdditionalInfo),
		PrimaryZSidePortUUID: l2connection.ZSidePortUUID,
		PrimaryZSideVlanSTag: l2connection.ZSideVlanSTag,
		PrimaryZSideVlanCTag: l2connection.ZSideVlanCTag,
		SellerRegion:         l2connection.SellerRegion,
		SellerMetroCode:      l2connection.SellerMetroCode,
		AuthorizationKey:     l2connection.AuthorizationKey}
}

func createL2RedundantConnectionRequest(primary L2Connection, secondary L2Connection) api.L2ConnectionRequest {
	connReq := createL2ConnectionRequest(primary)
	connReq.SecondaryName = secondary.Name
	connReq.SecondaryPortUUID = secondary.PortUUID
	if StringValue(primary.DeviceUUID) != StringValue(secondary.DeviceUUID) {
		connReq.SecondaryVirtualDeviceUUID = secondary.DeviceUUID
	}
	connReq.SecondaryVlanSTag = secondary.VlanSTag
	connReq.SecondaryVlanCTag = secondary.VlanCTag
	connReq.SecondaryZSidePortUUID = secondary.ZSidePortUUID
	connReq.SecondaryZSideVlanSTag = secondary.ZSideVlanSTag
	connReq.SecondaryZSideVlanCTag = secondary.ZSideVlanCTag
	connReq.SecondarySpeed = secondary.Speed
	connReq.SecondarySpeedUnit = secondary.SpeedUnit
	connReq.SecondaryProfileUUID = secondary.ProfileUUID
	connReq.SecondaryAuthorizationKey = secondary.AuthorizationKey
	connReq.SecondarySellerMetroCode = secondary.SellerMetroCode
	connReq.SecondarySellerRegion = secondary.SellerRegion
	connReq.SecondaryInterfaceID = secondary.DeviceInterfaceID
	return connReq
}

func mapAdditionalInfoDomainToAPI(info []L2ConnectionAdditionalInfo) []api.L2ConnectionAdditionalInfo {
	apiInfo := make([]api.L2ConnectionAdditionalInfo, len(info))
	for i, v := range info {
		apiInfo[i] = api.L2ConnectionAdditionalInfo{
			Name:  v.Name,
			Value: v.Value,
		}
	}
	return apiInfo
}

func mapAdditionalInfoAPIToDomain(apiInfo []api.L2ConnectionAdditionalInfo) []L2ConnectionAdditionalInfo {
	info := make([]L2ConnectionAdditionalInfo, len(apiInfo))
	for i, v := range apiInfo {
		info[i] = L2ConnectionAdditionalInfo{
			Name:  v.Name,
			Value: v.Value,
		}
	}
	return info
}

func mapL2ConnectionActionsAPIToDomain(apiActions []api.L2ConnectionActionDetail) []L2ConnectionAction {
	transformed := make([]L2ConnectionAction, len(apiActions))
	for i := range apiActions {
		transformed[i] = L2ConnectionAction{
			Type:         apiActions[i].ActionType,
			OperationID:  apiActions[i].OperationID,
			Message:      apiActions[i].ActionMessage,
			RequiredData: mapL2ConnectionActionDataAPIToDomain(apiActions[i].ActionRequiredData),
		}
	}
	return transformed
}

func mapL2ConnectionActionDataAPIToDomain(apiActionData []api.L2ConnectionActionRequiredData) []L2ConnectionActionData {
	transformed := make([]L2ConnectionActionData, len(apiActionData))
	for i := range apiActionData {
		transformed[i] = L2ConnectionActionData{
			Key:               apiActionData[i].Key,
			Label:             apiActionData[i].Label,
			Value:             apiActionData[i].Value,
			IsEditable:        apiActionData[i].Editable,
			ValidationPattern: apiActionData[i].ValidationPattern,
		}
	}
	return transformed
}
