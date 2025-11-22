package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestAzureLocationFunction(t *testing.T) {
	t.Parallel()

	fn := NewAzureLocationFunction()
	ctx := context.Background()

	cases := []struct {
		name          string
		value         attr.Value
		expectError   bool
		expectUnknown bool
		expectTrue    bool
	}{
		{
			name:       "valid eastus",
			value:      types.StringValue("eastus"),
			expectTrue: true,
		},
		{
			name:       "valid westeurope",
			value:      types.StringValue("westeurope"),
			expectTrue: true,
		},
		{
			name:       "valid southeastasia",
			value:      types.StringValue("southeastasia"),
			expectTrue: true,
		},
		{
			name:       "valid australiaeast",
			value:      types.StringValue("australiaeast"),
			expectTrue: true,
		},
		{
			name:       "valid usgovvirginia",
			value:      types.StringValue("usgovvirginia"),
			expectTrue: true,
		},
		{
			name:        "invalid location",
			value:       types.StringValue("invalid-location"),
			expectError: true,
		},
		{
			name:        "aws-style region",
			value:       types.StringValue("us-east-1"),
			expectError: true,
		},
		{
			name:          "null",
			value:         types.StringNull(),
			expectUnknown: true,
		},
		{
			name:          "unknown",
			value:         types.StringUnknown(),
			expectUnknown: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resp := &function.RunResponse{}
			fn.Run(ctx, function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{tc.value}),
			}, resp)

			if tc.expectError {
				if resp.Error == nil {
					t.Fatalf("expected error for %q", tc.name)
				}
				return
			}

			if resp.Error != nil {
				t.Fatalf("unexpected error: %v", resp.Error)
			}

			boolVal, ok := resp.Result.Value().(basetypes.BoolValue)
			if !ok {
				t.Fatalf("unexpected result type %T", resp.Result.Value())
			}

			if tc.expectUnknown {
				if !boolVal.IsUnknown() {
					t.Fatalf("expected unknown result")
				}
				return
			}

			if boolVal.IsUnknown() {
				t.Fatalf("did not expect unknown result")
			}

			if boolVal.ValueBool() != tc.expectTrue {
				t.Fatalf("expected %t, got %t", tc.expectTrue, boolVal.ValueBool())
			}
		})
	}
}

func TestAzureLocationFunction_Metadata(t *testing.T) {
	fn := NewAzureLocationFunction()
	resp := &function.MetadataResponse{}
	fn.Metadata(context.Background(), function.MetadataRequest{}, resp)

	if resp.Name != "azure_location" {
		t.Errorf("expected name 'azure_location', got %q", resp.Name)
	}
}

func TestAzureLocationFunction_Definition(t *testing.T) {
	fn := NewAzureLocationFunction()
	resp := &function.DefinitionResponse{}
	fn.Definition(context.Background(), function.DefinitionRequest{}, resp)

	if resp.Definition.Return == nil {
		t.Fatal("expected return definition, got nil")
	}
	if len(resp.Definition.Parameters) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}
}

func TestAzureLocationFunction_AllLocations(t *testing.T) {
	t.Parallel()

	fn := NewAzureLocationFunction()
	ctx := context.Background()

	// Test all major location types
	locations := []string{
		"eastus", "westus2", "centralus",
		"canadacentral",
		"brazilsouth",
		"northeurope", "westeurope", "uksouth", "francecentral",
		"germanywestcentral", "swedencentral",
		"eastasia", "southeastasia", "japaneast", "koreacentral",
		"australiaeast", "centralindia",
		"uaenorth", "qatarcentral",
		"southafricanorth",
		"chinaeast", "usgovvirginia",
	}

	for _, location := range locations {
		location := location
		t.Run(location, func(t *testing.T) {
			t.Parallel()

			resp := &function.RunResponse{}
			fn.Run(ctx, function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{types.StringValue(location)}),
			}, resp)

			if resp.Error != nil {
				t.Fatalf("unexpected error for location %q: %v", location, resp.Error)
			}

			boolVal, ok := resp.Result.Value().(basetypes.BoolValue)
			if !ok || !boolVal.ValueBool() {
				t.Fatalf("expected true for valid location %q", location)
			}
		})
	}
}
