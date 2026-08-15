# internal/data 层规则

- 写操作用 `r.db.Master(c)`，读操作用 `r.db.Slave(c)`。
- 所有错误用 `errorx.Wrap(err)` 包裹后返回。
- 过滤条件用生成的 model 常量，不写裸 SQL 字符串。
- 字段映射在 data 层完成，model 不泄漏给上层。
