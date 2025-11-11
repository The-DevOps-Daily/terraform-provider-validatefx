package validators

import (
	"context"
	"regexp"
	"strings"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// FQDN returns a validator that ensures the value is a fully qualified domain name.
func FQDN() frameworkvalidator.String { return fqdnValidator{} }

type fqdnValidator struct{}

var _ frameworkvalidator.String = (*fqdnValidator)(nil)

// RFC-like constraints: labels 1-63 chars, alnum and hyphen, no leading/trailing hyphen.
var (
	fqdnLabelRe = regexp.MustCompile(`^[A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?$`)
)

func (fqdnValidator) Description(_ context.Context) string {
	return "value must be a fully qualified domain name"
}

func (v fqdnValidator) MarkdownDescription(ctx context.Context) string { return v.Description(ctx) }

func (fqdnValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	raw := strings.TrimSpace(req.ConfigValue.ValueString())
	if raw == "" {
		return
	}

	// Must have at least one dot and no empty labels
	parts := strings.Split(raw, ".")
	if len(parts) < 2 {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid FQDN", "FQDN must contain at least one dot (e.g., example.com)")
		return
	}

	if len(raw) > 253 {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid FQDN", "FQDN is too long (maximum 253 characters)")
		return
	}

	for _, label := range parts {
		if label == "" {
			resp.Diagnostics.AddAttributeError(req.Path, "Invalid FQDN", "FQDN must not contain empty labels")
			return
		}
		if !fqdnLabelRe.MatchString(label) {
			resp.Diagnostics.AddAttributeError(req.Path, "Invalid FQDN", "Each label must start/end with alphanumeric and contain only letters, digits, or hyphens with length 1-63")
			return
		}
	}
}
