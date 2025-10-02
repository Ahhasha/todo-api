# Todo API

Учебный проект для управления списками (Todo Lists) через REST API.

## Описание
Этот API предоставляет функциональность для создания, получения, обновления и удаления списков задач. Он также поддерживает пагинацию для получения всех списков и проверку состояния сервера через эндпоинт health.

### Структура проекта

.
├── cmd/
│ └── todo-api/
│ └── main.go # Точка входа HTTP-сервиса
├── internal/
│ ├── domain/
│ │ └── list.go # Модель List (entity)
│ ├── service/
│ │ └── list_service.go # Бизнес-логика (интерфейс + реализация)
│ ├── storage/
│ │ └── mem/
│ │ └── list_repo.go # Репозиторий в памяти (CRUD)
│ └── http/
│ ├── router.go # Маршруты и middleware
│ ├── handlers/
│ │ └── lists.go # Обработчики для /api/v1/lists
│ └── middleware/
│ ├── request_id.go # X-Request-Id
│ └── logging.go # Логирование запросов
├── pkg/
│ └── logger/
│ └── logger.go # Обертка над логгером
├── docs/
│ └── openapi.yaml # Спецификация API
├── Makefile # Команды сборки/запуска
├── README.md # Документация проекта
├── go.mod
└── go.sum


#### Требования
- Go 1.22+

##### Установка

1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/Ahhasha/todo-api.git

 Перейдите в директорию проекта:
cd todo-api

Установите зависимости:
 go mod tidy

##### # Запуск

Чтобы запустить сервер, используйте команду:

go run ./cmd/todo-api

Сервер будет доступен по адресу http://localhost:8080

##### ## Эндпоинты

1. POST /api/v1/lists - Создать список

Пример:

curl -sS -X POST http://localhost:8080/api/v1/lists \
  -H "Content-Type: application/json" \
  -d '{"title":"Домашние дела"}'

Ответ:

{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Домашние дела",
  "created_at": "2025-09-20T14:00:00Z"
}

2. GET /api/v1/lists - Получить все списки (с пагинацией)

Пример:

curl -sS "http://localhost:8080/api/v1/lists?limit=10&offset=0"

Ответ:

[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Домашние дела",
    "created_at": "2025-09-20T14:00:00Z"
  },
  ...
]

В ответе будет заголовок X-Total-Count, который содержит общее количество элементов.

3. GET /api/v1/lists/{id} - Получить список по ID

Пример:

curl -sS http://localhost:8080/api/v1/lists/550e8400-e29b-41d4-a716-446655440000

Ответ:

{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Домашние дела",
  "created_at": "2025-09-20T14:00:00Z"
}

4. PATCH /api/v1/lists/{id} - Обновить название списка

Пример:

curl -sS -X PATCH http://localhost:8080/api/v1/lists/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{"title":"Покупки"}'

Ответ:

{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Покупки",
  "created_at": "2025-09-20T14:00:00Z"
}

5. DELETE /api/v1/lists/{id} - Удалить список

Пример:

curl -sS -X DELETE http://localhost:8080/api/v1/lists/550e8400-e29b-41d4-a716-446655440000

Ответ:

(no content)

6. GET /health - Проверка состояния сервера

Пример:

curl -sS http://localhost:8080/health

Ответ:

OK

##### #### Пагинация

Параметры запроса limit и offset позволяют управлять пагинацией:

    limit — количество элементов для возврата (по умолчанию 20).

    offset — начальная позиция (по умолчанию 0).

OpenAPI

Спецификация API доступна в файле openapi.yaml

Вы также можете получить ее по адресу:
curl http://localhost:8080/openapi.yaml
Логирование

Все запросы логируются в файле логов. Используется стандартный логгер, который может быть настроен для вывода в файл.