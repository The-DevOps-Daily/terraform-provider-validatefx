package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestMapKeysMatchFunction(t *testing.T) {
	t.Parallel()

	fn := NewMapKeysMatchFunction()
	ctx := context.Background()

	testCases := []struct {
		name          string
		inputMap      attr.Value
		allowedKeys   attr.Value
		requiredKeys  attr.Value
		expectError   bool
		expectUnknown bool
		expectTrue    bool
	}{
		// Valid scenarios
		{
			name:         "all keys allowed",
			inputMap:     mapValue(map[string]string{"a": "1", "b": "2"}),
			allowedKeys:  listValue([]string{"a", "b", "c"}),
			requiredKeys: listValue([]string{}),
			expectTrue:   true,
		},
		{
			name:         "required keys present",
			inputMap:     mapValue(map[string]string{"a": "1", "b": "2"}),
			allowedKeys:  listValue([]string{"a", "b", "c"}),
			requiredKeys: listValue([]string{"a"}),
			expectTrue:   true,
		},
		{
			name:         "empty allowed means all allowed",
			inputMap:     mapValue(map[string]string{"a": "1", "b": "2"}),
			allowedKeys:  listValue([]string{}),
			requiredKeys: listValue([]string{"a"}),
			expectTrue:   true,
		},

		// Invalid scenarios
		{
			name:         "disallowed key",
			inputMap:     mapValue(map[string]string{"a": "1", "c": "3"}),
			allowedKeys:  listValue([]string{"a", "b"}),
			requiredKeys: listValue([]string{}),
			expectError:  true,
		},
		{
			name:         "missing required key",
			inputMap:     mapValue(map[string]string{"a": "1"}),
			allowedKeys:  listValue([]string{"a", "b", "c"}),
			requiredKeys: listValue([]string{"a", "b"}),
			expectError:  true,
		},

		// Null and unknown
		{
			name:          "null map",
			inputMap:      types.MapNull(basetypes.StringType{}),
			allowedKeys:   listValue([]string{"a"}),
			requiredKeys:  listValue([]string{}),
			expectError:   true,
		},
		{
			name:          "unknown map",
			inputMap:      types.MapUnknown(basetypes.StringType{}),
			allowedKeys:   listValue([]string{"a"}),
			requiredKeys:  listValue([]string{}),
			expectUnknown: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resp := &function.RunResponse{}
			fn.Run(ctx, function.RunRequest{Arguments: function.NewArgumentsData([]attr.Value{tc.inputMap, tc.allowedKeys, tc.requiredKeys})}, resp)

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

func mapValue(m map[string]string) types.Map {
	elements := make(map[string]attr.Value, len(m))
	for k, v := range m {
		elements[k] = types.StringValue(v)
	}
	val, _ := types.MapValue(basetypes.StringType{}, elements)
	return val
}

func listValue(items []string) types.List {
	elements := make([]attr.Value, len(items))
	for i, item := range items {
		elements[i] = types.StringValue(item)
	}
	val, _ := types.ListValue(basetypes.StringType{}, elements)
	return val
}
