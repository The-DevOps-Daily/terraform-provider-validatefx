package functions

import (
	"context"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework/function"
)

type semverRangeFunction struct{}

var _ function.Function = (*semverRangeFunction)(nil)

func NewSemVerRangeFunction() function.Function { return &semverRangeFunction{} }

func (semverRangeFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "semver_range"
}

func (semverRangeFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	// Reuse the common single-string function wrapper to keep consistency
	fn := newStringValidationFunction(
		"semver_range",
		"Validate that a string is a valid semantic version range expression.",
		"Returns true when the input string represents a valid SemVer range, e.g., '>=1.0.0,<2.0.0'.",
		validators.SemVerRange(),
	)
	fn.Definition(context.Background(), function.DefinitionRequest{}, resp)
}

func (semverRangeFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	fn := newStringValidationFunction(
		"semver_range",
		"Validate that a string is a valid semantic version range expression.",
		"Returns true when the input string represents a valid SemVer range, e.g., '>=1.0.0,<2.0.0'.",
		validators.SemVerRange(),
	)
	fn.Run(ctx, req, resp)
}
