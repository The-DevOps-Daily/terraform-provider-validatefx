package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestUsernameValidatorValid(t *testing.T) {
	t.Parallel()

	validator := DefaultUsernameValidator()
	req := frameworkvalidator.StringRequest{
		Path:        path.Root("username"),
		ConfigValue: types.StringValue("user_123"),
	}
	resp := &frameworkvalidator.StringResponse{}

	validator.ValidateString(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no diagnostics for valid username: %v", resp.Diagnostics)
	}
}

func TestUsernameValidatorInvalid(t *testing.T) {
	t.Parallel()

	validator := DefaultUsernameValidator()
	req := frameworkvalidator.StringRequest{
		Path:        path.Root("username"),
		ConfigValue: types.StringValue("invalid-user!"),
	}
	resp := &frameworkvalidator.StringResponse{}

	validator.ValidateString(context.Background(), req, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatalf("expected diagnostics for invalid username")
	}

	diagnostic := resp.Diagnostics[0]

	if diagnostic.Severity() != diag.SeverityError {
		t.Fatalf("expected error diagnostic, got severity: %s", diagnostic.Severity())
	}

	if diagnostic.Summary() != "Invalid Username" {
		t.Fatalf("unexpected diagnostic summary: %s", diagnostic.Summary())
	}
}

func TestUsernameValidatorHandlesNullAndUnknown(t *testing.T) {
	t.Parallel()

	validator := DefaultUsernameValidator()

	testCases := map[string]frameworkvalidator.StringRequest{
		"null": {
			Path:        path.Root("username"),
			ConfigValue: types.StringNull(),
		},
		"unknown": {
			Path:        path.Root("username"),
			ConfigValue: types.StringUnknown(),
		},
	}

	for name, req := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &frameworkvalidator.StringResponse{}
			validator.ValidateString(context.Background(), req, resp)

			if resp.Diagnostics.HasError() {
				t.Fatalf("expected no diagnostics for %s input", name)
			}
		})
	}
}

func TestUsernameValidatorLengthBounds(t *testing.T) {
	t.Parallel()

	custom := Username(5, 8)

	tests := []struct {
		name        string
		value       string
		shouldError bool
	}{
		{name: "too short", value: "user", shouldError: true},
		{name: "too long", value: "user_name", shouldError: true},
		{name: "within bounds", value: "user123", shouldError: false},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := frameworkvalidator.StringRequest{
				Path:        path.Root("username"),
				ConfigValue: types.StringValue(tc.value),
			}
			resp := &frameworkvalidator.StringResponse{}

			custom.ValidateString(context.Background(), req, resp)

			if tc.shouldError && !resp.Diagnostics.HasError() {
				t.Fatalf("expected error for value %q", tc.value)
			}

			if !tc.shouldError && resp.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostic for value %q: %v", tc.value, resp.Diagnostics)
			}
		})
	}
}
