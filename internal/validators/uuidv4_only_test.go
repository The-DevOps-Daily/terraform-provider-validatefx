package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestUUIDv4OnlyValidator(t *testing.T) {
	tests := []struct {
		name        string
		value       types.String
		expectError bool
	}{
		// Valid UUIDv4
		{"valid v4 lowercase", types.StringValue("550e8400-e29b-41d4-a716-446655440000"), false},
		{"valid v4 uppercase", types.StringValue("550E8400-E29B-41D4-A716-446655440000"), false},
		{"valid v4 mixed case", types.StringValue("f47ac10b-58cc-4372-a567-0e02b2c3d479"), false},
		{"valid v4 another", types.StringValue("123e4567-e89b-42d3-a456-426614174000"), false},
		{"valid v4 zeros", types.StringValue("00000000-0000-4000-8000-000000000000"), false},
		// Invalid - wrong UUID version
		{"invalid v1", types.StringValue("6ba7b810-9dad-11d1-80b4-00c04fd430c8"), true},
		{"invalid v3", types.StringValue("6ba7b810-9dad-31d1-80b4-00c04fd430c8"), true},
		{"invalid v5", types.StringValue("6ba7b810-9dad-51d1-80b4-00c04fd430c8"), true},
		// Invalid format
		{"invalid format", types.StringValue("not-a-uuid"), true},
		{"invalid short", types.StringValue("550e8400-e29b-41d4"), true},
		{"invalid no dashes", types.StringValue("550e8400e29b41d4a716446655440000"), false}, // uuid.Parse accepts this
		// Empty, null, unknown
		{"empty string", types.StringValue(""), false},
		{"null value", types.StringNull(), false},
		{"unknown value", types.StringUnknown(), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := UUIDv4Only()
			req := validator.StringRequest{
				Path:        path.Root("test"),
				ConfigValue: tt.value,
			}
			resp := &validator.StringResponse{}

			v.ValidateString(context.Background(), req, resp)

			if tt.expectError && !resp.Diagnostics.HasError() {
				t.Errorf("expected error for value %v but got none", tt.value)
			}
			if !tt.expectError && resp.Diagnostics.HasError() {
				t.Errorf("unexpected error for value %v: %v", tt.value, resp.Diagnostics.Errors())
			}
		})
	}
}

func TestUUIDv4OnlyValidator_Description(t *testing.T) {
	v := UUIDv4Only()
	ctx := context.Background()

	desc := v.Description(ctx)
	if desc == "" {
		t.Error("Description should not be empty")
	}

	if desc != "value must be a valid UUID version 4" {
		t.Errorf("unexpected description: %s", desc)
	}

	markdown := v.MarkdownDescription(ctx)
	if markdown == "" {
		t.Error("MarkdownDescription should not be empty")
	}

	if desc != markdown {
		t.Error("Description and MarkdownDescription should match")
	}
}
