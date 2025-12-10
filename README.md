# SalesTracker

## Описание

Сервис предназначен для **учeта продаж и финансовых транзакций** с поддержкой аналитики и агрегирования данных.

**Основные возможности:**
- CRUD-операции с записями (доходы/расходы)
- Аналитика с расчeтом суммы, среднего, медианы и перцентилей
- Группировка данных по дням, неделям и категориям
- Фильтрация и сортировка записей
- Экспорт данных в CSV
- Веб-интерфейс для управления записями и просмотра аналитики

## HTTP API

- POST /api/items - создание записи
- GET /api/items - получение списка записей с фильтрами
- GET /api/items/{id} - получение записи по ID
- PUT /api/items/{id} - обновление записи
- DELETE /api/items/{id} - удаление записи
- GET /api/analytics - получение аналитики за период
- GET /api/export - экспорт записей в CSV

## Установка и запуск проекта

### 1. Клонирование репозитория

```bash
git clone https://github.com/kstsm/wb-sales-tracker
```

### 2. Настройка переменных окружения

Создайте `.env` файл, скопировав в него значения из `.example.env`:

```bash
cp .example.env .env
```

Отредактируйте `.env` файл, указав необходимые значения:

```bash
# Server
SRV_HOST=localhost
SRV_PORT=8080

# Postgres
POSTGRES_CONTAINER_NAME=sales-tracker-db
POSTGRES_USER=admin
POSTGRES_PASSWORD=admin
POSTGRES_DB=sales_tracker
POSTGRES_PORT=5432
POSTGRES_HOST=localhost
POSTGRES_SSL=disable
POSTGRES_VOLUME_NAME=sales_tracker_data

# Goose
DB_URL=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSL}
MIGRATIONS_DIR=./migrations
```

### 3. Запуск зависимостей (Docker)

```bash
make up
```

Это запустит PostgreSQL в контейнере Docker.

### 4. Миграция базы данных

```bash
make migrate-up
```

### 5. Запуск сервиса

```bash
make run
```
Сервис будет доступен по адресу: http://localhost:8080
___
## Линтер

Проект использует **golangci-lint** для проверки качества кода. 

### Запуск линтера

```bash
make linter
```
___
# API запросы

## POST /api/items - Создание записи
**URL:** `http://localhost:8080/api/items`

**Content-Type:** `application/json`

**Параметры:**

- `type` (обязательно) - тип записи: "income" (доход) или "expense" (расход)
- `amount` (обязательно) - сумма типа int, которая разделяет на рубли и копейки.
- `date` (обязательно) - дата и время в формате RFC3339
- `category` (обязательно) - категория (минимум 3 символа)

**Body:**

```json
{
  "type": "income",
  "amount": 1500000,
  "date": "2025-12-04T19:00:00Z",
  "category": "Оперативная память"
}
```

**Ожидаемый ответ (201 Created):**

```json
{
  "id": "b9ab5b36-444a-47c4-b7b1-7067a4977e67",
  "type": "income",
  "amount": "15000.00",
  "date": "2025-12-04T19:00:00Z",
  "category": "Оперативная память",
  "created_at": "2025-12-09T19:43:11Z",
  "updated_at": "2025-12-09T19:43:11Z"
}
```

### Ошибки:

**Некорректный JSON (400 Bad Request):**

```json
{
  "error": "invalid request body"
}
```

**Ошибки валидации (400 Bad Request):**

```json
{
  "error": "validation for 'Type' failed on the 'item_type' tag"
}
```

```json
{
  "error": "validation for 'Amount' failed on the 'required' tag"
}
```

```json
{
  "error": "validation for 'Date' failed on the 'rfc3339' tag"
}
```

```json
{
  "error": "validation for 'Category' failed on the 'required' tag"
}
```

**Внутренняя ошибка сервера (500 Internal Server Error):**

```json
{
  "error": "internal server error"
}
```
---

## GET /api/items/{id} - Получение записи по ID

**URL:** `http://localhost:8080/api/items/{id}`

**Параметры:**

- `{id}` (обязательно) - UUID записи

**Ожидаемый ответ (200 OK):**

```json
{
  "id": "b9ab5b36-444a-47c4-b7b1-7067a4977e67",
  "type": "income",
  "amount": "15000.00",
  "date": "2025-12-04T19:00:00Z",
  "category": "Оперативная память",
  "created_at": "2025-12-09T19:43:11Z",
  "updated_at": "2025-12-09T19:43:11Z"
}
```

### Ошибки:

**Некорректный ID (400 Bad Request):**

```json
{
  "error": "invalid id"
}
```

или

```json
{
  "error": "id is required"
}
```

**Запись не найдена (404 Not Found):**

```json
{
  "error": "item not found"
}
```

**Внутренняя ошибка сервера (500 Internal Server Error):**

```json
{
  "error": "internal server error"
}
```

---

## GET /api/items - Получение списка записей

**URL:** `http://localhost:8080/api/items`

**Параметры:**

- `from` (опционально) - фильтр по дате начала (RFC3339)
- `to` (опционально) - фильтр по дате окончания (RFC3339)
- `type` (опционально) - фильтр по типу ("income" или "expense")
- `category` (опционально) - фильтр по категории
- `sort_by` (опционально) - сортировка: "date", "amount", "category"
- `sort_order` (опционально) - порядок сортировки: "asc" или "desc"

**Пример запроса:**

```
GET /api/items?from=2025-12-10T03:34:00+06:00Z&type=income&sort_by=date&sort_order=desc
```

**Ожидаемый ответ (200 OK):**

```json
{
  "items": [
    {
      "id": "e633d1de-5838-4424-8a3f-59e9d155c6a7",
      "type": "income",
      "amount": "0.03",
      "date": "2025-12-04T19:00:00Z",
      "category": "Оперативная память",
      "created_at": "2025-12-10T05:15:13Z",
      "updated_at": "2025-12-10T05:15:13Z"
    },
    {
      "id": "a3b6b89f-f129-4341-aab3-72efb40f8f9a",
      "type": "income",
      "amount": "0.02",
      "date": "2025-12-04T19:00:00Z",
      "category": "Оперативная память",
      "created_at": "2025-12-10T05:15:10Z",
      "updated_at": "2025-12-10T05:15:10Z"
    },
    {
      "id": "7097bd26-37c1-4ac8-8d9d-572e329c321a",
      "type": "income",
      "amount": "0.01",
      "date": "2025-12-04T19:00:00Z",
      "category": "Оперативная память",
      "created_at": "2025-12-10T05:15:08Z",
      "updated_at": "2025-12-10T05:15:08Z"
    }
  ],
  "total": 3
}
```

### Ошибки:

**Неизвестный параметр (400 Bad Request):**

```json
{
  "error": "unknown parameter 'invalid_param'"
}
```

**Пустое значение параметра (400 Bad Request):**

```json
{
  "error": "parameter 'type' cannot be empty"
}
```

```json
{
  "error": "parameter 'sort_by' cannot be empty"
}
```

```json
{
  "error": "parameter 'sort_order' cannot be empty"
}
```

**Некорректный формат даты (400 Bad Request):**

```json
{
  "error": "invalid date format"
}
```

```json
{
  "error": "parameter 'from' cannot be after 'to'"
}
```

**Ошибки валидации (400 Bad Request):**

```json
{
  "error": "validation for 'Type' failed on the 'item_type' tag"
}
```

```json
{
  "error": "validation for 'SortBy' failed on the 'sort_by' tag"
}
```

```json
{
  "error": "validation for 'SortOrder' failed on the 'sort_order' tag"
}
```

**Внутренняя ошибка сервера (500 Internal Server Error):**

```json
{
  "error": "internal server error"
}
```

---

## PUT /api/items/{id} - Обновление записи

**URL:** `http://localhost:8080/api/items/{id}`

**Content-Type:** `application/json`

**Параметры:**

- `{id}` (обязательно) - UUID записи
- `type` (опционально) - тип записи: "income" (доход) или "expense" (расход)
- `amount` (опционально) - сумма типа int, которая разделяет на рубли и копейки
- `date` (опционально) - дата и время в формате RFC3339
- `category` (опционально) - категория (минимум 3 символа)

**Body:**

```json
{
  "type": "expense",
  "amount": 2000000,
  "date": "2026-12-20T10:00:00Z",
  "category": "Процессоры"
}
```

**Ожидаемый ответ (200 OK):**

```json
{
  "id": "e9534410-a7e9-4e62-bd5a-0a73ece08bdf",
  "type": "expense",
  "amount": "20000.00",
  "date": "2026-12-20T10:00:00Z",
  "category": "Процессоры",
  "created_at": "2025-12-10T03:34:39Z",
  "updated_at": "2025-12-10T06:38:15Z"
}
```

### Ошибки:

**Некорректный ID (400 Bad Request):**

```json
{
  "error": "invalid id"
}
```

**Некорректный JSON (400 Bad Request):**

```json
{
  "error": "invalid request body"
}
```

**Ошибки валидации (400 Bad Request):**

```json
{
  "error": "validation for 'Type' failed on the 'item_type' tag"
}
```

```json
{
  "error": "validation for 'Amount' failed on the 'required' tag"
}
```

```json
{
  "error": "validation for 'Date' failed on the 'rfc3339' tag"
}
```

```json
{
  "error": "validation for 'Category' failed on the 'required' tag"
}
```

**Запись не найдена (404 Not Found):**

```json
{
  "error": "item not found"
}
```

**Внутренняя ошибка сервера (500 Internal Server Error):**

```json
{
  "error": "internal server error"
}
```

---

## DELETE /api/items/{id} - Удаление записи

**URL:** `http://localhost:8080/api/items/{id}`

**Параметры:**

- `{id}` (обязательно) - UUID записи

**Ожидаемый ответ (200 OK):**

```json
{
  "message": "item deleted successfully"
}
```

### Ошибки:

**Некорректный ID (400 Bad Request):**

```json
{
  "error": "invalid id"
}
```

**Запись не найдена (404 Not Found):**

```json
{
  "error": "item not found"
}
```

**Внутренняя ошибка сервера (500 Internal Server Error):**

```json
{
  "error": "internal server error"
}
```

---

## GET /api/analytics - Получение аналитики

**URL:** `http://localhost:8080/api/analytics`

**Параметры:**

- `from` (обязательно) - дата начала периода (RFC3339)
- `to` (обязательно) - дата окончания периода (RFC3339)
- `group_by` (опционально) - группировка: "day", "week", "category"

**Пример запроса:**

```
GET /api/analytics?from=2025-09-01T03:34:39Z&to=2025-12-18T03:34:39Z&group_by=day
```

**Ожидаемый ответ без группировки (200 OK):**

```json
{
  "from": "2025-12-01",
  "to": "2025-12-31",
  "sum": 15000.50,
  "avg": 5000.17,
  "count": 3,
  "median": 5000.00,
  "percentile_90": 8000.00
}
```

**Ожидаемый ответ с группировкой (200 OK):**

```json
{
  "from": "2025-09-01T03:34:39Z",
  "to": "2025-12-18T03:34:39Z",
  "sum": 180000.09,
  "avg": 18000.009,
  "count": 10,
  "grouped": [
    {
      "group": "2025-09-04",
      "sum": 30000,
      "avg": 30000,
      "count": 1,
      "median": 30000,
      "percentile90": 30000
    },
    {
      "group": "2025-10-04",
      "sum": 30000,
      "avg": 30000,
      "count": 1,
      "median": 30000,
      "percentile90": 30000
    },
    {
      "group": "2025-12-04",
      "sum": 120000.09,
      "avg": 15000.01125,
      "count": 8,
      "median": 15000.015000000001,
      "percentile90": 30000
    }
  ]
}
```

### Ошибки:

**Отсутствуют обязательные параметры (400 Bad Request):**

```json
{
  "error": "parameter 'from' is required"
}
```

```json
{
  "error": "parameter 'to' is required"
}
```

**Некорректный формат даты (400 Bad Request):**

```json
{
  "error": "invalid 'from' date format, expected YYYY-MM-DD: ..."
}
```

```json
{
  "error": "parameter 'from' cannot be after 'to'"
}
```

**Внутренняя ошибка сервера (500 Internal Server Error):**

```json
{
  "error": "internal server error"
}
```

---

## GET /api/export - Экспорт записей в CSV

**URL:** `http://localhost:8080/api/export`

**Параметры:**

- `from` (опционально) - фильтр по дате начала (RFC3339)
- `to` (опционально) - фильтр по дате окончания (RFC3339)
- `type` (опционально) - фильтр по типу ("income" или "expense")
- `category` (опционально) - фильтр по категории
- `sort_by` (опционально) - сортировка: "date", "amount", "category"
- `sort_order` (опционально) - порядок сортировки: "asc" или "desc"

**Пример запроса:**

```
GET /api/export?from=2025-12-01T00:00:00Z&to=2025-12-31T23:59:59Z&type=income&sort_by=amount&sort_order=asc
```

**Ожидаемый ответ (200 OK):**

Файл CSV с заголовками и данными:

```csv
id,type,amount,date,category,created_at,updated_at
7097bd26-37c1-4ac8-8d9d-572e329c321a,income,0.01,2025-12-04T19:00:00Z,Оперативная память,2025-12-10T05:15:08Z,2025-12-10T05:15:08Z
a3b6b89f-f129-4341-aab3-72efb40f8f9a,income,0.02,2025-12-04T19:00:00Z,Оперативная память,2025-12-10T05:15:10Z,2025-12-10T05:15:10Z
e633d1de-5838-4424-8a3f-59e9d155c6a7,income,0.03,2025-12-04T19:00:00Z,Оперативная память,2025-12-10T05:15:13Z,2025-12-10T05:15:13Z
55564dc6-ddcc-46c6-87cc-efff166620ec,income,0.03,2025-12-04T19:00:00Z,Оперативная память,2025-12-10T07:10:51Z,2025-12-10T07:10:51Z
```

**Content-Type:** `text/csv`

**Content-Disposition:** `attachment; filename=items.csv`

### Ошибки:

**Неизвестный параметр (400 Bad Request):**

```json
{
  "error": "unknown parameter 'invalid_param'"
}
```

**Пустое значение параметра (400 Bad Request):**

```json
{
  "error": "parameter 'type' cannot be empty"
}
```

```json
{
  "error": "parameter 'sort_by' cannot be empty"
}
```

```json
{
  "error": "parameter 'sort_order' cannot be empty"
}
```

**Некорректный формат даты (400 Bad Request):**

```json
{
  "error": "invalid date format"
}
```

```json
{
  "error": "parameter 'from' cannot be after 'to'"
}
```

**Ошибки валидации (400 Bad Request):**

```json
{
  "error": "validation for 'Type' failed on the 'item_type' tag"
}
```

```json
{
  "error": "validation for 'SortBy' failed on the 'sort_by' tag"
}
```

```json
{
  "error": "validation for 'SortOrder' failed on the 'sort_order' tag"
}
```

**Внутренняя ошибка сервера (500 Internal Server Error):**

```json
{
  "error": "internal server error"
}
```





