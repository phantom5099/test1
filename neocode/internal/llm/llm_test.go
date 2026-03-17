package llm

import (
	cfgpkg "github.com/yourname/neocode/config"
	"testing"
)

func TestMockClient_Generate(t *testing.T) {
	cfg := &cfgpkg.Config{Mock: true}
	c := NewClient(cfg)
	resp, err := c.Generate("test prompt")
	if err != nil {
		t.Fatalf("Mock Generate returned error: %v", err)
	}
	if resp.Description == "" {
		t.Fatalf("mock response missing description")
	}
	if len(resp.Edits) == 0 {
		t.Fatalf("mock response should include at least one edit")
	}
}
