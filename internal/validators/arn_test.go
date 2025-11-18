package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestARNValidatorValid(t *testing.T) {
	t.Parallel()
	v := ARN()
	cases := []string{
		"arn:aws:iam::123456789012:role/Admin",
		"arn:aws:s3:::my-bucket",
		"arn:aws:lambda:us-east-1:123456789012:function:my-func",
	}
	for _, s := range cases {
		req := frameworkvalidator.StringRequest{Path: path.Root("arn"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("expected no error for %q: %v", s, resp.Diagnostics)
		}
	}
}

func TestARNValidatorInvalid(t *testing.T) {
	t.Parallel()
	v := ARN()
	cases := []string{
		"not-an-arn",
		"arn:aws:ec2:us-east-1:abc:function:x", // bad account
		"arn:aws:::::",
	}
	for _, s := range cases {
		req := frameworkvalidator.StringRequest{Path: path.Root("arn"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		if !resp.Diagnostics.HasError() {
			t.Fatalf("expected error for %q", s)
		}
	}
}

func TestARNValidatorNullUnknown(t *testing.T) {
	t.Parallel()
	v := ARN()
	cases := []frameworkvalidator.StringRequest{
		{Path: path.Root("arn"), ConfigValue: types.StringNull()},
		{Path: path.Root("arn"), ConfigValue: types.StringUnknown()},
		{Path: path.Root("arn"), ConfigValue: types.StringValue("")},
	}
	for _, req := range cases {
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("expected no error for %v", req.ConfigValue)
		}
	}
}
