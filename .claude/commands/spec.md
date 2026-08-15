---
description: 创建/推进一个 Spec 变更集（SDD 流程）
---

按 `docs/specs/README.md` 的流程工作：

1. 新变更：在 `docs/specs/changes/` 下建 `<日期>-<slug>/`，写 `.spec.yaml`（status: draft/active、capability）、`proposal.md`（五问）、`spec.md`（delta：`## ADDED/MODIFIED/REMOVED Requirements` + `### Requirement: R#` + Gherkin `#### Scenario:`）、`tasks.md`（`- [ ] T1（R#）…`）。
2. 逐条按 tasks.md 实现并勾选，每轮改动后 `make loop`，直到 `make fitness` 全绿（含 spec-validate / spec-tasks）。
3. 完成后 `make spec-archive change=<日期>-<slug>`：delta 并入 `docs/specs/lib/`，状态置 archived，变更集移入 `changes/archive/`。
