package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestAWSRegionValidator_Valid(t *testing.T) {
	t.Parallel()
	v := AWSRegion()

	validRegions := []string{
		// US regions
		"us-east-1",
		"us-east-2",
		"us-west-1",
		"us-west-2",
		// US GovCloud
		"us-gov-east-1",
		"us-gov-west-1",
		// Canada
		"ca-central-1",
		// Europe
		"eu-central-1",
		"eu-central-2",
		"eu-west-1",
		"eu-west-2",
		"eu-west-3",
		"eu-north-1",
		"eu-south-1",
		"eu-south-2",
		// Asia Pacific
		"ap-east-1",
		"ap-south-1",
		"ap-south-2",
		"ap-northeast-1",
		"ap-northeast-2",
		"ap-northeast-3",
		"ap-southeast-1",
		"ap-southeast-2",
		"ap-southeast-3",
		"ap-southeast-4",
		// South America
		"sa-east-1",
		// Middle East
		"me-central-1",
		"me-south-1",
		// Africa
		"af-south-1",
		// China
		"cn-north-1",
		"cn-northwest-1",
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

func TestAWSRegionValidator_Invalid(t *testing.T) {
	t.Parallel()
	v := AWSRegion()

	invalidRegions := []string{
		"invalid-region",
		"us-east-3",
		"eu-west-5",
		"ap-south-3",
		"not-a-region",
		"us_east_1",
		"US-EAST-1",
		"us-east",
		"east-1",
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

func TestAWSRegionValidator_NullUnknownEmpty(t *testing.T) {
	t.Parallel()
	v := AWSRegion()

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

func TestAWSRegionValidator_Description(t *testing.T) {
	v := AWSRegion()
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
