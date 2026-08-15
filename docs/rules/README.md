# Rules（约束层）

规则是 AI Coding 的**入口**，不是护栏。真正的护栏是外部信号链（`go build` / `go test` / `go vet` / `gofmt` / fitness 门禁），它们负责把越界产出挡在合入之前。

## 分层

- 根目录 `AGENTS.md` — 全 agent 通用入口，放仓库级硬性规范；`CLAUDE.md` 是它的指针。
- `docs/rules/00-global.md` — 仓库级 NEVER / DO NOT 清单。
- `docs/rules/internal-*.md` — 包级规则，离它约束的代码更近。

## 何时新增规则

同一类错误**重复出现**才回写进规则。发生一次，修在当次任务里；再次发生，写进对应规则文件，并同步到 `scripts/fitness.py` 里的机器信号。无法被外部信号裁决的规则，最多只是建议。

## 排序原则

- NEVER / DO NOT 写在建议之前。
- 先约束高代价错误（契约漂移、生成代码失同步、绕过评审），再讨论风格。
