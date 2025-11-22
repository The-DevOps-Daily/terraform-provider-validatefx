package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestGCPRegionValidator_Valid(t *testing.T) {
	t.Parallel()
	v := GCPRegion()

	validRegions := []string{
		// Americas
		"us-central1",
		"us-east1",
		"us-east4",
		"us-west1",
		"us-west2",
		"northamerica-northeast1",
		"southamerica-east1",
		// Europe
		"europe-central2",
		"europe-north1",
		"europe-west1",
		"europe-west2",
		"europe-west3",
		"europe-west4",
		// Asia Pacific
		"asia-east1",
		"asia-east2",
		"asia-northeast1",
		"asia-northeast2",
		"asia-south1",
		"asia-southeast1",
		// Australia
		"australia-southeast1",
		"australia-southeast2",
		// Middle East
		"me-central1",
		// Africa
		"africa-south1",
	}

	for _, region := range validRegions {
		req := frameworkvalidator.StringRequest{
			Path:        path.Root("region"),
			ConfigValue: types.StringValue(region),
		}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)

		if resp.Diagnostics.HasError() {
			t.Errorf("expected no error for valid region %q, got: %v", region, resp.Diagnostics)
		}
	}
}

func TestGCPRegionValidator_Invalid(t *testing.T) {
	t.Parallel()
	v := GCPRegion()

	invalidRegions := []string{
		"invalid-region",
		"us-east-1",
		"us-central2",
		"eu-west-1",
		"asia-east-1",
		"not-a-region",
		"us_central1",
		"US-CENTRAL1",
		"us-central",
	}

	for _, region := range invalidRegions {
		req := frameworkvalidator.StringRequest{
			Path:        path.Root("region"),
			ConfigValue: types.StringValue(region),
		}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)

		if !resp.Diagnostics.HasError() {
			t.Errorf("expected error for invalid region %q", region)
		}
	}
}

func TestGCPRegionValidator_NullUnknownEmpty(t *testing.T) {
	t.Parallel()
	v := GCPRegion()

	cases := []struct {
		name  string
		value types.String
	}{
		{"null", types.StringNull()},
		{"unknown", types.StringUnknown()},
		{"empty", types.StringValue("")},
	}

	for _, tc := range cases {
		req := frameworkvalidator.StringRequest{
			Path:        path.Root("region"),
			ConfigValue: tc.value,
		}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)

		if resp.Diagnostics.HasError() {
			t.Errorf("expected no error for %s value, got: %v", tc.name, resp.Diagnostics)
		}
	}
}

func TestGCPRegionValidator_Description(t *testing.T) {
	v := GCPRegion()
	ctx := context.Background()

	desc := v.Description(ctx)
	if desc == "" {
		t.Error("expected non-empty description")
	}

	markdownDesc := v.MarkdownDescription(ctx)
	if markdownDesc == "" {
		t.Error("expected non-empty markdown description")
	}

	if desc != markdownDesc {
		t.Errorf("expected description %q to match markdown description %q", desc, markdownDesc)
	}
}
