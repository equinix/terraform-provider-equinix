package schema_validators

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure our implementation satisfies the validator.String interface.
var _ validator.String = &StringIsJSONValidator{}

type StringIsJSONValidator struct{}

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v StringIsJSONValidator) Description(ctx context.Context) string {
    return "string must be valid JSON"
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v StringIsJSONValidator) MarkdownDescription(ctx context.Context) string {
    return "string must be valid JSON"
}

// ValidateString runs the main validation logic of the validator, reading configuration data out of `req` and updating `resp` with diagnostics.
func (v StringIsJSONValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
    // If the value is unknown or null, there is nothing to validate.
    if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
        return
    }

    var js json.RawMessage
    err := json.Unmarshal([]byte(req.ConfigValue.ValueString()), &js)
    if err != nil {
        resp.Diagnostics.AddAttributeError(
            req.Path,
            "Invalid JSON Format",
            fmt.Sprintf("String must be valid JSON: %s", err),
        )
    }
}

// StringIsJSON returns a new StringIsJSONValidator.
func StringIsJSON() StringIsJSONValidator {
    return StringIsJSONValidator{}
}
