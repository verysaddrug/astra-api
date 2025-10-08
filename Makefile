SHELL := /bin/bash

# Переменные
GO ?= go
GOBIN ?= $(shell $(GO) env GOPATH)/bin
GOOSE ?= $(GOBIN)/goose
MOCKGEN ?= $(GOBIN)/mockgen

# Параметры БД (можно переопределить через env)
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= postgres
DB_PASSWORD ?= postgres
DB_NAME ?= astra
DB_SSLMODE ?= disable

# Строка подключения для goose
GOOSE_DRIVER ?= postgres
GOOSE_DBSTRING ?= host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=$(DB_SSLMODE)

.PHONY: help tidy deps build test cover cover-html lint mockgen-install mocks-gen goose-install migrate-up migrate-down migrate-status migrate-redo clean

# Показать справку
help:
	@echo "Доступные команды:"
	@echo "  help           - показать эту справку"
	@echo "  tidy           - очистить go.mod"
	@echo "  deps           - скачать зависимости"
	@echo "  build          - собрать проект"
	@echo "  test           - запустить тесты"
	@echo "  cover          - показать покрытие тестами"
	@echo "  cover-html     - создать HTML отчёт покрытия"
	@echo "  lint           - запустить линтер"
	@echo "  mockgen-install- установить mockgen (gomock)"
	@echo "  mocks-gen      - сгенерировать gomock моки"
	@echo "  goose-install  - установить goose"
	@echo "  migrate-up     - применить миграции"
	@echo "  migrate-down   - откатить миграции"
	@echo "  migrate-status - показать статус миграций"
	@echo "  migrate-redo   - повторить последнюю миграцию"
	@echo "  clean          - очистить временные файлы"

# Очистить go.mod
tidy:
	$(GO) mod tidy

# Скачать зависимости
deps:
	$(GO) mod download

# Собрать проект
build:
	$(GO) build ./...

# Запустить тесты
test:
	$(GO) test ./... -run . -v

# Показать покрытие тестами
cover:
	$(GO) test ./... -coverprofile=coverage.out -covermode=atomic
	@echo "Покрытие тестами:"
	@$(GO) tool cover -func=coverage.out | tail -n 1

# Создать HTML отчёт покрытия
cover-html: cover
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "HTML отчёт покрытия создан: coverage.html"

# Запустить линтер
lint:
	golangci-lint run --timeout=5m

# Установить mockgen для генерации gomock моков
mockgen-install:
	$(GO) install go.uber.org/mock/mockgen@v0.5.0

# Сгенерировать gomock моки для всех интерфейсов
mocks-gen: mockgen-install
	@echo "Генерация gomock моков..."
	mkdir -p internal/mocks/gomock
	$(MOCKGEN) -source=internal/service/interface.go -destination=internal/mocks/gomock/service_mocks.go -package=mocks
	$(MOCKGEN) -source=internal/repository/interface.go -destination=internal/mocks/gomock/repository_mocks.go -package=mocks
	@echo "Gomock моки сгенерированы в internal/mocks/gomock/"

# Установить goose для миграций
goose-install:
	$(GO) install github.com/pressly/goose/v3/cmd/goose@v3.25.0

# Применить миграции
migrate-up: goose-install
	@echo "Применение миграций..."
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$(GOOSE_DBSTRING)" $(GOOSE) -dir ./migrations up

# Откатить миграции
migrate-down: goose-install
	@echo "Откат миграций..."
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$(GOOSE_DBSTRING)" $(GOOSE) -dir ./migrations down

# Показать статус миграций
migrate-status: goose-install
	@echo "Статус миграций:"
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$(GOOSE_DBSTRING)" $(GOOSE) -dir ./migrations status

# Повторить последнюю миграцию
migrate-redo: goose-install
	@echo "Повтор последней миграции..."
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$(GOOSE_DBSTRING)" $(GOOSE) -dir ./migrations redo

# Очистить временные файлы
clean:
	rm -f coverage.out coverage.html
	@echo "Временные файлы очищены"