package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestResourceNameMetadata(t *testing.T) {
	t.Parallel()

	fn := NewResourceNameFunction()
	resp := &function.MetadataResponse{}
	fn.Metadata(context.Background(), function.MetadataRequest{}, resp)

	if resp.Name != "resource_name" {
		t.Fatalf("expected name resource_name, got %s", resp.Name)
	}
}

func TestResourceNameDefinition(t *testing.T) {
	t.Parallel()

	fn := NewResourceNameFunction()
	resp := &function.DefinitionResponse{}
	fn.Definition(context.Background(), function.DefinitionRequest{}, resp)

	if resp.Definition.Summary == "" {
		t.Fatal("expected non-empty Summary")
	}
	if len(resp.Definition.Parameters) != 1 {
		t.Fatalf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}
}

func TestResourceNameFunction(t *testing.T) {
	t.Parallel()

	fn := NewResourceNameFunction()
	ctx := context.Background()

	testCases := []struct {
		name          string
		value         attr.Value
		expectError   bool
		expectUnknown bool
		expectTrue    bool
	}{
		// Valid resource names
		{
			name:       "simple lowercase",
			value:      types.StringValue("myresource"),
			expectTrue: true,
		},
		{
			name:       "with underscores",
			value:      types.StringValue("my_resource_name"),
			expectTrue: true,
		},
		{
			name:       "with hyphens",
			value:      types.StringValue("my-resource-name"),
			expectTrue: true,
		},
		{
			name:       "with digits",
			value:      types.StringValue("resource123"),
			expectTrue: true,
		},
		{
			name:       "starts with underscore",
			value:      types.StringValue("_private_resource"),
			expectTrue: true,
		},
		{
			name:       "aws resource name",
			value:      types.StringValue("aws_s3_bucket_2024"),
			expectTrue: true,
		},

		// Invalid resource names
		{
			name:        "empty string",
			value:       types.StringValue(""),
			expectError: true,
		},
		{
			name:        "uppercase letters",
			value:       types.StringValue("MyResource"),
			expectError: true,
		},
		{
			name:        "starts with digit",
			value:       types.StringValue("1resource"),
			expectError: true,
		},
		{
			name:        "contains spaces",
			value:       types.StringValue("my resource"),
			expectError: true,
		},

		// Null and unknown
		{
			name:          "null value",
			value:         types.StringNull(),
			expectUnknown: true,
		},
		{
			name:          "unknown value",
			value:         types.StringUnknown(),
			expectUnknown: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resp := &function.RunResponse{}
			fn.Run(ctx, function.RunRequest{Arguments: function.NewArgumentsData([]attr.Value{tc.value})}, resp)

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
