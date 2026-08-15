# PR Review Checklist

本清单与 `AGENTS.md` / `docs/rules/` 绑定，是 Harness 层的人工 review 环节。fitness 门禁通过是合入前提，但不等于人工评审。

## 每次合并前必须确认

- [ ] `make fitness`（`python3 scripts/fitness.py`）全绿，block 项零失败
- [ ] 变更面与链接的 spec（`docs/specs/`）一致：没碰「不解决什么」里的项
- [ ] 冻结 contract 未动：proto message、服务接口、DB schema、wire provider 集、路由路径
- [ ] 生成代码与源同步：proto 变更后 `cd proto && make all` 已重新生成；`wire_gen.go` 已重新生成
- [ ] 无绕过机制：无 `--no-verify`、无跳过测试的提交
- [ ] 无敏感信息（密钥、真实凭据）

## 发现越界时

按高代价优先回退：契约破坏 > 生成代码漂移 > 范围扩散 > 风格。回退后必须重跑 `make fitness`。
