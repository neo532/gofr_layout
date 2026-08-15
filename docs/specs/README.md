# Specs（规范层）

Spec 把一次模糊意图压缩成可执行、可审查、可验证的变更。它不再是散文——结构由 `scripts/fitness.py` 的 `spec-validate` / `spec-tasks` 机器裁决（block 级），架构参考 OpenSpec：**能力库 + 变更集 delta + 归档**。

## 目录

```
docs/specs/
├── README.md                 # 本说明
├── _template.md              # 变更集 spec.md（delta）模板
├── lib/<capability>.md       # 能力库：当前系统行为的真相（SHALL 需求）
└── changes/
    ├── archive/              # 已归档变更（归档后移入）
    └── <日期>-<slug>/        # 一个进行中的变更集
        ├── .spec.yaml        # status: draft|active|done|archived；capability；created
        ├── proposal.md       # 五问 + 完成条件
        ├── spec.md           # delta：ADDED/MODIFIED/REMOVED Requirements + Gherkin Scenario
        └── tasks.md          # 编号任务，引用 R 需求号
```

## 工作流

1. **propose**：`docs/specs/changes/` 下建 `<日期>-<slug>/`，写 `proposal.md`（五问，见下）。
2. **spec**：写 `spec.md`（delta：`## ADDED/MODIFIED/REMOVED Requirements`，每条 `### Requirement: R#` 用 SHALL，配 ≥1 个 `#### Scenario:` 的 GIVEN/WHEN/THEN/AND）。模板：`_template.md`。
3. **tasks**：写 `tasks.md`，`- [ ] T1（R4）…`，每条需求 R# 至少要有一个任务引用它。
4. **apply**：按 tasks.md 逐条实现，每轮改动后跑 `make loop`，直到 `make fitness` 全绿。
5. **archive**：全部完成后 `make spec-archive change=<日期>-<slug>`——delta 并入 `lib/<capability>.md`，状态置 archived，变更集移入 `changes/archive/`。

## 五问（proposal.md 骨架）

1. 解决什么（一句话可验收）
2. 不解决什么（显式排除，防止顺手扩写）
3. 允许改的 surface（文件 / 包级变更面）
4. 不能动的 contract（冻结项：proto message、服务接口、DB schema、wire provider 集、路由路径）
5. 完成条件（验证命令；可测试断言由 spec.md 的 Scenario 承担）

## 机器门禁（fitness）

- `spec-validate`（block）：校验 lib 与变更集 spec 结构——每条 Requirement 有 R 编号、≥1 个 Scenario、Scenario 含 THEN；`.spec.yaml` 的 status 合法。`draft` 状态宽松放行，`active/done` 全量检查。
- `spec-tasks`（block）：每条需求 R# 在 tasks.md 有对应任务；任务编号 T# 不重复。
- `make loop` 会把当前 active 的变更集名记录进 `.claude/state/loop.json`（`active_specs`）。

## 生命周期

`draft`（写草稿）→ `active`（实现中）→ 归档时置 `archived`（`done` 可作为中间态）。归档前，lib 不包含本次变更的新需求——这就是 delta 的意义。
