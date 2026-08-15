# internal/biz 层规则

- 多步写操作必须包在 `d.txUser(c, func(c context.Context) (err error) {...})` 里；不自己开裸事务。
- 事务提交后需要通知外部的消息（如 kafka）用 producer（`d.pdcUser.Send`）发送，在事务内只做 DB 写。
- 所有错误用 `errorx.Wrap(err)` 包裹后返回。
- biz 只调 repo / connect 接口，不直接接触 DB 或 model。
