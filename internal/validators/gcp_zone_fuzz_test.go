package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FuzzGCPZone(f *testing.F) {
	// Seed with valid zones
	f.Add("us-central1-a")
	f.Add("us-central1-b")
	f.Add("us-central1-c")
	f.Add("us-east1-b")
	f.Add("us-west1-a")
	f.Add("europe-west1-b")
	f.Add("europe-west4-a")
	f.Add("asia-east1-a")
	f.Add("asia-northeast1-a")
	f.Add("australia-southeast1-a")
	f.Add("me-central1-a")
	f.Add("africa-south1-a")

	// Seed with invalid zones
	f.Add("")
	f.Add("invalid-zone")
	f.Add("us-central1")
	f.Add("us-central1-z")
	f.Add("us-east-1a")
	f.Add("not-a-zone")

	f.Fuzz(func(t *testing.T, input string) {
		v := GCPZone()
		req := validator.StringRequest{
			Path:        path.Root("test"),
			ConfigValue: types.StringValue(input),
		}
		resp := &validator.StringResponse{}

		// Should not panic on any input
		v.ValidateString(context.Background(), req, resp)
	})
}
