package functions

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestSetProviderVersion(t *testing.T) {
	t.Parallel()

	SetProviderVersion("1.2.3")

	if !providerVersion.Equal(basetypes.NewStringValue("1.2.3")) {
		t.Fatalf("expected version 1.2.3, got %s", providerVersion.ValueString())
	}

	SetProviderVersion("")

	if !providerVersion.Equal(basetypes.NewStringValue("dev")) {
		t.Fatalf("expected fallback dev version, got %s", providerVersion.ValueString())
	}
}

func TestVersionFunctionRun(t *testing.T) {
	t.Parallel()

	SetProviderVersion("9.9.9")

	fn := NewVersionFunction()

	resp := &function.RunResponse{}
	fn.Run(context.Background(), function.RunRequest{}, resp)

	if resp.Error != nil {
		t.Fatalf("unexpected error: %s", resp.Error)
	}

	result := resp.Result.Value()

	strVal, ok := result.(basetypes.StringValue)
	if !ok {
		t.Fatalf("unexpected result type %T", result)
	}

	if strVal.ValueString() != "9.9.9" {
		t.Fatalf("expected version 9.9.9, got %s", strVal.ValueString())
	}
}

func TestProviderVersionGetter(t *testing.T) {
	// Test that ProviderVersion() returns the current version
	original := providerVersion
	defer func() { providerVersion = original }()

	SetProviderVersion("1.0.0")
	got := ProviderVersion()
	expected := basetypes.NewStringValue("1.0.0")

	if diff := cmp.Diff(expected, got); diff != "" {
		t.Errorf("ProviderVersion mismatch (-expected +got): %s", diff)
	}
}

func TestVersionFunctionMetadata(t *testing.T) {
	fn := NewVersionFunction()
	resp := &function.MetadataResponse{}
	fn.Metadata(context.Background(), function.MetadataRequest{}, resp)

	if resp.Name != "version" {
		t.Errorf("expected name 'version', got %q", resp.Name)
	}
}

func TestVersionFunctionDefinition(t *testing.T) {
	fn := NewVersionFunction()
	resp := &function.DefinitionResponse{}
	fn.Definition(context.Background(), function.DefinitionRequest{}, resp)

	if resp.Definition.Return == nil {
		t.Fatal("expected return definition, got nil")
	}
	if resp.Definition.Summary == "" {
		t.Error("expected non-empty summary")
	}
}
