package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestARNFunction(t *testing.T) {
	t.Parallel()
	fn := NewARNFunction()
	ctx := context.Background()

	cases := []struct {
		name          string
		value         attr.Value
		expectError   bool
		expectUnknown bool
		expectTrue    bool
	}{
		{name: "valid iam role", value: types.StringValue("arn:aws:iam::123456789012:role/Admin"), expectTrue: true},
		{name: "invalid format", value: types.StringValue("not-an-arn"), expectError: true},
		{name: "null", value: types.StringNull(), expectUnknown: true},
		{name: "unknown", value: types.StringUnknown(), expectUnknown: true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			resp := &function.RunResponse{}
			fn.Run(ctx, function.RunRequest{Arguments: function.NewArgumentsData([]attr.Value{tc.value})}, resp)

			if tc.expectError {
				if resp.Error == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if resp.Error != nil {
				t.Fatalf("unexpected error: %v", resp.Error)
			}
			b, ok := resp.Result.Value().(basetypes.BoolValue)
			if !ok {
				t.Fatalf("unexpected result type %T", resp.Result.Value())
			}
			if tc.expectUnknown {
				if !b.IsUnknown() {
					t.Fatalf("expected unknown")
				}
				return
			}
			if b.IsUnknown() {
				t.Fatalf("did not expect unknown")
			}
			if !b.ValueBool() {
				t.Fatalf("expected true result")
			}
		})
	}
}

func TestARNFunction_Metadata(t *testing.T) {
	fn := NewARNFunction()
	resp := &function.MetadataResponse{}
	fn.Metadata(context.Background(), function.MetadataRequest{}, resp)

	if resp.Name != "arn" {
		t.Errorf("expected name 'arn', got %q", resp.Name)
	}
}

func TestARNFunction_Definition(t *testing.T) {
	fn := NewARNFunction()
	resp := &function.DefinitionResponse{}
	fn.Definition(context.Background(), function.DefinitionRequest{}, resp)

	if resp.Definition.Return == nil {
		t.Fatal("expected return definition, got nil")
	}
	if len(resp.Definition.Parameters) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}
}

func TestARNFunction_AdditionalValidCases(t *testing.T) {
	t.Parallel()

	fn := NewARNFunction()
	ctx := context.Background()

	cases := []struct {
		name  string
		value string
	}{
		{name: "s3 bucket", value: "arn:aws:s3:::my-bucket"},
		{name: "s3 object", value: "arn:aws:s3:::my-bucket/path/to/object"},
		{name: "lambda function", value: "arn:aws:lambda:us-east-1:123456789012:function:my-func"},
		{name: "lambda with version", value: "arn:aws:lambda:us-west-2:123456789012:function:my-func:1"},
		{name: "iam user", value: "arn:aws:iam::123456789012:user/alice"},
		{name: "iam group", value: "arn:aws:iam::123456789012:group/developers"},
		{name: "iam policy", value: "arn:aws:iam::123456789012:policy/ReadOnlyAccess"},
		{name: "iam instance profile", value: "arn:aws:iam::123456789012:instance-profile/WebServer"},
		{name: "gov partition", value: "arn:aws-us-gov:iam::123456789012:role/Admin"},
		{name: "china partition", value: "arn:aws-cn:iam::123456789012:role/Admin"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			resp := &function.RunResponse{}
			fn.Run(ctx, function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{types.StringValue(tc.value)}),
			}, resp)

			if resp.Error != nil {
				t.Fatalf("unexpected error for %q: %v", tc.value, resp.Error)
			}

			b, ok := resp.Result.Value().(basetypes.BoolValue)
			if !ok {
				t.Fatalf("unexpected result type %T", resp.Result.Value())
			}
			if !b.ValueBool() {
				t.Fatalf("expected true for valid ARN %q", tc.value)
			}
		})
	}
}

func TestARNFunction_InvalidCases(t *testing.T) {
	t.Parallel()

	fn := NewARNFunction()
	ctx := context.Background()

	cases := []struct {
		name  string
		value string
	}{
		{name: "empty resource", value: "arn:aws:s3:::"},
		{name: "missing parts", value: "arn:aws:s3"},
		{name: "s3 with region", value: "arn:aws:s3:us-east-1::my-bucket"},
		{name: "s3 with account", value: "arn:aws:s3::123456789012:my-bucket"},
		{name: "iam with region", value: "arn:aws:iam:us-east-1:123456789012:role/Admin"},
		{name: "iam invalid account", value: "arn:aws:iam::abc:role/Admin"},
		{name: "iam invalid resource", value: "arn:aws:iam::123456789012:invalid/Admin"},
		{name: "lambda no region", value: "arn:aws:lambda::123456789012:function:my-func"},
		{name: "lambda no account", value: "arn:aws:lambda:us-east-1::function:my-func"},
		{name: "lambda invalid resource", value: "arn:aws:lambda:us-east-1:123456789012:layer:my-layer"},
		{name: "invalid region format", value: "arn:aws:ec2:invalid:123456789012:instance/i-123"},
		{name: "invalid account digits", value: "arn:aws:ec2:us-east-1:abc:instance/i-123"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			resp := &function.RunResponse{}
			fn.Run(ctx, function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{types.StringValue(tc.value)}),
			}, resp)

			if resp.Error == nil {
				t.Fatalf("expected error for invalid ARN %q", tc.value)
			}
		})
	}
}
