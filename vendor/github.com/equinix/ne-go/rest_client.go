package ne

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/equinix/rest-go"
)

//RestClient describes REST implementation of Network Edge Client
type RestClient struct {
	*rest.Client
}

//NewClient creates new REST Network Edge client with a given baseURL, context and httpClient
func NewClient(ctx context.Context, baseURL string, httpClient *http.Client) *RestClient {
	rest := rest.NewClient(ctx, baseURL, httpClient)
	rest.SetHeader("User-agent", "equinix/ne-go")
	return &RestClient{rest}
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Unexported package methods
//_______________________________________________________________________

const (
	changeTypeCreate = "Add"
	changeTypeUpdate = "Update"
	changeTypeDelete = "Delete"
)

type headerProvider interface {
	Header() http.Header
}

func getLocationHeaderValue(provider headerProvider) (*string, error) {
	locationValues, ok := provider.Header()["Location"]
	if !ok {
		return nil, fmt.Errorf("location header not found")
	}
	if len(locationValues) != 1 {
		return nil, fmt.Errorf("only one location header value is expected")
	}
	return &locationValues[0], nil
}

func parseResourceIDFromLocationHeader(header string) (*string, error) {
	re := regexp.MustCompile(".+/([^/]+)$")
	res := re.FindAllStringSubmatch(header, -1)
	if len(res) < 1 || len(res[0]) != 2 {
		return nil, fmt.Errorf("could not parse resource identifier from location header value %q", header)
	}
	return &res[0][1], nil
}

func getResourceIDFromLocationHeader(provider headerProvider) (*string, error) {
	locHeaderValue, err := getLocationHeaderValue(provider)
	if err != nil {
		return nil, err
	}
	return parseResourceIDFromLocationHeader(*locHeaderValue)
}
