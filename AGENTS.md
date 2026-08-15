# AGENTS.md（全 agent 通用规则）

> 本文件是所有 AI 编码 agent（Claude Code、Codex、Copilot 等）在本仓库工作时必须遵守的规则，每次会话开始时读取，改动后立即生效。
> 形式：硬性规则用祈使句，附一行"为什么"，方便判断边界。
> `CLAUDE.md` 只是本文件的指针；后续新增规则按"重复错误才追加"下沉到 `docs/rules/` / `docs/specs/`。

## 一句话背景

本项目（gofr 系）的目标：**一套 proto 定义 API，HTTP / gRPC / rpcx / WebSocket 等所有协议都从同一份 `google.api.http` 注解派生**，路径映射到对应方法。`gofr_layout` 是这套体系的**项目布局模板**：拿来即用的目录结构、wire 注入、配置热加载与示例服务（user 域 CRUD）。

## 硬性规范（必须遵守）

- 代码注释一律用英文；与用户对话用中文。
- 不新建与现有模式相悖的抽象；优先改现有代码。
- 修改代码后必须运行所在模块的 `go build` 和 `go test`。

## 风格约定

- 一个文件尽量收拢成一个结构体，通用工具名（如 `headerSet`、`shellQuote`）只做方法，不做包级函数，避免污染包作用域。
- 不写多余的注释；只在 WHY 非显而易见时才写，且不超过一行。
- 不要过度防御：只为系统边界（用户输入、外部 API）做校验。

## 模块布局

- 本仓库 `gofr_layout` 是 gofr 项目的布局模板（示例服务：user 域 CRUD）。
- 引用的 `../gofr`、`../gokit` 是独立模块（独立 go.mod，暂由 gofr_layout/go.mod 的 `replace` 指向本地）。改动它们时需分别进入各自目录验证。
- `./proto` 是仓库内 proto 子模块（`github.com/neo532/gofr_layout/proto`），生成的 `*.pb.go` 被 gitignore，需 `cd proto && make all` 现场生成。
- 开发用根目录 `go.work` 把本仓库与 `../gofr`、`../gokit`、`./proto` 组成 workspace；生成/CI 用 `GOWORK=off` 走 go.mod。

## 验证流程（Golden Path）

- 每次改动后运行：
  ```sh
  go build ./... && go test ./internal/... ./cmd/api ./cmd/consumer ./cmd/script
  ```
- 集成测试 `go test ./cmd/test` 需要本地 DB，可选。
- 每轮编辑后 `make loop`，直到 `make fitness` 全绿再宣布完成。

## 数据与生成代码

- 多步 DB 写操作走 `txUser` 事务；写用 `r.db.Master(c)`，读用 `r.db.Slave(c)`。
- 事务提交后需要通知外部的消息（如 kafka）用 producer（`pdcUser.Send`）发送。
- 生成代码（proto `_pb.go`、`wire_gen.go`）只由生成器产出，禁止手改。

## AI Coding 治理（Rule / Spec / Loop / Harness）

- 分层规则见 `docs/rules/`；本文件放仓库级硬性规范，具体规则按"重复错误才追加"下沉到 `docs/rules/`。
- 单次变更先建 Spec 变更集（`docs/specs/changes/<日期>-<slug>/`，结构见 `docs/specs/README.md`）。
- 完成条件由 `scripts/fitness.py` 机器裁决（`make fitness`，含 spec-validate / spec-tasks）；一轮编辑闭环 `make loop`。
- 变更完成后 `make spec-archive change=<slug>` 归档：delta 并入 `docs/specs/lib/`。
- PR 门禁见 `.github/`；评审清单与 `AGENTS.md` / `docs/rules/` 绑定。

## 开始前必读

- `AGENTS_README.md` — 治理骨架的使用指南（命令速查 + 三种用法）。
- `docs/rules/00-global.md` — 仓库级 NEVER / DO NOT 清单。
- `docs/rules/README.md` — 规则如何分层、何时可以新增规则。
- `docs/specs/README.md` — Spec 变更集工作流（能力库 + delta + 归档）。
