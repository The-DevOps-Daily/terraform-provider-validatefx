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

func TestCIDROverlapFunction_NullList(t *testing.T) {
	fn := NewCIDROverlapFunction()
	ctx := context.Background()

	args := function.NewArgumentsData([]attr.Value{basetypes.NewListNull(basetypes.StringType{})})

	runResp := &function.RunResponse{}
	fn.Run(ctx, function.RunRequest{Arguments: args}, runResp)
	if runResp.Error == nil {
		t.Fatalf("expected error for null list")
	}
}

func TestCIDROverlapFunction_UnknownList(t *testing.T) {
	fn := NewCIDROverlapFunction()
	ctx := context.Background()

	args := function.NewArgumentsData([]attr.Value{basetypes.NewListUnknown(basetypes.StringType{})})

	runResp := &function.RunResponse{}
	fn.Run(ctx, function.RunRequest{Arguments: args}, runResp)
	if runResp.Error != nil {
		t.Fatalf("unexpected error for unknown list: %v", runResp.Error)
	}
	// Should return unknown result
	result, ok := runResp.Result.Value().(basetypes.BoolValue)
	if !ok {
		t.Fatalf("expected BoolValue result")
	}
	if !result.IsUnknown() {
		t.Fatalf("expected unknown result for unknown list")
	}
}

func TestCIDROverlapFunction_UnknownElement(t *testing.T) {
	fn := NewCIDROverlapFunction()
	ctx := context.Background()

	args := function.NewArgumentsData([]attr.Value{basetypes.NewListValueMust(
		basetypes.StringType{},
		[]attr.Value{basetypes.NewStringValue("10.0.0.0/24"), basetypes.NewStringUnknown()},
	)})

	runResp := &function.RunResponse{}
	fn.Run(ctx, function.RunRequest{Arguments: args}, runResp)
	if runResp.Error != nil {
		t.Fatalf("unexpected error for unknown element: %v", runResp.Error)
	}
	// Should return unknown result
	result, ok := runResp.Result.Value().(basetypes.BoolValue)
	if !ok {
		t.Fatalf("expected BoolValue result")
	}
	if !result.IsUnknown() {
		t.Fatalf("expected unknown result when list contains unknown element")
	}
}

func TestCIDROverlapFunction_NullElement(t *testing.T) {
	fn := NewCIDROverlapFunction()
	ctx := context.Background()

	args := function.NewArgumentsData([]attr.Value{basetypes.NewListValueMust(
		basetypes.StringType{},
		[]attr.Value{basetypes.NewStringValue("10.0.0.0/24"), basetypes.NewStringNull()},
	)})

	runResp := &function.RunResponse{}
	fn.Run(ctx, function.RunRequest{Arguments: args}, runResp)
	if runResp.Error == nil {
		t.Fatalf("expected error for null element in list")
	}
}

func TestCIDROverlapFunction_InvalidCIDR(t *testing.T) {
	fn := NewCIDROverlapFunction()
	ctx := context.Background()

	args := function.NewArgumentsData([]attr.Value{basetypes.NewListValueMust(
		basetypes.StringType{},
		[]attr.Value{basetypes.NewStringValue("not-a-cidr"), basetypes.NewStringValue("10.0.0.0/24")},
	)})

	runResp := &function.RunResponse{}
	fn.Run(ctx, function.RunRequest{Arguments: args}, runResp)
	if runResp.Error == nil {
		t.Fatalf("expected error for invalid CIDR")
	}
}

func TestCIDROverlapFunction_EmptyList(t *testing.T) {
	fn := NewCIDROverlapFunction()
	ctx := context.Background()

	args := function.NewArgumentsData([]attr.Value{basetypes.NewListValueMust(
		basetypes.StringType{},
		[]attr.Value{},
	)})

	runResp := &function.RunResponse{}
	fn.Run(ctx, function.RunRequest{Arguments: args}, runResp)
	if runResp.Error != nil {
		t.Fatalf("unexpected error for empty list: %v", runResp.Error)
	}
	// Empty list should return true (no overlaps)
	result, ok := runResp.Result.Value().(basetypes.BoolValue)
	if !ok {
		t.Fatalf("expected BoolValue result")
	}
	if !result.ValueBool() {
		t.Fatalf("expected true for empty list")
	}
}

func TestCIDROverlapFunction_SingleCIDR(t *testing.T) {
	fn := NewCIDROverlapFunction()
	ctx := context.Background()

	args := function.NewArgumentsData([]attr.Value{basetypes.NewListValueMust(
		basetypes.StringType{},
		[]attr.Value{basetypes.NewStringValue("10.0.0.0/24")},
	)})

	runResp := &function.RunResponse{}
	fn.Run(ctx, function.RunRequest{Arguments: args}, runResp)
	if runResp.Error != nil {
		t.Fatalf("unexpected error for single CIDR: %v", runResp.Error)
	}
	// Single CIDR should return true (no overlaps possible)
	result, ok := runResp.Result.Value().(basetypes.BoolValue)
	if !ok {
		t.Fatalf("expected BoolValue result")
	}
	if !result.ValueBool() {
		t.Fatalf("expected true for single CIDR")
	}
}

func TestCIDROverlapFunction_Metadata(t *testing.T) {
	fn := NewCIDROverlapFunction()
	resp := &function.MetadataResponse{}
	fn.Metadata(context.Background(), function.MetadataRequest{}, resp)

	if resp.Name != "cidr_overlap" {
		t.Errorf("expected name 'cidr_overlap', got %q", resp.Name)
	}
}

func TestCIDROverlapFunction_Definition(t *testing.T) {
	fn := NewCIDROverlapFunction()
	resp := &function.DefinitionResponse{}
	fn.Definition(context.Background(), function.DefinitionRequest{}, resp)

	if resp.Definition.Return == nil {
		t.Fatal("expected return definition, got nil")
	}
	if len(resp.Definition.Parameters) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}
}
