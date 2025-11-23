package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestSlugValidator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		value     string
		wantError bool
	}{
		// Valid slugs
		{"simple slug", "hello", false},
		{"slug with hyphen", "hello-world", false},
		{"slug with digits", "hello-123", false},
		{"slug starting with digit", "123-hello", false},
		{"slug all digits", "123", false},
		{"slug multiple hyphens", "my-awesome-slug", false},
		{"slug with version", "my-app-v2", false},
		{"slug resource name", "web-server-01", false},

		// Invalid slugs
		{"empty string", "", true},
		{"uppercase", "Hello", true},
		{"leading hyphen", "-hello", true},
		{"trailing hyphen", "hello-", true},
		{"consecutive hyphens", "hello--world", true},
		{"underscore", "hello_world", true},
		{"space", "hello world", true},
		{"special chars", "hello@world", true},
		{"dots", "hello.world", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := Slug()
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

func TestSlugValidator_NullAndUnknown(t *testing.T) {
	t.Parallel()

	v := Slug()

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

func TestSlugValidator_Description(t *testing.T) {
	t.Parallel()

	v := Slug()

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
