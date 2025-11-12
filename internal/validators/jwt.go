package validators

import (
	"context"
	"encoding/base64"
	"strings"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// JWT validates that a string is a well-formed JSON Web Token (three base64url segments).
func JWT() frameworkvalidator.String { return jwtValidator{} }

type jwtValidator struct{}

var _ frameworkvalidator.String = (*jwtValidator)(nil)

func (jwtValidator) Description(_ context.Context) string {
	return "value must be a well-formed JWT (three base64url segments)"
}

func (v jwtValidator) MarkdownDescription(ctx context.Context) string { return v.Description(ctx) }

func (jwtValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	raw := strings.TrimSpace(req.ConfigValue.ValueString())
	if raw == "" {
		return
	}

	parts := strings.Split(raw, ".")
	if len(parts) != 3 {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid JWT", "JWT must have three segments separated by dots")
		return
	}

	for _, p := range parts {
		if p == "" {
			resp.Diagnostics.AddAttributeError(req.Path, "Invalid JWT", "JWT contains empty segment")
			return
		}
		// base64url decode without padding
		if _, err := base64.RawURLEncoding.DecodeString(p); err != nil {
			resp.Diagnostics.AddAttributeError(req.Path, "Invalid JWT", "JWT segments must be base64url encoded")
			return
		}
	}
}
