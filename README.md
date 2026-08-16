# 家庭缴费提醒 API

一个从零实现的 Go 后端 API，用于记录水费、电费、燃气费、物业费等周期性账单。数据以本地 JSON 文件持久化，接口采用统一的 JSON 成功/错误结构。

## 目录结构

```text
.
├── cmd/server/main.go
├── internal
│   ├── bill
│   │   ├── handler.go
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── reminder
│   │   ├── handler.go
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── statistics
│   │   ├── handler.go
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── config
│   │   └── config.go
│   ├── httpapi
│   │   └── response.go
│   └── store
│       └── json_store.go
├── Dockerfile
├── go.mod
└── README.md
```

## 启动方式

### Docker

```bash
docker build -t bill-reminder-api .
docker run -d --name bill-reminder-api -p 18089:8080 bill-reminder-api
```

服务启动后监听 `http://localhost:18089`。验证完成后清理：

```bash
docker rm -f bill-reminder-api
```

### 本机 Go

如果本机安装了 Go 1.23+：

```bash
go run ./cmd/server
```

## 环境变量

| 变量名 | 默认值 | 说明 |
| --- | --- | --- |
| `PORT` | `8080` | HTTP 服务监听端口 |
| `DATA_FILE` | `data/bills.json` | 本地 JSON 数据文件路径 |

## API

所有响应均为 JSON。成功响应结构为：

```json
{"success":true,"data":{}}
```

失败响应结构为：

```json
{"success":false,"error":{"code":"validation_error","message":"..."}}
```

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| `GET` | `/healthz` | 健康检查 |
| `POST` | `/api/v1/bills` | 新增账单 |
| `GET` | `/api/v1/bills` | 查看全部账单，可用 `?status=unpaid` 或 `?status=paid` 筛选 |
| `PATCH` | `/api/v1/bills/{id}/pay` | 标记已缴 |
| `PATCH` | `/api/v1/bills/{id}/due-date` | 修改到期日 |
| `GET` | `/api/v1/reminders/upcoming?days=7` | 返回即将到期账单 |
| `GET` | `/api/v1/statistics/summary` | 返回账单统计汇总 |

## curl 示例

新增账单：

```bash
curl -sS -X POST http://localhost:18089/api/v1/bills \
  -H 'Content-Type: application/json' \
  -d '{"name":"8月水费","type":"water","amount":42.50,"due_date":"2026-08-20"}'
```

查看全部账单：

```bash
curl -sS http://localhost:18089/api/v1/bills
```

按状态筛选：

```bash
curl -sS 'http://localhost:18089/api/v1/bills?status=unpaid'
```

标记已缴：

```bash
curl -sS -X PATCH http://localhost:18089/api/v1/bills/{id}/pay
```

修改到期日：

```bash
curl -sS -X PATCH http://localhost:18089/api/v1/bills/{id}/due-date \
  -H 'Content-Type: application/json' \
  -d '{"due_date":"2026-09-01"}'
```

返回即将到期账单：

```bash
curl -sS 'http://localhost:18089/api/v1/reminders/upcoming?days=7'
```

查看统计汇总：

```bash
curl -sS http://localhost:18089/api/v1/statistics/summary
```
