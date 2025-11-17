package functions

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

type ipRangeSizeFunction struct{}

var _ function.Function = (*ipRangeSizeFunction)(nil)

// NewIPRangeSizeFunction exposes the ip_range_size validator as a Terraform function.
func NewIPRangeSizeFunction() function.Function { return &ipRangeSizeFunction{} }

func (ipRangeSizeFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "ip_range_size"
}

func (ipRangeSizeFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Validate that a CIDR's prefix length falls within an allowed inclusive range.",
		MarkdownDescription: "Returns true when the input is a valid CIDR whose prefix length is within the provided [min,max] bounds.",
		Return:              function.BoolReturn{},
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:               "cidr",
				AllowNullValue:     true,
				AllowUnknownValues: true,
				Description:        "CIDR block to validate (IPv4 or IPv6).",
			},
			function.Int64Parameter{
				Name:        "min_prefix",
				Description: "Minimum allowed prefix length (inclusive).",
			},
			function.Int64Parameter{
				Name:        "max_prefix",
				Description: "Maximum allowed prefix length (inclusive).",
			},
		},
	}
}

func (ipRangeSizeFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var cidr types.String
	var min types.Int64
	var max types.Int64

	if err := req.Arguments.GetArgument(ctx, 0, &cidr); err != nil {
		resp.Error = err
		return
	}
	if err := req.Arguments.GetArgument(ctx, 1, &min); err != nil {
		resp.Error = err
		return
	}
	if err := req.Arguments.GetArgument(ctx, 2, &max); err != nil {
		resp.Error = err
		return
	}

	if cidr.IsNull() || cidr.IsUnknown() {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	// Build validator and execute
	minV := int(min.ValueInt64())
	maxV := int(max.ValueInt64())
	v := validators.NewIPRangeSizeValidator(minV, maxV)

	vr := frameworkvalidator.StringResponse{}
	v.ValidateString(ctx, frameworkvalidator.StringRequest{
		Path:        path.Root("cidr"),
		ConfigValue: cidr,
	}, &vr)

	if vr.Diagnostics.HasError() {
		diags := diag.Diagnostics{}
		diags.Append(vr.Diagnostics...)
		resp.Error = function.FuncErrorFromDiags(ctx, diags)
		return
	}

	resp.Result = function.NewResultData(basetypes.NewBoolValue(true))
}
