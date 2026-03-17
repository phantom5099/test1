package edit

import (
	"fmt"
	"github.com/yourname/neocode/internal/fs"
	"github.com/yourname/neocode/internal/llm"
	"github.com/yourname/neocode/internal/meta"
	"os"
	"strings"
)

type Editor struct{}

func NewEditor() *Editor { return &Editor{} }

// ApplyEdits 将 LLM 派生的改动应用到本地文件系统
func (e *Editor) ApplyEdits(plan llm.LLMResponse) ([]string, error) {
	var applied []string
	for _, ed := range plan.Edits {
		path := strings.TrimSpace(ed.Path)
		switch ed.Op {
		case "create":
			if err := fs.WriteFileAtomic(path, []byte(ed.Content)); err != nil {
				return applied, err
			}
			applied = append(applied, fmt.Sprintf("创建了 %s", path))
		case "update":
			if err := fs.WriteFileAtomic(path, []byte(ed.Content)); err != nil {
				return applied, err
			}
			applied = append(applied, fmt.Sprintf("更新了 %s", path))
		case "delete":
			if err := os.Remove(path); err != nil {
				return applied, err
			}
			applied = append(applied, fmt.Sprintf("删除了 %s", path))
		default:
			// Unknown operation; skip
		}
	}
	// 将此次应用的摘要记录到历史中，便于追踪与回滚
	if len(applied) > 0 {
		entry := "应用改动: " + strings.Join(applied, "; ")
		_ = meta.AppendHistory(entry)
	}
	return applied, nil
}
