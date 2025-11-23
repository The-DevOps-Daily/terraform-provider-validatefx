package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestSlugFunction_Metadata(t *testing.T) {
	t.Parallel()

	fn := NewSlugFunction()

	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	fn.Metadata(context.Background(), req, resp)

	if resp.Name != "slug" {
		t.Errorf("expected function name 'slug', got %q", resp.Name)
	}
}

func TestSlugFunction_Definition(t *testing.T) {
	t.Parallel()

	fn := NewSlugFunction()

	req := function.DefinitionRequest{}
	resp := &function.DefinitionResponse{}
	fn.Definition(context.Background(), req, resp)

	if resp.Definition.Summary == "" {
		t.Error("expected non-empty summary")
	}

	if resp.Definition.MarkdownDescription == "" {
		t.Error("expected non-empty description")
	}

	if len(resp.Definition.Parameters) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}
}

func TestSlugFunction_Run(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		value         attr.Value
		expectError   bool
		expectUnknown bool
		expectTrue    bool
	}{
		// Valid slugs
		{name: "simple slug", value: types.StringValue("hello"), expectTrue: true},
		{name: "slug with hyphen", value: types.StringValue("hello-world"), expectTrue: true},
		{name: "slug with digits", value: types.StringValue("hello-123"), expectTrue: true},
		{name: "slug starting with digit", value: types.StringValue("123-hello"), expectTrue: true},
		{name: "slug all digits", value: types.StringValue("123"), expectTrue: true},
		{name: "slug multiple hyphens", value: types.StringValue("my-awesome-slug"), expectTrue: true},
		{name: "slug with version", value: types.StringValue("my-app-v2"), expectTrue: true},
		{name: "slug resource name", value: types.StringValue("web-server-01"), expectTrue: true},

		// Invalid slugs
		{name: "empty string", value: types.StringValue(""), expectError: true},
		{name: "uppercase", value: types.StringValue("Hello"), expectError: true},
		{name: "leading hyphen", value: types.StringValue("-hello"), expectError: true},
		{name: "trailing hyphen", value: types.StringValue("hello-"), expectError: true},
		{name: "consecutive hyphens", value: types.StringValue("hello--world"), expectError: true},
		{name: "underscore", value: types.StringValue("hello_world"), expectError: true},
		{name: "space", value: types.StringValue("hello world"), expectError: true},
		{name: "special chars", value: types.StringValue("hello@world"), expectError: true},

		// Null and unknown
		{name: "null input", value: types.StringNull(), expectUnknown: true},
		{name: "unknown input", value: types.StringUnknown(), expectUnknown: true},
	}

	fn := NewSlugFunction()
	ctx := context.Background()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resp := &function.RunResponse{}
			fn.Run(ctx, function.RunRequest{Arguments: function.NewArgumentsData([]attr.Value{tt.value})}, resp)

			if tt.expectError {
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

			if tt.expectUnknown {
				if !boolVal.IsUnknown() {
					t.Fatalf("expected unknown result")
				}
				return
			}

			if boolVal.IsUnknown() {
				t.Fatalf("did not expect unknown result")
			}

			if boolVal.ValueBool() != tt.expectTrue {
				t.Fatalf("expected %t, got %t", tt.expectTrue, boolVal.ValueBool())
			}
		})
	}
}
