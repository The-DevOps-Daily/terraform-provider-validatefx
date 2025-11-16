package validators

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure interface compliance.
var _ frameworkvalidator.String = (*uriValidator)(nil)

// URI returns a validator ensuring a string is a valid URI.
func URI() frameworkvalidator.String {
	return &uriValidator{}
}

type uriValidator struct{}

func (v *uriValidator) Description(_ context.Context) string {
	return "string must be a valid URI with a non-empty scheme and host (when required)"
}

func (v *uriValidator) MarkdownDescription(ctx context.Context) string { return v.Description(ctx) }

func (v *uriValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	s := strings.TrimSpace(req.ConfigValue.ValueString())
	u, err := url.Parse(s)
	if err != nil || u.Scheme == "" {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid URI", fmt.Sprintf("Value %q is not a valid URI: %v", s, err))
		return
	}

	// For common hierarchical schemes, require non-empty host.
	switch strings.ToLower(u.Scheme) {
	case "http", "https", "ftp", "ssh", "postgres", "postgresql", "mysql", "amqp":
		if u.Host == "" {
			resp.Diagnostics.AddAttributeError(req.Path, "Invalid URI", fmt.Sprintf("Value %q must include a host for scheme %q.", s, u.Scheme))
			return
		}
	}
}
