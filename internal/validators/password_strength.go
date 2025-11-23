package validators

import (
	"context"
	"regexp"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// PasswordStrengthValidator validates passwords for minimal complexity.
func PasswordStrengthValidator() frameworkvalidator.String { return passwordStrength{} }

type passwordStrength struct{}

var _ frameworkvalidator.String = (*passwordStrength)(nil)

func (passwordStrength) Description(_ context.Context) string {
	return "value must be a strong password (min 8, upper, lower, number, special)"
}

func (v passwordStrength) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (passwordStrength) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() || req.ConfigValue.ValueString() == "" {
		return
	}

	s := req.ConfigValue.ValueString()
	if len(s) < 8 {
		resp.Diagnostics.AddAttributeError(req.Path, "Weak Password", "password must be at least 8 characters long")
		return
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(s)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(s)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(s)
	hasSpecial := regexp.MustCompile(`[!@#\$%\^&\*(),.?":{}|<>]`).MatchString(s)

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		resp.Diagnostics.AddAttributeError(req.Path, "Weak Password", "password must contain upper, lower, number, and special character")
	}
}
