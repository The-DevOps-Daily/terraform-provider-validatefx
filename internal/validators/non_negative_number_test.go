package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestNonNegativeNumberValidator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		value     string
		wantError bool
	}{
		// Valid non-negative numbers
		{"zero", "0", false},
		{"positive integer", "1", false},
		{"positive decimal", "1.5", false},
		{"large positive", "9999.99", false},
		{"small positive", "0.001", false},
		{"zero decimal", "0.0", false},
		{"positive with plus sign", "+5", false},
		{"positive decimal with plus", "+3.14", false},

		// Invalid values
		{"negative integer", "-1", true},
		{"negative decimal", "-1.5", true},
		{"negative zero", "-0.1", true},
		{"non-numeric", "abc", true},
		{"empty string", "", true},
		{"mixed characters", "123abc", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := NonNegativeNumber()
			req := validator.StringRequest{
				Path:        path.Root("test"),
				ConfigValue: types.StringValue(tt.value),
			}
			resp := &validator.StringResponse{}

			v.ValidateString(context.Background(), req, resp)

			if tt.wantError && !resp.Diagnostics.HasError() {
				t.Errorf("expected error for %q, got none", tt.value)
			}
			if !tt.wantError && resp.Diagnostics.HasError() {
				t.Errorf("unexpected error for %q: %v", tt.value, resp.Diagnostics)
			}
		})
	}
}

func TestNonNegativeNumberValidator_NullAndUnknown(t *testing.T) {
	t.Parallel()

	v := NonNegativeNumber()

	// Null value
	req := validator.StringRequest{
		Path:        path.Root("test"),
		ConfigValue: types.StringNull(),
	}
	resp := &validator.StringResponse{}
	v.ValidateString(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Errorf("unexpected error for null value: %v", resp.Diagnostics)
	}

	// Unknown value
	req = validator.StringRequest{
		Path:        path.Root("test"),
		ConfigValue: types.StringUnknown(),
	}
	resp = &validator.StringResponse{}
	v.ValidateString(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Errorf("unexpected error for unknown value: %v", resp.Diagnostics)
	}
}

func TestNonNegativeNumberValidator_Description(t *testing.T) {
	t.Parallel()

	v := NonNegativeNumber()

	desc := v.Description(context.Background())
	if desc == "" {
		t.Error("Description should not be empty")
	}

	markdownDesc := v.MarkdownDescription(context.Background())
	if markdownDesc == "" {
		t.Error("MarkdownDescription should not be empty")
	}

	if desc != markdownDesc {
		t.Errorf("Description and MarkdownDescription should match: %q != %q", desc, markdownDesc)
	}
}
