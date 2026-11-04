package internetaccess

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
)

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

func searchServices(ctx context.Context, client *fabricv4.APIClient, request fabricv4.InternetAccessSearchRequest) (*fabricv4.InternetAccessServices, *http.Response, error) {
	services, resp, err := client.InternetAccessServicesApi.SearchEiaServices(ctx).InternetAccessSearchRequest(request).Execute()
	if err == nil {
		return services, resp, nil
	}

	if resp == nil || resp.StatusCode != http.StatusOK {
		return services, resp, err
	}
	var apiErr *fabricv4.GenericOpenAPIError
	if !errors.As(err, &apiErr) {
		return services, resp, err
	}
	repaired, repairErr := decodeSearchWithDefaultOrder(apiErr.Body())
	if repairErr != nil {
		return services, resp, err
	}
	return repaired, resp, nil
}

func decodeSearchWithDefaultOrder(body []byte) (*fabricv4.InternetAccessServices, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	var items []map[string]json.RawMessage
	if err := json.Unmarshal(raw["data"], &items); err != nil {
		return nil, err
	}
	for _, item := range items {
		if _, ok := item["order"]; !ok {
			item["order"] = json.RawMessage(`{}`)
		}
	}
	data, err := json.Marshal(items)
	if err != nil {
		return nil, err
	}
	raw["data"] = data

	fixed, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}
	var services fabricv4.InternetAccessServices
	if err := json.Unmarshal(fixed, &services); err != nil {
		return nil, err
	}
	return &services, nil
}
