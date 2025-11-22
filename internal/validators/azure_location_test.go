package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestAzureLocationValidator_Valid(t *testing.T) {
	t.Parallel()
	v := AzureLocation()

	validLocations := []string{
		// Americas
		"eastus",
		"eastus2",
		"westus",
		"westus2",
		"centralus",
		"canadacentral",
		"brazilsouth",
		// Europe
		"northeurope",
		"westeurope",
		"uksouth",
		"francecentral",
		"germanywestcentral",
		"swedencentral",
		// Asia Pacific
		"eastasia",
		"southeastasia",
		"australiaeast",
		"japaneast",
		"koreacentral",
		"centralindia",
		// Middle East
		"uaenorth",
		"qatarcentral",
		// Africa
		"southafricanorth",
		// Special regions
		"chinaeast",
		"usgovvirginia",
	}

	for _, location := range validLocations {
		req := frameworkvalidator.StringRequest{
			Path:        path.Root("location"),
			ConfigValue: types.StringValue(location),
		}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)

		if resp.Diagnostics.HasError() {
			t.Errorf("expected no error for valid location %q, got: %v", location, resp.Diagnostics)
		}
	}
}

func TestAzureLocationValidator_Invalid(t *testing.T) {
	t.Parallel()
	v := AzureLocation()

	invalidLocations := []string{
		"invalid-location",
		"us-east-1",
		"east-us",
		"EastUS",
		"not-a-location",
		"east_us",
		"useast1",
		"eastus3",
	}

	for _, location := range invalidLocations {
		req := frameworkvalidator.StringRequest{
			Path:        path.Root("location"),
			ConfigValue: types.StringValue(location),
		}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)

		if !resp.Diagnostics.HasError() {
			t.Errorf("expected error for invalid location %q", location)
		}
	}
}

func TestAzureLocationValidator_NullUnknownEmpty(t *testing.T) {
	t.Parallel()
	v := AzureLocation()

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
			Path:        path.Root("location"),
			ConfigValue: tc.value,
		}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)

		if resp.Diagnostics.HasError() {
			t.Errorf("expected no error for %s value, got: %v", tc.name, resp.Diagnostics)
		}
	}
}

func TestAzureLocationValidator_Description(t *testing.T) {
	v := AzureLocation()
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
