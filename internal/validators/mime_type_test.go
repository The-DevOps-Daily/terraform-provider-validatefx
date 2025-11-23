package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestMIMETypeValidator(t *testing.T) {
	tests := []struct {
		name        string
		value       types.String
		expectError bool
	}{
		// Valid MIME types
		{"valid application/json", types.StringValue("application/json"), false},
		{"valid text/html", types.StringValue("text/html"), false},
		{"valid text/plain", types.StringValue("text/plain"), false},
		{"valid image/png", types.StringValue("image/png"), false},
		{"valid image/svg+xml", types.StringValue("image/svg+xml"), false},
		{"valid application/xml", types.StringValue("application/xml"), false},
		{"valid video/mp4", types.StringValue("video/mp4"), false},
		{"valid audio/mpeg", types.StringValue("audio/mpeg"), false},
		{"valid text/html; charset=utf-8", types.StringValue("text/html; charset=utf-8"), false},
		{"valid application/json; charset=utf-8", types.StringValue("application/json; charset=utf-8"), false},
		{"valid application/vnd.api+json", types.StringValue("application/vnd.api+json"), false},
		{"valid application/x-www-form-urlencoded", types.StringValue("application/x-www-form-urlencoded"), false},
		// Invalid MIME types
		{"invalid missing slash", types.StringValue("applicationjson"), true},
		{"invalid no subtype", types.StringValue("application/"), true},
		{"invalid no type", types.StringValue("/json"), true},
		{"invalid just text", types.StringValue("notamimetype"), true},
		{"invalid with spaces", types.StringValue("application / json"), true},
		{"invalid multiple slashes", types.StringValue("application/json/extra"), true},
		// Empty, null, unknown
		{"empty string", types.StringValue(""), false},
		{"null value", types.StringNull(), false},
		{"unknown value", types.StringUnknown(), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := MIMEType()
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

func TestMIMETypeValidator_Description(t *testing.T) {
	v := MIMEType()
	ctx := context.Background()

	desc := v.Description(ctx)
	if desc == "" {
		t.Error("Description should not be empty")
	}

	markdown := v.MarkdownDescription(ctx)
	if markdown == "" {
		t.Error("MarkdownDescription should not be empty")
	}

	if desc != markdown {
		t.Error("Description and MarkdownDescription should match")
	}
}
