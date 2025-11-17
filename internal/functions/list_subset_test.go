package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestListSubsetFunction(t *testing.T) {
	t.Parallel()

	fn := NewListSubsetFunction()
	ctx := context.Background()

	list := basetypes.NewListValueMust(basetypes.StringType{}, []attr.Value{
		basetypes.NewStringValue("read"), basetypes.NewStringValue("write"),
	})
	allowed := basetypes.NewListValueMust(basetypes.StringType{}, []attr.Value{
		basetypes.NewStringValue("read"), basetypes.NewStringValue("write"), basetypes.NewStringValue("admin"),
	})

	// success
	resp := &function.RunResponse{}
	fn.Run(ctx, function.RunRequest{Arguments: function.NewArgumentsData([]attr.Value{list, allowed})}, resp)
	if resp.Error != nil {
		t.Fatalf("unexpected error: %s", resp.Error)
	}

	// failure
	bad := basetypes.NewListValueMust(basetypes.StringType{}, []attr.Value{
		basetypes.NewStringValue("read"), basetypes.NewStringValue("unknown"),
	})
	resp = &function.RunResponse{}
	fn.Run(ctx, function.RunRequest{Arguments: function.NewArgumentsData([]attr.Value{bad, allowed})}, resp)
	if resp.Error == nil {
		t.Fatalf("expected error for disallowed element")
	}
}
