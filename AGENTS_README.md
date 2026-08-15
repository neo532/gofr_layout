# AGENTS_README.md — AI Coding 治理使用指南

> 本文档说明本仓库内置的 **Rule / Spec / Loop / Harness** 治理骨架怎么用。上手只需要记三件事：**建 spec → 改代码跑 `make loop` → 完成归档**。

## 这套东西是什么

四层架构（参考 Phodal 的《从 Rule、Spec 到 Harness》，目录结构借鉴 OpenSpec / Spec-Kit）：

| 层 | 职责 | 落到哪 |
|---|---|---|
| **Rule** | 约束：让 agent 先不越界 | `AGENTS.md`（全 agent 通用入口）+ `docs/rules/`（分层 NEVER / DO NOT） |
| **Spec** | 规范：一次变更范围钉住 | `docs/specs/`（能力库 `lib/` + 变更集 `changes/`，delta 归档） |
| **Loop** | 闭环：每轮改完验证 + 外置状态 | `make loop` → `scripts/fitness.py` + `.claude/state/loop.json` |
| **Harness** | 治理：机器裁决 + PR 门禁 | `scripts/fitness.py` + `.github/workflows/ci.yml` + 评审清单 |

## 命令速查

| 命令 | 干什么 | 什么时候用 |
|---|---|---|
| `/spec` | 引导 agent 走完整个变更集流程 | 要开始一个新功能时，直接喊它 |
| `make fitness` | 一键门禁：build / vet / 单测 / 格式 / wire / proto / spec 结构 / 任务覆盖，共 9 项 | 想快速自检，合入前必跑 |
| `make loop` | 一轮闭环：跑 fitness + 把当前 active spec 记进 `.claude/state/loop.json` | 每改完一轮代码，代替手动跑 build/test |
| `make spec-archive change=<日期>-<slug>` | 归档：delta 并入 `docs/specs/lib/`，变更集移入 `changes/archive/` | 功能全部完成、fitness 全绿之后 |

## 三种典型用法

### 1. 大改动 / 新功能（走完整流程）

```
/spec 我要加一个"按用户名搜索用户"的接口
```

agent 会帮你建 `docs/specs/changes/<日期>-<slug>/`，写好：

```
docs/specs/changes/<日期>-<slug>/
├── .spec.yaml        # status: active；capability；created
├── proposal.md       # 五问（解决什么/不解决什么/允许改的 surface/不能动的 contract/完成条件）
├── spec.md           # delta：ADDED/MODIFIED/REMOVED Requirements + Gherkin Scenario
└── tasks.md          # 编号任务，引用 R 需求号
```

然后它按 tasks.md 逐条实现、每轮跑 `make loop`，全绿后归档。你只需要：

- 开工时把需求说清楚，看一遍 spec 的五问有没有跑偏；
- 收尾时确认 fitness 绿，跑一次 `make spec-archive change=<slug>`。

### 2. 小改动 / 修 bug（可以轻量）

直接改代码 + `make loop`。spec 流程可以省，但 `make loop` 别省——它是你唯一不需要盯的"我是不是真的改完了"的信号。

### 3. 自己先想清楚，再丢给 AI

```
docs/specs/changes/2026-08-13-xxx/
├── proposal.md   # 你写：五问
├── spec.md       # 你写：R1 SHALL 能按用户名搜到用户；Scenario 给 GIVEN/WHEN/THEN
└── tasks.md      # 你写：T1（R1）repo 加 GetByUsername；T2（R1）service 接上
```

然后把 `make fitness` 和 tasks.md 一起甩给 agent："按 tasks.md 做，做完 `make loop` 全绿再回来。" 这时 spec 里的 Scenario 就是验收标准——agent 说"做完了"之前，fitness 会先替你把关。

## 兜底机制（你不需要盯的）

- 我写错 spec 结构（缺 Scenario、状态非法）→ `spec-validate` 红灯，`make loop` 直接停。
- 我改了 proto 或 wire 没重新生成 → `wire-sync` / `proto-sync` 红灯。
- 每轮状态在 `.claude/state/loop.json`（gitignored，不污染仓库），跨会话也能接着收敛。

## 其他 agent（Codex / Copilot）

规则都在 `AGENTS.md`（全 agent 通用），它们也会读。它们没有 `/spec` 命令，就把 `docs/specs/README.md` 指给它，它会按同样的流程走；`make fitness` / `make loop` / `make spec-archive` 是仓库级 make 目标，谁都调得动。

## 一句话

**`make loop` 全绿 = 收敛，归档 = 收尾。** 剩下的交给 agent 自己循环。
