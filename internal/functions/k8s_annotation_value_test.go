package functions

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestK8sAnnotationValueFunction(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		wantError bool
	}{
		{
			name:      "empty annotation",
			value:     "",
			wantError: false,
		},
		{
			name:      "simple annotation",
			value:     "This is a simple annotation",
			wantError: false,
		},
		{
			name:      "annotation with special chars",
			value:     "annotation@example.com: value!",
			wantError: false,
		},
		{
			name:      "annotation with newlines",
			value:     "line1\nline2\nline3",
			wantError: false,
		},
		{
			name:      "large annotation",
			value:     strings.Repeat("a", 100000),
			wantError: false,
		},
		{
			name:      "annotation too large",
			value:     strings.Repeat("a", 262145),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := NewK8sAnnotationValueFunction()
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
