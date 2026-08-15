# Fitness Functions（Harness 层）

Harness 回答「这次变更为什么足够可信」。fitness function 是仓库级裁决的完成条件机制：每条 fitness 定义一个外部信号，gate 等级决定它是否拦截 PR。

## 裁决机制

- 执行器：`scripts/fitness.py`（`make fitness`）
- 门禁等级：
  - `block`：失败即整体失败（exit 1），PR 不允许合入
  - `warn`：只报告、不拦截（如需要本地 DB 的集成测试）
  - `skip`：条件不满足（如 proto 目录不存在）
- CI：`.github/workflows/ci.yml` 在 PR 上跑同一执行器

## 当前 fitness 清单

| name | 命令 / 判定 | gate | 说明 |
|---|---|---|---|
| build | `go build ./...`（go.work） | block | 编译 |
| vet | `GOWORK=off go vet ./...` | block | 静态检查 |
| test-unit | `go test ./internal/... ./cmd/api ./cmd/consumer ./cmd/script` | block | 单测（不依赖外部服务） |
| test-integration | `go test ./cmd/test` | warn | 集成测试，需本地 DB |
| fmt | `gofmt -l` 本次变更的 .go 文件 | warn | 只查变更面，避免被预存在的格式问题绊住 |
| wire-sync | diff：wire.go / wireProviderSet.go 变更 ⇒ 必须带 wire_gen.go | block | 守住 wire 生成契约 |
| proto-sync | 变更的 .proto ⇒ 磁盘上对应 `*.pb.go` 产物必须比 .proto 新 | warn | 守住 proto 契约（产物 gitignored，用 mtime 判定新鲜度） |
| spec-validate | 结构：每条 Requirement 有 R 编号 + ≥1 Scenario，Scenario 含 THEN；status 合法 | block | 守住 spec 结构（draft 宽松） |
| spec-tasks | 覆盖：每条需求 R# 有任务，任务 T# 不重复 | block | 守住 spec ⇔ tasks 一致 |

> 注意：本仓库 proto 的 `*.pb.go` 是生成产物并被 gitignore，proto-sync 用 **mtime 新鲜度** 判定（`.proto` 变更后必须 `cd proto && make all` 重新生成），而不是 git diff。spec 类检查是内容型，随时生效。

## 如何加一条 fitness

1. 在 `scripts/fitness.py` 加一个 `check_xxx()`，返回 `(status, label, output)`，并挂进 `CHECKS`。
2. 按高代价优先决定 gate：破坏契约用 `block`，仅提示用 `warn`。
3. 在 `docs/rules/` 对应规则里同步一句话，保证「规则 ⇔ 机器信号」一致。
