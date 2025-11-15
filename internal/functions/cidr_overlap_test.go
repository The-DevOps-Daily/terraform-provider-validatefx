package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestCIDROverlapFunction_NoOverlap(t *testing.T) {
	fn := NewCIDROverlapFunction()
	ctx := context.Background()

	defResp := &function.DefinitionResponse{}
	fn.Definition(ctx, function.DefinitionRequest{}, defResp)

	args := function.NewArgumentsData([]attr.Value{basetypes.NewListValueMust(
		basetypes.StringType{},
		[]attr.Value{basetypes.NewStringValue("10.0.0.0/24"), basetypes.NewStringValue("10.0.1.0/24")},
	)})

	runResp := &function.RunResponse{}
	fn.Run(ctx, function.RunRequest{Arguments: args}, runResp)
	if runResp.Error != nil {
		t.Fatalf("unexpected error: %v", runResp.Error)
	}
}

func TestCIDROverlapFunction_Overlap(t *testing.T) {
	fn := NewCIDROverlapFunction()
	ctx := context.Background()

	args := function.NewArgumentsData([]attr.Value{basetypes.NewListValueMust(
		basetypes.StringType{},
		[]attr.Value{basetypes.NewStringValue("192.168.0.0/24"), basetypes.NewStringValue("192.168.0.0/25")},
	)})

	runResp := &function.RunResponse{}
	fn.Run(ctx, function.RunRequest{Arguments: args}, runResp)
	if runResp.Error == nil {
		t.Fatalf("expected error for overlapping CIDRs")
	}
}
