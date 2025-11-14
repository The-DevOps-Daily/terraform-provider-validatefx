package functions

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// createTempReadme writes the given content to a temporary README and returns its path and a cleanup func.
func createTempReadme(t *testing.T, content string) (string, func()) {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "README.md")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp readme: %v", err)
	}
	return p, func() { _ = os.RemoveAll(dir) }
}

func TestUpdateReadmeFunctionsTable_UpdatesSection(t *testing.T) {
	t.Parallel()

	// Minimal README content with the expected anchors used by UpdateReadmeFunctionsTable.
	original := strings.Join([]string{
		"# ValidateFX",
		"",
		"## 🧩 Available Functions",
		"",
		"| Function | Description |",
		"| -------------------------- | ------------------------------------------------ |",
		"| `old` | placeholder |",
		"",
		"---",
		"",
		"Other content...",
		"",
	}, "\n")

	path, _ := createTempReadme(t, original)

	// Run update
	if err := UpdateReadmeFunctionsTable(context.Background(), path); err != nil {
		t.Fatalf("update table: %v", err)
	}

	// Verify file content updated and includes some known function names.
	updated, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read updated: %v", err)
	}

	got := string(updated)
	if strings.Contains(got, "`old`") {
		t.Fatalf("expected placeholder row to be replaced")
	}
	// Spot-check a couple of known functions from the registry.
	if !strings.Contains(got, "`email`") || !strings.Contains(got, "`uuid`") {
		t.Fatalf("expected known functions in table, got: %s", got)
	}
}

func TestUpdateReadmeFunctionsTable_Idempotent(t *testing.T) {
	t.Parallel()

	base := strings.Join([]string{
		"# ValidateFX",
		"",
		"## 🧩 Available Functions",
		"",
		"| Function | Description |",
		"| -------------------------- | ------------------------------------------------ |",
		"| `placeholder` | will be replaced |",
		"",
		"---",
		"",
	}, "\n")

	path, _ := createTempReadme(t, base)

	if err := UpdateReadmeFunctionsTable(context.Background(), path); err != nil {
		t.Fatalf("first update: %v", err)
	}

	first, _ := os.ReadFile(path)

	if err := UpdateReadmeFunctionsTable(context.Background(), path); err != nil {
		t.Fatalf("second update: %v", err)
	}

	second, _ := os.ReadFile(path)
	if string(first) != string(second) {
		t.Fatalf("expected idempotent update; content changed on second run")
	}
}

func TestUpdateReadmeFunctionsTable_NoSection_NoChange(t *testing.T) {
	t.Parallel()

	content := "# Readme without functions section\n\nJust text.\n"
	path, _ := createTempReadme(t, content)

	if err := UpdateReadmeFunctionsTable(context.Background(), path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	after, _ := os.ReadFile(path)
	if string(after) != content {
		t.Fatalf("expected unchanged content when anchors absent")
	}
}
