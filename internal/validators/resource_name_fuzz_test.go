//go:build gofuzz
// +build gofuzz

package validators

import (
	"context"
	"testing"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FuzzResourceName(f *testing.F) {
	// Valid seeds
	f.Add("myresource")
	f.Add("my_resource")
	f.Add("my-resource")
	f.Add("_private")
	f.Add("resource123")
	f.Add("aws_s3_bucket")

	// Invalid seeds
	f.Add("MyResource")
	f.Add("1resource")
	f.Add("-resource")
	f.Add("my resource")
	f.Add("my.resource")
	f.Add("")

	f.Fuzz(func(t *testing.T, value string) {
		validator := ResourceName()
		resp := &frameworkvalidator.StringResponse{}
		validator.ValidateString(context.Background(), frameworkvalidator.StringRequest{
			ConfigValue: types.StringValue(value),
		}, resp)
		// Just ensure it doesn't panic
	})
}
