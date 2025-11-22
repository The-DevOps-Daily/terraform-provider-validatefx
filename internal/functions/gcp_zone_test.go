package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestGCPZoneFunction(t *testing.T) {
	t.Parallel()

	fn := NewGCPZoneFunction()
	ctx := context.Background()

	cases := []struct {
		name          string
		value         attr.Value
		expectError   bool
		expectUnknown bool
		expectTrue    bool
	}{
		{
			name:       "valid us-central1-a",
			value:      types.StringValue("us-central1-a"),
			expectTrue: true,
		},
		{
			name:       "valid europe-west1-b",
			value:      types.StringValue("europe-west1-b"),
			expectTrue: true,
		},
		{
			name:       "valid asia-southeast1-a",
			value:      types.StringValue("asia-southeast1-a"),
			expectTrue: true,
		},
		{
			name:       "valid australia-southeast1-a",
			value:      types.StringValue("australia-southeast1-a"),
			expectTrue: true,
		},
		{
			name:       "valid africa-south1-a",
			value:      types.StringValue("africa-south1-a"),
			expectTrue: true,
		},
		{
			name:        "invalid zone",
			value:       types.StringValue("invalid-zone-x"),
			expectError: true,
		},
		{
			name:        "region only",
			value:       types.StringValue("us-central1"),
			expectError: true,
		},
		{
			name:        "aws-style zone",
			value:       types.StringValue("us-east-1a"),
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

func TestGCPZoneFunction_Metadata(t *testing.T) {
	fn := NewGCPZoneFunction()
	resp := &function.MetadataResponse{}
	fn.Metadata(context.Background(), function.MetadataRequest{}, resp)

	if resp.Name != "gcp_zone" {
		t.Errorf("expected name 'gcp_zone', got %q", resp.Name)
	}
}

func TestGCPZoneFunction_Definition(t *testing.T) {
	fn := NewGCPZoneFunction()
	resp := &function.DefinitionResponse{}
	fn.Definition(context.Background(), function.DefinitionRequest{}, resp)

	if resp.Definition.Return == nil {
		t.Fatal("expected return definition, got nil")
	}
	if len(resp.Definition.Parameters) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}
}

func TestGCPZoneFunction_AllZoneTypes(t *testing.T) {
	t.Parallel()

	fn := NewGCPZoneFunction()
	ctx := context.Background()

	// Test representative zones from all major regions
	zones := []string{
		"us-central1-a", "us-east1-b", "us-west1-a",
		"northamerica-northeast1-a",
		"southamerica-east1-a",
		"europe-west1-b", "europe-north1-a", "europe-central2-a",
		"asia-east1-a", "asia-northeast1-a", "asia-southeast1-a", "asia-south1-a",
		"australia-southeast1-a",
		"me-central1-a",
		"africa-south1-a",
	}

	for _, zone := range zones {
		zone := zone
		t.Run(zone, func(t *testing.T) {
			t.Parallel()

			resp := &function.RunResponse{}
			fn.Run(ctx, function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{types.StringValue(zone)}),
			}, resp)

			if resp.Error != nil {
				t.Fatalf("unexpected error for zone %q: %v", zone, resp.Error)
			}

			boolVal, ok := resp.Result.Value().(basetypes.BoolValue)
			if !ok || !boolVal.ValueBool() {
				t.Fatalf("expected true for valid zone %q", zone)
			}
		})
	}
}
