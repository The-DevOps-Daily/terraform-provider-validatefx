package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestListLengthBetweenMetadata(t *testing.T) {
	t.Parallel()

	fn := NewListLengthBetweenFunction()
	resp := &function.MetadataResponse{}
	fn.Metadata(context.Background(), function.MetadataRequest{}, resp)

	if resp.Name != "list_length_between" {
		t.Fatalf("expected name list_length_between, got %s", resp.Name)
	}
}

func TestListLengthBetweenDefinition(t *testing.T) {
	t.Parallel()

	fn := NewListLengthBetweenFunction()
	resp := &function.DefinitionResponse{}
	fn.Definition(context.Background(), function.DefinitionRequest{}, resp)

	if resp.Definition.Summary == "" {
		t.Fatal("expected non-empty Summary")
	}
	if len(resp.Definition.Parameters) != 3 {
		t.Fatalf("expected 3 parameters, got %d", len(resp.Definition.Parameters))
	}
}

func TestListLengthBetweenFunction(t *testing.T) {
	t.Parallel()

	fn := NewListLengthBetweenFunction()
	ctx := context.Background()

	testCases := []struct {
		name          string
		values        attr.Value
		min           attr.Value
		max           attr.Value
		expectError   bool
		expectUnknown bool
		expectTrue    bool
	}{
		// Valid lengths
		{
			name:       "length at minimum",
			values:     listValue([]string{"a", "b"}),
			min:        types.StringValue("2"),
			max:        types.StringValue("5"),
			expectTrue: true,
		},
		{
			name:       "length at maximum",
			values:     listValue([]string{"a", "b", "c", "d", "e"}),
			min:        types.StringValue("2"),
			max:        types.StringValue("5"),
			expectTrue: true,
		},
		{
			name:       "length in middle",
			values:     listValue([]string{"a", "b", "c"}),
			min:        types.StringValue("2"),
			max:        types.StringValue("5"),
			expectTrue: true,
		},
		{
			name:       "exact length when min equals max",
			values:     listValue([]string{"a", "b", "c"}),
			min:        types.StringValue("3"),
			max:        types.StringValue("3"),
			expectTrue: true,
		},
		{
			name:       "empty list with zero minimum",
			values:     listValue([]string{}),
			min:        types.StringValue("0"),
			max:        types.StringValue("5"),
			expectTrue: true,
		},

		// Invalid lengths - too short
		{
			name:        "length below minimum",
			values:      listValue([]string{"a"}),
			min:         types.StringValue("2"),
			max:         types.StringValue("5"),
			expectError: true,
		},
		{
			name:        "empty list with positive minimum",
			values:      listValue([]string{}),
			min:         types.StringValue("1"),
			max:         types.StringValue("5"),
			expectError: true,
		},

		// Invalid lengths - too long
		{
			name:        "length above maximum",
			values:      listValue([]string{"a", "b", "c", "d", "e", "f"}),
			min:         types.StringValue("2"),
			max:         types.StringValue("5"),
			expectError: true,
		},

		// Invalid parameters
		{
			name:        "negative minimum",
			values:      listValue([]string{"a"}),
			min:         types.StringValue("-1"),
			max:         types.StringValue("5"),
			expectError: true,
		},
		{
			name:        "min greater than max",
			values:      listValue([]string{"a"}),
			min:         types.StringValue("5"),
			max:         types.StringValue("2"),
			expectError: true,
		},

		// Null and unknown
		{
			name:          "null list",
			values:        types.ListNull(basetypes.StringType{}),
			min:           types.StringValue("2"),
			max:           types.StringValue("5"),
			expectUnknown: true,
		},
		{
			name:          "unknown list",
			values:        types.ListUnknown(basetypes.StringType{}),
			min:           types.StringValue("2"),
			max:           types.StringValue("5"),
			expectUnknown: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resp := &function.RunResponse{}
			fn.Run(ctx, function.RunRequest{Arguments: function.NewArgumentsData([]attr.Value{tc.values, tc.min, tc.max})}, resp)

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
