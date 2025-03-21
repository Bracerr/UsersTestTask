# User Management API

REST API для управления пользователями с использованием Clean Architecture и PostgreSQL.

## Запуск проекта

### Локальный запуск

1. Убедитесь, что у вас установлен Go 1.22 или выше
2. Установите и запустите PostgreSQL
3. Создайте файл `.env` в корне проекта со следующими параметрами (или используйте значения по умолчанию):
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=users_db
DB_SSLMODE=disable
```
4. Запустите приложение:
```bash
go run src/cmd/api/main.go
```

Значения по умолчанию для базы данных (если нет .env файла и переменных окружения. Настроить значения можно в `src/internal/config`):
- Host: localhost
- Port: 5432
- User: postgres
- Password: postgres
- Database: users_db
- SSL Mode: disable

### Запуск в Docker

1. Убедитесь, что у вас установлены Docker и Docker Compose
2. Запустите приложение:
```bash
docker-compose up --build
```

При запуске в Docker:
- Приложение будет доступно на порту 8080
- База данных будет доступна на порту 5432
- Все необходимые переменные окружения уже настроены в docker-compose.yml
- Данные базы данных сохраняются в Docker volume

### Остановка контейнеров
```bash
docker-compose down
```

### Просмотр логов
```bash
docker-compose logs -f
```

## API Endpoints

### Создание пользователя
```http
POST /users
Content-Type: application/json

{
    "name": "Ivan",
    "email": "ivan@example.com"
}
```

Ответ в случае успеха (201 Created):
```json
{
    "id": 1,
    "name": "Ivan",
    "email": "ivan@example.com",
    "created_at": "2025-03-21T13:45:30Z",
    "updated_at": "2025-03-21T13:45:30Z"
}
```

### Получение пользователя
```http
GET /users?id=1
```

Ответ в случае успеха (200 OK):
```json
{
    "id": 1,
    "name": "Ivan",
    "email": "ivan@example.com",
    "created_at": "2025-03-21T13:45:30Z",
    "updated_at": "2025-03-21T13:45:30Z"
}
```

### Обновление пользователя
```http
PUT /users
Content-Type: application/json

{
    "id": 1,
    "name": "Ivan Updated",
    "email": "ivan.updated@example.com"
}
```

Ответ в случае успеха (200 OK):
```json
{
    "id": 1,
    "name": "Ivan Updated",
    "email": "ivan.updated@example.com",
    "created_at": "2025-03-21T13:45:30Z",
    "updated_at": "2025-03-21T13:46:15Z"
}
```

Можно обновлять отдельные поля:
```json
{
    "id": 1,
    "email": "new.email@example.com"
}
```

### Возможные ошибки

#### Невалидный email (400 Bad Request):
```json
{
    "error": "invalid email format"
}
```

#### Пользователь не найден (404 Not Found):
```json
{
    "error": "user not found"
}
```

#### Пустые обязательные поля (400 Bad Request):
```json
{
    "error": "invalid input"
}
```

## Конфигурация

Приложение использует следующие переменные окружения (можно задать в `.env` файле):

- `DB_HOST` - хост базы данных (по умолчанию: localhost)
- `DB_PORT` - порт базы данных (по умолчанию: 5432)
- `DB_USER` - пользователь базы данных (по умолчанию: postgres)
- `DB_PASSWORD` - пароль базы данных (по умолчанию: postgres)
- `DB_NAME` - имя базы данных (по умолчанию: users_db)
- `DB_SSLMODE` - режим SSL для подключения к базе данных (по умолчанию: disable)

## Миграции

Миграции базы данных применяются автоматически при запуске приложения. Файлы миграций находятся в директории `src/migrations/`.

При необходимости можно откатить миграции с помощью файлов `.down.sql`.

## Тестирование

В проекте реализованы модульные тесты для всех ключевых компонентов:

### Тесты обработчиков (`src/internal/delivery/handlers/handlers_test.go`)

Тестируют HTTP-обработчики с использованием `httptest` и моков:

- `TestCreateUser`:
  - Проверка успешного создания пользователя
  - Проверка обработки невалидного JSON в теле запроса

- `TestGetUser`:
  - Проверка получения существующего пользователя
  - Проверка обработки невалидного ID пользователя
  - Проверка случая, когда пользователь не найден

- `TestUpdateUser`:
  - Проверка успешного обновления пользователя
  - Проверка обработки невалидного JSON
  - Проверка обновления несуществующего пользователя

### Тесты репозитория (`src/internal/repository/postgres/user_repository_test.go`)

Тестируют слой работы с базой данных с использованием `go-sqlmock`:

- Тесты CRUD операций с пользователями
- Проверка корректности SQL-запросов
- Проверка обработки ошибок базы данных
- Моки для изоляции от реальной базы данных

### Тесты сервиса (`src/internal/service/user_service_test.go`)

Тестируют бизнес-логику:

- Валидация данных пользователя
- Проверка обработки ошибок
- Тестирование взаимодействия с репозиторием через моки

### Запуск тестов

Запуск всех тестов:
```bash
go test ./src/...
```

Запуск тестов конкретного пакета:
```bash
go test ./src/internal/delivery/handlers/  # тесты обработчиков
go test ./src/internal/service/           # тесты сервиса
go test ./src/internal/repository/postgres/ # тесты репозитория
```

Запуск тестов с подробным выводом:
```bash
go test -v ./src/...
``` 