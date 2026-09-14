package internetaccess

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
)

// searchInternetAccessServicesRaw makes the POST /fabric/v4/internetAccessServices/search call
// without going through the SDK's strict UnmarshalJSON, which requires "order" to be present in
// every InternetAccessService item even though the API does not always return it.
func searchInternetAccessServicesRaw(ctx context.Context, client *fabricv4.APIClient, req fabricv4.InternetAccessSearchRequest) (*fabricv4.InternetAccessServices, error) {
	baseURL := client.GetConfig().Servers[0].URL
	url := baseURL + "/fabric/v4/internetAccessServices/search"

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("X-SOURCE", "API")

	resp, err := client.GetConfig().HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%s", string(respBody))
	}

	patched, err := patchMissingOrderInSearchResponse(respBody)
	if err != nil {
		return nil, fmt.Errorf("patch response: %w", err)
	}

	var services fabricv4.InternetAccessServices
	if err := json.Unmarshal(patched, &services); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return &services, nil
}

// getInternetAccessServiceRaw retrieves a single EIA service by UUID without going through
// the SDK's strict UnmarshalJSON that requires the "order" field.
func getInternetAccessServiceRaw(ctx context.Context, client *fabricv4.APIClient, id string) (*fabricv4.InternetAccessService, error) {
	baseURL := client.GetConfig().Servers[0].URL
	url := baseURL + "/fabric/v4/internetAccessServices/" + id

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("X-SOURCE", "API")

	resp, err := client.GetConfig().HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("not found")
	}
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%s", string(respBody))
	}

	patched := patchOrderField(respBody)

	var svc fabricv4.InternetAccessService
	if err := json.Unmarshal(patched, &svc); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return &svc, nil
}

// patchMissingOrderInSearchResponse adds "order":{} to each item in the "data" array
// if that field is absent. The fabricv4.InternetAccessService.UnmarshalJSON requires it.
func patchMissingOrderInSearchResponse(data []byte) ([]byte, error) {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(data, &m); err != nil {
		return data, nil
	}

	dataField, ok := m["data"]
	if !ok {
		return data, nil
	}

	var items []json.RawMessage
	if err := json.Unmarshal(dataField, &items); err != nil {
		return data, nil
	}

	for i, item := range items {
		items[i] = patchOrderField(item)
	}

	patchedItems, err := json.Marshal(items)
	if err != nil {
		return data, nil
	}
	m["data"] = patchedItems

	patched, err := json.Marshal(m)
	if err != nil {
		return data, nil
	}
	return patched, nil
}

// patchOrderField adds "order":{} to the JSON object if the key is absent.
func patchOrderField(data []byte) []byte {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(data, &m); err != nil {
		return data
	}
	if _, ok := m["order"]; !ok {
		m["order"] = json.RawMessage(`{}`)
	}
	patched, err := json.Marshal(m)
	if err != nil {
		return data
	}
	return patched
}
