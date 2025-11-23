package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestSizeBetweenValidator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		min       string
		max       string
		value     string
		wantError bool
	}{
		// Valid values within range
		{"value in range", "1", "10", "5", false},
		{"value at min", "1", "10", "1", false},
		{"value at max", "1", "10", "10", false},
		{"decimal in range", "0", "1", "0.5", false},
		{"negative range", "-10", "-1", "-5", false},
		{"large numbers", "100", "1000", "500", false},

		// Invalid values
		{"value below min", "1", "10", "0", true},
		{"value above max", "1", "10", "11", true},
		{"non-numeric value", "1", "10", "abc", true},
		{"empty string", "1", "10", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := SizeBetween(tt.min, tt.max)
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

func TestSizeBetweenValidator_NullAndUnknown(t *testing.T) {
	t.Parallel()

	v := SizeBetween("1", "10")

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

func TestSizeBetweenValidator_Description(t *testing.T) {
	t.Parallel()

	v := SizeBetween("1", "10")

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
