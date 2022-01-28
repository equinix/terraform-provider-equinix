package ne

import (
	"net/http"
	"net/url"

	"github.com/equinix/ne-go/internal/api"
)

//GetSSHPublicKeys retrieves list of available SSH public keys
func (c RestClient) GetSSHPublicKeys() ([]SSHPublicKey, error) {
	path := "/ne/v1/publicKeys"
	respBody := make([]api.SSHPublicKey, 0)
	req := c.R().SetResult(&respBody)
	if err := c.Execute(req, http.MethodGet, path); err != nil {
		return nil, err
	}
	return mapSSHPublicKeysAPIToDomain(respBody), nil
}

//GetSSHPublicKey retrieves SSH public key with a given identifier
func (c RestClient) GetSSHPublicKey(uuid string) (*SSHPublicKey, error) {
	path := "/ne/v1/publicKeys/" + url.PathEscape(uuid)
	respBody := api.SSHPublicKey{}
	req := c.R().SetResult(&respBody)
	if err := c.Execute(req, http.MethodGet, path); err != nil {
		return nil, err
	}
	mapped := mapSSHPublicKeyAPIToDomain(respBody)
	return &mapped, nil
}

//CreateSSHPublicKey creates new SSH public key with a given details
func (c RestClient) CreateSSHPublicKey(key SSHPublicKey) (*string, error) {
	path := "/ne/v1/publicKeys"
	reqBody := mapSSHPublicKeyDomainToAPI(key)
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

//DeleteSSHPublicKey removes SSH Public key with given identifier
func (c RestClient) DeleteSSHPublicKey(uuid string) error {
	path := "/ne/v1/publicKeys/" + url.PathEscape(uuid)
	if err := c.Execute(c.R(), http.MethodDelete, path); err != nil {
		return err
	}
	return nil
}

func mapSSHPublicKeysAPIToDomain(apiKeys []api.SSHPublicKey) []SSHPublicKey {
	transformed := make([]SSHPublicKey, len(apiKeys))
	for i := range apiKeys {
		transformed[i] = mapSSHPublicKeyAPIToDomain(apiKeys[i])
	}
	return transformed
}

func mapSSHPublicKeyAPIToDomain(apiKey api.SSHPublicKey) SSHPublicKey {
	return SSHPublicKey{
		UUID:  apiKey.UUID,
		Name:  apiKey.KeyName,
		Value: apiKey.KeyValue,
	}
}

func mapSSHPublicKeyDomainToAPI(key SSHPublicKey) api.SSHPublicKey {
	return api.SSHPublicKey{
		UUID:     key.UUID,
		KeyName:  key.Name,
		KeyValue: key.Value,
	}
}
