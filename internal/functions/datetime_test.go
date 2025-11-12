package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestDateTimeFunction(t *testing.T) {
	t.Parallel()

	fn := NewDateTimeFunction()
	ctx := context.Background()

	noLayouts := basetypes.NewListNull(basetypes.StringType{})
	customLayouts := basetypes.NewListValueMust(
		basetypes.StringType{},
		[]attr.Value{types.StringValue("2006-01-02 15:04:05")},
	)

	cases := []struct {
		name          string
		args          []attr.Value
		expectError   bool
		expectUnknown bool
		expectTrue    bool
	}{
		{
			name: "rfc3339",
			args: []attr.Value{
				types.StringValue("2025-11-02T15:04:05Z"),
				noLayouts,
			},
			expectTrue: true,
		},
		{
			name: "custom layout",
			args: []attr.Value{
				types.StringValue("2025-11-02 15:04:05"),
				customLayouts,
			},
			expectTrue: true,
		},
		{
			name: "invalid date",
			args: []attr.Value{
				types.StringValue("2025-13-02T15:04:05Z"),
				noLayouts,
			},
			expectError: true,
		},
		{
			name: "invalid layouts",
			args: []attr.Value{
				types.StringValue("2025-11-02T15:04:05Z"),
				types.ListValueMust(
					types.BoolType,
					[]attr.Value{types.BoolValue(true)},
				),
			},
			expectError: true,
		},
		{
			name: "unknown value",
			args: []attr.Value{
				types.StringUnknown(),
				noLayouts,
			},
			expectUnknown: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resp := &function.RunResponse{}
			fn.Run(ctx, function.RunRequest{Arguments: function.NewArgumentsData(tc.args)}, resp)

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

func TestDateTimeFunction_ProviderDefaults(t *testing.T) {
	// Do not run in parallel: adjusts global provider configuration.
	orig := GetProviderConfiguration()
	t.Cleanup(func() { SetProviderConfiguration(orig) })
	SetProviderConfiguration(ProviderConfiguration{DatetimeLayouts: []string{"2006-01-02 15:04"}})

	fn := NewDateTimeFunction()
	ctx := context.Background()

	args := []attr.Value{
		types.StringValue("2025-11-02 15:04"),
		basetypes.NewListNull(basetypes.StringType{}),
	}

	resp := &function.RunResponse{}
	fn.Run(ctx, function.RunRequest{Arguments: function.NewArgumentsData(args)}, resp)

	if resp.Error != nil {
		t.Fatalf("unexpected error: %s", resp.Error)
	}

	boolVal, ok := resp.Result.Value().(basetypes.BoolValue)
	if !ok || boolVal.IsUnknown() || !boolVal.ValueBool() {
		t.Fatalf("expected true result using provider defaults, got %v", resp.Result.Value())
	}
}
