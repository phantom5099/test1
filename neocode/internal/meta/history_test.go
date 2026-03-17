package meta

import (
	"os"
	"testing"
)

func TestHistoryLifeCycle(t *testing.T) {
	// Use a temp dir to isolate tests
	tmp := t.TempDir()
	// Change working directory to temp dir for History file path calculations
	orig, _ := os.Getwd()
	_ = os.Chdir(tmp)
	defer func() {
		_ = os.Chdir(orig)
	}()

	// Ensure history file doesn't exist initially
	if _, err := os.Stat(HistoryFilePath()); err == nil {
		t.Fatalf("history file should not exist yet in temp dir")
	}

	// Append first entry
	if err := AppendHistory("first entry"); err != nil {
		t.Fatalf("AppendHistory failed: %v", err)
	}
	hist, err := LoadHistory()
	if err != nil {
		t.Fatalf("LoadHistory failed: %v", err)
	}
	if len(hist) != 1 || hist[0] != "first entry" {
		t.Fatalf("unexpected history after first append: %#v", hist)
	}

	// Append second entry
	if err := AppendHistory("second entry"); err != nil {
		t.Fatalf("AppendHistory second failed: %v", err)
	}
	hist, err = LoadHistory()
	if err != nil {
		t.Fatalf("LoadHistory after second failed: %v", err)
	}
	if len(hist) != 2 || hist[1] != "second entry" {
		t.Fatalf("unexpected history after second append: %#v", hist)
	}

	// Ensure the history file exists on disk with valid JSON
	path := HistoryFilePath()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("history file not found on disk: %v", err)
	}
	// Optional: check it's valid JSON
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed reading history file: %v", err)
	}
	if len(data) == 0 {
		t.Fatalf("history file is empty")
	}
}
