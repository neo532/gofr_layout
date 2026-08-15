# Fitness: proto-sync

## 目标

proto 契约（`*.proto`）变更后必须重新生成 `*.pb.go` 系列产物，禁止手改生成代码造成契约漂移。

## 背景：产物被 gitignore

本仓库的 proto 是**仓库内子模块**（`./proto`，`github.com/neo532/gofr_layout/proto`），生成的 `*.pb.go` 被 `proto/.gitignore` 忽略，`cd proto && make all` 现场生成。因此 git diff 看不到生成产物，proto-sync 改用 **mtime 新鲜度** 判定。

## 判定

- 取 `git merge-base HEAD <origin/main|master|main|HEAD~1>` 到 HEAD 的变更集；
- 若变更集含 `proto/` 下的 `.proto`，则检查该 `.proto` 同目录是否存在 `*.pb.go` 产物，且产物 mtime ≥ `.proto` mtime；
- 违反则 `warn` 并列出未重新生成的文件。

## 为什么是 warn 而不是 block

生成的 `*.pb.go` 不会进入提交，仓库内无法强制"提交时同步"；且 `go build ./...` 本身就是新鲜度的硬门槛（不生成就编译不过）。这里只作提醒信号。

## 手动命令

```sh
cd proto && make all   # 生成 api + client 产物
git status             # 只会看到 .proto / 结构文件的变更
```
