# BUG_REPRO

## Bug 是什么

创建账单、查询即将到期提醒和统计汇总等操作没有检查或向下传播已取消的 `context`。

## 如何触发

先取消一个 `context.Context`，再调用这些服务方法。

## 错误信息

调用没有返回 `context canceled`，而是继续执行并返回成功结果。
