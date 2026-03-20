package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"go-llm-demo/internal/server/infra/provider"
	"go-llm-demo/internal/tui/infra"
)

const defaultHistoryTurns = 6

func main() {
	if err := loadDotEnv(".env"); err != nil {
		fmt.Printf("加载 .env 失败：%v\n", err)
		return
	}

	activeModel := strings.TrimSpace(os.Getenv("AI_MODEL"))
	if activeModel == "" {
		activeModel = provider.DefaultModel()
	}

	personaPrompt, err := loadPersonaPrompt(os.Getenv("PERSONA_FILE_PATH"))
	if err != nil {
		fmt.Printf("加载人设文件失败：%v\n", err)
		return
	}

	if activeModel == "" {
		fmt.Println("未配置可用模型")
		return
	}

	fmt.Println("=== NeoCode ===")
	fmt.Println("多行输入: 输入 ``` 开始代码块，输入 /send 发送，/cancel 取消")
	fmt.Println("命令: /switch <model> 切换模型, /run <code> 执行代码, /explain <code> 解释代码, /help 查看帮助")

	scanner := bufio.NewScanner(os.Stdin)
	ctx := context.Background()
	historyTurns := envInt("SHORT_TERM_HISTORY_TURNS", defaultHistoryTurns)
	history := initialHistory(personaPrompt, historyTurns)

	apiClient, err := infra.NewLocalChatClient()
	if err != nil {
		fmt.Printf("初始化失败：%v\n", err)
		return
	}

	for {
		fmt.Printf("[%s] > ", activeModel)

		input, err := readMultilineInput(scanner)
		if err != nil {
			if err == errEof {
				fmt.Println("\nExiting NeoCode")
				break
			}
			fmt.Printf("\n读取输入失败：%v\n", err)
			continue
		}

		line := strings.TrimSpace(input)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "/") {
			historyChanged := false
			shouldExit, err := handleCommand(ctx, line, &activeModel, &history, &historyChanged, personaPrompt, historyTurns, apiClient)
			if err != nil {
				fmt.Println(err)
			}
			if historyChanged {
				continue
			}
			if shouldExit {
				fmt.Println("Exiting NeoCode")
				break
			}
			continue
		}

		fmt.Println("Thinking...")
		messages := append([]infra.Message(nil), history...)
		messages = append(messages, infra.Message{Role: "user", Content: line})

		rep, err := apiClient.Chat(ctx, messages, activeModel)
		if err != nil {
			fmt.Printf("生成失败：%v\n", err)
			continue
		}

		var replyBuilder strings.Builder
		for msg := range rep {
			replyBuilder.WriteString(msg)
			fmt.Print(msg)
		}
		if replyBuilder.Len() > 0 {
			history = append(history, infra.Message{Role: "assistant", Content: replyBuilder.String()})
			history = trimHistory(history, historyTurns)
		}
		fmt.Println()
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("输入错误：%v\n", err)
	}
}

var errEof = fmt.Errorf("EOF")

type multilineState struct {
	active    bool
	lines     []string
	codeBlock bool
}

func readMultilineInput(scanner *bufio.Scanner) (string, error) {
	state := &multilineState{}
	var inputLines []string

	for {
		if !scanner.Scan() {
			if len(inputLines) > 0 {
				return strings.Join(inputLines, "\n"), nil
			}
			if err := scanner.Err(); err != nil {
				return "", err
			}
			return "", errEof
		}

		line := scanner.Text()

		if !state.active && strings.HasPrefix(line, "```") {
			state.active = true
			state.codeBlock = true
			state.lines = append(state.lines, line)
			fmt.Println(line)
			continue
		}

		if state.active && state.codeBlock {
			state.lines = append(state.lines, line)
			fmt.Println(line)

			if strings.TrimSpace(line) == "```" {
				state.active = false
				state.codeBlock = false
				return strings.Join(state.lines, "\n"), nil
			}
			continue
		}

		if !state.active {
			lower := strings.ToLower(strings.TrimSpace(line))
			if lower == "/send" {
				if len(inputLines) > 0 {
					return strings.Join(inputLines, "\n"), nil
				}
				continue
			}
			if lower == "/cancel" {
				inputLines = nil
				state.active = false
				state.lines = nil
				fmt.Println("已取消输入")
				continue
			}

			inputLines = append(inputLines, line)
			fmt.Println(line)

			if !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") && len(inputLines) == 1 {
				continue
			}

			if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
				continue
			}

			return strings.Join(inputLines, "\n"), nil
		}
	}
}

func handleCommand(ctx context.Context, input string, activeModel *string, history *[]infra.Message, historyChanged *bool, personaPrompt string, historyTurns int, client infra.ChatClient) (bool, error) {
	fields := strings.Fields(input)
	if len(fields) == 0 {
		return false, nil
	}

	switch fields[0] {
	case "/switch":
		if len(fields) < 2 {
			printAvailableModels()
			return false, fmt.Errorf("用法：/switch <model>")
		}

		target := fields[1]
		if !provider.IsSupportedModel(target) {
			printAvailableModels()
			return false, fmt.Errorf("模型不受支持：%s", target)
		}

		*activeModel = target
		fmt.Printf("已切换到模型：%s\n", target)

	case "/models":
		printAvailableModels()

	case "/run":
		if len(fields) < 2 {
			return false, fmt.Errorf("用法：/run <代码> 或粘贴代码后按回车发送")
		}
		code := strings.Join(fields[1:], " ")
		return false, runCode(code)

	case "/explain":
		if len(fields) < 2 {
			return false, fmt.Errorf("用法：/explain <代码> 或粘贴代码后按回车发送")
		}
		code := strings.Join(fields[1:], " ")
		return false, explainCode(ctx, code, client)

	case "/memory":
		stats, err := client.GetMemoryStats(ctx)
		if err != nil {
			return false, err
		}
		fmt.Printf("记忆条目: %d, TopK: %d, 最小分数: %.2f, 文件: %s\n",
			stats.Items, stats.TopK, stats.MinScore, stats.Path)

	case "/clear-memory":
		if err := client.ClearMemory(ctx); err != nil {
			return false, err
		}
		fmt.Println("已清空本地长期记忆")

	case "/clear-context":
		*history = initialHistory(personaPrompt, historyTurns)
		*historyChanged = true
		fmt.Println("已清空当前会话上下文")

	case "/help":
		printHelp()

	case "/exit":
		return true, nil

	case "/send":
		return false, nil

	case "/cancel":
		fmt.Println("无正在进行的输入")

	default:
		fmt.Printf("无法识别的命令：%s，输入 /help 查看帮助\n", fields[0])
	}

	return false, nil
}

func runCode(code string) error {
	ext, runner := detectLanguage(code)
	if ext == "" {
		return fmt.Errorf("无法识别代码语言，请使用 /explain 让 AI 解释")
	}

	tmpFile, err := os.CreateTemp("", "neocode-*."+ext)
	if err != nil {
		return fmt.Errorf("创建临时文件失败：%w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(code); err != nil {
		return fmt.Errorf("写入临时文件失败：%w", err)
	}
	tmpFile.Close()

	fmt.Printf("\n--- 运行 %s 代码 ---\n", ext)

	if runner != "" {
		cmd := exec.Command(runner, tmpFile.Name())
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		return cmd.Run()
	}

	cmd := exec.Command("go", "run", tmpFile.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func explainCode(ctx context.Context, code string, client infra.ChatClient) error {
	prompt := fmt.Sprintf("请解释以下代码的功能和工作原理（用中文回答，简洁清晰）：\n\n```\n%s\n```", code)

	messages := []infra.Message{
		{Role: "system", Content: "你是一个专业的编程助手，擅长解释代码逻辑。回答要简洁清晰，必要时可以给出示例。"},
		{Role: "user", Content: prompt},
	}

	rep, err := client.Chat(ctx, messages, provider.DefaultModel())
	if err != nil {
		return err
	}

	fmt.Println("\n--- 代码解释 ---")
	for msg := range rep {
		fmt.Print(msg)
	}
	fmt.Println("\n--- 解释结束 ---")
	return nil
}

func detectLanguage(code string) (string, string) {
	code = strings.TrimSpace(code)

	if strings.HasPrefix(code, "#!/bin/bash") || strings.HasPrefix(code, "#!/bin/sh") {
		return "sh", "bash"
	}
	if strings.HasPrefix(code, "package main") || strings.Contains(code, "func main()") {
		return "go", ""
	}
	if strings.HasPrefix(code, "<!DOCTYPE") || strings.HasPrefix(code, "<html") || strings.HasPrefix(code, "<div") {
		return "html", ""
	}
	if strings.HasPrefix(code, "<?php") {
		return "php", "php"
	}
	if strings.HasPrefix(code, "SELECT ") || strings.HasPrefix(code, "INSERT ") || strings.HasPrefix(code, "CREATE ") {
		return "sql", ""
	}
	if strings.HasPrefix(code, "def ") || strings.HasPrefix(code, "class ") || strings.Contains(code, "import ") && !strings.Contains(code, "fmt") && !strings.Contains(code, "go-llm") {
		return "py", "python"
	}
	if strings.HasPrefix(code, "fn ") || strings.HasPrefix(code, "let mut") || strings.HasPrefix(code, "impl ") {
		return "rs", "rustc"
	}
	if strings.HasPrefix(code, "console.log") || strings.HasPrefix(code, "const ") || strings.HasPrefix(code, "let ") && strings.Contains(code, "=>") {
		return "js", "node"
	}

	return "", ""
}

func printAvailableModels() {
	fmt.Println("可用模型:")
	for _, model := range provider.SupportedModels {
		fmt.Printf("  %s\n", model)
	}
}

func printHelp() {
	fmt.Println("命令:")
	fmt.Println("  /switch <model>   切换模型")
	fmt.Println("  /models           列出可用模型")
	fmt.Println("  /run <代码>       执行代码（支持 Go, Bash, Python, JS, PHP, Rust）")
	fmt.Println("  /explain <代码>   解释代码功能")
	fmt.Println("  /memory           显示本地记忆统计")
	fmt.Println("  /clear-memory     清空本地长期记忆")
	fmt.Println("  /clear-context    清空当前会话上下文")
	fmt.Println("  /exit             退出程序")
	fmt.Println("  /help             显示帮助")
	fmt.Println("")
	fmt.Println("多行输入:")
	fmt.Println("  输入 ``` 开始代码块输入")
	fmt.Println("  输入 /send 发送当前输入")
	fmt.Println("  输入 /cancel 取消输入")
}

func trimHistory(history []infra.Message, maxTurns int) []infra.Message {
	var systemMessages []infra.Message
	start := 0
	for start < len(history) && history[start].Role == "system" {
		systemMessages = append(systemMessages, history[start])
		start++
	}

	conversation := history[start:]
	maxMessages := maxTurns * 2
	if maxTurns <= 0 || len(conversation) <= maxMessages {
		return history
	}

	trimmed := append([]infra.Message(nil), systemMessages...)
	trimmed = append(trimmed, conversation[len(conversation)-maxMessages:]...)
	return trimmed
}

func initialHistory(personaPrompt string, historyTurns int) []infra.Message {
	history := make([]infra.Message, 0, historyTurns*2+1)
	if personaPrompt != "" {
		history = append(history, infra.Message{Role: "system", Content: personaPrompt})
	}
	return history
}

func loadPersonaPrompt(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}

func envInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed := 0
	for _, ch := range value {
		if ch < '0' || ch > '9' {
			return fallback
		}
		parsed = parsed*10 + int(ch-'0')
	}
	if parsed <= 0 {
		return fallback
	}
	return parsed
}

func loadDotEnv(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" || os.Getenv(key) != "" {
			continue
		}

		value = strings.Trim(value, `"'`)
		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}

	return scanner.Err()
}
