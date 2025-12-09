# SalesTracker

## Описание

Сервис предназначен для **учёта продаж и финансовых транзакций** с поддержкой аналитики и агрегирования данных.

**Основные возможности:**
- CRUD-операции с записями (доходы/расходы)
- Аналитика с расчётом суммы, среднего, медианы и перцентилей
- Группировка данных по дням, неделям и категориям
- Фильтрация и сортировка записей
- Экспорт данных в CSV
- Веб-интерфейс для управления записями и просмотра аналитики

## Линтер

Проект использует **golangci-lint** для проверки качества кода. Все файлы должны соответствовать стандартам Go и правилам линтера.

### Запуск линтера

```bash
make linter
```
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

# API запросы

## POST /api/items - Создание записи

**URL:** `http://localhost:8080/api/items`

**Content-Type:** `application/json`

**Параметры:**

- `type` (обязательно) - тип записи: "income" (доход) или "expense" (расход)
- `amount` (обязательно) - сумма (должна быть больше 0)
- `date` (обязательно) - дата и время в формате RFC3339
- `category` (обязательно) - категория (минимум 3 символа)

**Body:**

```json
{
  "type": "income",
  "amount": 5000.50,
  "date": "2025-12-15T19:00:00Z",
  "category": "Зарплата"
}
```

**Ожидаемый ответ (201 Created):**

```json
{
  "id": "fcdcf25c-fbc1-4941-a3b7-40a24bb71446",
  "type": "income",
  "amount": 5000.50,
  "date": "2025-12-15T19:00:00Z",
  "category": "Зарплата",
  "created_at": "2025-12-02T16:44:42.788089761Z",
  "updated_at": "2025-12-02T16:44:42.788089761Z"
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
  "error": "type is required"
}
```

```json
{
  "error": "amount must be greater than 0"
}
```

```json
{
  "error": "invalid date format"
}
```

```json
{
  "error": "category must be at least 3 characters"
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

**Query параметры (все опциональны):**

- `from` - фильтр по дате начала (RFC3339)
- `to` - фильтр по дате окончания (RFC3339)
- `type` - фильтр по типу ("income" или "expense")
- `category` - фильтр по категории
- `sort_by` - сортировка: "date", "amount", "category"
- `sort_order` - порядок сортировки: "asc" или "desc"

**Пример запроса:**

```
GET /api/items?from=2025-12-01T00:00:00Z&to=2025-12-31T23:59:59Z&type=income&sort_by=date&sort_order=desc
```

**Ожидаемый ответ (200 OK):**

```json
{
  "items": [
    {
      "id": "fcdcf25c-fbc1-4941-a3b7-40a24bb71446",
      "type": "income",
      "amount": 5000.50,
      "date": "2025-12-15T19:00:00Z",
      "category": "Зарплата",
      "created_at": "2025-12-02T16:44:42.788089761Z",
      "updated_at": "2025-12-02T16:44:42.788089761Z"
    }
  ],
  "total": 1
}
```

### Ошибки:

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
  "id": "fcdcf25c-fbc1-4941-a3b7-40a24bb71446",
  "type": "income",
  "amount": 5000.50,
  "date": "2025-12-15T19:00:00Z",
  "category": "Зарплата",
  "created_at": "2025-12-02T16:44:42.788089761Z",
  "updated_at": "2025-12-02T16:44:42.788089761Z"
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
  "error": "Item not found"
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
- Все поля опциональны (обновляются только переданные)

**Body:**

```json
{
  "type": "expense",
  "amount": 1500.00,
  "date": "2025-12-16T10:00:00Z",
  "category": "Продукты"
}
```

**Ожидаемый ответ (200 OK):**

```json
{
  "id": "fcdcf25c-fbc1-4941-a3b7-40a24bb71446",
  "type": "expense",
  "amount": 1500.00,
  "date": "2025-12-16T10:00:00Z",
  "category": "Продукты",
  "created_at": "2025-12-02T16:44:42.788089761Z",
  "updated_at": "2025-12-02T17:30:00.000000000Z"
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
  "error": "amount must be greater than 0"
}
```

**Запись не найдена (404 Not Found):**

```json
{
  "error": "Item not found"
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
null
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
  "error": "Item not found"
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

**Query параметры:**

- `from` (обязательно) - дата начала периода (RFC3339)
- `to` (обязательно) - дата окончания периода (RFC3339)
- `group_by` (опционально) - группировка: "day", "week", "category"

**Пример запроса:**

```
GET /api/analytics?from=2025-12-01T00:00:00Z&to=2025-12-31T23:59:59Z&group_by=day
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
  "from": "2025-12-01",
  "to": "2025-12-31",
  "sum": 15000.50,
  "avg": 5000.17,
  "count": 3,
  "grouped": [
    {
      "group": "2025-12-15",
      "sum": 5000.50,
      "avg": 5000.50,
      "count": 1,
      "median": 5000.50,
      "percentile90": 5000.50
    },
    {
      "group": "2025-12-20",
      "sum": 10000.00,
      "avg": 5000.00,
      "count": 2,
      "median": 5000.00,
      "percentile90": 7500.00
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

**Неподдерживаемая группировка (500 Internal Server Error):**

```json
{
  "error": "internal server error"
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

**Query параметры (все опциональны, аналогично GET /api/items):**

- `from` - фильтр по дате начала (RFC3339)
- `to` - фильтр по дате окончания (RFC3339)
- `type` - фильтр по типу ("income" или "expense")
- `category` - фильтр по категории
- `sort_by` - сортировка: "date", "amount", "category"
- `sort_order` - порядок сортировки: "asc" или "desc"

**Пример запроса:**

```
GET /api/export?from=2025-12-01T00:00:00Z&to=2025-12-31T23:59:59Z&type=income
```

**Ожидаемый ответ (200 OK):**

Файл CSV с заголовками и данными:

```csv
id,type,amount,date,category,created_at,updated_at
fcdcf25c-fbc1-4941-a3b7-40a24bb71446,income,5000.50,2025-12-15T19:00:00Z,Зарплата,2025-12-02T16:44:42.788089761Z,2025-12-02T16:44:42.788089761Z
```

**Content-Type:** `text/csv`

**Content-Disposition:** `attachment; filename=items.csv`

### Ошибки:

**Некорректный формат даты (400 Bad Request):**

```json
{
  "error": "invalid date format"
}
```

**Внутренняя ошибка сервера (500 Internal Server Error):**

```json
{
  "error": "internal server error"
}
```





