package validators

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestBetweenValidator(t *testing.T) {
	t.Parallel()

	validator := Between("1.5", "10")

	testCases := map[string]struct {
		value       types.String
		expectError bool
		message     string
	}{
		"within range": {
			value:       types.StringValue("5"),
			expectError: false,
		},
		"below min": {
			value:       types.StringValue("1.4"),
			expectError: true,
			message:     "less than minimum",
		},
		"above max": {
			value:       types.StringValue("11"),
			expectError: true,
			message:     "greater than maximum",
		},
		"non numeric": {
			value:       types.StringValue("abc"),
			expectError: true,
			message:     "not a valid decimal",
		},
		"empty": {
			value:       types.StringValue(""),
			expectError: false,
		},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := frameworkvalidator.StringRequest{
				Path:        path.Root("value"),
				ConfigValue: tc.value,
			}
			resp := &frameworkvalidator.StringResponse{}

			validator.ValidateString(context.Background(), req, resp)

			if resp.Diagnostics.HasError() != tc.expectError {
				t.Fatalf("expected error=%t, got=%t", tc.expectError, resp.Diagnostics.HasError())
			}

			if tc.expectError && tc.message != "" {
				if !containsDetail(resp, tc.message) {
					t.Fatalf("expected diagnostic to contain %q, diagnostics=%v", tc.message, resp.Diagnostics)
				}
			}
		})
	}
}

func TestBetweenValidatorHandlesNullUnknown(t *testing.T) {
	t.Parallel()

	validator := Between("0", "1")

	cases := map[string]types.String{
		"null":    types.StringNull(),
		"unknown": types.StringUnknown(),
	}

	for name, val := range cases {
		name, val := name, val
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := frameworkvalidator.StringRequest{
				Path:        path.Root("value"),
				ConfigValue: val,
			}
			resp := &frameworkvalidator.StringResponse{}

			validator.ValidateString(context.Background(), req, resp)

			if resp.Diagnostics.HasError() {
				t.Fatalf("expected no diagnostics for %s", name)
			}
		})
	}
}

func TestBetweenValidatorInvalidBounds(t *testing.T) {
	t.Parallel()

	validator := Between("invalid", "1")

	req := frameworkvalidator.StringRequest{
		Path:        path.Root("value"),
		ConfigValue: types.StringValue("0"),
	}
	resp := &frameworkvalidator.StringResponse{}

	validator.ValidateString(context.Background(), req, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatalf("expected diagnostics for invalid min")
	}
}

func containsDetail(resp *frameworkvalidator.StringResponse, needle string) bool {
	for _, diag := range resp.Diagnostics {
		if strings.Contains(diag.Detail(), needle) {
			return true
		}
	}
	return false
}
