package validators

import (
	"context"
	"testing"

	frameworkdiag "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestMACAddressValidatorValid(t *testing.T) {
	t.Parallel()

	testCases := []string{
		"00:1A:2b:3C:4d:5e",
		"00-1A-2B-3C-4D-5E",
		"001a2b3c4d5e",
	}

	validator := MACAddress()

	for _, tc := range testCases {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			t.Parallel()

			req := frameworkvalidator.StringRequest{
				Path:        path.Root("mac_address"),
				ConfigValue: types.StringValue(tc),
			}

			resp := &frameworkvalidator.StringResponse{}

			validator.ValidateString(context.Background(), req, resp)

			if resp.Diagnostics.HasError() {
				t.Fatalf("expected no diagnostics for valid mac %q, got: %v", tc, resp.Diagnostics)
			}
		})
	}
}

func TestMACAddressValidatorInvalid(t *testing.T) {
	t.Parallel()

	testCases := map[string]string{
		"too_short":        "00:11:22:33:44",
		"too_long":         "00:11:22:33:44:55:66",
		"invalid_chars":    "00:ZZ:22:33:44:55",
		"mixed_separators": "00:11-22:33:44:55",
		"bad_compact":      "00112233445G",
	}

	validator := MACAddress()

	for name, value := range testCases {
		name := name
		value := value
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := frameworkvalidator.StringRequest{
				Path:        path.Root("mac_address"),
				ConfigValue: types.StringValue(value),
			}
			resp := &frameworkvalidator.StringResponse{}

			validator.ValidateString(context.Background(), req, resp)

			if !resp.Diagnostics.HasError() {
				t.Fatalf("expected diagnostics for invalid case %q", name)
			}

			diag := resp.Diagnostics[0]
			if diag.Severity() != frameworkdiag.SeverityError {
				t.Fatalf("expected error severity for case %q, got %s", name, diag.Severity())
			}

			if diag.Summary() != "Invalid MAC Address" {
				t.Fatalf("unexpected summary for case %q: %s", name, diag.Summary())
			}
		})
	}
}

func TestMACAddressValidatorHandlesEmptyAndNull(t *testing.T) {
	t.Parallel()

	testCases := map[string]frameworkvalidator.StringRequest{
		"empty": {
			Path:        path.Root("mac_address"),
			ConfigValue: types.StringValue(""),
		},
		"null": {
			Path:        path.Root("mac_address"),
			ConfigValue: types.StringNull(),
		},
		"unknown": {
			Path:        path.Root("mac_address"),
			ConfigValue: types.StringUnknown(),
		},
	}

	validator := MACAddress()

	for name, req := range testCases {
		name := name
		req := req
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &frameworkvalidator.StringResponse{}
			validator.ValidateString(context.Background(), req, resp)

			if resp.Diagnostics.HasError() {
				t.Fatalf("expected no diagnostics for case %q, got: %v", name, resp.Diagnostics)
			}
		})
	}
}

func TestNormalizeMACAddress(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		input       string
		expected    string
		shouldError bool
	}{
		"colon format": {
			input:    "aa:bb:cc:dd:ee:ff",
			expected: "AABBCCDDEEFF",
		},
		"dash format": {
			input:    "aa-bb-cc-dd-ee-ff",
			expected: "AABBCCDDEEFF",
		},
		"compact format": {
			input:    "aabbccddeeff",
			expected: "AABBCCDDEEFF",
		},
		"mixed separators": {
			input:       "aa:bb-cc:dd:ee:ff",
			shouldError: true,
		},
		"invalid length": {
			input:       "aabbccddee",
			shouldError: true,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := normalizeMACAddress(tc.input)
			if tc.shouldError {
				if err == nil {
					t.Fatalf("expected error for case %q", name)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error for case %q: %v", name, err)
			}

			if result != tc.expected {
				t.Fatalf("expected %q for case %q, got %q", tc.expected, name, result)
			}
		})
	}
}
