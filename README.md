# Astra API — кеширующий веб‑сервер документов (Go + PostgreSQL)

Проект реализует REST API для сохранения и раздачи электронных документов с поддержкой:
- регистрации пользователей (по админ‑токену),
- аутентификации и сессионных токенов,
- загрузки документов (файл/JSON),
- выдачи списка/одного документа,
- удаления документа,
- кеширования GET/HEAD и выборочной инвалидации.

## Стек
- Go 1.24
- PostgreSQL 15
- sqlx, pq
- goose (миграции)
- swag + http-swagger (документация)

## Архитектура
- `cmd/` — запуск приложения
- `internal/config` — конфигурация
- `internal/model` — модели
- `internal/interfaces` — интерфейсы для слоёв (DI)
- `internal/repository` — доступ к БД (с поддержкой транзакций)
- `internal/service` — бизнес‑логика (auth, docs, session)
- `internal/handler` — HTTP‑эндпоинты, единый формат ответа
- `internal/middleware` — мидлвари (auth, logging, CORS)
- `internal/cache` — in‑memory кэш (GET/HEAD)
- `migrations/` — SQL миграции
- `uploads/` — сохранённые файлы

Принципы: слои разделены (SOLID), используются интерфейсы для dependency injection, поддержка транзакций в БД, мидлвари для безопасности и логирования.

## Конфигурация
Переменные окружения в `.enb` (можно переименовать в `.env`):
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=astra
ADMIN_TOKEN=supersecrettoken
# Автоматические миграции при старте (по умолчанию отключены для безопасности)
AUTO_MIGRATE=true
JWT_SECRET=unused
```
Примечание: JWT не используется, по ТЗ авторизация через сессионный токен.

### Безопасность миграций
- `AUTO_MIGRATE=true` — автоматические миграции при старте (для разработки)
- `AUTO_MIGRATE=false` или не указано — миграции отключены (для продакшна)
- В продакшне рекомендуется запускать миграции вручную

## Быстрый старт
1) Поднять PostgreSQL через Docker:
```bash
docker-compose up -d
```

2) Установить зависимости и сгенерировать Swagger (один раз):
```bash
go mod download
export PATH=$PATH:$(go env GOPATH)/bin
swag init --generalInfo cmd/main.go --output docs
```

3) Запуск приложения (миграции применяются автоматически при старте):
```bash
go run ./cmd
```

4) Swagger UI:
- UI: http://localhost:8080/swagger/index.html
- Спецификация: http://localhost:8080/docs/swagger.json

## Миграции
- Миграции **НЕ** применяются автоматически при старте (безопасность в проде).
- Автоматические миграции: установите `AUTO_MIGRATE=true` в переменных окружения.
- Ручной запуск (рекомендуется для продакшна):
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
DB_URL='host=localhost port=5432 user=postgres password=postgres dbname=astra sslmode=disable'
goose -dir ./migrations postgres "$DB_URL" up
```

## Авторизация и сессии
- Регистрация доступна только с `ADMIN_TOKEN`.
- Аутентификация создаёт in‑memory сессию и возвращает `token` (UUID).
- Все методы документов требуют `token`:
  - либо как `?token=...` в query,
  - либо как поле формы `token` (для multipart).
- Логаут: инвалидирует сессию (удаляет токен).

## Формат ответа
Всегда JSON с опциональными блоками:
```json
{
  "error": {"code": 123, "text": "so sad"},
  "response": { ... },
  "data": { ... }
}
```
Коды: 200, 400, 401, 403, 405, 500, 501.

## Эндпоинты (основное)

### Регистрация [POST] /api/register
Вход (JSON):
```json
{"token":"<ADMIN_TOKEN>", "login":"TestUser01", "pswd":"Qwerty123!"}
```
Ответ:
```json
{"response":{"login":"TestUser01"}}
```

### Логин [POST] /api/auth
Вход:
```json
{"login":"TestUser01", "pswd":"Qwerty123!"}
```
Ответ:
```json
{"response":{"token":"<SESSION_TOKEN>"}}
```

### Логаут [DELETE] /api/auth/{token}
Ответ:
```json
{"response":{"<token>":true}}
```

### Загрузка документа [POST] /api/docs
Требует `token` (query или поле формы). Форм‑данные (multipart):
- `meta` — JSON строки с метаданными:
```json
{
  "name": "image.png",
  "file": true,
  "public": false,
  "mime": "image/png",
  "grants": ["login1", "login2"]
}
```
- `file` — бинарный файл (если `file=true`)
- `json` — произвольный JSON контент (опционально)

Ответ (пример):
```json
{"data":{"id":"<doc_id>","file":"image.png","json":null}}
```

Пример curl:
```bash
curl -X POST "http://localhost:8080/api/docs?token=<TOKEN>" \
  -F 'meta={"name":"image.png","file":true,"public":false,"mime":"image/png","grants":["login1","login2"]}' \
  -F 'file=@/path/image.png;type=image/png'
```
Для JSON‑документа:
```bash
curl -X POST "http://localhost:8080/api/docs?token=<TOKEN>" \
  -F 'meta={"name":"data.json","file":false,"public":false,"mime":"application/json"}' \
  -F 'json={"key":"value","n":123}'
```

### Список документов [GET|HEAD] /api/docs
Параметры:
- `token` — обязателен
- `login` — опционально (если не указан — список своих)
- `limit` — опционально

Ответ:
```json
{"data":{"docs":[{ "id":"...","name":"...","mime":"...","file":true,"public":false,"created":"...","grants":["login1"]}]}}
```

### Документ по id [GET|HEAD] /api/docs/{id}
- Если `file=true` — отдаётся файл с `Content-Type: mime`.
- Если JSON — отдаётся `data: { ... }`.

### Удаление [DELETE] /api/docs/{id}
Ответ:
```json
{"response":{"<id>":true}}
```

## Кэширование
- GET/HEAD /api/docs и /api/docs/{id} читают из in‑memory кэша.
- POST/DELETE инвалидируют кэш (выборочно и/или целиком для упрощения консистентности).
- TTL кэша: 5 минут (настраивается в `cmd/main.go`).

## Хранение файлов
- Файлы сохраняются в `./uploads`.
- Имя файла берётся из загруженного файла (в проде стоит добавить генерацию уникального имени).

## БД (схема)
- `users(id uuid pk, login unique, password, created_at)`
- `documents(id uuid pk, name, mime, file bool, public bool, owner uuid fk users, created_at, grants text[], json_data jsonb)`

## Как тестировать (пример последовательности)
1) Регистрация:
```bash
curl -X POST http://localhost:8080/api/register \
  -H 'Content-Type: application/json' \
  -d '{"token":"supersecrettoken","login":"TestUser01","pswd":"Qwerty123!"}'
```
2) Логин → получить `token`:
```bash
curl -X POST http://localhost:8080/api/auth \
  -H 'Content-Type: application/json' \
  -d '{"login":"TestUser01","pswd":"Qwerty123!"}'
```
3) Загрузка файла:
```bash
curl -X POST "http://localhost:8080/api/docs?token=<TOKEN>" \
  -F 'meta={"name":"image.png","file":true,"public":false,"mime":"image/png"}' \
  -F 'file=@./image.png;type=image/png'
```
4) Список:
```bash
curl "http://localhost:8080/api/docs?token=<TOKEN>&limit=10"
```
5) Получение по id:
```bash
curl "http://localhost:8080/api/docs/<DOC_ID>?token=<TOKEN>"
```
6) Удаление:
```bash
curl -X DELETE "http://localhost:8080/api/docs/<DOC_ID>?token=<TOKEN>"
```
7) Логаут:
```bash
curl -X DELETE "http://localhost:8080/api/auth/<TOKEN>"
```

## Нюансы и оговорки
- **Безопасность**: Миграции теперь не выполняются автоматически в продакшне.
- **Архитектура**: Используются интерфейсы для dependency injection и лучшей тестируемости.
- **Транзакции**: Добавлена поддержка транзакций в репозиториях для атомарных операций.
- **Мидлвари**: Реализованы мидлвари для аутентификации, логирования и CORS.
- Сессии In‑Memory: для продакшна имеет смысл вынести в Redis.
- Файлы: нет дедупликации, ограничений размера, проверки контента — добавляются при необходимости.
- Кэш: простой in‑memory, без распределённости и метрик.

## Улучшения безопасности и архитектуры

### Решенные проблемы:
1. ✅ **Безусловное выполнение миграций** — добавлен флаг `AUTO_MIGRATE`
2. ✅ **Отсутствие интерфейсов** — реализованы интерфейсы для всех слоёв
3. ✅ **Отсутствие транзакций** — добавлена поддержка транзакций в репозиториях
4. ✅ **Отсутствие мидлварей** — полная система мидлварей (auth, logging, CORS)




PR/Issues: предлагайте улучшения по безопасности, правам, транзакциям