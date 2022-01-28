package ne

import (
	"net/http"
	"net/url"

	"github.com/equinix/ne-go/internal/api"
)

type restBGPConfigurationUpdateRequest struct {
	uuid              string
	localIPAddress    *string
	localASN          *int
	remoteIPAddress   *string
	remoteASN         *int
	authenticationKey *string
	c                 RestClient
}

//CreateBGPConfiguration creates new Network Edge BGP configuration
//with a given model. Configuration's UUID is returned on successful
//creation
func (c RestClient) CreateBGPConfiguration(config BGPConfiguration) (*string, error) {
	path := "/ne/v1/bgp"
	reqBody := mapBGPConfigurationDomainToAPI(config)
	respBody := api.BGPConfigurationCreateResponse{}
	req := c.R().SetBody(&reqBody).SetResult(&respBody)
	if err := c.Execute(req, http.MethodPost, path); err != nil {
		return nil, err
	}
	return respBody.UUID, nil
}

//GetBGPConfiguration retrieves BGP configuration with a given UUID
func (c RestClient) GetBGPConfiguration(uuid string) (*BGPConfiguration, error) {
	path := "/ne/v1/bgp/" + url.PathEscape(uuid)
	respBody := api.BGPConfiguration{}
	req := c.R().SetResult(&respBody)
	if err := c.Execute(req, http.MethodGet, path); err != nil {
		return nil, err
	}
	return mapBGPConfigurationAPIToDomain(respBody), nil
}

//GetBGPConfigurationForConnection retreive BGP configuration for
//a connection with a given connection UUID
func (c RestClient) GetBGPConfigurationForConnection(uuid string) (*BGPConfiguration, error) {
	path := "/ne/v1/bgp/connection/" + url.PathEscape(uuid)
	respBody := api.BGPConfiguration{}
	req := c.R().SetResult(&respBody)
	if err := c.Execute(req, http.MethodGet, path); err != nil {
		return nil, err
	}
	return mapBGPConfigurationAPIToDomain(respBody), nil
}

//NewBGPConfigurationUpdateRequest creates new BGP configuration update
//request for a configuration with given UUID
func (c RestClient) NewBGPConfigurationUpdateRequest(uuid string) BGPUpdateRequest {
	return &restBGPConfigurationUpdateRequest{
		uuid: uuid,
		c:    c,
	}
}

func (req *restBGPConfigurationUpdateRequest) WithLocalIPAddress(localIPAddress string) BGPUpdateRequest {
	req.localIPAddress = &localIPAddress
	return req
}

func (req *restBGPConfigurationUpdateRequest) WithLocalASN(localASN int) BGPUpdateRequest {
	req.localASN = &localASN
	return req
}

func (req *restBGPConfigurationUpdateRequest) WithRemoteASN(remoteASN int) BGPUpdateRequest {
	req.remoteASN = &remoteASN
	return req
}

func (req *restBGPConfigurationUpdateRequest) WithRemoteIPAddress(remoteIPAddress string) BGPUpdateRequest {
	req.remoteIPAddress = &remoteIPAddress
	return req
}

func (req *restBGPConfigurationUpdateRequest) WithAuthenticationKey(authenticationKey string) BGPUpdateRequest {
	req.authenticationKey = &authenticationKey
	return req
}

func (req *restBGPConfigurationUpdateRequest) Execute() error {
	path := "/ne/v1/bgp/" + url.PathEscape(req.uuid)
	reqBody := api.BGPConfiguration{
		LocalIPAddress:    req.localIPAddress,
		LocalASN:          req.localASN,
		RemoteIPAddress:   req.remoteIPAddress,
		RemoteASN:         req.remoteASN,
		AuthenticationKey: req.authenticationKey,
	}
	respBody := api.BGPConfigurationCreateResponse{}
	restReq := req.c.R().SetBody(&reqBody).SetResult(&respBody)
	if err := req.c.Execute(restReq, http.MethodPut, path); err != nil {
		return err
	}
	return nil
}

func mapBGPConfigurationDomainToAPI(config BGPConfiguration) api.BGPConfiguration {
	return api.BGPConfiguration{
		UUID:              config.UUID,
		ConnectionUUID:    config.ConnectionUUID,
		LocalIPAddress:    config.LocalIPAddress,
		LocalASN:          config.LocalASN,
		RemoteASN:         config.RemoteASN,
		RemoteIPAddress:   config.RemoteIPAddress,
		AuthenticationKey: config.AuthenticationKey,
	}
}

func mapBGPConfigurationAPIToDomain(apiConfig api.BGPConfiguration) *BGPConfiguration {
	return &BGPConfiguration{
		UUID:               apiConfig.UUID,
		ConnectionUUID:     apiConfig.ConnectionUUID,
		DeviceUUID:         apiConfig.VirtualDeviceUUID,
		LocalIPAddress:     apiConfig.LocalIPAddress,
		LocalASN:           apiConfig.LocalASN,
		RemoteIPAddress:    apiConfig.RemoteIPAddress,
		RemoteASN:          apiConfig.RemoteASN,
		AuthenticationKey:  apiConfig.AuthenticationKey,
		State:              apiConfig.State,
		ProvisioningStatus: apiConfig.ProvisioningStatus,
	}
}
