package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestListUniqueMetadata(t *testing.T) {
	t.Parallel()

	fn := NewListUniqueFunction()
	resp := &function.MetadataResponse{}
	fn.Metadata(context.Background(), function.MetadataRequest{}, resp)

	if resp.Name != "list_unique" {
		t.Fatalf("expected name list_unique, got %s", resp.Name)
	}
}

func TestListUniqueDefinition(t *testing.T) {
	t.Parallel()

	fn := NewListUniqueFunction()
	resp := &function.DefinitionResponse{}
	fn.Definition(context.Background(), function.DefinitionRequest{}, resp)

	if resp.Definition.Summary == "" {
		t.Fatal("expected non-empty Summary")
	}
	if len(resp.Definition.Parameters) != 1 {
		t.Fatalf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}
}

func TestListUniqueFunction(t *testing.T) {
	t.Parallel()

	fn := NewListUniqueFunction()
	ctx := context.Background()

	testCases := []struct {
		name          string
		values        attr.Value
		expectError   bool
		expectUnknown bool
		expectTrue    bool
	}{
		// Valid - all unique
		{
			name:       "all unique elements",
			values:     listValue([]string{"a", "b", "c"}),
			expectTrue: true,
		},
		{
			name:       "single element",
			values:     listValue([]string{"a"}),
			expectTrue: true,
		},
		{
			name:       "empty list",
			values:     listValue([]string{}),
			expectTrue: true,
		},
		{
			name:       "many unique elements",
			values:     listValue([]string{"apple", "banana", "cherry", "date", "elderberry"}),
			expectTrue: true,
		},
		{
			name:       "numbers as strings",
			values:     listValue([]string{"1", "2", "3", "4", "5"}),
			expectTrue: true,
		},

		// Invalid - contains duplicates
		{
			name:        "simple duplicate",
			values:      listValue([]string{"a", "b", "a"}),
			expectError: true,
		},
		{
			name:        "multiple duplicates",
			values:      listValue([]string{"a", "b", "a", "b"}),
			expectError: true,
		},
		{
			name:        "duplicate at end",
			values:      listValue([]string{"a", "b", "c", "a"}),
			expectError: true,
		},
		{
			name:        "consecutive duplicates",
			values:      listValue([]string{"a", "a"}),
			expectError: true,
		},
		{
			name:        "all same elements",
			values:      listValue([]string{"a", "a", "a", "a"}),
			expectError: true,
		},

		// Null and unknown
		{
			name:          "null list",
			values:        types.ListNull(basetypes.StringType{}),
			expectUnknown: true,
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
