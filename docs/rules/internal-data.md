# internal/data 层规则

- 里面除了wireProviderSet.go外，所有文件内的对外执行方法，必须是一个结构体+New方法导出的
  - 结构体命名规则：XxxBiz
  - New方法命名：NewXxxBiz
  - 文件名命名：xxx.go
- 写操作用 `r.db.Master(c)`，读操作用 `r.db.Slave(c)`。
- 所有错误用 `errorx.Wrap(err)` 包裹后返回。
- 过滤条件用 `fmt.Sprintf("%s=?", model.XxxColumn.ID)` + 生成的 column 常量，不写裸 SQL 字符串。
- 字段映射在 data 层完成，model 不泄漏给上层。
