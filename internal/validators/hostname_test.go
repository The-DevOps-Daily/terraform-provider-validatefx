package validators

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestHostnameValidator(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		value       types.String
		expectError bool
	}{
		"simple hostname": {
			value:       types.StringValue("example"),
			expectError: false,
		},
		"multi label": {
			value:       types.StringValue("api.service.local"),
			expectError: false,
		},
		"punycode label": {
			value:       types.StringValue("xn--bcher-kva.example"),
			expectError: false,
		},
		"trailing dot": {
			value:       types.StringValue("example.com."),
			expectError: false,
		},
		"single label with hyphen middle": {
			value:       types.StringValue("foo-bar"),
			expectError: false,
		},
		"numeric label": {
			value:       types.StringValue("12345"),
			expectError: false,
		},
		"underscore": {
			value:       types.StringValue("bad_name"),
			expectError: true,
		},
		"label starting hyphen": {
			value:       types.StringValue("-bad.example"),
			expectError: true,
		},
		"label ending hyphen": {
			value:       types.StringValue("bad-.example"),
			expectError: true,
		},
		"double dot": {
			value:       types.StringValue("bad..example"),
			expectError: true,
		},
		"empty": {
			value:       types.StringValue(""),
			expectError: false,
		},
		"label too long": {
			value:       types.StringValue(strings.Repeat("a", 64)),
			expectError: true,
		},
		"hostname too long": {
			value:       types.StringValue(strings.Repeat("a.", 127) + "a"),
			expectError: true,
		},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := frameworkvalidator.StringRequest{
				Path:        path.Root("hostname"),
				ConfigValue: tc.value,
			}
			resp := &frameworkvalidator.StringResponse{}

			Hostname().ValidateString(context.Background(), req, resp)

			hasError := resp.Diagnostics.HasError()
			if hasError != tc.expectError {
				t.Fatalf("expected error=%t, got=%t diagnostics=%v", tc.expectError, hasError, resp.Diagnostics)
			}
		})
	}
}

func TestHostnameValidatorNullUnknown(t *testing.T) {
	t.Parallel()

	cases := map[string]types.String{
		"null":    types.StringNull(),
		"unknown": types.StringUnknown(),
	}

	for name, val := range cases {
		name, val := name, val
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := frameworkvalidator.StringRequest{
				Path:        path.Root("hostname"),
				ConfigValue: val,
			}
			resp := &frameworkvalidator.StringResponse{}

			Hostname().ValidateString(context.Background(), req, resp)

			if resp.Diagnostics.HasError() {
				t.Fatalf("expected no diagnostics for %s, got %v", name, resp.Diagnostics)
			}
		})
	}
}
