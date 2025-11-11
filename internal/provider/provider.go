package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/functions"
)

var (
	_ provider.Provider              = &validateFXProvider{}
	_ provider.ProviderWithFunctions = &validateFXProvider{}
)

// validateFXProvider defines the ValidateFX Terraform provider implementation.
type validateFXProvider struct {
	version string
}

// New returns a new instance of the ValidateFX provider factory function.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		functions.SetProviderVersion(version)
		return &validateFXProvider{
			version: version,
		}
	}
}

func (p *validateFXProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "validatefx"
	resp.Version = p.version
}

func (p *validateFXProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The validatefx provider exposes a suite of reusable validation functions that can be invoked from Terraform expressions using the `provider::validatefx::<name>` syntax.",
		// Non-intrusive: leave provider schema minimal for now.
	}
}

func (p *validateFXProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

func (p *validateFXProvider) Resources(ctx context.Context) []func() resource.Resource {
	return nil
}

func (p *validateFXProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

func (p *validateFXProvider) Functions(ctx context.Context) []func() function.Function {
	return functions.ProviderFunctionFactories()
}
