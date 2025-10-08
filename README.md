# Astra API

REST API для хранения и раздачи документов с кэшированием и простым управлением сессиями.

## Возможности
- Регистрация пользователя (через админ-токен)
- Аутентификация и сессии (in-memory)
- Загрузка документов (файл или JSON), список, получение по id, удаление
- Кэширование ответов в памяти для ускорения повторных запросов
- Документация Swagger и статическая раздача JSON спецификации

## Требования
- Go 1.22+
- PostgreSQL 13+

## Конфигурация
Настройки читаются из файла `.enb` в корне проекта:

Пример `.enb`:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=astra
ADMIN_TOKEN=changeme
AUTO_MIGRATE=true
```

Поля:
- DB_* — параметры подключения к БД
- ADMIN_TOKEN — токен, требуемый для регистрации пользователей
- AUTO_MIGRATE — если `true`, миграции применяются при старте

## Запуск
1) База данных (локально через docker-compose):
```bash
docker-compose up -d
```

2) Приложение:
```bash
go run ./cmd
```

При `AUTO_MIGRATE=true` миграции из каталога `./migrations` применятся автоматически.

## Архитектура
- `internal/repository` — доступ к данным (Postgres, sqlx)
  - `interface.go` — интерфейсы репозиториев (`UserRepositoryInterface`, `DocumentRepositoryInterface`)
  - `user.go`, `document.go` — реализации
- `internal/service` — бизнес-логика
  - `interface.go` — интерфейсы сервисов (`AuthServiceInterface`, `DocsServiceInterface`, `SessionServiceInterface`)
  - `auth.go`, `docs.go`, `session.go` — реализации
- `internal/handler` — HTTP-обработчики
- `internal/middleware` — middleware (логирование запросов, проверка авторизации)
- `internal/cache` — простой in-memory кэш с TTL и инвалидацией
- `cmd/main.go` — точка входа, DI, роутинг


## Маршруты API (кратко)

Аутентификация:
- POST `/api/register` — регистрация
  - body: `{ "login": string, "pswd": string, "token": string }` (token = ADMIN_TOKEN)
  - 200: `{ "response": { "login": string } }`
- POST `/api/auth` — логин
  - body: `{ "login": string, "pswd": string }`
  - 200: `{ "response": { "token": string } }`
- DELETE `/api/auth/{token}` — логаут
  - 200: `{ "response": { "<token>": true } }`

Документы (требуется токен):
- POST `/api/docs` — загрузка
  - form-data: `token`, `meta` (json c описанием: name, file, public, mime, grants[]), `file` (опционально), `json` (опционально)
  - 200: `{ "data": { "id": string, "file": string, "json": any|null } }`
- GET|HEAD `/api/docs` — список
  - query: `token`, `login` (опц.), `limit` (опц.)
  - 200: `{ "data": { "docs": Document[] } }`
- GET|HEAD `/api/docs/{id}` — получить по id
  - query: `token`
  - 200: если `file=true` — отдаётся файл, иначе JSON из поля `json_data`
- DELETE `/api/docs/{id}` — удалить по id
  - query: `token`
  - 200: `{ "response": { "<id>": true } }`

Формат ошибки:
```json
{"response":{"<id>":true}}
```

## Миграции
Файлы миграций: `./migrations/*.sql`
Применение на старте управляется `AUTO_MIGRATE`. Для ручного применения можно использовать `goose`.

## Swagger
Статика доступна по:
- `/docs/` — каталог со swagger.json
- `/swagger/` — Swagger UI

## CI
Файл: `.github/workflows/ci.yml`
- checkout, setup-go 1.22.x, cache
- `go build ./...`, `go vet ./...`, `go test ./...`

## Структура проекта (сводно)
```
cmd/
  main.go
internal/
  cache/
  config/
  handler/
  middleware/
  model/
  repository/
    interface.go
    user.go
    document.go
  service/
    interface.go
    auth.go
    docs.go
    session.go
migrations/
docs/
uploads/
```
