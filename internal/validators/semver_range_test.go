package validators

import (
    "context"
    "testing"

    "github.com/hashicorp/terraform-plugin-framework/path"
    frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

func TestSemVerRangeValidator(t *testing.T) {
    t.Parallel()
    v := SemVerRange()

    run := func(s types.String) *frameworkvalidator.StringResponse {
        req := frameworkvalidator.StringRequest{Path: path.Root("range"), ConfigValue: s}
        resp := &frameworkvalidator.StringResponse{}
        v.ValidateString(context.Background(), req, resp)
        return resp
    }

    cases := []struct {
        name    string
        val     types.String
        wantErr bool
    }{
        {"single comparator", types.StringValue(">=1.2.3"), false},
        {"multiple comparators", types.StringValue(">=1.0.0, <2.0.0"), false},
        {"with leading v", types.StringValue(">=v1.0.0,<=v1.5.0"), false},
        {"bad operator", types.StringValue("~1.0.0"), true},
        {"bad version", types.StringValue(">=1.0"), true},
        {"empty comparator", types.StringValue(">=1.0.0, , <2.0.0"), true},
        {"null", types.StringNull(), false},
        {"unknown", types.StringUnknown(), false},
    }

    for _, tc := range cases {
        tc := tc
        t.Run(tc.name, func(t *testing.T) {
            t.Parallel()
            resp := run(tc.val)
            if tc.wantErr && !resp.Diagnostics.HasError() {
                t.Fatalf("expected error")
            }
            if !tc.wantErr && resp.Diagnostics.HasError() {
                t.Fatalf("unexpected error: %v", resp.Diagnostics)
            }
        })
    }
}

