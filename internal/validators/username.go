package validators

import (
	"context"
	"fmt"
	"unicode"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

const (
	defaultUsernameMinLength = 3
	defaultUsernameMaxLength = 20
)

// Ensure interface compliance.
var _ frameworkvalidator.String = (*usernameValidator)(nil)

// Username returns a validator enforcing ValidateFX username rules.
// Callers may customize the length bounds; characters are limited to
// ASCII letters, digits, and underscore to align with common username rules.
func Username(minLength, maxLength int) frameworkvalidator.String {
	return &usernameValidator{
		minLength:   minLength,
		maxLength:   maxLength,
		description: fmt.Sprintf("%d-%d characters using letters, digits, or underscores", minLength, maxLength),
	}
}

type usernameValidator struct {
	minLength   int
	maxLength   int
	description string
}

func (v *usernameValidator) Description(_ context.Context) string {
	return v.description
}

func (v *usernameValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *usernameValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	runeCount := len([]rune(value))

	if runeCount < v.minLength || runeCount > v.maxLength || !isValidUsername(value) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Username",
			fmt.Sprintf("Username must be %s.", v.description),
		)
	}
}

// DefaultUsernameValidator returns the project defaults.
func DefaultUsernameValidator() frameworkvalidator.String {
	return Username(defaultUsernameMinLength, defaultUsernameMaxLength)
}

func isValidUsername(value string) bool {
	if value == "" {
		return false
	}

	for _, r := range value {
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r):
			continue
		case r == '_':
			continue
		default:
			return false
		}
	}

	return true
}
