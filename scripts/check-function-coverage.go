package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/function"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/functions"
)

func main() {
	log.SetFlags(0)

	flag.Parse()

	targets := flag.Args()
	if len(targets) == 0 {
		log.Fatal("at least one target must be specified (examples, integration)")
	}

	checks := map[string]func([]string) []string{
		"examples":    func(names []string) []string { return validateExamples("examples/functions", names) },
		"integration": func(names []string) []string { return validateIntegration("integration/main.tf", names) },
	}

	ctx := context.Background()
	names, err := loadFunctionNames(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var failures []string
	for _, target := range targets {
		fn, ok := checks[target]
		if !ok {
			log.Fatalf("unknown target %q", target)
		}
		failures = append(failures, fn(names)...)
	}

	if len(failures) > 0 {
		log.Fatalf("function coverage validation failed:\n  - %s", strings.Join(failures, "\n  - "))
	}
}

func loadFunctionNames(ctx context.Context) ([]string, error) {
	factories := functions.ProviderFunctionFactories()
	names := make([]string, 0, len(factories))

	for _, factory := range factories {
		fn := factory()

		metaResp := &function.MetadataResponse{}
		fn.Metadata(ctx, function.MetadataRequest{}, metaResp)

		names = append(names, metaResp.Name)
	}

	sort.Strings(names)

	return names, nil
}

func validateExamples(root string, names []string) []string {
	var failures []string

	for _, name := range names {
		dir := filepath.Join(root, name)
		info, err := os.Stat(dir)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				failures = append(failures, fmt.Sprintf("missing examples directory for function %q (%s)", name, dir))
			} else {
				failures = append(failures, fmt.Sprintf("failed to stat %s for function %q: %v", dir, name, err))
			}
			continue
		}

		if !info.IsDir() {
			failures = append(failures, fmt.Sprintf("expected %s to be a directory for function %q", dir, name))
			continue
		}

		tfPath := filepath.Join(dir, "function.tf")
		if _, err := os.Stat(tfPath); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				failures = append(failures, fmt.Sprintf("missing function.tf for function %q (%s)", name, tfPath))
			} else {
				failures = append(failures, fmt.Sprintf("failed to stat %s for function %q: %v", tfPath, name, err))
			}
		}
	}

	return failures
}

func validateIntegration(path string, names []string) []string {
	contents, err := os.ReadFile(path)
	if err != nil {
		return []string{fmt.Sprintf("failed to read integration file %s: %v", path, err)}
	}

	text := string(contents)
	var failures []string

	for _, name := range names {
		needle := fmt.Sprintf("provider::validatefx::%s", name)
		if !strings.Contains(text, needle) {
			failures = append(failures, fmt.Sprintf("integration scenario missing function %q (%s)", name, needle))
		}
	}

	return failures
}
