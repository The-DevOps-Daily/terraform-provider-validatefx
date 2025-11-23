package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FuzzUUIDv4Only(f *testing.F) {
	// Seed with valid UUIDv4
	f.Add("550e8400-e29b-41d4-a716-446655440000")
	f.Add("f47ac10b-58cc-4372-a567-0e02b2c3d479")
	f.Add("123e4567-e89b-42d3-a456-426614174000")
	f.Add("00000000-0000-4000-8000-000000000000")

	// Seed with invalid UUIDs (other versions)
	f.Add("6ba7b810-9dad-11d1-80b4-00c04fd430c8") // v1
	f.Add("6ba7b810-9dad-31d1-80b4-00c04fd430c8") // v3
	f.Add("6ba7b810-9dad-51d1-80b4-00c04fd430c8") // v5

	// Seed with invalid formats
	f.Add("")
	f.Add("not-a-uuid")
	f.Add("550e8400-e29b-41d4")

	f.Fuzz(func(t *testing.T, input string) {
		v := UUIDv4Only()
		req := validator.StringRequest{
			Path:        path.Root("test"),
			ConfigValue: types.StringValue(input),
		}
		resp := &validator.StringResponse{}

		// Should not panic on any input
		v.ValidateString(context.Background(), req, resp)
	})
}
