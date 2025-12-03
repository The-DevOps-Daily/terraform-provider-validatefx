package functions

import (
	"context"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ function.Function = &K8sLabelKeyFunction{}

// K8sLabelKeyFunction validates Kubernetes label keys.
type K8sLabelKeyFunction struct{}

func NewK8sLabelKeyFunction() function.Function {
	return &K8sLabelKeyFunction{}
}

func (f *K8sLabelKeyFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "k8s_label_key"
}

func (f *K8sLabelKeyFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Validates Kubernetes label key format",
		Description:         "Validates that a string is a valid Kubernetes label key, optionally with a DNS subdomain prefix.",
		MarkdownDescription: "Validates that a string is a valid Kubernetes label key, optionally with a DNS subdomain prefix.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "value",
				Description: "The label key to validate (e.g., 'app' or 'example.com/app')",
			},
		},
		Return: function.BoolReturn{},
	}
}

func (f *K8sLabelKeyFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var value string

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &value))
	if resp.Error != nil {
		return
	}

	if err := validators.ValidateLabelKey(value); err != nil {
		resp.Error = function.NewArgumentFuncError(0, "Invalid Kubernetes Label Key: "+err.Error())
		return
	}

	resp.Result = function.NewResultData(basetypes.NewBoolValue(true))
}
