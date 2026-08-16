# internal/biz 层规则

- 里面除了wireProviderSet.go外，所有文件内的对外执行方法，必须是一个结构体+New方法导出的。
  - 命名规则：
    - 结构体命名规则：XxxBiz
    - New方法命名：NewXxxBiz
    - 文件名命名：xxx.go
  - 修改完后要执行 make generate
- 多步写操作必须包在 `d.txUser(c, func(c context.Context) (err error) {...})` 里；不自己开裸事务。
- 事务提交后需要通知外部的消息（如 kafka）用 producer（`d.pdcUser.Send`）发送，在事务内只做 DB 写。
- 所有错误用 `errorx.Wrap(err)` 包裹后返回。
- biz 只调 repo / connect 接口，不直接接触 DB 或 model。
- New函数对应结构体参数命名规则：
  - repo变量用r+名字,
  - 事务对象用tx+名字
  - lock对象用lock+名字
  - producer对象用pdc前缀+名字
  - biz对象简写用b+名字
