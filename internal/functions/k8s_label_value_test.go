package functions

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestK8sLabelValueFunction(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		wantError bool
	}{
		{
			name:      "valid simple value",
			value:     "production",
			wantError: false,
		},
		{
			name:      "empty value is valid",
			value:     "",
			wantError: false,
		},
		{
			name:      "valid with dash",
			value:     "prod-env",
			wantError: false,
		},
		{
			name:      "valid with underscore",
			value:     "prod_env",
			wantError: false,
		},
		{
			name:      "valid with dot",
			value:     "v1.0",
			wantError: false,
		},
		{
			name:      "valid with numbers",
			value:     "app123",
			wantError: false,
		},
		{
			name:      "valid with uppercase",
			value:     "Production",
			wantError: false,
		},
		{
			name:      "value too long",
			value:     strings.Repeat("a", 64),
			wantError: true,
		},
		{
			name:      "invalid special chars",
			value:     "prod@env",
			wantError: true,
		},
		{
			name:      "starts with dash",
			value:     "-production",
			wantError: true,
		},
		{
			name:      "ends with dash",
			value:     "production-",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := NewK8sLabelValueFunction()
			ctx := context.Background()
			req := function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue(tt.value),
				}),
			}
			resp := &function.RunResponse{}

			fn.Run(ctx, req, resp)

			if tt.wantError {
				if resp.Error == nil {
					t.Errorf("expected error for value %q, but got none", tt.value)
				}
			} else {
				if resp.Error != nil {
					t.Errorf("unexpected error for value %q: %v", tt.value, resp.Error)
				}
			}
		})
	}
}
