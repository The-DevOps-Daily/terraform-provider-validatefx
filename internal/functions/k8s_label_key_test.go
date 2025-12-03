package functions

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestK8sLabelKeyFunction(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		wantError bool
	}{
		{
			name:      "valid simple key",
			value:     "app",
			wantError: false,
		},
		{
			name:      "valid key with prefix",
			value:     "example.com/app",
			wantError: false,
		},
		{
			name:      "valid kubernetes.io prefix",
			value:     "kubernetes.io/name",
			wantError: false,
		},
		{
			name:      "valid with dashes",
			value:     "app-name",
			wantError: false,
		},
		{
			name:      "empty key",
			value:     "",
			wantError: true,
		},
		{
			name:      "too many slashes",
			value:     "a/b/c",
			wantError: true,
		},
		{
			name:      "invalid prefix format",
			value:     "Example.com/app",
			wantError: true,
		},
		{
			name:      "name too long",
			value:     strings.Repeat("a", 64),
			wantError: true,
		},
		{
			name:      "invalid special chars",
			value:     "app@name",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := NewK8sLabelKeyFunction()
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
