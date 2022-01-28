package ne

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/equinix/ne-go/internal/api"
	"github.com/equinix/rest-go"
)

const (
	associateDevice   = "ADD"
	unassociateDevice = "DELETE"
)

type restSSHUserUpdateRequest struct {
	uuid        string
	newPassword string
	oldDevices  []string
	newDevices  []string
	c           RestClient
}

//CreateSSHUser creates new Network Edge SSH user with a given parameters and returns its UUID upon successful creation
func (c RestClient) CreateSSHUser(username string, password string, device string) (*string, error) {
	path := "/ne/v1/sshUsers"
	reqBody := api.SSHUserRequest{
		Username:   &username,
		Password:   &password,
		DeviceUUID: &device,
	}
	req := c.R().SetBody(&reqBody)
	resp, err := c.Do(http.MethodPost, path, req)
	if err != nil {
		return nil, err
	}
	uuid, err := getResourceIDFromLocationHeader(resp)
	if err != nil {
		return nil, err
	}
	return uuid, nil
}

//GetSSHUsers retrieves list of all SSH users (with details)
func (c RestClient) GetSSHUsers() ([]SSHUser, error) {
	path := "/ne/v1/sshUsers"
	content, err := c.GetOffsetPaginated(path, &api.SSHUsersResponse{},
		rest.DefaultOffsetPagingConfig().
			SetAdditionalParams(map[string]string{"verbose": "true"}))
	if err != nil {
		return nil, err
	}
	transformed := make([]SSHUser, len(content))
	for i := range content {
		transformed[i] = *mapSSHUserAPIToDomain(content[i].(api.SSHUser))
	}
	return transformed, nil
}

//GetSSHUser fetches details of a SSH user with a given UUID
func (c RestClient) GetSSHUser(uuid string) (*SSHUser, error) {
	path := "/ne/v1/sshUsers/" + url.PathEscape(uuid)
	respBody := api.SSHUser{}
	req := c.R().SetResult(&respBody)
	if err := c.Execute(req, http.MethodGet, path); err != nil {
		return nil, err
	}
	return mapSSHUserAPIToDomain(respBody), nil
}

//NewSSHUserUpdateRequest creates new composite update request for a user with a given UUID
func (c RestClient) NewSSHUserUpdateRequest(uuid string) SSHUserUpdateRequest {
	return &restSSHUserUpdateRequest{
		uuid: uuid,
		c:    c}
}

//DeleteSSHUser deletes ssh user with a given UUID
func (c RestClient) DeleteSSHUser(uuid string) error {
	user, err := c.GetSSHUser(uuid)
	if err != nil {
		return err
	}
	updateErr := UpdateError{}
	for _, dev := range user.DeviceUUIDs {
		if err := c.changeDeviceAssociation(unassociateDevice, uuid, dev); err != nil {
			updateErr.AddChangeError(changeTypeDelete, "devices", dev, err)
		}
	}
	if updateErr.ChangeErrorsCount() > 0 {
		return updateErr
	}
	return nil
}

func (req *restSSHUserUpdateRequest) WithNewPassword(password string) SSHUserUpdateRequest {
	req.newPassword = password
	return req
}

func (req *restSSHUserUpdateRequest) WithDeviceChange(old []string, new []string) SSHUserUpdateRequest {
	req.oldDevices = old
	req.newDevices = new
	return req
}

func (req *restSSHUserUpdateRequest) Execute() error {
	updateErr := UpdateError{}
	if req.newPassword != "" {
		if err := req.c.changeUserPassword(req.uuid, req.newPassword); err != nil {
			updateErr.AddChangeError(changeTypeUpdate, "password", req.newPassword, err)
		}
	}
	removed, added := diffStringSlices(req.oldDevices, req.newDevices)
	for _, dev := range added {
		if err := req.c.changeDeviceAssociation(associateDevice, req.uuid, dev); err != nil {
			updateErr.AddChangeError(changeTypeCreate, "devices", dev, err)
		}
	}
	for _, dev := range removed {
		if err := req.c.changeDeviceAssociation(unassociateDevice, req.uuid, dev); err != nil {
			updateErr.AddChangeError(changeTypeDelete, "devices", dev, err)
		}
	}
	if updateErr.ChangeErrorsCount() > 0 {
		return updateErr
	}
	return nil
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Unexported package methods
//_______________________________________________________________________

func (c RestClient) changeUserPassword(userID string, newPassword string) error {
	path := "/ne/v1/sshUsers/" + url.PathEscape(userID)
	reqBody := api.SSHUserUpdateRequest{Password: &newPassword}
	req := c.R().SetBody(&reqBody)
	if err := c.Execute(req, http.MethodPut, path); err != nil {
		return err
	}
	return nil
}

func (c RestClient) changeDeviceAssociation(changeType string, userID string, deviceID string) error {
	path := fmt.Sprintf("/ne/v1/sshUsers/%s/devices/%s",
		url.PathEscape(userID), url.PathEscape(deviceID))
	var method string
	switch changeType {
	case associateDevice:
		method = http.MethodPost
	case unassociateDevice:
		method = http.MethodDelete
	default:
		return fmt.Errorf("unsupported association change type")
	}
	req := c.R().
		//due to bug in NE API that requires content type and content len = 0 altough there is no content needed in any case
		SetHeader("Content-Type", "application/json").
		SetBody("{}")
	if err := c.Execute(req, method, path); err != nil {
		return err
	}
	return nil
}

func mapSSHUserAPIToDomain(apiUser api.SSHUser) *SSHUser {
	return &SSHUser{
		UUID:        apiUser.UUID,
		Username:    apiUser.Username,
		DeviceUUIDs: apiUser.DeviceUUIDs}
}

func diffStringSlices(a, b []string) (extraA, extraB []string) {
	visited := make([]bool, len(b))
	for i := range a {
		found := false
		for j := range b {
			if visited[j] {
				continue
			}
			if a[i] == b[j] {
				visited[j] = true
				found = true
				break
			}
		}
		if !found {
			extraA = append(extraA, a[i])
		}
	}
	for j := range b {
		if visited[j] {
			continue
		}
		extraB = append(extraB, b[j])
	}
	return
}
