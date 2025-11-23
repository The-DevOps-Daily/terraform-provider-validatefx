package validators

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ frameworkvalidator.String = URL()

// URL returns a schema.String validator that ensures the value is a well formed URL with scheme and host.
func URL() frameworkvalidator.String {
	return urlValidator{}
}

type urlValidator struct{}

func (urlValidator) Description(context.Context) string {
	return "value must be a valid URL including scheme and host"
}

func (v urlValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (urlValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := strings.TrimSpace(req.ConfigValue.ValueString())
	if value == "" {
		return
	}

	parsed, err := url.ParseRequestURI(value)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid URL",
			fmt.Sprintf("Value %q is not a valid URL including scheme and host", value),
		)
		return
	}

	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Unsupported URL Scheme",
			fmt.Sprintf("Value %q uses unsupported URL scheme %q. Only http and https are permitted.", value, parsed.Scheme),
		)
	}
}
