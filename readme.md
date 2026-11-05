# golang-ToDoApp

REST API для управления задачами (todo-лист) на Go. Поддерживает создание, чтение, обновление и удаление задач с хранением в PostgreSQL. Документация API — через Swagger.

## Стек

- **Язык:** Go
- **База данных:** PostgreSQL
- **Инфраструктура:** Docker, docker-compose
- **Документация:** Swagger (OpenAPI)
- **Миграции:** SQL-миграции в `migrations/`
- **Автоматизация:** Makefile

## Структура проекта

```
.
├── cmd/todoapp/        # точка входа приложения
├── internal/           # внутренняя логика (handlers, service, repository)
├── migrations/         # SQL-миграции схемы БД
├── docs/               # Swagger-документация
├── docker-compose.yaml # запуск приложения и PostgreSQL
├── Makefile            # команды для сборки и запуска
└── go.mod
```


## Запуск

### Через Docker (рекомендуется)

```bash
docker-compose up --build
