package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestDependentValueMetadata(t *testing.T) {
	t.Parallel()

	fn := NewDependentValueFunction()
	resp := &function.MetadataResponse{}
	fn.Metadata(context.Background(), function.MetadataRequest{}, resp)

	if resp.Name != "dependent_value" {
		t.Fatalf("expected name dependent_value, got %s", resp.Name)
	}
}

func TestDependentValueDefinition(t *testing.T) {
	t.Parallel()

	fn := NewDependentValueFunction()
	resp := &function.DefinitionResponse{}
	fn.Definition(context.Background(), function.DefinitionRequest{}, resp)

	if resp.Definition.Summary == "" {
		t.Fatal("expected non-empty Summary")
	}
	if len(resp.Definition.Parameters) != 2 {
		t.Fatalf("expected 2 parameters, got %d", len(resp.Definition.Parameters))
	}
}

func TestDependentValueFunction(t *testing.T) {
	t.Parallel()

	fn := NewDependentValueFunction()
	ctx := context.Background()

	testCases := []struct {
		name          string
		condition     attr.Value
		dependent     attr.Value
		expectError   bool
		expectUnknown bool
		expectTrue    bool
	}{
		// Valid cases
		{
			name:       "both empty strings",
			condition:  types.StringValue(""),
			dependent:  types.StringValue(""),
			expectTrue: true,
		},
		{
			name:       "both set",
			condition:  types.StringValue("enabled"),
			dependent:  types.StringValue("config"),
			expectTrue: true,
		},
		{
			name:       "condition empty, dependent set",
			condition:  types.StringValue(""),
			dependent:  types.StringValue("value"),
			expectTrue: true,
		},
		{
			name:       "both null",
			condition:  types.StringNull(),
			dependent:  types.StringNull(),
			expectTrue: true,
		},
		{
			name:       "condition null, dependent set",
			condition:  types.StringNull(),
			dependent:  types.StringValue("value"),
			expectTrue: true,
		},

		// Invalid cases
		{
			name:        "condition set, dependent empty",
			condition:   types.StringValue("enabled"),
			dependent:   types.StringValue(""),
			expectError: true,
		},
		{
			name:        "condition set, dependent null",
			condition:   types.StringValue("enabled"),
			dependent:   types.StringNull(),
			expectError: true,
		},
		{
			name:        "condition true, dependent empty",
			condition:   types.StringValue("true"),
			dependent:   types.StringValue(""),
			expectError: true,
		},

		// Unknown values
		{
			name:          "condition unknown",
			condition:     types.StringUnknown(),
			dependent:     types.StringValue("value"),
			expectUnknown: true,
		},
		{
			name:          "dependent unknown",
			condition:     types.StringValue("enabled"),
			dependent:     types.StringUnknown(),
			expectUnknown: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resp := &function.RunResponse{}
			fn.Run(ctx, function.RunRequest{Arguments: function.NewArgumentsData([]attr.Value{tc.condition, tc.dependent})}, resp)

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
