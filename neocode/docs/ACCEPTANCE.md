# Acceptance Criteria for neocode MVP (Phase 1 & Phase 2)

- Phase 1 (Local CLI MVP)
  - Build and run: go mod tidy, go test ./..., go run ./cmd/neocode works locally.
  - REPL works: user can input natural language prompts; system returns a textual Description and Edits; user can preview and apply changes locally.
  - All changes are performed locally (Mock LLM by default) with no network calls unless explicitly enabled.
  - File system operations (create, update, delete) work atomically with basic backups (.bak) and idempotent behavior.
  - Tests cover core filesystem operations and mock LLM behavior; documentation is present.

- Phase 2 (Enhancements)
  - REPL supports explicit plan/preview/apply flow with deterministic results and robust error handling.
  - History/logging: a local history mechanism records applied changes for auditing and rollback.
  - End-to-end demo scripts and quickstart are provided to demonstrate end-to-end flows.
  - CI workflow ensures on push/PR that tests pass and builds succeed.

- Acceptance process
  - Self-check: ensure all steps in QuickStart.md can be followed to achieve the same end-to-end outcome.
  - Review: verify that Phase 2 patch sets are consistent, with clear tests and documentation.
