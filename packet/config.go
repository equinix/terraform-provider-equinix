package packet

import (
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-sdk/helper/logging"
	"github.com/packethost/packngo"
)

const (
	consumerToken = "aZ9GmqHTPtxevvFq9SK3Pi2yr9YCbRzduCSXF2SNem5sjB91mDq7Th3ZwTtRqMWZ"
)

type Config struct {
	AuthToken string
}

// Client returns a new client for accessing Packet's API.
func (c *Config) Client() *packngo.Client {
	httpClient := retryablehttp.NewClient()
	httpClient.RetryWaitMin = time.Second
	httpClient.RetryWaitMax = 30 * time.Second
	httpClient.RetryMax = 10
	httpClient.HTTPClient.Transport = logging.NewTransport(
		"Packet",
		httpClient.HTTPClient.Transport)

	return packngo.NewClientWithAuth(consumerToken, c.AuthToken, httpClient)
}
