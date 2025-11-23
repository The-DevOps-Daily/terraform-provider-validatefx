package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestUUIDv4OnlyFunction(t *testing.T) {
	t.Parallel()

	fn := NewUUIDv4OnlyFunction()
	ctx := context.Background()

	cases := []struct {
		name          string
		value         attr.Value
		expectError   bool
		expectUnknown bool
		expectTrue    bool
	}{
		{
			name:       "valid v4 lowercase",
			value:      types.StringValue("550e8400-e29b-41d4-a716-446655440000"),
			expectTrue: true,
		},
		{
			name:       "valid v4 uppercase",
			value:      types.StringValue("550E8400-E29B-41D4-A716-446655440000"),
			expectTrue: true,
		},
		{
			name:       "valid v4 mixed case",
			value:      types.StringValue("f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			expectTrue: true,
		},
		{
			name:       "valid v4 another",
			value:      types.StringValue("123e4567-e89b-42d3-a456-426614174000"),
			expectTrue: true,
		},
		{
			name:        "invalid v1",
			value:       types.StringValue("6ba7b810-9dad-11d1-80b4-00c04fd430c8"),
			expectError: true,
		},
		{
			name:        "invalid v3",
			value:       types.StringValue("6ba7b810-9dad-31d1-80b4-00c04fd430c8"),
			expectError: true,
		},
		{
			name:        "invalid v5",
			value:       types.StringValue("6ba7b810-9dad-51d1-80b4-00c04fd430c8"),
			expectError: true,
		},
		{
			name:        "invalid format",
			value:       types.StringValue("not-a-uuid"),
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

func TestUUIDv4OnlyFunction_Metadata(t *testing.T) {
	fn := NewUUIDv4OnlyFunction()
	resp := &function.MetadataResponse{}
	fn.Metadata(context.Background(), function.MetadataRequest{}, resp)

	if resp.Name != "uuidv4_only" {
		t.Errorf("expected name 'uuidv4_only', got %q", resp.Name)
	}
}

func TestUUIDv4OnlyFunction_Definition(t *testing.T) {
	fn := NewUUIDv4OnlyFunction()
	resp := &function.DefinitionResponse{}
	fn.Definition(context.Background(), function.DefinitionRequest{}, resp)

	if resp.Definition.Return == nil {
		t.Fatal("expected return definition, got nil")
	}
	if len(resp.Definition.Parameters) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}
}
