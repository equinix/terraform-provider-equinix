package internetaccess

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
)

// getService retrieves an Internet Access service.
//
// The SDK model marks "order" as required, but the API leaves it out while a service is still
// being provisioned, which makes the SDK fail to decode an otherwise valid response. In that case
// the response body is decoded again with an empty "order" so that callers can see the state.
func getService(ctx context.Context, client *fabricv4.APIClient, uuid string) (*fabricv4.InternetAccessService, *http.Response, error) {
	svc, resp, err := client.InternetAccessServicesApi.GetEiaService(ctx, uuid).Execute()
	if err == nil {
		return svc, resp, nil
	}

	if resp == nil || resp.StatusCode != http.StatusOK {
		return svc, resp, err
	}
	var apiErr *fabricv4.GenericOpenAPIError
	if !errors.As(err, &apiErr) {
		return svc, resp, err
	}
	repaired, repairErr := decodeWithDefaultOrder(apiErr.Body())
	if repairErr != nil {
		return svc, resp, err
	}
	return repaired, resp, nil
}

func decodeWithDefaultOrder(body []byte) (*fabricv4.InternetAccessService, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	if _, ok := raw["order"]; ok {
		return nil, errors.New("order is present, the decode error has another cause")
	}
	raw["order"] = json.RawMessage(`{}`)

	fixed, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}
	var svc fabricv4.InternetAccessService
	if err := json.Unmarshal(fixed, &svc); err != nil {
		return nil, err
	}
	return &svc, nil
}
