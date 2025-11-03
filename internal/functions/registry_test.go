package functions

import (
	"context"
	"testing"
)

func TestProviderFunctionFactories(t *testing.T) {
	t.Parallel()

	factories := ProviderFunctionFactories()

	if len(factories) == 0 {
		t.Fatal("expected registered function factories")
	}

	for i, factory := range factories {
		if factory == nil {
			t.Fatalf("factory %d is nil", i)
		}

		if factory() == nil {
			t.Fatalf("factory %d returned nil function", i)
		}
	}
}

func TestAvailableFunctionDocs(t *testing.T) {
	t.Parallel()

	docs, err := AvailableFunctionDocs(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(docs) == 0 {
		t.Fatal("expected function docs metadata")
	}

	for _, doc := range docs {
		if doc.Name == "" {
			t.Fatal("function doc missing name")
		}
		if doc.Description == "" {
			t.Fatalf("function %s missing description", doc.Name)
		}
	}
}
