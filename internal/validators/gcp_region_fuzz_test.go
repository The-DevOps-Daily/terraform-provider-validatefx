package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FuzzGCPRegion(f *testing.F) {
	// Seed with valid regions
	f.Add("us-central1")
	f.Add("europe-west1")
	f.Add("asia-southeast1")
	// Seed with invalid regions
	f.Add("invalid-region")
	f.Add("us-east-1")
	f.Add("")

	v := GCPRegion()
	f.Fuzz(func(t *testing.T, region string) {
		req := frameworkvalidator.StringRequest{
			Path:        path.Root("region"),
			ConfigValue: types.StringValue(region),
		}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		// Should not panic
	})
}
