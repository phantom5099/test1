# neocode End-to-End Demonstration (Phase 2 Plan)

This document demonstrates a complete end-to-end flow for Phase 2 enhancements: plan/preview/apply flow, robust prompts, and local file edits.

Prerequisites
- Go 1.20+ toolchain
- Local repository with neocode project
- No network required (Mock LLM by default)

Scenario: Create a file via natural language and validate the full cycle
1) User input: 创建一个 demo.txt，内容为 Phase 2 演示
2) System generates a plan via LLM: Description + Edits (e.g., create demo.txt with the provided content)
3) User executes: plan / preview to view, then apply to modify the local file system
4) Validation: demo.txt exists, content matches, and a backup demo.txt.bak is created if pre-existing

Notes
- This demo is designed for CI-friendly demonstration and to show the end-to-end cycle in a single session.
- The LLM is currently a local mock; on Phase 2 we integrate explicit prompts and more robust validation.
