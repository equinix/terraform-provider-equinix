package ecx

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/equinix/rest-go"
)

//RestClient describes Equinix Fabric client that uses REST API
type RestClient struct {
	*rest.Client
}

//NewClient creates new Equinix Fabric REST API client with a given baseURL and http.Client
func NewClient(ctx context.Context, baseURL string, httpClient *http.Client) *RestClient {
	rest := rest.NewClient(ctx, baseURL, httpClient)
	rest.SetHeader("User-agent", "equinix/ecx-go")
	return &RestClient{rest}
}

func buildQueryParamValueString(values []string) string {
	var sb strings.Builder
	for i := range values {
		sb.WriteString(url.QueryEscape(values[i]))
		if i < len(values)-1 {
			sb.WriteString(",")
		}
	}
	return sb.String()
}
