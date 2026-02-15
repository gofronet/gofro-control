# gofro-control

`gofro-control` - HTTP control-plane сервис для управления Xray на удаленных нодах через gRPC.

Сервис:
- поднимает REST API на `:8080`;
- подключает ноды по адресу `host:port` (gRPC);
- хранит список нод в памяти процесса;
- проксирует операции `start/stop/restart`, `get/update config`, `get node info`.

## Стек

- Go (версия из `go.mod`, сейчас `1.26`)
- HTTP: `chi/v5`
- gRPC клиент: `google.golang.org/grpc`
- protobuf-генерация: `buf` + `protoc-gen-go` + `protoc-gen-go-grpc`

## Быстрый старт

### 1. Установить зависимости

```bash
go mod download
```

### 2. Запустить сервис

```bash
go run ./cmd
```

Сервис слушает `http://localhost:8080`.

## API

Базовый префикс: `/v1`

### Ноды

`GET /v1/nodes/`  
Вернет список добавленных нод.

Пример ответа:
```json
[
  {
    "node_name": "node-1",
    "is_xray_running": true
  }
]
```

`POST /v1/nodes/`  
Добавит ноду по gRPC адресу.

Тело запроса:
```json
{
  "node_address": "127.0.0.1:50051"
}
```

Пример ответа:
```json
{
  "node_name": "node-1",
  "is_xray_running": false
}
```

### Конфиг ноды

`GET /v1/nodes/{node_name}/config`  
Вернет текущий конфиг Xray.

Пример ответа:
```json
{
  "node_name": "node-1",
  "current_config": "{...}"
}
```

`PUT /v1/nodes/{node_name}/config`  
Обновит конфиг Xray.

Тело запроса:
```json
{
  "new_config": "{...}"
}
```

Пример ответа:
```json
{
  "status": "ok"
}
```

### Управление процессом Xray

`POST /v1/nodes/{node_name}/start`  
`POST /v1/nodes/{node_name}/stop`  
`POST /v1/nodes/{node_name}/restart`  

Пример ответа:
```json
{
  "status": "ok"
}
```

### Формат ошибок

Ошибки возвращаются в JSON:
```json
{
  "error": "error message"
}
```

## Примеры `curl`

Добавить ноду:
```bash
curl -X POST http://localhost:8080/v1/nodes/ \
  -H "Content-Type: application/json" \
  -d '{"node_address":"127.0.0.1:50051"}'
```

Получить список нод:
```bash
curl http://localhost:8080/v1/nodes/
```

Получить конфиг:
```bash
curl http://localhost:8080/v1/nodes/node-1/config
```

Обновить конфиг:
```bash
curl -X PUT http://localhost:8080/v1/nodes/node-1/config \
  -H "Content-Type: application/json" \
  -d '{"new_config":"{...json...}"}'
```

Рестарт Xray:
```bash
curl -X POST http://localhost:8080/v1/nodes/node-1/restart
```

## Тесты

```bash
go test ./...
```

Важно: текущие тесты в `nodes/nodes_mananger_test.go` используют реальный внешний адрес (`147.45.214.213:50051`) и требуют сетевого доступа до этой ноды.

## Генерация protobuf

Конфиг генерации: `buf.gen.yaml`.  
Входной proto-репозиторий: `ssh://git@gitlab.com/gofronet/specs/node-proto.git`.

Сгенерированные файлы лежат в `gen/go/xray_managment/api/v1`.

## Структура проекта

- `cmd/main.go` - входная точка HTTP сервера.
- `delivery/` - HTTP слой (роутинг, модели запросов/ответов, handlers).
- `nodes/` - менеджер нод и gRPC-клиент к удаленным нодам.
- `gen/go/` - сгенерированный protobuf/gRPC код.

## Ограничения текущей реализации

- Нет персистентности: список нод хранится только в памяти процесса.
- Нет аутентификации/авторизации API.
- Подключение к нодам идет через `insecure` gRPC transport credentials (без TLS).
