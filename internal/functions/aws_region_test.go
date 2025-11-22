package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestAWSRegionFunction(t *testing.T) {
	t.Parallel()

	fn := NewAWSRegionFunction()
	ctx := context.Background()

	cases := []struct {
		name          string
		value         attr.Value
		expectError   bool
		expectUnknown bool
		expectTrue    bool
	}{
		{
			name:       "valid us-east-1",
			value:      types.StringValue("us-east-1"),
			expectTrue: true,
		},
		{
			name:       "valid eu-west-1",
			value:      types.StringValue("eu-west-1"),
			expectTrue: true,
		},
		{
			name:       "valid ap-southeast-1",
			value:      types.StringValue("ap-southeast-1"),
			expectTrue: true,
		},
		{
			name:       "valid us-gov-west-1",
			value:      types.StringValue("us-gov-west-1"),
			expectTrue: true,
		},
		{
			name:       "valid cn-north-1",
			value:      types.StringValue("cn-north-1"),
			expectTrue: true,
		},
		{
			name:        "invalid region",
			value:       types.StringValue("invalid-region"),
			expectError: true,
		},
		{
			name:        "non-existent us-east-3",
			value:       types.StringValue("us-east-3"),
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

func TestAWSRegionFunction_Metadata(t *testing.T) {
	fn := NewAWSRegionFunction()
	resp := &function.MetadataResponse{}
	fn.Metadata(context.Background(), function.MetadataRequest{}, resp)

	if resp.Name != "aws_region" {
		t.Errorf("expected name 'aws_region', got %q", resp.Name)
	}
}

func TestAWSRegionFunction_Definition(t *testing.T) {
	fn := NewAWSRegionFunction()
	resp := &function.DefinitionResponse{}
	fn.Definition(context.Background(), function.DefinitionRequest{}, resp)

	if resp.Definition.Return == nil {
		t.Fatal("expected return definition, got nil")
	}
	if len(resp.Definition.Parameters) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}
}

func TestAWSRegionFunction_AllRegions(t *testing.T) {
	t.Parallel()

	fn := NewAWSRegionFunction()
	ctx := context.Background()

	// Test all major region types
	regions := []string{
		"us-east-1", "us-east-2", "us-west-1", "us-west-2",
		"ca-central-1",
		"eu-central-1", "eu-west-1", "eu-north-1",
		"ap-northeast-1", "ap-southeast-1", "ap-south-1",
		"sa-east-1",
		"me-south-1",
		"af-south-1",
		"us-gov-west-1", "us-gov-east-1",
		"cn-north-1", "cn-northwest-1",
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
