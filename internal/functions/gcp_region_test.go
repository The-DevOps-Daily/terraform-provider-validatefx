package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestGCPRegionFunction(t *testing.T) {
	t.Parallel()

	fn := NewGCPRegionFunction()
	ctx := context.Background()

	cases := []struct {
		name          string
		value         attr.Value
		expectError   bool
		expectUnknown bool
		expectTrue    bool
	}{
		{
			name:       "valid us-central1",
			value:      types.StringValue("us-central1"),
			expectTrue: true,
		},
		{
			name:       "valid europe-west1",
			value:      types.StringValue("europe-west1"),
			expectTrue: true,
		},
		{
			name:       "valid asia-southeast1",
			value:      types.StringValue("asia-southeast1"),
			expectTrue: true,
		},
		{
			name:       "valid australia-southeast1",
			value:      types.StringValue("australia-southeast1"),
			expectTrue: true,
		},
		{
			name:       "valid africa-south1",
			value:      types.StringValue("africa-south1"),
			expectTrue: true,
		},
		{
			name:        "invalid region",
			value:       types.StringValue("invalid-region"),
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

func TestGCPRegionFunction_Metadata(t *testing.T) {
	fn := NewGCPRegionFunction()
	resp := &function.MetadataResponse{}
	fn.Metadata(context.Background(), function.MetadataRequest{}, resp)

	if resp.Name != "gcp_region" {
		t.Errorf("expected name 'gcp_region', got %q", resp.Name)
	}
}

func TestGCPRegionFunction_Definition(t *testing.T) {
	fn := NewGCPRegionFunction()
	resp := &function.DefinitionResponse{}
	fn.Definition(context.Background(), function.DefinitionRequest{}, resp)

	if resp.Definition.Return == nil {
		t.Fatal("expected return definition, got nil")
	}
	if len(resp.Definition.Parameters) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}
}

func TestGCPRegionFunction_AllRegions(t *testing.T) {
	t.Parallel()

	fn := NewGCPRegionFunction()
	ctx := context.Background()

	// Test all major region types
	regions := []string{
		"us-central1", "us-east1", "us-west1",
		"northamerica-northeast1",
		"southamerica-east1",
		"europe-west1", "europe-north1", "europe-central2",
		"asia-east1", "asia-northeast1", "asia-southeast1", "asia-south1",
		"australia-southeast1",
		"me-central1",
		"africa-south1",
	}

	for _, region := range regions {
		region := region
		t.Run(region, func(t *testing.T) {
			t.Parallel()

			resp := &function.RunResponse{}
			fn.Run(ctx, function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{types.StringValue(region)}),
			}, resp)

			if resp.Error != nil {
				t.Fatalf("unexpected error for region %q: %v", region, resp.Error)
			}

			boolVal, ok := resp.Result.Value().(basetypes.BoolValue)
			if !ok || !boolVal.ValueBool() {
				t.Fatalf("expected true for valid region %q", region)
			}
		})
	}
}
