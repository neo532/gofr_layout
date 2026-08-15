# <变更标题> — Spec Delta

> 新建变更：复制本文件为 `docs/specs/changes/<日期>-<slug>/spec.md`，并在同目录补 `.spec.yaml`、`proposal.md`、`tasks.md`。流程见 `README.md`。
> 每条需求必须有 R 编号与 ≥1 个 Gherkin Scenario（`make fitness` 的 spec-validate / spec-tasks 强制）。

## ADDED Requirements

### Requirement: R<N> — <需求名>

<行为 SHALL 描述>。

#### Scenario: <场景名>
- **GIVEN** <前提>
- **WHEN** <动作>
- **THEN** <可观测结果>
- **AND** <附加断言>（可选）

## MODIFIED Requirements

<!-- 必须包含完整的更新后需求文本；归档脚本按 R 编号整体替换 lib 中的旧块 -->

### Requirement: R<N> — <需求名>

## REMOVED Requirements

<!-- 只需 R 编号与名称；归档脚本按 R 编号从 lib 删除 -->

### Requirement: R<N> — <需求名>
