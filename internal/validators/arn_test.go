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
		// additional valid permutations
		"arn:aws:iam::123456789012:user/alice",
		"arn:aws:iam::123456789012:policy/ReadOnlyAccess",
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
		"arn:aws:ec2:us-east-1:abc:function:x", // bad account digits
		"arn:aws:::::",
		"arn:aws:iam:us-east-1:123456789012:role/Admin", // iam must have empty region
		"arn:aws:s3:us-east-1:123456789012:bucket",      // s3 must have empty region/account
		"arn:aws:lambda::123456789012:function:x",       // lambda must have region
		"arn:aws:lambda:us-east-1::function:x",          // lambda must have account
		"arn:aws:lambda:us-east-1:123456789012:layer:x", // wrong resource prefix for lambda
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
