# 设计草案

- 结构：命令行入口 -> REPL -> LLM 交互 -> 本地文件系统改动
- 数据：LLMResponse 包含 Description 与 Edits（Op/Path/Content）
- 行为：预览再执行，支持简单备份与回滚
- 未来：引入轻量级元数据存储、版本历史、TUI 界面
