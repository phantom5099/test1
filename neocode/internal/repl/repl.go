package repl

import (
	"bufio"
	"fmt"
	"github.com/yourname/neocode/config"
	"github.com/yourname/neocode/internal/edit"
	llm "github.com/yourname/neocode/internal/llm"
	"os"
	"strings"
)

// 运行 neocode 的 REPL
func Run(cfg *config.Config) error {
	client := llm.NewClient(cfg)
	editor := edit.NewEditor()
	reader := bufio.NewReader(os.Stdin)
	// 新的交互状态：允许用户分步控制计划与执行
	fmt.Println("neocode – 本地 CLI，按需访问 LLM（Mock 模式就绪）。输入 exit 退出。")
	var lastResp llm.LLMResponse
	hasPending := false
	for {
		fmt.Print("neocode> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		input := strings.TrimSpace(line)
		if input == "" {
			fmt.Println("输入为空，请输入自然语言描述。")
			continue
		}
		switch strings.ToLower(input) {
		case "exit", "quit":
			fmt.Println("bye")
			return nil
		case "plan", "preview":
			if !hasPending {
				fmt.Println("没有待执行的计划，请输入自然语言描述以生成计划。")
				continue
			}
			// 展示当前计划描述与改动
			fmt.Println("描述:")
			fmt.Println(lastResp.Description)
			if len(lastResp.Edits) == 0 {
				fmt.Println("拟执行改动为空。")
				continue
			}
			fmt.Println("拟执行的改动:")
			for i, e := range lastResp.Edits {
				fmt.Printf("  %d) %s %s\n", i+1, e.Op, e.Path)
			}
			continue
		case "apply":
			if !hasPending {
				fmt.Println("没有待执行的计划，无法应用。请输入描述以生成计划。")
				continue
			}
			applied, err := editor.ApplyEdits(lastResp)
			if err != nil {
				fmt.Println("执行错误:", err)
				// 清理状态，允许重新尝试
				hasPending = false
				continue
			}
			fmt.Println("已应用:")
			fmt.Println(strings.Join(applied, ", "))
			// 清空待执行计划
			hasPending = false
			lastResp = llm.LLMResponse{}
			continue
		default:
			// 将新输入视为自然语言需求，触发 LLM 生成
			resp, err := client.Generate(input)
			if err != nil {
				fmt.Println("LLM 错误:", err)
				continue
			}
			lastResp = resp
			hasPending = true
			// 展示描述与拟执行改动
			fmt.Println("描述:")
			fmt.Println(resp.Description)
			if len(resp.Edits) == 0 {
				fmt.Println("未提出改动。")
				continue
			}
			fmt.Println("拟执行的改动:")
			for i, e := range resp.Edits {
				fmt.Printf("  %d) %s %s\n", i+1, e.Op, e.Path)
			}
			fmt.Println("提示：输入 'plan'/'preview' 查看计划，输入 'apply' 立即应用。")
		}
	}
}
