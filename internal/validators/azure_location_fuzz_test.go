package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FuzzAzureLocation(f *testing.F) {
	// Seed with valid locations
	f.Add("eastus")
	f.Add("westeurope")
	f.Add("southeastasia")
	// Seed with invalid locations
	f.Add("invalid-location")
	f.Add("us-east-1")
	f.Add("")

	v := AzureLocation()
	f.Fuzz(func(t *testing.T, location string) {
		req := frameworkvalidator.StringRequest{
			Path:        path.Root("location"),
			ConfigValue: types.StringValue(location),
		}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		// Should not panic
	})
}
