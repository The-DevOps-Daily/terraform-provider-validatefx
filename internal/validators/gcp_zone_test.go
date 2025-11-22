package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestGCPZoneValidator(t *testing.T) {
	tests := []struct {
		name        string
		value       types.String
		expectError bool
	}{
		// Valid zones - US regions
		{"valid us-central1-a", types.StringValue("us-central1-a"), false},
		{"valid us-central1-c", types.StringValue("us-central1-c"), false},
		{"valid us-east1-b", types.StringValue("us-east1-b"), false},
		{"valid us-east4-a", types.StringValue("us-east4-a"), false},
		{"valid us-west1-a", types.StringValue("us-west1-a"), false},
		{"valid us-west2-b", types.StringValue("us-west2-b"), false},
		// Valid zones - North America
		{"valid northamerica-northeast1-a", types.StringValue("northamerica-northeast1-a"), false},
		// Valid zones - South America
		{"valid southamerica-east1-a", types.StringValue("southamerica-east1-a"), false},
		// Valid zones - Europe
		{"valid europe-west1-b", types.StringValue("europe-west1-b"), false},
		{"valid europe-west4-a", types.StringValue("europe-west4-a"), false},
		{"valid europe-north1-a", types.StringValue("europe-north1-a"), false},
		// Valid zones - Asia
		{"valid asia-east1-a", types.StringValue("asia-east1-a"), false},
		{"valid asia-northeast1-a", types.StringValue("asia-northeast1-a"), false},
		{"valid asia-south1-a", types.StringValue("asia-south1-a"), false},
		{"valid asia-southeast1-a", types.StringValue("asia-southeast1-a"), false},
		// Valid zones - Australia
		{"valid australia-southeast1-a", types.StringValue("australia-southeast1-a"), false},
		// Valid zones - Middle East
		{"valid me-central1-a", types.StringValue("me-central1-a"), false},
		{"valid me-west1-a", types.StringValue("me-west1-a"), false},
		// Valid zones - Africa
		{"valid africa-south1-a", types.StringValue("africa-south1-a"), false},
		// Invalid zones
		{"invalid empty", types.StringValue(""), false},
		{"invalid zone", types.StringValue("invalid-zone-x"), true},
		{"invalid region only", types.StringValue("us-central1"), true},
		{"invalid zone suffix", types.StringValue("us-central1-z"), true},
		{"invalid format", types.StringValue("not-a-zone"), true},
		{"invalid aws zone", types.StringValue("us-east-1a"), true},
		// Null and unknown
		{"null value", types.StringNull(), false},
		{"unknown value", types.StringUnknown(), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := GCPZone()
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

func TestGCPZoneValidator_Description(t *testing.T) {
	v := GCPZone()
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
