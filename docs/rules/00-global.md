
# 00 — Global（NEVER / DO NOT）

## NEVER

- NEVER 手改生成代码：proto 的 `_pb.go` / `_http.pb.go`、`wire_gen.go`、gorm model和wireProviderSet_gen.go一律从源重新生成。
- NEVER 改 `.proto` 契约却不重新生成（`cd proto && make all`）——`*.pb.go` 被 gitignore，不生成就会用过期产物。
- NEVER 绕过校验或评审：不用 `--no-verify`、不 force-push、不用 TODO 跳过关口。
- NEVER 提交密钥或带真实凭据的配置。
- NEVER 删除或覆盖别人未提交的工作；删除前先调查，别拿破坏性操作当捷径。
- NEVER AI 未经用户明确执行自动执行 `git commit `、`git push` 等git写操作。只有用户明确要求提交或推送时才能执行。

## DO

- DO 一次任务一个变更面；修 bug 不捎带无关重构。
- DO 每次改动后跑 `go build ./... && go test ./internal/... ./cmd/api ./cmd/consumer ./cmd/script`，再宣布完成。
- DO 所有 DB 写操作走 data 层；biz/service 不裸用 db。
- DO 检查失败就停下来修，不带着失败状态往前滚。
