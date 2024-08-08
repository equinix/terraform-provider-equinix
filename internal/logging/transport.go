package logging

import (
	"net/http"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
)

type HTTPLoggingTransport struct {
	subsystem       string
	transport       http.RoundTripper
	sensitiveFields []string
	level           hclog.Level
}

var _ http.RoundTripper = (*HTTPLoggingTransport)(nil)

func NewHTTPLoggingTransport(subsystem string, transport http.RoundTripper, sensitiveFields []string, level hclog.Level) *HTTPLoggingTransport {
	return &HTTPLoggingTransport{
		subsystem:       subsystem,
		transport:       logging.NewSubsystemLoggingHTTPTransport(subsystem, transport),
		sensitiveFields: sensitiveFields,
		level:           level,
	}
}

func (t *HTTPLoggingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	ctx := tflog.NewSubsystem(r.Context(), t.subsystem, tflog.WithLevel(t.level))
	ctx = tflog.SubsystemMaskFieldValuesWithFieldKeys(ctx, t.subsystem, t.sensitiveFields...)

	requestWithSubsystem := r.WithContext(ctx)

	return t.transport.RoundTrip(requestWithSubsystem)
}
