package llm

import (
	"github.com/yourname/neocode/config"
	"time"
)

// Edit 描述由 LLM 产生的对文件的操作
type Edit struct {
	Op      string // create|update|delete
	Path    string
	Content string
}

// LLMResponse 是 LLM 的结构化响应
type LLMResponse struct {
	Description string
	Edits       []Edit
}

// LLMIClient 定义一个从提示生成 LLM 响应的接口
type LLMIClient interface {
	Generate(prompt string) (LLMResponse, error)
}

// NewClient 返回一个 LLM 客户端。在 mock 模式下返回本地模拟实现
func NewClient(cfg *config.Config) LLMIClient {
	if cfg.Mock {
		return &MockLLM{}
	}
	return &HTTPClientLLM{Endpoint: cfg.LLMEndpoint, APIKey: cfg.APIKey, Timeout: 10 * time.Second}
}

// MockLLM 提供一个用于离线使用的简单本地 LLM 模拟
type MockLLM struct{}

func (m *MockLLM) Generate(prompt string) (LLMResponse, error) {
	// Very naive deterministic response for demonstration purposes
	return LLMResponse{
		Description: "Mock plan: create a sample.txt with content 'Hello neocode'",
		Edits: []Edit{
			{Op: "create", Path: "sample.txt", Content: "Hello neocode"},
		},
	}, nil
}

// HTTPClientLLM 是一个用于真实远程 LLM 的占位实现
type HTTPClientLLM struct {
	Endpoint string
	APIKey   string
	Timeout  time.Duration
}

func (h *HTTPClientLLM) Generate(prompt string) (LLMResponse, error) {
	// 在此 MVP 的本地实现中不进行网络调用
	return LLMResponse{Description: "Remote LLM not enabled in local build.", Edits: nil}, nil
}
