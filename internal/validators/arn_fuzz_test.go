package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FuzzARNValidator(f *testing.F) {
	v := ARN()
	seeds := []string{
		"arn:aws:iam::123456789012:role/Admin",
		"arn:aws:s3:::bucket",
		"arn:aws:lambda:us-west-2:123456789012:function:func",
		"not-an-arn",
		"arn:aws:::::",
		"",
	}
	for _, s := range seeds {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		req := frameworkvalidator.StringRequest{Path: path.Root("arn"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
	})
}
