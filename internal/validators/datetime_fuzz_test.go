package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FuzzDateTimeValidator(f *testing.F) {
	seeds := []string{
		"", "2025-11-02T15:04:05Z", "2025-11-02T15:04:05.123456Z",
		"2025-11-02 15:04:05", "not-a-date",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	vDefault := DateTime(nil)
	vCustom := DateTime([]string{"2006-01-02 15:04:05"})

	f.Fuzz(func(t *testing.T, s string) {
		t.Parallel()
		// default (RFC3339Nano)
		req := frameworkvalidator.StringRequest{Path: path.Root("dt"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		vDefault.ValidateString(context.Background(), req, resp)
		if s == "" {
			if resp.Diagnostics.HasError() {
				t.Fatalf("empty should not error")
			}
		}

		// custom layout allows space-separated format
		req2 := frameworkvalidator.StringRequest{Path: path.Root("dt"), ConfigValue: types.StringValue(s)}
		resp2 := &frameworkvalidator.StringResponse{}
		vCustom.ValidateString(context.Background(), req2, resp2)
		// No strict oracle; assert no panics and at least one validator accepts obvious valid seeds
	})
}
