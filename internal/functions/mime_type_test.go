package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestMIMETypeFunction(t *testing.T) {
	t.Parallel()

	fn := NewMIMETypeFunction()
	ctx := context.Background()

	cases := []struct {
		name          string
		value         attr.Value
		expectError   bool
		expectUnknown bool
		expectTrue    bool
	}{
		{
			name:       "valid application/json",
			value:      types.StringValue("application/json"),
			expectTrue: true,
		},
		{
			name:       "valid text/html",
			value:      types.StringValue("text/html"),
			expectTrue: true,
		},
		{
			name:       "valid image/svg+xml",
			value:      types.StringValue("image/svg+xml"),
			expectTrue: true,
		},
		{
			name:       "valid with parameters",
			value:      types.StringValue("text/html; charset=utf-8"),
			expectTrue: true,
		},
		{
			name:        "invalid missing slash",
			value:       types.StringValue("applicationjson"),
			expectError: true,
		},
		{
			name:        "invalid no subtype",
			value:       types.StringValue("application/"),
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

func TestMIMETypeFunction_Metadata(t *testing.T) {
	fn := NewMIMETypeFunction()
	resp := &function.MetadataResponse{}
	fn.Metadata(context.Background(), function.MetadataRequest{}, resp)

	if resp.Name != "mime_type" {
		t.Errorf("expected name 'mime_type', got %q", resp.Name)
	}
}

func TestMIMETypeFunction_Definition(t *testing.T) {
	fn := NewMIMETypeFunction()
	resp := &function.DefinitionResponse{}
	fn.Definition(context.Background(), function.DefinitionRequest{}, resp)

	if resp.Definition.Return == nil {
		t.Fatal("expected return definition, got nil")
	}
	if len(resp.Definition.Parameters) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}
}
