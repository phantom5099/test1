# neocode 快速上手（Phase 1：本地 CLI MVP）

本快速指南帮助你在本地快速验证 Phase 1 的全面稳定实现。

前提
- 已安装 Go 1.20+，以及 Git、基本的终端环境。
- 代码已在仓库中，路径以 neocode/ 为根。请在 neocode 目录下执行后续命令。

步骤 1：编译与运行
- 进入工作目录：`cd neocode`
- 安装依赖并测试：`go mod tidy && go test ./...`
- 启动 CLI：`go run ./cmd/neocode`

步骤 2：REPL 使用示例

注意
- 由于当前实现采用严格的本地 Mock LLM，所有操作均在本地文件系统中体现，无需网络。
- 如需演示不同场景，可修改 Hello 内容、路径等来验证不同的 Edits。

阶段性结果与可验证性
- Phase 1 验收点：REPL 启动、自然语言输入、Mock LLM 返回描述+Edits、本地落地、简单备份回滚。
- 验收步骤同阶段性验收清单中的描述。

扩展提示
- 下一阶段将引入显式的 plan/preview/apply 流程、增强错误处理、增加历史/回滚等。
