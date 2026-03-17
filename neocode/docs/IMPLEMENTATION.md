# neocode 实现手册 (MVP 本地 CLI + 本地化 LLM)

本手册面向开发者，给出可落地的实现细纲，帮助在本地快速构建并迭代 neocode MVP。

## 1. 项目目标
- 构建一个纯本地的命令行工具 neocode，核心能力：接收自然语言描述，通过提示词组装触发本地化的 LLM；LLM 返回文本描述及对本地文件的改动指令，CLI 在本地应用改动。无远程后端依赖，LLM 仅在需要时调用（默认本地 Mock）。
- 最小化实现成本、可快速部署、易于扩展。未来可逐步接入更丰富的 UX（如 TUI）和离线推理能力。

## 2. 总体架构
- 前端/CLI：命令行入口，REPL 模式，负责采集用户输入并与其他组件协作。
- LLM 层：提供 LLMIClient 接口，当前实现为 MockLLM（离线）+ HTTPClientLLM（未来接入真实远端）。
- 文件系统层：封装原子写入、备份、路径存在性检查等，确保改动落地的幂等性与可回滚性。
- 编辑执行层：将 LLM 的 Edits 应用到本地文件系统，支持 create/update/delete。
- 元数据层：占位，便于后续扩展历史、版本等能力。

数据流示意：用户输入 -> REPL -> LLM Client.Generate -> LLMResponse（Description + Edits） -> Editor.ApplyEdits -> FS 修改 -> 输出摘要

## 3. 数据模型与接口契约
- Edit: 由 LLM 返回，用于描述对单个文件的操作
  - Op: create|update|delete
  - Path: 文件路径
  - Content: 针对 create/update 的新内容
- LLMResponse: LLM 的结构化响应
  - Description: 任务描述性文本
  - Edits: Edit[]
- LLMIClient: 定义从提示生成 LLM 响应的接口
  - Generate(prompt string) (LLMResponse, error)
- Editor: 将 Edits 应用到本地文件系统
  - ApplyEdits(plan LLMResponse) ([]string, error) // 返回应用摘要
- FS（文件系统封装）
  - WriteFileAtomic(path string, data []byte) error
  - BackupFile(path string) (string, error)
  - ReadFile(path string) ([]byte, error)

示例代码片段（参考，实际实现见代码库相应文件）
```go
// neocode/internal/llm/llm.go（接口定义示例）
type LLMIClient interface {
    Generate(prompt string) (LLMResponse, error)
}
```

```go
// neocode/internal/llm/llm.go（LLMResponse/Edits 示例）
type Edit struct {
    Op      string
    Path    string
    Content string
}
type LLMResponse struct {
    Description string
    Edits       []Edit
}
```

## 4. 交互流程（CLI 流程）
- 启动 neocode，进入 REPL。
- 用户输入自然语言描述，例如：创建一个 hello.txt，内容为 Hello neocode。
- 系统将输入封装成提示并调用 LLM 客户端 Generate。
- LLM 返回 Description 与 Edits，CLI 展示描述与拟执行改动。
- 用户确认执行，Editor.ApplyEdits 将 Edits 应用到本地文件系统，输出执行摘要。
- 支持“预览/计划”模式，降低误修改风险。

## 5. 本地化实现要点
- 默认开启 Mock 模式，确保所有演示在离线环境下可用。
- 未来支持切换到离线推理或离线模型的扩展。
- 对关键操作实现简单回滚：保留备份文件（.bak），必要时扩展历史记录。

## 6. 架构演进路线
- 阶段 1（已实现）：REPL、Mock LLM、文件落地、基础提示模板。
- 阶段 2：接入离线推理/离线模型；增强错误处理与回滚。
- 阶段 3：版本历史、变更审计、私有/公开访问控制。
- 阶段 4：引入 TUI 提升视觉与交互体验。

## 7. 构建、运行与测试
- 构建：
  - go mod tidy
  - go build ./...
- 运行：
  - go run ./cmd/neocode
- 测试：对 llm.go、fs.go、editor.go 编写单元测试，当前以 Mock 为核心测试路径。

## 8. 风险与缓解
- 风险：LLM 给出不可落地的改动，需要有“预览再执行”的保护。
- 缓解：实现简单的文件备份与回滚策略，严格执行预览模式。
- 风险：全本地实现可能影响 UX。
- 缓解：未来引入 TUI 与更丰富的交互策略。

## 9. 里程碑与验收标准
- 验收标准：REPL 能接收自然语言，Mock LLM 给出描述与 Edits，能在本地落地修改文件且输出摘要。
- 验收方法：测试用例覆盖创建/修改/删除文件，验证回滚与备份策略。

## 10. 快速起步清单（可直接落地执行）
- 新增 neocode/docs/IMPLEMENTATION.md 文件（已创建）
- 按照仓库现有结构对照实现，确保各核心模块已导出接口。
- 运行本地演示：使用 Mock 模式，输入自然语言描述，观察 Description 与 Edits。

## 11. 附录
- 术语表、常见问题、扩展点清单。

注：本文档中示例代码仅用于说明接口约束，实际实现请以代码库为准。
