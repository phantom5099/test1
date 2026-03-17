package integration

import (
	"os"
	"testing"

	editpkg "github.com/yourname/neocode/internal/edit"
	llm "github.com/yourname/neocode/internal/llm"
	meta "github.com/yourname/neocode/internal/meta"
)

func TestIntegrationBasicWorkflow(t *testing.T) {
	tmp := t.TempDir()
	orig, _ := os.Getwd()
	_ = os.Chdir(tmp)
	defer func() { _ = os.Chdir(orig) }()

	// Prepare a mock LLM response with a create operation
	resp := llm.LLMResponse{
		Description: "Mock create demo.txt",
		Edits:       []llm.Edit{{Op: "create", Path: "demo.txt", Content: "Hello Neocode"}},
	}

	editor := editpkg.NewEditor()
	applied, err := editor.ApplyEdits(resp)
	if err != nil {
		t.Fatalf("ApplyEdits failed: %v", err)
	}
	if len(applied) != 1 || applied[0] != "创建了 demo.txt" {
		t.Fatalf("Unexpected apply summary: %v", applied)
	}

	// Validate file created with expected content
	data, err := os.ReadFile("demo.txt")
	if err != nil {
		t.Fatalf("Failed reading demo.txt: %v", err)
	}
	if string(data) != "Hello Neocode" {
		t.Fatalf("Unexpected file content: %q", string(data))
	}

	// Validate history was appended
	if err := meta.AppendHistory("integration test"); err != nil {
		t.Fatalf("AppendHistory failed: %v", err)
	}
	hist, err := meta.LoadHistory()
	if err != nil {
		t.Fatalf("LoadHistory failed: %v", err)
	}
	if len(hist) == 0 {
		t.Fatalf("History should not be empty after append")
	}
}
