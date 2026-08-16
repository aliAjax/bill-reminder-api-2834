# BUG_REPRO

## Bug 是什么

新增账单接口在 `due_date` 为空时，错误分类逻辑不再识别 `ErrInvalidInput`，导致本应返回 `400 validation_error` 的请求被错误返回为 `500 internal_error`。

## 如何触发

向 `POST /api/v1/bills` 提交缺少 `due_date` 的 JSON：

```json
{"name":"8月水费","type":"water","amount":42.5,"due_date":""}
```

## 错误信息

接口返回状态码 `500`，响应体为：

```json
{"success":false,"error":{"code":"internal_error","message":"internal server error"}}
```
