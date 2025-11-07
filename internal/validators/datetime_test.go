package validators

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestDateTimeValidator_DefaultLayout(t *testing.T) {
	t.Parallel()

	validator := DateTime(nil)

	tests := map[string]struct {
		value       types.String
		expectError bool
	}{
		"valid RFC3339": {
			value:       types.StringValue("2025-11-02T15:04:05Z"),
			expectError: false,
		},
		"valid RFC3339Nano": {
			value:       types.StringValue("2025-11-02T15:04:05.123456789Z"),
			expectError: false,
		},
		"invalid date": {
			value:       types.StringValue("2025-13-02T15:04:05Z"),
			expectError: true,
		},
		"invalid time": {
			value:       types.StringValue("2025-11-02T25:04:05Z"),
			expectError: true,
		},
		"empty": {
			value:       types.StringValue(""),
			expectError: false,
		},
	}

	for name, tc := range tests {
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
				t.Fatalf("expected error=%t got=%t diagnostics=%v", tc.expectError, resp.Diagnostics.HasError(), resp.Diagnostics)
			}
		})
	}
}

func TestDateTimeValidator_CustomLayouts(t *testing.T) {
	t.Parallel()

	validator := DateTime([]string{
		"2006-01-02 15:04:05",
	})

	req := frameworkvalidator.StringRequest{
		Path:        path.Root("value"),
		ConfigValue: types.StringValue("2025-11-02 21:30:00"),
	}
	resp := &frameworkvalidator.StringResponse{}

	validator.ValidateString(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected custom layout to pass: %v", resp.Diagnostics)
	}
}

func TestDateTimeValidator_InvalidCustomLayout(t *testing.T) {
	t.Parallel()

	validator := DateTime([]string{
		"2006-01-02 15:04:05",
	})

	req := frameworkvalidator.StringRequest{
		Path:        path.Root("value"),
		ConfigValue: types.StringValue("02-11-2025"),
	}
	resp := &frameworkvalidator.StringResponse{}

	validator.ValidateString(context.Background(), req, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatalf("expected custom layout to fail")
	}

	if resp.Diagnostics[0].Summary() != "Invalid Datetime" {
		t.Fatalf("unexpected diagnostic summary: %s", resp.Diagnostics[0].Summary())
	}
}

func TestDateTimeValidatorHandlesNullUnknown(t *testing.T) {
	t.Parallel()

	validator := DateTime(nil)

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

func TestDateTimeValidatorWithLocation(t *testing.T) {
	t.Parallel()

	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("load location: %v", err)
	}

	validator := DateTimeWithLocation([]string{"2006-01-02 15:04"}, location)

	req := frameworkvalidator.StringRequest{
		Path:        path.Root("value"),
		ConfigValue: types.StringValue("2025-11-02 07:30"),
	}
	resp := &frameworkvalidator.StringResponse{}

	validator.ValidateString(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected time to parse with location: %v", resp.Diagnostics)
	}
}
