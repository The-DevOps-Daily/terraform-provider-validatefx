package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FuzzBetweenValidator(f *testing.F) {
	seeds := []struct{ val, min, max string }{
		{"7.5", "5", "10"},
		{"-1", "-5", "5"},
		{"11", "5", "10"},
		{"3.14", "", ""},
		{"notnum", "0", "1"},
	}
	for _, s := range seeds {
		f.Add(s.val, s.min, s.max)
	}

	f.Fuzz(func(t *testing.T, val, min, max string) {
		t.Parallel()
		v := Between(min, max)
		req := frameworkvalidator.StringRequest{Path: path.Root("between"), ConfigValue: types.StringValue(val)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		// Robustness check only; EvaluateBetween tested separately.
	})
}
