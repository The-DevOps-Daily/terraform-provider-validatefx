package functions

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestCreditCardExpiryFunction(t *testing.T) {
	t.Parallel()

	fn := NewCreditCardExpiryFunction()
	ctx := context.Background()

	// Get current date for testing
	now := time.Now()
	currentYear := now.Year()
	currentMonth := int(now.Month())

	// Future dates
	nextMonth := currentMonth + 1
	nextYear := currentYear
	if nextMonth > 12 {
		nextMonth = 1
		nextYear++
	}

	cases := []struct {
		name          string
		value         attr.Value
		expectError   bool
		expectUnknown bool
		expectTrue    bool
	}{
		{
			name:       "valid MM/YY format - next month",
			value:      types.StringValue(fmt.Sprintf("%02d/%02d", nextMonth, nextYear%100)),
			expectTrue: true,
		},
		{
			name:       "valid MM/YYYY format - next year",
			value:      types.StringValue(fmt.Sprintf("%02d/%d", currentMonth, currentYear+1)),
			expectTrue: true,
		},
		{
			name:       "valid far future",
			value:      types.StringValue("12/2099"),
			expectTrue: true,
		},
		{
			name:        "invalid format - single digit month",
			value:       types.StringValue("1/25"),
			expectError: true,
		},
		{
			name:        "invalid format - no slash",
			value:       types.StringValue("0125"),
			expectError: true,
		},
		{
			name:        "invalid month - 00",
			value:       types.StringValue("00/25"),
			expectError: true,
		},
		{
			name:        "invalid month - 13",
			value:       types.StringValue("13/25"),
			expectError: true,
		},
		{
			name:        "past date",
			value:       types.StringValue("12/2020"),
			expectError: true,
		},
		{
			name:          "null",
			value:         types.StringNull(),
			expectUnknown: true,
		},
		{
			name:          "unknown",
			value:         types.StringUnknown(),
			expectUnknown: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resp := &function.RunResponse{}
			fn.Run(ctx, function.RunRequest{Arguments: function.NewArgumentsData([]attr.Value{tc.value})}, resp)

			if tc.expectError {
				if resp.Error == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if resp.Error != nil {
				t.Fatalf("unexpected error: %s", resp.Error)
			}

			boolVal, ok := resp.Result.Value().(basetypes.BoolValue)
			if !ok {
				t.Fatalf("unexpected result type %T", resp.Result.Value())
			}

			if tc.expectUnknown {
				if !boolVal.IsUnknown() {
					t.Fatalf("expected unknown result")
				}
				return
			}

			if boolVal.IsUnknown() {
				t.Fatalf("did not expect unknown result")
			}

			if boolVal.ValueBool() != tc.expectTrue {
				t.Fatalf("expected %t, got %t", tc.expectTrue, boolVal.ValueBool())
			}
		})
	}
}
