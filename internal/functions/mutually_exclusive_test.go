package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestMutuallyExclusiveMetadata(t *testing.T) {
	t.Parallel()

	fn := NewMutuallyExclusiveFunction()
	resp := &function.MetadataResponse{}
	fn.Metadata(context.Background(), function.MetadataRequest{}, resp)

	if resp.Name != "mutually_exclusive" {
		t.Fatalf("expected name mutually_exclusive, got %s", resp.Name)
	}
}

func TestMutuallyExclusiveDefinition(t *testing.T) {
	t.Parallel()

	fn := NewMutuallyExclusiveFunction()
	resp := &function.DefinitionResponse{}
	fn.Definition(context.Background(), function.DefinitionRequest{}, resp)

	if resp.Definition.Summary == "" {
		t.Fatal("expected non-empty Summary")
	}
	if len(resp.Definition.Parameters) != 1 {
		t.Fatalf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}
}

func TestMutuallyExclusiveFunction(t *testing.T) {
	t.Parallel()

	fn := NewMutuallyExclusiveFunction()
	ctx := context.Background()

	testCases := []struct {
		name          string
		values        attr.Value
		expectError   bool
		expectUnknown bool
		expectTrue    bool
	}{
		// Valid scenarios - exactly one value set
		{
			name:       "one non-empty value",
			values:     listValue([]string{"value1", "", ""}),
			expectTrue: true,
		},
		{
			name:       "one value at start",
			values:     listValue([]string{"value", ""}),
			expectTrue: true,
		},
		{
			name:       "one value at end",
			values:     listValue([]string{"", "value"}),
			expectTrue: true,
		},
		{
			name:       "one value among many empties",
			values:     listValue([]string{"", "", "only", "", ""}),
			expectTrue: true,
		},

		// Invalid scenarios - zero values set
		{
			name:        "all empty strings",
			values:      listValue([]string{"", "", ""}),
			expectError: true,
		},
		{
			name:        "empty list",
			values:      listValue([]string{}),
			expectError: true,
		},

		// Invalid scenarios - multiple values set
		{
			name:        "two values set",
			values:      listValue([]string{"value1", "value2"}),
			expectError: true,
		},
		{
			name:        "multiple values with empties",
			values:      listValue([]string{"value1", "", "value2"}),
			expectError: true,
		},
		{
			name:        "all values set",
			values:      listValue([]string{"a", "b", "c"}),
			expectError: true,
		},

		// Null and unknown
		{
			name:        "null list",
			values:      types.ListNull(basetypes.StringType{}),
			expectError: true,
		},
		{
			name:          "unknown list",
			values:        types.ListUnknown(basetypes.StringType{}),
			expectUnknown: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resp := &function.RunResponse{}
			fn.Run(ctx, function.RunRequest{Arguments: function.NewArgumentsData([]attr.Value{tc.values})}, resp)

			if tc.expectError {
				if resp.Error == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if resp.Error != nil {
				t.Fatalf("unexpected error: %s", resp.Error)
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
