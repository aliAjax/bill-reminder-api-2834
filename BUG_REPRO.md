# BUG_REPRO

## Bug 是什么

账单标记为已缴后，`PaidAt` 被设置为 nil；统计汇总随后无条件解引用该指针。

## 如何触发

创建账单，调用标记已缴接口，再调用统计汇总接口。

## 错误信息

`runtime error: invalid memory address or nil pointer dereference`，堆栈指向账单支付时间取值和统计汇总路径。
