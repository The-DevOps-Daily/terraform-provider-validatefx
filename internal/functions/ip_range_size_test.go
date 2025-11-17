package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestIPRangeSizeFunction_Valid(t *testing.T) {
	t.Parallel()
	fn := NewIPRangeSizeFunction()
	ctx := context.Background()

	// cidr, min, max
	args := function.NewArgumentsData([]attr.Value{
		types.StringValue("10.0.0.0/16"),
		types.Int64Value(8),
		types.Int64Value(28),
	})

	runResp := &function.RunResponse{}
	fn.Run(ctx, function.RunRequest{Arguments: args}, runResp)
	if runResp.Error != nil {
		t.Fatalf("expected no error, got %v", runResp.Error)
	}
}

func TestIPRangeSizeFunction_OutOfRange(t *testing.T) {
	t.Parallel()
	fn := NewIPRangeSizeFunction()
	ctx := context.Background()

	args := function.NewArgumentsData([]attr.Value{
		types.StringValue("10.0.0.0/30"),
		types.Int64Value(8),
		types.Int64Value(28),
	})
	runResp := &function.RunResponse{}
	fn.Run(ctx, function.RunRequest{Arguments: args}, runResp)
	if runResp.Error == nil {
		t.Fatalf("expected error for /30 outside allowed range")
	}
}

func TestIPRangeSizeFunction_NullUnknown(t *testing.T) {
	t.Parallel()
	fn := NewIPRangeSizeFunction()
	ctx := context.Background()

	args := function.NewArgumentsData([]attr.Value{
		types.StringUnknown(),
		types.Int64Value(8),
		types.Int64Value(24),
	})
	runResp := &function.RunResponse{}
	fn.Run(ctx, function.RunRequest{Arguments: args}, runResp)
	if runResp.Error != nil {
		t.Fatalf("expected no error for unknown input, got %v", runResp.Error)
	}
}
